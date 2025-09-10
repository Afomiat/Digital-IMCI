// usecase/signup_usecase.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/internal/userutil"
	"github.com/google/uuid"
)

type SignupUsecase struct {
	medicalProfessionalRepo domain.MedicalProfessionalRepository
	otpRepo                 domain.OtpRepository
	telegramService         domain.TelegramService
	contextTimeout          time.Duration
	env                     *config.Env
}
// (medicalProfessionalRepo, otpRepo, smsService, timeout, env, telegramService)
func NewSignupUsecase(
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
	otpRepo domain.OtpRepository,
	telegramService domain.TelegramService,
	timeout time.Duration,
	env *config.Env,
) domain.SignupUsecase {
	return &SignupUsecase{
		medicalProfessionalRepo: medicalProfessionalRepo,
		otpRepo:                 otpRepo,
		telegramService:         telegramService,

		contextTimeout:          timeout,
		env:                     env,
	}
}
// Add this missing method
func (su *SignupUsecase) GetMedicalProfessionalByPhone(ctx context.Context, phone string) (*domain.MedicalProfessional, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	return su.medicalProfessionalRepo.GetByPhone(ctx, phone)
}

func (su *SignupUsecase) RegisterMedicalProfessional(ctx context.Context, form *domain.SignupForm) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	hashedPass, err := userutil.HashPassword(form.Password)
	if err != nil {
		return uuid.Nil, err
	}

	professional := domain.MedicalProfessional{
		FullName:     form.FullName,
		Phone:        form.Phone,
		PasswordHash: hashedPass,
		Role:         form.Role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = su.medicalProfessionalRepo.Create(ctx, &professional)
	if err != nil {
		return uuid.Nil, err
	}

	return professional.ID, nil
}


func (su *SignupUsecase) SendOtp(ctx context.Context, professional *domain.MedicalProfessional) error {
	storedOTP, err := su.otpRepo.GetOtpByPhone(ctx, professional.Phone)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	if storedOTP != nil {
		if time.Now().Before(storedOTP.ExpiresAt) {
			return errors.New("OTP already sent")
		}
		if err := su.otpRepo.DeleteOTP(ctx, storedOTP.Phone); err != nil {
			return err
		}
	}

	otp := domain.OTP{
		Phone:     professional.Phone,
		Code:      userutil.GenerateOTP(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 5),
	}

	if err := su.otpRepo.SaveOTP(ctx, &otp); err != nil {
		return err
	}

	// ✅ ONLY Telegram - No SMS fallback
	fmt.Print("proffesional************************", professional)
	if professional.TelegramUsername == "" {
		return errors.New("Telegram username is required for OTP delivery")
	}

	// Send OTP via Telegram only
	if err := su.telegramService.SendOTP(ctx, professional.TelegramUsername, otp.Code); err != nil {
		return fmt.Errorf("failed to send Telegram OTP: %w", err)
	}

	log.Printf("Telegram OTP sent successfully to @%s", professional.TelegramUsername)
	return nil
}
func (su *SignupUsecase) GetOtpByPhone(ctx context.Context, phone string) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()
	
	return su.otpRepo.GetOtpByPhone(ctx, phone)
}

func (su *SignupUsecase) VerifyOtp(ctx context.Context, otp *domain.VerifyOtp) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	storedOTP, err := su.GetOtpByPhone(ctx, otp.Phone)
	if err != nil || storedOTP == nil {
		if storedOTP == nil {
			return nil, errors.New("OTP not found for the provided phone. Please sign up again.")
		}
		return nil, err
	}

	if storedOTP.Code != otp.Code {
		return nil, errors.New("invalid OTP")
	}

	if time.Now().After(storedOTP.ExpiresAt) {
		return nil, errors.New("OTP has expired")
	}

	err = su.otpRepo.DeleteOTP(ctx, storedOTP.Phone)
	if err != nil {
		return nil, err
	}

	return storedOTP, nil
}

// usecase/signup_usecase.go
