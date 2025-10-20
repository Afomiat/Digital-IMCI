package engine

import (
	"testing"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/google/uuid"
)

func TestGetFeedingAssessmentTree(t *testing.T) {
	tree := GetFeedingAssessmentTree()

	// Test basic tree structure
	if tree.AssessmentID != "feeding_assessment" {
		t.Errorf("Expected AssessmentID 'feeding_assessment', got '%s'", tree.AssessmentID)
	}

	if tree.Title != "Feeding Assessment for Children Under 2 Years" {
		t.Errorf("Expected correct title, got '%s'", tree.Title)
	}

	// Test questions count
	expectedQuestions := 20
	if len(tree.QuestionsFlow) != expectedQuestions {
		t.Errorf("Expected %d questions, got %d", expectedQuestions, len(tree.QuestionsFlow))
	}

	// Test outcomes count
	expectedOutcomes := 2
	if len(tree.Outcomes) != expectedOutcomes {
		t.Errorf("Expected %d outcomes, got %d", expectedOutcomes, len(tree.Outcomes))
	}
}

func TestFeedingAssessmentOutcomes(t *testing.T) {
	tree := GetFeedingAssessmentTree()

	// Test FEEDING_PROBLEM outcome
	feedingProblem, exists := tree.Outcomes["FEEDING_PROBLEM"]
	if !exists {
		t.Fatal("FEEDING_PROBLEM outcome not found")
	}

	if feedingProblem.Classification != "FEEDING PROBLEM" {
		t.Errorf("Expected classification 'FEEDING PROBLEM', got '%s'", feedingProblem.Classification)
	}

	if feedingProblem.Color != "yellow" {
		t.Errorf("Expected color 'yellow', got '%s'", feedingProblem.Color)
	}

	if feedingProblem.Emergency != false {
		t.Errorf("Expected Emergency to be false, got %v", feedingProblem.Emergency)
	}

	// Test NO_FEEDING_PROBLEM outcome
	noFeedingProblem, exists := tree.Outcomes["NO_FEEDING_PROBLEM"]
	if !exists {
		t.Fatal("NO_FEEDING_PROBLEM outcome not found")
	}

	if noFeedingProblem.Classification != "NO FEEDING PROBLEM" {
		t.Errorf("Expected classification 'NO FEEDING PROBLEM', got '%s'", noFeedingProblem.Classification)
	}

	if noFeedingProblem.Color != "green" {
		t.Errorf("Expected color 'green', got '%s'", noFeedingProblem.Color)
	}

	if noFeedingProblem.Emergency != false {
		t.Errorf("Expected Emergency to be false, got %v", noFeedingProblem.Emergency)
	}
}

func TestFeedingAssessmentQuestionFlow(t *testing.T) {
	tree := GetFeedingAssessmentTree()

	// Test start node
	if tree.StartNode != "breastfeeding_check" {
		t.Errorf("Expected start node 'breastfeeding_check', got '%s'", tree.StartNode)
	}

	// Test key questions exist
	questionMap := make(map[string]bool)
	for _, q := range tree.QuestionsFlow {
		questionMap[q.NodeID] = true
	}

	expectedQuestions := []string{
		"breastfeeding_check", "breastfeeding_frequency", "night_breastfeeding",
		"other_food_check", "other_food_types", "food_quantity", "food_frequency",
		"feeding_method", "replacement_milk_check", "replacement_milk_type",
		"replacement_frequency", "replacement_quantity", "milk_preparation",
		"utensil_cleaning", "mam_specific_check", "serving_size", "own_serving",
		"feeding_person", "feeding_changes", "finalize_feeding_classification",
	}

	for _, expectedQ := range expectedQuestions {
		if !questionMap[expectedQ] {
			t.Errorf("Expected question '%s' not found", expectedQ)
		}
	}

	// Test feeding method question has correct answers
	var feedingMethodQuestion *domain.Question
	for _, q := range tree.QuestionsFlow {
		if q.NodeID == "feeding_method" {
			feedingMethodQuestion = &q
			break
		}
	}

	if feedingMethodQuestion == nil {
		t.Fatal("feeding_method question not found")
	}

	expectedMethods := []string{"cup", "bottle", "both", "other"}
	for _, method := range expectedMethods {
		if _, exists := feedingMethodQuestion.Answers[method]; !exists {
			t.Errorf("Expected feeding method '%s' not found", method)
		}
	}
}

