package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/google/uuid"
)

type RuleEngine struct {
	trees map[string]*domain.DecisionTree
}

func NewRuleEngine() (*RuleEngine, error) {
	engine := &RuleEngine{
		trees: make(map[string]*domain.DecisionTree),
	}

	// Load decision trees
	if err := engine.loadTrees(); err != nil {
		return nil, fmt.Errorf("failed to load decision trees: %w", err)
	}

	return engine, nil
}

func (e *RuleEngine) loadTrees() error {
	// Get the project root directory (one level up from cmd)
	projectRoot, err := e.getProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to get project root: %w", err)
	}

	treeFiles := map[string]string{
		"child":        filepath.Join(projectRoot, "ruleengine", "data", "child_tree.json"),
		"young_infant": filepath.Join(projectRoot, "ruleengine", "data", "young_infant_tree.json"),
	}

	fmt.Printf("🔍 Project root: %s\n", projectRoot)
	
	for treeType, filePath := range treeFiles {
		fmt.Printf("📁 Loading tree: %s from %s\n", treeType, filePath)
		
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("tree file does not exist: %s", filePath)
		}
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read tree file %s: %w", filePath, err)
		}

		var tree domain.DecisionTree
		if err := json.Unmarshal(data, &tree); err != nil {
			return fmt.Errorf("failed to parse tree file %s: %w", filePath, err)
		}

		e.trees[treeType] = &tree
		fmt.Printf("✅ Successfully loaded %s tree\n", treeType)
	}

	return nil
}

func (e *RuleEngine) getProjectRoot() (string, error) {
	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0) // Gets the path of this source file
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// This file is in: Backend-imci/ruleengine/engine/rule_engine.go
	// We want: Backend-imci/
	engineDir := filepath.Dir(filename)                    // ruleengine/engine
	ruleengineDir := filepath.Dir(engineDir)               // ruleengine
	projectRoot := filepath.Dir(ruleengineDir)             // Backend-imci/
	
	// Verify the project root contains both cmd and ruleengine directories
	cmdPath := filepath.Join(projectRoot, "cmd")
	ruleenginePath := filepath.Join(projectRoot, "ruleengine")
	
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		return "", fmt.Errorf("cmd directory not found in project root: %s", cmdPath)
	}
	if _, err := os.Stat(ruleenginePath); os.IsNotExist(err) {
		return "", fmt.Errorf("ruleengine directory not found in project root: %s", ruleenginePath)
	}
	
	fmt.Printf("✅ Project structure verified:\n")
	fmt.Printf("   - Project root: %s\n", projectRoot)
	fmt.Printf("   - CMD dir: %s\n", cmdPath)
	fmt.Printf("   - RuleEngine dir: %s\n", ruleenginePath)
	
	return projectRoot, nil
}

func (e *RuleEngine) StartSession(assessmentID uuid.UUID, assessmentType string) (*domain.AssessmentSession, error) {
	tree, exists := e.trees[assessmentType]
	if !exists {
		return nil, fmt.Errorf("no decision tree found for type: %s", assessmentType)
	}

	session := &domain.AssessmentSession{
		SessionID:       uuid.New(),
		AssessmentID:    assessmentID,
		AssessmentType:  assessmentType,
		CurrentNodeID:   tree.StartNode,
		Answers:         make(map[string]interface{}),
		ClinicalFindings: &domain.ClinicalFindings{},
	}

	return session, nil
}

// ruleengine/engine/rule_engine.go
func getTreeTypeFromAssessment(session *domain.AssessmentSession) string {
    // Directly use the stored assessment type
    return session.AssessmentType
}

