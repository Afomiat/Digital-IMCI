package controller

import (
	"net/http"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/internal/userutil"
	"github.com/gin-gonic/gin"
)

type TelegramController struct {
	SignupUsecase   domain.SignupUsecase
	TelegramService domain.TelegramService
}

func NewTelegramController(signupUsecase domain.SignupUsecase, telegramService domain.TelegramService) *TelegramController {
	return &TelegramController{
		SignupUsecase:   signupUsecase,
		TelegramService: telegramService,
	}
}

func (tc *TelegramController) HandleSignup(ctx *gin.Context, form *domain.SignupForm) {
	if tc.TelegramService == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Telegram signup is not available"})
		return
	}

	if _, err := tc.SignupUsecase.PrepareSignupOTP(ctx.Request.Context(), form); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !tc.TelegramService.IsRunning() {
		if err := tc.TelegramService.StartPolling(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start Telegram bot"})
			return
		}
	}

	phoneForDisplay := userutil.FormatPhoneE164(form.Phone)
	if phoneForDisplay == "" {
		phoneForDisplay = form.Phone
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Open the Telegram bot to receive your OTP",
		"method":    "telegram",
		"bot_link":  tc.TelegramService.GetStartLink(),
		"phone":     phoneForDisplay,
		"full_name": form.FullName,
	})
}

func (tc *TelegramController) GetStartLink(ctx *gin.Context) {
	if tc.TelegramService == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Telegram service not configured"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"start_link": tc.TelegramService.GetStartLink(),
		"message":    "Use this link to start the Telegram bot",
	})
}
