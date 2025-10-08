package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
)

type AssessmentUsecase struct {
	assessmentRepo domain.AssessmentRepository
	patientRepo    domain.PatientRepository
	contextTimeout time.Duration
}

func NewAssessmentUsecase(
	assessmentRepo domain.AssessmentRepository,
	patientRepo domain.PatientRepository,
	timeout time.Duration,
) domain.AssessmentUsecase {
	return &AssessmentUsecase{
		assessmentRepo: assessmentRepo,
		patientRepo:    patientRepo,
		contextTimeout: timeout,
	}
}

func (uc *AssessmentUsecase) CreateAssessment(ctx context.Context, req *domain.CreateAssessmentRequest, medicalProfessionalID uuid.UUID) (*domain.Assessment, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// PatientID already UUID, just fetch
	patient, err := uc.patientRepo.GetByID(ctx, req.PatientID)
	if err != nil {
		return nil, domain.ErrPatientNotFound/// what if new patient is being assessed
	}

	assessmentTime := time.Now()
	ageMonths, assessmentType, err := uc.assessmentRepo.CalculateAgeInfo(ctx, req.PatientID, assessmentTime)
	if err != nil {
		return nil, err
	}

	if err := uc.validateWeight(req.WeightKg, ageMonths); err != nil {
		return nil, err
	}

	mainSymptoms := domain.JSONB{}
	for _, symptom := range req.MainSymptoms {
		mainSymptoms[symptom] = true
	}

	assessment := &domain.Assessment{
		ID:                   uuid.New(), // Generate new UUID
		MedicalProfessionalID: medicalProfessionalID,
		PatientID:            req.PatientID,
		AssessmentType:       assessmentType,
		Status:               domain.StatusDraft,
		WeightKg:             req.WeightKg,
		Temperature:          req.Temperature,
		MainSymptoms:         mainSymptoms,
		MUAC:                 req.MUAC,
		RespiratoryRate:      req.RespiratoryRate,
		AgeMonths:            ageMonths,
		GuidelineVersion:     "2014",
		StartTime:            assessmentTime,
		IsOffline:            req.IsOffline,
		Patient:              patient,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := uc.assessmentRepo.Create(ctx, assessment); err != nil {
		return nil, fmt.Errorf("failed to create assessment: %w", err)
	}

	return assessment, nil
}

func (uc *AssessmentUsecase) validateWeight(weightKg float64, ageMonths int) error {
	if ageMonths < 2 {
		if weightKg < 0.5 || weightKg > 6.0 {
			return fmt.Errorf("weight %.2f kg outside valid range for infants (0.5–6.0 kg): %w", weightKg, domain.ErrInvalidWeight)
		}
	} else {
		if weightKg < 0.5 || weightKg > 30.0 {
			return fmt.Errorf("weight %.2f kg outside valid range for children (0.5–30.0 kg): %w", weightKg, domain.ErrInvalidWeight)
		}
	}
	return nil
}

func (uc *AssessmentUsecase) GetAssessment(ctx context.Context, assessmentID uuid.UUID, medicalProfessionalID uuid.UUID) (*domain.Assessment, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.assessmentRepo.GetByID(ctx, assessmentID, medicalProfessionalID)
}

func (uc *AssessmentUsecase) GetAssessmentsByPatient(ctx context.Context, patientID uuid.UUID, medicalProfessionalID uuid.UUID) ([]*domain.Assessment, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.assessmentRepo.GetByPatientID(ctx, patientID, medicalProfessionalID)
}

func (uc *AssessmentUsecase) UpdateAssessment(ctx context.Context, assessment *domain.Assessment) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.assessmentRepo.Update(ctx, assessment)
}

func (uc *AssessmentUsecase) DeleteAssessment(ctx context.Context, assessmentID uuid.UUID, medicalProfessionalID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.assessmentRepo.Delete(ctx, assessmentID, medicalProfessionalID)
}
