// ruleengine/engine/feeding_problem_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetFeedingProblemUnderweightTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "feeding_problem_underweight_check",
		Title:        "Assess Feeding Problems and Underweight",
		Instructions: "Assess infant feeding practices and nutritional status for breastfeeding infants.",
		StartNode:    "feeding_difficulty",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "feeding_difficulty",
				Question:     "Is there any difficulty of feeding?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Assess if the infant has any feeding difficulties",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "breastfeeding_status",
						Color:    "yellow",
					},
					"no": {
						NextNode: "breastfeeding_status",
					},
				},
			},
			{
				NodeID:       "breastfeeding_status",
				Question:     "Is the infant breastfed?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Determine if the infant is receiving breast milk",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "breastfeeding_frequency",
					},
					"no": {
						NextNode: "weight_age_assessment",
					},
				},
			},
			{
				NodeID:       "breastfeeding_frequency",
				Question:     "How many times in 24 hours does the infant breastfeed?",
				QuestionType: "number_input",
				Required:     true,
				Level:        2,
				ParentNode:   "breastfeeding_status",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "Count the number of breastfeeding sessions in 24 hours",
				Validation: &domain.Validation{
					Min:  0,
					Max:  24,
					Step: 1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "empty_breast_before_switching",
					},
				},
			},
			{
				NodeID:       "empty_breast_before_switching",
				Question:     "Do you empty one breast before switching to the other?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "breastfeeding_status",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "Determine if mother allows infant to fully empty one breast before switching",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "increase_frequency_illness",
					},
					"no": {
						NextNode: "increase_frequency_illness",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "increase_frequency_illness",
				Question:     "Do you increase frequency of breastfeeding during illness?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "breastfeeding_status",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "Assess feeding practices during illness",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "other_foods_drinks",
					},
					"no": {
						NextNode: "other_foods_drinks",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "other_foods_drinks",
				Question:     "Does the infant receive any other foods or drinks?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "breastfeeding_status",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "Assess if infant receives supplemental foods or fluids",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "assess_breastfeeding_observation",
						Color:    "yellow",
					},
					"no": {
						NextNode: "assess_breastfeeding_observation",
					},
				},
			},
			{
				NodeID:       "assess_breastfeeding_observation",
				Question:     "Has the infant breastfed in the previous hour?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "breastfeeding_status",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "If yes: Ask mother to wait and tell you when infant is willing to feed again. If no: Ask mother to put infant to breast and observe for 4 minutes.",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "observe_positioning",
					},
					"no": {
						NextNode: "observe_positioning",
					},
				},
			},
			{
				NodeID:       "observe_positioning",
				Question:     "Is the infant well positioned?",
				QuestionType: "yes_no", 
				Required:     true,
				Level:        3,
				ParentNode:   "assess_breastfeeding_observation",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "Observe: Infant's head and body straight, facing breast, body close to mother, whole body supported",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "observe_attachment",
					},
					"no": {
						NextNode: "observe_attachment",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "observe_attachment",
				Question:     "Is the infant able to attach well?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "observe_positioning",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "Observe: Chin touching breast, mouth wide open, lower lip turned outward, more areola above than below mouth",
				Answers: map[string]domain.Answer{
					"yes": { 
						NextNode: "observe_suckling",
					},
					"no": { 
						NextNode: "observe_suckling",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "observe_suckling",
				Question:     "Is the infant suckling effectively?",
				QuestionType: "yes_no", 
				Required:     true,
				Level:        3,
				ParentNode:   "observe_attachment",
				ShowCondition: "breastfeeding_status.yes",
				Instructions: "Observe for slow deep sucks with occasional pausing",
				Answers: map[string]domain.Answer{
					"yes": { 
						NextNode: "weight_age_assessment",
					},
					"no": { 
						NextNode: "weight_age_assessment",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "weight_age_assessment",
				Question:     "Determine weight for age (WFA) Z-score",
				QuestionType: "number_input",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Use WFA charts to determine Z-score. If visible severe wasting or edema present, classify as Severe Acute Malnutrition.",
				Validation: &domain.Validation{
					Min:  -5,
					Max:  5,
					Step: 0.1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "oral_thrush_check",
					},
				},
			},
			{
				NodeID:       "oral_thrush_check",
				Question:     "Look for ulcers or white patches in the mouth (thrush)",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Examine mouth for signs of oral thrush",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "FEEDING_PROBLEM_OR_UNDERWEIGHT",
						Color:          "yellow",
					},
					"no": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"FEEDING_PROBLEM_OR_UNDERWEIGHT": {
				Classification: "FEEDING PROBLEM OR UNDERWEIGHT",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Advise the mother to breastfeed as often and for as long as the infant wants, day and night",
					"If baby not sucking, show her how to express breast milk",
					"If not well positioned, attached or not suckling effectively, teach correct positioning and attachment",
					"If breastfeeding less than 8 times in 24 hours, advise to increase frequency of feeding",
					"Empty one breast completely before switching to the other",
					"Increase frequency of feeding during and after illness",
					"If receiving other foods or drinks, counsel mother on exclusive breastfeeding",
					"If not breastfeeding at all: Counsel on breastfeeding and relactation",
					"If no possibility of breastfeeding: Advise about correct preparation of breast milk substitutes and using a cup",
					"If thrush, teach the mother to treat thrush at home",
					"Advise mother to give home care for the young infant",
					"Ensure infant is tested for HIV",
				},
				TreatmentPlan: "Feeding counseling and management",
				FollowUp: []string{
					"Follow-up any feeding problem or thrush in 2 days",
					"Follow-up for underweight in 14 days",
				},
				MotherAdvice: "Advise on proper feeding techniques and when to return immediately",
				Notes:        "Not well positioned, not well attached, not suckling effectively, <8 breastfeeds/24h, switching breasts frequently, not increasing during illness, other foods/drinks, not breastfeeding at all, WFA < -2Z, or thrush",
			},
			"NO_FEEDING_PROBLEM_NOT_UNDERWEIGHT": {
				Classification: "NO FEEDING PROBLEM AND NOT UNDERWEIGHT",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Advise mother to give home care for the young infant",
					"Praise the mother for feeding the infant well",
				},
				TreatmentPlan: "Continue current feeding practices",
				FollowUp: []string{
					"Routine follow-up as scheduled",
				},
				MotherAdvice: "Continue good feeding practices",
				Notes:        "WFA ≥ -2Z and no signs of feeding problems",
			},
			"SEVERE_ACUTE_MALNUTRITION": {
				Classification: "SEVERE ACUTE MALNUTRITION",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Refer urgently to hospital for SAM management",
					"Use sick child acute malnutrition assessment",
				},
				TreatmentPlan: "Hospital management for severe acute malnutrition",
				FollowUp: []string{
					"Immediate hospital referral",
				},
				MotherAdvice: "Go to hospital immediately - this is an emergency",
				Notes:        "Visible severe wasting or edema present",
			},
			"AUTO_CLASSIFY": {
				Classification: "AUTO_CLASSIFY",
				Color:          "white",
				Emergency:      false,
				Actions: []string{
					"Review all assessment findings for classification",
				},
				TreatmentPlan: "Further assessment needed",
				FollowUp: []string{
					"Determine classification based on collected signs",
				},
				MotherAdvice: "Wait for complete assessment",
				Notes:        "Automatic classification based on collected data",
			},
		},
	}
}