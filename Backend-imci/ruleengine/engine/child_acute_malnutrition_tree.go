// ruleengine/engine/acute_malnutrition_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetAcuteMalnutritionTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "acute_malnutrition",
		Title:        "Assess Acute Malnutrition in Children 6 months to 5 years",
		Instructions: "LOOK AND FEEL: Check for pitting oedema, measure WFL/H Z-score and MUAC. If signs of malnutrition present, assess for complications and perform appetite test if needed.",
		StartNode:    "pitting_edema",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "pitting_edema",
				Question:     "Look for pitting oedema on both feet",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "LOOK AND FEEL: Check for pitting oedema on both feet. Grade as +, ++, or +++",
				Answers: map[string]domain.Answer{
					"none": {
						NextNode: "wfl_z_score",
					},
					"plus": {
						NextNode: "wfl_z_score",
						Color:    "yellow",
					},
					"plus_plus": {
						NextNode: "wfl_z_score",
						Color:    "yellow",
					},
					"plus_plus_plus": {
						NextNode: "wfl_z_score",
						Color:    "pink",
					},
				},
			},
			{
				NodeID:       "wfl_z_score",
				Question:     "Determine Weight-for-Length/Height Z-score",
				QuestionType: "number_input",
				Required:     true,
				Level:        2,
				ParentNode:   "pitting_edema",
				Instructions: "MEASURE: Determine WFL/H Z-score using growth charts",
				Validation: &domain.Validation{
					Min:  -5,
					Max:  5,
					Step: 0.1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "muac_measurement",
					},
				},
			},
			{
				NodeID:       "muac_measurement",
				Question:     "Measure Mid-Upper Arm Circumference (MUAC)",
				QuestionType: "number_input",
				Required:     true,
				Level:        3,
				ParentNode:   "wfl_z_score",
				Instructions: "MEASURE: Measure MUAC in centimeters. If WFL/H and MUAC measurements are discordant, use the worse measurement for classification.",
				Validation: &domain.Validation{
					Min:  8,
					Max:  20,
					Step: 0.1,
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "malnutrition_risk_assessment",
					},
				},
			},
			{
				NodeID:       "medical_complications_multi",
				Question:     "Select medical complications present (if any)",
				QuestionType: "multiple_choice",
				Required:     false,
				Level:        4,
				ParentNode:   "muac_measurement",
				Instructions: "Select all that apply: Any General Danger Sign; Any severe classification; Pneumonia; Dehydration; Persistent diarrhea; Dysentery; Measles (now/eye/mouth complications); Fever ≥38.5°C; Low body temperature <35°C; Dermatosis +++; Vitamin A deficiency eye signs",
				Answers: map[string]domain.Answer{
					"any_present": {
						NextNode: "finalize_classification",
					},
					"none": {
						NextNode: "appetite_test",
					},
				},
			},
			{
				NodeID:        "severe_wasting_with_edema_check",
				Question:      "Does the child have severe wasting with oedema?",
				QuestionType:  "yes_no",
				Required:      true,
				Level:         5,
				ParentNode:    "medical_complications_multi",
				ShowCondition: "medical_complications_multi.none",
				Instructions:  "ASSESS: Child with WFH <-3Z plus oedema, OR with MUAC<11.5cm plus oedema",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "COMPLICATED_SEVERE_ACUTE_MALNUTRITION",
						Color:          "pink",
					},
					"no": {
						NextNode: "medical_complications_check",
					},
				},
			},
			{
				NodeID:        "appetite_test",
				Question:      "Perform appetite test",
				QuestionType:  "single_choice",
				Required:      true,
				Level:         7,
				ParentNode:    "severe_wasting_with_edema_check",
				ShowCondition: "severe_wasting_with_edema_check.no AND medical_complications_multi.none",
				Instructions:  "TEST: Offer RUTF or F-75. If child takes adequate amount, appetite test is PASSED. If child refuses or takes very little, appetite test is FAILED.",
				Answers: map[string]domain.Answer{
					"passed": {
						NextNode: "finalize_classification",
					},
					"failed": {
						NextNode: "finalize_classification",
					},
				},
			},
			{
				NodeID:       "finalize_classification",
				Question:     "Finalize classification",
				QuestionType: "single_choice",
				Required:     true,
				Level:        8,
				ParentNode:   "appetite_test",
				Instructions: "System will compute classification based on measurements, oedema, complications, and appetite test.",
				Answers: map[string]domain.Answer{
					"compute": {
						Classification: "AUTO_CLASSIFY",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"COMPLICATED_SEVERE_ACUTE_MALNUTRITION": {
				Classification: "COMPLICATED SEVERE ACUTE MALNUTRITION",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Admit to inpatient care (Stabilization Center) or Refer urgently to hospital",
					"Give 1st dose of Ampicillin and Gentamicin IM",
					"Treat the child to prevent low blood sugar",
					"Advise the mother to feed and keep the child warm",
					"Advise mother on the need of referral",
				},
				TreatmentPlan: "Inpatient care with stabilization and antibiotic treatment",
				FollowUp: []string{
					"Refer urgently to hospital",
					"Monitor in stabilization center",
				},
				MotherAdvice: "This is an emergency. Your child needs immediate hospital care. Go to the hospital now.",
				Notes:        "WFL/H <-3Z OR MUAC <11.5cm OR oedema of both feet (+, ++) with medical complications OR failed appetite test OR +++ oedema OR severe wasting with oedema",
			},
			"UNCOMPLICATED_SEVERE_ACUTE_MALNUTRITION": {
				Classification: "UNCOMPLICATED SEVERE ACUTE MALNUTRITION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"If OTP available: Admit child to OTP and follow standard OTP treatment",
					"Give RUTF for 7 days",
					"Give oral Amoxicillin for 5 days",
					"Counsel on how to feed RUTF to the child",
					"Advise when to return immediately",
					"Assess for TB infection",
					"Follow-up in 7 days",
					"If no OTP in facility: Refer child to OTP service",
					"If social problem at home: Treat child as in-patient",
				},
				TreatmentPlan: "Outpatient Therapeutic Program (OTP) with RUTF and antibiotics",
				FollowUp: []string{
					"Follow-up in 7 days",
					"Assess for TB infection",
					"Monitor RUTF consumption",
				},
				MotherAdvice: "Your child needs special feeding and care. Follow the treatment plan carefully and return in 7 days.",
				Notes:        "WFL/H <-3Z OR MUAC <11.5cm OR oedema of both feet (+, ++) with no medical complications and passed appetite test",
			},
			"MODERATE_ACUTE_MALNUTRITION": {
				Classification: "MODERATE ACUTE MALNUTRITION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Admit or Refer to Supplementary Feeding Program (TSFP)",
					"Follow TSFP care protocol",
					"Assess for feeding and counsel the mother accordingly",
					"Assess for TB infection",
					"If feeding problem, follow up in 5 days",
					"Follow up in 30 days",
				},
				TreatmentPlan: "Supplementary Feeding Program with nutritional support",
				FollowUp: []string{
					"Follow up in 5 days if feeding problem",
					"Follow up in 30 days",
					"Assess for TB infection",
				},
				MotherAdvice: "Your child needs better nutrition. Follow the feeding advice and return for follow-up as scheduled.",
				Notes:        "WFL/H ≥-3Z to <-2Z OR MUAC 11.5cm to <12.5cm with no oedema of both feet",
			},
			"NO_ACUTE_MALNUTRITION": {
				Classification: "NO ACUTE MALNUTRITION",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Assess feeding and advise the mother on feeding",
					"Follow up in 5 days if feeding problem",
					"If no feeding problem - praise the mother",
				},
				TreatmentPlan: "Nutritional counseling and monitoring",
				FollowUp: []string{
					"Follow up in 5 days if feeding problem",
					"Routine follow-up if no feeding problem",
				},
				MotherAdvice: "Your child's nutrition is good. Continue with current feeding practices.",
				Notes:        "WFL/H ≥-2Z OR MUAC ≥12.5cm with no oedema of both feet",
			},
		},
	}
}
