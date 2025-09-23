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
		return nil, fmt.Errorf("failed to get patient by ID: %w", err)
	}
	return patient, nil
}

func (p *PatientRepo) GetAll(ctx context.Context) ([]*domain.Patient, error) {
	query := `SELECT id, name, date_of_birth, gender, is_offline, created_at, updated_at 
	          FROM patients ORDER BY created_at DESC`
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all patients: %w", err)
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
			return nil, fmt.Errorf("failed to scan patient: %w", err)
		}
		patients = append(patients, &patient)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating patients: %w", err)
	}
	
	return patients, nil
}

func (p *PatientRepo) Update(ctx context.Context, patient *domain.Patient) error {
	query := `
	UPDATE patients 
	SET name=$1, date_of_birth=$2, gender=$3, is_offline=$4, updated_at=$5 
	WHERE id=$6
	`
	_, err := p.db.Exec(ctx, query,
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
	return nil
}

func (p *PatientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM patients WHERE id=$1`
	_, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}
	return nil
}