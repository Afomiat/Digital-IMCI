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
	ruleEngine                      *engine.ChildRuleEngine
	assessmentRepo                  domain.AssessmentRepository
	medicalProfessionalAnswerRepo   domain.MedicalProfessionalAnswerRepository
	clinicalFindingsRepo            domain.ClinicalFindingsRepository
	classificationRepo              domain.ClassificationRepository
	treatmentPlanRepo               domain.TreatmentPlanRepository
	counselingRepo                  domain.CounselingRepository
	contextTimeout                  time.Duration
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

func (uc *ChildRuleEngineUsecase) saveClassificationResults(ctx context.Context, assessment *domain.Assessment, classification *ruleenginedomain.ClassificationResult) error {
	if classification == nil {
		return nil
	}

	class := &domain.Classification{
		ID:                    uuid.New(),
		AssessmentID:          assessment.ID,
		Disease:               classification.Classification,
		Color:                 classification.Color,
		Details:               classification.TreatmentPlan,
		RuleVersion:           "imnci_2021_v1",
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

func (uc *ChildRuleEngineUsecase) getTreatmentPriority(classification string) int {
	switch classification {
	case "VERY SEVERE DISEASE", "SEVERE PNEUMONIA OR VERY SEVERE DISEASE", "SEVERE DEHYDRATION", "SEVERE MALNUTRITION":
		return 1
	case "PNEUMONIA", "SOME DEHYDRATION", "FEVER - MALARIA RISK", "ACUTE EAR INFECTION":
		return 2
	case "NO COUGH OR DIFFICULT BREATHING", "COUGH OR COLD", "NO DEHYDRATION", "NO MALNUTRITION", "NO MALARIA RISK":
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
	case "NO COUGH OR DIFFICULT BREATHING":
		// No specific treatment plans needed
		return nil
	}

	for _, plan := range plans {
		if err := uc.treatmentPlanRepo.Create(ctx, plan); err != nil {
			return err
		}
	}

	return nil
}