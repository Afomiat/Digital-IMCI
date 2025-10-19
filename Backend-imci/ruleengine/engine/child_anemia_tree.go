package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildAnemiaTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "child_anemia_check",
		Title:        "Check for Anemia",
		Instructions: "Look for palmar pallor. If some or severe pallor and child ≥ 6 months, do blood test and measure Haemoglobin (Hb) or Haematocrit (Hct)",
		StartNode:    "palmar_pallor_present",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "palmar_pallor_present",
				Question:     "Look for palmar pallor. Is there palmar pallor?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Look for palmar pallor",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "palmar_pallor_severity",
					},
					"no": {
						Classification: "NO_ANEMIA",
						Color:          "green",
					},
				},
			},
			{
				NodeID:       "palmar_pallor_severity",
				Question:     "Is the palmar pallor:",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "palmar_pallor_present",
				ShowCondition: "palmar_pallor_present.yes",
				Instructions: "Determine severity of palmar pallor",
				Answers: map[string]domain.Answer{
					"severe_palmar_pallor": {
						NextNode: "child_age",
					},
					"some_palmar_pallor": {
						NextNode: "child_age",
					},
					"no_palmar_pallor": {
						Classification: "NO_ANEMIA",
						Color:          "green",
					},
				},
			},
			{
				NodeID:       "child_age",
				Question:     "Is the child ≥ 6 months old?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "palmar_pallor_severity",
				ShowCondition: "palmar_pallor_severity.severe_palmar_pallor OR palmar_pallor_severity.some_palmar_pallor",
				Instructions: "Check if child is 6 months or older for blood test",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "blood_test_done",
					},
					"no": {
						NextNode: "classify_by_pallor_only",
					},
				},
			},
			{
				NodeID:       "blood_test_done",
				Question:     "Was blood test done to measure Haemoglobin (Hb) or Haematocrit (Hct)?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "child_age",
				ShowCondition: "child_age.yes",
				Instructions: "If some or severe pallor and child ≥ 6 months, do blood test",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "hb_value",
					},
					"no": {
						NextNode: "classify_by_pallor_only",
					},
				},
			},
			{
				NodeID:       "classify_by_pallor_only",
				Question:     "Classify anemia based on palmar pallor:",
				QuestionType: "single_choice",
				Required:     true,
				Level:        5,
				ParentNode:   "child_age",
				ShowCondition: "child_age.no OR blood_test_done.no",
				Instructions: "Classify based on palmar pallor when blood test not done",
				Answers: map[string]domain.Answer{
					"severe_pallor": {
						Classification: "SEVERE_ANEMIA",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"some_pallor": {
						Classification: "ANEMIA",
						Color:          "yellow",
					},
				},
			},
			{
				NodeID:       "hb_value",
				Question:     "What is the Haemoglobin (Hb) value in gm/dL?",
				QuestionType: "number_input",
				Required:     false,
				Level:        5,
				ParentNode:   "blood_test_done",
				ShowCondition: "blood_test_done.yes",
				Instructions: "Enter Hb value in gm/dL (e.g., 6.5, 8.0, 12.0)",
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "hct_value",
					},
				},
			},
			{
				NodeID:       "hct_value",
				Question:     "What is the Haematocrit (Hct) value in %?",
				QuestionType: "number_input",
				Required:     false,
				Level:        5,
				ParentNode:   "hb_value",
				ShowCondition: "blood_test_done.yes",
				Instructions: "Enter Hct value in percentage (e.g., 20, 25, 35)",
				Answers: map[string]domain.Answer{
					"value_based": {
						Classification: "AUTO_CLASSIFY",
						Color:          "white",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"SEVERE_ANEMIA": {
				Classification: "SEVERE ANEMIA",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Refer URGENTLY to hospital",
				},
				TreatmentPlan: "Urgent hospital referral for severe anemia",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Go to hospital immediately - this is severe anemia requiring urgent care",
				Notes:        "Hb < 7gm/dL OR Hct < 21% OR Severe palmar pallor",
			},
			"ANEMIA": {
				Classification: "ANEMIA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Assess the child's feeding and counsel the mother on feeding according to the FOOD box on the COUNSEL THE MOTHER chart",
					"Give Iron",
					"Do blood film for malaria, if malaria risk is high or has travel history to malarious area in last 30 days",
					"Give Albendazole or Mebendazole, if the child is ≥ 1 year old and has not had a dose in the previous six months",
					"Advise mother when to return immediately",
				},
				TreatmentPlan: "Iron supplementation, feeding counseling, and malaria screening",
				FollowUp: []string{
					"Follow-up in 14 days",
				},
				MotherAdvice: "Give iron supplements as prescribed and improve feeding. Return immediately if child becomes more pale or weak",
				Notes:        "Hb 7—<11 gm/dL OR Hct 21—<33% OR Some palmar pallor",
			},
			"NO_ANEMIA": {
				Classification: "NO ANEMIA",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"No additional treatment",
					"Counsel the mother on feeding recommendations",
				},
				TreatmentPlan: "No anemia treatment needed",
				FollowUp: []string{
					"Routine follow-up",
				},
				MotherAdvice: "Continue with normal feeding practices as recommended",
				Notes:        "Hb ≥ 11gm/dL OR Hct ≥ 33% OR No palmar pallor",
			},
		},
	}
}