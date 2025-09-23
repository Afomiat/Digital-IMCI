package domain

import "context"

type WhatsAppService interface {
    SendOTP(ctx context.Context, phoneNumber, code string) error

}

