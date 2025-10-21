// ruleengine/engine/child_tb_assessment_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildTBAssessmentTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "tb_assessment",
		Title:        "TB Infection Classification",
		Instructions: "ASSESS: TB symptoms, contact history, and clinical signs to classify TB infection",
		StartNode:    "tb_symptoms_check",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "tb_symptoms_check",
				Question:     "Does the child have any of these TB symptoms for more than 14 days?",
				QuestionType: "multi_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "ASK: About persistent TB symptoms lasting more than 14 days. If no symptoms, select 'None of the above'",
				Options: []domain.Option{
					{Value: "cough_14_days", DisplayText: "Cough for more than 14 days"},
					{Value: "fever_night_sweats_14_days", DisplayText: "Fever and/or night sweats for more than 14 days"},
					{Value: "none", DisplayText: "None of the above"},
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "tb_contact_history",
					},
				},
			},
			{
				NodeID:       "tb_contact_history",
				Question:     "Does the child have contact history with known Bacteriologically Confirmed Pulmonary TB (BC-PTB) patient?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "tb_symptoms_check",
				Instructions: "ASK: About contact with confirmed TB patient",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "tb_signs_check",
					},
					"no": {
						NextNode: "tb_signs_check",
					},
				},
			},
			{
				NodeID:       "tb_signs_check",
				Question:     "Does the child have any of these TB signs?",
				QuestionType: "multi_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "tb_contact_history",
				Instructions: "LOOK and FEEL: For TB clinical signs. If no signs, select 'None of the above'",
				Options: []domain.Option{
					{Value: "weight_loss_failure_gain", DisplayText: "Loss of weight or failure to gain weight or MAM or SAM"},
					{Value: "swelling_discharging_wound", DisplayText: "Swelling or discharging wound in the neck or armpit"},
					{Value: "none", DisplayText: "None of the above"},
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "check_investigation_need",
					},
				},
			},
			{
				NodeID:       "check_investigation_need",
				Question:     "Check if TB investigation is needed",
				QuestionType: "single_choice",
				Required:     true,
				Level:        4,
				ParentNode:   "tb_signs_check",
				Instructions: "System determines if investigation is needed based on symptoms, contact, or signs",
				Answers: map[string]domain.Answer{
					"investigate": {
						NextNode: "sputum_collection_method",
					},
					"no_investigation": {
						NextNode: "finalize_tb_classification",
					},
				},
			},
			{
				NodeID:       "sputum_collection_method",
				Question:     "What method was used to collect the sample for Gene Xpert test?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        5,
				ParentNode:   "check_investigation_need",
				Instructions: "SELECT: Sample collection method for Gene Xpert testing",
				Options: []domain.Option{
					{Value: "gastric_aspiration", DisplayText: "Sputum collected from Gastric Aspiration (NG Tube)"},
					{Value: "sputum_production", DisplayText: "Sputum collected from production"},
					{Value: "other_sites", DisplayText: "Sample collected from other sites"},
				},
				Answers: map[string]domain.Answer{
					"gastric_aspiration": {
						NextNode: "gene_xpert_test",
					},
					"sputum_production": {
						NextNode: "gene_xpert_test",
					},
					"other_sites": {
						NextNode: "gene_xpert_test",
					},
				},
			},
			{
				NodeID:       "gene_xpert_test",
				Question:     "What is the Gene Xpert test result?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        6,
				ParentNode:   "sputum_collection_method",
				Instructions: "TEST: Gene Xpert result from the collected sample",
				Options: []domain.Option{
					{Value: "positive", DisplayText: "Positive"},
					{Value: "negative", DisplayText: "Negative"},
					{Value: "not_done", DisplayText: "Not Done"},
				},
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "chest_xray_available",
					},
					"negative": {
						NextNode: "chest_xray_available",
					},
					"not_done": {
						NextNode: "chest_xray_available",
					},
				},
			},
			{
				NodeID:       "chest_xray_available",
				Question:     "Is Chest X-ray available and what is the result?",
				QuestionType: "single_choice",
				Required:     false,
				Level:        7,
				ParentNode:   "gene_xpert_test",
				Instructions: "CHECK: Chest X-ray if available",
				Options: []domain.Option{
					{Value: "suggestive_tb", DisplayText: "Suggestive of TB"},
					{Value: "not_suggestive", DisplayText: "Not suggestive of TB"},
					{Value: "not_available", DisplayText: "Not Available"},
				},
				Answers: map[string]domain.Answer{
					"suggestive_tb": {
						NextNode: "hiv_testing",
					},
					"not_suggestive": {
						NextNode: "hiv_testing",
					},
					"not_available": {
						NextNode: "hiv_testing",
					},
				},
			},
			{
				NodeID:       "hiv_testing",
				Question:     "Has provider-initiated HIV testing and counseling been done?",
				QuestionType: "yes_no",
				Required:     false,
				Level:        8,
				ParentNode:   "chest_xray_available",
				Instructions: "TEST: Provider-initiated HIV testing and counseling",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "finalize_tb_classification",
					},
					"no": {
						NextNode: "finalize_tb_classification",
					},
				},
			},
			{
				NodeID:       "finalize_tb_classification",
				Question:     "Finalize TB classification based on all collected information",
				QuestionType: "single_choice",
				Required:     true,
				Level:        9,
				ParentNode:   "hiv_testing",
				Instructions: "System will compute TB classification based on test results, symptoms, contact history, and clinical signs",
				Answers: map[string]domain.Answer{
					"compute": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"TB_INFECTION": {
				Classification: "TB INFECTION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Advise mother on the need for TB prevention treatment",
					"Ensure that mother is escorted and linked to TB clinic, for TB prevention treatment and follow up",
				},
				TreatmentPlan: "TB Infection - require TB prevention treatment",
				FollowUp: []string{
					"Referral to TB clinic for prevention treatment",
					"Regular monitoring during prevention treatment",
					"Watch for development of TB symptoms",
				},
				MotherAdvice: "Your child has been exposed to TB and has TB infection. Preventive treatment is needed to stop the infection from developing into active TB disease.",
				Notes:        "Contact history with BC-PTB patient but no TB symptoms and signs",
			},
			"NO_TB_INFECTION": {
				Classification: "NO TB INFECTION",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Continue and complete assessment and classification for other problems",
				},
				TreatmentPlan: "No TB Infection - routine care",
				FollowUp: []string{
					"Routine child health follow-up",
					"Return if TB symptoms develop",
				},
				MotherAdvice: "Based on current assessment, your child shows no signs of TB infection. Continue with good feeding practices and return if any TB symptoms develop.",
				Notes:        "No contact with known BC-PTB patient AND no TB signs/symptoms",
			},
		},
	}
}
