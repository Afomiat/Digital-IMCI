// ruleengine/engine/jaundice_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetJaundiceTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "jaundice_check",
		Title:        "Check for Jaundice in Young Infant",
		Instructions: "Assess and classify jaundice in sick young infant from birth up to 2 months. Look for jaundice in natural light.",
		StartNode:    "skin_yellow",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "skin_yellow",
				Question:     "Is there yellowish discoloration of the skin or eyes?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Look for jaundice in natural light. Check skin and eyes carefully.",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "palms_soles_yellow",
					},
					"no": {
						Classification: "NO_JAUNDICE",
						Color:          "green",
					},
				},
			},
			{
				NodeID:       "palms_soles_yellow",
				Question:     "Are the palms and/or soles yellow?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "skin_yellow",
				ShowCondition: "skin_yellow.yes",
				Instructions: "Check palms and soles carefully in natural light",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "SEVERE_JAUNDICE_URGENT",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "infant_age",
					},
				},
			},
			{
				NodeID:       "infant_age",
				Question:     "How old is the infant in days?",
				QuestionType: "number_input",
				Required:     true,
				Level:        3,
				ParentNode:   "palms_soles_yellow",
				ShowCondition: "palms_soles_yellow.no",
				Instructions: "Enter age in days (0 for less than 24 hours)",
				Validation: &domain.Validation{
					Min:  0,
					Max:  60,
					Step: 1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"SEVERE_JAUNDICE_URGENT": {
				Classification: "SEVERE JAUNDICE",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Treat to prevent low blood sugar",
					"Warm the young infant by skin-to-skin contact if temperature is less than 36.5°C while arranging referral",
					"Advise mother how to keep the young infant warm on the way to the hospital",
				},
				TreatmentPlan: "Urgent hospital referral",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Go to hospital immediately - keep infant warm during transport",
				Notes:        "Severe jaundice requiring urgent medical attention",
			},
			"JAUNDICE": {
				Classification: "JAUNDICE",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Advise mother to give home care for the young infant",
					"Advise the mother to expose and check in natural light daily",
				},
				TreatmentPlan: "Home care with monitoring",
				FollowUp: []string{
					"Follow-up after 2 days",
					"Advise mother to return immediately if the infant's palms or soles appear yellow",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "Mild jaundice in appropriate age range - monitor closely",
			},
			"NO_JAUNDICE": {
				Classification: "NO JAUNDICE",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Advise mother to give home care for the infant",
				},
				TreatmentPlan: "Routine home care",
				FollowUp: []string{
					"Routine follow-up as scheduled",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "No signs of jaundice detected",
			},
		},
	}
}