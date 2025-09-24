package domain

import (
	"time"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrPatientNotFound     = errors.New("patient not found")
	ErrInvalidDateOfBirth  = errors.New("invalid date of birth")
	ErrInvalidGender       = errors.New("invalid gender")
	ErrNameRequired        = errors.New("patient name is required")
)

type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderUnknown Gender = "unknown"
)

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderUnknown:
		return true
	default:
		return false
	}
}

type Patient struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" binding:"required"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
	Gender      Gender    `json:"gender" binding:"required"`
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

type PaginationRequest struct {
	Page    int `form:"page,default=1" binding:"min=1"`
	PerPage int `form:"per_page,default=10" binding:"min=1,max=100"`
}

type PaginatedPatientsResponse struct {
	Patients   []*Patient `json:"patients"`
	TotalCount int        `json:"total_count"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}