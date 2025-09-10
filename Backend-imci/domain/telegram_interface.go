package domain

import "context"

// TelegramService defines the contract for sending messages via Telegram
type TelegramService interface {
    // SendOTP sends an OTP code to a user via Telegram.
    // It requires the user's Telegram username and the OTP code.
    SendOTP(ctx context.Context, telegramUsername, code string) error
	GetStartLink() string
}

