// domain/errors.go
package domain

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidOTP          = errors.New("invalid OTP")
	ErrInvalidToken        = errors.New("invalid token")
	ErrDatabaseOperation   = errors.New("database operation failed")
	ErrTelegramSendFailed  = errors.New("failed to send telegram message")
	ErrOTPExpired          = errors.New("OTP has expired")
	ErrOTPNotFound         = errors.New("OTP not found")
)