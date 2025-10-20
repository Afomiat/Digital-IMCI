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
	engine.RegisterAssessmentTree(GetChildFeverTree())
	engine.RegisterAssessmentTree(GetChildEarProblemTree())
	engine.RegisterAssessmentTree(GetChildAnemiaTree())
	engine.RegisterAssessmentTree(GetAcuteMalnutritionTree())
	engine.RegisterAssessmentTree(GetFeedingAssessmentTree())
	engine.RegisterAssessmentTree(GetChildHIVAssessmentTree())

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
		case "child_fever":
			finalClassification = re.classifyFever(flow.Answers)
		case "acute_malnutrition":
			finalClassification = re.classifyAcuteMalnutrition(flow.Answers)
		case "feeding_assessment":
			finalClassification = re.classifyFeedingAssessment(flow.Answers)
		case "hiv_assessment": // NEW HIV CASE
			finalClassification = re.classifyHIVAssessment(flow.Answers)
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
	case "child_fever":
		finalClassification = re.classifyFever(answers)
	case "child_ear_problem":
		finalClassification = re.classifyEarProblem(answers)
	case "child_anemia_check":
		finalClassification = re.classifyAnemia(answers)
	case "acute_malnutrition":
		finalClassification = re.classifyAcuteMalnutrition(answers)
	case "feeding_assessment":
		finalClassification = re.classifyFeedingAssessment(answers)
	case "hiv_assessment": 
		finalClassification = re.classifyHIVAssessment(answers)
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

	duration := 0
	if dur, ok := diarrheaDuration.(int); ok {
		duration = dur
	} else if durStr, ok := diarrheaDuration.(string); ok {
		if dur, err := strconv.Atoi(durStr); err == nil {
			duration = dur
		}
	}

	if bloodInStool == "yes" {
		return "DYSENTERY"
	}

	if duration >= 14 {
		if lethargicUnconscious == "yes" || restlessIrritable == "yes" || sunkenEyes == "yes" ||
			drinkingAbility == "no" || drinkingThirsty == "yes" || skinPinchVerySlow == "yes" || skinPinchSlow == "yes" {
			return "SEVERE_PERSISTENT_DIARRHEA"
		}
		return "PERSISTENT_DIARRHEA"
	}

	if (lethargicUnconscious == "yes" && sunkenEyes == "yes") ||
		(lethargicUnconscious == "yes" && drinkingAbility == "no") ||
		(lethargicUnconscious == "yes" && skinPinchVerySlow == "yes") ||
		(sunkenEyes == "yes" && drinkingAbility == "no") ||
		(sunkenEyes == "yes" && skinPinchVerySlow == "yes") ||
		(drinkingAbility == "no" && skinPinchVerySlow == "yes") {
		return "SEVERE_DEHYDRATION"
	}

	if (restlessIrritable == "yes" && sunkenEyes == "yes") ||
		(restlessIrritable == "yes" && drinkingThirsty == "yes") ||
		(restlessIrritable == "yes" && skinPinchSlow == "yes") ||
		(sunkenEyes == "yes" && drinkingThirsty == "yes") ||
		(sunkenEyes == "yes" && skinPinchSlow == "yes") ||
		(drinkingThirsty == "yes" && skinPinchSlow == "yes") {
		return "SOME_DEHYDRATION"
	}

	return "NO_DEHYDRATION"
}

func (re *ChildRuleEngine) classifyFever(answers map[string]interface{}) string {
	malariaRisk := answers["malaria_risk"]
	bloodFilmResult := answers["blood_film_result"]
	stiffNeck := answers["stiff_neck"]
	bulgingFontanelle := answers["bulging_fontanelle"]
	cloudingCornea := answers["clouding_cornea"]
	mouthUlcersSeverity := answers["mouth_ulcers_severity"]
	pusDrainingEye := answers["eye_pus"]
	measlesNow := answers["current_measles"]
	measlesHistory := answers["measles_history"]
	generalDangerSign := answers["any_general_danger_sign"]

	if generalDangerSign == "yes" || stiffNeck == "yes" || bulgingFontanelle == "yes" {
		return "VERY_SEVERE_FEBRILE_DISEASE"
	}

	if measlesNow == "yes" || measlesHistory == "yes" {
		if cloudingCornea == "yes" || mouthUlcersSeverity == "deep_extensive" {
			return "SEVERE_COMPLICATED_MEASLES"
		}
		if pusDrainingEye == "yes" || mouthUlcersSeverity == "not_deep_extensive" {
			return "MEASLES_WITH_EYE_MOUTH_COMPLICATIONS"
		}
		return "MEASLES_NO_COMPLICATIONS"
	}

	if malariaRisk == "high" {
		if bloodFilmResult == "positive" || bloodFilmResult == "not_available" {
			return "MALARIA_HIGH_RISK"
		}
	} else if malariaRisk == "low" {
		if bloodFilmResult == "positive" {
			return "MALARIA_LOW_RISK"
		}
	}

	if malariaRisk == "no" || bloodFilmResult == "negative" {
		return "FEVER_NO_MALARIA"
	}

	return "FEVER_NO_MALARIA"
}

