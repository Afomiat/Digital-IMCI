package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"

)

type OtpRepository struct {
	db *pgxpool.Pool
}

func NewOtpRepository(db *pgxpool.Pool) *OtpRepository {
	return &OtpRepository{db: db}
}

// Get OTP by phone (updated from email)
func (o *OtpRepository) GetOtpByPhone(ctx context.Context, phone string) (*domain.OTP, error) {
	otp := &domain.OTP{}
	query := `SELECT id, phone, code, created_at, expires_at FROM otp WHERE phone=$1`
	
	err := o.db.QueryRow(ctx, query, phone).Scan(
		&otp.ID,
		&otp.Phone,
		&otp.Code,
		&otp.CreatedAt,
		&otp.ExpiresAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			// No OTP found for this phone - this is normal for new signups
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get OTP by phone: %w", err)
	}
	
	return otp, nil
}

// Save OTP with phone (updated from email)
func (o *OtpRepository) SaveOTP(ctx context.Context, otp *domain.OTP) error {
	query := `INSERT INTO otp (phone, code, created_at, expires_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := o.db.QueryRow(ctx, query, otp.Phone, otp.Code, time.Now(), otp.ExpiresAt).Scan(&otp.ID)
	return err
}

// Delete OTP by phone (updated from email)
func (o *OtpRepository) DeleteOTP(ctx context.Context, phone string) error {
	query := `DELETE FROM otp WHERE phone=$1`
	_, err := o.db.Exec(ctx, query, phone)
	return err
}