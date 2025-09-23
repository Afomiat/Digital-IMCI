package domain

import (
	"context"

	"github.com/google/uuid"
)



type PatientRepository interface {
	Create(ctx context.Context, patient *Patient) error
	GetByID(ctx context.Context, id uuid.UUID) (*Patient, error)
	GetAll(ctx context.Context) ([]*Patient, error)
	Update(ctx context.Context, patient *Patient) error
	Delete(ctx context.Context, id uuid.UUID) error
}

