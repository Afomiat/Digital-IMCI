// usecase/password_reset_usecase.go
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

type passwordResetUsecase struct {
	medicalProfessionalRepo domain.MedicalProfessionalRepository
	passwordResetRepo       domain.PasswordResetRepository
	telegramRepo            domain.TelegramRepository
	telegramService         domain.TelegramService
	whatsappService         domain.WhatsAppService
	contextTimeout          time.Duration
	env                     *config.Env
}

func NewPasswordResetUsecase(
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
	passwordResetRepo domain.PasswordResetRepository,
	telegramRepo domain.TelegramRepository,
	telegramService domain.TelegramService,
	whatsappService domain.WhatsAppService,
	timeout time.Duration,
	env *config.Env,
) domain.PasswordResetUsecase {
	return &passwordResetUsecase{
		medicalProfessionalRepo: medicalProfessionalRepo,
		passwordResetRepo:       passwordResetRepo,
		telegramRepo:           telegramRepo,
		telegramService:         telegramService,
		whatsappService:         whatsappService,
		contextTimeout:          timeout,
		env:                     env,
	}
}

func (u *passwordResetUsecase) InitiatePasswordReset(ctx context.Context, phone string, useWhatsApp bool) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	_, err := u.medicalProfessionalRepo.GetByPhone(ctx, phone)
	if err != nil {
		log.Printf("Password reset requested for non-existent phone: %s", phone)
		return nil 
	}

	existingRequest, err := u.passwordResetRepo.GetPasswordResetOTP(ctx, phone)
	if err == nil && existingRequest != nil {
		if time.Now().Before(existingRequest.ExpiresAt) {
			if time.Since(existingRequest.CreatedAt) < time.Minute {
				return errors.New("please wait before requesting another reset")
			}
		} else {
			u.passwordResetRepo.DeletePasswordResetOTP(ctx, phone)
		}
	}

	otpCode := userutil.GenerateOTP()
	resetRequest := &domain.PasswordResetRequest{
		ID:          uuid.New().String(),
		Phone:       phone,
		OTPCode:     otpCode,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
		CreatedAt:   time.Now(),
		IsVerified:  false,
		Attempts:    0,
		UseWhatsApp: useWhatsApp,
	}

	if err := u.passwordResetRepo.SavePasswordResetOTP(ctx, resetRequest); err != nil {
		return fmt.Errorf("failed to save reset request: %w", err)
	}

	if useWhatsApp {
		log.Printf("Sending password reset OTP via WhatsApp to %s", phone)
			
		if err := u.whatsappService.SendOTP(ctx, phone, otpCode); err != nil {
			return fmt.Errorf("failed to send WhatsApp OTP: %w", err)
		}
	} else {
		telegramUsername, err := u.telegramRepo.GetUsernameByPhone(ctx, phone)
		if err != nil || telegramUsername == "" {
			return fmt.Errorf("Telegram account not linked. Please link your Telegram account first or use WhatsApp")
		}
		
		log.Printf("Sending password reset OTP via Telegram to @%s", telegramUsername)
		if err := u.telegramService.SendPasswordResetOTP(ctx, telegramUsername, otpCode); err != nil {
			return fmt.Errorf("failed to send Telegram OTP: %w", err)
		}
	}
	return nil
}
func (u *passwordResetUsecase) VerifyPasswordResetOTP(ctx context.Context, phone, otpCode string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	resetRequest, err := u.passwordResetRepo.GetPasswordResetOTP(ctx, phone)
	if err != nil {
		return false, errors.New("reset request not found or expired")
	}

	if resetRequest.Attempts >= 3 {
		return false, errors.New("too many failed attempts. Please request a new OTP")
	}

	if resetRequest.OTPCode != otpCode {
		u.passwordResetRepo.IncrementAttempts(ctx, phone)
		return false, errors.New("invalid OTP code")
	}

	if time.Now().After(resetRequest.ExpiresAt) {
		return false, errors.New("OTP has expired")
	}

	resetRequest.IsVerified = true
	if err := u.passwordResetRepo.SavePasswordResetOTP(ctx, resetRequest); err != nil {
		return false, fmt.Errorf("failed to update reset request: %w", err)
	}

	return true, nil
}

func (u *passwordResetUsecase) ResetPassword(ctx context.Context, phone, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	resetRequest, err := u.passwordResetRepo.GetPasswordResetOTP(ctx, phone)
	if err != nil || !resetRequest.IsVerified {
		return errors.New("OTP not verified. Please verify OTP first")
	}

	// Hash new password
	hashedPassword, err := userutil.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	professional, err := u.medicalProfessionalRepo.GetByPhone(ctx, phone)
	if err != nil {
		return errors.New("user not found")
	}

	professional.PasswordHash = hashedPassword
	if err := u.medicalProfessionalRepo.Update(ctx, professional); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Clean up reset request
	u.passwordResetRepo.DeletePasswordResetOTP(ctx, phone)

	return nil
}