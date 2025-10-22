package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildDevelopmentalAssessmentTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "developmental_assessment",
		Title:        "Child Development Assessment",
		Instructions: "ASSESS: Child's developmental milestones, risk factors, and parental concerns to classify developmental status",
		StartNode:    "severe_classification_check",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "severe_classification_check",
				Question:     "Does the child have any severe classification from other assessments?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "CHECK: If child has severe classification, don't do the assessment of development. Complete other assessments first.",
				Options: []domain.Option{
					{Value: "yes", DisplayText: "Yes, child has severe classification"},
					{Value: "no", DisplayText: "No severe classification"},
				},
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "",
					},
					"no": {
						NextNode: "child_age_group",
					},
				},
			},
			{
				NodeID:       "child_age_group",
				Question:     "What is the child's current age?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        2,
				ParentNode:   "severe_classification_check",
				Instructions: "ASK: Child's current age in months to determine appropriate milestone assessment",
				Options: []domain.Option{
					{Value: "0_24_months", DisplayText: "0-24 months"},
					{Value: "24_60_months", DisplayText: "24-60 months"},
				},
				Answers: map[string]domain.Answer{
					"0_24_months": {
						NextNode: "risk_factors",
					},
					"24_60_months": {
						NextNode: "risk_factors",
					},
				},
			},
			{
				NodeID:       "risk_factors",
				Question:     "Are there any risk factors that can affect how this child is developing?",
				QuestionType: "multi_choice",
				Required:     true,
				Level:        3,
				ParentNode:   "child_age_group",
				Instructions: "ASK: About biological and environmental risk factors that may impact development",
				Options: []domain.Option{
					{Value: "difficult_birth", DisplayText: "Difficult birth or any neonatal admission"},
					{Value: "prematurity_low_birth_weight", DisplayText: "Prematurity or low birth weight"},
					{Value: "malnutrition", DisplayText: "Malnutrition"},
					{Value: "head_circumference_abnormal", DisplayText: "Head circumference too large or too small"},
					{Value: "hiv_exposure", DisplayText: "HIV or exposure to HIV"},
					{Value: "serious_infection", DisplayText: "Serious infection or illness"},
					{Value: "young_elderly_caregiver", DisplayText: "Very young or elderly caregiver"},
					{Value: "substance_abuse", DisplayText: "Abuse of drugs or alcohol"},
					{Value: "violence_neglect", DisplayText: "Signs of violence or neglect"},
					{Value: "unresponsive_caregiver", DisplayText: "Lack of caregiver responsiveness to the child"},
					{Value: "poverty", DisplayText: "Poverty"},
					{Value: "none", DisplayText: "None of the above"},
				},
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "parental_concerns",
					},
				},
			},
			{
				NodeID:       "parental_concerns",
				Question:     "How do you think your child is developing? Do you have any concerns?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        4,
				ParentNode:   "risk_factors",
				Instructions: "ASK: About parental concerns regarding child's development",
				Options: []domain.Option{
					{Value: "yes", DisplayText: "Yes, I have concerns"},
					{Value: "no", DisplayText: "No concerns"},
				},
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "current_milestones_achieved",
					},
					"no": {
						NextNode: "current_milestones_achieved",
					},
				},
			},
			{
				NodeID:       "current_milestones_achieved",
				Question:     "Has the child achieved all important milestones for their current age group?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        5,
				ParentNode:   "parental_concerns",
				Instructions: "OBSERVE: Test each milestone for the child's current age group. If child is sleeping, shy or too sick, ask the mother if the child does this action at home",
				Options: []domain.Option{
					{Value: "yes", DisplayText: "Yes, all milestones achieved"},
					{Value: "no", DisplayText: "No, some milestones missing"},
				},
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "",
					},
					"no": {
						NextNode: "earlier_milestones_achieved",
					},
				},
			},
			{
				NodeID:       "earlier_milestones_achieved",
				Question:     "Has the child achieved all milestones for the earlier age group?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        6,
				ParentNode:   "current_milestones_achieved",
				Instructions: "OBSERVE: Test milestones for the earlier age group if current age milestones were not achieved",
				Options: []domain.Option{
					{Value: "yes", DisplayText: "Yes, all earlier milestones achieved"},
					{Value: "no", DisplayText: "No, some earlier milestones missing"},
				},
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "regression_signs",
					},
					"no": {
						NextNode: "regression_signs",
					},
				},
			},
			{
				NodeID:       "regression_signs",
				Question:     "Are there any signs of regression in previously achieved milestones?",
				QuestionType: "single_choice",
				Required:     true,
				Level:        7,
				ParentNode:   "earlier_milestones_achieved",
				Instructions: "ASK: About any loss of previously achieved developmental skills",
				Options: []domain.Option{
					{Value: "yes", DisplayText: "Yes, there are regression signs"},
					{Value: "no", DisplayText: "No regression signs"},
				},
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "",
					},
					"no": {
						NextNode: "",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"ASSESSMENT_NOT_APPLICABLE": {
				Classification: "ASSESSMENT NOT APPLICABLE",
				Color:          "gray",
				Emergency:      false,
				Actions: []string{
					"Complete other assessments first",
					"Address severe conditions immediately",
					"Return for developmental assessment when child is stable",
				},
				TreatmentPlan: "Assessment not applicable - complete other assessments first",
				FollowUp: []string{
					"Complete all other IMCI assessments",
					"Address severe conditions",
					"Return for developmental assessment when stable",
				},
				MotherAdvice: "Your child needs immediate attention for other health conditions. We will complete the developmental assessment once your child is stable.",
				Notes:        "Child has severe classification from other assessments - developmental assessment not applicable",
			},
			"CONFIRMED_DEVELOPMENTAL_DELAY": {
				Classification: "CONFIRMED DEVELOPMENTAL DELAY",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Counsel caregiver on play & communication, Responsive caregiving activities to do at home",
					"Refer for psychomotor evaluation",
					"Screen for mothers health needs and risk factors and other possible causes including Malnutrition, TB disease and hyperthyroidism",
					"Advise to continue with follow up consultations",
				},
				TreatmentPlan: "Confirmed Developmental Delay - requires comprehensive evaluation and intervention",
				FollowUp: []string{
					"Immediate referral to psychomotor evaluation",
					"Screen for underlying medical conditions",
					"Regular developmental monitoring",
					"Family counseling and support",
				},
				MotherAdvice: "Your child shows signs of developmental delay that need immediate attention. We will refer you to specialists who can provide comprehensive evaluation and support. Continue with responsive caregiving activities at home.",
				Notes:        "Absence of one or more milestones from current age group AND absence of one or more milestones from earlier age group, OR regression of milestones signs",
			},
			"SUSPECTED_DEVELOPMENTAL_DELAY": {
				Classification: "SUSPECTED DEVELOPMENTAL DELAY",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Praise caregiver on milestones achieved",
					"Counsel caregiver on play & communication, Responsive caregiving activities to do at home",
					"Advise to return for follow up in 30 days",
					"Screen for other possible causes including malnutrition, TB disease",
				},
				TreatmentPlan: "Suspected Developmental Delay - requires monitoring and early intervention",
				FollowUp: []string{
					"Return for follow-up in 30 days",
					"Continue responsive caregiving activities",
					"Monitor for improvement or further delays",
					"Screen for underlying causes",
				},
				MotherAdvice: "Your child may have some developmental concerns that we need to monitor. Continue with the activities we discussed and return in 30 days for re-evaluation. Your child has achieved many milestones which is great!",
				Notes:        "Absence of one or more milestones from current age but has reached all milestones for earlier age, OR if there is risk factors, OR parental concern",
			},
			"NO_DEVELOPMENTAL_DELAY": {
				Classification: "NO DEVELOPMENTAL DELAY",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Praise caregiver on milestones achieved",
					"Advice the caregiver on the importance of responsive caregiving, talking to the child, reading, singing and play with the child on daily basis",
					"Encourage caregiver to exercise more challenging activities of the next age group",
					"Advise to continue with follow up consultations",
					"Share Key message for caregiver",
				},
				TreatmentPlan: "No Developmental Delay - continue with age-appropriate activities and regular monitoring",
				FollowUp: []string{
					"Continue regular child health follow-up",
					"Encourage next age group activities",
					"Maintain responsive caregiving practices",
					"Monitor for any future concerns",
				},
				MotherAdvice: "Excellent! Your child is developing well and has achieved all the important milestones for their age. Continue with the activities we discussed to support their continued development.",
				Notes:        "All the important milestones for the current age group achieved",
			},
		},
	}
}
