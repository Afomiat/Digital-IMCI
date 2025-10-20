// ruleengine/engine/child_hiv_assessment_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildHIVAssessmentTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "hiv_assessment",
		Title:        "HIV Infection Classification",
		Instructions: "ASSESS: HIV status of mother and child, test results, breastfeeding status, and clinical signs to classify HIV infection",
		StartNode:    "mother_hiv_status",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "mother_hiv_status",
				Question:     "What is the HIV status of the mother?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "ASK: About mother's HIV test status",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "child_antibody_test",
					},
					"negative": {
						NextNode: "child_antibody_test",
					},
					"unknown": {
						NextNode: "child_antibody_test",
					},
				},
			},
			{
				NodeID:       "child_antibody_test",
				Question:     "What is the HIV antibody test result of the sick child?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "mother_hiv_status",
				Instructions: "ASK: About child's HIV antibody test result",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "child_dna_pcr_test",
					},
					"negative": {
						NextNode: "child_dna_pcr_test",
					},
					"unknown": {
						NextNode: "child_dna_pcr_test",
					},
				},
			},
			{
				NodeID:       "child_dna_pcr_test",
				Question:     "What is the DNA PCR test result of the sick child?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "child_antibody_test",
				Instructions: "ASK: About child's DNA PCR test result",
				Answers: map[string]domain.Answer{
					"positive": {
						NextNode: "child_breastfeeding",
					},
					"negative": {
						NextNode: "child_breastfeeding",
					},
					"unknown": {
						NextNode: "clinical_signs_check",
					},
				},
			},
			{
				NodeID:        "clinical_signs_check",
				Question:      "Does the child have clinical signs for presumptive HIV? (Oral thrush, Severe pneumonia, or Very Severe Disease)",
				QuestionType:  "multi_choice",
				Required:      false,
				Level:         4,
				ParentNode:    "child_dna_pcr_test",
				ShowCondition: "child_dna_pcr_test.unknown AND child_antibody_test.positive",
				Instructions:  "CHECK: For clinical signs of presumptive severe HIV disease",
				Options: []string{
					"oral_thrush",
					"severe_pneumonia", 
					"very_severe_disease",
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "child_breastfeeding",
					},
				},
			},
			{
				NodeID:       "child_breastfeeding",
				Question:     "Is child on breastfeeding?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        5,
				ParentNode:   "child_dna_pcr_test",
				Instructions: "ASK: About current breastfeeding status",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "finalize_hiv_classification",
					},
					"no": {
						NextNode: "breastfed_last_6weeks",
					},
				},
			},
			{
				NodeID:        "breastfed_last_6weeks",
				Question:      "If no, was child breastfed in the last 6 weeks?",
				QuestionType:  "yes_no",
				Required:      true,
				Level:         6,
				ParentNode:    "child_breastfeeding",
				ShowCondition: "child_breastfeeding.no",
				Instructions:  "ASK: About breastfeeding in the last 6 weeks",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "finalize_hiv_classification",
					},
					"no": {
						NextNode: "finalize_hiv_classification",
					},
				},
			},
			{
				NodeID:       "finalize_hiv_classification",
				Question:     "Finalize HIV classification based on all collected information",
				QuestionType: "single_choice",
				Required:     true,
				Level:        7,
				ParentNode:   "child_breastfeeding",
				Instructions: "System will compute HIV classification based on test results, breastfeeding status, and clinical signs",
				Answers: map[string]domain.Answer{
					"compute": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"HIV_INFECTED_DNA_PCR": {
				Classification: "HIV INFECTED",
				Color:          "red",
				Emergency:      true,
				Actions: []string{
					"Give Cotrimoxazole prophylaxis",
					"Assess feeding and counsel",
					"Assess for TB infection ",
					"Ensure mother is tested & enrolled in HIV care & treatment",
					"Advise on home care",
					"Refer/Link to ART clinic for ART initiation and other components of care",
					"Ensure child has appropriate follow up",
				},
				TreatmentPlan: "HIV infected (DNA PCR confirmed) - immediate ART initiation and comprehensive care",
				FollowUp: []string{
					"Immediate referral to ART clinic",
					"Regular HIV care follow-up",
					"Monitor treatment adherence",
					"CD4 and viral load monitoring",
				},
				MotherAdvice: "Your child's DNA PCR test confirms HIV infection. Immediate antiretroviral treatment is essential. Follow all medical advice strictly.",
				Notes:        "Confirmed by DNA PCR positive test",
			},
			"HIV_INFECTED_ANTIBODY": {
				Classification: "HIV INFECTED",
				Color:          "red",
				Emergency:      true,
				Actions: []string{
					"Consider Cotrimoxazole prophylaxis",
					"Assess feeding and counsel",
					"Advise on home care",
					"Refer/Link to ART clinic for ART initiation and other components of care",
					"Ensure mother is tested & enrolled in HIV care & treatment",
				},
				TreatmentPlan: "HIV infected (antibody positive) - ART initiation and comprehensive care",
				FollowUp: []string{
					"Immediate referral to ART clinic",
					"Confirmatory testing if needed",
					"Regular HIV care follow-up",
				},
				MotherAdvice: "Your child's HIV antibody test is positive, indicating HIV infection. Immediate treatment is needed. Follow all medical advice.",
				Notes:        "Confirmed by antibody positive test (definitive in children >18 months)",
			},
			"PRESUMPTIVE_SEVERE_HIV": {
				Classification: "PRESUMPTIVE SEVERE HIV DISEASE",
				Color:          "red",
				Emergency:      true,
				Actions: []string{
					"Give Cotrimoxazole prophylaxis",
					"Assess feeding and counsel",
					"Assess for TB infection",
					"Refer to ART clinic and treat as HIV INFECTED",
					"Ensure mother is tested & enrolled in HIV care",
					"Advise on home care",
					"Ensure child has appropriate follow up",
				},
				TreatmentPlan: "Presumptive severe HIV - treat as HIV infected while arranging confirmatory tests",
				FollowUp: []string{
					"Immediate ART initiation",
					"Arrange DNA PCR confirmatory testing",
					"Close clinical monitoring",
				},
				MotherAdvice: "Your child shows signs of severe HIV disease and needs immediate treatment. We will treat as HIV infected while waiting for confirmatory tests.",
				Notes:        "DNA PCR unavailable + antibody positive + clinical signs present",
			},
			"HIV_EXPOSED": {
				Classification: "HIV EXPOSED",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give Cotrimoxazole prophylaxis",
					"Assess feeding and counsel",
					"Assess for TB infection",
					"If child antibody/DNA PCR test is unknown, test as soon as possible",
					"If child antibody test is negative, repeat 6 weeks after complete cessation of breastfeeding",
					"Ensure both mother and baby are enrolled in mother-baby cohort follow up at ANC/PMTCT clinic",
					"Ensure provisions of other components of care",
				},
				TreatmentPlan: "HIV exposed infant - require prophylaxis and regular monitoring",
				FollowUp: []string{
					"Regular PMTCT follow-up",
					"Repeat HIV testing 6 weeks after breastfeeding cessation",
					"Monitor for seroconversion",
					"Continue cotrimoxazole prophylaxis",
				},
				MotherAdvice: "Your child has been exposed to HIV but infection status is not confirmed. Continue with preventive medications and regular testing as advised.",
				Notes:        "Mother HIV positive + child test negative/unknown + breastfeeding or recent exposure",
			},
			"HIV_STATUS_UNKNOWN": {
				Classification: "HIV STATUS UNKNOWN",
				Color:          "orange",
				Emergency:      false,
				Actions: []string{
					"Counsel the mother for HIV testing for herself & the child",
					"Test the child if mother is not available (e.g., orphan)",
					"Advise the mother to give home care",
					"Assess feeding and counsel",
				},
				TreatmentPlan: "HIV status unknown - require testing and counseling",
				FollowUp: []string{
					"Arrange for HIV testing",
					"Follow-up after test results",
					"Provide post-test counseling",
				},
				MotherAdvice: "We need to determine your and your child's HIV status through testing. This is important for proper healthcare.",
				Notes:        "Incomplete testing information for classification",
			},
			"HIV_INFECTION_UNLIKELY": {
				Classification: "HIV INFECTION UNLIKELY",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Advise on home care",
					"Assess feeding and counsel",
					"Advise on HIV prevention",
					"If mother HIV status is unknown, encourage mother to be tested",
				},
				TreatmentPlan: "HIV infection unlikely - routine care with prevention counseling",
				FollowUp: []string{
					"Routine child health follow-up",
					"HIV prevention counseling",
				},
				MotherAdvice: "Based on current information, HIV infection is unlikely for your child. Continue with good feeding practices and preventive healthcare.",
				Notes:        "Adequate negative testing with no exposure risk",
			},
		},
	}
}