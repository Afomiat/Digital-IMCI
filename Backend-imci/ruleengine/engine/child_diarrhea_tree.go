// ruleengine/engine/child_diarrhea_tree.go
package engine

import "github.com/Afomiat/Digital-IMCI/ruleengine/domain"

func GetChildDiarrheaTree() *domain.AssessmentTree {
	return &domain.AssessmentTree{
		AssessmentID: "child_diarrhea",
		Title:        "Assess Child with Diarrhea",
		Instructions: "ASK: Does the child have diarrhea? IF YES, LOOK AND FEEL: Look at the child's general condition, Look for sunken eyes, Check for blood in stool, Assess drinking ability, Pinch skin for elasticity",
		StartNode:    "diarrhea_present",
		QuestionsFlow: []domain.Question{
			{
				NodeID:       "diarrhea_present",
				Question:     "Does the child have diarrhea?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        1,
				ParentNode:   "",
				Instructions: "ASK: Does the child have diarrhea?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "how_long_diarrhea",
					},
					"no": {
						Classification: "NO_DIARRHEA",
						Color:          "green",
					},
				},
			},
			{
				NodeID:       "how_long_diarrhea",
				Question:     "For how long has the child had diarrhea?",
				QuestionType: "number",
				Required:     true,
				Level:        2,
				ParentNode:   "diarrhea_present",
				ShowCondition: "diarrhea_present.yes",
				Instructions: "ASK: For how many days?",
				Answers: map[string]domain.Answer{
					"*": {
						NextNode: "blood_in_stool",
					},
				},
			},
			{
				NodeID:       "blood_in_stool",
				Question:     "Is there blood in the stool?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        3,
				ParentNode:   "how_long_diarrhea",
				ShowCondition: "how_long_diarrhea.*",
				Instructions: "LOOK: Check for blood in the stool",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "lethargic_unconscious",
					},
					"no": {
						NextNode: "lethargic_unconscious",
					},
				},
			},
			{
				NodeID:       "lethargic_unconscious",
				Question:     "Is the child lethargic or unconscious?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "blood_in_stool",
				ShowCondition: "blood_in_stool.*",
				Instructions: "LOOK: Look at the child's general condition - is child lethargic or unconscious?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "sunken_eyes",
					},
					"no": {
						NextNode: "restless_irritable",
					},
				},
			},
			{
				NodeID:       "restless_irritable",
				Question:     "Is the child restless and irritable?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        4,
				ParentNode:   "lethargic_unconscious",
				ShowCondition: "lethargic_unconscious.no",
				Instructions: "LOOK: Look at the child's general condition - is child restless and irritable?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "sunken_eyes",
					},
					"no": {
						NextNode: "sunken_eyes",
					},
				},
			},
			{
				NodeID:       "sunken_eyes",
				Question:     "Are the eyes sunken?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        5,
				ParentNode:   "lethargic_unconscious",
				ShowCondition: "lethargic_unconscious.*",
				Instructions: "LOOK: Look for sunken eyes",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "drinking_ability",
					},
					"no": {
						NextNode: "drinking_ability",
					},
				},
			},
			{
				NodeID:       "drinking_ability",
				Question:     "Is the child able to drink or breastfeed?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        6,
				ParentNode:   "sunken_eyes",
				ShowCondition: "sunken_eyes.*",
				Instructions: "OFFER FLUID: Is the child able to drink or breastfeed?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "drinking_eagerly",
					},
					"no": {
						NextNode: "skin_pinch",
					},
				},
			},
			{
				NodeID:       "drinking_eagerly",
				Question:     "Is the child drinking eagerly and thirsty?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        7,
				ParentNode:   "drinking_ability",
				ShowCondition: "drinking_ability.yes",
				Instructions: "OFFER FLUID: Is the child drinking eagerly and thirsty?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "skin_pinch",
					},
					"no": {
						NextNode: "skin_pinch",
					},
				},
			},
			{
				NodeID:       "skin_pinch",
				Question:     "Does the skin pinch go back very slowly (more than 2 seconds)?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        8,
				ParentNode:   "drinking_ability",
				ShowCondition: "drinking_ability.*",
				Instructions: "PINCH SKIN: Pinch the skin of abdomen - does it go back very slowly (more than 2 seconds)?",
				Answers: map[string]domain.Answer{
					"yes": {
						NextNode: "skin_pinch_slow",
					},
					"no": {
						NextNode: "skin_pinch_slow",
					},
				},
			},
			{
				NodeID:       "skin_pinch_slow",
				Question:     "Does the skin pinch go back slowly?",
				QuestionType: "yes_no",
				Required:     true,
				Level:        9,
				ParentNode:   "skin_pinch",
				ShowCondition: "skin_pinch.no",
				Instructions: "PINCH SKIN: Does the skin pinch go back slowly?",
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
			"NO_DIARRHEA": {
				Classification: "NO DIARRHEA",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Continue with assessment of other symptoms",
				},
				TreatmentPlan: "No treatment needed for diarrhea",
				FollowUp: []string{
					"Assess other symptoms as needed",
				},
				MotherAdvice: "Child has no diarrhea.",
				Notes:        "No diarrhea detected",
			},
			"SEVERE_DEHYDRATION": {
				Classification: "SEVERE DEHYDRATION",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Give fluid for severe dehydration (Plan C)",
					"Advise mother to continue breastfeeding",
				},
				TreatmentPlan: "Urgent fluid replacement and possible referral",
				FollowUp: []string{
					"If child has no other severe classification: treat with Plan C",
					"If child also has another severe classification: Refer URGENTLY to hospital with mother giving frequent sips of ORS on the way",
					"If child is 2 years or older and there is cholera in your area, give antibiotic for cholera",
				},
				MotherAdvice: "Child has severe dehydration. Give fluids urgently and continue breastfeeding. Return immediately if condition worsens.",
				Notes:        "Two or more signs: Lethargic/unconscious, Sunken eyes, Not able to drink, Skin pinch very slow",
			},
			"SOME_DEHYDRATION": {
				Classification: "SOME DEHYDRATION",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Give fluid, Zinc supplements and food for some dehydration (Plan B)",
					"Advise mother to continue breastfeeding",
				},
				TreatmentPlan: "Oral rehydration and zinc supplementation",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up in 5 days if not improving",
					"If child also has a severe classification: Refer URGENTLY to hospital with mother giving frequent sips of ORS on the way",
				},
				MotherAdvice: "Child has some dehydration. Give ORS, zinc supplements and continue feeding. Return immediately if child cannot drink, becomes lethargic, or condition worsens.",
				Notes:        "Two or more signs: Restless/irritable, Sunken eyes, Drinking eagerly/thirsty, Skin pinch slow",
			},
			"NO_DEHYDRATION": {
				Classification: "NO DEHYDRATION",
				Color:          "green",
				Emergency:      false,
				Actions: []string{
					"Give fluid, Zinc supplements and food to treat diarrhea at home (Plan A)",
				},
				TreatmentPlan: "Home management with ORS and zinc",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up in 5 days if not improving",
				},
				MotherAdvice: "Child has diarrhea but no dehydration. Give ORS, zinc supplements and continue feeding. Return immediately if child develops signs of dehydration or condition worsens.",
				Notes:        "Not enough signs to classify as some or severe dehydration",
			},
			"SEVERE_PERSISTENT_DIARRHEA": {
				Classification: "SEVERE PERSISTENT DIARRHOEA",
				Color:          "pink",
				Emergency:      true,
				Actions: []string{
					"Treat dehydration before referral unless the child has another severe classification",
					"Give Vitamin A",
					"Refer to hospital",
				},
				TreatmentPlan: "Hospital referral with vitamin A supplementation",
				FollowUp: []string{
					"Refer to hospital for management",
				},
				MotherAdvice: "Child has severe persistent diarrhoea. Refer to hospital for specialized care. Continue breastfeeding during transport.",
				Notes:        "Diarrhoea 14 days or more with dehydration present",
			},
			"PERSISTENT_DIARRHEA": {
				Classification: "PERSISTENT DIARRHOEA",
				Color:          "yellow",
				Emergency:      false,
				Actions: []string{
					"Advise mother on feeding recommendations for a child who has PERSISTENT DIARRHOEA",
					"Give Vitamin A therapeutic dose", 
					"Give zinc for 10 days",
				},
				TreatmentPlan: "Nutritional management and micronutrient supplementation",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up in 5 days",
				},
				MotherAdvice: "Child has persistent diarrhoea. Give special feeding, vitamin A and zinc. Return immediately if child develops signs of dehydration or condition worsens.",
				Notes:        "Diarrhoea 14 days or more with no dehydration",
			},
			"DYSENTERY": {
				Classification: "DYSENTERY",
				Color:          "yellow", 
				Emergency:      false,
				Actions: []string{
					"Treat for 3 days with Ciprofloxacin",
				},
				TreatmentPlan: "Antibiotic treatment for dysentery",
				FollowUp: []string{
					"Advise mother when to return immediately",
					"Follow-up after 2 days of antibiotic treatment",
				},
				MotherAdvice: "Child has dysentery (blood in stool). Give antibiotics for 3 days. Return immediately if child develops high fever, severe abdominal pain, or condition worsens.",
				Notes:        "Blood in stool present",
			},
		},
	}
}