package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	ruleenginedomain "github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/Afomiat/Digital-IMCI/ruleengine/engine"
	"github.com/google/uuid"
)

type RuleEngineUsecase struct {
	ruleEngine                    *engine.RuleEngine
	assessmentRepo                domain.AssessmentRepository
	medicalProfessionalAnswerRepo domain.MedicalProfessionalAnswerRepository
	clinicalFindingsRepo          domain.ClinicalFindingsRepository
	classificationRepo            domain.ClassificationRepository
	counselingRepo                domain.CounselingRepository
	treatmentPlanRepo             domain.TreatmentPlanRepository
	contextTimeout                time.Duration
}

func NewRuleEngineUsecase(
	ruleEngine *engine.RuleEngine,
	assessmentRepo domain.AssessmentRepository,
	medicalProfessionalAnswerRepo domain.MedicalProfessionalAnswerRepository,
	clinicalFindingsRepo domain.ClinicalFindingsRepository,
	classificationRepo domain.ClassificationRepository,
	treatmentPlanRepo domain.TreatmentPlanRepository,
	counselingRepo domain.CounselingRepository, // Add this
	timeout time.Duration,
) *RuleEngineUsecase {
	return &RuleEngineUsecase{
		ruleEngine:                    ruleEngine,
		assessmentRepo:                assessmentRepo,
		medicalProfessionalAnswerRepo: medicalProfessionalAnswerRepo,
		clinicalFindingsRepo:          clinicalFindingsRepo,
		classificationRepo:            classificationRepo,
		treatmentPlanRepo:             treatmentPlanRepo,
		counselingRepo:                counselingRepo, // Add this
		contextTimeout:                timeout,
	}
}
// ruleengine/usecase/rule_engine_usecase.go
func (uc *RuleEngineUsecase) StartAssessmentFlow(ctx context.Context, assessmentID uuid.UUID, medicalProfessionalID uuid.UUID) (*domain.AssessmentFlowResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
    defer cancel()

    // Get assessment with already-calculated type
    assessment, err := uc.assessmentRepo.GetByID(ctx, assessmentID, medicalProfessionalID)
    if err != nil {
        return nil, fmt.Errorf("failed to get assessment: %w", err)
    }

    // Convert domain assessment type to tree type
    treeType := "child"
    if assessment.AssessmentType == domain.TypeYoungInfant {
        treeType = "young_infant"
    }

    // Start session with the correct tree type
    ruleEngineSession, err := uc.ruleEngine.StartSession(assessmentID, treeType)
    if err != nil {
        return nil, fmt.Errorf("failed to start rule engine session: %w", err)
    }

    // Store session in database
    answer := &domain.MedicalProfessionalAnswer{
        ID:                 ruleEngineSession.SessionID,
        AssessmentID:       assessmentID,
        Answers:            domain.JSONB{"current_node": ruleEngineSession.CurrentNodeID},
        QuestionSetVersion: "2021",
        ClinicalFindings:   domain.JSONB{},
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    if err := uc.medicalProfessionalAnswerRepo.Upsert(ctx, answer); err != nil {
        return nil, fmt.Errorf("failed to create medical professional answer: %w", err)
    }

    // Get first question
    question, err := uc.getCurrentQuestion(ruleEngineSession)
    if err != nil {
        return nil, fmt.Errorf("failed to get first question: %w", err)
    }

    return &domain.AssessmentFlowResponse{
        SessionID:   answer.ID,
        Question:    question,
        IsComplete:  false,
    }, nil
}
func (uc *RuleEngineUsecase) SubmitAnswer(ctx context.Context, assessmentID uuid.UUID, medicalProfessionalID uuid.UUID, nodeID string, answer interface{}) (*domain.AssessmentFlowResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
    defer cancel()

    // Get assessment and session
    assessment, err := uc.assessmentRepo.GetByID(ctx, assessmentID, medicalProfessionalID)
    if err != nil {
        return nil, fmt.Errorf("failed to get assessment: %w", err)
    }

    mpAnswer, err := uc.medicalProfessionalAnswerRepo.GetByAssessmentID(ctx, assessmentID)
    if err != nil {
        return nil, fmt.Errorf("failed to get medical professional answer: %w", err)
    }

    ruleEngineSession, err := uc.reconstructSessionFromDB(mpAnswer, medicalProfessionalID)
    if err != nil {
        return nil, fmt.Errorf("failed to reconstruct session: %w", err)
    }

    // PROCESS ANSWER AND CHECK FOR CLASSIFICATION
    classification, err := uc.ruleEngine.SubmitAnswer(ruleEngineSession, answer)
    if err != nil {
        // Enhanced medical error messaging with specific guidance
        errorMsg := fmt.Sprintf("clinical assessment error: %v. Please review the assessment findings and ensure all classification nodes are properly configured.", err)
        log.Printf("❌ Rule Engine Error: %v", err)
        return nil, fmt.Errorf(errorMsg)
    }

    // Update session state
    mpAnswer.Answers["current_node"] = ruleEngineSession.CurrentNodeID
    mpAnswer.Answers[nodeID] = answer
    mpAnswer.ClinicalFindings = uc.mapClinicalFindingsToJSON(ruleEngineSession.ClinicalFindings)
    mpAnswer.UpdatedAt = time.Now()

    if err := uc.medicalProfessionalAnswerRepo.Update(ctx, mpAnswer); err != nil {
        return nil, fmt.Errorf("failed to update medical professional answer: %w", err)
    }

    // Handle classification if reached
    if classification != nil {
        log.Printf("🎯 Classification reached: %s (%s)", classification.Name, classification.Color)
        return uc.handleClassification(ctx, assessment, mpAnswer, ruleEngineSession, classification)
    }

    // Continue with next question
    nextQuestion, err := uc.getCurrentQuestion(ruleEngineSession)
    if err != nil {
        return nil, fmt.Errorf("failed to get next question: %w", err)
    }

    log.Printf("➡️  Next question: %s", nextQuestion.NodeID)
    return &domain.AssessmentFlowResponse{
        SessionID:  mpAnswer.ID,
        Question:   nextQuestion,
        IsComplete: false,
    }, nil
}

func (uc *RuleEngineUsecase) handleClassification(ctx context.Context, assessment *domain.Assessment, mpAnswer *domain.MedicalProfessionalAnswer, ruleEngineSession *ruleenginedomain.AssessmentSession, classification *ruleenginedomain.Classification) (*domain.AssessmentFlowResponse, error) {
    // Save clinical findings first
    clinicalFindings := uc.mapToDomainClinicalFindings(assessment.ID, ruleEngineSession.ClinicalFindings)
    if err := uc.clinicalFindingsRepo.Upsert(ctx, clinicalFindings); err != nil {
        return nil, fmt.Errorf("failed to save clinical findings: %w", err)
    }

    // Save classification with proper upsert
    dbClassification := &domain.Classification{
        ID:                     uuid.New(),
        AssessmentID:           assessment.ID,
        Disease:                classification.Name,
        Color:                  classification.Color,
        Details:                "IMCI 2021 Classification",
        RuleVersion:            "2021",
        IsCriticalIllness:      classification.Color == "PINK",
        RequiresUrgentReferral: classification.Color == "PINK",
        TreatmentPriority:      uc.getTreatmentPriority(classification.Color),
        CreatedAt:              time.Now(),
    }

    // Use Upsert to avoid duplicate key errors
    if err := uc.classificationRepo.Upsert(ctx, dbClassification); err != nil {
        return nil, fmt.Errorf("failed to save classification: %w", err)
    }

    // Save treatment plans if any
    if classification.TreatmentPlan != nil {
        if err := uc.saveTreatmentPlans(ctx, assessment.ID, dbClassification.ID, classification.TreatmentPlan); err != nil {
            log.Printf("Warning: Failed to save treatment plans: %v", err)
            // Continue without treatment plans
        }
    }

    // Generate counseling - make it non-blocking
    if err := uc.generateCounseling(ctx, assessment.ID, classification); err != nil {
        log.Printf("Warning: Failed to generate counseling (non-critical): %v", err)
        // Continue without counseling - don't fail the whole assessment
    }

    // Update assessment status
    assessment.Status = domain.StatusCompleted
    endTime := time.Now()
    assessment.EndTime = &endTime

    // Update assessment flags based on classification
    assessment.IsCriticalIllness = dbClassification.IsCriticalIllness
    assessment.RequiresUrgentReferral = dbClassification.RequiresUrgentReferral

    if err := uc.assessmentRepo.Update(ctx, assessment); err != nil {
        return nil, fmt.Errorf("failed to update assessment: %w", err)
    }

    return &domain.AssessmentFlowResponse{
        SessionID:        mpAnswer.ID,
        Classification:   dbClassification,
        ClinicalFindings: clinicalFindings,
        IsComplete:       true,
    }, nil
}

func (uc *RuleEngineUsecase) generateCounseling(ctx context.Context, assessmentID uuid.UUID, classification *ruleenginedomain.Classification) error {
    // First, verify the assessment exists
    _, err := uc.assessmentRepo.GetByID(ctx, assessmentID, uuid.MustParse("medical-professional-id")) // You need to pass the actual medical professional ID
    if err != nil {
        return fmt.Errorf("assessment not found for counseling: %w", err)
    }

    counselingPoints := uc.getCounselingPoints(assessmentID, classification)

    for _, counseling := range counselingPoints {
        if err := uc.counselingRepo.Create(ctx, counseling); err != nil {
            return fmt.Errorf("failed to create counseling: %w", err)
        }
    }

    return nil
}


func (uc *RuleEngineUsecase) getCounselingPoints(assessmentID uuid.UUID, classification *ruleenginedomain.Classification) []*domain.Counseling {
    var counselings []*domain.Counseling

    switch classification.Color {
    case "PINK":
        counselings = append(counselings, &domain.Counseling{
            ID:           uuid.New(),
            AssessmentID: assessmentID, // ✅ Use the actual assessment ID
            AdviceType:   "urgent_referral",
            Details:      "Go to hospital immediately. Do not delay. Keep child warm during transport.",
            Language:     "en",
        })
    case "YELLOW":
        counselings = append(counselings, &domain.Counseling{
            ID:           uuid.New(),
            AssessmentID: assessmentID, // ✅ Use the actual assessment ID
            AdviceType:   "treatment_instructions", 
            Details:      "Give all medications as directed. Return if child gets worse or develops new symptoms.",
            Language:     "en",
        })
    case "GREEN":
        counselings = append(counselings, &domain.Counseling{
            ID:           uuid.New(),
            AssessmentID: assessmentID, // ✅ Use the actual assessment ID
            AdviceType:   "home_care",
            Details:      "Continue breastfeeding and normal feeding. Return if child develops any danger signs.",
            Language:     "en",
        })
    }

    return counselings
}
func (uc *RuleEngineUsecase) mapToDomainClinicalFindings(assessmentID uuid.UUID, findings *ruleenginedomain.ClinicalFindings) *domain.ClinicalFindings {
	return &domain.ClinicalFindings{
		ID:                   uuid.New(),
		AssessmentID:         assessmentID,
		UnableToDrink:        findings.UnableToDrink,
		VomitsEverything:     findings.VomitsEverything,
		ConvulsingNow:        findings.ConvulsingNow,
		LethargicUnconscious: findings.LethargicUnconscious,
		FastBreathing:        findings.FastBreathing,
		ChestIndrawing:       findings.ChestIndrawing,
		Stridor:              findings.Stridor,
		OxygenSaturation:     findings.OxygenSaturation,
		// Map other fields as needed...
	}
}

func (uc *RuleEngineUsecase) mapClinicalFindingsToJSON(findings *ruleenginedomain.ClinicalFindings) domain.JSONB {
	return domain.JSONB{
		"unable_to_drink":       findings.UnableToDrink,
		"vomits_everything":     findings.VomitsEverything,
		"convulsing_now":        findings.ConvulsingNow,
		"lethargic_unconscious": findings.LethargicUnconscious,
		"fast_breathing":        findings.FastBreathing,
		"chest_indrawing":       findings.ChestIndrawing,
		"stridor":               findings.Stridor,
		"oxygen_saturation":     findings.OxygenSaturation,
	}
}

func (uc *RuleEngineUsecase) getTreatmentPriority(color string) int {
	switch color {
	case "PINK":
		return 1
	case "YELLOW":
		return 2
	case "GREEN":
		return 3
	default:
		return 3
	}
}

func (uc *RuleEngineUsecase) saveTreatmentPlans(ctx context.Context, assessmentID, classificationID uuid.UUID, treatmentPlan *ruleenginedomain.TreatmentPlan) error {
	if treatmentPlan == nil {
		return nil
	}

	// Save pre-referral treatments
	for _, treatment := range treatmentPlan.PreReferralTreatments {
		plan := &domain.TreatmentPlan{
			ID:                  uuid.New(),
			AssessmentID:        assessmentID,
			ClassificationID:    classificationID,
			DrugName:            treatment, // For pre-referral, this might be a procedure
			Dosage:              "single dose",
			Frequency:           "once",
			Duration:            "immediate",
			AdministrationRoute: "varies",
			IsPreReferral:       true,
			Instructions:        "Pre-referral treatment",
		}
		if err := uc.treatmentPlanRepo.Create(ctx, plan); err != nil {
			return fmt.Errorf("failed to create pre-referral treatment: %w", err)
		}
	}

	// Save drug treatments
	for _, drug := range treatmentPlan.Drugs {
		plan := &domain.TreatmentPlan{
			ID:                  uuid.New(),
			AssessmentID:        assessmentID,
			ClassificationID:    classificationID,
			DrugName:            drug.Name,
			Dosage:              drug.Dosage,
			Frequency:           "varies", // This should come from your rule engine
			Duration:            drug.Duration,
			AdministrationRoute: "oral", // Default, should come from rules
			IsPreReferral:       false,
			Instructions:        "Follow IMCI 2021 guidelines",
		}
		if err := uc.treatmentPlanRepo.Create(ctx, plan); err != nil {
			return fmt.Errorf("failed to create drug treatment: %w", err)
		}
	}

	return nil
}

func (uc *RuleEngineUsecase) getCurrentQuestion(session *ruleenginedomain.AssessmentSession) (*domain.FlowQuestion, error) {
	node, err := uc.ruleEngine.GetCurrentNode(session)
	if err != nil {
		return nil, err
	}

	question := &domain.FlowQuestion{
		NodeID:      node.ID,
		Type:        string(node.Type),
		Question:    node.Question,
		Instruction: node.Instruction,
	}

	for _, opt := range node.Options {
		question.Options = append(question.Options, domain.Option{
			Value: opt.Value,
			Text:  opt.Text,
		})
	}

	return question, nil
}

// ruleengine/usecase/rule_engine_usecase.go

func (uc *RuleEngineUsecase) reconstructSessionFromDB(mpAnswer *domain.MedicalProfessionalAnswer, medicalProfessionalID uuid.UUID) (*ruleenginedomain.AssessmentSession, error) {
    // Get current node
    currentNodeID, ok := mpAnswer.Answers["current_node"].(string)
    if !ok {
        // Start from beginning if no current node
        assessment, err := uc.assessmentRepo.GetByID(context.Background(), mpAnswer.AssessmentID, medicalProfessionalID) // FIXED: Use actual MP ID
        if err != nil {
            return nil, fmt.Errorf("failed to get assessment for session reconstruction: %w", err)
        }

        treeType := "child"
        if assessment.AssessmentType == domain.TypeYoungInfant {
            treeType = "young_infant"
        }

        return uc.ruleEngine.StartSession(mpAnswer.AssessmentID, treeType)
    }

    // Get assessment to determine type for existing session
    assessment, err := uc.assessmentRepo.GetByID(context.Background(), mpAnswer.AssessmentID, medicalProfessionalID) // FIXED: Use actual MP ID
    if err != nil {
        return nil, fmt.Errorf("failed to get assessment for session reconstruction: %w", err)
    }

    treeType := "child"
    if assessment.AssessmentType == domain.TypeYoungInfant {
        treeType = "young_infant"
    }

    session := &ruleenginedomain.AssessmentSession{
        SessionID:        mpAnswer.ID,
        AssessmentID:     mpAnswer.AssessmentID,
        AssessmentType:   treeType, 
        CurrentNodeID:    currentNodeID,
        Answers:          make(map[string]interface{}),
        ClinicalFindings: &ruleenginedomain.ClinicalFindings{},
    }

    // Reconstruct answers
    for key, value := range mpAnswer.Answers {
        if key != "current_node" {
            session.Answers[key] = value
        }
    }

    // Reconstruct clinical findings
    if mpAnswer.ClinicalFindings != nil {
        findings, err := uc.mapJSONToClinicalFindings(mpAnswer.ClinicalFindings)
        if err != nil {
            return nil, err
        }
        session.ClinicalFindings = findings
    }

    return session, nil
}

func (uc *RuleEngineUsecase) mapJSONToClinicalFindings(jsonData domain.JSONB) (*ruleenginedomain.ClinicalFindings, error) {
    findings := &ruleenginedomain.ClinicalFindings{}
    
    // Map all boolean fields
    boolFields := map[string]*bool{
        "unable_to_drink": &findings.UnableToDrink,
        "vomits_everything": &findings.VomitsEverything,
        "had_convulsions": &findings.HadConvulsions,
        "lethargic_unconscious": &findings.LethargicUnconscious,
        "convulsing_now": &findings.ConvulsingNow,
        "fast_breathing": &findings.FastBreathing,
        "chest_indrawing": &findings.ChestIndrawing,
        "stridor": &findings.Stridor,
        "wheezing": &findings.Wheezing,
        "palms_soles_yellow": &findings.PalmsSolesYellow,
        "skin_eyes_yellow": &findings.SkinEyesYellow,
        "blood_in_stool": &findings.BloodInStool,
        "measles_now": &findings.MeaslesNow,
        "measles_last_3_months": &findings.MeaslesLast3Months,
        "stiff_neck": &findings.StiffNeck,
        "bulging_fontanelle": &findings.BulgingFontanelle,
        "bilateral_edema": &findings.BilateralEdema,
        "cough_present": &findings.CoughPresent,
        "diarrhea_present": &findings.DiarrheaPresent,
        "sunken_eyes": &findings.SunkenEyes,
        "skin_pinch_slow": &findings.SkinPinchSlow,
        "skin_pinch_very_slow": &findings.SkinPinchVerySlow,
        "restless_irritable": &findings.RestlessIrritable,
        "drinking_eagerly": &findings.DrinkingEagerly,
        "drinking_poorly": &findings.DrinkingPoorly,
        "fever_present": &findings.FeverPresent,
        "runny_nose": &findings.RunnyNose,
        "red_eyes": &findings.RedEyes,
        "generalized_rash": &findings.GeneralizedRash,
        "ear_pain": &findings.EarPain,
        "ear_discharge": &findings.EarDischarge,
        "tender_swelling_behind_ear": &findings.TenderSwellingBehindEar,
        "visible_severe_wasting": &findings.VisibleSevereWasting,
        "severe_palmar_pallor": &findings.SeverePalmarPallor,
        "some_palmar_pallor": &findings.SomePalmarPallor,
        "unable_to_feed": &findings.UnableToFeed,
        "not_feeding_well": &findings.NotFeedingWell,
        "movement_only_when_stimulated": &findings.MovementOnlyWhenStimulated,
        "no_movement": &findings.NoMovement,
        "umbilicus_red": &findings.UmbilicusRed,
        "umbilicus_draining_pus": &findings.UmbilicusDrainingPus,
        "skin_pustules": &findings.SkinPustules,
        "low_body_temperature": &findings.LowBodyTemperature,
        "hiv_exposed": &findings.HIVExposed,
        "hiv_status_known": &findings.HIVStatusKnown,
        "tb_contact_history": &findings.TBContactHistory,
        "tb_weight_loss": &findings.TBWeightLoss,
        "night_sweats": &findings.NightSweats,
        "suspected_developmental_delay": &findings.SuspectedDevelopmentalDelay,
        "breastfeeding": &findings.Breastfeeding,
        "complementary_foods": &findings.ComplementaryFoods,
        "feeding_problem": &findings.FeedingProblem,
        "underweight": &findings.Underweight,
    }
    
    for field, ptr := range boolFields {
        if val, exists := jsonData[field]; exists {
            if b, ok := val.(bool); ok {
                *ptr = b
            }
        }
    }
    
intFields := map[string]**int{
    "oxygen_saturation": &findings.OxygenSaturation,
    "respiratory_rate": &findings.RespiratoryRate,
    "cough_duration_days": &findings.CoughDurationDays,
    "diarrhea_duration_days": &findings.DiarrheaDurationDays,
    "fever_duration_days": &findings.FeverDurationDays,
    "jaundice_age_hours": &findings.JaundiceAgeHours,
    "ear_discharge_duration_days": &findings.EarDischargeDurationDays,
    "tb_cough_duration_days": &findings.TBCoughDurationDays,
    "breastfeeding_frequency": &findings.BreastfeedingFrequency,
}

    
    for field, ptr := range intFields {
        if val, exists := jsonData[field]; exists && val != nil {
            if f, ok := val.(float64); ok {
                intVal := int(f)
                *ptr = &intVal
            }
        }
    }
    
  floatFields := map[string]**float64{
    "muac": &findings.MUAC,
    "weight_for_height_z_score": &findings.WeightForHeightZScore,
    "hb_level": &findings.HbLevel,
    "hct_level": &findings.HctLevel,
    "body_temperature": &findings.BodyTemperature,
}
    
    for field, ptr := range floatFields {
        if val, exists := jsonData[field]; exists && val != nil {
            if f, ok := val.(float64); ok {
                *ptr = &f
            }
        }
    }
    
    // Map string array fields
    if val, exists := jsonData["milestones_absent"]; exists {
        if milestones, ok := val.([]interface{}); ok {
            var strMilestones []string
            for _, m := range milestones {
                if str, ok := m.(string); ok {
                    strMilestones = append(strMilestones, str)
                }
            }
            findings.MilestonesAbsent = strMilestones
        }
    }
    
    if val, exists := jsonData["risk_factors_present"]; exists {
        if risks, ok := val.([]interface{}); ok {
            var strRisks []string
            for _, r := range risks {
                if str, ok := r.(string); ok {
                    strRisks = append(strRisks, str)
                }
            }
            findings.RiskFactorsPresent = strRisks
        }
    }
    
    // Map other findings
    if val, exists := jsonData["other_findings"]; exists {
        if other, ok := val.(string); ok {
            findings.OtherFindings = other
        }
    }
    
    return findings, nil
}