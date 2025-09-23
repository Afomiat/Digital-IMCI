package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
)

var (
	TestEnv     *config.Env
	TestTimeout time.Duration
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// Load test environment variables
	TestEnv = &config.Env{
		LocalServerPort:         ":8081",
		ContextTimeout:          30,
		AccessTokenSecret:       "test-access-secret",
		RefreshTokenSecret:      "test-refresh-secret",
		AccessTokenExpiryMinute: 15,
		RefreshTokenExpiryDay:   7,
		TelegramBotToken:        "test-bot-token",
		MetaWhatsAppAccessToken: "test-whatsapp-token",
		MetaWhatsAppPhoneNumberID: "test-phone-number-id",
	}

	TestTimeout = time.Duration(TestEnv.ContextTimeout) * time.Second
}

func teardown() {
	// Cleanup if needed
}

func GetTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), TestTimeout)
}