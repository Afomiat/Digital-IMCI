// repository/treatment_plan_repo.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TreatmentPlanRepo struct {
	db *pgxpool.Pool
}

func NewTreatmentPlanRepo(db *pgxpool.Pool) domain.TreatmentPlanRepository {
	return &TreatmentPlanRepo{db: db}
}

func (r *TreatmentPlanRepo) Create(ctx context.Context, plan *domain.TreatmentPlan) error {
	query := `
		INSERT INTO treatment_plans (
			id, assessment_id, classification_id, drug_name, dosage, frequency,
			duration, administration_route, is_pre_referral, instructions,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	now := time.Now()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		plan.ID,
		plan.AssessmentID,
		plan.ClassificationID,
		plan.DrugName,
		plan.Dosage,
		plan.Frequency,
		plan.Duration,
		plan.AdministrationRoute,
		plan.IsPreReferral,
		plan.Instructions,
		plan.CreatedAt,
		plan.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create treatment plan: %w", err)
	}

	return nil
}

func (r *TreatmentPlanRepo) GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) ([]*domain.TreatmentPlan, error) {
	query := `
		SELECT id, assessment_id, classification_id, drug_name, dosage, frequency,
			duration, administration_route, is_pre_referral, instructions,
			created_at, updated_at
		FROM treatment_plans 
		WHERE assessment_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.Query(ctx, query, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query treatment plans: %w", err)
	}
	defer rows.Close()

	var plans []*domain.TreatmentPlan
	for rows.Next() {
		var plan domain.TreatmentPlan
		err := rows.Scan(
			&plan.ID,
			&plan.AssessmentID,
			&plan.ClassificationID,
			&plan.DrugName,
			&plan.Dosage,
			&plan.Frequency,
			&plan.Duration,
			&plan.AdministrationRoute,
			&plan.IsPreReferral,
			&plan.Instructions,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan treatment plan: %w", err)
		}
		plans = append(plans, &plan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating treatment plans: %w", err)
	}

	return plans, nil
}