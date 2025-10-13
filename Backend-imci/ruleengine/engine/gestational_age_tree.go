// ruleengine/engine/gestation_classification_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetGestationClassificationTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "gestation_classification",
		Title:        "Gestation and Birth Weight Classification",
		Instructions: "Assess gestational age and birth weight for newborn classification",
		StartNode:    "know_gestational_age",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "know_gestational_age",
				Question:     "Do you know the gestational age?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "gestational_age_weeks",
					},
					"no": {
						NextNode: "know_birth_weight",
					},
				},
			},
			{
				NodeID:       "gestational_age_weeks",
				Question:     "What is the gestational age in weeks?",
				QuestionType: "number_input",
				Required:     true,
				Level:        2,
				ParentNode:   "know_gestational_age",
				Instructions: "Enter number of weeks",
				Answers: map[string]domain.Answer{
					"value_based": { 
						NextNode: "know_birth_weight_ga",
					},
				},
			},
			{
				NodeID:       "know_birth_weight",
				Question:     "Do you know the birth weight?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "know_gestational_age",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "birth_weight_grams",
					},
					"no": {
						NextNode: "can_weigh_baby",
					},
				},
			},
			{
				NodeID:       "know_birth_weight_ga",
				Question:     "Do you know the birth weight?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "gestational_age_weeks",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "birth_weight_grams_ga",
					},
					"no": {
						NextNode: "can_weigh_baby_ga",
					},
				},
			},
			{
				NodeID:       "birth_weight_grams",
				Question:     "What is the birth weight in grams?",
				QuestionType: "number_input",
				Required:     true,
				Level:        3,
				ParentNode:   "know_birth_weight",
				Instructions: "Enter weight in grams",
				Answers: map[string]domain.Answer{
					"value_based": { 
						NextNode: "classify_by_weight_only",
					},
				},
			},
			{
				NodeID:       "birth_weight_grams_ga",
				Question:     "What is the birth weight in grams?",
				QuestionType: "number_input",
				Required:     true,
				Level:        4,
				ParentNode:   "know_birth_weight_ga",
				Instructions: "Enter weight in grams",
				Answers: map[string]domain.Answer{
					"value_based": { 
						NextNode: "classify_by_ga_and_weight",
					},
				},
			},
			{
				NodeID:       "can_weigh_baby",
				Question:     "Can you weigh the baby now?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "know_birth_weight",
				Instructions: "Weigh within 7 days of life",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "current_weight_grams",
					},
					"no": {
						Classification: "WEIGHT_UNKNOWN",
					},
				},
			},
			{
				NodeID:       "can_weigh_baby_ga",
				Question:     "Can you weigh the baby now?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "know_birth_weight_ga",
				Instructions: "Weigh within 7 days of life",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "current_weight_grams_ga",
					},
					"no": {
						NextNode: "classify_by_ga_only",
					},
				},
			},
			{
				NodeID:       "current_weight_grams",
				Question:     "What is the current weight in grams?",
				QuestionType: "number_input",
				Required:     true,
				Level:        4,
				ParentNode:   "can_weigh_baby",
				Instructions: "Enter weight in grams",
				Answers: map[string]domain.Answer{
					"value_based": { 
						NextNode: "classify_by_weight_only",
					},
				},
			},
			{
				NodeID:       "current_weight_grams_ga",
				Question:     "What is the current weight in grams?",
				QuestionType: "number_input",
				Required:     true,
				Level:        5,
				ParentNode:   "can_weigh_baby_ga",
				Instructions: "Enter weight in grams",
				Answers: map[string]domain.Answer{
					"value_based": { 
						NextNode: "classify_by_ga_and_weight",
					},
				},
			},
			{
				NodeID:       "classify_by_weight_only",
				Question:     "AUTO_CLASSIFY",
				QuestionType: "auto_classify",
				Required:     false,
				Level:        4,
				ParentNode:   "",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"auto": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
			{
				NodeID:       "classify_by_ga_and_weight",
				Question:     "AUTO_CLASSIFY",
				QuestionType: "auto_classify",
				Required:     false,
				Level:        5,
				ParentNode:   "",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"auto": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
			{
				NodeID:       "classify_by_ga_only",
				Question:     "AUTO_CLASSIFY",
				QuestionType: "auto_classify",
				Required:     false,
				Level:        5,
				ParentNode:   "can_weigh_baby_ga",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"auto": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"VERY_LOW_BIRTH_WEIGHT": {
				Classification: "VERY LOW BIRTH WEIGHT AND/OR VERY PRETERM",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Continue breastfeeding (if not sucking feed expressed breast milk by cup)",
					"Start Kangaroo Mother Care (KMC)",
					"Give Vitamin K 0.5mg IM on anterior mid lateral thigh, if not already given",
					"Refer URGENTLY with mother to hospital with KMC position",
				},
				TreatmentPlan: "Immediate referral and specialized newborn care",
				FollowUp: []string{
					"Immediate hospital admission",
					"Specialized neonatal care",
				},
				MotherAdvice: "Your baby needs immediate specialized hospital care. Do not delay referral.",
				Notes:        "Weight < 1,500gm OR Gestational age < 32 weeks",
			},
			"LOW_BIRTH_WEIGHT": {
				Classification: "LOW BIRTH WEIGHT AND/OR PRETERM",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"KMC if <2,000gm (in the HF or Hospital)",
					"Counsel on optimal breastfeeding",
					"Counsel mother on prevention of infection",
					"Give Vitamin K 1mg IM (if GA < 34 wks, 0.5 mg IM) on anterior mid lateral thigh if not already given",
					"Provide follow-up for KMC",
					"If baby ≥ 2,000 gms follow-up visits at age 6–24 hrs, 3 days, 7 days & 6 weeks",
					"Give 1st dose of vaccine",
					"Advise mother when to return immediately",
				},
				TreatmentPlan: "Low birth weight care and monitoring",
				FollowUp: []string{
					"Follow-up visits at age 6-24 hrs, 3 days, 7 days & 6 weeks",
					"Monitor weight gain",
					"Kangaroo Mother Care monitoring",
				},
				MotherAdvice: "Your baby needs special care and regular follow-up. Practice Kangaroo Mother Care if weight < 2000g.",
				Notes:        "Weight 1,500 - 2,500 gm OR Gestational age 32-37 weeks",
			},
			"NORMAL_BIRTH_WEIGHT": {
				Classification: "NORMAL BIRTH WEIGHT AND/OR TERM",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Counsel on optimal breastfeeding",
					"Counsel mother/family on prevention of infection",
					"Provide follow-up visits at age 6-24 hrs, 3 days, 7 days & 6 weeks",
					"Give 1st dose of vaccine",
					"Give Vitamin K 1mg IM on anterior mid thigh if not already given",
					"Advise mother when to return immediately",
				},
				TreatmentPlan: "Routine newborn care",
				FollowUp: []string{
					"Follow-up visits at age 6-24 hrs, 3 days, 7 days & 6 weeks",
					"Routine immunization schedule",
				},
				MotherAdvice: "Continue with routine newborn care and breastfeeding. Return for scheduled follow-up visits.",
				Notes:        "Weight ≥ 2,500 gm OR Gestational age ≥ 37 weeks",
			},
			"WEIGHT_UNKNOWN": {
				Classification: "INCOMPLETE ASSESSMENT",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Weigh the baby to complete assessment",
					"Provide basic newborn care meanwhile",
					"Schedule return visit for weight measurement",
				},
				TreatmentPlan: "Need weight measurement for complete classification",
				FollowUp: []string{
					"Return when baby can be weighed",
					"Complete assessment at that time",
				},
				MotherAdvice: "We need to weigh your baby to provide the right care. Please return when you can get the baby weighed.",
				Notes:        "Weight unknown - cannot complete classification",
			},
		},
	}
}