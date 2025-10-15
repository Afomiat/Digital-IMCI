package contract_tests

import (
	"testing"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/Afomiat/Digital-IMCI/ruleengine/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBirthAsphyxiaTreeContract(t *testing.T) {
	tree := engine.GetBirthAsphyxiaTree()

	// Contract: Tree must have valid structure
	t.Run("Tree structure contract", func(t *testing.T) {
		assert.NotEmpty(t, tree.AssessmentID)
		assert.NotEmpty(t, tree.Title)
		assert.NotEmpty(t, tree.Instructions)
		assert.NotEmpty(t, tree.StartNode)
		assert.NotEmpty(t, tree.QuestionsFlow)
		assert.NotEmpty(t, tree.Outcomes)
	})

	// Contract: Start node must exist in questions
	t.Run("Start node exists contract", func(t *testing.T) {
		startNodeExists := false
		for _, question := range tree.QuestionsFlow {
			if question.NodeID == tree.StartNode {
				startNodeExists = true
				break
			}
		}
		assert.True(t, startNodeExists, "Start node must exist in questions flow")
	})

	// Contract: All next nodes must exist
	t.Run("Next nodes validity contract", func(t *testing.T) {
		nodes := make(map[string]bool)
		for _, question := range tree.QuestionsFlow {
			nodes[question.NodeID] = true
		}

		for _, question := range tree.QuestionsFlow {
			for _, answer := range question.Answers {
				if answer.NextNode != "" && answer.NextNode != "assessment_complete_normal" {
					assert.True(t, nodes[answer.NextNode], "Next node %s must exist", answer.NextNode)
				}
			}
		}
	})

	// Contract: All classifications must have outcomes
	t.Run("Classifications have outcomes contract", func(t *testing.T) {
		usedClassifications := make(map[string]bool)
		
		// Collect classifications from answers
		for _, question := range tree.QuestionsFlow {
			for _, answer := range question.Answers {
				if answer.Classification != "" {
					// Your tree uses underscores: "BIRTH_ASPHYXIA" and "NO_BIRTH_ASPHYXIA"
					usedClassifications[answer.Classification] = true
				}
			}
		}

		// Verify all used classifications have outcomes
		for classification := range usedClassifications {
			_, exists := tree.Outcomes[classification]
			assert.True(t, exists, "Classification %s must have an outcome defined", classification)
		}
	})

	// Contract: Emergency paths must have emergency outcomes
	t.Run("Emergency paths contract", func(t *testing.T) {
		for _, question := range tree.QuestionsFlow {
			for _, answer := range question.Answers {
				if answer.EmergencyPath {
					outcome, exists := tree.Outcomes[answer.Classification]
					require.True(t, exists, "Emergency path must lead to existing classification: %s", answer.Classification)
					assert.True(t, outcome.Emergency, "Emergency path must lead to emergency outcome: %s", answer.Classification)
				}
			}
		}
	})

	// Contract: Verify specific outcome keys used in your tree
	t.Run("Specific outcome keys contract", func(t *testing.T) {
		// These are the actual keys used in your birth_asphyxia_tree.go (with underscores)
		expectedOutcomeKeys := []string{"BIRTH_ASPHYXIA", "NO_BIRTH_ASPHYXIA"}
		
		for _, expectedKey := range expectedOutcomeKeys {
			_, exists := tree.Outcomes[expectedKey]
			assert.True(t, exists, "Expected outcome key %s must exist", expectedKey)
		}
	})
}

func TestDevelopmentalAssessmentTreeContract(t *testing.T) {
	tree := engine.GetDevelopmentalAssessmentTree()

	// Contract: Tree must have valid structure
	t.Run("Tree structure contract", func(t *testing.T) {
		assert.NotEmpty(t, tree.AssessmentID)
		assert.NotEmpty(t, tree.Title)
		assert.NotEmpty(t, tree.Instructions)
		assert.NotEmpty(t, tree.StartNode)
		assert.NotEmpty(t, tree.QuestionsFlow)
		assert.NotEmpty(t, tree.Outcomes)
	})

	// Contract: Start node must exist in questions
	t.Run("Start node exists contract", func(t *testing.T) {
		startNodeExists := false
		for _, question := range tree.QuestionsFlow {
			if question.NodeID == tree.StartNode {
				startNodeExists = true
				break
			}
		}
		assert.True(t, startNodeExists, "Start node must exist in questions flow")
	})

	// Contract: All next nodes must exist
	t.Run("Next nodes validity contract", func(t *testing.T) {
		nodes := make(map[string]bool)
		for _, question := range tree.QuestionsFlow {
			nodes[question.NodeID] = true
		}

		for _, question := range tree.QuestionsFlow {
			for _, answer := range question.Answers {
				if answer.NextNode != "" {
					assert.True(t, nodes[answer.NextNode], "Next node %s must exist", answer.NextNode)
				}
			}
		}
	})

	// Contract: All classifications must have outcomes
	t.Run("Classifications have outcomes contract", func(t *testing.T) {
		usedClassifications := make(map[string]bool)
		
		// Collect classifications from answers
		for _, question := range tree.QuestionsFlow {
			for _, answer := range question.Answers {
				if answer.Classification != "" && answer.Classification != "AUTO_CLASSIFY_MILESTONES" {
					usedClassifications[answer.Classification] = true
				}
			}
		}

		// Verify all used classifications have outcomes
		for classification := range usedClassifications {
			_, exists := tree.Outcomes[classification]
			assert.True(t, exists, "Classification %s must have an outcome defined", classification)
		}
	})

	// Contract: Emergency paths must have emergency outcomes
	t.Run("Emergency paths contract", func(t *testing.T) {
		for _, question := range tree.QuestionsFlow {
			for _, answer := range question.Answers {
				if answer.EmergencyPath {
					outcome, exists := tree.Outcomes[answer.Classification]
					require.True(t, exists, "Emergency path must lead to existing classification: %s", answer.Classification)
					assert.True(t, outcome.Emergency, "Emergency path must lead to emergency outcome: %s", answer.Classification)
				}
			}
		}
	})

	// Contract: Verify specific outcome keys used in your tree
	t.Run("Specific outcome keys contract", func(t *testing.T) {
		expectedOutcomeKeys := []string{
			"SEVERE_CLASSIFICATION_NO_ASSESSMENT", 
			"SUSPECTED_DEVELOPMENTAL_DELAY", 
			"NO_DEVELOPMENTAL_DELAY",
		}
		
		for _, expectedKey := range expectedOutcomeKeys {
			_, exists := tree.Outcomes[expectedKey]
			assert.True(t, exists, "Expected outcome key %s must exist", expectedKey)
		}
	})

	// Contract: Milestone assessment should have proper answer mappings
	t.Run("Milestone assessment contract", func(t *testing.T) {
		milestoneQuestion := findQuestionByNodeID(tree, "assess_milestones")
		require.NotNil(t, milestoneQuestion)
		
		// Should have the three expected answer options
		assert.Contains(t, milestoneQuestion.Answers, "all_achieved")
		assert.Contains(t, milestoneQuestion.Answers, "one_missing")
		assert.Contains(t, milestoneQuestion.Answers, "multiple_missing")
		
		// All should lead to classifications
		assert.Equal(t, "NO_DEVELOPMENTAL_DELAY", milestoneQuestion.Answers["all_achieved"].Classification)
		assert.Equal(t, "SUSPECTED_DEVELOPMENTAL_DELAY", milestoneQuestion.Answers["one_missing"].Classification)
		assert.Equal(t, "SUSPECTED_DEVELOPMENTAL_DELAY", milestoneQuestion.Answers["multiple_missing"].Classification)
	})
}

// Helper function to find a question by node ID
func findQuestionByNodeID(tree *domain.AssessmentTree, nodeID string) *domain.Question {
	for _, question := range tree.QuestionsFlow {
		if question.NodeID == nodeID {
			return &question
		}
	}
	return nil
}