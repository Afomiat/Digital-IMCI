package engine

import (
	"testing"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBirthAsphyxiaTree(t *testing.T) {
	tree := GetBirthAsphyxiaTree()

	// Basic structure validation
	assert.Equal(t, "birth_asphyxia_check", tree.AssessmentID)
	assert.Equal(t, "Check for Birth Asphyxia", tree.Title)
	assert.Contains(t, tree.Instructions, "Golden Minute")
	assert.Equal(t, "check_birth_asphyxia", tree.StartNode)
	
	// Validate questions flow
	assert.Len(t, tree.QuestionsFlow, 5)
	
	// Validate outcomes
	assert.Len(t, tree.Outcomes, 2)
	assert.Contains(t, tree.Outcomes, "BIRTH_ASPHYXIA")
	assert.Contains(t, tree.Outcomes, "NO_BIRTH_ASPHYXIA")
}

func TestBirthAsphyxiaTree_QuestionStructure(t *testing.T) {
	tree := GetBirthAsphyxiaTree()

	tests := []struct {
		nodeID       string
		question     string
		questionType string
		required     bool
		level        int
	}{
		{
			nodeID:       "check_birth_asphyxia",
			question:     "Check for Birth Asphyxia?",
			questionType: "yes_no",
			required:     true,
			level:        1,
		},
		{
			nodeID:       "not_breathing",
			question:     "Is baby not breathing?",
			questionType: "yes_no",
			required:     true,
			level:        2,
		},
		// Add other questions...
	}

	for _, tt := range tests {
		t.Run(tt.nodeID, func(t *testing.T) {
			question := findQuestionByNodeID(tree, tt.nodeID)
			require.NotNil(t, question, "Question %s should exist", tt.nodeID)
			assert.Equal(t, tt.question, question.Question)
			assert.Equal(t, tt.questionType, question.QuestionType)
			assert.Equal(t, tt.required, question.Required)
			assert.Equal(t, tt.level, question.Level)
		})
	}
}

// Helper function
func findQuestionByNodeID(tree *domain.AssessmentTree, nodeID string) *domain.Question {
	for _, question := range tree.QuestionsFlow {
		if question.NodeID == nodeID {
			return &question
		}
	}
	return nil
}