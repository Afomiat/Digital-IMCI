package usecase

import (
	"context"
	"fmt"
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

// FIXED: Add AssessmentID back but without JSON tag since it comes from URL
type StartFlowRequest struct {
	AssessmentID uuid.UUID
	TreeID       string `json:"tree_id" binding:"required"`
}

type StartFlowResponse struct {
	SessionID   uuid.UUID                  `json:"session_id"`
	Question    *ruleenginedomain.Question `json:"question,omitempty"`
	IsComplete  bool                       `json:"is_complete"`
	CurrentNode string                     `json:"current_node"`
}

// FIXED: Remove AssessmentID from JSON tag since it comes from URL
type SubmitAnswerRequest struct {
	AssessmentID uuid.UUID
	NodeID       string      `json:"node_id" binding:"required"`
	Answer       interface{} `json:"answer" binding:"required"`
}

type SubmitAnswerResponse struct {
	SessionID      uuid.UUID                           `json:"session_id"`
	Question       *ruleenginedomain.Question          `json:"question,omitempty"`
	Classification *ruleenginedomain.ClassificationResult `json:"classification,omitempty"`
	IsComplete     bool                                `json:"is_complete"`
	CurrentNode    string                              `json:"current_node"`
	Status         ruleenginedomain.FlowStatus         `json:"status"`
}

func (uc *RuleEngineUsecase) StartAssessmentFlow(ctx context.Context, req StartFlowRequest, medicalProfessionalID uuid.UUID) (*StartFlowResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Verify assessment exists and belongs to medical professional
	assessment, err := uc.assessmentRepo.GetByID(ctx, req.AssessmentID, medicalProfessionalID)
	if err != nil {
		return nil, err
	}

	// Start flow in rule engine
	flow, err := uc.ruleEngine.StartAssessmentFlow(req.AssessmentID, req.TreeID)
	if err != nil {
		return nil, err
	}

	// Save initial flow state
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

	// Update assessment status
	assessment.Status = domain.StatusInProgress
	if err := uc.assessmentRepo.Update(ctx, assessment); err != nil {
		return nil, fmt.Errorf("failed to update assessment status: %w", err)
	}

	// Get current question
	currentQuestion, err := uc.ruleEngine.GetCurrentQuestion(flow)
	if err != nil {
		return nil, err
	}

	return &StartFlowResponse{
		SessionID:   medicalProfessionalAnswer.ID,
		Question:    currentQuestion,
		IsComplete:  flow.Status == ruleenginedomain.FlowStatusCompleted,
		CurrentNode: flow.CurrentNode,
	}, nil
}

func (uc *RuleEngineUsecase) SubmitAnswer(ctx context.Context, req SubmitAnswerRequest, medicalProfessionalID uuid.UUID) (*SubmitAnswerResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Verify assessment exists and belongs to medical professional
	assessment, err := uc.assessmentRepo.GetByID(ctx, req.AssessmentID, medicalProfessionalID)
	if err != nil {
		return nil, err
	}

	// Get existing flow
	medicalProfessionalAnswer, err := uc.medicalProfessionalAnswerRepo.GetByAssessmentID(ctx, req.AssessmentID)
	if err != nil {
		return nil, domain.ErrMedicalProfessionalAnswerNotFound
	}

	// Reconstruct flow from saved data
	flow := &ruleenginedomain.AssessmentFlow{
		AssessmentID: req.AssessmentID,
		Answers:      map[string]interface{}(medicalProfessionalAnswer.Answers),
		Status:       ruleenginedomain.FlowStatusInProgress,
		CreatedAt:    medicalProfessionalAnswer.CreatedAt,
		UpdatedAt:    medicalProfessionalAnswer.UpdatedAt,
	}

	// Set current node from request or use stored one
	if flow.CurrentNode == "" {
		flow.CurrentNode = "check_birth_asphyxia"
	}

	// Submit answer to rule engine
	updatedFlow, nextQuestion, err := uc.ruleEngine.SubmitAnswer(flow, req.NodeID, req.Answer)
	if err != nil {
		return nil, err
	}

	// Update stored flow
	medicalProfessionalAnswer.Answers = domain.JSONB(updatedFlow.Answers)
	medicalProfessionalAnswer.UpdatedAt = time.Now()

	if err := uc.medicalProfessionalAnswerRepo.Upsert(ctx, medicalProfessionalAnswer); err != nil {
		return nil, fmt.Errorf("failed to update assessment flow: %w", err)
	}

	// If flow is complete, save classification and update assessment
	if updatedFlow.Status == ruleenginedomain.FlowStatusCompleted || updatedFlow.Status == ruleenginedomain.FlowStatusEmergency {
		if err := uc.saveClassificationResults(ctx, assessment, updatedFlow.Classification); err != nil {
			return nil, fmt.Errorf("failed to save classification results: %w", err)
		}

		assessment.Status = domain.StatusCompleted
		if err := uc.assessmentRepo.Update(ctx, assessment); err != nil {
			return nil, fmt.Errorf("failed to update assessment status: %w", err)
		}
	}

	return &SubmitAnswerResponse{
		SessionID:      medicalProfessionalAnswer.ID,
		Question:       nextQuestion,
		Classification: updatedFlow.Classification,
		IsComplete:     updatedFlow.Status != ruleenginedomain.FlowStatusInProgress,
		CurrentNode:    updatedFlow.CurrentNode,
		Status:         updatedFlow.Status,
	}, nil
}

func (uc *RuleEngineUsecase) saveClassificationResults(ctx context.Context, assessment *domain.Assessment, classification *ruleenginedomain.ClassificationResult) error {
	if classification == nil {
		return nil
	}

	// Save classification
	class := &domain.Classification{
		ID:                    uuid.New(),
		AssessmentID:          assessment.ID,
		Disease:               classification.Classification,
		Color:                 classification.Color,
		Details:               classification.TreatmentPlan,
		RuleVersion:           "birth_asphyxia_v1",
		IsCriticalIllness:     classification.Emergency,
		RequiresUrgentReferral: classification.Emergency,
		TreatmentPriority:     1,
		CreatedAt:             time.Now(),
	}

	if err := uc.classificationRepo.Create(ctx, class); err != nil {
		return err
	}

	// Save counseling/advice
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

	// Save follow-up as additional counseling
	followUpCounseling := &domain.Counseling{
		ID:           uuid.New(),
		AssessmentID: assessment.ID,
		AdviceType:   "follow_up_schedule",
		Details:      fmt.Sprintf("Follow-up schedule: %v", classification.FollowUp),
		Language:     "en",
		CreatedAt:    time.Now(),
	}

	return uc.counselingRepo.Create(ctx, followUpCounseling)
}

func (uc *RuleEngineUsecase) GetAssessmentTree(treeID string) (*ruleenginedomain.AssessmentTree, error) {
	return uc.ruleEngine.GetAssessmentTree(treeID)
}