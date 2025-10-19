// ruleengine/engine/child_engine.go
package engine

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/google/uuid"
)

type ChildRuleEngine struct {
	trees map[string]*domain.AssessmentTree
}

func NewChildRuleEngine() (*ChildRuleEngine, error) {
	engine := &ChildRuleEngine{
		trees: make(map[string]*domain.AssessmentTree),
	}
	
	engine.RegisterAssessmentTree(GetChildGeneralDangerSignsTree())
	engine.RegisterAssessmentTree(GetChildCoughDifficultBreathingTree())
	engine.RegisterAssessmentTree(GetChildDiarrheaTree()) 
	
	
	return engine, nil
}

func (re *ChildRuleEngine) RegisterAssessmentTree(tree *domain.AssessmentTree) {
	re.trees[tree.AssessmentID] = tree
}

func (re *ChildRuleEngine) GetAssessmentTree(assessmentID string) (*domain.AssessmentTree, error) {
	tree, exists := re.trees[assessmentID]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrTreeNotFound, assessmentID)
	}
	return tree, nil
}

func (re *ChildRuleEngine) StartAssessmentFlow(assessmentID uuid.UUID, treeID string) (*domain.AssessmentFlow, error) {
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

func (re *ChildRuleEngine) SubmitAnswer(flow *domain.AssessmentFlow, nodeID string, answer interface{}) (*domain.AssessmentFlow, *domain.Question, error) {
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
		case "child_general_danger_signs":
			finalClassification = re.classifyChildGeneralDangerSigns(flow.Answers)
		case "child_cough_difficult_breathing":
			finalClassification = re.classifyChildCoughDifficultBreathing(flow.Answers)
		case "child_diarrhea":
			finalClassification = re.classifyChildDiarrhea(flow.Answers)
		default:
			finalClassification = "NO_DANGER_SIGNS"
		}
		
		outcome, exists := tree.Outcomes[finalClassification]
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

func (re *ChildRuleEngine) ProcessBatchAssessment(assessmentID uuid.UUID, treeID string, answers map[string]interface{}) (*domain.AssessmentFlow, error) {
	tree, err := re.GetAssessmentTree(treeID)
	if err != nil {
		return nil, err
	}

	flow := &domain.AssessmentFlow{
		AssessmentID: assessmentID,
		TreeID:       treeID,
		CurrentNode:  "",
		Status:       domain.FlowStatusInProgress,
		Answers:      answers,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	var finalClassification string
	
	switch treeID {
	case "child_general_danger_signs":
		finalClassification = re.classifyChildGeneralDangerSigns(answers)
	case "child_cough_difficult_breathing":
		finalClassification = re.classifyChildCoughDifficultBreathing(answers)
	case "child_diarrhea":
		finalClassification = re.classifyChildDiarrhea(answers)
	default:
		finalClassification = "NO_DANGER_SIGNS"
		
	}

	outcome, exists := tree.Outcomes[finalClassification]
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
	}

	return flow, nil
}

func (re *ChildRuleEngine) GetCurrentQuestion(flow *domain.AssessmentFlow) (*domain.Question, error) {
	if flow.Status != domain.FlowStatusInProgress {
		return nil, nil
	}

	tree, err := re.GetAssessmentTree(flow.TreeID)
	if err != nil {
		return nil, err
	}

	return re.findQuestion(tree, flow.CurrentNode)
}

func (re *ChildRuleEngine) findQuestion(tree *domain.AssessmentTree, nodeID string) (*domain.Question, error) {
	for _, question := range tree.QuestionsFlow {
		if question.NodeID == nodeID {
			return &question, nil
		}
	}
	return nil, ErrQuestionNotFound
}

func (re *ChildRuleEngine) classifyChildGeneralDangerSigns(answers map[string]interface{}) string {
	unableToDrink := answers["unable_to_drink_breastfeed"]
	vomitsEverything := answers["vomits_everything"]
	convulsionsHistory := answers["convulsions_history"]
	lethargicUnconscious := answers["lethargic_unconscious"]
	convulsingNow := answers["convulsing_now"]

	if unableToDrink == "no" || 
	   vomitsEverything == "yes" ||
	   convulsionsHistory == "yes" ||
	   lethargicUnconscious == "yes" ||
	   convulsingNow == "yes" {
		return "VERY_SEVERE_DISEASE"
	}

	return "NO_GENERAL_DANGER_SIGNS"
}

func (re *ChildRuleEngine) classifyChildCoughDifficultBreathing(answers map[string]interface{}) string {
	generalDangerSigns := answers["general_danger_signs"]
	stridor := answers["stridor"]
	oxygenSaturation := answers["oxygen_saturation"]
	
	if generalDangerSigns == "yes" || stridor == "yes" || oxygenSaturation == "yes" {
		return "SEVERE_PNEUMONIA_OR_VERY_SEVERE_DISEASE"
	}
	
	chestIndrawing := answers["chest_indrawing"]
	hivExposed := answers["hiv_exposed"]
	
	if chestIndrawing == "yes" && hivExposed == "yes" {
		return "CHEST_INDRAWING_HIV_EXPOSED"
	}
	
	fastBreathing := answers["fast_breathing"]
	wheezing := answers["wheezing"]
	
	if fastBreathing == "yes" || chestIndrawing == "yes" {
		if wheezing == "yes" {
			return "PNEUMONIA_WITH_WHEEZING"
		}
		return "PNEUMONIA"
	}
	
	if wheezing == "yes" {
		return "COUGH_OR_COLD_WITH_WHEEZING"
	}
	
	return "COUGH_OR_COLD"
}

func (re *ChildRuleEngine) classifyChildDiarrhea(answers map[string]interface{}) string {
	diarrheaDuration := answers["how_long_diarrhea"]
	bloodInStool := answers["blood_in_stool"]
	lethargicUnconscious := answers["lethargic_unconscious"]
	restlessIrritable := answers["restless_irritable"]
	sunkenEyes := answers["sunken_eyes"]
	drinkingAbility := answers["drinking_ability"]
	drinkingThirsty := answers["drinking_thirsty"]
	skinPinchVerySlow := answers["skin_pinch"]
	skinPinchSlow := answers["skin_pinch_slow"]
	
	// Convert duration
	duration := 0
	if dur, ok := diarrheaDuration.(int); ok {
		duration = dur
	} else if durStr, ok := diarrheaDuration.(string); ok {
		if dur, err := strconv.Atoi(durStr); err == nil {
			duration = dur
		}
	}
	
	// DYSENTERY first (blood in stool)
	if bloodInStool == "yes" {
		return "DYSENTERY"
	}
	
	// PERSISTENT DIARRHOEA (14+ days)
	if duration >= 14 {
		// Check for any dehydration signs
		if lethargicUnconscious == "yes" || restlessIrritable == "yes" || sunkenEyes == "yes" || 
		   drinkingAbility == "no" || drinkingThirsty == "yes" || skinPinchVerySlow == "yes" || skinPinchSlow == "yes" {
			return "SEVERE_PERSISTENT_DIARRHEA"
		}
		return "PERSISTENT_DIARRHEA"
	}
	
	// SEVERE DEHYDRATION - need 2+ signs from: lethargic, sunken eyes, not able to drink, skin very slow
	if (lethargicUnconscious == "yes" && sunkenEyes == "yes") ||
	   (lethargicUnconscious == "yes" && drinkingAbility == "no") ||
	   (lethargicUnconscious == "yes" && skinPinchVerySlow == "yes") ||
	   (sunkenEyes == "yes" && drinkingAbility == "no") ||
	   (sunkenEyes == "yes" && skinPinchVerySlow == "yes") ||
	   (drinkingAbility == "no" && skinPinchVerySlow == "yes") {
		return "SEVERE_DEHYDRATION"
	}
	
	// SOME DEHYDRATION - need 2+ signs from: restless, sunken eyes, drinking eagerly, skin slow
	if (restlessIrritable == "yes" && sunkenEyes == "yes") ||
	   (restlessIrritable == "yes" && drinkingThirsty == "yes") ||
	   (restlessIrritable == "yes" && skinPinchSlow == "yes") ||
	   (sunkenEyes == "yes" && drinkingThirsty == "yes") ||
	   (sunkenEyes == "yes" && skinPinchSlow == "yes") ||
	   (drinkingThirsty == "yes" && skinPinchSlow == "yes") {
		return "SOME_DEHYDRATION"
	}
	
	// NO DEHYDRATION
	return "NO_DEHYDRATION"
}

func (re *ChildRuleEngine) formatAnswer(question *domain.Question, answer interface{}) string {
	switch question.QuestionType {
	case "number_input":
		return "value_based"
	default:
		return fmt.Sprintf("%v", answer)
	}
}

func (re *ChildRuleEngine) ShouldShowQuestion(flow *domain.AssessmentFlow, question domain.Question) bool {
	if question.ShowCondition == "" {
		return true
	}

	tree, err := re.GetAssessmentTree(flow.TreeID)
	if err != nil {
		return false
	}

	return re.evaluateCondition(flow, tree, question.ShowCondition)
}

func (re *ChildRuleEngine) evaluateCondition(flow *domain.AssessmentFlow, tree *domain.AssessmentTree, condition string) bool {
	conditions := strings.Split(condition, " AND ")
	for _, condition := range conditions {
		condition = strings.TrimSpace(condition)
		if !re.evaluateSingleCondition(flow, tree, condition) {
			return false
		}
	}
	return true
}

func (re *ChildRuleEngine) evaluateSingleCondition(flow *domain.AssessmentFlow, tree *domain.AssessmentTree, condition string) bool {
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

func (re *ChildRuleEngine) GetAvailableTrees() []string {
	treeIDs := make([]string, 0, len(re.trees))
	for treeID := range re.trees {
		treeIDs = append(treeIDs, treeID)
	}
	return treeIDs
}

func (re *ChildRuleEngine) GetTreeQuestions(treeID string) (*domain.AssessmentTree, error) {
	return re.GetAssessmentTree(treeID)
}