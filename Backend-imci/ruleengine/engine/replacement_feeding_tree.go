// ruleengine/engine/replacement_feeding_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetReplacementFeedingTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "replacement_feeding_check",
		Title:        "Assess Replacement Feeding for HIV-Positive Mothers",
		Instructions: "Assess feeding practices for infants not receiving breast milk due to HIV-positive mother or other reasons",
		StartNode:    "feeding_difficulty_non_bf",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "feeding_difficulty_non_bf",
				Question:     "Is there any difficulty in feeding?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Assess if the infant has any feeding difficulties",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "milk_type",
						Color:    "yellow",
					},
					"no": {
						NextNode: "milk_type",
					},
				},
			},
			{
				NodeID:       "milk_type",
				Question:     "What milk are you giving?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Select the type of replacement milk being used",
				Answers: map[string]domain.Answer{
					"infant_formula": {
						NextNode: "feeding_frequency_non_bf",
					},
					"animal_milk": {
						NextNode: "feeding_frequency_non_bf",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "feeding_frequency_non_bf",
				Question:     "How many times during the day and night?",
				QuestionType: "number_input",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Count total feeds in 24 hours",
				Validation: &domain.Validation{
					Min:  0,
					Max:  24,
					Step: 1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "amount_per_feed_non_bf",
					},
				},
			},
			{
				NodeID:       "amount_per_feed_non_bf",
				Question:     "How much is given at each feed?  (ml)",
				QuestionType: "number_input",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Enter the amount in milliliters per feed",
				Validation: &domain.Validation{
					Min:  0,
					Max:  500,
					Step: 10,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "preparation_method_non_bf",
					},
				},
			},
			{
				NodeID:       "preparation_method_non_bf",
				Question:     "How are you preparing the milk?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Let the mother demonstrate or explain how a feed is prepared",
				Answers: map[string]domain.Answer{
					"correct_hygienic": {
						NextNode: "breast_milk_given_non_bf",
					},
					"incorrect_unhygienic": {
						NextNode: "breast_milk_given_non_bf",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "breast_milk_given_non_bf",
				Question:     "Are you giving any breast milk?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Assess if mixed feeding is occurring",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "additional_foods_fluids_non_bf",
						Color:    "yellow",
					},
					"no": {
						NextNode: "additional_foods_fluids_non_bf",
					},
				},
			},
			{
				NodeID:       "additional_foods_fluids_non_bf",
				Question:     "What foods or fluids in addition to replacement feeding are given?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Select any supplemental foods or fluids being given",
				Answers: map[string]domain.Answer{
					"none": {
						NextNode: "feeding_method_non_bf",
					},
					"inappropriate_foods": {
						NextNode: "feeding_method_non_bf",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "feeding_method_non_bf",
				Question:     "How is the milk being given?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Observe feeding method",
				Answers: map[string]domain.Answer{
					"cup": {
						NextNode: "utensil_cleaning_non_bf",
					},
					"bottle": {
						NextNode: "utensil_cleaning_non_bf",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "utensil_cleaning_non_bf",
				Question:     "How are you cleaning the utensils?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Assess hygiene practices for feeding equipment",
				Answers: map[string]domain.Answer{
					"proper_cleaning": {
						NextNode: "weight_age_assessment_non_bf",
					},
					"improper_cleaning": {
						NextNode: "weight_age_assessment_non_bf",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "weight_age_assessment_non_bf",
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
						NextNode: "oral_thrush_check_non_bf",
					},
				},
			},
			{
				NodeID:       "oral_thrush_check_non_bf",
				Question:     "Look for mouth ulcers or white patches in the mouth (oral thrush)",
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
					"Counsel on optimal replacement feeding",
					"Identify concerns of the mother and the family about feeding",
					"Help the mother gradually withdraw other foods or fluids",
					"If mother is using a bottle, teach cup feeding",
					"If thrush, teach the mother to treat thrush at home",
					"Advise mother how to feed and keep the young infant warm at home",
				},
				TreatmentPlan: "Replacement feeding counseling and management",
				FollowUp: []string{
					"Follow-up any feeding problem or thrush in 2 days",
					"Follow-up underweight in 14 days",
				},
				MotherAdvice: "Advise on proper replacement feeding techniques and when to return immediately",
				Notes:        "Milk incorrectly/unhygienically prepared, inappropriate replacement milk, insufficient feeds, mixed feeding, bottle use, WFA < -2Z, or thrush",
			},
			"NO_FEEDING_PROBLEM_NOT_UNDERWEIGHT": {
				Classification: "NO FEEDING PROBLEM AND NOT UNDERWEIGHT",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Advise mother to give home care for the young infant",
					"Praise the mother for feeding the infant well",
				},
				TreatmentPlan: "Continue current replacement feeding practices",
				FollowUp: []string{
					"Routine follow-up as scheduled",
				},
				MotherAdvice: "Continue good replacement feeding practices",
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