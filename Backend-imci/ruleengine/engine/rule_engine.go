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
	ErrTreeNotFound          = errors.New("assessment tree not found")
)

type RuleEngine struct {
	trees map[string]*domain.AssessmentTree
}

func NewRuleEngine() (*RuleEngine, error) {
	engine := &RuleEngine{
		trees: make(map[string]*domain.AssessmentTree),
	}
	
	engine.RegisterAssessmentTree(GetBirthAsphyxiaTree())
	engine.RegisterAssessmentTree(GetVerySevereDiseaseTree())
	engine.RegisterAssessmentTree(GetJaundiceTree())
	engine.RegisterAssessmentTree(GetDiarrheaTree())
	
	return engine, nil
}

func (re *RuleEngine) RegisterAssessmentTree(tree *domain.AssessmentTree) {
	re.trees[tree.AssessmentID] = tree
}

func (re *RuleEngine) GetAssessmentTree(assessmentID string) (*domain.AssessmentTree, error) {
	tree, exists := re.trees[assessmentID]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrTreeNotFound, assessmentID)
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
		TreeID:       treeID, 
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

	tree, err := re.GetAssessmentTree(flow.TreeID)
	if err != nil {
		return nil, nil, err
	}

	question, err := re.findQuestion(tree, nodeID)
	if err != nil {
		return nil, nil, err
	}

	answerStr := re.formatAnswer(question, answer)
	if _, valid := question.Answers[answerStr]; !valid {
		return nil, nil, ErrInvalidAnswer
	}

	flow.Answers[nodeID] = answer
	flow.UpdatedAt = time.Now()

	answerConfig := question.Answers[answerStr]
	
	if answerConfig.Classification == "AUTO_CLASSIFY" {
		var finalClassification string
		
		switch flow.TreeID {
		case "very_severe_disease_check":
			finalClassification = re.classifyVerySevereDisease(flow.Answers)
		case "jaundice_check":
			finalClassification = re.classifyJaundice(flow.Answers) 
		case "diarrhea_check":
			finalClassification = re.classifyDehydration(flow.Answers)
			fmt.Printf("DEBUG: Diarrhea classification called, result: %s\n", finalClassification) // ADD THIS
		default:
			finalClassification = "SEVERE_INFECTION_UNLIKELY"
		}
		
		fmt.Printf("DEBUG: TreeID=%s, Classification=%s\n", flow.TreeID, finalClassification) // ADD THIS
		
		outcome, exists := tree.Outcomes[finalClassification]
    	fmt.Printf("DEBUG: Outcome exists=%v\n", exists) // ADD THIS
    
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
	
	if answerConfig.Classification != "" && answerConfig.Classification != "AUTO_CLASSIFY" {
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

	if answerConfig.NextNode != "" {
		flow.CurrentNode = answerConfig.NextNode
		nextQuestion, _ := re.findQuestion(tree, answerConfig.NextNode)
		return flow, nextQuestion, nil
	}

	flow.Status = domain.FlowStatusCompleted
	return flow, nil, nil
}
	


func (re *RuleEngine) GetCurrentQuestion(flow *domain.AssessmentFlow) (*domain.Question, error) {
	if flow.Status != domain.FlowStatusInProgress {
		return nil, nil
	}

	tree, err := re.GetAssessmentTree(flow.TreeID)
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
func (re *RuleEngine) classifyVerySevereDisease(answers map[string]interface{}) string {
	feedingAbility := answers["feeding_ability_detail"]
	convulsions := answers["convulsions_history"]
	movements := answers["check_movements"]
	breathingRate, _ := answers["breathing_rate"].(float64)
	chestIndrawing := answers["chest_indrawing"]
	umbilicus := answers["umbilicus_check"]
	skinPustules := answers["skin_pustules"]
	temperature, _ := answers["temperature_measurement"].(float64)

	if movements == "no_movement_even_stimulated" {
		return "CRITICAL_ILLNESS"
	}
	if feedingAbility == "unable_to_feed" {
		return "CRITICAL_ILLNESS" 
	}
	if convulsions == "yes" {
		return "CRITICAL_ILLNESS"
	}

	if feedingAbility == "not_feeding_well" {
		return "VERY_SEVERE_DISEASE"
	}
	if movements == "moves_only_when_stimulated" {
		return "VERY_SEVERE_DISEASE"
	}
	if chestIndrawing == "yes" {
		return "VERY_SEVERE_DISEASE"
	}
	if temperature >= 37.5 || temperature < 35.5 {
		return "VERY_SEVERE_DISEASE"
	}
	if breathingRate >= 60 {
		return "VERY_SEVERE_DISEASE"
	}

	if umbilicus == "yes" || skinPustules == "yes" {
		return "LOCAL_BACTERIAL_INFECTION"
	}

	return "SEVERE_INFECTION_UNLIKELY"
}

func (re *RuleEngine) classifyJaundice(answers map[string]interface{}) string {
	skinYellow := answers["skin_yellow"]
	palmsSolesYellow := answers["palms_soles_yellow"]
	age, _ := answers["infant_age"].(float64)
	
	hasJaundice := skinYellow == "yes"
	hasSevereSigns := palmsSolesYellow == "yes"
	
	if !hasJaundice {
		return "NO_JAUNDICE"
	}
	
	if hasSevereSigns {
		return "SEVERE_JAUNDICE_URGENT"
	}
	
	if age < 1 { 
		return "SEVERE_JAUNDICE_URGENT"
	}
	
	if age >= 14 {
		return "SEVERE_JAUNDICE_URGENT"
	}
	
	if age >= 1 && age < 14 { 
		return "JAUNDICE"
	}
	
	return "NO_JAUNDICE"
}

func (re *RuleEngine) classifyDehydration(answers map[string]interface{}) string {
    movementCondition := answers["movement_condition"]
    skinPinch := answers["skin_pinch"]
    otherSevere := answers["assess_other_severe"]
    
    var severity string
    if movementCondition == "no_movement_even_when_stimulated" || skinPinch == "very_slowly_more_than_2_seconds" {
        severity = "SEVERE_DEHYDRATION"
    } else if skinPinch == "slowly" {
        severity = "SOME_DEHYDRATION"
    } else {
        return "NO_DEHYDRATION"
    }
    
    if otherSevere == "yes" {
        if severity == "SEVERE_DEHYDRATION" {
            return "SEVERE_DEHYDRATION_WITH_OTHER_SEVERE"
        } else if severity == "SOME_DEHYDRATION" {
            return "SOME_DEHYDRATION_WITH_OTHER_SEVERE"
        }
    } else {
        if severity == "SEVERE_DEHYDRATION" {
            return "SEVERE_DEHYDRATION_ALONE"
        } else if severity == "SOME_DEHYDRATION" {
            return "SOME_DEHYDRATION_ALONE"
        }
    }
    
    return "NO_DEHYDRATION"
}

func (re *RuleEngine) formatAnswer(question *domain.Question, answer interface{}) string {
	switch question.QuestionType {
	case "number_input":
		return "value_based"
	default:
		return fmt.Sprintf("%v", answer)
	}
}

func stringToFloat(s string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	return result, err
}

func (re *RuleEngine) ShouldShowQuestion(flow *domain.AssessmentFlow, question domain.Question) bool {
	if question.ShowCondition == "" {
		return true
	}

	tree, err := re.GetAssessmentTree(flow.TreeID)
	if err != nil {
		return false
	}

	return re.evaluateCondition(flow, tree, question.ShowCondition)
}

func (re *RuleEngine) evaluateCondition(flow *domain.AssessmentFlow, tree *domain.AssessmentTree, condition string) bool {
	conditions := strings.Split(condition, " AND ")
	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		if !re.evaluateSingleCondition(flow, tree, condition) {
			return false
		}
	}
	return true
}

func (re *RuleEngine) evaluateSingleCondition(flow *domain.AssessmentFlow, tree *domain.AssessmentTree, condition string) bool {
	parts := strings.Split(condition, ".")
	if len(parts) != 2 {
		return false
	}

	nodeID := strings.TrimSpace(parts[0])
	expectedAnswer := strings.TrimSpace(parts[1])

	actualAnswer, exists := flow.Answers[nodeID]
	if !exists {
		return false
	}

	return fmt.Sprintf("%v", actualAnswer) == expectedAnswer
}

func (re *RuleEngine) GetAvailableTrees() []string {
	treeIDs := make([]string, 0, len(re.trees))
	for treeID := range re.trees {
		treeIDs = append(treeIDs, treeID)
	}
	return treeIDs
}