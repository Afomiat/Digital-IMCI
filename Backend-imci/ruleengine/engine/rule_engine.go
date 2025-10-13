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
	engine.RegisterAssessmentTree(GetFeedingProblemUnderweightTree())
	engine.RegisterAssessmentTree(GetReplacementFeedingTree())
	engine.RegisterAssessmentTree(GetHIVAssessmentTree())
	engine.RegisterAssessmentTree(GetGestationClassificationTree())



	
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
			fmt.Printf("DEBUG: Diarrhea classification called, result: %s\n", finalClassification)
		case "feeding_problem_underweight_check": 
			finalClassification = re.classifyFeedingProblem(flow.Answers)
			fmt.Printf("DEBUG: Feeding problem classification called, result: %s\n", finalClassification)
		case "replacement_feeding_check":
			finalClassification = re.classifyReplacementFeeding(flow.Answers)
			fmt.Printf("DEBUG: Replacement feeding classification called, result: %s\n", finalClassification)
		default:
			finalClassification = "SEVERE_INFECTION_UNLIKELY"
		}
		
		fmt.Printf("DEBUG: TreeID=%s, Classification=%s\n", flow.TreeID, finalClassification) 
		
		outcome, exists := tree.Outcomes[finalClassification]
    	fmt.Printf("DEBUG: Outcome exists=%v\n", exists) 
    
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
	