func (e *RuleEngine) GetCurrentNode(session *domain.AssessmentSession) (*domain.DecisionNode, error) {
    treeType := getTreeTypeFromAssessment(session)
    
    fmt.Printf("🔍 Getting current node - Tree Type: %s, Node ID: %s\n", treeType, session.CurrentNodeID)
    
    tree, exists := e.trees[treeType]
    if !exists {
        return nil, fmt.Errorf("decision tree not found for type: %s", treeType)
    }

    node, exists := tree.Nodes[session.CurrentNodeID]
    if !exists {
        return nil, fmt.Errorf("node not found: %s", session.CurrentNodeID)
    }

    fmt.Printf("🔍 Node Found - Type: %s, ID: %s\n", node.Type, node.ID)
    if node.Type == domain.NodeTypeClassification {
        fmt.Printf("🔍 Classification Node Details: %+v\n", node.Classification)
    }

    return &node, nil
}
func (e *RuleEngine) SubmitAnswer(session *domain.AssessmentSession, answer interface{}) (*domain.Classification, error) {
    fmt.Printf("🚨 ========== SUBMIT ANSWER - START ==========\n")
    fmt.Printf("🔍 Current Node ID: %s\n", session.CurrentNodeID)
    
    currentNode, err := e.GetCurrentNode(session)
    if err != nil {
        fmt.Printf("❌ Error getting current node: %v\n", err)
        return nil, err
    }

    fmt.Printf("🔍 Current Node Type: %s\n", currentNode.Type)
    fmt.Printf("🔍 Current Node Question: %s\n", currentNode.Question)
    fmt.Printf("🔍 Answer Type: %T, Value: %v\n", answer, answer)

    // Store answer
    session.Answers[currentNode.ID] = answer

    // Update clinical findings
    e.updateClinicalFindings(session, currentNode, answer)

    // Find next node
    nextNodeID, err := e.evaluateRules(session, currentNode, answer)
    if err != nil {
        fmt.Printf("❌ Error evaluating rules: %v\n", err)
        return nil, err
    }

    fmt.Printf("🔍 Next Node ID: %s\n", nextNodeID)

    // Move to next node
    session.CurrentNodeID = nextNodeID

    // Check if we reached a classification
    nextNode, err := e.GetCurrentNode(session)
    if err != nil {
        fmt.Printf("❌ Error getting next node: %v\n", err)
        return nil, err
    }

    fmt.Printf("🔍 Next Node Type: %s\n", nextNode.Type)
    fmt.Printf("🔍 Next Node ID: %s\n", nextNode.ID)

    if nextNode.Type == domain.NodeTypeClassification {
        fmt.Printf("🎯 REACHED CLASSIFICATION NODE!\n")
        
        // FIX: Check if Classification is not nil before accessing it
        if nextNode.Classification == nil {
            fmt.Printf("❌ ERROR: Classification is nil for node: %s\n", nextNode.ID)
            return nil, fmt.Errorf("classification data missing for node: %s", nextNode.ID)
        }
        
        fmt.Printf("🎯 Classification Name: %s\n", nextNode.Classification.Name)
        fmt.Printf("🎯 Classification Color: %s\n", nextNode.Classification.Color)
        fmt.Printf("🎯 Treatment Plan: %+v\n", nextNode.Classification.TreatmentPlan)
        
        session.Classification = nextNode.Classification
        fmt.Printf("✅ ========== SUBMIT ANSWER - COMPLETE WITH CLASSIFICATION ==========\n")
        return nextNode.Classification, nil
    }

    fmt.Printf("➡️  Continuing to next question: %s\n", nextNode.ID)
    fmt.Printf("✅ ========== SUBMIT ANSWER - CONTINUING ==========\n")
    return nil, nil
}

