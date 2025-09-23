package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
)

type PatientUsecase struct {
	patientRepo    domain.PatientRepository
	contextTimeout time.Duration
}

func NewPatientUsecase(
	patientRepo domain.PatientRepository,
	timeout time.Duration,
) domain.PatientUsecase {
	return &PatientUsecase{
		patientRepo:    patientRepo,
		contextTimeout: timeout,
	}
}

func (pu *PatientUsecase) CreatePatient(ctx context.Context, patient *domain.Patient) error {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	// Validate required fields
	if patient.Name == "" {
		return fmt.Errorf("patient name is required")
	}
	if patient.Gender == "" {
		return fmt.Errorf("patient gender is required")
	}
	if patient.DateOfBirth.IsZero() {
		return fmt.Errorf("patient date of birth is required")
	}

	// Set default values
	if patient.IsOffline {
		patient.IsOffline = false
	}

	return pu.patientRepo.Create(ctx, patient)
}

func (pu *PatientUsecase) GetPatient(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	return pu.patientRepo.GetByID(ctx, id)
}

func (pu *PatientUsecase) GetAllPatients(ctx context.Context) ([]*domain.Patient, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	return pu.patientRepo.GetAll(ctx)
}

func (pu *PatientUsecase) UpdatePatient(ctx context.Context, patient *domain.Patient) error {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	// Check if patient exists
	existingPatient, err := pu.patientRepo.GetByID(ctx, patient.ID)
	if err != nil {
		return fmt.Errorf("patient not found: %w", err)
	}

	// Validate required fields
	if patient.Name == "" {
		return fmt.Errorf("patient name is required")
	}
	if patient.Gender == "" {
		return fmt.Errorf("patient gender is required")
	}
	if patient.DateOfBirth.IsZero() {
		return fmt.Errorf("patient date of birth is required")
	}

	// Update only allowed fields
	existingPatient.Name = patient.Name
	existingPatient.DateOfBirth = patient.DateOfBirth
	existingPatient.Gender = patient.Gender
	existingPatient.IsOffline = patient.IsOffline

	return pu.patientRepo.Update(ctx, existingPatient)
}

func (pu *PatientUsecase) DeletePatient(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	// Check if patient exists
	_, err := pu.patientRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("patient not found: %w", err)
	}

	return pu.patientRepo.Delete(ctx, id)
}