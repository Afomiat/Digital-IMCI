// usecase/login_usecase.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/internal/userutil"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type LoginUsecase struct {
    medicalProfessionalRepo domain.MedicalProfessionalRepository
    contextTimeout          time.Duration
    env                     *config.Env
}

func NewLoginUsecase(
    medicalProfessionalRepo domain.MedicalProfessionalRepository,
    timeout time.Duration,
    env *config.Env,
) domain.LoginUsecase {
    return &LoginUsecase{
        medicalProfessionalRepo: medicalProfessionalRepo,
        contextTimeout:          timeout,
        env:                     env,
    }
}

func (lu *LoginUsecase) Login(ctx context.Context, request *domain.LoginRequest) (*domain.LoginResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
    defer cancel()

    // Find medical professional by phone
    professional, err := lu.medicalProfessionalRepo.GetByPhone(ctx, request.Phone)
    if err != nil {
        fmt.Printf("User not found for phone: %s, error: %v\n", request.Phone, err)
        return nil, errors.New("invalid phone or password")
    }

    // Debug: Print what's stored in the database
    fmt.Printf("Stored password hash: %s\n", professional.PasswordHash)
    fmt.Printf("Input password: %s\n", request.Password)

    // Verify password - Use the same pattern as working code
    err = userutil.ComparePassword(professional.PasswordHash, request.Password)
    if err != nil {
        fmt.Printf("Password comparison failed: %v\n", err)
        return nil, errors.New("invalid credentials")
    }

    // Generate JWT tokens
    accessToken, err := lu.generateAccessToken(professional)
    if err != nil {
        return nil, errors.New("failed to generate token")
    }

    refreshToken, err := lu.generateRefreshToken(professional)
    if err != nil {
        return nil, errors.New("failed to generate refresh token")
    }

    return &domain.LoginResponse{
        ID:           professional.ID,
        FullName:     professional.FullName,
        Phone:        professional.Phone,
        Role:         professional.Role,
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}
// usecase/login_usecase.go
// usecase/login_usecase.go
func (lu *LoginUsecase) generateAccessToken(professional *domain.MedicalProfessional) (string, error) {
    claims := jwt.MapClaims{
        "id":    professional.ID.String(),
        "phone": professional.Phone,
        "role":  professional.Role,
        "type":  "access", // Add token type
        "exp":   time.Now().Add(time.Minute * time.Duration(lu.env.AccessTokenExpiryMinute)).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(lu.env.AccessTokenSecret))
}

func (lu *LoginUsecase) generateRefreshToken(professional *domain.MedicalProfessional) (string, error) {
    claims := jwt.MapClaims{
        "id":   professional.ID.String(),
        "type": "refresh", // Add token type
        "exp":  time.Now().Add(time.Hour * 24 * time.Duration(lu.env.RefreshTokenExpiryDay)).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(lu.env.RefreshTokenSecret))
}
// usecase/login_usecase.go
func (lu *LoginUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
    defer cancel()

    // 1. Verify refresh token signature and expiration
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(lu.env.RefreshTokenSecret), nil
    })

    if err != nil || !token.Valid {
        return nil, errors.New("invalid or expired refresh token")
    }

    // 2. Extract claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }

    // 3. Get user ID from claims
    userIDStr, ok := claims["id"].(string)
    if !ok {
        return nil, errors.New("user ID not found in token")
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return nil, errors.New("invalid user ID format")
    }

    // 4. Get user from database
    professional, err := lu.medicalProfessionalRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    // 5. TODO: Check if refresh token is revoked (optional enhancement)
    // if lu.isRefreshTokenRevoked(ctx, refreshToken) {
    //     return nil, errors.New("refresh token has been revoked")
    // }

    // 6. Generate new access token
    newAccessToken, err := lu.generateAccessToken(professional)
    if err != nil {
        return nil, errors.New("failed to generate new access token")
    }

    // 7. Optionally generate new refresh token (rotation)
    newRefreshToken, err := lu.generateRefreshToken(professional)
    if err != nil {
        return nil, errors.New("failed to generate new refresh token")
    }

    return &domain.LoginResponse{
        ID:           professional.ID,
        FullName:     professional.FullName,
        Phone:        professional.Phone,
        Role:         professional.Role,
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken, // Return new refresh token (rotation)
    }, nil
}
