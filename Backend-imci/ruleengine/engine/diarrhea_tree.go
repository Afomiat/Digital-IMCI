package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetDiarrheaTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "diarrhea_check",
		Title:        "Check for Diarrhea and Dehydration in Young Infant",
		Instructions: "Assess young infant for diarrhea and classify dehydration severity. Check for blood in stool and dehydration signs.",
		StartNode:    "diarrhea_present",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "diarrhea_present",
				Question:     "Does the young infant have diarrhea?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Assess for loose or watery stools",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "diarrhea_duration",
					},
					"no": {
						Classification: "NO_DIARRHEA",
						Color:          "green",
					},
				},
			},
			{
				NodeID:       "diarrhea_duration",
				Question:     "For how long has the infant had diarrhea?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "diarrhea_present",
				ShowCondition: "diarrhea_present.yes",
				Instructions: "Determine duration of diarrhea episode",
				Answers: map[string]domain.Answer{
					"less_than_14_days": {
						NextNode: "blood_in_stool",
					},
					"14_days_or_more": {
						Classification: "SEVERE_PERSISTENT_DIARRHEA",
						Color:          "pink",
						EmergencyPath:  true,
					},
				},
			},
			{
				NodeID:       "blood_in_stool",
				Question:     "Is there blood in the stool?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "diarrhea_duration",
				ShowCondition: "diarrhea_present.yes",
				Instructions: "Check stool for visible blood",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "DYSENTERY",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "movement_condition",
					},
				},
			},
			{
				NodeID:       "movement_condition",
				Question:     "Look at the young infant's general condition - how does the infant move?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        4,
				ParentNode:   "blood_in_stool",
				ShowCondition: "blood_in_stool.no",
				Instructions: "Observe infant's spontaneous movement and response to stimulation",
				Answers: map[string]domain.Answer{
					"moves_on_own": {
						NextNode: "sunken_eyes",
					},
					"moves_only_when_stimulated": {
						NextNode: "sunken_eyes",
						Color:    "pink",
					},
					"no_movement_even_when_stimulated": {
						Classification: "SEVERE_DEHYDRATION",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"restless_irritable": {
						NextNode: "sunken_eyes",
						Color:    "yellow",
					},
				},
			},
			{
				NodeID:       "sunken_eyes",
				Question:     "Look for sunken eyes?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        5,
				ParentNode:   "movement_condition",
				ShowCondition: "movement_condition.moves_on_own OR movement_condition.moves_only_when_stimulated OR movement_condition.restless_irritable",
				Instructions: "Check for sunken eyes as sign of dehydration",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "skin_pinch",
					},
					"no": {
						NextNode: "skin_pinch",
					},
				},
			},
			{
				NodeID:       "skin_pinch",
				Question:     "Pinch the skin of the abdomen. How does it go back?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        6,
				ParentNode:   "sunken_eyes",
				ShowCondition: "sunken_eyes.yes OR sunken_eyes.no",
				Instructions: "Pinch skin and observe return time",
				Answers: map[string]domain.Answer{
					"very_slowly_more_than_2_seconds": {
						Classification: "SEVERE_DEHYDRATION",
						Color:          "pink",
					},
					"slowly": {
						Classification: "SOME_DEHYDRATION",
						Color:          "yellow",
					},
					"immediately": {
						Classification: "NO_DEHYDRATION",
						Color:          "green",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"DYSENTERY": {
				Classification: "DYSENTERY",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give first dose of IM Ampicillin and Gentamicin",
					"Treat to prevent low blood sugar",
					"Advise how to keep infant warm on the way to the hospital",
					"Refer to hospital",
				},
				TreatmentPlan: "Urgent hospital referral with antibiotic treatment",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Go to hospital immediately - keep infant warm during transport",
				Notes:        "Blood in stool indicates dysentery requiring antibiotics",
			},
			"SEVERE_PERSISTENT_DIARRHEA": {
				Classification: "SEVERE PERSISTENT DIARRHOEA",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give first dose of IM Ampicillin and Gentamicin",
					"Treat to prevent low blood sugar",
					"Advise how to keep infant warm on the way to the hospital",
					"Refer to hospital",
				},
				TreatmentPlan: "Urgent hospital referral for persistent diarrhea",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Go to hospital immediately - diarrhea lasting 14 days or more",
				Notes:        "Diarrhea lasting 14 days or more requires hospital management",
			},
			"SEVERE_DEHYDRATION": {
				Classification: "SEVERE DEHYDRATION",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"If infant has another severe classification:",
					"- Refer URGENTLY to hospital with mother giving frequent sips of ORS on the way",
					"- Advise mother to continue breastfeeding more frequently",
					"- Advise mother how to keep the young infant warm on the way to hospital",
					"",
					"If infant does not have any other severe classification:",
					"- Give fluid for severe dehydration (Plan C)",
				},
				TreatmentPlan: "Plan C - Severe dehydration treatment",
				FollowUp: []string{
					"Follow referral instructions if other severe conditions present",
					"Follow-up in 2 days if treating with Plan C",
				},
				MotherAdvice: "Give ORS frequently and continue breastfeeding",
				Notes:        "Severe dehydration - treatment depends on presence of other severe classifications",
			},
			"SOME_DEHYDRATION": {
				Classification: "SOME DEHYDRATION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"If infant has another severe classification:",
					"- Refer URGENTLY to hospital with mother giving frequent sips of ORS on the way",
					"- Advise mother to continue breastfeeding more frequently",
					"- Advise mother how to keep the young infant warm on the way to hospital",
					"",
					"If infant does not have any other severe classification:",
					"- Give fluid for some dehydration and Zinc supplement (Plan B)",
					"- Advise mother when to return immediately",
				},
				TreatmentPlan: "Plan B - Some dehydration treatment",
				FollowUp: []string{
					"Follow referral instructions if other severe conditions present",
					"Follow-up in 2 days if treating with Plan B",
				},
				MotherAdvice: "Give ORS with zinc and continue breastfeeding",
				Notes:        "Some dehydration - treatment depends on presence of other severe classifications",
			},
			"NO_DEHYDRATION": {
				Classification: "NO DEHYDRATION",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Give fluids to treat diarrhoea at home and Zinc supplement (Plan A)",
					"Advise mother to continue breastfeeding",
					"Advise mother when to return immediately",
				},
				TreatmentPlan: "Home care with ORS and Zinc",
				FollowUp: []string{
					"Follow-up in 5 days if not improving",
				},
				MotherAdvice: "Continue breastfeeding and give ORS with zinc at home",
				Notes:        "No signs of dehydration",
			},
			"NO_DIARRHEA": {
				Classification: "NO DIARRHEA",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Continue routine care",
				},
				TreatmentPlan: "Routine infant care",
				FollowUp: []string{
					"Routine follow-up as scheduled",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "No diarrhea present",
			},
		},
	}
}