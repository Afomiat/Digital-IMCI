// ruleengine/engine/child_fever_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildFeverTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "child_fever",
		Title:        "Assess Child with Fever",
		Instructions: "ASK: Does the child have fever? IF YES: Determine malaria risk, assess duration, check for measles, examine for danger signs",
		StartNode:    "fever_present",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "fever_present",
				Question:     "Does the child have fever?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "ASK: Does the child have fever?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "malaria_risk",
					},
					"no": {
						Classification: "NO_FEVER",
						Color:          "green",
					},
				},
			},
			{
				NodeID:       "malaria_risk",
				Question:     "What is the malaria risk in this area?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "fever_present",
				ShowCondition: "fever_present.yes",
				Instructions: "Decide Malaria Risk: High/Low or No",
				Answers: map[string]domain.Answer{
					"high": {
						NextNode: "do_blood_film_high",
					},
					"low": {
						NextNode: "do_blood_film_low",
					},
					"no": {
						NextNode: "travel_outside",
					},
				},
			},
			{
				NodeID:       "travel_outside",
				Question:     "Has the child traveled outside this area during the previous 30 days?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "malaria_risk",
				ShowCondition: "malaria_risk.no",
				Instructions: "ASK: Has the child traveled outside this area during the previous 30 days?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "travel_malarious_area",
					},
					"no": {
						NextNode: "fever_duration",
					},
				},
			},
			{
				NodeID:       "travel_malarious_area",
				Question:     "Has the child been to a malarious area?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "travel_outside",
				ShowCondition: "travel_outside.yes",
				Instructions: "ASK: If yes, has he been to a malarious area?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "do_blood_film_travel",
					},
					"no": {
						NextNode: "fever_duration",
					},
				},
			},
			{
				NodeID:       "do_blood_film_high",
				Question:     "Perform blood film for malaria",
				QuestionType: "single_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "malaria_risk",
				ShowCondition: "malaria_risk.high",
				Instructions: "Do blood film for all children in high malaria risk areas",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "fever_duration",
					},
					"negative": {
						NextNode: "fever_duration",
					},
					"not_available": {
						NextNode: "fever_duration",
					},
				},
			},
			{
				NodeID:       "do_blood_film_low",
				Question:     "Perform blood film for malaria",
				QuestionType: "single_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "malaria_risk",
				ShowCondition: "malaria_risk.low",
				Instructions: "Do blood film if no other obvious cause of fever",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "fever_duration",
					},
					"negative": {
						NextNode: "fever_duration",
					},
					"not_done": {
						NextNode: "fever_duration",
					},
				},
			},
			{
				NodeID:       "do_blood_film_travel",
				Question:     "Perform blood film for malaria",
				QuestionType: "single_choice",
				Required:     true,
				Level:        5,
				ParentNode:   "travel_malarious_area",
				ShowCondition: "travel_malarious_area.yes",
				Instructions: "Do blood film for children who traveled to malarious areas",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "fever_duration",
					},
					"negative": {
						NextNode: "fever_duration",
					},
					"not_available": {
						NextNode: "fever_duration",
					},
				},
			},
			{
				NodeID:       "fever_duration",
				Question:     "For how long has the child had fever?",
				QuestionType: "number",
				Required:     true,
				Level:        6,
				ParentNode:   "malaria_risk",
				ShowCondition: "malaria_risk.*",
				Instructions: "ASK: For how long has the child had fever?",
				Answers: map[string]domain.Answer{
					"*": {
						NextNode: "measles_history",
					},
				},
			},
			{
				NodeID:       "measles_history",
				Question:     "Has the child had measles within the last 3 months?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        7,
				ParentNode:   "fever_duration",
				ShowCondition: "fever_duration.*",
				Instructions: "ASK: Has the child had measles within the last 3 months?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "current_measles",
					},
					"no": {
						NextNode: "current_measles",
					},
				},
			},
			{
				NodeID:       "current_measles",
				Question:     "Does the child have measles now?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        8,
				ParentNode:   "measles_history",
				ShowCondition: "measles_history.*",
				Instructions: "LOOK: Generalized rash AND one of these: cough, runny nose or red eyes",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "stiff_neck",
					},
					"no": {
						NextNode: "stiff_neck",
					},
				},
			},
			{
				NodeID:       "stiff_neck",
				Question:     "Look or feel for stiff neck",
				QuestionType: "yes_no",
				Required:     true,
				Level:        9,
				ParentNode:   "current_measles",
				ShowCondition: "current_measles.*",
				Instructions: "LOOK AND FEEL: Look or feel for stiff neck",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "bulging_fontanelle",
					},
					"no": {
						NextNode: "bulging_fontanelle",
					},
				},
			},
			{
				NodeID:       "bulging_fontanelle",
				Question:     "Look or feel for bulging fontanelle",
				QuestionType: "yes_no",
				Required:     true,
				Level:        10,
				ParentNode:   "stiff_neck",
				ShowCondition: "stiff_neck.*",
				Instructions: "LOOK AND FEEL: Look or feel for bulging fontanelle (< 1 year of age)",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "runny_nose",
					},
					"no": {
						NextNode: "runny_nose",
					},
				},
			},
			{
				NodeID:       "runny_nose",
				Question:     "Look for runny nose",
				QuestionType: "yes_no",
				Required:     true,
				Level:        11,
				ParentNode:   "bulging_fontanelle",
				ShowCondition: "bulging_fontanelle.*",
				Instructions: "LOOK: Look for runny nose",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "other_bacterial_cause",
					},
					"no": {
						NextNode: "other_bacterial_cause",
					},
				},
			},
			{
				NodeID:       "other_bacterial_cause",
				Question:     "Look for any obvious other bacterial causes of fever",
				QuestionType: "yes_no",
				Required:     true,
				Level:        12,
				ParentNode:   "runny_nose",
				ShowCondition: "runny_nose.*",
				Instructions: "LOOK: Look for any obvious other bacterial causes of fever",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "mouth_ulcers",
					},
					"no": {
						NextNode: "mouth_ulcers",
					},
				},
			},
			{
				NodeID:       "mouth_ulcers",
				Question:     "Look for mouth ulcers",
				QuestionType: "single_choice",
				Required:     true,
				Level:        13,
				ParentNode:   "other_bacterial_cause",
				ShowCondition: "other_bacterial_cause.*",
				Instructions: "LOOK: Look for mouth ulcers: Are they deep or extensive? Are they not deep or extensive?",
				Answers: map[string]domain.Answer{
					"deep_extensive": {
						NextNode: "eye_pus",
					},
					"not_deep_extensive": {
						NextNode: "eye_pus",
					},
					"none": {
						NextNode: "eye_pus",
					},
				},
			},
			{
				NodeID:       "eye_pus",
				Question:     "Look for pus draining from the eye",
				QuestionType: "yes_no",
				Required:     true,
				Level:        14,
				ParentNode:   "mouth_ulcers",
				ShowCondition: "mouth_ulcers.*",
				Instructions: "LOOK: Look for pus draining from the eye",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "clouding_cornea",
					},
					"no": {
						NextNode: "clouding_cornea",
					},
				},
			},
			{
				NodeID:       "clouding_cornea",
				Question:     "Look for clouding of the cornea",
				QuestionType: "yes_no",
				Required:     true,
				Level:        15,
				ParentNode:   "eye_pus",
				ShowCondition: "eye_pus.*",
				Instructions: "LOOK: Look for clouding of the cornea",
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
			"NO_FEVER": {
				Classification: "NO FEVER",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Continue with assessment of other symptoms",
				},
				TreatmentPlan: "No treatment needed for fever",
				FollowUp: []string{
					"Assess other symptoms as needed",
				},
				MotherAdvice: "Child has no fever.",
				Notes:        "No fever detected",
			},
			"VERY_SEVERE_FEBRILE_DISEASE": {
				Classification: "VERY SEVERE FEBRILE DISEASE",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give first dose IV/IM Artesunate for severe malaria",
					"Give first dose of IV/IM Ampicillin and Gentamicin",
					"Give sugar to prevent low blood sugar",
					"Give Paracetamol in health facility for high fever (≥38.5°C)",
					"Refer URGENTLY to hospital",
				},
				TreatmentPlan: "Urgent referral to hospital with pre-referral treatments",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Refer URGENTLY to hospital",
				Notes:        "Any general danger sign OR stiff neck OR bulging fontanelle",
			},
			"MALARIA_HIGH_RISK": {
				Classification: "MALARIA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Treat with Artemisinin-Lumefantrine (AL) and Primaquine for P. falciparum or mixed or no confirmatory test done",
					"Treat with Chloroquine and Primaquine for confirmed P. vivax",
					"Give Paracetamol in health facility for high fever (38.5°C or above)",
					"Give an appropriate antibiotic for identified bacterial cause of fever",
				},
				TreatmentPlan: "Antimalarial treatment based on malaria risk",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up after 2 days of antimalarial if fever persists or if on Primaquine",
					"If fever is present every day for more than 7 days, refer for assessment",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "High malaria risk with positive blood film or blood film not available",
			},
			"MALARIA_LOW_RISK": {
				Classification: "MALARIA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Treat with Artemisinin-Lumefantrine (AL) and Primaquine for P. falciparum or mixed or no confirmatory test done",
					"Treat with Chloroquine and Primaquine for confirmed P. vivax",
					"Give Paracetamol in health facility for high fever (38.5°C or above)",
					"Give an appropriate antibiotic for identified bacterial cause of fever",
				},
				TreatmentPlan: "Antimalarial treatment with antibiotics for bacterial causes",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up after 2 days of antimalarial if fever persists or if on Primaquine",
					"If fever is present every day for more than 7 days, refer for assessment",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "Low malaria risk with positive blood film",
			},
			"FEVER_NO_MALARIA": {
				Classification: "FEVER: NO MALARIA",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Give one dose of Paracetamol in health facility for high fever (≥38.5°C)",
					"Give an appropriate antibiotic for identified bacterial cause of fever",
				},
				TreatmentPlan: "Symptomatic treatment with antibiotics for bacterial causes",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up after 2 days of antibiotics if fever persists",
					"If fever is present every day for more than 7 days, refer for assessment",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "Negative blood film OR other obvious cause of fever present",
			},
			"SEVERE_COMPLICATED_MEASLES": {
				Classification: "SEVERE COMPLICATED MEASLES",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give Vitamin A, first dose",
					"Give first dose of IV/IM Ampicillin and Gentamicin",
					"If clouding of the cornea or pus draining from the eye, apply Tetracycline eye ointment",
					"Refer URGENTLY to hospital",
				},
				TreatmentPlan: "Urgent hospital referral with vitamin A and antibiotic pre-treatment",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Refer URGENTLY to hospital",
				Notes:        "Measles now or within last 3 months AND (any general danger sign OR clouding of cornea OR deep/extensive mouth ulcers)",
			},
			"MEASLES_WITH_EYE_MOUTH_COMPLICATIONS": {
				Classification: "MEASLES WITH EYE OR MOUTH COMPLICATIONS",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give Vitamin A, therapeutic dose",
					"If pus draining from the eye, treat eye infection with Tetracycline eye ointment",
					"If mouth ulcers, treat with gentian violet",
				},
				TreatmentPlan: "Vitamin A supplementation and local treatment for complications",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up after 2 days",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "Measles now or within last 3 months AND (pus draining from eye OR mouth ulcers not deep/extensive)",
			},
			"MEASLES_NO_COMPLICATIONS": {
				Classification: "MEASLES",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Give Vitamin A, therapeutic dose",
				},
				TreatmentPlan: "Vitamin A supplementation",
				FollowUp: []string{
					"Advise mother when to return immediately",
				},
				MotherAdvice: "Advise mother when to return immediately",
				Notes:        "Measles now or within last 3 months with no complications",
			},
		},
	}
}