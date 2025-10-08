// usecase/patient_usecase.go
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
	if err := pu.validatePatient(patient); err != nil {
		return err
	}

	return pu.patientRepo.Create(ctx, patient)
}

func (pu *PatientUsecase) GetPatient(ctx context.Context, id uuid.UUID) (*domain.Patient, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	return pu.patientRepo.GetByID(ctx, id)
}

func (pu *PatientUsecase) GetAllPatients(ctx context.Context, page, perPage int) ([]*domain.Patient, int, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return pu.patientRepo.GetAll(ctx, page, perPage)
}

func (pu *PatientUsecase) UpdatePatient(ctx context.Context, patient *domain.Patient) error {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	if err := pu.validatePatient(patient); err != nil {
		return err
	}

	return pu.patientRepo.Update(ctx, patient)
}

func (pu *PatientUsecase) DeletePatient(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	return pu.patientRepo.Delete(ctx, id)
}

func (pu *PatientUsecase) validatePatient(patient *domain.Patient) error {
	if patient.Name == "" {
		return domain.ErrNameRequired
	}
	
	if patient.DateOfBirth.IsZero() {
		return domain.ErrInvalidDateOfBirth
	}
	
	// Validate date is not in the future
	if patient.DateOfBirth.After(time.Now()) {
		return fmt.Errorf("date of birth cannot be in the future: %w", domain.ErrInvalidDateOfBirth)
	}

	// Validate gender enum
	if !patient.Gender.IsValid() {
		return fmt.Errorf("%s: %w", patient.Gender, domain.ErrInvalidGender)
	}

	return nil
}