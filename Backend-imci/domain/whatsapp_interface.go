package domain

import "context"

type WhatsAppService interface {
    SendOTP(ctx context.Context, phoneNumber, code string) error
    SendPasswordResetOTP(ctx context.Context, phone, otpCode string) error

}

