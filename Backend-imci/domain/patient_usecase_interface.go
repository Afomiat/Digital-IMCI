package domain

import (
	"context"

	"github.com/google/uuid"
)


type PatientUsecase interface {
	CreatePatient(ctx context.Context, patient *Patient) error
	GetPatient(ctx context.Context, id uuid.UUID) (*Patient, error)
	GetAllPatients(ctx context.Context) ([]*Patient, error)
	UpdatePatient(ctx context.Context, patient *Patient) error
	DeletePatient(ctx context.Context, id uuid.UUID) error
}