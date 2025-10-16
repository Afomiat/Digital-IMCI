// ruleengine/engine/child_cough_difficult_breathing_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildCoughDifficultBreathingTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "child_cough_difficult_breathing",
		Title:        "Assess Child with Cough or Difficult Breathing",
		Instructions: "ASK: Does the child have cough or difficult breathing? IF YES, LOOK, LISTEN, FEEL: Count breaths in one minute, Look for chest-indrawing, Look and listen for stridor, Look and listen for wheezing, Measure oxygen saturation. CHILD MUST BE CALM.",
		StartNode:    "cough_difficult_breathing",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "cough_difficult_breathing",
				Question:     "Does the child have cough or difficult breathing?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "ASK the mother: Does the child have cough or difficult breathing?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "how_long",
					},
					"no": {
						Classification: "NO_COUGH_DIFFICULT_BREATHING",
						Color:          "green",
					},
				},
			},
			{
				NodeID:       "how_long",
				Question:     "For how long has the child had cough or difficult breathing?",
				QuestionType: "number",
				Required:     true,
				Level:        2,
				ParentNode:   "cough_difficult_breathing",
				ShowCondition: "cough_difficult_breathing.yes",
				Instructions: "ASK: For how many days?",
				Answers: map[string]domain.Answer{
					"*": {
						NextNode: "child_calm",
					},
				},
			},
			{
				NodeID:       "child_calm",
				Question:     "Is the child calm?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "how_long",
				ShowCondition: "how_long.*",
				Instructions: "LOOK: Ensure child is calm before proceeding with assessment",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "general_danger_signs",
					},
					"no": {
						NextNode: "wait_for_calm",
					},
				},
			},
			{
				NodeID:       "wait_for_calm",
				Question:     "Wait until child is calm, then proceed?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "child_calm",
				ShowCondition: "child_calm.no",
				Instructions: "WAIT: Child must be calm for accurate assessment",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "general_danger_signs",
					},
					"no": {
						NextNode: "child_calm",
					},
				},
			},
			{
				NodeID:       "general_danger_signs",
				Question:     "Any general danger sign present?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "child_calm",
				ShowCondition: "child_calm.yes",
				Instructions: "CHECK: Any general danger sign (unable to drink, vomits everything, convulsions, lethargic, unconscious)",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "stridor",
					},
					"no": {
						NextNode: "stridor",
					},
				},
			},
			{
				NodeID:       "stridor",
				Question:     "Is stridor present in calm child?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        5,
				ParentNode:   "general_danger_signs",
				ShowCondition: "general_danger_signs.*",
				Instructions: "LISTEN: Check for stridor in calm child",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "oxygen_saturation",
					},
					"no": {
						NextNode: "fast_breathing",
					},
				},
			},
			{
				NodeID:       "oxygen_saturation",
				Question:     "Is oxygen saturation <90%?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        6,
				ParentNode:   "stridor",
				ShowCondition: "stridor.yes",
				Instructions: "MEASURE: Oxygen saturation with pulse oximeter",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "AUTO_CLASSIFY",
						Color:          "",
					},
					"no": {
						Classification: "AUTO_CLASSIFY",
						Color:          "",
					},
				},
			},
			{
				NodeID:       "fast_breathing",
				Question:     "Is there fast breathing for age?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        6,
				ParentNode:   "stridor",
				ShowCondition: "stridor.no",
				Instructions: "COUNT: Count breaths in one minute. Fast breathing: ≥50/min (2-12 months) or ≥40/min (12 months-5 years)",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "chest_indrawing",
					},
					"no": {
						NextNode: "chest_indrawing",
					},
				},
			},
			{
				NodeID:       "chest_indrawing",
				Question:     "Is chest indrawing present?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        7,
				ParentNode:   "fast_breathing",
				ShowCondition: "fast_breathing.*",
				Instructions: "LOOK: Check for chest indrawing",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "hiv_exposed",
					},
					"no": {
						NextNode: "wheezing",
					},
				},
			},
			{
				NodeID:       "hiv_exposed",
				Question:     "Is the child HIV exposed?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        8,
				ParentNode:   "chest_indrawing",
				ShowCondition: "chest_indrawing.yes",
				Instructions: "ASK: Is the child HIV exposed (mother HIV positive or unknown status)?",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "AUTO_CLASSIFY",
						Color:          "",
					},
					"no": {
						NextNode: "wheezing",
					},
				},
			},
			{
				NodeID:       "wheezing",
				Question:     "Is wheezing present?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        8,
				ParentNode:   "chest_indrawing",
				ShowCondition: "chest_indrawing.no",
				Instructions: "LISTEN: Check for wheezing in calm child",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "AUTO_CLASSIFY",
						Color:          "",
					},
					"no": {
						Classification: "AUTO_CLASSIFY",
						Color:          "",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"NO_COUGH_DIFFICULT_BREATHING": {
				Classification: "NO COUGH OR DIFFICULT BREATHING",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Continue with assessment of other symptoms",
				},
				TreatmentPlan: "No treatment needed for cough/difficult breathing",
				FollowUp: []string{
					"Assess other symptoms as needed",
				},
				MotherAdvice: "Child has no cough or difficult breathing.",
				Notes:        "No cough or difficult breathing detected",
			},
			"SEVERE_PNEUMONIA_OR_VERY_SEVERE_DISEASE": {
				Classification: "SEVERE PNEUMONIA OR VERY SEVERE DISEASE",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give first dose of IV/IM Ampicillin and Gentamicin",
				},
				TreatmentPlan: "Urgent pre-referral treatments and referral",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Child has severe pneumonia or very severe disease. Refer urgently to hospital.",
				Notes:        "Any general danger sign OR Stridor in calm child OR Oxygen saturation <90%",
			},
			"CHEST_INDRAWING_HIV_EXPOSED": {
				Classification: "PNEUMONIA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give first dose of amoxicillin",
					"Refer to hospital",
				},
				TreatmentPlan: "First dose antibiotic and referral for HIV exposed child with chest indrawing",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Refer for further management",
				},
				MotherAdvice: "Child has chest indrawing and is HIV exposed. Give first dose of amoxicillin and refer to hospital.",
				Notes:        "Chest indrawing in HIV exposed child - give first dose amoxicillin and refer",
			},
			"PNEUMONIA": {
				Classification: "PNEUMONIA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give oral Amoxicillin for 5 days",
					"Soothe the throat and relieve the cough with a safe remedy",
				},
				TreatmentPlan: "Oral antibiotics and symptomatic relief",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up after 2 days antibiotic treatment",
				},
				MotherAdvice: "Child has pneumonia. Give antibiotics for 5 days. Soothe throat and relieve cough. Return immediately if symptoms worsen.",
				Notes:        "Fast breathing OR Chest indrawing, no wheezing",
			},
			"PNEUMONIA_WITH_WHEEZING": {
				Classification: "PNEUMONIA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give oral Amoxicillin for 5 days",
					"Give rapid acting inhaled bronchodilator for up to 3 times, 15-20 minutes apart",
					"If wheezing disappears after bronchodilator, give inhaled bronchodilator for 5 days",
					"Soothe the throat and relieve the cough with a safe remedy",
				},
				TreatmentPlan: "Oral antibiotics, bronchodilator and symptomatic relief",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up after 2 days antibiotic treatment",
				},
				MotherAdvice: "Child has pneumonia with wheezing. Give antibiotics and bronchodilator. Soothe throat and relieve cough. Return immediately if symptoms worsen.",
				Notes:        "Fast breathing OR Chest indrawing, with wheezing",
			},
			"COUGH_OR_COLD": {
				Classification: "COUGH OR COLD",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Soothe the throat and relieve the cough with a safe remedy",
				},
				TreatmentPlan: "Symptomatic relief and monitoring",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up in 5 days if not improving",
					"If coughing > 14 days, assess for TB",
				},
				MotherAdvice: "Child has cough or cold. Soothe throat and relieve cough. Return immediately if symptoms worsen or if coughing continues for more than 14 days.",
				Notes:        "No signs of very severe disease AND no pneumonia",
			},
			"COUGH_OR_COLD_WITH_WHEEZING": {
				Classification: "COUGH OR COLD",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Give an inhaled bronchodilator for 5 days",
					"Soothe the throat and relieve the cough with a safe remedy",
				},
				TreatmentPlan: "Bronchodilator and symptomatic relief",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up in 5 days if not improving",
					"If coughing > 14 days, assess for TB",
				},
				MotherAdvice: "Child has cough or cold with wheezing. Give bronchodilator for 5 days. Soothe throat and relieve cough. Return immediately if symptoms worsen.",
				Notes:        "Wheezing present but no fast breathing or chest indrawing",
			},
		},
	}
}