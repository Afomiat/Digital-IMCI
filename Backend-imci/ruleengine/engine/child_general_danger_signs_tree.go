// ruleengine/engine/child_general_danger_signs_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildGeneralDangerSignsTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "child_general_danger_signs",
		Title:        "Check for General Danger Signs in Child",
		Instructions: "ASK: Is the child able to drink or breastfeed? Does the child vomit everything? Has the child had convulsions? LOOK: See if the child is lethargic or unconscious. Is the child convulsing now?",
		StartNode:    "unable_to_drink_breastfeed",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "unable_to_drink_breastfeed",
				Question:     "Is the child able to drink or breastfeed?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "ASK: Is the child able to drink or breastfeed?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "vomits_everything",
					},
					"no": {
						Classification: "VERY_SEVERE_DISEASE",
						Color:          "pink",
						EmergencyPath:  true,
					},
				},
			},
			{
				NodeID:       "vomits_everything",
				Question:     "Does the child vomit everything?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "unable_to_drink_breastfeed",
				ShowCondition: "unable_to_drink_breastfeed.yes",
				Instructions: "ASK: Does the child vomit everything?",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "VERY_SEVERE_DISEASE",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "convulsions_history",
					},
				},
			},
			{
				NodeID:       "convulsions_history",
				Question:     "Has the child had convulsions?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "vomits_everything",
				ShowCondition: "vomits_everything.no",
				Instructions: "ASK: Has the child had convulsions?",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "VERY_SEVERE_DISEASE",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "lethargic_unconscious",
					},
				},
			},
			{
				NodeID:       "lethargic_unconscious",
				Question:     "Is the child lethargic or unconscious?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "convulsions_history",
				ShowCondition: "convulsions_history.no",
				Instructions: "LOOK: See if the child is lethargic or unconscious",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "VERY_SEVERE_DISEASE",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "convulsing_now",
					},
				},
			},
			{
				NodeID:       "convulsing_now",
				Question:     "Is the child convulsing now?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        5,
				ParentNode:   "lethargic_unconscious",
				ShowCondition: "lethargic_unconscious.no",
				Instructions: "LOOK: Is the child convulsing now?",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "VERY_SEVERE_DISEASE",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						Classification: "NO_GENERAL_DANGER_SIGNS",
						Color:          "green",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"VERY_SEVERE_DISEASE": {
				Classification: "VERY SEVERE DISEASE",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give diazepam if convulsing now",
					"Quickly complete the assessment",
					"Give appropriate pre-referral treatment immediately, based on other severe classifications",
					"Treat to prevent low blood sugar",
					"Keep the child warm",
					"Refer URGENTLY",
				},
				TreatmentPlan: "Urgent pre-referral treatments and referral",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Child has very severe disease. Refer urgently to hospital. Keep child warm during transport.",
				Notes:        "Child has one or more general danger signs: Unable to drink/breastfeed, vomits everything, had convulsions, lethargic, unconscious, or convulsing now",
			},
			"NO_GENERAL_DANGER_SIGNS": {
				Classification: "NO GENERAL DANGER SIGNS",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Continue with assessment of main symptoms",
				},
				TreatmentPlan: "No urgent referral needed",
				FollowUp: []string{
					"Assess other symptoms as needed",
				},
				MotherAdvice: "Child has no general danger signs. Continue with assessment of other symptoms.",
				Notes:        "No general danger signs detected",
			},
		},
	}
}