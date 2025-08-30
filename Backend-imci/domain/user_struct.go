// domain/medical_professional.go
package domain

import (
	"time"
	"github.com/google/uuid"

)

type MedicalProfessional struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MedicalProfessionalRole string

const (
	DoctorRole    MedicalProfessionalRole = "doctor"
	NurseRole     MedicalProfessionalRole = "nurse"
	TechnicianRole MedicalProfessionalRole = "technician"
	AdminRole     MedicalProfessionalRole = "admin"
)

type SignupForm struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type LoginRequest struct {
    Phone    string `json:"phone" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`           // Change to uuid.UUID
    FullName     string    `json:"full_name"`
    Phone        string    `json:"phone"`
    Role         string    `json:"role"`
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token,omitempty"`
}