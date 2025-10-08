// domain/counseling.go
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Counseling struct {
	ID                    uuid.UUID  `json:"id"`
	AssessmentID          uuid.UUID  `json:"assessment_id"`
	AdviceType            string     `json:"advice_type"`
	Details               string     `json:"details"`
	Language              string     `json:"language"`
	UnderstoodByCaregiver *string    `json:"understood_by_caregiver,omitempty"`
	QuestionsAsked        *string    `json:"questions_asked,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
}

type CounselingRepository interface {
	Create(ctx context.Context, counseling *Counseling) error
	GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) ([]*Counseling, error)
}