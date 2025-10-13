// ruleengine/engine/developmental_assessment_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetDevelopmentalAssessmentTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "developmental_assessment",
		Title:        "Developmental Milestones Assessment",
		Instructions: "Assess child's developmental milestones. Skip if child has severe classification.",
		StartNode:    "check_severe_classification",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "check_severe_classification",
				Question:     "Does the child have any severe classification?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "If the infant has severe classification, don't do the assessment of development",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "SEVERE_CLASSIFICATION_NO_ASSESSMENT",
					},
					"no": {
						NextNode: "child_age_months",
					},
				},
			},
			{
				NodeID:       "child_age_months",
				Question:     "What is the child's current age in months?",
				QuestionType: "number_input",
				Required:     true,
				Level:        2,
				ParentNode:   "check_severe_classification",
				Instructions: "Ask child's current age",
				Answers: map[string]domain.Answer{
					"value_based": {
						NextNode: "risk_factors_present",
					},
				},
			},
			{
				NodeID:       "risk_factors_present",
				Question:     "Are there any risk factors that can affect how this child is developing?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "child_age_months",
				Instructions: "Consider medical, environmental, and social risk factors",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "medical_risk_factors",
					},
					"no": {
						NextNode: "parental_concerns",
					},
				},
			},
			{
				NodeID:       "medical_risk_factors",
				Question:     "Select medical risk factors present:",
				QuestionType: "multiple_choice",
				Required:     false,
				Level:        4,
				ParentNode:   "risk_factors_present",
				Instructions: "Risk factors: Difficult birth or any neonatal admission, Prematurity or low birth weight, Malnutrition, Head circumference issues, HIV or exposure to HIV, Serious infection or illness",
				Answers: map[string]domain.Answer{
					"difficult_birth": {
						NextNode: "environmental_risk_factors",
					},
					"prematurity_low_bw": {
						NextNode: "environmental_risk_factors",
					},
					"malnutrition": {
						NextNode: "environmental_risk_factors",
					},
					"head_circumference": {
						NextNode: "environmental_risk_factors",
					},
					"hiv_exposure": {
						NextNode: "environmental_risk_factors",
					},
					"serious_infection": {
						NextNode: "environmental_risk_factors",
					},
					"none": {
						NextNode: "environmental_risk_factors",
					},
				},
			},
			{
				NodeID:       "environmental_risk_factors",
				Question:     "Select environmental risk factors present:",
				QuestionType: "multiple_choice",
				Required:     false,
				Level:        5,
				ParentNode:   "medical_risk_factors",
				Instructions: "Environmental factors: Very young/elderly caregiver, substance abuse, maternal depression, violence/neglect, lack of caregiver responsiveness, poverty",
				Answers: map[string]domain.Answer{
					"young_elderly_caregiver": {
						NextNode: "parental_concerns",
					},
					"substance_abuse": {
						NextNode: "parental_concerns",
					},
					"maternal_depression": {
						NextNode: "parental_concerns",
					},
					"violence_neglect": {
						NextNode: "parental_concerns",
					},
					"lack_responsiveness": {
						NextNode: "parental_concerns",
					},
					"poverty": {
						NextNode: "parental_concerns",
					},
					"none": {
						NextNode: "parental_concerns",
					},
				},
			},
			{
				NodeID:       "parental_concerns",
				Question:     "Do you have any concerns about your child's development?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        6,
				ParentNode:   "",
				Instructions: "How do you think your child is developing? Consider parental concerns when observing the child's development",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "assess_milestones",
					},
					"no": {
						NextNode: "assess_milestones",
					},
				},
			},
			{
				NodeID:       "assess_milestones",
				Question:     "For infants less than 2 months: Assess birth milestones",
				QuestionType: "milestone_assessment",
				Required:     true,
				Level:        7,
				ParentNode:   "parental_concerns",
				Instructions: "Look: For all infants less than two months; assess if the infant achieved the at birth development milestones",
				Answers: map[string]domain.Answer{
					"all_achieved": {
						Classification: "NO_DEVELOPMENTAL_DELAY",
					},
					"one_missing": {
						Classification: "SUSPECTED_DEVELOPMENTAL_DELAY",
					},
					"multiple_missing": {
						Classification: "SUSPECTED_DEVELOPMENTAL_DELAY",
					},
				},
			},
			{
				NodeID:       "milestone_flexed_position",
				Question:     "Does the infant remain flexed in supine position?",
				QuestionType: "yes_no",
				Required:     false,
				Level:        8,
				ParentNode:   "assess_milestones",
				Instructions: "At birth milestone: Remains flexed in supine position",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "milestone_grasp_reflex",
					},
					"no": {
						NextNode: "milestone_grasp_reflex",
					},
				},
			},
			{
				NodeID:       "milestone_grasp_reflex",
				Question:     "Does the infant grasp with fingers and toes when touched on palm or sole?",
				QuestionType: "yes_no",
				Required:     false,
				Level:        8,
				ParentNode:   "milestone_flexed_position",
				Instructions: "At birth milestone: Grasps with fingers and toe when touched on the palm or sole",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "milestone_prefers_faces",
					},
					"no": {
						NextNode: "milestone_prefers_faces",
					},
				},
			},
			{
				NodeID:       "milestone_prefers_faces",
				Question:     "Does the infant prefer facial features (looks at faces)?",
				QuestionType: "yes_no",
				Required:     false,
				Level:        8,
				ParentNode:   "milestone_grasp_reflex",
				Instructions: "At birth milestone: Prefers facial features (looks at faces)",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "milestone_suckle_reflex",
					},
					"no": {
						NextNode: "milestone_suckle_reflex",
					},
				},
			},
			{
				NodeID:       "milestone_suckle_reflex",
				Question:     "Does the infant suckle when touched on mouth with finger?",
				QuestionType: "yes_no",
				Required:     false,
				Level:        8,
				ParentNode:   "milestone_prefers_faces",
				Instructions: "At birth milestone: Suckles when touched on mouth with finger",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "milestone_visual_tracking",
					},
					"no": {
						NextNode: "milestone_visual_tracking",
					},
				},
			},
			{
				NodeID:       "milestone_visual_tracking",
				Question:     "Does the infant follow light or moving object and startle to sound?",
				QuestionType: "yes_no",
				Required:     false,
				Level:        8,
				ParentNode:   "milestone_suckle_reflex",
				Instructions: "At birth milestone: Follows light or moving object in line of vision and startle to sound",
				Answers: map[string]domain.Answer{
					"yes": {
						Classification: "AUTO_CLASSIFY_MILESTONES",
					},
					"no": {
						Classification: "AUTO_CLASSIFY_MILESTONES",
					},
				},
			},
		},
		Outcomes: map[string]domain.Outcome{
			"SEVERE_CLASSIFICATION_NO_ASSESSMENT": {
				Classification: "SEVERE CLASSIFICATION - NO DEVELOPMENTAL ASSESSMENT",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Address severe medical condition first",
					"Defer developmental assessment until child is stable",
					"Provide urgent medical care as needed",
				},
				TreatmentPlan: "Priority medical care - developmental assessment deferred",
				FollowUp: []string{
					"Re-assess development after medical condition stabilizes",
					"Follow up for ongoing medical care",
				},
				MotherAdvice: "Your child needs immediate medical attention first. We will assess development once your child is stable.",
				Notes:        "Child has severe classification - developmental assessment not performed",
			},
			"SUSPECTED_DEVELOPMENTAL_DELAY": {
				Classification: "SUSPECTED DEVELOPMENTAL DELAY",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Praise caregiver on milestones achieved",
					"Counsel caregiver on play & communication, responsive caregiving activities to do at home",
					"Screen for other possible causes including malnutrition, TB disease",
					"Advise to return for follow up in 30 days",
				},
				TreatmentPlan: "Developmental monitoring and support",
				FollowUp: []string{
					"Return for follow up in 30 days",
					"Monitor developmental progress",
					"Screen for underlying conditions",
				},
				MotherAdvice: "Your child may need some extra support with development. Practice the activities we discussed daily and return in 30 days for follow-up.",
				Notes:        "Absence of one or more milestones from current age group",
			},
			"NO_DEVELOPMENTAL_DELAY": {
				Classification: "NO DEVELOPMENTAL DELAY",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Praise caregiver on milestones achieved",
					"Advice the care giver on the importance of responsive caregiving, talking to the child, reading, singing and play with the child on daily basis",
					"Encourage caregiver to exercise more challenging activities of the next age group",
					"Advise to continue with follow up consultations",
					"Share Key message for care giver",
				},
				TreatmentPlan: "Routine developmental promotion",
				FollowUp: []string{
					"Continue with routine follow-up consultations",
					"Monitor developmental progress at next visit",
				},
				MotherAdvice: "Your child is developing well! Continue talking, reading, singing and playing with your child every day to support their ongoing development.",
				Notes:        "All the important milestones for the current age group achieved",
			},
		},
	}
}