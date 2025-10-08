// repository/medical_professional_repo.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalProfessionalRepo struct {
	db *pgxpool.Pool
}

func NewMedicalProfessionalRepo(db *pgxpool.Pool) domain.MedicalProfessionalRepository {
	return &MedicalProfessionalRepo{db: db}
}

func (m *MedicalProfessionalRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.MedicalProfessional, error) {
	professional := &domain.MedicalProfessional{}
	query := `
		SELECT id, full_name, phone, password_hash, role, telegram_username, use_whatsapp, facility_name, created_at, updated_at 
		FROM medical_professionals 
		WHERE id=$1
	`
	err := m.db.QueryRow(ctx, query, id).Scan(
		&professional.ID,
		&professional.FullName,
		&professional.Phone,
		&professional.PasswordHash,
		&professional.Role,
		&professional.TelegramUsername,
		&professional.UseWhatsApp,
		&professional.FacilityName,
		&professional.CreatedAt,
		&professional.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical professional by ID: %w", err)
	}
	return professional, nil
}

func (m *MedicalProfessionalRepo) GetByPhone(ctx context.Context, phone string) (*domain.MedicalProfessional, error) {
	professional := &domain.MedicalProfessional{}
	query := `
		SELECT id, full_name, phone, password_hash, role, telegram_username, use_whatsapp, facility_name, created_at, updated_at 
		FROM medical_professionals 
		WHERE phone=$1
	`
	err := m.db.QueryRow(ctx, query, phone).Scan(
		&professional.ID,
		&professional.FullName,
		&professional.Phone,
		&professional.PasswordHash,
		&professional.Role,
		&professional.TelegramUsername,
		&professional.UseWhatsApp,
		&professional.FacilityName,
		&professional.CreatedAt,
		&professional.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical professional by phone: %w", err)
	}
	return professional, nil
}

func (m *MedicalProfessionalRepo) Create(ctx context.Context, professional *domain.MedicalProfessional) error {
	query := `
		INSERT INTO medical_professionals 
		(full_name, phone, password_hash, role, telegram_username, use_whatsapp, facility_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := m.db.QueryRow(ctx, query,
		professional.FullName,
		professional.Phone,
		professional.PasswordHash,
		professional.Role,
		professional.TelegramUsername,
		professional.UseWhatsApp,
		professional.FacilityName,
		time.Now(),
		time.Now(),
	).Scan(&professional.ID)

	if err != nil {
		return fmt.Errorf("failed to create medical professional: %w", err)
	}
	return nil
}

func (m *MedicalProfessionalRepo) Update(ctx context.Context, professional *domain.MedicalProfessional) error {
	query := `
		UPDATE medical_professionals 
		SET full_name=$1, phone=$2, password_hash=$3, role=$4, telegram_username=$5, 
		    use_whatsapp=$6, facility_name=$7, updated_at=$8 
		WHERE id=$9
	`
	_, err := m.db.Exec(ctx, query,
		professional.FullName,
		professional.Phone,
		professional.PasswordHash,
		professional.Role,
		professional.TelegramUsername,
		professional.UseWhatsApp,
		professional.FacilityName,
		time.Now(),
		professional.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update medical professional: %w", err)
	}
	return nil
}

func (m *MedicalProfessionalRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM medical_professionals WHERE id=$1`
	_, err := m.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete medical professional: %w", err)
	}
	return nil
}

func (m *MedicalProfessionalRepo) GetAll(ctx context.Context) ([]*domain.MedicalProfessional, error) {
	query := `
		SELECT id, full_name, phone, password_hash, role, telegram_username, use_whatsapp, facility_name, created_at, updated_at 
		FROM medical_professionals
	`
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all medical professionals: %w", err)
	}
	defer rows.Close()

	var professionals []*domain.MedicalProfessional
	for rows.Next() {
		var professional domain.MedicalProfessional
		if err := rows.Scan(
			&professional.ID,
			&professional.FullName,
			&professional.Phone,
			&professional.PasswordHash,
			&professional.Role,
			&professional.TelegramUsername,
			&professional.UseWhatsApp,
			&professional.FacilityName,
			&professional.CreatedAt,
			&professional.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan medical professional: %w", err)
		}
		professionals = append(professionals, &professional)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating medical professionals: %w", err)
	}

	return professionals, nil
}
