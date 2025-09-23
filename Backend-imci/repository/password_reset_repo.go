package repository

import (
	"context"
	"fmt"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type passwordResetRepo struct {
	db *pgxpool.Pool
}

func NewPasswordResetRepository(db *pgxpool.Pool) domain.PasswordResetRepository {
	return &passwordResetRepo{db: db}
}

func (r *passwordResetRepo) SavePasswordResetOTP(ctx context.Context, resetRequest *domain.PasswordResetRequest) error {
	query := `
		INSERT INTO password_reset_requests (id, phone, otp_code, expires_at, created_at, is_verified, attempts, use_whatsapp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (phone) 
		DO UPDATE SET 
			otp_code = EXCLUDED.otp_code,
			expires_at = EXCLUDED.expires_at,
			created_at = EXCLUDED.created_at,
			is_verified = EXCLUDED.is_verified,
			attempts = EXCLUDED.attempts,
			use_whatsapp = EXCLUDED.use_whatsapp
	`
	
	_, err := r.db.Exec(ctx, query,
		resetRequest.ID,
		resetRequest.Phone,
		resetRequest.OTPCode,
		resetRequest.ExpiresAt,
		resetRequest.CreatedAt,
		resetRequest.IsVerified,
		resetRequest.Attempts,
		resetRequest.UseWhatsApp, 
	)
	
	if err != nil {
		return fmt.Errorf("failed to save password reset OTP: %w", err)
	}
	return nil
}
func (r *passwordResetRepo) GetPasswordResetOTP(ctx context.Context, phone string) (*domain.PasswordResetRequest, error) {
	query := `
		SELECT id, phone, otp_code, expires_at, created_at, is_verified, attempts, use_whatsapp
		FROM password_reset_requests 
		WHERE phone = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	resetRequest := &domain.PasswordResetRequest{}
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&resetRequest.ID,
		&resetRequest.Phone,
		&resetRequest.OTPCode,
		&resetRequest.ExpiresAt,
		&resetRequest.CreatedAt,
		&resetRequest.IsVerified,
		&resetRequest.Attempts,
		&resetRequest.UseWhatsApp, 
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get password reset OTP: %w", err)
	}
	
	return resetRequest, nil
}
func (r *passwordResetRepo) DeletePasswordResetOTP(ctx context.Context, phone string) error {
	query := `DELETE FROM password_reset_requests WHERE phone = $1`
	_, err := r.db.Exec(ctx, query, phone)
	if err != nil {
		return fmt.Errorf("failed to delete password reset OTP: %w", err)
	}
	return nil
}

func (r *passwordResetRepo) IncrementAttempts(ctx context.Context, phone string) error {
	query := `UPDATE password_reset_requests SET attempts = attempts + 1 WHERE phone = $1`
	_, err := r.db.Exec(ctx, query, phone)
	if err != nil {
		return fmt.Errorf("failed to increment attempts: %w", err)
	}
	return nil
}