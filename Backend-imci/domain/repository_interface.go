// domain/repository.go or domain/interfaces/repository_interface.go
package domain

import (
	"context"
	"github.com/google/uuid"
)

type MedicalProfessionalRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*MedicalProfessional, error)      
	GetByPhone(ctx context.Context, phone string) (*MedicalProfessional, error)
	Create(ctx context.Context, professional *MedicalProfessional) error
	Update(ctx context.Context, professional *MedicalProfessional) error
	Delete(ctx context.Context, id uuid.UUID) error                              
	GetAll(ctx context.Context) ([]*MedicalProfessional, error)
}

type OtpRepository interface {
	GetOtpByPhone(ctx context.Context, phone string) (*OTP, error)
	SaveOTP(ctx context.Context, otp *OTP) error
	DeleteOTP(ctx context.Context, phone string) error
}


type PasswordResetRepository interface {
	SavePasswordResetOTP(ctx context.Context, resetRequest *PasswordResetRequest) error
	GetPasswordResetOTP(ctx context.Context, phone string) (*PasswordResetRequest, error)
	DeletePasswordResetOTP(ctx context.Context, phone string) error
	IncrementAttempts(ctx context.Context, phone string) error
}

