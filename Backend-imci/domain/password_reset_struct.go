package domain

import (
	"time"
)


type PasswordResetRequest struct {
	ID         string    `json:"id"`
	Phone      string    `json:"phone"`
	OTPCode    string    `json:"otp_code"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	IsVerified bool      `json:"is_verified"`
	Attempts   int       `json:"attempts"`
	UseWhatsApp bool      `json:"use_whatsapp"`
}

type ForgotPasswordRequest struct {
	Phone       string `json:"phone" binding:"required"`
	UseWhatsApp bool   `json:"use_whatsapp"`
}

type VerifyResetOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type ResetPasswordRequest struct {
	Phone           string `json:"phone" binding:"required"`
	OTP             string `json:"otp" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}