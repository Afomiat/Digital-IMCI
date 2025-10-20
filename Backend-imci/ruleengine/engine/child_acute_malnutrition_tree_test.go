package engine

import (
	"testing"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
)

func TestGetAcuteMalnutritionTree(t *testing.T) {
	tree := GetAcuteMalnutritionTree()

	// Test basic tree structure
	if tree.AssessmentID != "acute_malnutrition" {
		t.Errorf("Expected AssessmentID 'acute_malnutrition', got '%s'", tree.AssessmentID)
	}

	if tree.StartNode != "pitting_edema" {
		t.Errorf("Expected StartNode 'pitting_edema', got '%s'", tree.StartNode)
	}

	// Test that all expected questions are present
	expectedQuestions := []string{
		"pitting_edema",
		"wfl_z_score",
		"muac_measurement",
		"medical_complications_multi",
		"severe_wasting_with_edema_check",
		"appetite_test",
		"finalize_classification",
	}

	questionMap := make(map[string]bool)
	for _, q := range tree.QuestionsFlow {
		questionMap[q.NodeID] = true
	}

	for _, expectedQ := range expectedQuestions {
		if !questionMap[expectedQ] {
			t.Errorf("Expected question '%s' not found in tree", expectedQ)
		}
	}

	// Test that all expected outcomes are present
	expectedOutcomes := []string{
		"COMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		"UNCOMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		"MODERATE_ACUTE_MALNUTRITION",
		"NO_ACUTE_MALNUTRITION",
	}

	for _, expectedOutcome := range expectedOutcomes {
		if _, exists := tree.Outcomes[expectedOutcome]; !exists {
			t.Errorf("Expected outcome '%s' not found in tree", expectedOutcome)
		}
	}
}

