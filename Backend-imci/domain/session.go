package domain

import (
    "sync"
    "time"
)

// TelegramSession represents a Telegram session with an expiry.
type TelegramSession struct {
    Username  string
    Token     string
    ExpiresAt time.Time
}

var (
    telegramSessions      = make(map[string]*TelegramSession)
    telegramSessionsMutex sync.RWMutex
)

// GetTelegramSession retrieves a Telegram session by username.
func GetTelegramSession(username string) (*TelegramSession, bool) {
    telegramSessionsMutex.RLock()
    defer telegramSessionsMutex.RUnlock()
    session, exists := telegramSessions[username]
    return session, exists
}