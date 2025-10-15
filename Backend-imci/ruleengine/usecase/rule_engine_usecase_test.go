package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	ruleenginedomain "github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/Afomiat/Digital-IMCI/ruleengine/engine"
	"github.com/Afomiat/Digital-IMCI/ruleengine/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRuleEngineUsecase_StartBirthAsphyxiaAssessment(t *testing.T) {
	// Create mock instances
	mockAssessmentRepo := &mocks.AssessmentRepository{}
	mockAnswerRepo := &mocks.MedicalProfessionalAnswerRepository{}
	mockClinicalFindingsRepo := &mocks.ClinicalFindingsRepository{}
	mockClassificationRepo := &mocks.ClassificationRepository{}
	mockTreatmentPlanRepo := &mocks.TreatmentPlanRepository{}
	mockCounselingRepo := &mocks.CounselingRepository{}

	// Create rule engine
	ruleEngine, err := engine.NewRuleEngine()
	require.NoError(t, err)

	// Create usecase instance
	usecase := NewRuleEngineUsecase(
		ruleEngine,
		mockAssessmentRepo,
		mockAnswerRepo,
		mockClinicalFindingsRepo,
		mockClassificationRepo,
		mockTreatmentPlanRepo,
		mockCounselingRepo,
		time.Second*30,
	)

	assessmentID := uuid.New()
	medicalProfessionalID := uuid.New()

	t.Run("Successfully start birth asphyxia assessment", func(t *testing.T) {
		// Setup mock expectations
		mockAssessment := &domain.Assessment{
			ID:     assessmentID,
			Status: domain.StatusInProgress,
		}

		mockAssessmentRepo.On("GetByID", mock.Anything, assessmentID, medicalProfessionalID).
			Return(mockAssessment, nil)
		mockAnswerRepo.On("Upsert", mock.Anything, mock.AnythingOfType("*domain.MedicalProfessionalAnswer")).
			Return(nil)
		mockAssessmentRepo.On("Update", mock.Anything, mockAssessment).
			Return(nil)

		// Execute
		req := ruleenginedomain.StartFlowRequest{
			AssessmentID: assessmentID,
			TreeID:       "birth_asphyxia_check",
		}

		response, err := usecase.StartAssessmentFlow(context.Background(), req, medicalProfessionalID)

		// Verify
		require.NoError(t, err)
		assert.NotNil(t, response.SessionID)
		assert.NotNil(t, response.Question)
		assert.Equal(t, "check_birth_asphyxia", response.Question.NodeID)
		assert.False(t, response.IsComplete)
	})
}