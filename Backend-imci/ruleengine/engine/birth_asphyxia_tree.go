// ruleengine/engine/birth_asphyxia_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetBirthAsphyxiaTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "birth_asphyxia_check",
		Title:        "Check for Birth Asphyxia",
		Instructions: "Assess and check for Birth Asphyxia while drying and wrapping with dry cloth. No crying is considered as no breathing. Complete resuscitation within the 1st minute - 'Golden Minute' of life",
		StartNode:    "check_birth_asphyxia",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "check_birth_asphyxia",
				Question:     "Check for Birth Asphyxia?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode:      "not_breathing",
						Color:         "pink",
						Action:        "Start resuscitation immediately - Golden Minute",
						EmergencyPath: true,
					},
					"no": {
						NextNode: "assessment_complete_normal",
						Color:    "green",
					},
				},
			},
			{
				NodeID:       "not_breathing",
				Question:     "Is baby not breathing?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "check_birth_asphyxia",
				ShowCondition: "check_birth_asphyxia.yes",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "BIRTH_ASPHYXIA",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "gasping",
					},
				},
			},
			{
				NodeID:       "gasping",
				Question:     "Is baby gasping?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "check_birth_asphyxia",
				ShowCondition: "check_birth_asphyxia.yes AND not_breathing.no",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "BIRTH_ASPHYXIA",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "breathing_poorly",
					},
				},
			},
			{
				NodeID:       "breathing_poorly",
				Question:     "Is baby breathing poorly (<30 breaths/minute)?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "check_birth_asphyxia",
				ShowCondition: "check_birth_asphyxia.yes AND not_breathing.no AND gasping.no",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "BIRTH_ASPHYXIA",
						EmergencyPath:  true,
					},
					"no": {
						NextNode: "breathing_normally",
					},
				},
			},
			{
				NodeID:       "breathing_normally",
				Question:     "Is baby breathing normally (crying or ≥30 breaths/minute)?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        2,
				ParentNode:   "check_birth_asphyxia",
				ShowCondition: "check_birth_asphyxia.yes AND not_breathing.no AND gasping.no AND breathing_poorly.no",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "NO_BIRTH_ASPHYXIA",
						Color:          "green",
					},
					"no": {
						Classification: "BIRTH_ASPHYXIA",
						EmergencyPath:  true,
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"BIRTH_ASPHYXIA": {
				Classification: "BIRTH ASPHYXIA",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Clear mouth first, then nose with bulb syringe",
					"Clamp/tie and cut the cord immediately",
					"Position the newborn supine with neck slightly extended",
					"Ventilate with appropriate size bag & mask",
				},
				TreatmentPlan: "Complete resuscitation within 1st minute - Golden Minute",
				FollowUp: []string{
					"Follow after 12 hrs",
					"24 hrs (in the facility)",
					"3 days",
					"7 days",
					"6 weeks",
				},
				MotherAdvice: "Advise mother when to return immediately",
			},
			"NO_BIRTH_ASPHYXIA": {
				Classification: "NO BIRTH ASPHYXIA",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Give cord care",
					"Initiate skin-to-skin contact",
					"Initiate breastfeeding",
					"Give eye care",
					"Give Vitamin K",
					"Apply Chlorhexidine gel",
					"Give HepB, BCG and OPV 0",
				},
				TreatmentPlan: "Essential Newborn Care",
				FollowUp: []string{
					"Follow after 6 hrs (in the facility)",
					"3 days",
					"7 days",
					"6 weeks",
				},
				MotherAdvice: "Advise mother when to return immediately",
			},
		},
	}
}