package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
)

// var tempSignupData = struct {
// 	sync.RWMutex
// 	data map[string]domain.SignupForm
// }{
// 	data: make(map[string]domain.SignupForm),
// }

type SignupController struct {
	SignupUsecase      domain.SignupUsecase
	TelegramController *TelegramController
	WhatsAppController *WhatsAppController
}

func NewSignupController(signupUsecase domain.SignupUsecase, telegramController *TelegramController, whatsappController *WhatsAppController) *SignupController {
	return &SignupController{
		SignupUsecase:      signupUsecase,
		TelegramController: telegramController,
		WhatsAppController: whatsappController,
	}
}

func (sc *SignupController) Signup(ctx *gin.Context) {
	log.Println("1. Signup endpoint called")

	var form domain.SignupForm

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("2. Form parsed for phone:", form.Phone)

	existingProfessional, _ := sc.SignupUsecase.GetMedicalProfessionalByPhone(ctx.Request.Context(), form.Phone)
	if existingProfessional != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already registered!"})
		return
	}

	if form.UseWhatsApp {
		if sc.WhatsAppController == nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "WhatsApp signup is not available"})
			return
		}
		sc.WhatsAppController.HandleSignup(ctx, &form)
		return
	}

	if sc.TelegramController == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Telegram signup is not available"})
		return
	}

	sc.TelegramController.HandleSignup(ctx, &form)
}
func (sc *SignupController) Verify(ctx *gin.Context) {
	var otp domain.VerifyOtp

	if err := ctx.ShouldBindJSON(&otp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	OtpResponse, err := sc.SignupUsecase.VerifyOtp(ctx.Request.Context(), &otp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := domain.SignupForm{
		FullName:     OtpResponse.FullName,
		Phone:        OtpResponse.Phone,
		Password:     OtpResponse.Password,
		Role:         OtpResponse.Role,
		FacilityName: OtpResponse.FacilityName,
	}
	fmt.Printf("**********Verified OTP for phone %s, full name: %s\n", OtpResponse.Role, OtpResponse.FacilityName)
	professionalID, err := sc.SignupUsecase.RegisterMedicalProfessional(ctx.Request.Context(), &user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Registration failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Registration completed successfully",
		"user_id": professionalID,
	})
}
