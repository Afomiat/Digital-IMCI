package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	ruleenginedomain "github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/Afomiat/Digital-IMCI/ruleengine/engine"
	"github.com/google/uuid"
)

type RuleEngineUsecase struct {
	ruleEngine                      *engine.RuleEngine
	assessmentRepo                  domain.AssessmentRepository
	medicalProfessionalAnswerRepo   domain.MedicalProfessionalAnswerRepository
	clinicalFindingsRepo            domain.ClinicalFindingsRepository
	classificationRepo              domain.ClassificationRepository
	treatmentPlanRepo               domain.TreatmentPlanRepository
	counselingRepo                  domain.CounselingRepository
	contextTimeout                  time.Duration
}

func NewRuleEngineUsecase(
	ruleEngine *engine.RuleEngine,
	assessmentRepo domain.AssessmentRepository,
	medicalProfessionalAnswerRepo domain.MedicalProfessionalAnswerRepository,
	clinicalFindingsRepo domain.ClinicalFindingsRepository,
	classificationRepo domain.ClassificationRepository,
	treatmentPlanRepo domain.TreatmentPlanRepository,
	counselingRepo domain.CounselingRepository,
	timeout time.Duration,
) *RuleEngineUsecase {
	return &RuleEngineUsecase{
		ruleEngine:                    ruleEngine,
		assessmentRepo:                assessmentRepo,
		medicalProfessionalAnswerRepo: medicalProfessionalAnswerRepo,
		clinicalFindingsRepo:          clinicalFindingsRepo,
		classificationRepo:            classificationRepo,
		treatmentPlanRepo:             treatmentPlanRepo,
		counselingRepo:                counselingRepo,
		contextTimeout:                timeout,
	}
}

