// domain/token_blacklist.go
package domain

import (
	"context"
	"time"
)

type TokenBlacklistRepository interface {
    BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
    IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type BlacklistedToken struct {
    Token      string    `json:"token"`
    ExpiresAt  time.Time `json:"expires_at"`
    UserID     string    `json:"user_id"`
    BlacklistedAt time.Time `json:"blacklisted_at"`
}