package route

import (
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/repository"
	"github.com/Afomiat/Digital-IMCI/service"
	"github.com/Afomiat/Digital-IMCI/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPasswordResetRouter(
	env *config.Env,
	timeout time.Duration,
	db *pgxpool.Pool,
	group *gin.RouterGroup,
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
) {
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	otpRepo := repository.NewOtpRepository(db)
	telegramRepo := repository.NewTelegramRepository(db)
	
	var telegramService domain.TelegramService
	if env.TelegramBotToken != "" {
		telegramSvc, err := service.NewTelegramBotService(
			env.TelegramBotToken, 
			telegramRepo, 
			otpRepo,
		)
		if err != nil {
			log.Printf("Warning: Telegram service not available for password reset: %v", err)
		} else {
			telegramService = telegramSvc
		}
	}

	var whatsappService domain.WhatsAppService
	if env.MetaWhatsAppAccessToken != "" && env.MetaWhatsAppPhoneNumberID != "" {
		whatsappService = service.NewMetaWhatsAppService(
			env.MetaWhatsAppAccessToken,
			env.MetaWhatsAppPhoneNumberID,
		)
	}

	passwordResetUsecase := usecase.NewPasswordResetUsecase(
		medicalProfessionalRepo,
		passwordResetRepo,
		telegramRepo,
		telegramService,
		whatsappService,
		timeout,
		env,
	)
	
	passwordResetController := controller.NewPasswordResetController(passwordResetUsecase)
	
	group.POST("/forgot-password", passwordResetController.ForgotPassword)
	group.POST("/verify-reset-otp", passwordResetController.VerifyResetOTP)
	group.POST("/reset-password", passwordResetController.ResetPassword)
}