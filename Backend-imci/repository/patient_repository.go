// repository/patient_repo.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientRepo struct {
	db *pgxpool.Pool
}

func NewPatientRepo(db *pgxpool.Pool) domain.PatientRepository {
	return &PatientRepo{db: db}
}

func (p *PatientRepo) Create(ctx context.Context, patient *domain.Patient) error {
	query := `
	INSERT INTO patients (id, name, date_of_birth, gender, is_offline, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`
	
	err := p.db.QueryRow(ctx, query,
		uuid.New(),
		patient.Name,
		patient.DateOfBirth,
		patient.Gender,
		patient.IsOffline,
		time.Now(),
		time.Now(),
	).Scan(&patient.ID)
	
	if err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}
	return nil
}

func (p *PatientRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
	patient := &domain.Patient{}
	query := `SELECT id, name, date_of_birth, gender, is_offline, created_at, updated_at 
	          FROM patients WHERE id=$1`
	err := p.db.QueryRow(ctx, query, id).Scan(
		&patient.ID,
		&patient.Name,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.IsOffline,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)
	if err != nil {
		return nil, domain.ErrPatientNotFound
	}
	return patient, nil
}

func (p *PatientRepo) GetAll(ctx context.Context, page, perPage int) ([]*domain.Patient, int, error) {
	offset := (page - 1) * perPage
	
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM patients`
	err := p.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count patients: %w", err)
	}

	query := `SELECT id, name, date_of_birth, gender, is_offline, created_at, updated_at 
	          FROM patients 
			  ORDER BY created_at DESC 
			  LIMIT $1 OFFSET $2`
	
	rows, err := p.db.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get patients: %w", err)
	}
	defer rows.Close()

	var patients []*domain.Patient
	for rows.Next() {
		var patient domain.Patient
		if err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&patient.DateOfBirth,
			&patient.Gender,
			&patient.IsOffline,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan patient: %w", err)
		}
		patients = append(patients, &patient)
	}
	
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating patients: %w", err)
	}
	
	return patients, totalCount, nil
}

func (p *PatientRepo) Update(ctx context.Context, patient *domain.Patient) error {
	query := `
	UPDATE patients 
	SET name=$1, date_of_birth=$2, gender=$3, is_offline=$4, updated_at=$5 
	WHERE id=$6
	`
	result, err := p.db.Exec(ctx, query,
		patient.Name,
		patient.DateOfBirth,
		patient.Gender,
		patient.IsOffline,
		time.Now(),
		patient.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPatientNotFound
	}
	return nil
}

func (p *PatientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM patients WHERE id=$1`
	result, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPatientNotFound
	}
	return nil
}