func (e *RuleEngine) evaluateRules(session *domain.AssessmentSession, node *domain.DecisionNode, answer interface{}) (string, error) {
    fmt.Printf("🚨 ========== RULE ENGINE DEBUG START ==========\n")
    fmt.Printf("🔍 Node: %s\n", node.ID)
    fmt.Printf("🔍 Answer Type: %T\n", answer)
    fmt.Printf("🔍 Answer Value: %v\n", answer)
    
	 if node.Type == domain.NodeTypeClassification {
        fmt.Printf("🎯 ALREADY AT CLASSIFICATION NODE - No rules to evaluate\n")
        fmt.Printf("✅ ========== RULE ENGINE DEBUG END ==========\n")
        return node.ID, nil // Stay on the same node
    }
    var selected []string
    var numericValue float64
    var hasNumeric bool
    
    // Handle different answer types based on node type
    switch node.Type {
    case domain.NodeTypeMultipleChoice:
        switch v := answer.(type) {
        case []string:
            selected = v
        case []interface{}:
            for _, item := range v {
                if str, ok := item.(string); ok {
                    selected = append(selected, str)
                }
            }
        }
    case domain.NodeTypeYesNo:
        if str, ok := answer.(string); ok {
            selected = []string{str}
        }
    case domain.NodeTypeAssessment:
        // Handle numeric assessments
        if num, ok := answer.(float64); ok {
            numericValue = num
            hasNumeric = true
        } else if num, ok := answer.(int); ok {
            numericValue = float64(num)
            hasNumeric = true
        }
    }

    fmt.Printf("🔍 Number of Rules: %d\n", len(node.Rules))
    for i, rule := range node.Rules {
        fmt.Printf("   Rule %d:\n", i+1)
        fmt.Printf("      Condition: '%s'\n", rule.Condition)
        fmt.Printf("      Selected Options: %v\n", rule.SelectedOptions)
		fmt.Printf("      Value: '%s'\n", rule.Value)  // ← ADD THIS LINE

        fmt.Printf("      Thresholds: %v\n", rule.Thresholds)
        fmt.Printf("      Next Node: %s\n", rule.NextNode)
    }

    for i, rule := range node.Rules {
        fmt.Printf("\n🎯 Processing Rule %d: '%s'\n", i+1, rule.Condition)
        
        switch rule.Condition {
        case "any_selected":
            fmt.Printf("   Processing 'any_selected' condition\n")
            fmt.Printf("   Selected options: %v\n", selected)
            fmt.Printf("   Looking for any of: %v\n", rule.SelectedOptions)
            
            for _, selectedOpt := range selected {
                for _, requiredOpt := range rule.SelectedOptions {
                    if selectedOpt == requiredOpt {
                        fmt.Printf("   ✅ MATCH FOUND: '%s' in required options\n", selectedOpt)
                        fmt.Printf("   ➡️  Moving to: %s\n", rule.NextNode)
                        fmt.Printf("✅ ========== RULE ENGINE DEBUG END ==========\n")
                        return rule.NextNode, nil
                    }
                }
            }
            fmt.Printf("   ❌ No matches found\n")

        case "none_selected":
            fmt.Printf("   Processing 'none_selected' condition\n")
            fmt.Printf("   Selected options: %v\n", selected)
            fmt.Printf("   Looking for exactly: %v\n", rule.SelectedOptions)
            
            // For "none_selected", we expect exactly the specified options and no others
            if len(selected) == len(rule.SelectedOptions) {
                allMatch := true
                for i, expectedOpt := range rule.SelectedOptions {
                    if i >= len(selected) || selected[i] != expectedOpt {
                        allMatch = false
                        break
                    }
                }
                if allMatch {
                    fmt.Printf("   ✅ EXACT MATCH: Selected options match required 'none' condition\n")
                    fmt.Printf("   ➡️  Moving to: %s\n", rule.NextNode)
                    fmt.Printf("✅ ========== RULE ENGINE DEBUG END ==========\n")
                    return rule.NextNode, nil
                }
            }
            fmt.Printf("   ❌ No exact match for 'none_selected' condition\n")

        case "all_selected":
            fmt.Printf("   Processing 'all_selected' condition\n")
            fmt.Printf("   Selected options: %v\n", selected)
            fmt.Printf("   Looking for all of: %v\n", rule.SelectedOptions)
            
            allFound := true
            for _, requiredOpt := range rule.SelectedOptions {
                found := false
                for _, selectedOpt := range selected {
                    if selectedOpt == requiredOpt {
                        found = true
                        break
                    }
                }
                if !found {
                    allFound = false
                    break
                }
            }
            if allFound {
                fmt.Printf("   ✅ ALL REQUIRED OPTIONS FOUND\n")
                fmt.Printf("   ➡️  Moving to: %s\n", rule.NextNode)
                fmt.Printf("✅ ========== RULE ENGINE DEBUG END ==========\n")
                return rule.NextNode, nil
            }
            fmt.Printf("   ❌ Not all required options found\n")

        case "yes", "no":
            fmt.Printf("   Processing '%s' condition\n", rule.Condition)
            if len(selected) > 0 && selected[0] == rule.Condition {
                fmt.Printf("   ✅ YES/NO MATCH: %s\n", rule.Condition)
                fmt.Printf("   ➡️  Moving to: %s\n", rule.NextNode)
                fmt.Printf("✅ ========== RULE ENGINE DEBUG END ==========\n")
                return rule.NextNode, nil
            }
            fmt.Printf("   ❌ Yes/No condition not met\n")

        case "above_threshold":
            if hasNumeric {
                fmt.Printf("   Processing 'above_threshold' condition\n")
                fmt.Printf("   Value: %.2f, Thresholds: %v\n", numericValue, rule.Thresholds)
                
                if threshold, exists := rule.Thresholds["value"]; exists {
                    if numericValue > float64(threshold) {
                        fmt.Printf("   ✅ ABOVE THRESHOLD: %.2f > %d\n", numericValue, threshold)
                        fmt.Printf("   ➡️  Moving to: %s\n", rule.NextNode)
                        return rule.NextNode, nil
                    }
                }
            }

        case "below_threshold":
            if hasNumeric {
                fmt.Printf("   Processing 'below_threshold' condition\n")
                fmt.Printf("   Value: %.2f, Thresholds: %v\n", numericValue, rule.Thresholds)
                
                if threshold, exists := rule.Thresholds["value"]; exists {
                    if numericValue < float64(threshold) {
                        fmt.Printf("   ✅ BELOW THRESHOLD: %.2f < %d\n", numericValue, threshold)
                        fmt.Printf("   ➡️  Moving to: %s\n", rule.NextNode)
                        return rule.NextNode, nil
                    }
                }
            }

        case "within_range":
            if hasNumeric {
                fmt.Printf("   Processing 'within_range' condition\n")
                fmt.Printf("   Value: %.2f, Thresholds: %v\n", numericValue, rule.Thresholds)
                
                if min, max := rule.Thresholds["min"], rule.Thresholds["max"]; min > 0 && max > 0 {
                    if numericValue >= float64(min) && numericValue <= float64(max) {
                        fmt.Printf("   ✅ WITHIN RANGE: %d ≤ %.2f ≤ %d\n", min, numericValue, max)
                        fmt.Printf("   ➡️  Moving to: %s\n", rule.NextNode)
                        return rule.NextNode, nil
                    }
                }
            }

        default:
            fmt.Printf("   ❌ UNKNOWN CONDITION TYPE: '%s'\n", rule.Condition)
        }
    }

    fmt.Printf("❌ ========== NO MATCHING RULE FOUND ==========\n")
    fmt.Printf("❌ Final Answer: %v\n", selected)
    fmt.Printf("❌ ========== RULE ENGINE DEBUG END ==========\n")
    return "", fmt.Errorf("no clinically appropriate rule found for assessment findings")
}

