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
	whatsappService         domain.WhatsAppService
	contextTimeout          time.Duration
	env                     *config.Env
}

func NewSignupUsecase(
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
	otpRepo domain.OtpRepository,
	telegramService domain.TelegramService,
	whatsappService domain.WhatsAppService,
	timeout time.Duration,
	env *config.Env,
) domain.SignupUsecase {
	return &SignupUsecase{
		medicalProfessionalRepo: medicalProfessionalRepo,
		otpRepo:                 otpRepo,
		telegramService:         telegramService,
		whatsappService:         whatsappService,

		contextTimeout: timeout,
		env:            env,
	}
}
func (su *SignupUsecase) GetMedicalProfessionalByPhone(ctx context.Context, phone string) (*domain.MedicalProfessional, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	normalized := userutil.NormalizePhone(phone)
	if normalized == "" {
		return nil, errors.New("invalid phone number")
	}

	return su.medicalProfessionalRepo.GetByPhone(ctx, normalized)
}

func (su *SignupUsecase) RegisterMedicalProfessional(ctx context.Context, form *domain.SignupForm) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	normalizedPhone := userutil.NormalizePhone(form.Phone)
	if normalizedPhone == "" {
		return uuid.Nil, errors.New("invalid phone number")
	}

	form.Phone = normalizedPhone

	hashedPass, err := userutil.HashPassword(form.Password)
	if err != nil {
		return uuid.Nil, err
	}
	professional := domain.MedicalProfessional{
		FullName:     form.FullName,
		Phone:        normalizedPhone,
		PasswordHash: hashedPass,
		Role:         form.Role,
		UseWhatsApp:  form.UseWhatsApp,
		FacilityName: form.FacilityName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = su.medicalProfessionalRepo.Create(ctx, &professional)
	if err != nil {
		return uuid.Nil, err
	}

	return professional.ID, nil
}

func (su *SignupUsecase) PrepareSignupOTP(ctx context.Context, form *domain.SignupForm) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	form.Phone = userutil.NormalizePhone(form.Phone)
	if form.Phone == "" {
		return nil, errors.New("invalid phone number")
	}

	otp, err := su.generateAndStoreOTP(ctx, form)
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (su *SignupUsecase) SendWhatsAppOTP(ctx context.Context, form *domain.SignupForm) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	if su.whatsappService == nil {
		return nil, errors.New("whatsapp service not configured")
	}

	form.Phone = userutil.NormalizePhone(form.Phone)
	if form.Phone == "" {
		return nil, errors.New("invalid phone number")
	}

	otp, err := su.generateAndStoreOTP(ctx, form)
	if err != nil {
		return nil, err
	}

	log.Printf("Sending WhatsApp OTP to %s", form.Phone)
	phoneE164 := userutil.FormatPhoneE164(form.Phone)
	if phoneE164 == "" {
		return nil, errors.New("invalid phone number for WhatsApp")
	}

	if err := su.whatsappService.SendOTP(ctx, phoneE164, otp.Code); err != nil {
		log.Printf("Failed to send WhatsApp OTP: %v", err)
		return nil, fmt.Errorf("failed to send WhatsApp OTP: %w", err)
	}

	log.Printf("WhatsApp OTP sent successfully to %s", form.Phone)
	return otp, nil
}
func (su *SignupUsecase) GetOtpByPhone(ctx context.Context, phone string) (*domain.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	normalized := userutil.NormalizePhone(phone)
	if normalized == "" {
		return nil, errors.New("invalid phone number")
	}

	return su.otpRepo.GetOtpByPhone(ctx, normalized)
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
	fmt.Printf("in signup_usecase   Stored OTP********************************8: %+v\n", storedOTP)
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

func (su *SignupUsecase) generateAndStoreOTP(ctx context.Context, form *domain.SignupForm) (*domain.OTP, error) {
	if form == nil {
		return nil, errors.New("signup form is required")
	}

	form.Phone = userutil.NormalizePhone(form.Phone)
	if form.Phone == "" {
		return nil, errors.New("invalid phone number")
	}

	otp := domain.OTP{
		FullName:     form.FullName,
		Phone:        form.Phone,
		Role:         form.Role,
		FacilityName: form.FacilityName,
		Code:         userutil.GenerateOTP(),
		Password:     form.Password,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}

	log.Printf("Generated signup OTP for phone %s", otp.Phone)

	if err := su.otpRepo.SaveOTP(ctx, &otp); err != nil {
		log.Printf("Error saving OTP: %v", err)
		return nil, err
	}

	return &otp, nil
}
