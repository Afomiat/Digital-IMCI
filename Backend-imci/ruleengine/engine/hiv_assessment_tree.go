// ruleengine/engine/hiv_assessment_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetHIVAssessmentTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "hiv_status_assessment",
		Title:        "HIV Status Assessment and Classification",
		Instructions: "Assess HIV status of mother and young infant. Classify based on test results.",
		StartNode:    "mother_hiv_status",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "mother_hiv_status",
				Question:     "What is the HIV status of the mother?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "infant_antibody_status",
					},
					"negative": {
						Classification: "HIV_INFECTION_UNLIKELY",
					},
					"unknown": {
						NextNode: "infant_antibody_status",
					},
				},
			},
			{
				NodeID:       "infant_antibody_status",
				Question:     "What is the HIV antibody status of the young infant?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "mother_hiv_status",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "infant_dna_pcr_status",
					},
					"negative": {
						NextNode: "infant_dna_pcr_status",
					},
					"unknown": {
						NextNode: "infant_dna_pcr_status",
					},
				},
			},
			{
				NodeID:       "infant_dna_pcr_status",
				Question:     "What is the DNA PCR test result of the young infant?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "infant_antibody_status",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"positive": {
						Classification: "HIV_INFECTED",
						Color:          "pink",
						EmergencyPath:  true,
					},
					"negative": {
						NextNode: "breastfeeding_status",
					},
					"unknown": {
						Classification: "HIV_EXPOSED",
						Color:          "yellow",
					},
				},
			},
			{
				NodeID:       "breastfeeding_status",
				Question:     "Is the infant breastfeeding now?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "infant_dna_pcr_status",
				Instructions: "",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "HIV_EXPOSED",
						Color:          "yellow",
					},
					"no": {
						Classification: "HIV_INFECTION_UNLIKELY",
						Color:          "green",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"HIV_INFECTED": {
				Classification: "HIV INFECTED",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Start Cotrimoxazole Prophylaxis from 6 weeks of age",
					"Assess feeding and counsel",
					"Assess for TB infection",
					"Refer / Link to ART clinic for immediate ART initiation and other care",
					"Ensure mother is tested and enrolled for HIV care, treatment and follow up",
				},
				TreatmentPlan: "Immediate ART initiation and comprehensive HIV care",
				FollowUp: []string{
					"Regular follow-up at ART clinic",
					"Monitor treatment adherence",
					"Ensure mother receives HIV care",
				},
				MotherAdvice: "Infant requires immediate antiretroviral treatment. Ensure regular clinic visits.",
				Notes:        "DNA PCR positive confirms HIV infection",
			},
			"HIV_EXPOSED": {
				Classification: "HIV EXPOSED",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Start Cotrimoxazole Prophylaxis from 6 weeks of age",
					"Assess feeding and counsel",
					"If DNA PCR test is unknown, test as soon as possible starting from 6 weeks of age",
					"Ensure both mother and baby are enrolled in mother-baby cohort follow up at ANC/PMTCT clinic",
					"Ensure provisions of other components of care",
				},
				TreatmentPlan: "HIV exposed infant care and monitoring",
				FollowUp: []string{
					"Regular DNA PCR testing until definitive status",
					"Monitor growth and development",
					"Mother-baby cohort follow-up",
				},
				MotherAdvice: "Infant is HIV-exposed and requires regular monitoring. Continue with PMTCT services.",
				Notes:        "Mother HIV positive with infant DNA PCR unknown/negative but breastfeeding",
			},
			"HIV_STATUS_UNKNOWN": {
				Classification: "HIV STATUS UNKNOWN",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Initiate HIV testing and counselling",
					"Conduct HIV test for the mother and if positive, a virological test for the infant",
					"Conduct virological test for the infant if mother is not available (Eg orphan)",
				},
				TreatmentPlan: "HIV testing and status determination",
				FollowUp: []string{
					"Return for test results",
					"Enroll in PMTCT services if positive",
				},
				MotherAdvice: "HIV status needs to be determined through testing.",
				Notes:        "Mother and infant HIV status unknown",
			},
			"HIV_INFECTION_UNLIKELY": {
				Classification: "HIV INFECTION UNLIKELY",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Advise on home care of infant",
					"Assess feeding and counsel",
					"Advise the mother on HIV prevention",
				},
				TreatmentPlan: "Routine infant care with HIV prevention",
				FollowUp: []string{
					"Routine child health services",
					"HIV testing if new risk factors emerge",
				},
				MotherAdvice: "Continue with routine infant care and practice HIV prevention measures.",
				Notes:        "No evidence of HIV infection or exposure",
			},
		},
	}
}