func (e *RuleEngine) updateClinicalFindings(session *domain.AssessmentSession, node *domain.DecisionNode, answer interface{}) {
    var selected []string
    var numericValue float64
    var hasNumeric bool
    var stringValue string
    
    // Handle different answer types
    switch node.Type {
    case domain.NodeTypeMultipleChoice:
        switch v := answer.(type) {
        case []string:
            selected = v
        case []interface{}:
            for _, item := range v {
                if str, ok := item.(string); ok {
                    selected = append(selected, str)
                }
            }
        }
    case domain.NodeTypeYesNo:
        if str, ok := answer.(string); ok {
            selected = []string{str}
        }
    case domain.NodeTypeAssessment:
        // Handle numeric assessments
        if num, ok := answer.(float64); ok {
            numericValue = num
            hasNumeric = true
        } else if num, ok := answer.(int); ok {
            numericValue = float64(num)
            hasNumeric = true
        } else if str, ok := answer.(string); ok {
            stringValue = str
        }
    }

    // Update findings based on node type and answers
    switch node.ID {
    
    // === GENERAL DANGER SIGNS ===
    case "check_danger_signs":
        for _, opt := range selected {
            switch opt {
            case "unable_to_drink":
                session.ClinicalFindings.UnableToDrink = true
            case "vomits_everything":
                session.ClinicalFindings.VomitsEverything = true
            case "had_convulsions":
                session.ClinicalFindings.HadConvulsions = true
            case "lethargic_unconscious":
                session.ClinicalFindings.LethargicUnconscious = true
            case "convulsing_now":
                session.ClinicalFindings.ConvulsingNow = true
            }
        }
    
    // === RESPIRATORY ASSESSMENT ===
    case "ask_about_cough":
        if len(selected) > 0 {
            session.ClinicalFindings.CoughPresent = (selected[0] == "yes")
        }
    
    case "check_cough_duration":
        if hasNumeric {
            days := int(numericValue)
            session.ClinicalFindings.CoughDurationDays = &days
        }
    
    case "check_respiratory_rate":
        if hasNumeric {
            rate := int(numericValue)
            session.ClinicalFindings.RespiratoryRate = &rate
            
            // Determine fast breathing based on age (IMCI 2021 thresholds)
            ageMonths := getAgeFromSession(session)
            threshold := getFastBreathingThreshold(ageMonths)
            session.ClinicalFindings.FastBreathing = (rate >= threshold)
        }
    
    case "check_chest_indrawing":
        if len(selected) > 0 {
            session.ClinicalFindings.ChestIndrawing = (selected[0] == "yes")
        }
    
    case "check_stridor":
        if len(selected) > 0 {
            session.ClinicalFindings.Stridor = (selected[0] == "yes")
        }
    
    case "check_wheezing":
        if len(selected) > 0 {
            session.ClinicalFindings.Wheezing = (selected[0] == "yes")
        }
    
    case "check_oxygen_saturation":
        if hasNumeric {
            saturation := int(numericValue)
            session.ClinicalFindings.OxygenSaturation = &saturation
            // IMCI 2021: SpO2 < 90% is severe pneumonia
            session.ClinicalFindings.FastBreathing = session.ClinicalFindings.FastBreathing || (saturation < 90)
        }
    
    // === DIARRHEA ASSESSMENT ===
    case "ask_about_diarrhea":
        if len(selected) > 0 {
            session.ClinicalFindings.DiarrheaPresent = (selected[0] == "yes")
        }
    
    case "check_diarrhea_duration":
        if hasNumeric {
            days := int(numericValue)
            session.ClinicalFindings.DiarrheaDurationDays = &days
            // IMCI 2021: Diarrhea ≥14 days is persistent
        }
    
    case "check_blood_in_stool":
        if len(selected) > 0 {
            session.ClinicalFindings.BloodInStool = (selected[0] == "yes")
        }
    
    case "check_dehydration_signs":
        for _, opt := range selected {
            switch opt {
            case "sunken_eyes":
                session.ClinicalFindings.SunkenEyes = true
            case "skin_pinch_slow":
                session.ClinicalFindings.SkinPinchSlow = true
            case "skin_pinch_very_slow":
                session.ClinicalFindings.SkinPinchVerySlow = true
            case "restless_irritable":
                session.ClinicalFindings.RestlessIrritable = true
            case "drinking_eagerly":
                session.ClinicalFindings.DrinkingEagerly = true
            case "drinking_poorly":
                session.ClinicalFindings.DrinkingPoorly = true
            case "lethargic_unconscious": // Can also indicate severe dehydration
                session.ClinicalFindings.LethargicUnconscious = true
            }
        }
    
    // === FEVER ASSESSMENT ===
    case "ask_about_fever":
        if len(selected) > 0 {
            session.ClinicalFindings.FeverPresent = (selected[0] == "yes")
        }
    
    case "check_fever_duration":
        if hasNumeric {
            days := int(numericValue)
            session.ClinicalFindings.FeverDurationDays = &days
        }
    
    case "check_stiff_neck":
        if len(selected) > 0 {
            session.ClinicalFindings.StiffNeck = (selected[0] == "yes")
        }
    
    case "check_bulging_fontanelle":
        if len(selected) > 0 {
            session.ClinicalFindings.BulgingFontanelle = (selected[0] == "yes")
        }
    
    case "check_measles_signs":
        for _, opt := range selected {
            switch opt {
            case "generalized_rash":
                session.ClinicalFindings.GeneralizedRash = true
            case "runny_nose":
                session.ClinicalFindings.RunnyNose = true
            case "red_eyes":
                session.ClinicalFindings.RedEyes = true
            case "cough":
                session.ClinicalFindings.CoughPresent = true
            case "measles_now":
                session.ClinicalFindings.MeaslesNow = true
            case "measles_last_3_months":
                session.ClinicalFindings.MeaslesLast3Months = true
            }
        }
    
    // === EAR PROBLEMS ===
    case "ask_about_ear_pain":
        if len(selected) > 0 {
            session.ClinicalFindings.EarPain = (selected[0] == "yes")
        }
    
    case "check_ear_discharge":
        if len(selected) > 0 {
            session.ClinicalFindings.EarDischarge = (selected[0] == "yes")
        }
    
    case "check_ear_discharge_duration":
        if hasNumeric {
            days := int(numericValue)
            session.ClinicalFindings.EarDischargeDurationDays = &days
            // IMCI 2021: Discharge ≥14 days is chronic
        }
    
    case "check_tender_swelling_behind_ear":
        if len(selected) > 0 {
            session.ClinicalFindings.TenderSwellingBehindEar = (selected[0] == "yes")
        }
    
    // === NUTRITION ASSESSMENT ===
    case "check_muac":
        if hasNumeric {
            muac := numericValue
            session.ClinicalFindings.MUAC = &muac
            // IMCI 2021: MUAC < 11.5 cm is severe wasting
        }
    
    case "check_edema":
        if len(selected) > 0 {
            session.ClinicalFindings.BilateralEdema = (selected[0] == "yes")
        }
    
    case "check_wasting":
        if len(selected) > 0 {
            session.ClinicalFindings.VisibleSevereWasting = (selected[0] == "yes")
        }
    
    // === ANEMIA ASSESSMENT ===
    case "check_palmar_pallor":
        for _, opt := range selected {
            switch opt {
            case "severe_pallor":
                session.ClinicalFindings.SeverePalmarPallor = true
            case "some_pallor":
                session.ClinicalFindings.SomePalmarPallor = true
            }
        }
    
    case "check_hb_level":
        if hasNumeric {
            hb := numericValue
            session.ClinicalFindings.HbLevel = &hb
            // IMCI 2021: Hb < 7 g/dL is severe anemia
        }
    
    // === JAUNDICE ASSESSMENT ===
    case "check_jaundice":
        for _, opt := range selected {
            switch opt {
            case "palms_soles_yellow":
                session.ClinicalFindings.PalmsSolesYellow = true
            case "skin_eyes_yellow":
                session.ClinicalFindings.SkinEyesYellow = true
            }
        }
    
    case "check_jaundice_age":
        if hasNumeric {
            hours := int(numericValue)
            session.ClinicalFindings.JaundiceAgeHours = &hours
            // IMCI 2021: Jaundice <24h or ≥14 days is severe
        }
    
    // === YOUNG INFANT SPECIFIC ===
    case "check_very_severe_disease_young_infant":
        for _, opt := range selected {
            switch opt {
            case "unable_to_feed":
                session.ClinicalFindings.UnableToFeed = true
            case "not_feeding_well":
                session.ClinicalFindings.NotFeedingWell = true
            case "convulsing_now":
                session.ClinicalFindings.ConvulsingNow = true
            case "no_movement":
                session.ClinicalFindings.NoMovement = true
            case "movement_only_when_stimulated":
                session.ClinicalFindings.MovementOnlyWhenStimulated = true
            case "fast_breathing":
                session.ClinicalFindings.FastBreathing = true
            case "severe_chest_indrawing":
                session.ClinicalFindings.ChestIndrawing = true
            case "fever":
                session.ClinicalFindings.FeverPresent = true
            case "low_temperature":
                session.ClinicalFindings.LowBodyTemperature = true
            }
        }
    
    case "check_local_infection_young_infant":
        for _, opt := range selected {
            switch opt {
            case "umbilicus_red":
                session.ClinicalFindings.UmbilicusRed = true
            case "umbilicus_draining_pus":
                session.ClinicalFindings.UmbilicusDrainingPus = true
            case "skin_pustules":
                session.ClinicalFindings.SkinPustules = true
            }
        }
    
    case "check_temperature_young_infant":
        if hasNumeric {
            temp := numericValue
            session.ClinicalFindings.BodyTemperature = &temp
            session.ClinicalFindings.LowBodyTemperature = (temp < 35.5)
            session.ClinicalFindings.FeverPresent = (temp >= 37.5)
        }
    
    // === HIV/TB ASSESSMENT ===
    case "check_hiv_exposure":
        if len(selected) > 0 {
            session.ClinicalFindings.HIVExposed = (selected[0] == "yes")
        }
    
    case "check_tb_symptoms":
        for _, opt := range selected {
            switch opt {
            case "cough_14_days":
                if session.ClinicalFindings.CoughDurationDays != nil {
                    session.ClinicalFindings.TBCoughDurationDays = session.ClinicalFindings.CoughDurationDays
                }
            case "weight_loss":
                session.ClinicalFindings.TBWeightLoss = true
            case "night_sweats":
                session.ClinicalFindings.NightSweats = true
            case "tb_contact":
                session.ClinicalFindings.TBContactHistory = true
            }
        }
    
    // === DEVELOPMENT ASSESSMENT ===
    case "check_development_milestones":
        if len(selected) > 0 {
            session.ClinicalFindings.MilestonesAbsent = selected
            session.ClinicalFindings.SuspectedDevelopmentalDelay = len(selected) > 0
        }
    
    case "check_development_risk_factors":
        if len(selected) > 0 {
            session.ClinicalFindings.RiskFactorsPresent = selected
        }
    
    // === FEEDING ASSESSMENT ===
    case "check_breastfeeding":
        if len(selected) > 0 {
            session.ClinicalFindings.Breastfeeding = (selected[0] == "yes")
        }
    
    case "check_breastfeeding_frequency":
        if hasNumeric {
            freq := int(numericValue)
            session.ClinicalFindings.BreastfeedingFrequency = &freq
        }
    
    case "check_complementary_foods":
        if len(selected) > 0 {
            session.ClinicalFindings.ComplementaryFoods = (selected[0] == "yes")
        }
    
    case "check_feeding_problem":
        if len(selected) > 0 {
            session.ClinicalFindings.FeedingProblem = (selected[0] == "yes")
        }
    
    case "check_underweight":
        if len(selected) > 0 {
            session.ClinicalFindings.Underweight = (selected[0] == "yes")
        }
    
    // === OTHER FINDINGS ===
    case "other_findings":
        if stringValue != "" {
            session.ClinicalFindings.OtherFindings = stringValue
        }
    }
}

// Helper functions
func getAgeFromSession(session *domain.AssessmentSession) int {
    // This should retrieve age from your assessment data
    // For now, return a default - you'll need to implement this based on your data structure
    return 12 // Default to 12 months
}

func getFastBreathingThreshold(ageMonths int) int {
    // IMCI 2021 thresholds
    if ageMonths < 2 {
        return 60 // Young infants: ≥60 breaths/min
    } else if ageMonths < 12 {
        return 50 // Infants 2-12 months: ≥50 breaths/min
    }
    return 40 // Children 12-59 months: ≥40 breaths/min
}

func intPtr(i int) *int {
    return &i
}

func floatPtr(f float64) *float64 {
    return &f
}