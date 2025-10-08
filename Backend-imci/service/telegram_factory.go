package service

import (
	"sync"

	"github.com/Afomiat/Digital-IMCI/domain"
)

var (
	telegramInstance domain.TelegramService
	telegramOnce     sync.Once
	telegramError    error
)

func GetTelegramService(token string, telegramRepo domain.TelegramRepository, otpRepo domain.OtpRepository) (domain.TelegramService, error) {
	telegramOnce.Do(func() {
		telegramInstance, telegramError = NewTelegramBotService(token, telegramRepo, otpRepo)
	})
	return telegramInstance, telegramError
}