func (re *RuleEngine) ProcessBatchAssessment(assessmentID uuid.UUID, treeID string, answers map[string]interface{}) (*domain.AssessmentFlow, error) {
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
	case "very_severe_disease_check":
		finalClassification = re.classifyVerySevereDisease(answers)
	case "jaundice_check":
		finalClassification = re.classifyJaundice(answers)
	case "diarrhea_check":
		finalClassification = re.classifyDehydration(answers)
	case "feeding_problem_underweight_check":
		finalClassification = re.classifyFeedingProblem(answers)
	case "replacement_feeding_check":
		finalClassification = re.classifyReplacementFeeding(answers)
	case "hiv_status_assessment": 
		finalClassification = re.ClassifyHIV(answers)
	case "birth_asphyxia_check":
        finalClassification = re.classifyBirthAsphyxia(answers)
	case "gestation_classification": 
		finalClassification = re.classifyGestation(answers)
	default:
		finalClassification = "SEVERE_INFECTION_UNLIKELY"
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
func (re *RuleEngine) classifyBirthAsphyxia(answers map[string]interface{}) string {
    checkAsphyxia := answers["check_birth_asphyxia"]
    notBreathing := answers["not_breathing"]
    gasping := answers["gasping"]
    breathingPoorly := answers["breathing_poorly"]
    breathingNormally := answers["breathing_normally"]

    if checkAsphyxia == "no" {
        return "NO_BIRTH_ASPHYXIA"
    }

    if notBreathing == "yes" {
        return "BIRTH_ASPHYXIA"
    }
    if gasping == "yes" {
        return "BIRTH_ASPHYXIA"
    }
    if breathingPoorly == "yes" {
        return "BIRTH_ASPHYXIA"
    }
    if breathingNormally == "no" {
        return "BIRTH_ASPHYXIA"
    }

    return "NO_BIRTH_ASPHYXIA"
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

func (re *RuleEngine) classifyFeedingProblem(answers map[string]interface{}) string {
	breastfeedingStatus := answers["breastfeeding_status"]
	breastfeedingFrequency, _ := answers["breastfeeding_frequency"].(float64)
	emptyBreast := answers["empty_breast_before_switching"]
	increaseDuringIllness := answers["increase_frequency_illness"]
	otherFoodsDrinks := answers["other_foods_drinks"]
	positioning := answers["observe_positioning"]
	attachment := answers["observe_attachment"]
	suckling := answers["observe_suckling"]
	weightAge, _ := answers["weight_age_assessment"].(float64)
	thrush := answers["oral_thrush_check"]

	hasFeedingProblem := false

	if breastfeedingStatus == "no" {
		hasFeedingProblem = true
	}

	if breastfeedingStatus == "yes" {
		if breastfeedingFrequency < 8 {
			hasFeedingProblem = true
		}

		if positioning == "no" || attachment == "no" {
			hasFeedingProblem = true
		}

		if suckling == "no" {
			hasFeedingProblem = true
		}

		if emptyBreast == "no" || increaseDuringIllness == "no" || otherFoodsDrinks == "yes" {
			hasFeedingProblem = true
		}
	}

	isUnderweight := weightAge < -2
	hasThrush := thrush == "yes"

	if hasFeedingProblem || isUnderweight || hasThrush {
		return "FEEDING_PROBLEM_OR_UNDERWEIGHT"
	}

	return "NO_FEEDING_PROBLEM_NOT_UNDERWEIGHT"
}

func (re *RuleEngine) classifyReplacementFeeding(answers map[string]interface{}) string {
	milkType := answers["milk_type"]
	preparationMethod := answers["preparation_method_non_bf"]
	breastMilkGiven := answers["breast_milk_given_non_bf"]
	additionalFoods := answers["additional_foods_fluids_non_bf"]
	feedingMethod := answers["feeding_method_non_bf"]
	utensilCleaning := answers["utensil_cleaning_non_bf"]
	weightAge, _ := answers["weight_age_assessment_non_bf"].(float64)
	thrush := answers["oral_thrush_check_non_bf"]
	
	amountPerFeed, _ := answers["amount_per_feed_non_bf"].(float64)

	hasFeedingProblem := false

	if milkType == "animal_milk" {
		hasFeedingProblem = true
	}

	if preparationMethod == "incorrect_unhygienic" {
		hasFeedingProblem = true
	}

	if breastMilkGiven == "yes" {
		hasFeedingProblem = true
	}

	if additionalFoods == "inappropriate_foods" {
		hasFeedingProblem = true
	}

	if feedingMethod == "bottle" {
		hasFeedingProblem = true
	}

	if utensilCleaning == "improper_cleaning" {
		hasFeedingProblem = true
	}
	
	if amountPerFeed < 80 {
		hasFeedingProblem = true
	}

	isUnderweight := weightAge < -2
	hasThrush := thrush == "yes"

	if hasFeedingProblem || isUnderweight || hasThrush {
		return "FEEDING_PROBLEM_OR_UNDERWEIGHT"
	}

	return "NO_FEEDING_PROBLEM_NOT_UNDERWEIGHT"
}

func (re *RuleEngine) ClassifyHIV(answers map[string]interface{}) string {
	motherStatus := answers["mother_hiv_status"]
	antibodyStatus := answers["infant_antibody_status"]
	pcrStatus := answers["infant_dna_pcr_status"]
	breastfeeding := answers["breastfeeding_status"]

	if pcrStatus == "positive" {
		return "HIV_INFECTED"
	}

	if motherStatus == "positive" {
		if pcrStatus == "unknown" {
			return "HIV_EXPOSED"
		}
		if pcrStatus == "negative" {
			if breastfeeding == "yes" {
				return "HIV_EXPOSED"
			} else {
				return "HIV_INFECTION_UNLIKELY"
			}
		}
	}

	if motherStatus == "unknown" || motherStatus == "negative" {
		if antibodyStatus == "negative" {
			return "HIV_INFECTION_UNLIKELY"
		}
		if motherStatus == "unknown" && antibodyStatus == "unknown" {
			return "HIV_STATUS_UNKNOWN"
		}
	}

	return "HIV_INFECTION_UNLIKELY"
}

func (re *RuleEngine) classifyGestation(answers map[string]interface{}) string {
	gestationalAge := re.getNumericValue(answers["gestational_age_weeks"])
	birthWeight := re.getNumericValue(answers["birth_weight_grams"])
	currentWeight := re.getNumericValue(answers["current_weight_grams"])
	birthWeightGA := re.getNumericValue(answers["birth_weight_grams_ga"])
	currentWeightGA := re.getNumericValue(answers["current_weight_grams_ga"])
	
	knowGestationalAge := answers["know_gestational_age"]
	knowBirthWeight := answers["know_birth_weight"]
	knowBirthWeightGA := answers["know_birth_weight_ga"]
	canWeighBaby := answers["can_weigh_baby"]
	canWeighBabyGA := answers["can_weigh_baby_ga"]

	var effectiveWeight float64
	
	if currentWeightGA > 0 {
		effectiveWeight = currentWeightGA
	} else if birthWeightGA > 0 {
		effectiveWeight = birthWeightGA
	} else if currentWeight > 0 {
		effectiveWeight = currentWeight
	} else if birthWeight > 0 {
		effectiveWeight = birthWeight
	}

	if effectiveWeight > 0 && effectiveWeight < 1500 {
		return "VERY_LOW_BIRTH_WEIGHT"
	}

	if knowGestationalAge == "yes" && gestationalAge > 0 {
		if gestationalAge < 32 {
			return "VERY_LOW_BIRTH_WEIGHT"
		} else if gestationalAge >= 32 && gestationalAge < 37 {
			return "LOW_BIRTH_WEIGHT"
		} else if gestationalAge >= 37 {
			return "NORMAL_BIRTH_WEIGHT"
		}
	}

	if knowGestationalAge == "no" && effectiveWeight > 0 {
		if effectiveWeight < 1500 {
			return "VERY_LOW_BIRTH_WEIGHT"
		} else if effectiveWeight >= 1500 && effectiveWeight < 2500 {
			return "LOW_BIRTH_WEIGHT"
		} else if effectiveWeight >= 2500 {
			return "NORMAL_BIRTH_WEIGHT"
		}
	}

	if (knowBirthWeight == "no" && canWeighBaby == "no") || (knowBirthWeightGA == "no" && canWeighBabyGA == "no") {
		return "WEIGHT_UNKNOWN"
	}

	return "LOW_BIRTH_WEIGHT"
}

func (re *RuleEngine) getNumericValue(value interface{}) float64 {
	if value == nil {
		return 0
	}
	
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		var result float64
		_, err := fmt.Sscanf(v, "%f", &result)
		if err == nil {
			return result
		}
		return 0
	default:
		return 0
	}
}
func (re *RuleEngine) formatAnswer(question *domain.Question, answer interface{}) string {
	switch question.QuestionType {
	case "number_input":
		return "value_based"
	default:
		return fmt.Sprintf("%v", answer)
	}
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