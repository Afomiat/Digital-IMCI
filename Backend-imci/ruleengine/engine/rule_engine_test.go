package engine

import (
	"testing"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleEngine_BirthAsphyxiaClassification(t *testing.T) {
	engine, err := NewRuleEngine()
	require.NoError(t, err)

	tests := []struct {
		name           string
		answers        map[string]interface{}
		expectedClass  string
		expectedStatus domain.FlowStatus
	}{
		{
			name: "No birth asphyxia - immediate no",
			answers: map[string]interface{}{
				"check_birth_asphyxia": "no",
			},
			expectedClass:  "NO_BIRTH_ASPHYXIA", // Use underscore to match tree definition
			expectedStatus: domain.FlowStatusCompleted,
		},
		{
			name: "Birth asphyxia - not breathing",
			answers: map[string]interface{}{
				"check_birth_asphyxia": "yes",
				"not_breathing":        "yes",
			},
			expectedClass:  "BIRTH_ASPHYXIA", // Use underscore to match tree definition
			expectedStatus: domain.FlowStatusEmergency,
		},
		{
			name: "Birth asphyxia - gasping",
			answers: map[string]interface{}{
				"check_birth_asphyxia": "yes",
				"not_breathing":        "no",
				"gasping":              "yes",
			},
			expectedClass:  "BIRTH_ASPHYXIA",
			expectedStatus: domain.FlowStatusEmergency,
		},
		{
			name: "Birth asphyxia - breathing poorly",
			answers: map[string]interface{}{
				"check_birth_asphyxia": "yes",
				"not_breathing":        "no",
				"gasping":              "no",
				"breathing_poorly":     "yes",
			},
			expectedClass:  "BIRTH_ASPHYXIA",
			expectedStatus: domain.FlowStatusEmergency,
		},
		{
			name: "Birth asphyxia - not breathing normally",
			answers: map[string]interface{}{
				"check_birth_asphyxia": "yes",
				"not_breathing":        "no",
				"gasping":              "no",
				"breathing_poorly":     "no",
				"breathing_normally":   "no",
			},
			expectedClass:  "BIRTH_ASPHYXIA",
			expectedStatus: domain.FlowStatusEmergency,
		},
		{
			name: "No birth asphyxia - breathing normally",
			answers: map[string]interface{}{
				"check_birth_asphyxia": "yes",
				"not_breathing":        "no",
				"gasping":              "no",
				"breathing_poorly":     "no",
				"breathing_normally":   "yes",
			},
			expectedClass:  "NO_BIRTH_ASPHYXIA",
			expectedStatus: domain.FlowStatusCompleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assessmentID := uuid.New()
			flow, err := engine.ProcessBatchAssessment(assessmentID, "birth_asphyxia_check", tt.answers)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, flow.Status)
			require.NotNil(t, flow.Classification)
			assert.Equal(t, tt.expectedClass, flow.Classification.Classification)
			
			// Verify outcome details match the tree definition
			tree, err := engine.GetAssessmentTree("birth_asphyxia_check")
			require.NoError(t, err)
			expectedOutcome := tree.Outcomes[tt.expectedClass]
			
			// Now compare the actual outcome details
			assert.Equal(t, expectedOutcome.Classification, flow.Classification.Classification)
			assert.Equal(t, expectedOutcome.Color, flow.Classification.Color)
			assert.Equal(t, expectedOutcome.Emergency, flow.Classification.Emergency)
			assert.Equal(t, expectedOutcome.Actions, flow.Classification.Actions)
			assert.Equal(t, expectedOutcome.TreatmentPlan, flow.Classification.TreatmentPlan)
		})
	}
}

func TestRuleEngine_BirthAsphyxiaSequentialFlow(t *testing.T) {
	engine, err := NewRuleEngine()
	require.NoError(t, err)

	assessmentID := uuid.New()

	// Start assessment
	flow, err := engine.StartAssessmentFlow(assessmentID, "birth_asphyxia_check")
	require.NoError(t, err)
	assert.Equal(t, "check_birth_asphyxia", flow.CurrentNode)
	assert.Equal(t, domain.FlowStatusInProgress, flow.Status)

	// Get first question
	question, err := engine.GetCurrentQuestion(flow)
	require.NoError(t, err)
	assert.Equal(t, "check_birth_asphyxia", question.NodeID)

	// Submit "yes" to first question
	updatedFlow, nextQuestion, err := engine.SubmitAnswer(flow, "check_birth_asphyxia", "yes")
	require.NoError(t, err)
	assert.Equal(t, "not_breathing", updatedFlow.CurrentNode)
	require.NotNil(t, nextQuestion)
	assert.Equal(t, "not_breathing", nextQuestion.NodeID)

	// Submit "yes" to not breathing - should complete with BIRTH_ASPHYXIA
	updatedFlow, nextQuestion, err = engine.SubmitAnswer(updatedFlow, "not_breathing", "yes")
	require.NoError(t, err)
	assert.Nil(t, nextQuestion) // No next question - flow completed
	assert.Equal(t, domain.FlowStatusEmergency, updatedFlow.Status)
	require.NotNil(t, updatedFlow.Classification)
	assert.Equal(t, "BIRTH_ASPHYXIA", updatedFlow.Classification.Classification)
}

func TestRuleEngine_BirthAsphyxiaInvalidAnswers(t *testing.T) {
	engine, err := NewRuleEngine()
	require.NoError(t, err)

	assessmentID := uuid.New()
	flow, err := engine.StartAssessmentFlow(assessmentID, "birth_asphyxia_check")
	require.NoError(t, err)

	// Test invalid answer type
	_, _, err = engine.SubmitAnswer(flow, "check_birth_asphyxia", "invalid_answer")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAnswer, err)

	// Test invalid node ID
	_, _, err = engine.SubmitAnswer(flow, "invalid_node", "yes")
	assert.Error(t, err)
	assert.Equal(t, ErrQuestionNotFound, err)
}

func TestRuleEngine_GetBirthAsphyxiaTree(t *testing.T) {
	engine, err := NewRuleEngine()
	require.NoError(t, err)

	tree, err := engine.GetAssessmentTree("birth_asphyxia_check")
	require.NoError(t, err)

	assert.Equal(t, "birth_asphyxia_check", tree.AssessmentID)
	assert.Equal(t, "Check for Birth Asphyxia", tree.Title)
	assert.Len(t, tree.QuestionsFlow, 5)
	assert.Len(t, tree.Outcomes, 2)

	// Verify the actual outcome keys used in the tree
	outcomeKeys := make([]string, 0, len(tree.Outcomes))
	for key := range tree.Outcomes {
		outcomeKeys = append(outcomeKeys, key)
	}
	t.Logf("Actual outcome keys in tree: %v", outcomeKeys)

	// Test non-existent tree
	_, err = engine.GetAssessmentTree("non_existent_tree")
	assert.Error(t, err)
	// Just check that it's an error, don't check the exact error message
	assert.Contains(t, err.Error(), "assessment tree not found")
}
