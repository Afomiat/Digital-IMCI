// domain/repository.go or domain/interfaces/repository_interface.go
package domain

import (
	"context"
	"github.com/google/uuid"
)

type MedicalProfessionalRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*MedicalProfessional, error)      // Changed to uuid.UUID
	GetByPhone(ctx context.Context, phone string) (*MedicalProfessional, error)
	Create(ctx context.Context, professional *MedicalProfessional) error
	Update(ctx context.Context, professional *MedicalProfessional) error
	Delete(ctx context.Context, id uuid.UUID) error                              // Changed to uuid.UUID
	GetAll(ctx context.Context) ([]*MedicalProfessional, error)
}

type OtpRepository interface {
	GetOtpByPhone(ctx context.Context, phone string) (*OTP, error)
	SaveOTP(ctx context.Context, otp *OTP) error
	DeleteOTP(ctx context.Context, phone string) error
}

// Remove the SignupUsecase interface from here - it should be in a separate file
// since it's a usecase interface, not a repository interface