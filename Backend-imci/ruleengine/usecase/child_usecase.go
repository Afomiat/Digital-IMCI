// ruleengine/usecase/child_usecase.go
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

type ChildRuleEngineUsecase struct {
	ruleEngine                    *engine.ChildRuleEngine
	assessmentRepo                domain.AssessmentRepository
	medicalProfessionalAnswerRepo domain.MedicalProfessionalAnswerRepository
	clinicalFindingsRepo          domain.ClinicalFindingsRepository
	classificationRepo            domain.ClassificationRepository
	treatmentPlanRepo             domain.TreatmentPlanRepository
	counselingRepo                domain.CounselingRepository
	contextTimeout                time.Duration
}

func NewChildRuleEngineUsecase(
	ruleEngine *engine.ChildRuleEngine,
	assessmentRepo domain.AssessmentRepository,
	medicalProfessionalAnswerRepo domain.MedicalProfessionalAnswerRepository,
	clinicalFindingsRepo domain.ClinicalFindingsRepository,
	classificationRepo domain.ClassificationRepository,
	treatmentPlanRepo domain.TreatmentPlanRepository,
	counselingRepo domain.CounselingRepository,
	timeout time.Duration,
) *ChildRuleEngineUsecase {
	return &ChildRuleEngineUsecase{
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

func (uc *ChildRuleEngineUsecase) StartAssessmentFlow(ctx context.Context, req ruleenginedomain.StartFlowRequest, medicalProfessionalID uuid.UUID) (*ruleenginedomain.StartFlowResponse, error) {
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

func (uc *ChildRuleEngineUsecase) SubmitAnswer(ctx context.Context, req ruleenginedomain.SubmitAnswerRequest, medicalProfessionalID uuid.UUID) (*ruleenginedomain.SubmitAnswerResponse, error) {
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

func (uc *ChildRuleEngineUsecase) ProcessBatchAssessment(ctx context.Context, req ruleenginedomain.BatchProcessRequest, medicalProfessionalID uuid.UUID) (*ruleenginedomain.BatchProcessResponse, error) {
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
		AssessmentID:   req.AssessmentID,
		Classification: flow.Classification,
		Status:         flow.Status,
	}, nil
}

func (uc *ChildRuleEngineUsecase) GetTreeQuestions(treeID string) (*ruleenginedomain.AssessmentTree, error) {
	tree, err := uc.ruleEngine.GetAssessmentTree(treeID)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func (uc *ChildRuleEngineUsecase) GetAssessmentTree(treeID string) (*ruleenginedomain.AssessmentTree, error) {
	return uc.ruleEngine.GetAssessmentTree(treeID)
}

func (uc *ChildRuleEngineUsecase) GetAvailableTrees() []string {
	return uc.ruleEngine.GetAvailableTrees()
}

func (uc *ChildRuleEngineUsecase) saveClassificationResults(ctx context.Context, assessment *domain.Assessment, classification *ruleenginedomain.ClassificationResult) error {
	if classification == nil {
		return nil
	}

	class := &domain.Classification{
		ID:                     uuid.New(),
		AssessmentID:           assessment.ID,
		Disease:                classification.Classification,
		Color:                  classification.Color,
		Details:                classification.TreatmentPlan,
		RuleVersion:            "imnci_2021_v1",
		IsCriticalIllness:      classification.Emergency,
		RequiresUrgentReferral: classification.Emergency,
		TreatmentPriority:      uc.getTreatmentPriority(classification.Classification),
		CreatedAt:              time.Now(),
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

func (uc *ChildRuleEngineUsecase) getTreatmentPriority(classification string) int {
	switch classification {
	case "VERY SEVERE DISEASE", "SEVERE PNEUMONIA OR VERY SEVERE DISEASE", "SEVERE DEHYDRATION", "SEVERE PERSISTENT DIARRHOEA", "SEVERE MALNUTRITION", "VERY SEVERE FEBRILE DISEASE", "SEVERE COMPLICATED MEASLES", "MASTOIDITIS", "SEVERE ANEMIA", "COMPLICATED SEVERE ACUTE MALNUTRITION":
		return 1
	case "PNEUMONIA", "SOME DEHYDRATION", "PERSISTENT DIARRHOEA", "DYSENTERY", "FEVER - MALARIA RISK", "ACUTE EAR INFECTION", "CHRONIC EAR INFECTION", "MALARIA_HIGH_RISK", "MALARIA_LOW_RISK", "MEASLES WITH EYE OR MOUTH COMPLICATIONS", "ANEMIA", "UNCOMPLICATED SEVERE ACUTE MALNUTRITION", "MODERATE ACUTE MALNUTRITION":
		return 2
	case "NO COUGH OR DIFFICULT BREATHING", "COUGH OR COLD", "NO DEHYDRATION", "NO MALNUTRITION", "NO MALARIA RISK", "FEVER_NO_MALARIA", "MEASLES_NO_COMPLICATIONS", "NO EAR INFECTION", "NO ANEMIA", "NO ACUTE MALNUTRITION":
		return 3
	default:
		return 3
	}
}
func (uc *ChildRuleEngineUsecase) saveTreatmentPlans(ctx context.Context, classification *domain.Classification, result *ruleenginedomain.ClassificationResult) error {
	var plans []*domain.TreatmentPlan

	switch result.Classification {
	case "VERY SEVERE DISEASE":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "First dose antibiotic",
				Dosage:              "Based on weight",
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
				DrugName:            "Vitamin A",
				Dosage:              "Based on age",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Give if not given in last month",
			},
		}
	case "SEVERE PNEUMONIA OR VERY SEVERE DISEASE":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "First dose of IV/IM Ampicillin and Gentamicin",
				Dosage:              "Based on weight",
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
				DrugName:            "Amoxicillin",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "5 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give oral Amoxicillin for 5 days",
			},
		}
	case "PNEUMONIA WITH WHEEZING":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Amoxicillin",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "5 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give oral Amoxicillin for 5 days",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Rapid acting inhaled bronchodilator",
				Dosage:              "Based on weight",
				Frequency:           "Up to 3 times, 15-20 minutes apart",
				Duration:            "As needed",
				AdministrationRoute: "Inhaled",
				IsPreReferral:       false,
				Instructions:        "Give rapid acting inhaled bronchodilator for up to 3 times, 15-20 minutes apart",
			},
		}
	case "CHEST INDRAWING HIV EXPOSED":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "First dose of amoxicillin",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Give first dose before referral",
			},
		}
	case "COUGH OR COLD":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Symptomatic relief",
				Dosage:              "As needed",
				Frequency:           "As directed",
				Duration:            "Until symptoms resolve",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Soothe throat and relieve cough with safe remedy",
			},
		}
	case "COUGH OR COLD WITH WHEEZING":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Inhaled bronchodilator",
				Dosage:              "Based on weight",
				Frequency:           "As needed for 5 days",
				Duration:            "5 days",
				AdministrationRoute: "Inhaled",
				IsPreReferral:       false,
				Instructions:        "Give inhaled bronchodilator for 5 days",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Symptomatic relief",
				Dosage:              "As needed",
				Frequency:           "As directed",
				Duration:            "Until symptoms resolve",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Soothe throat and relieve cough with safe remedy",
			},
		}
	case "SEVERE DEHYDRATION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "ORS Plan C",
				Dosage:              "Based on weight",
				Frequency:           "During transport",
				Duration:            "Until hospital arrival",
				AdministrationRoute: "Oral/NG",
				IsPreReferral:       true,
				Instructions:        "Give fluid for severe dehydration (Plan C)",
			},
		}
	case "SOME DEHYDRATION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "ORS Plan B",
				Dosage:              "Based on weight",
				Frequency:           "As directed",
				Duration:            "Until diarrhea stops",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give fluid for some dehydration (Plan B)",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Zinc sulfate",
				Dosage:              "20mg daily",
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
				DrugName:            "ORS Plan A",
				Dosage:              "After each loose stool",
				Frequency:           "As needed",
				Duration:            "Until diarrhea stops",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give fluid to treat diarrhea at home (Plan A)",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Zinc sulfate",
				Dosage:              "20mg daily",
				Frequency:           "Once daily",
				Duration:            "10-14 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give zinc supplement",
			},
		}
	case "SEVERE PERSISTENT DIARRHOEA":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin A",
				Dosage:              "Based on age",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Give Vitamin A before referral",
			},
		}
	case "PERSISTENT DIARRHOEA":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin A",
				Dosage:              "Therapeutic dose based on age",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give Vitamin A therapeutic dose",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Zinc sulfate",
				Dosage:              "20mg daily",
				Frequency:           "Once daily",
				Duration:            "10 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give zinc for 10 days",
			},
		}
	case "DYSENTERY":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Ciprofloxacin",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "3 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Treat for 3 days with Ciprofloxacin",
			},
		}
	case "VERY SEVERE FEBRILE DISEASE":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "First dose of IV/IM Ampicillin and Gentamicin",
				Dosage:              "Based on weight",
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
				DrugName:            "First dose IV/IM Artesunate",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM/IV",
				IsPreReferral:       true,
				Instructions:        "Give for severe malaria if high malaria risk",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Paracetamol",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Give for high fever (≥38.5°C) in health facility",
			},
		}
	case "MALARIA_HIGH_RISK", "MALARIA_LOW_RISK":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Artemisinin-Lumefantrine (AL)",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "3 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Treat for P. falciparum or mixed infection",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Primaquine",
				Dosage:              "Based on weight",
				Frequency:           "Once daily",
				Duration:            "14 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give for P. falciparum gametocytes",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Paracetamol",
				Dosage:              "Based on weight",
				Frequency:           "As needed",
				Duration:            "Until fever resolves",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give for high fever (≥38.5°C)",
			},
		}
	case "FEVER_NO_MALARIA":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Paracetamol",
				Dosage:              "Based on weight",
				Frequency:           "As needed",
				Duration:            "Until fever resolves",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give one dose for high fever (≥38.5°C)",
			},
		}
	case "SEVERE COMPLICATED MEASLES":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin A",
				Dosage:              "Based on age",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Give first dose before referral",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "IV/IM Ampicillin and Gentamicin",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM/IV",
				IsPreReferral:       true,
				Instructions:        "Give first dose before referral",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Tetracycline eye ointment",
				Dosage:              "Apply to both eyes",
				Frequency:           "4 times daily",
				Duration:            "7 days",
				AdministrationRoute: "Topical",
				IsPreReferral:       true,
				Instructions:        "Apply if clouding cornea or pus draining from eye",
			},
		}
	case "MEASLES WITH EYE OR MOUTH COMPLICATIONS":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin A",
				Dosage:              "Therapeutic dose based on age",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give therapeutic dose",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Tetracycline eye ointment",
				Dosage:              "Apply to affected eye",
				Frequency:           "3 times daily",
				Duration:            "7 days",
				AdministrationRoute: "Topical",
				IsPreReferral:       false,
				Instructions:        "Apply if pus draining from eye",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Gentian Violet",
				Dosage:              "Apply to mouth ulcers",
				Frequency:           "Twice daily",
				Duration:            "Until healed",
				AdministrationRoute: "Topical",
				IsPreReferral:       false,
				Instructions:        "Apply to mouth ulcers",
			},
		}
	case "MEASLES_NO_COMPLICATIONS":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Vitamin A",
				Dosage:              "Therapeutic dose based on age",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give therapeutic dose",
			},
		}
	case "MASTOIDITIS":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Ceftriaxone",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IV/IM",
				IsPreReferral:       true,
				Instructions:        "Give first dose before referral to hospital",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Paracetamol",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Give for pain relief before referral",
			},
		}
	case "ACUTE EAR INFECTION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Amoxicillin",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "5 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give oral Amoxicillin for 5 days",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Paracetamol",
				Dosage:              "Based on weight",
				Frequency:           "As needed",
				Duration:            "Until pain resolves",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give for pain relief",
			},
		}
	case "CHRONIC EAR INFECTION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Quinolone eardrops",
				Dosage:              "3-4 drops",
				Frequency:           "Twice daily",
				Duration:            "2 weeks",
				AdministrationRoute: "Topical",
				IsPreReferral:       false,
				Instructions:        "Apply topical quinolone eardrops for 2 weeks",
			},
		}
	case "NO EAR INFECTION":
		plans = []*domain.TreatmentPlan{}
	case "SEVERE ANEMIA":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Urgent referral",
				Dosage:              "N/A",
				Frequency:           "Immediate",
				Duration:            "N/A",
				AdministrationRoute: "N/A",
				IsPreReferral:       true,
				Instructions:        "Refer URGENTLY to hospital for severe anemia management",
			},
		}
	case "ANEMIA":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Iron supplement",
				Dosage:              "Based on weight and age",
				Frequency:           "Once daily",
				Duration:            "3 months",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give iron supplementation for anemia",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Albendazole/Mebendazole",
				Dosage:              "Based on age",
				Frequency:           "Single dose",
				Duration:            "Once",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give if child ≥ 1 year and no dose in previous 6 months",
			},
		}
	case "NO ANEMIA":

		plans = []*domain.TreatmentPlan{}

	case "COMPLICATED SEVERE ACUTE MALNUTRITION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Ampicillin and Gentamicin",
				Dosage:              "Based on weight",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "IM",
				IsPreReferral:       true,
				Instructions:        "Give 1st dose of Ampicillin and Gentamicin IM before referral",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Sugar solution",
				Dosage:              "10ml/kg",
				Frequency:           "Stat",
				Duration:            "Single dose",
				AdministrationRoute: "Oral",
				IsPreReferral:       true,
				Instructions:        "Treat the child to prevent low blood sugar",
			},
		}
	case "UNCOMPLICATED SEVERE ACUTE MALNUTRITION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "RUTF (Ready-to-Use Therapeutic Food)",
				Dosage:              "Based on weight",
				Frequency:           "Multiple times daily",
				Duration:            "7 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give RUTF for 7 days as per OTP protocol",
			},
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Amoxicillin",
				Dosage:              "Based on weight",
				Frequency:           "Twice daily",
				Duration:            "5 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Give oral Amoxicillin for 5 days",
			},
		}
	case "MODERATE ACUTE MALNUTRITION":
		plans = []*domain.TreatmentPlan{
			{
				ID:                  uuid.New(),
				AssessmentID:        classification.AssessmentID,
				ClassificationID:    classification.ID,
				DrugName:            "Supplementary feeding",
				Dosage:              "As per TSFP protocol",
				Frequency:           "Daily",
				Duration:            "30 days",
				AdministrationRoute: "Oral",
				IsPreReferral:       false,
				Instructions:        "Follow TSFP care protocol for nutritional support",
			},
		}
	case "NO ACUTE MALNUTRITION":
		plans = []*domain.TreatmentPlan{}

	default:
		plans = []*domain.TreatmentPlan{}
	}

	for _, plan := range plans {
		if err := uc.treatmentPlanRepo.Create(ctx, plan); err != nil {
			return err
		}
	}

	return nil
}
