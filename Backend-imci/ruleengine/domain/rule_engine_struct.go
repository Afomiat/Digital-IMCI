// ruleengine/domain/rule_engine.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type AssessmentFlow struct {
	AssessmentID       uuid.UUID              `json:"assessment_id"`
	TreeID             string                 `json:"tree_id"`
	CurrentNode        string                 `json:"current_node"`
	Status             FlowStatus             `json:"status"`
	Answers            map[string]interface{} `json:"answers"`
	Classification     *ClassificationResult  `json:"classification,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}


type AgeGroup string

const (
    AgeGroupYoungInfant AgeGroup = "young_infant"
    AgeGroupChild       AgeGroup = "child"
)
type FlowStatus string

const (
	FlowStatusInProgress FlowStatus = "in_progress"
	FlowStatusCompleted  FlowStatus = "completed"
	FlowStatusEmergency  FlowStatus = "emergency"
)

type ClassificationResult struct {
	Classification string   `json:"classification"`
	Color          string   `json:"color"`
	Emergency      bool     `json:"emergency"`
	Actions        []string `json:"actions"`
	TreatmentPlan  string   `json:"treatment_plan"`
	FollowUp       []string `json:"follow_up"`
	MotherAdvice   string   `json:"mother_advice"`
}

type Question struct {
	NodeID         string            `json:"node_id"`
	Question       string            `json:"question"`
	QuestionType   string            `json:"question_type"` // yes_no, single_choice, number_input, classification
	Required       bool              `json:"required"`
	Level          int               `json:"level"`
	ParentNode     string            `json:"parent_node,omitempty"`
	ShowCondition  string            `json:"show_condition,omitempty"`
	Instructions   string            `json:"instructions,omitempty"`    // ADD THIS
	Validation     *Validation       `json:"validation,omitempty"`      // ADD THIS
	Answers        map[string]Answer `json:"answers"`
}

type Answer struct {
	NextNode      string `json:"next_node,omitempty"`
	Color         string `json:"color,omitempty"`
	Action        string `json:"action,omitempty"`
	Classification string `json:"classification,omitempty"`
	EmergencyPath bool   `json:"emergency_path,omitempty"`
}

type AssessmentTree struct {
	AssessmentID   string              `json:"assessment_id"`
	Title          string              `json:"title"`
	Instructions   string              `json:"instructions"`
	QuestionsFlow  []Question          `json:"questions_flow"`
	Outcomes       map[string]Outcome  `json:"outcomes"`
	StartNode      string              `json:"start_node"`
}

type Outcome struct {
	Classification string   `json:"classification"`
	Color          string   `json:"color"`
	Emergency      bool     `json:"emergency"`
	Actions        []string `json:"actions"`
	TreatmentPlan  string   `json:"treatment_plan"`
	FollowUp       []string `json:"follow_up"`
	MotherAdvice   string   `json:"mother_advice"`
	Notes          string   `json:"notes,omitempty"`
}

type Validation struct {
	Min  float64 `json:"min,omitempty"`
	Max  float64 `json:"max,omitempty"`
	Step float64 `json:"step,omitempty"`
}

// Batch Processing
type BatchProcessRequest struct {
	AssessmentID uuid.UUID            `json:"assessment_id" binding:"required"`
	TreeID       string               `json:"tree_id" binding:"required"`
	Answers      map[string]interface{} `json:"answers" binding:"required"`
}

type BatchProcessResponse struct {
	AssessmentID  uuid.UUID            `json:"assessment_id"`
	Classification *ClassificationResult `json:"classification,omitempty"`
	Status         FlowStatus          `json:"status"`
}

// Sequential Processing
type StartFlowRequest struct {
	AssessmentID uuid.UUID
	TreeID       string `json:"tree_id" binding:"required"`
}

type StartFlowResponse struct {
	SessionID   uuid.UUID `json:"session_id"`
	Question    *Question `json:"question,omitempty"`
	IsComplete  bool      `json:"is_complete"`
	CurrentNode string    `json:"current_node"`
}

type SubmitAnswerRequest struct {
	AssessmentID uuid.UUID
	NodeID       string      `json:"node_id" binding:"required"`
	Answer       interface{} `json:"answer" binding:"required"`
}

type SubmitAnswerResponse struct {
	SessionID      uuid.UUID            `json:"session_id"`
	Question       *Question            `json:"question,omitempty"`
	Classification *ClassificationResult `json:"classification,omitempty"`
	IsComplete     bool                 `json:"is_complete"`
	CurrentNode    string               `json:"current_node"`
	Status         FlowStatus           `json:"status"`
}