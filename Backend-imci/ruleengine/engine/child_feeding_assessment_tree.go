// ruleengine/engine/child_feeding_assessment_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetFeedingAssessmentTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "feeding_assessment",
		Title:        "Feeding Assessment for Children Under 2 Years",
		Instructions: "ASSESS: If child is < 2 years old, or has Anemia or MAM; AND Has no severe classification - Do feeding assessment. Ask about breastfeeding, other foods, and feeding practices.",
		StartNode:    "breastfeeding_check",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "breastfeeding_check",
				Question:     "Do you breastfeed your child?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "ASK: About current breastfeeding practices",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "breastfeeding_frequency",
					},
					"no": {
						NextNode: "other_food_check",
					},
				},
			},
			{
				NodeID:        "breastfeeding_frequency",
				Question:      "How many times in 24 hours do you breastfeed?",
				QuestionType:  "number_input",
				Required:      true,
				Level:         2,
				ParentNode:    "breastfeeding_check",
				ShowCondition: "breastfeeding_check.yes",
				Instructions:  "ASK: Count breastfeeding sessions in 24 hours",
				Validation: &domain.Validation{
					Min:  0,
					Max:  20,
					Step: 1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "night_breastfeeding",
					},
				},
			},
			{
				NodeID:        "night_breastfeeding",
				Question:      "Do you breastfeed during the night?",
				QuestionType:  "yes_no",
				Required:      true,
				Level:         3,
				ParentNode:    "breastfeeding_frequency",
				ShowCondition: "breastfeeding_check.yes",
				Instructions:  "ASK: About night breastfeeding practices",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "other_food_check",
					},
					"no": {
						NextNode: "other_food_check",
					},
				},
			},
			{
				NodeID:       "other_food_check",
				Question:     "Does the child take any other food or fluids?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "breastfeeding_check",
				Instructions: "ASK: About complementary feeding and other foods/fluids",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "other_food_types",
					},
					"no": {
						NextNode: "replacement_milk_check",
					},
				},
			},
			{
				NodeID:        "other_food_types",
				Question:      "What food or fluids does the child take?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         5,
				ParentNode:    "other_food_check",
				ShowCondition: "other_food_check.yes",
				Instructions:  "ASK: Select the main foods and fluids given to the child",
				Answers: map[string]domain.Answer{
					"porridge": {
						NextNode: "food_quantity",
					},
					"vegetables": {
						NextNode: "food_quantity",
					},
					"fruits": {
						NextNode: "food_quantity",
					},
					"meat_fish": {
						NextNode: "food_quantity",
					},
					"milk": {
						NextNode: "food_quantity",
					},
					"water": {
						NextNode: "food_quantity",
					},
					"other": {
						NextNode: "food_quantity",
					},
				},
			},
			{
				NodeID:        "food_quantity",
				Question:      "How much is given at each feed?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         6,
				ParentNode:    "other_food_types",
				ShowCondition: "other_food_check.yes",
				Instructions:  "ASK: About quantity of food/fluids per feeding",
				Answers: map[string]domain.Answer{
					"small": {
						NextNode: "food_frequency",
					},
					"medium": {
						NextNode: "food_frequency",
					},
					"large": {
						NextNode: "food_frequency",
					},
					"varies": {
						NextNode: "food_frequency",
					},
				},
			},
			{
				NodeID:        "food_frequency",
				Question:      "How many times in 24 hours is other food given?",
				QuestionType:  "number_input",
				Required:      true,
				Level:         7,
				ParentNode:    "food_quantity",
				ShowCondition: "other_food_check.yes",
				Instructions:  "ASK: Count complementary feeding sessions in 24 hours",
				Validation: &domain.Validation{
					Min:  0,
					Max:  10,
					Step: 1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "feeding_method",
					},
				},
			},
			{
				NodeID:        "feeding_method",
				Question:      "What do you use to feed the child?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         8,
				ParentNode:    "food_frequency",
				ShowCondition: "other_food_check.yes",
				Instructions:  "ASK: About feeding utensils and methods",
				Answers: map[string]domain.Answer{
					"cup": {
						NextNode: "replacement_milk_check",
					},
					"bottle": {
						NextNode: "replacement_milk_check",
					},
					"both": {
						NextNode: "replacement_milk_check",
					},
					"other": {
						NextNode: "replacement_milk_check",
					},
				},
			},
			{
				NodeID:       "replacement_milk_check",
				Question:     "Is the child on replacement milk?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        9,
				ParentNode:   "other_food_check",
				Instructions: "ASK: About replacement milk feeding",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "replacement_milk_type",
					},
					"no": {
						NextNode: "mam_specific_check",
					},
				},
			},
			{
				NodeID:        "replacement_milk_type",
				Question:      "What replacement milk are you giving?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         10,
				ParentNode:    "replacement_milk_check",
				ShowCondition: "replacement_milk_check.yes",
				Instructions:  "ASK: About type of replacement milk used",
				Answers: map[string]domain.Answer{
					"infant_formula": {
						NextNode: "replacement_frequency",
					},
					"cow_milk": {
						NextNode: "replacement_frequency",
					},
					"goat_milk": {
						NextNode: "replacement_frequency",
					},
					"condensed_milk": {
						NextNode: "replacement_frequency",
					},
					"evaporated_milk": {
						NextNode: "replacement_frequency",
					},
					"other": {
						NextNode: "replacement_frequency",
					},
				},
			},
			{
				NodeID:        "replacement_frequency",
				Question:      "How many times in 24 hours is replacement milk given?",
				QuestionType:  "number_input",
				Required:      true,
				Level:         11,
				ParentNode:    "replacement_milk_type",
				ShowCondition: "replacement_milk_check.yes",
				Instructions:  "ASK: Count replacement milk feeds in 24 hours",
				Validation: &domain.Validation{
					Min:  0,
					Max:  12,
					Step: 1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "replacement_quantity",
					},
				},
			},
			{
				NodeID:        "replacement_quantity",
				Question:      "How much replacement milk is given at each feed?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         12,
				ParentNode:    "replacement_frequency",
				ShowCondition: "replacement_milk_check.yes",
				Instructions:  "ASK: About quantity of replacement milk per feed",
				Answers: map[string]domain.Answer{
					"small": {
						NextNode: "milk_preparation",
					},
					"medium": {
						NextNode: "milk_preparation",
					},
					"large": {
						NextNode: "milk_preparation",
					},
					"varies": {
						NextNode: "milk_preparation",
					},
				},
			},
			{
				NodeID:        "milk_preparation",
				Question:      "How is the replacement milk prepared?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         13,
				ParentNode:    "replacement_quantity",
				ShowCondition: "replacement_milk_check.yes",
				Instructions:  "ASK: About milk preparation method",
				Answers: map[string]domain.Answer{
					"according_to_instructions": {
						NextNode: "utensil_cleaning",
					},
					"diluted_with_water": {
						NextNode: "utensil_cleaning",
					},
					"mixed_with_other_foods": {
						NextNode: "utensil_cleaning",
					},
					"other": {
						NextNode: "utensil_cleaning",
					},
				},
			},
			{
				NodeID:        "utensil_cleaning",
				Question:      "How are you cleaning the feeding utensils?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         14,
				ParentNode:    "milk_preparation",
				ShowCondition: "replacement_milk_check.yes",
				Instructions:  "ASK: About hygiene practices for feeding utensils",
				Answers: map[string]domain.Answer{
					"washed_with_soap": {
						NextNode: "mam_specific_check",
					},
					"washed_with_water_only": {
						NextNode: "mam_specific_check",
					},
					"not_cleaned_properly": {
						NextNode: "mam_specific_check",
					},
					"other": {
						NextNode: "mam_specific_check",
					},
				},
			},
			{
				NodeID:       "mam_specific_check",
				Question:     "Does the child have Moderate Acute Malnutrition (MAM)?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        15,
				ParentNode:   "replacement_milk_check",
				Instructions: "ASK: About MAM-specific feeding practices",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "serving_size",
					},
					"no": {
						NextNode: "feeding_changes",
					},
				},
			},
			{
				NodeID:        "serving_size",
				Question:      "How large are the servings given to the child?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         16,
				ParentNode:    "mam_specific_check",
				ShowCondition: "mam_specific_check.yes",
				Instructions:  "ASK: About serving sizes for MAM child",
				Answers: map[string]domain.Answer{
					"small": {
						NextNode: "own_serving",
					},
					"normal": {
						NextNode: "own_serving",
					},
					"large": {
						NextNode: "own_serving",
					},
					"varies": {
						NextNode: "own_serving",
					},
				},
			},
			{
				NodeID:        "own_serving",
				Question:      "Does the child receive his/her own serving?",
				QuestionType:  "yes_no",
				Required:      true,
				Level:         17,
				ParentNode:    "serving_size",
				ShowCondition: "mam_specific_check.yes",
				Instructions:  "ASK: About individual serving for MAM child",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "feeding_person",
					},
					"no": {
						NextNode: "feeding_person",
					},
				},
			},
			{
				NodeID:        "feeding_person",
				Question:      "Who feeds the child and how?",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         18,
				ParentNode:    "own_serving",
				ShowCondition: "mam_specific_check.yes",
				Instructions:  "ASK: About feeding person and method for MAM child",
				Answers: map[string]domain.Answer{
					"mother_feeds_child": {
						NextNode: "feeding_changes",
					},
					"caregiver_feeds_child": {
						NextNode: "feeding_changes",
					},
					"child_feeds_self": {
						NextNode: "feeding_changes",
					},
					"other": {
						NextNode: "feeding_changes",
					},
				},
			},
			{
				NodeID:       "feeding_changes",
				Question:     "During the illness, has the child's feeding changed?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        19,
				ParentNode:   "mam_specific_check",
				Instructions: "ASK: About feeding changes during illness",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "finalize_feeding_classification",
					},
					"no": {
						NextNode: "finalize_feeding_classification",
					},
				},
			},
			{
				NodeID:       "finalize_feeding_classification",
				Question:     "Finalize feeding assessment classification",
				QuestionType: "single_choice",
				Required:     true,
				Level:        20,
				ParentNode:   "feeding_changes",
				Instructions: "System will compute feeding classification based on all collected information",
				Answers: map[string]domain.Answer{
					"compute": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"FEEDING_PROBLEM": {
				Classification: "FEEDING PROBLEM",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Advise mother on appropriate age specific feeding recommendations",
					"Advise mother on recommendations about child's specific feeding problem",
					"Follow-up of feeding problem in 5 days",
				},
				TreatmentPlan: "Nutritional counseling and feeding support",
				FollowUp: []string{
					"Follow-up of feeding problem in 5 days",
					"Monitor feeding practices",
					"Provide ongoing support",
				},
				MotherAdvice: "Your child needs better feeding practices. Follow the feeding advice and return in 5 days for follow-up.",
				Notes:        "One or more feeding problem signs identified: infrequent breastfeeding, not breastfeeding at night, inappropriate complementary feeding, bottle feeding, incorrect replacement milk practices, child refuses food, shares meals, no active feeding, or feeding decreased during illness",
			},
			"NO_FEEDING_PROBLEM": {
				Classification: "NO FEEDING PROBLEM",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Praise and encourage the mother for feeding the infant well",
				},
				TreatmentPlan: "Positive reinforcement and continued good practices",
				FollowUp: []string{
					"Continue current feeding practices",
					"Routine follow-up as needed",
				},
				MotherAdvice: "Excellent! You are feeding your child well. Continue with these good practices.",
				Notes:        "No signs of feeding problems identified. Child is receiving appropriate feeding for their age and condition.",
			},
		},
	}
}
