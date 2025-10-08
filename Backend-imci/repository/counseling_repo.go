// repository/counseling_repo.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CounselingRepo struct {
	db *pgxpool.Pool
}

func NewCounselingRepo(db *pgxpool.Pool) domain.CounselingRepository {
	return &CounselingRepo{db: db}
}

func (r *CounselingRepo) Create(ctx context.Context, counseling *domain.Counseling) error {
	query := `
		INSERT INTO counselings (
			id, assessment_id, advice_type, details, language, 
			understood_by_caregiver, questions_asked, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	counseling.CreatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		counseling.ID,
		counseling.AssessmentID,
		counseling.AdviceType,
		counseling.Details,
		counseling.Language,
		counseling.UnderstoodByCaregiver,
		counseling.QuestionsAsked,
		counseling.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create counseling: %w", err)
	}

	return nil
}

func (r *CounselingRepo) GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) ([]*domain.Counseling, error) {
	query := `
		SELECT id, assessment_id, advice_type, details, language, 
			understood_by_caregiver, questions_asked, created_at
		FROM counselings 
		WHERE assessment_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.Query(ctx, query, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query counselings: %w", err)
	}
	defer rows.Close()

	var counselings []*domain.Counseling
	for rows.Next() {
		var counseling domain.Counseling
		var understoodByCaregiver, questionsAsked *string
		
		err := rows.Scan(
			&counseling.ID,
			&counseling.AssessmentID,
			&counseling.AdviceType,
			&counseling.Details,
			&counseling.Language,
			&understoodByCaregiver,
			&questionsAsked,
			&counseling.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan counseling: %w", err)
		}

		// Handle nullable fields
		if understoodByCaregiver != nil {
			counseling.UnderstoodByCaregiver = understoodByCaregiver
		}
		if questionsAsked != nil {
			counseling.QuestionsAsked = questionsAsked
		}

		counselings = append(counselings, &counseling)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating counselings: %w", err)
	}

	return counselings, nil
}