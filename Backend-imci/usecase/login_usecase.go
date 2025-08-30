// usecase/login_usecase.go
package usecase

import (
    "context"
    "errors"
    "time"

    "github.com/Afomiat/Digital-IMCI/config"
    "github.com/Afomiat/Digital-IMCI/domain"
    "github.com/Afomiat/Digital-IMCI/internal/userutil"
    "github.com/golang-jwt/jwt/v4"
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
        return nil, errors.New("invalid phone or password")
    }

    // Verify password
    if !userutil.CheckPasswordHash(request.Password, professional.PasswordHash) {
        return nil, errors.New("invalid phone or password")
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
        ID:          professional.ID,
        FullName:    professional.FullName,
        Phone:       professional.Phone,
        Role:        professional.Role,
        AccessToken: accessToken,
        RefreshToken: refreshToken,
    }, nil
}

// usecase/login_usecase.go
func (lu *LoginUsecase) generateAccessToken(professional *domain.MedicalProfessional) (string, error) {
	claims := jwt.MapClaims{
		"id":    professional.ID.String(), // Convert UUID to string for JWT
		"phone": professional.Phone,
		"role":  professional.Role,
		"exp":   time.Now().Add(time.Hour * time.Duration(lu.env.AccessTokenExpiryHour)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(lu.env.AccessTokenSecret))
}

func (lu *LoginUsecase) generateRefreshToken(professional *domain.MedicalProfessional) (string, error) {
	claims := jwt.MapClaims{
		"id":  professional.ID.String(), // Convert UUID to string for JWT
		"exp": time.Now().Add(time.Hour * time.Duration(lu.env.RefreshTokenExpiryHour)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(lu.env.RefreshTokenSecret))
}

func (lu *LoginUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
    // Implement token refresh logic
    return nil, errors.New("not implemented")
}

func (lu *LoginUsecase) Logout(ctx context.Context, token string) error {
    // Implement logout logic (token blacklisting)
    return errors.New("not implemented")
}