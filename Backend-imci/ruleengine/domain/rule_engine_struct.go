// ruleengine/domain/rule_engine.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type AssessmentFlow struct {
	AssessmentID       uuid.UUID              `json:"assessment_id"`
	CurrentNode        string                 `json:"current_node"`
	Status             FlowStatus             `json:"status"`
	Answers            map[string]interface{} `json:"answers"`
	Classification     *ClassificationResult  `json:"classification,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

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
	NodeID       string            `json:"node_id"`
	Question     string            `json:"question"`
	QuestionType string            `json:"question_type"`
	Required     bool              `json:"required"`
	Level        int               `json:"level"`
	ParentNode   string            `json:"parent_node,omitempty"`
	ShowCondition string           `json:"show_condition,omitempty"`
	Answers      map[string]Answer `json:"answers"`
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
}