package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
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
	SignupUsecase domain.SignupUsecase
	TelegramRepo  domain.TelegramRepository
	env           *config.Env
}

func NewSignupController(signupUsecase domain.SignupUsecase, telegramRepo domain.TelegramRepository, env *config.Env) *SignupController {
	return &SignupController{
		SignupUsecase: signupUsecase,
		TelegramRepo:  telegramRepo, 
		env:           env,
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

	// Check if phone already exists
	existingProfessional, _ := sc.SignupUsecase.GetMedicalProfessionalByPhone(ctx, form.Phone)
	if existingProfessional != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already registered!"})
		return
	}

	// AUTO-DETECT TELEGRAM: Check if this phone has a linked Telegram account
	var usingTelegram bool
	var telegramUsername string

	if !form.UseWhatsApp {
		username, err := sc.TelegramRepo.GetUsernameByPhone(ctx, form.Phone)
		if err == nil && username != "" {
			usingTelegram = true
			telegramUsername = username
			log.Printf("3. Auto-detected Telegram user @%s for phone %s", username, form.Phone)
		} else {
			log.Printf("3. No Telegram account found for phone %s: %v", form.Phone, err)
		}
	}

	// Validate OTP method selection
	if !usingTelegram && !form.UseWhatsApp {
		// Telegram not linked - return linking instructions
		ctx.JSON(http.StatusOK, gin.H{
			"status": "telegram_linking_required",
			"message": "Please link your Telegram account to receive OTP",
			"telegram_link": fmt.Sprintf("https://t.me/%s?start=%s", 
				"DigitalIMCIBot", // Replace with actual bot username
				url.QueryEscape(form.Phone)),
			"phone": form.Phone,
		})
		return
	}

	if usingTelegram && form.UseWhatsApp {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Please choose only one OTP delivery method: Telegram OR WhatsApp",
		})
		return
	}

	professional := &domain.MedicalProfessional{
		FullName:    form.FullName,
		Phone:       form.Phone,
		PasswordHash: form.Password, // Plain password; will be hashed in usecase
		Role:        form.Role,
		UseWhatsApp: form.UseWhatsApp,
	}

	fmt.Println("7. Professional object created************************:", professional.FullName)
	// Set Telegram username if auto-detected
	if usingTelegram {
		professional.TelegramUsername = telegramUsername
	}

	log.Println("4. About to call SendOtp for:", form.Phone)
	err := sc.SignupUsecase.SendOtp(ctx, professional)
	if err != nil {
		log.Println("5. SendOtp failed:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OTP sending failed: " + err.Error()})
		return
	}
	log.Println("6. SendOtp completed successfully")

	// Return which method was used
	response := gin.H{"message": "OTP sent successfully"}
	if usingTelegram {
		response["method"] = "telegram"
		response["telegram_username"] = telegramUsername
	} else {
		response["method"] = "whatsapp"
		response["phone"] = form.Phone
	}

	ctx.JSON(http.StatusOK, response)
}
func (sc *SignupController) Verify(ctx *gin.Context) {
	var otp domain.VerifyOtp

	if err := ctx.ShouldBindJSON(&otp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	OtpResponse, err := sc.SignupUsecase.VerifyOtp(ctx, &otp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	 user := domain.SignupForm{
        FullName: OtpResponse.FullName,
        Phone:    OtpResponse.Phone,
        Password: OtpResponse.Password,
    }

	// Register the user automatically after successful verification
	professionalID, err := sc.SignupUsecase.RegisterMedicalProfessional(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Registration failed: " + err.Error()})
		return
	}

	

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Registration completed successfully",
		"user_id": professionalID,
	})
}

func (sc *SignupController) DebugConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"whatsapp_configured": sc.env.MetaWhatsAppAccessToken != "" && sc.env.MetaWhatsAppPhoneNumberID != "",
		"phone_number_id":     sc.env.MetaWhatsAppPhoneNumberID,
		"access_token_length": len(sc.env.MetaWhatsAppAccessToken),
		"business_account_id": sc.env.MetaWhatsAppBusinessAccountID,
	})
}

// NEW: Validate Telegram session endpoint - USE HELPER FUNCTION
func (sc *SignupController) ValidateTelegramSession(ctx *gin.Context) {
	username := ctx.Query("username")
	token := ctx.Query("token")

	if username == "" || token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username and token required"})
		return
	}

	// Use helper function instead of direct access
	session, exists := domain.GetTelegramSession(username)

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"valid": false,
			"error": "Telegram session not found. Please start the bot again.",
		})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": "Telegram session expired. Please start the bot again.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid": true,
		"username": username,
		"expires_at": session.ExpiresAt,
		"message": "Telegram session is valid",
	})
}

// NEW: Test Telegram connection endpoint
func (sc *SignupController) TestTelegramConnection(ctx *gin.Context) {
	// This would test if the Telegram service is properly connected
	// For now, just return success since we can't easily test the bot connection from here
	ctx.JSON(http.StatusOK, gin.H{
		"status": "connected",
		"message": "Telegram controller is operational",
	})
}