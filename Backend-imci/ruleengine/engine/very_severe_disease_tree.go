// ruleengine/engine/very_severe_disease_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetVerySevereDiseaseTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "very_severe_disease_check",
		Title:        "Check for Very Severe Disease and Local Bacterial Infection",
		Instructions: "Assess and classify sick young infant from birth up to 2 months. Count the breaths in one minute. Repeat if count is ≥60 breaths/minute. Look, listen and feel for all signs.",
		StartNode:    "difficulty_feeding",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "difficulty_feeding",
				Question:     "Is the infant having difficulty in feeding?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "If yes, assess feeding ability in detail. See if the infant is not feeding.",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "feeding_ability_detail",
						Color:    "yellow",
					},
					"no": {
						NextNode: "convulsions_history",
					},
				},
			},
			{
				NodeID:       "feeding_ability_detail",
				Question:     "How would you describe the feeding difficulty?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "difficulty_feeding",
				ShowCondition: "difficulty_feeding.yes",
				Instructions: "Assess the infant's ability to feed properly",
				Answers: map[string]domain.Answer{
					"unable_to_feed": {
						Classification: "CRITICAL_ILLNESS",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"not_feeding_well": {
						NextNode: "convulsions_history",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "convulsions_history",
				Question:     "Has the infant had convulsions or is convulsing now?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Look if the infant is convulsing now",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "CRITICAL_ILLNESS",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "check_movements",
					},
				},
			},
			{
				NodeID:       "check_movements",
				Question:     "Look at the young infant's movements (if sleeping, ask mother to wake him/her)",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Observe the infant's spontaneous movement",
				Answers: map[string]domain.Answer{
					"moves_on_own": {
						NextNode: "breathing_rate",
					},
					"moves_only_when_stimulated": {
						NextNode: "breathing_rate",
						Color:    "yellow",
					},
					"no_movement_even_stimulated": {
						Classification: "CRITICAL_ILLNESS",
						Color:          "pink",
						EmergencyPath:  true,
					},
				},
			},
			{
				NodeID:       "breathing_rate",
				Question:     "Count the breaths in one minute",
				QuestionType: "number_input",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Repeat the count if ≥60 breaths/minute",
				Validation: &domain.Validation{
					Min:  0,
					Max:  100,
					Step: 1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "chest_indrawing",
					},
				},
			},
			{
				NodeID:       "chest_indrawing",
				Question:     "Look for severe chest indrawing?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Observe for severe chest indentation",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "VERY_SEVERE_DISEASE",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "umbilicus_check",
					},
				},
			},
			{
				NodeID:       "umbilicus_check",
				Question:     "Look at the umbilicus - is it red or draining pus?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Examine the umbilical area carefully",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "skin_pustules",
						Color:    "yellow",
					},
					"no": {
						NextNode: "skin_pustules",
					},
				},
			},
			{
				NodeID:       "skin_pustules",
				Question:     "Look for skin pustules?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Check entire skin surface for pustules",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "temperature_measurement",
						Color:    "yellow",
					},
					"no": {
						NextNode: "temperature_measurement",
					},
				},
			},
			{
				NodeID:       "temperature_measurement",
				Question:     "Measure axillary temperature (°C)",
				QuestionType: "number_input",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Note: Rectal temperature thresholds are approximately 0.5°C higher. Consider malaria in young infant with fever based on associated symptoms.",
				Validation: &domain.Validation{
					Min:  30.0,
					Max:  42.0,
					Step: 0.1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"CRITICAL_ILLNESS": {
				Classification: "CRITICAL ILLNESS",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give first dose of Ampicillin and Gentamicin",
					"Advise mother how to keep the infant warm on the way to the hospital",
					"Reinforce referral and administer URGENTLY to hospital",
				},
				TreatmentPlan: "Immediate hospital referral",
				FollowUp: []string{
					"Refer immediately to hospital",
					"Ensure proper referral process",
				},
				MotherAdvice: "Go to hospital immediately - keep infant warm during transport",
				Notes:        "If referral is not possible, see 'Where Referral is not Possible' protocol",
			},
			"VERY_SEVERE_DISEASE": {
				Classification: "VERY SEVERE DISEASE",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give first dose of Ampicillin and Gentamicin",
					"Treat for low blood sugar",
					"Warm infant by skin-to-skin contact if temperature <35.5°C while continuing referral",
					"Advise mother how to keep the infant warm on the way to the hospital",
				},
				TreatmentPlan: "Urgent hospital referral",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Go to hospital urgently - keep infant warm during transport",
				Notes:        "If referral is not possible, see 'Where Referral is not Possible' protocol",
			},
			"PNEUMONIA": {
				Classification: "PNEUMONIA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give Ampicillin for 7 days",
				},
				TreatmentPlan: "Outpatient antibiotic treatment",
				FollowUp: []string{
					"Follow-up after 2 days of Ampicillin",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "For infants ≥7 days old with fast breathing (≥60 breaths/minute)",
			},
			"LOCAL_BACTERIAL_INFECTION": {
				Classification: "LOCAL BACTERIAL INFECTION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give Ampicillin for 5 days",
					"Teach mother to treat local infections at home",
				},
				TreatmentPlan: "Outpatient antibiotic treatment and home care",
				FollowUp: []string{
					"Follow-up after 2 days of Ampicillin",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "Red umbilicus, draining pus, or skin pustules",
			},
			"SEVERE_INFECTION_UNLIKELY": {
				Classification: "SEVERE INFECTION UNLIKELY",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Warm the infant using skin-to-skin contact for one hour if temperature 35.5°C - 36.4°C and reassess",
					"If same temperature after warming, advise mother on how to keep infant warm at home",
					"Advise mother to give home care for the infant",
				},
				TreatmentPlan: "Home care with temperature management",
				FollowUp: []string{
					"Routine follow-up as scheduled",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "No signs of critical illness, very severe disease, pneumonia, or local bacterial infection",
			},
			"AUTO_CLASSIFY": {
				Classification: "SEVERE INFECTION UNLIKELY",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Review all clinical findings for proper classification",
					"Consider reassessment if any doubts",
				},
				TreatmentPlan: "Further assessment needed",
				FollowUp: []string{
					"Reassess if condition changes",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "Automatic classification based on collected signs - review clinical findings",
			},
		},
	}
}