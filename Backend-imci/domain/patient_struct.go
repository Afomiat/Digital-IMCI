package domain

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" binding:"required"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
	Gender      string    `json:"gender" binding:"required"`
	IsOffline   bool      `json:"is_offline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreatePatientRequest struct {
	Name        string `json:"name" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"` 
	Gender      string `json:"gender" binding:"required"`
	IsOffline   bool   `json:"is_offline"`
}

type UpdatePatientRequest struct {
	Name        string `json:"name" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"` 
	Gender      string `json:"gender" binding:"required"`
	IsOffline   bool   `json:"is_offline"`
}