// ruleengine/engine/rule_engine.go
package engine

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/google/uuid"
)

var (
	ErrAssessmentFlowNotFound = errors.New("assessment flow not found")
	ErrInvalidAnswer         = errors.New("invalid answer for question")
	ErrQuestionNotFound      = errors.New("question not found")
	ErrFlowAlreadyCompleted  = errors.New("assessment flow already completed")
)

type RuleEngine struct {
	trees map[string]*domain.AssessmentTree
}

func NewRuleEngine() (*RuleEngine, error) {
	engine := &RuleEngine{
		trees: make(map[string]*domain.AssessmentTree),
	}
	
	// Register assessment trees
	engine.RegisterAssessmentTree(GetBirthAsphyxiaTree())
	
	return engine, nil
}

func (re *RuleEngine) RegisterAssessmentTree(tree *domain.AssessmentTree) {
	re.trees[tree.AssessmentID] = tree
}

func (re *RuleEngine) GetAssessmentTree(assessmentID string) (*domain.AssessmentTree, error) {
	tree, exists := re.trees[assessmentID]
	if !exists {
		return nil, fmt.Errorf("assessment tree not found: %s", assessmentID)
	}
	return tree, nil
}

func (re *RuleEngine) StartAssessmentFlow(assessmentID uuid.UUID, treeID string) (*domain.AssessmentFlow, error) {
	tree, err := re.GetAssessmentTree(treeID)
	if err != nil {
		return nil, err
	}

	flow := &domain.AssessmentFlow{
		AssessmentID: assessmentID,
		CurrentNode:  tree.StartNode,
		Status:       domain.FlowStatusInProgress,
		Answers:      make(map[string]interface{}),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return flow, nil
}

func (re *RuleEngine) SubmitAnswer(flow *domain.AssessmentFlow, nodeID string, answer interface{}) (*domain.AssessmentFlow, *domain.Question, error) {
	if flow.Status == domain.FlowStatusCompleted {
		return nil, nil, ErrFlowAlreadyCompleted
	}

	tree, err := re.GetAssessmentTree("birth_asphyxia_check")
	if err != nil {
		return nil, nil, err
	}

	// Find the current question
	question, err := re.findQuestion(tree, nodeID)
	if err != nil {
		return nil, nil, err
	}

	// Validate answer
	answerStr := fmt.Sprintf("%v", answer)
	if _, valid := question.Answers[answerStr]; !valid {
		return nil, nil, ErrInvalidAnswer
	}

	// Store answer
	flow.Answers[nodeID] = answer
	flow.UpdatedAt = time.Now()

	// Process answer
	answerConfig := question.Answers[answerStr]
	
	// Check if this leads to classification
	if answerConfig.Classification != "" {
		outcome, exists := tree.Outcomes[answerConfig.Classification]
		if exists {
			flow.Classification = &domain.ClassificationResult{
				Classification: outcome.Classification,
				Color:          outcome.Color,
				Emergency:      outcome.Emergency,
				Actions:        outcome.Actions,
				TreatmentPlan:  outcome.TreatmentPlan,
				FollowUp:       outcome.FollowUp,
				MotherAdvice:   outcome.MotherAdvice,
			}
			flow.Status = domain.FlowStatusCompleted
			if outcome.Emergency {
				flow.Status = domain.FlowStatusEmergency
			}
			return flow, nil, nil
		}
	}

	// Move to next node
	if answerConfig.NextNode != "" {
		flow.CurrentNode = answerConfig.NextNode
		nextQuestion, _ := re.findQuestion(tree, answerConfig.NextNode)
		return flow, nextQuestion, nil
	}

	// If no next node and no classification, complete flow
	flow.Status = domain.FlowStatusCompleted
	return flow, nil, nil
}

func (re *RuleEngine) GetCurrentQuestion(flow *domain.AssessmentFlow) (*domain.Question, error) {
	if flow.Status != domain.FlowStatusInProgress {
		return nil, nil
	}

	tree, err := re.GetAssessmentTree("birth_asphyxia_check")
	if err != nil {
		return nil, err
	}

	return re.findQuestion(tree, flow.CurrentNode)
}

func (re *RuleEngine) findQuestion(tree *domain.AssessmentTree, nodeID string) (*domain.Question, error) {
	for _, question := range tree.QuestionsFlow {
		if question.NodeID == nodeID {
			return &question, nil
		}
	}
	return nil, ErrQuestionNotFound
}

func (re *RuleEngine) ShouldShowQuestion(flow *domain.AssessmentFlow, question domain.Question) bool {
	if question.ShowCondition == "" {
		return true
	}

	// Simple condition parser - you might want to use a more sophisticated one
	conditions := strings.Split(question.ShowCondition, " AND ")
	for _, condition := range conditions {
		parts := strings.Split(condition, ".")
		if len(parts) != 2 {
			continue
		}

		nodeID := parts[0]
		expectedAnswer := parts[1]

		actualAnswer, exists := flow.Answers[nodeID]
		if !exists {
			return false
		}

		if fmt.Sprintf("%v", actualAnswer) != expectedAnswer {
			return false
		}
	}

	return true
}