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

// GetOtpByPhone retrieves the latest valid (non-expired) OTP entry for a phone number.
func (o *OtpRepository) GetOtpByPhone(ctx context.Context, phone string) (*domain.OTP, error) {
	otp := &domain.OTP{}
	// Query finds the most recent, non-expired OTP for the given phone.
	query := `
        SELECT id, phone, code, full_name, password, created_at, expires_at 
        FROM otp 
        WHERE phone = $1 AND expires_at > $2
        ORDER BY created_at DESC 
        LIMIT 1
    `
	now := time.Now()

	err := o.db.QueryRow(ctx, query, phone, now).Scan(
		&otp.ID,
		&otp.Phone,
		&otp.Code,
		&otp.FullName,
		&otp.Password,
		&otp.CreatedAt,
		&otp.ExpiresAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get OTP by phone: %w", err)
	}

	return otp, nil
}

// SaveOTP stores a new OTP record.
// Common practice: Delete any existing OTPs for this phone before saving a new one.
// repository/otp_repository.go
func (o *OtpRepository) SaveOTP(ctx context.Context, otp *domain.OTP) error {
    // Hash the password before storing in OTP
    

    tx, err := o.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    _, err = tx.Exec(ctx, `DELETE FROM otp WHERE phone = $1`, otp.Phone)
    if err != nil {
        return fmt.Errorf("failed to delete existing OTPs: %w", err)
    }

    query := `
        INSERT INTO otp (phone, code, full_name, password, expires_at) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id, created_at
    `
    err = tx.QueryRow(
        ctx,
        query,
        otp.Phone,
        otp.Code,
        otp.FullName,
        otp.Password,
        otp.ExpiresAt,
    ).Scan(&otp.ID, &otp.CreatedAt)

    if err != nil {
        return fmt.Errorf("failed to insert new OTP: %w", err)
    }

    return tx.Commit(ctx)
}

// DeleteOTP removes an OTP record based on the phone number.
// This is called after a successful verification.
func (o *OtpRepository) DeleteOTP(ctx context.Context, phone string) error {
	query := `DELETE FROM otp WHERE phone = $1`
	_, err := o.db.Exec(ctx, query, phone)
	return err
}