func TestClassifyAcuteMalnutrition(t *testing.T) {
	engine := &ChildRuleEngine{}

	tests := []struct {
		name           string
		answers        map[string]interface{}
		expectedResult string
	}{
		{
			name: "Complicated Severe - +++ Oedema",
			answers: map[string]interface{}{
				"pitting_edema":                   "plus_plus_plus",
				"wfl_z_score":                     -2.0,
				"muac_measurement":                12.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "COMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		},
		{
			name: "Complicated Severe - Severe Wasting with Oedema",
			answers: map[string]interface{}{
				"pitting_edema":                   "plus",
				"wfl_z_score":                     -3.5,
				"muac_measurement":                10.0,
				"severe_wasting_with_edema_check": "yes",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "COMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		},
		{
			name: "Complicated Severe - Medical Complications Present",
			answers: map[string]interface{}{
				"pitting_edema":                   "plus",
				"wfl_z_score":                     -3.2,
				"muac_measurement":                11.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "any_present",
				"appetite_test":                   "passed",
			},
			expectedResult: "COMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		},
		{
			name: "Complicated Severe - Failed Appetite Test",
			answers: map[string]interface{}{
				"pitting_edema":                   "plus",
				"wfl_z_score":                     -3.2,
				"muac_measurement":                11.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "failed",
			},
			expectedResult: "COMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		},
		{
			name: "Uncomplicated Severe - Passed Appetite Test",
			answers: map[string]interface{}{
				"pitting_edema":                   "plus",
				"wfl_z_score":                     -3.2,
				"muac_measurement":                11.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "UNCOMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		},
		{
			name: "Uncomplicated Severe - WFL < -3Z",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -3.5,
				"muac_measurement":                12.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "UNCOMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		},
		{
			name: "Uncomplicated Severe - MUAC < 11.5cm",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -2.0,
				"muac_measurement":                11.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "UNCOMPLICATED_SEVERE_ACUTE_MALNUTRITION",
		},
		{
			name: "Moderate Acute Malnutrition - WFL -2.5Z",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -2.5,
				"muac_measurement":                12.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "MODERATE_ACUTE_MALNUTRITION",
		},
		{
			name: "Moderate Acute Malnutrition - MUAC 12.0cm",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -2.0,
				"muac_measurement":                12.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "MODERATE_ACUTE_MALNUTRITION",
		},
		{
			name: "No Acute Malnutrition - Normal WFL",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -1.5,
				"muac_measurement":                13.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "NO_ACUTE_MALNUTRITION",
		},
		{
			name: "No Acute Malnutrition - Normal MUAC",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -1.0,
				"muac_measurement":                13.5,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "NO_ACUTE_MALNUTRITION",
		},
		{
			name: "Edge Case - WFL exactly -3Z",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -3.0,
				"muac_measurement":                12.0,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "MODERATE_ACUTE_MALNUTRITION",
		},
		{
			name: "Edge Case - MUAC exactly 11.5cm",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -2.0,
				"muac_measurement":                11.5,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "MODERATE_ACUTE_MALNUTRITION",
		},
		{
			name: "Edge Case - MUAC exactly 12.5cm",
			answers: map[string]interface{}{
				"pitting_edema":                   "none",
				"wfl_z_score":                     -2.0,
				"muac_measurement":                12.5,
				"severe_wasting_with_edema_check": "no",
				"medical_complications_multi":     "none",
				"appetite_test":                   "passed",
			},
			expectedResult: "NO_ACUTE_MALNUTRITION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.classifyAcuteMalnutrition(tt.answers)
			if result != tt.expectedResult {
				t.Errorf("Expected '%s', got '%s'", tt.expectedResult, result)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	engine := &ChildRuleEngine{}

	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{"Float64", 3.14, 3.14},
		{"Int", 42, 42.0},
		{"String", "3.14", 3.14},
		{"Invalid String", "invalid", 0},
		{"Nil", nil, 0},
		{"Bool", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.parseFloat(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestAcuteMalnutritionOutcomes(t *testing.T) {
	tree := GetAcuteMalnutritionTree()

	// Test Complicated Severe Acute Malnutrition outcome
	complicated := tree.Outcomes["COMPLICATED_SEVERE_ACUTE_MALNUTRITION"]
	if complicated.Color != "pink" {
		t.Errorf("Expected color 'pink', got '%s'", complicated.Color)
	}
	if !complicated.Emergency {
		t.Error("Expected Emergency to be true")
	}
	if len(complicated.Actions) == 0 {
		t.Error("Expected actions to be present")
	}

	// Test Uncomplicated Severe Acute Malnutrition outcome
	uncomplicated := tree.Outcomes["UNCOMPLICATED_SEVERE_ACUTE_MALNUTRITION"]
	if uncomplicated.Color != "yellow" {
		t.Errorf("Expected color 'yellow', got '%s'", uncomplicated.Color)
	}
	if uncomplicated.Emergency {
		t.Error("Expected Emergency to be false")
	}

	// Test Moderate Acute Malnutrition outcome
	moderate := tree.Outcomes["MODERATE_ACUTE_MALNUTRITION"]
	if moderate.Color != "yellow" {
		t.Errorf("Expected color 'yellow', got '%s'", moderate.Color)
	}
	if moderate.Emergency {
		t.Error("Expected Emergency to be false")
	}

	// Test No Acute Malnutrition outcome
	noMalnutrition := tree.Outcomes["NO_ACUTE_MALNUTRITION"]
	if noMalnutrition.Color != "green" {
		t.Errorf("Expected color 'green', got '%s'", noMalnutrition.Color)
	}
	if noMalnutrition.Emergency {
		t.Error("Expected Emergency to be false")
	}
}

func TestAcuteMalnutritionQuestionFlow(t *testing.T) {
	tree := GetAcuteMalnutritionTree()

	// Test pitting oedema question
	oedemaQuestion := findQuestionByID(tree, "pitting_edema")
	if oedemaQuestion == nil {
		t.Fatal("pitting_edema question not found")
	}

	// Test oedema answer colors
	oedemaAnswers := oedemaQuestion.Answers
	if oedemaAnswers["plus"].Color != "yellow" {
		t.Error("Expected plus oedema to be yellow")
	}
	if oedemaAnswers["plus_plus"].Color != "yellow" {
		t.Error("Expected plus_plus oedema to be yellow")
	}
	if oedemaAnswers["plus_plus_plus"].Color != "pink" {
		t.Error("Expected plus_plus_plus oedema to be pink")
	}

	// Test medical complications multi-select
	complicationsQuestion := findQuestionByID(tree, "medical_complications_multi")
	if complicationsQuestion == nil {
		t.Fatal("medical_complications_multi question not found")
	}
	if complicationsQuestion.QuestionType != "multiple_choice" {
		t.Error("Expected medical_complications_multi to be multiple_choice")
	}

	// Test appetite test
	appetiteQuestion := findQuestionByID(tree, "appetite_test")
	if appetiteQuestion == nil {
		t.Fatal("appetite_test question not found")
	}
	if appetiteQuestion.QuestionType != "single_choice" {
		t.Error("Expected appetite_test to be single_choice")
	}
}

func findQuestionByID(tree *domain.AssessmentTree, nodeID string) *domain.Question {
	for _, question := range tree.QuestionsFlow {
		if question.NodeID == nodeID {
			return &question
		}
	}
	return nil
}

// Integration test for the complete flow
func TestAcuteMalnutritionIntegration(t *testing.T) {
	// This would test the complete flow from start to finish
	// including the engine registration and usecase integration

	childEngine, err := NewChildRuleEngine()
	if err != nil {
		t.Fatalf("Failed to create child engine: %v", err)
	}

	// Test that the acute malnutrition tree is registered
	tree, err := childEngine.GetAssessmentTree("acute_malnutrition")
	if err != nil {
		t.Fatalf("Failed to get acute malnutrition tree: %v", err)
	}

	if tree.AssessmentID != "acute_malnutrition" {
		t.Errorf("Expected AssessmentID 'acute_malnutrition', got '%s'", tree.AssessmentID)
	}

	// Test that the tree is in available trees
	availableTrees := childEngine.GetAvailableTrees()
	found := false
	for _, treeID := range availableTrees {
		if treeID == "acute_malnutrition" {
			found = true
			break
		}
	}
	if !found {
		t.Error("acute_malnutrition tree not found in available trees")
	}
}