func (re *ChildRuleEngine) classifyEarProblem(answers map[string]interface{}) string {
	earPain := answers["ear_pain"]
	pusDraining := answers["pus_draining"]
	dischargeDuration := answers["discharge_duration"]
	tenderSwelling := answers["tender_swelling"]

	if tenderSwelling == "yes" {
		return "MASTOIDITIS"
	}

	if earPain == "yes" {
		return "ACUTE_EAR_INFECTION"
	}

	if pusDraining == "yes" {
		if dischargeDuration == "less_than_14_days" {
			return "ACUTE_EAR_INFECTION"
		} else if dischargeDuration == "14_days_or_more" {
			return "CHRONIC_EAR_INFECTION"
		}
	}

	if earPain == "no" && pusDraining == "no" {
		return "NO_EAR_INFECTION"
	}

	return "NO_EAR_INFECTION"
}

func (re *ChildRuleEngine) classifyAnemia(answers map[string]interface{}) string {
	palmarPallorPresent := answers["palmar_pallor_present"]
	palmarPallorSeverity := answers["palmar_pallor_severity"]
	hbValue := answers["hb_value"]
	hctValue := answers["hct_value"]
	classifyByPallor := answers["classify_by_pallor_only"]

	if hbValue != nil {
		if hb, ok := hbValue.(float64); ok {
			if hb < 7.0 {
				return "SEVERE_ANEMIA"
			} else if hb >= 7.0 && hb < 11.0 {
				return "ANEMIA"
			} else if hb >= 11.0 {
				return "NO_ANEMIA"
			}
		}
	}

	if hctValue != nil {
		if hct, ok := hctValue.(float64); ok {
			if hct < 21.0 {
				return "SEVERE_ANEMIA"
			} else if hct >= 21.0 && hct < 33.0 {
				return "ANEMIA"
			} else if hct >= 33.0 {
				return "NO_ANEMIA"
			}
		}
	}

	if classifyByPallor != nil {
		if classifyByPallor == "severe_pallor" {
			return "SEVERE_ANEMIA"
		} else if classifyByPallor == "some_pallor" {
			return "ANEMIA"
		}
	}

	if palmarPallorPresent == "no" || palmarPallorSeverity == "no_palmar_pallor" {
		return "NO_ANEMIA"
	}

	return "NO_ANEMIA"
}

func (re *ChildRuleEngine) classifyAcuteMalnutrition(answers map[string]interface{}) string {

	oedema := answers["pitting_edema"]
	wfl := re.parseFloat(answers["wfl_z_score"])
	muac := re.parseFloat(answers["muac_measurement"])
	severeWastingWithEdema := answers["severe_wasting_with_edema_check"]
	complicationsAny := answers["medical_complications_multi"]
	appetite := answers["appetite_test"]

	oedemaAny := (oedema == "plus" || oedema == "plus_plus" || oedema == "plus_plus_plus")
	oedemaSevere := (oedema == "plus_plus_plus")
	severeByWfl := wfl < -3.0
	severeByMuac := muac < 11.5
	moderateByWfl := (wfl >= -3.0 && wfl < -2.0)
	moderateByMuac := (muac >= 11.5 && muac < 12.5)

	if oedemaSevere {
		return "COMPLICATED_SEVERE_ACUTE_MALNUTRITION"
	}

	if severeWastingWithEdema == "yes" {
		return "COMPLICATED_SEVERE_ACUTE_MALNUTRITION"
	}

	if complicationsAny == "any_present" {
		return "COMPLICATED_SEVERE_ACUTE_MALNUTRITION"
	}

	if oedemaAny || severeByWfl || severeByMuac {
		if appetite == "failed" {
			return "COMPLICATED_SEVERE_ACUTE_MALNUTRITION"
		}
		if appetite == "passed" {
			return "UNCOMPLICATED_SEVERE_ACUTE_MALNUTRITION"
		}
	}

	if (moderateByWfl || moderateByMuac) && !oedemaAny {
		return "MODERATE_ACUTE_MALNUTRITION"
	}

	return "NO_ACUTE_MALNUTRITION"
}
func (re *ChildRuleEngine) classifyFeedingAssessment(answers map[string]interface{}) string {
	breastfeeding := answers["breastfeeding_check"]
	breastfeedingFreq := re.parseInt(answers["breastfeeding_frequency"])
	nightBreastfeeding := answers["night_breastfeeding"]
	otherFood := answers["other_food_check"]
	otherFoodTypes := answers["other_food_types"]
	foodQuantity := answers["food_quantity"]
	foodFreq := re.parseInt(answers["food_frequency"])
	feedingMethod := answers["feeding_method"]
	replacementMilk := answers["replacement_milk_check"]
	replacementMilkType := answers["replacement_milk_type"]
	replacementFreq := re.parseInt(answers["replacement_frequency"])
	_ = answers["replacement_quantity"] 
	milkPreparation := answers["milk_preparation"]
	utensilCleaning := answers["utensil_cleaning"]
	mamChild := answers["mam_specific_check"]
	servingSize := answers["serving_size"]
	ownServing := answers["own_serving"]
	feedingPerson := answers["feeding_person"]
	feedingChanged := answers["feeding_changes"]

	feedingProblemSigns := []bool{}

	if breastfeeding == "yes" {
		if breastfeedingFreq < 6 {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}
		if nightBreastfeeding == "no" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}
	}

	if otherFood == "yes" {
		otherFoodStr := fmt.Sprintf("%v", otherFoodTypes)
		quantityStr := fmt.Sprintf("%v", foodQuantity)

		if strings.Contains(strings.ToLower(otherFoodStr), "milk") && (quantityStr == "small" || quantityStr == "varies") {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}

		if foodFreq < 3 {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}

		if feedingMethod == "bottle" || feedingMethod == "both" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}
	}

	if replacementMilk == "yes" {
		milkTypeStr := fmt.Sprintf("%v", replacementMilkType)
		preparationStr := fmt.Sprintf("%v", milkPreparation)
		cleaningStr := fmt.Sprintf("%v", utensilCleaning)

		if milkTypeStr == "condensed_milk" || milkTypeStr == "evaporated_milk" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}

		if replacementFreq < 6 {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}

		if preparationStr == "diluted_with_water" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}

		if cleaningStr == "not_cleaned_properly" || cleaningStr == "washed_with_water_only" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}
	}

	if mamChild == "yes" {
		servingSizeStr := fmt.Sprintf("%v", servingSize)
		feedingPersonStr := fmt.Sprintf("%v", feedingPerson)

		if servingSizeStr == "small" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}

		if ownServing == "no" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}

		if feedingPersonStr == "child_feeds_self" {
			feedingProblemSigns = append(feedingProblemSigns, true)
		}
	}

	if feedingChanged == "yes" {
		feedingProblemSigns = append(feedingProblemSigns, true)
	}

	if len(feedingProblemSigns) > 0 {
		return "FEEDING_PROBLEM"
	}

	return "NO_FEEDING_PROBLEM"
}

