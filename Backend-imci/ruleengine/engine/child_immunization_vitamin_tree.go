package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

// Immunization and Vitamin A/Deworming assessment for children
func GetChildImmunizationVitaminTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "immunization_vitamin_status",
		Title:        "Immunization and Vitamin A/Deworming Status",
		Instructions: "Check scheduled immunizations by age and whether Vitamin A and Deworming (if eligible) were given in the last 6 months.",
		StartNode:    "child_age_months",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "child_age_months",
				Question:     "What is the child's current age in months?",
				QuestionType: "number_input",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Enter child's age in completed months.",
				Answers: map[string]domain.Answer{
					"value_based": {NextNode: "immunization_missing"},
				},
			},
			{
				NodeID:       "immunization_missing",
				Question:     "Are any scheduled vaccines for the child's age missing?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "child_age_months",
				Instructions: "Review the EPI schedule and confirm if any doses due for this age are missing.",
				Options: []domain.Option{
					{Value: "yes", DisplayText: "Yes, some vaccines are missing"},
					{Value: "no", DisplayText: "No, vaccines are up to date"},
				},
				Answers: map[string]domain.Answer{
					"yes": {NextNode: "vitamin_a_last_6months"},
					"no":  {NextNode: "vitamin_a_last_6months"},
				},
			},
			{
				NodeID:       "vitamin_a_last_6months",
				Question:     "Vitamin A supplementation in the last 6 months? (If <6 months, choose Not applicable)",
				QuestionType: "single_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "immunization_missing",
				Instructions: "If the child is 6 months or older, check if a dose of Vitamin A was given during the previous 6 months.",
				Options: []domain.Option{
					{Value: "received", DisplayText: "Received in last 6 months"},
					{Value: "not_received", DisplayText: "Not received in last 6 months"},
					{Value: "not_applicable", DisplayText: "Not applicable (<6 months)"},
				},
				Answers: map[string]domain.Answer{
					"received":       {NextNode: "deworming_last_6months"},
					"not_received":   {NextNode: "deworming_last_6months"},
					"not_applicable": {NextNode: "deworming_last_6months"},
				},
			},
			{
				NodeID:       "deworming_last_6months",
				Question:     "Deworming (Mebendazole/Albendazole) in the last 6 months? (If <24 months, choose Not applicable)",
				QuestionType: "single_choice",
				Required:     true,
				Level:        4,
				ParentNode:   "vitamin_a_last_6months",
				Instructions: "If the child is 2 years or older, check if deworming was given during the previous 6 months.",
				Options: []domain.Option{
					{Value: "received", DisplayText: "Received in last 6 months"},
					{Value: "not_received", DisplayText: "Not received in last 6 months"},
					{Value: "not_applicable", DisplayText: "Not applicable (<24 months)"},
				},
				Answers: map[string]domain.Answer{
					"received":       {Classification: "AUTO_CLASSIFY"},
					"not_received":   {Classification: "AUTO_CLASSIFY"},
					"not_applicable": {Classification: "AUTO_CLASSIFY"},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"IMMUNIZATION_AND_SUPPLEMENTS_UP_TO_DATE": {
				Classification: "IMMUNIZATION AND SUPPLEMENTS UP TO DATE",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Record doses on the child's card",
					"Advise to continue routine schedule and follow-up",
				},
				TreatmentPlan: "No pending immunization, Vitamin A or Deworming due",
				FollowUp:      []string{"Return for next scheduled visit"},
				MotherAdvice:  "Your child is up to date on vaccines and supplements.",
				Notes:         "All due items up to date",
			},
			"MISSING_IMMUNIZATIONS": {
				Classification: "MISSING IMMUNIZATIONS",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give due vaccines today as per EPI schedule",
					"Record doses on the child's card",
				},
				TreatmentPlan: "Catch-up immunization per national schedule",
				FollowUp:      []string{"Return for next visit per schedule"},
				MotherAdvice:  "Some vaccines are missing; we will catch up today.",
				Notes:         "One or more scheduled vaccines for age missing",
			},
			"VITAMIN_A_DUE": {
				Classification: "VITAMIN A DUE",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give Vitamin A supplementation (if 6 months or older)",
					"Record the dose on the child's card",
				},
				TreatmentPlan: "Vitamin A supplementation due",
				FollowUp:      []string{"Supplement every 6 months up to 5 years"},
				MotherAdvice:  "Your child needs Vitamin A today.",
				Notes:         "No Vitamin A dose in the previous 6 months",
			},
			"DEWORMING_DUE": {
				Classification: "DEWORMING DUE",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give Mebendazole/Albendazole (if 2 years or older)",
					"Record the dose on the child's card",
				},
				TreatmentPlan: "Deworming due",
				FollowUp:      []string{"Deworm every 6 months"},
				MotherAdvice:  "Your child needs deworming today.",
				Notes:         "No deworming in the previous 6 months",
			},
			"MULTIPLE_DUE": {
				Classification: "MULTIPLE PREVENTIVE CARE DUE",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Provide all due vaccines and supplements today",
					"Record all doses on the child's card",
				},
				TreatmentPlan: "Catch-up immunization and provide due Vitamin A/Deworming",
				FollowUp:      []string{"Return per schedule; Vitamin A/Deworming every 6 months"},
				MotherAdvice:  "We will provide all due items today and schedule follow-up.",
				Notes:         "Two or more items are due (vaccines and/or supplements)",
			},
		},
	}
}
