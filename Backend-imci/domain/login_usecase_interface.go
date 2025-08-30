// domain/interfaces/auth_usecase.go
package domain

import "context"

type LoginUsecase interface {
    Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
    RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error)
    Logout(ctx context.Context, token string) error
}