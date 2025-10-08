// repository/classification_repo.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClassificationRepo struct {
	db *pgxpool.Pool
}

func NewClassificationRepo(db *pgxpool.Pool) domain.ClassificationRepository {
	return &ClassificationRepo{db: db}
}

func (r *ClassificationRepo) Create(ctx context.Context, classification *domain.Classification) error {
	query := `
		INSERT INTO classifications (
			id, assessment_id, disease, color, details, rule_version,
			confidence_score, is_critical_illness, requires_urgent_referral,
			treatment_priority, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	classification.CreatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		classification.ID,
		classification.AssessmentID,
		classification.Disease,
		classification.Color,
		classification.Details,
		classification.RuleVersion,
		classification.ConfidenceScore,
		classification.IsCriticalIllness,
		classification.RequiresUrgentReferral,
		classification.TreatmentPriority,
		classification.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create classification: %w", err)
	}

	return nil
}

func (r *ClassificationRepo) GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) (*domain.Classification, error) {
	query := `
		SELECT id, assessment_id, disease, color, details, rule_version,
			confidence_score, is_critical_illness, requires_urgent_referral,
			treatment_priority, created_at
		FROM classifications 
		WHERE assessment_id = $1
	`

	var classification domain.Classification
	err := r.db.QueryRow(ctx, query, assessmentID).Scan(
		&classification.ID,
		&classification.AssessmentID,
		&classification.Disease,
		&classification.Color,
		&classification.Details,
		&classification.RuleVersion,
		&classification.ConfidenceScore,
		&classification.IsCriticalIllness,
		&classification.RequiresUrgentReferral,
		&classification.TreatmentPriority,
		&classification.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("classification not found for assessment: %w", err)
		}
		return nil, fmt.Errorf("failed to get classification: %w", err)
	}

	return &classification, nil
}

// repository/classification_repo.go - Fix the Upsert method
func (r *ClassificationRepo) Upsert(ctx context.Context, classification *domain.Classification) error {
    // First try to update if exists
    query := `
        UPDATE classifications 
        SET disease = $1, color = $2, details = $3, rule_version = $4,
            confidence_score = $5, is_critical_illness = $6, 
            requires_urgent_referral = $7, treatment_priority = $8,
            created_at = $9
        WHERE assessment_id = $10
    `

    result, err := r.db.Exec(ctx, query,
        classification.Disease,
        classification.Color,
        classification.Details,
        classification.RuleVersion,
        classification.ConfidenceScore,
        classification.IsCriticalIllness,
        classification.RequiresUrgentReferral,
        classification.TreatmentPriority,
        classification.CreatedAt,
        classification.AssessmentID,
    )

    if err != nil {
        return fmt.Errorf("failed to update classification: %w", err)
    }

    // If no rows were updated, then insert
    if result.RowsAffected() == 0 {
        return r.Create(ctx, classification)
    }

    return nil
}
func (r *ClassificationRepo) Update(ctx context.Context, classification *domain.Classification) error {
	query := `
		UPDATE classifications 
		SET disease = $1, color = $2, details = $3, rule_version = $4,
			confidence_score = $5, is_critical_illness = $6, 
			requires_urgent_referral = $7, treatment_priority = $8,
			created_at = $9
		WHERE assessment_id = $10
	`

	_, err := r.db.Exec(ctx, query,
		classification.Disease,
		classification.Color,
		classification.Details,
		classification.RuleVersion,
		classification.ConfidenceScore,
		classification.IsCriticalIllness,
		classification.RequiresUrgentReferral,
		classification.TreatmentPriority,
		classification.CreatedAt,
		classification.AssessmentID,
	)

	if err != nil {
		return fmt.Errorf("failed to update classification: %w", err)
	}

	return nil
}