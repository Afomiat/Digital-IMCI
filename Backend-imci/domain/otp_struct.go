// domain/otp.go
package domain

import "time"

type OTP struct {
	ID        int       `json:"id"`
	Phone     string    `json:"phone"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type VerifyOtp struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}