func TestClassifyFeedingAssessment(t *testing.T) {
	engine := &ChildRuleEngine{}

	testCases := []struct {
		name           string
		answers        map[string]interface{}
		expectedResult string
	}{
		{
			name: "No Feeding Problem - Good Breastfeeding",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "porridge",
				"food_quantity":           "medium",
				"food_frequency":          3,
				"feeding_method":          "cup",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "no",
				"feeding_changes":         "no",
			},
			expectedResult: "NO_FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Infrequent Breastfeeding",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 4, // Less than 6
				"night_breastfeeding":     "yes",
				"other_food_check":        "no",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "no",
				"feeding_changes":         "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - No Night Breastfeeding",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "no", // Problem sign
				"other_food_check":        "no",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "no",
				"feeding_changes":         "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Bottle Feeding",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "milk",
				"food_quantity":           "medium",
				"food_frequency":          4,
				"feeding_method":          "bottle", // Problem sign
				"replacement_milk_check":  "no",
				"mam_specific_check":      "no",
				"feeding_changes":         "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Diluted Milk",
			answers: map[string]interface{}{
				"breastfeeding_check":    "no",
				"other_food_check":       "yes",
				"other_food_types":       "milk", // Problem sign
				"food_quantity":          "small",
				"food_frequency":         4,
				"feeding_method":         "cup",
				"replacement_milk_check": "no",
				"mam_specific_check":     "no",
				"feeding_changes":        "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Infrequent Complementary Food",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "porridge",
				"food_quantity":           "medium",
				"food_frequency":          2, // Less than 3
				"feeding_method":          "cup",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "no",
				"feeding_changes":         "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Inappropriate Replacement Milk",
			answers: map[string]interface{}{
				"breastfeeding_check":    "no",
				"other_food_check":       "no",
				"replacement_milk_check": "yes",
				"replacement_milk_type":  "condensed_milk", // Problem sign
				"replacement_frequency":  6,
				"replacement_quantity":   "1 cup",
				"milk_preparation":       "according_to_instructions",
				"utensil_cleaning":       "washed_with_soap",
				"mam_specific_check":     "no",
				"feeding_changes":        "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Insufficient Replacement Feeds",
			answers: map[string]interface{}{
				"breastfeeding_check":    "no",
				"other_food_check":       "no",
				"replacement_milk_check": "yes",
				"replacement_milk_type":  "infant_formula",
				"replacement_frequency":  4, // Less than 6
				"replacement_quantity":   "1 cup",
				"milk_preparation":       "according_to_instructions",
				"utensil_cleaning":       "washed_with_soap",
				"mam_specific_check":     "no",
				"feeding_changes":        "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Incorrect Milk Preparation",
			answers: map[string]interface{}{
				"breastfeeding_check":    "no",
				"other_food_check":       "no",
				"replacement_milk_check": "yes",
				"replacement_milk_type":  "infant_formula",
				"replacement_frequency":  6,
				"replacement_quantity":   "1 cup",
				"milk_preparation":       "diluted_with_water", // Problem sign
				"utensil_cleaning":       "washed_with_soap",
				"mam_specific_check":     "no",
				"feeding_changes":        "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Unhygienic Preparation",
			answers: map[string]interface{}{
				"breastfeeding_check":    "no",
				"other_food_check":       "no",
				"replacement_milk_check": "yes",
				"replacement_milk_type":  "infant_formula",
				"replacement_frequency":  6,
				"replacement_quantity":   "1 cup",
				"milk_preparation":       "according_to_instructions",
				"utensil_cleaning":       "not_cleaned_properly", // Problem sign
				"mam_specific_check":     "no",
				"feeding_changes":        "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - MAM Child Small Servings",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "porridge",
				"food_quantity":           "medium",
				"food_frequency":          4,
				"feeding_method":          "cup",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "yes",
				"serving_size":            "small", // Problem sign
				"own_serving":             "yes",
				"feeding_person":          "mother_feeds_child",
				"feeding_changes":         "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - MAM Child No Own Serving",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "porridge",
				"food_quantity":           "medium",
				"food_frequency":          4,
				"feeding_method":          "cup",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "yes",
				"serving_size":            "normal",
				"own_serving":             "no", // Problem sign
				"feeding_person":          "mother_feeds_child",
				"feeding_changes":         "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - MAM Child Self Feeding",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "porridge",
				"food_quantity":           "medium",
				"food_frequency":          4,
				"feeding_method":          "cup",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "yes",
				"serving_size":            "normal",
				"own_serving":             "yes",
				"feeding_person":          "child_feeds_self", // Problem sign
				"feeding_changes":         "no",
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "Feeding Problem - Feeding Changed During Illness",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "porridge",
				"food_quantity":           "medium",
				"food_frequency":          4,
				"feeding_method":          "cup",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "no",
				"feeding_changes":         "yes", // Problem sign
			},
			expectedResult: "FEEDING_PROBLEM",
		},
		{
			name: "No Feeding Problem - Good Replacement Milk",
			answers: map[string]interface{}{
				"breastfeeding_check":    "no",
				"other_food_check":       "no",
				"replacement_milk_check": "yes",
				"replacement_milk_type":  "infant_formula",
				"replacement_frequency":  6,
				"replacement_quantity":   "1 cup",
				"milk_preparation":       "according_to_instructions",
				"utensil_cleaning":       "washed_with_soap",
				"mam_specific_check":     "no",
				"feeding_changes":        "no",
			},
			expectedResult: "NO_FEEDING_PROBLEM",
		},
		{
			name: "No Feeding Problem - Good MAM Child",
			answers: map[string]interface{}{
				"breastfeeding_check":     "yes",
				"breastfeeding_frequency": 8,
				"night_breastfeeding":     "yes",
				"other_food_check":        "yes",
				"other_food_types":        "porridge",
				"food_quantity":           "medium",
				"food_frequency":          4,
				"feeding_method":          "cup",
				"replacement_milk_check":  "no",
				"mam_specific_check":      "yes",
				"serving_size":            "normal",
				"own_serving":             "yes",
				"feeding_person":          "mother_feeds_child",
				"feeding_changes":         "no",
			},
			expectedResult: "NO_FEEDING_PROBLEM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := engine.classifyFeedingAssessment(tc.answers)
			if result != tc.expectedResult {
				t.Errorf("Expected '%s', got '%s'", tc.expectedResult, result)
			}
		})
	}
}

func TestFeedingAssessmentIntegration(t *testing.T) {
	// Test engine registration
	engine, err := NewChildRuleEngine()
	if err != nil {
		t.Fatalf("Failed to create child rule engine: %v", err)
	}

	// Test tree availability
	tree, err := engine.GetAssessmentTree("feeding_assessment")
	if err != nil {
		t.Fatalf("Failed to get feeding assessment tree: %v", err)
	}

	if tree.AssessmentID != "feeding_assessment" {
		t.Errorf("Expected tree ID 'feeding_assessment', got '%s'", tree.AssessmentID)
	}

	// Test complete workflow
	flow, err := engine.StartAssessmentFlow(uuid.New(), "feeding_assessment")
	if err != nil {
		t.Fatalf("Failed to start feeding assessment flow: %v", err)
	}

	if flow.TreeID != "feeding_assessment" {
		t.Errorf("Expected flow tree ID 'feeding_assessment', got '%s'", flow.TreeID)
	}

	// Test classification with sample answers
	answers := map[string]interface{}{
		"breastfeeding_check":     "yes",
		"breastfeeding_frequency": 8,
		"night_breastfeeding":     "yes",
		"other_food_check":        "yes",
		"other_food_types":        "porridge",
		"food_quantity":           "1 cup",
		"food_frequency":          4,
		"feeding_method":          "cup",
		"replacement_milk_check":  "no",
		"mam_specific_check":      "no",
		"feeding_changes":         "no",
	}

	classification := engine.classifyFeedingAssessment(answers)
	if classification != "NO_FEEDING_PROBLEM" {
		t.Errorf("Expected 'NO_FEEDING_PROBLEM', got '%s'", classification)
	}
}

func TestParseInt(t *testing.T) {
	engine := &ChildRuleEngine{}

	testCases := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{"Int", 42, 42},
		{"Float64", 42.7, 42},
		{"String", "42", 42},
		{"Invalid String", "not_a_number", 0},
		{"Nil", nil, 0},
		{"Bool", true, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := engine.parseInt(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %d, got %d", tc.expected, result)
			}
		})
	}
}