func (re *ChildRuleEngine) classifyHIVAssessment(answers map[string]interface{}) string {
	getValue := func(key string) string {
		if val, exists := answers[key]; exists {
			return fmt.Sprintf("%v", val)
		}
		return ""
	}

	motherStatus := getValue("mother_hiv_status")
	childAntibody := getValue("child_antibody_test")
	childDNAPCR := getValue("child_dna_pcr_test")
	clinicalSigns := answers["clinical_signs_check"]
	breastfeeding := getValue("child_breastfeeding")
	breastfedLast6Weeks := getValue("breastfed_last_6weeks")

	if childDNAPCR == "positive" {
		return "HIV_INFECTED_DNA_PCR"
	}

	if childAntibody == "positive" {
		if childDNAPCR == "unknown" {
			if clinicalSigns != nil {
				if signs, ok := clinicalSigns.([]interface{}); ok && len(signs) >= 2 {
					return "PRESUMPTIVE_SEVERE_HIV"
				}
				if signStr, ok := clinicalSigns.(string); ok && signStr != "" {
					return "PRESUMPTIVE_SEVERE_HIV"
				}
			}
		}
		return "HIV_INFECTED_ANTIBODY"
	}

	if motherStatus == "positive" {
		childTestNegativeOrUnknown := (childAntibody == "negative" || childAntibody == "unknown" || 
									 childDNAPCR == "negative" || childDNAPCR == "unknown")
		
		if childTestNegativeOrUnknown && breastfeeding == "yes" {
			return "HIV_EXPOSED"
		}
		
		if childTestNegativeOrUnknown && breastfeeding == "no" && breastfedLast6Weeks == "yes" {
			return "HIV_EXPOSED"
		}
	}

	if motherStatus == "unknown" && (childAntibody == "unknown" && childDNAPCR == "unknown") {
		return "HIV_STATUS_UNKNOWN"
	}

	if motherStatus == "negative" {
		return "HIV_INFECTION_UNLIKELY"
	}

	if motherStatus == "positive" && 
	   (childDNAPCR == "negative" || (childAntibody == "negative" && childDNAPCR == "unknown")) && 
	   breastfeeding == "no" && 
	   breastfedLast6Weeks == "no" {
		return "HIV_INFECTION_UNLIKELY"
	}

	if motherStatus == "unknown" && (childAntibody == "negative" || childDNAPCR == "negative") {
		return "HIV_INFECTION_UNLIKELY"
	}

	return "HIV_STATUS_UNKNOWN"
}

func (re *ChildRuleEngine) parseInt(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case float64:
		return int(t)
	case string:
		i, err := strconv.Atoi(t)
		if err == nil {
			return i
		}
		return 0
	default:
		return 0
	}
}

func (re *ChildRuleEngine) parseFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err == nil {
			return f
		}
		return 0
	default:
		return 0
	}
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
