package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildEarProblemTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "child_ear_problem",
		Title:        "Check for Ear Problems",
		Instructions: "Assess for ear pain, discharge, and signs of infection. Check for mastoiditis.",
		StartNode:    "ear_pain",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "ear_pain",
				Question:     "Is there ear pain?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "Ask about ear pain",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "pus_draining",
					},
					"no": {
						NextNode: "pus_draining",
					},
				},
			},
			{
				NodeID:       "pus_draining",
				Question:     "Is pus seen draining from the ear?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "ear_pain",
				ShowCondition: "ear_pain.yes OR ear_pain.no",
				Instructions: "Look for pus draining from the ear",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "discharge_duration",
					},
					"no": {
						NextNode: "tender_swelling",
					},
				},
			},
			{
				NodeID:       "discharge_duration",
				Question:     "For how long has there been discharge?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "pus_draining",
				ShowCondition: "pus_draining.yes",
				Instructions: "Determine duration of ear discharge",
				Answers: map[string]domain.Answer{
					"less_than_14_days": {
						NextNode: "tender_swelling",
					},
					"14_days_or_more": {
						NextNode: "tender_swelling",
					},
				},
			},
			{
				NodeID:       "tender_swelling",
				Question:     "Is there tender swelling behind the ear?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "pus_draining",
				ShowCondition: "pus_draining.yes OR pus_draining.no",
				Instructions: "Feel for tender swelling behind the ear",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "MASTOIDITIS",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"no": {
						Classification: "CLASSIFY_BY_SYMPTOMS",
						Color:          "white",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"MASTOIDITIS": {
				Classification: "MASTOIDITIS",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give first dose of Ceftriaxone IV/IM OR Ampicillin and Chloramphenicol IV/IM",
					"Give first dose of Paracetamol for pain",
					"Refer URGENTLY to hospital",
				},
				TreatmentPlan: "Urgent antibiotic treatment and hospital referral",
				FollowUp: []string{
					"Refer URGENTLY to hospital",
				},
				MotherAdvice: "Go to hospital immediately - this is a serious infection",
				Notes:        "Tender swelling behind ear indicates mastoiditis",
			},
			"ACUTE_EAR_INFECTION": {
				Classification: "ACUTE EAR INFECTION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give Amoxicillin for 5 days",
					"Give Paracetamol for pain",
					"Dry the ear by wicking",
				},
				TreatmentPlan: "Antibiotic treatment and pain management",
				FollowUp: []string{
					"Follow-up in 5 days",
				},
				MotherAdvice: "Complete full course of antibiotics and keep ear dry",
				Notes:        "Ear pain OR pus draining for less than 14 days",
			},
			"CHRONIC_EAR_INFECTION": {
				Classification: "CHRONIC EAR INFECTION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Dry the ear by wicking",
					"Treat with topical Quinolone eardrops for 2 weeks",
				},
				TreatmentPlan: "Topical antibiotic eardrops and ear drying",
				FollowUp: []string{
					"Follow-up in 5 days",
				},
				MotherAdvice: "Use eardrops as prescribed and keep ear dry",
				Notes:        "Pus draining for 14 days or more",
			},
			"NO_EAR_INFECTION": {
				Classification: "NO EAR INFECTION",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"No additional treatment",
				},
				TreatmentPlan: "No treatment needed",
				FollowUp: []string{
					"Routine follow-up",
				},
				MotherAdvice: "No specific advice needed",
				Notes:        "No ear pain and no pus seen draining",
			},
			"CLASSIFY_BY_SYMPTOMS": {
				Classification: "CLASSIFY BY SYMPTOMS",
				Color:          "white",
				Emergency:      false,
				Actions: []string{
					"Classify based on symptoms presented",
				},
				TreatmentPlan: "Symptom-based classification",
				FollowUp: []string{
					"Based on final classification",
				},
				MotherAdvice: "Based on final classification",
				Notes:        "Intermediate classification - determine final based on symptoms",
			},
		},
	}
}