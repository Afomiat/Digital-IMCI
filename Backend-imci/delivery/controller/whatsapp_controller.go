package controller

import (
	"net/http"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/internal/userutil"
	"github.com/gin-gonic/gin"
)

type WhatsAppController struct {
	SignupUsecase domain.SignupUsecase
}

func NewWhatsAppController(signupUsecase domain.SignupUsecase) *WhatsAppController {
	return &WhatsAppController{
		SignupUsecase: signupUsecase,
	}
}

func (wc *WhatsAppController) HandleSignup(ctx *gin.Context, form *domain.SignupForm) {
	otp, err := wc.SignupUsecase.SendWhatsAppOTP(ctx.Request.Context(), form)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	phoneForDisplay := userutil.FormatPhoneE164(form.Phone)
	if phoneForDisplay == "" {
		phoneForDisplay = form.Phone
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP sent via WhatsApp",
		"method":  "whatsapp",
		"phone":   phoneForDisplay,
		"expires": otp.ExpiresAt,
	})
}
