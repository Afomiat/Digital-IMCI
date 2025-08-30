// domain/usecase.go or domain/interfaces/usecase_interface.go
package domain

import (
	"context"
	"github.com/google/uuid" 
)

type SignupUsecase interface {
	GetMedicalProfessionalByPhone(ctx context.Context, phone string) (*MedicalProfessional, error)
	SendOtp(ctx context.Context, professional *MedicalProfessional) error // Removed SMTP parameters
	GetOtpByPhone(ctx context.Context, phone string) (*OTP, error)
	VerifyOtp(ctx context.Context, otp *VerifyOtp) (*OTP, error)
	RegisterMedicalProfessional(ctx context.Context, form *SignupForm) (uuid.UUID, error) }