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
	whatsappService        domain.WhatsAppService
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
		whatsappService:        whatsappService,

		contextTimeout:          timeout,
		env:                     env,
	}
}
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
	fmt.Printf("Original password: %s\n", form.Password)
	fmt.Printf("Hashed password: %s\n", hashedPass)
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


// usecase/signup_usecase.go - Update the SendOtp method
func (su *SignupUsecase) SendOtp(ctx context.Context, professional *domain.MedicalProfessional) error {
	log.Printf("SendOtp called for professional: %+v", professional)
	
	storedOTP, err := su.otpRepo.GetOtpByPhone(ctx, professional.Phone)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		log.Printf("Error checking existing OTP: %v", err)
		return err
	}

	if storedOTP != nil {
		log.Printf("Found existing OTP: %+v", storedOTP)
		if time.Now().Before(storedOTP.ExpiresAt) {
			log.Printf("OTP already sent and not expired yet")
			return errors.New("OTP already sent")
		}
		if err := su.otpRepo.DeleteOTP(ctx, storedOTP.Phone); err != nil {
			log.Printf("Error deleting expired OTP: %v", err)
			return err
		}
		log.Printf("Deleted expired OTP")
	}

	otp := domain.OTP{
		FullName: professional.FullName,
		Phone:     professional.Phone,
		Code:      userutil.GenerateOTP(),
		Password:  professional.PasswordHash, // Store hashed password
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 5),
	}

	log.Printf("Generated new OTP: %s for phone: %s", otp.Code, otp.Phone)

	if err := su.otpRepo.SaveOTP(ctx, &otp); err != nil {
		log.Printf("Error saving OTP: %v", err)
		return err
	}
	log.Printf("OTP saved successfully")

	if professional.TelegramUsername != "" {
		log.Printf("Sending Telegram OTP to @%s", professional.TelegramUsername)
		if err := su.telegramService.SendOTP(ctx, professional.TelegramUsername, otp.Code); err != nil {
			log.Printf("Failed to send Telegram OTP: %v", err)
			return fmt.Errorf("failed to send Telegram OTP: %w", err)
		}
		log.Printf("Telegram OTP sent successfully to @%s", professional.TelegramUsername)
		
	} else if professional.UseWhatsApp {
		log.Printf("Sending WhatsApp OTP to %s", professional.Phone)
		if err := su.whatsappService.SendOTP(ctx, professional.Phone, otp.Code); err != nil {
			log.Printf("Failed to send WhatsApp OTP: %v", err)
			return fmt.Errorf("failed to send WhatsApp OTP: %w", err)
		}
		log.Printf("WhatsApp OTP sent successfully to %s", professional.Phone)
		
	} else {
		log.Printf("No OTP delivery method specified")
		return errors.New("no OTP delivery method specified. Use Telegram or WhatsApp")
	}

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

// usecase/signup_usecase.go
