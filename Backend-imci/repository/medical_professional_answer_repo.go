package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalProfessionalAnswerRepo struct {
	db *pgxpool.Pool
}

func NewMedicalProfessionalAnswerRepo(db *pgxpool.Pool) domain.MedicalProfessionalAnswerRepository {
	return &MedicalProfessionalAnswerRepo{db: db}
}

func (r *MedicalProfessionalAnswerRepo) Create(ctx context.Context, answer *domain.MedicalProfessionalAnswer) error {
	query := `
		INSERT INTO medical_professional_answers (
			id, assessment_id, answers, question_set_version, clinical_findings, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	answersJSON, err := json.Marshal(answer.Answers)
	if err != nil {
		return fmt.Errorf("failed to marshal answers: %w", err)
	}

	clinicalFindingsJSON, err := json.Marshal(answer.ClinicalFindings)
	if err != nil {
		return fmt.Errorf("failed to marshal clinical findings: %w", err)
	}

	now := time.Now()
	answer.CreatedAt = now
	answer.UpdatedAt = now

	_, err = r.db.Exec(ctx, query,
		answer.ID,
		answer.AssessmentID,
		answersJSON,
		answer.QuestionSetVersion,
		clinicalFindingsJSON,
		answer.CreatedAt,
		answer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create medical professional answer: %w", err)
	}

	return nil
}

func (r *MedicalProfessionalAnswerRepo) GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) (*domain.MedicalProfessionalAnswer, error) {
	query := `
		SELECT id, assessment_id, answers, question_set_version, clinical_findings, created_at, updated_at
		FROM medical_professional_answers 
		WHERE assessment_id = $1
	`

	var answer domain.MedicalProfessionalAnswer
	var answersData []byte
	var clinicalFindingsData []byte

	err := r.db.QueryRow(ctx, query, assessmentID).Scan(
		&answer.ID,
		&answer.AssessmentID,
		&answersData,
		&answer.QuestionSetVersion,
		&clinicalFindingsData,
		&answer.CreatedAt,
		&answer.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrMedicalProfessionalAnswerNotFound
		}
		return nil, fmt.Errorf("failed to get medical professional answer: %w", err)
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(answersData, &answer.Answers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal answers: %w", err)
	}
	if err := json.Unmarshal(clinicalFindingsData, &answer.ClinicalFindings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal clinical findings: %w", err)
	}

	return &answer, nil
}

func (r *MedicalProfessionalAnswerRepo) Update(ctx context.Context, answer *domain.MedicalProfessionalAnswer) error {
	query := `
		UPDATE medical_professional_answers 
		SET answers = $1, clinical_findings = $2, updated_at = $3
		WHERE id = $4
	`

	answersJSON, err := json.Marshal(answer.Answers)
	if err != nil {
		return fmt.Errorf("failed to marshal answers: %w", err)
	}

	clinicalFindingsJSON, err := json.Marshal(answer.ClinicalFindings)
	if err != nil {
		return fmt.Errorf("failed to marshal clinical findings: %w", err)
	}

	answer.UpdatedAt = time.Now()

	_, err = r.db.Exec(ctx, query,
		answersJSON,
		clinicalFindingsJSON,
		answer.UpdatedAt,
		answer.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update medical professional answer: %w", err)
	}

	return nil
}

func (r *MedicalProfessionalAnswerRepo) Upsert(ctx context.Context, answer *domain.MedicalProfessionalAnswer) error {
	existing, err := r.GetByAssessmentID(ctx, answer.AssessmentID)
	if err != nil {
		if err == domain.ErrMedicalProfessionalAnswerNotFound {
			return r.Create(ctx, answer)
		}
		return err
	}
	
	// Update existing
	answer.ID = existing.ID
	return r.Update(ctx, answer)
}