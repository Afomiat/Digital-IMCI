package domain

import "context"

type TelegramService interface {
    SendOTP(ctx context.Context, telegramUsername, code string) error
	GetStartLink() string
	SendPasswordResetOTP(ctx context.Context, username, otpCode string) error

}

type TelegramRepository interface {
	SaveChatID(ctx context.Context, username string, chatID int64, phone string) error
    GetChatIDByUsername(ctx context.Context, username string) (int64, error)
    GetUsernameByPhone(ctx context.Context, phone string) (string, error)
}