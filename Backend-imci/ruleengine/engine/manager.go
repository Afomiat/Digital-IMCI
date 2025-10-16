// ruleengine/engine/manager.go
package engine

import (
    "errors"
    "fmt"

    "github.com/Afomiat/Digital-IMCI/ruleengine/domain"
    "github.com/google/uuid"
)

var (
    ErrAgeGroupNotSupported = errors.New("age group not supported")
    ErrInvalidAgeGroup      = errors.New("invalid age group")
)

type RuleEngineManager struct {
    youngInfantEngine *YoungInfantRuleEngine
    childEngine       *ChildRuleEngine
}

func NewRuleEngineManager() (*RuleEngineManager, error) {
    youngInfantEngine, err := NewYoungInfantRuleEngine()
    if err != nil {
        return nil, fmt.Errorf("failed to create young infant engine: %w", err)
    }

    childEngine, err := NewChildRuleEngine()
    if err != nil {
        return nil, fmt.Errorf("failed to create child engine: %w", err)
    }

    return &RuleEngineManager{
        youngInfantEngine: youngInfantEngine,
        childEngine:       childEngine,
    }, nil
}

func (m *RuleEngineManager) GetEngineForAgeGroup(ageGroup domain.AgeGroup) (RuleEngineInterface, error) {
    switch ageGroup {
    case domain.AgeGroupYoungInfant:
        return m.youngInfantEngine, nil
    case domain.AgeGroupChild:
        return m.childEngine, nil
    default:
        return nil, ErrAgeGroupNotSupported
    }
}

func (m *RuleEngineManager) GetEngineForTree(treeID string) (RuleEngineInterface, error) {
    if _, err := m.youngInfantEngine.GetAssessmentTree(treeID); err == nil {
        return m.youngInfantEngine, nil
    }

    if _, err := m.childEngine.GetAssessmentTree(treeID); err == nil {
        return m.childEngine, nil
    }

    return nil, ErrTreeNotFound
}

func (m *RuleEngineManager) StartAssessmentFlow(assessmentID uuid.UUID, treeID string, ageGroup domain.AgeGroup) (*domain.AssessmentFlow, error) {
    engine, err := m.GetEngineForAgeGroup(ageGroup)
    if err != nil {
        return nil, err
    }
    return engine.StartAssessmentFlow(assessmentID, treeID)
}

func (m *RuleEngineManager) GetAllTrees() map[string][]string {
    return map[string][]string{
        string(domain.AgeGroupYoungInfant): m.youngInfantEngine.GetAvailableTrees(),
        string(domain.AgeGroupChild):       m.childEngine.GetAvailableTrees(),
    }
}

type RuleEngineInterface interface {
    StartAssessmentFlow(assessmentID uuid.UUID, treeID string) (*domain.AssessmentFlow, error)
    SubmitAnswer(flow *domain.AssessmentFlow, nodeID string, answer interface{}) (*domain.AssessmentFlow, *domain.Question, error)
    ProcessBatchAssessment(assessmentID uuid.UUID, treeID string, answers map[string]interface{}) (*domain.AssessmentFlow, error)
    GetAssessmentTree(assessmentID string) (*domain.AssessmentTree, error)
    GetCurrentQuestion(flow *domain.AssessmentFlow) (*domain.Question, error)
}