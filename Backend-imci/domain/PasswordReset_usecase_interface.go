package domain

import (
	"context"
)

type PasswordResetUsecase interface {
	InitiatePasswordReset(ctx context.Context, phone string, useWhatsApp bool) error
	VerifyPasswordResetOTP(ctx context.Context, phone, otpCode string) (bool, error)
	ResetPassword(ctx context.Context, phone, newPassword string) error
}