func (uc *RuleEngineUsecase) StartAssessmentFlow(ctx context.Context, req ruleenginedomain.StartFlowRequest, medicalProfessionalID uuid.UUID) (*ruleenginedomain.StartFlowResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	assessment, err := uc.assessmentRepo.GetByID(ctx, req.AssessmentID, medicalProfessionalID)
	if err != nil {
		return nil, err
	}

	flow, err := uc.ruleEngine.StartAssessmentFlow(req.AssessmentID, req.TreeID)
	if err != nil {
		return nil, err
	}

	medicalProfessionalAnswer := &domain.MedicalProfessionalAnswer{
		ID:                 uuid.New(),
		AssessmentID:       req.AssessmentID,
		Answers:            domain.JSONB(flow.Answers),
		QuestionSetVersion: req.TreeID,
		ClinicalFindings:   domain.JSONB{},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := uc.medicalProfessionalAnswerRepo.Upsert(ctx, medicalProfessionalAnswer); err != nil {
		return nil, fmt.Errorf("failed to save assessment flow: %w", err)
	}

	assessment.Status = domain.StatusInProgress
	if err := uc.assessmentRepo.Update(ctx, assessment); err != nil {
		return nil, fmt.Errorf("failed to update assessment status: %w", err)
	}

	currentQuestion, err := uc.ruleEngine.GetCurrentQuestion(flow)
	if err != nil {
		return nil, err
	}

	return &ruleenginedomain.StartFlowResponse{
		SessionID:   medicalProfessionalAnswer.ID,
		Question:    currentQuestion,
		IsComplete:  flow.Status == ruleenginedomain.FlowStatusCompleted,
		CurrentNode: flow.CurrentNode,
	}, nil
}

func (uc *RuleEngineUsecase) SubmitAnswer(ctx context.Context, req ruleenginedomain.SubmitAnswerRequest, medicalProfessionalID uuid.UUID) (*ruleenginedomain.SubmitAnswerResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	assessment, err := uc.assessmentRepo.GetByID(ctx, req.AssessmentID, medicalProfessionalID)
	if err != nil {
		return nil, err
	}

	medicalProfessionalAnswer, err := uc.medicalProfessionalAnswerRepo.GetByAssessmentID(ctx, req.AssessmentID)
	if err != nil {
		return nil, domain.ErrMedicalProfessionalAnswerNotFound
	}

	flow := &ruleenginedomain.AssessmentFlow{
		AssessmentID: req.AssessmentID,
		TreeID:       medicalProfessionalAnswer.QuestionSetVersion, 
		Answers:      map[string]interface{}(medicalProfessionalAnswer.Answers),
		Status:       ruleenginedomain.FlowStatusInProgress,
		CreatedAt:    medicalProfessionalAnswer.CreatedAt,
		UpdatedAt:    medicalProfessionalAnswer.UpdatedAt,
	}

	if flow.CurrentNode == "" {
		tree, err := uc.ruleEngine.GetAssessmentTree(flow.TreeID)
		if err != nil {
			return nil, err
		}
		flow.CurrentNode = tree.StartNode
	}

	updatedFlow, nextQuestion, err := uc.ruleEngine.SubmitAnswer(flow, req.NodeID, req.Answer)
	if err != nil {
		return nil, err
	}

	medicalProfessionalAnswer.Answers = domain.JSONB(updatedFlow.Answers)
	medicalProfessionalAnswer.UpdatedAt = time.Now()

	if err := uc.medicalProfessionalAnswerRepo.Upsert(ctx, medicalProfessionalAnswer); err != nil {
		return nil, fmt.Errorf("failed to update assessment flow: %w", err)
	}

	if updatedFlow.Status == ruleenginedomain.FlowStatusCompleted || updatedFlow.Status == ruleenginedomain.FlowStatusEmergency {
		if err := uc.saveClassificationResults(ctx, assessment, updatedFlow.Classification); err != nil {
			return nil, fmt.Errorf("failed to save classification results: %w", err)
		}

		assessment.Status = domain.StatusCompleted
		if err := uc.assessmentRepo.Update(ctx, assessment); err != nil {
			return nil, fmt.Errorf("failed to update assessment status: %w", err)
		}
	}

	return &ruleenginedomain.SubmitAnswerResponse{
		SessionID:      medicalProfessionalAnswer.ID,
		Question:       nextQuestion,
		Classification: updatedFlow.Classification,
		IsComplete:     updatedFlow.Status != ruleenginedomain.FlowStatusInProgress,
		CurrentNode:    updatedFlow.CurrentNode,
		Status:         updatedFlow.Status,
	}, nil
}

func (uc *RuleEngineUsecase) ProcessBatchAssessment(ctx context.Context, req ruleenginedomain.BatchProcessRequest, medicalProfessionalID uuid.UUID) (*ruleenginedomain.BatchProcessResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	assessment, err := uc.assessmentRepo.GetByID(ctx, req.AssessmentID, medicalProfessionalID)
	if err != nil {
		return nil, err
	}

	flow, err := uc.ruleEngine.ProcessBatchAssessment(req.AssessmentID, req.TreeID, req.Answers)
	if err != nil {
		return nil, err
	}

	medicalProfessionalAnswer := &domain.MedicalProfessionalAnswer{
		ID:                 uuid.New(),
		AssessmentID:       req.AssessmentID,
		Answers:            domain.JSONB(req.Answers),
		QuestionSetVersion: req.TreeID,
		ClinicalFindings:   domain.JSONB{},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := uc.medicalProfessionalAnswerRepo.Upsert(ctx, medicalProfessionalAnswer); err != nil {
		return nil, fmt.Errorf("failed to save assessment answers: %w", err)
	}

	if flow.Classification != nil {
		if err := uc.saveClassificationResults(ctx, assessment, flow.Classification); err != nil {
			return nil, fmt.Errorf("failed to save classification results: %w", err)
		}

		assessment.Status = domain.StatusCompleted
		if err := uc.assessmentRepo.Update(ctx, assessment); err != nil {
			return nil, fmt.Errorf("failed to update assessment status: %w", err)
		}
	}

	return &ruleenginedomain.BatchProcessResponse{
		AssessmentID:  req.AssessmentID,
		Classification: flow.Classification,
		Status:         flow.Status,
	}, nil
}

func (uc *RuleEngineUsecase) GetTreeQuestions(treeID string) (*ruleenginedomain.AssessmentTree, error) {
	tree, err := uc.ruleEngine.GetAssessmentTree(treeID)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func (uc *RuleEngineUsecase) GetAssessmentTree(treeID string) (*ruleenginedomain.AssessmentTree, error) {
	return uc.ruleEngine.GetAssessmentTree(treeID)
}


func (uc *RuleEngineUsecase) saveClassificationResults(ctx context.Context, assessment *domain.Assessment, classification *ruleenginedomain.ClassificationResult) error {
	if classification == nil {
		return nil
	}

	class := &domain.Classification{
		ID:                    uuid.New(),
		AssessmentID:          assessment.ID,
		Disease:               classification.Classification,
		Color:                 classification.Color,
		Details:               classification.TreatmentPlan,
		RuleVersion:           "imci_2021_v1", 
		IsCriticalIllness:     classification.Emergency,
		RequiresUrgentReferral: classification.Emergency,
		TreatmentPriority:     uc.getTreatmentPriority(classification.Classification),
		CreatedAt:             time.Now(),
	}

	if err := uc.classificationRepo.Create(ctx, class); err != nil {
		return err
	}

	if err := uc.saveTreatmentPlans(ctx, class, classification); err != nil {
		return err
	}

	counseling := &domain.Counseling{
		ID:           uuid.New(),
		AssessmentID: assessment.ID,
		AdviceType:   "mother_advice",
		Details:      classification.MotherAdvice,
		Language:     "en",
		CreatedAt:    time.Now(),
	}

	if err := uc.counselingRepo.Create(ctx, counseling); err != nil {
		return err
	}

	if len(classification.FollowUp) > 0 {
		followUpCounseling := &domain.Counseling{
			ID:           uuid.New(),
			AssessmentID: assessment.ID,
			AdviceType:   "follow_up_schedule",
			Details:      fmt.Sprintf("Follow-up schedule: %v", strings.Join(classification.FollowUp, ", ")),
			Language:     "en",
			CreatedAt:    time.Now(),
		}
		if err := uc.counselingRepo.Create(ctx, followUpCounseling); err != nil {
			return err
		}
	}

	return nil
}


func (uc *RuleEngineUsecase) getTreatmentPriority(classification string) int {
	switch classification {
	case "CRITICAL ILLNESS", "VERY SEVERE DISEASE", "VERY LOW BIRTH WEIGHT AND/OR VERY PRETERM":
		return 1
	case "PNEUMONIA", "LOCAL BACTERIAL INFECTION", "LOW BIRTH WEIGHT AND/OR PRETERM":
		return 2
	default:
		return 3
	}
}

func (uc *RuleEngineUsecase) saveTreatmentPlans(ctx context.Context, classification *domain.Classification, result *ruleenginedomain.ClassificationResult) error {
	var plans []*domain.TreatmentPlan
	
	switch result.Classification {
	case "CRITICAL ILLNESS", "VERY SEVERE DISEASE":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Ampicillin",
				Dosage:              "First dose",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM/IV",
				IsPreReferral:       true,
				Instructions:        "Give before referral to hospital",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Gentamicin",
				Dosage:              "First dose",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM/IV",
				IsPreReferral:       true,
				Instructions:        "Give before referral to hospital",
			},
		}
	case "PNEUMONIA":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Ampicillin",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "7 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Complete full course of antibiotics",
			},
		}
	case "LOCAL BACTERIAL INFECTION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Ampicillin",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "5 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Teach mother to treat local infections at home",
			},
		}
	case "SEVERE JAUNDICE":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Glucose",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral/NG",
				IsPreReferral:       true,
				Instructions:        "Treat to prevent low blood sugar before referral",
			},
		}
	case "SEVERE DISEASE - BLOOD IN STOOL":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Ampicillin",
				Dosage:              "First dose IM",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM",
				IsPreReferral:       true,
				Instructions:        "Give before referral to hospital",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Gentamicin",
				Dosage:              "First dose IM",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM",
				IsPreReferral:       true,
				Instructions:        "Give before referral to hospital",
			},
		}
	case "SEVERE DEHYDRATION": 
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "ORS",
				Dosage:              "Frequent sips",
				Frequency:           "During transport",
				Duration:            "Until hospital arrival",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Give frequent sips during transport to hospital",
			},
		}
	case "SOME DEHYDRATION", "PROLONGED DIARRHEA": 
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "ORS",
				Dosage:              "Plan B",
				Frequency:           "As directed",
				Duration:            "Until diarrhea stops",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give for some dehydration",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Zinc sulfate",
				Dosage:              "10-20mg daily",
				Frequency:           "Once daily",
				Duration:            "10-14 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give zinc supplement",
			},
		}
	case "NO DEHYDRATION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "ORS",
				Dosage:              "Plan A",
				Frequency:           "After each loose stool",
				Duration:            "Until diarrhea stops",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give to treat diarrhea at home",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Zinc sulfate",
				Dosage:              "10-20mg daily",
				Frequency:           "Once daily",
				Duration:            "10-14 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give zinc supplement",
			},
		}
	case "FEEDING PROBLEM OR UNDERWEIGHT": 
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Breastfeeding Counseling",
				Dosage:              "N/A",
				Frequency:           "As needed",
				Duration:            "Until resolved",
				AdministrationRoute: "Counseling",
				IsPreReferral:       false,
				Instructions:        "Teach correct positioning and attachment",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Nutritional Support",
				Dosage:              "N/A",
				Frequency:           "Daily",
				Duration:            "Until weight improves",
				AdministrationRoute: "Dietary",
				IsPreReferral:       false,
				Instructions:        "Increase feeding frequency and ensure adequate nutrition",
			},
		}
	case "HIV INFECTED":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Cotrimoxazole",
				Dosage:              "Based on weight",
				Frequency:           "Once daily",
				Duration:            "Until further evaluation",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Start prophylaxis from 6 weeks of age",
			},
		}
	case "HIV EXPOSED":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Cotrimoxazole",
				Dosage:              "Based on weight",
				Frequency:           "Once daily",
				Duration:            "Until HIV status confirmed negative",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Start prophylaxis from 6 weeks of age",
			},
		}
		if strings.Contains(result.MotherAdvice, "thrush") {
			plans = append(plans, &domain.TreatmentPlan{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Nystatin",
				Dosage:              "As prescribed",
				Frequency:           "As directed",
				Duration:            "7-14 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Treat oral thrush",
			})
		}
	
	case "VERY LOW BIRTH WEIGHT AND/OR VERY PRETERM":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin K",
				Dosage:              "0.5mg",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM",
				IsPreReferral:       true,
				Instructions:        "Give on anterior mid lateral thigh before referral",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Kangaroo Mother Care",
				Dosage:              "N/A",
				Frequency:           "Continuous",
				Duration:            "Until hospital transfer",
				AdministrationRoute: "Positioning",
				IsPreReferral:       true,
				Instructions:        "Start KMC and maintain during referral",
			},
		}
	case "LOW BIRTH WEIGHT AND/OR PRETERM":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin K",
				Dosage:              "1mg (0.5mg if GA <34 weeks)",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM",
				IsPreReferral:       false,
				Instructions:        "Give on anterior mid lateral thigh",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Kangaroo Mother Care",
				Dosage:              "If <2000g",
				Frequency:           "Continuous",
				Duration:            "Until weight ≥2500g",
				AdministrationRoute: "Positioning",
				IsPreReferral:       false,
				Instructions:        "Practice KMC in health facility or hospital",
			},
		}
	case "NORMAL BIRTH WEIGHT AND/OR TERM":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin K",
				Dosage:              "1mg",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM",
				IsPreReferral:       false,
				Instructions:        "Give on anterior mid thigh",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "First Vaccine",
				Dosage:              "As per schedule",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM",
				IsPreReferral:       false,
				Instructions:        "Give first dose of vaccine",
			},
		}

		
	}
	

	for _, plan := range plans {
		if err := uc.treatmentPlanRepo.Create(ctx, plan); err != nil {
			return err
		}
	}

	return nil
}

