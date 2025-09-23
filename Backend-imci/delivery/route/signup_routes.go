package route

import (
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSignUpRouter(
	env *config.Env,
	timeout time.Duration,
	db *pgxpool.Pool,
	Group *gin.RouterGroup,
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
	otpRepo domain.OtpRepository,
	telegramService domain.TelegramService,
	whatsappService domain.WhatsAppService,
	telegramRepo domain.TelegramRepository, // Add this parameter
) {
	signUsecase := usecase.NewSignupUsecase(
		medicalProfessionalRepo, 
		otpRepo, 
		telegramService, 
		whatsappService,
		timeout, 
		env,
	)
	
	// Pass telegramRepo to the controller
	signController := controller.NewSignupController(signUsecase, telegramRepo, env)
	Group.POST("/signup", signController.Signup)
	Group.POST("/verify", signController.Verify)
	Group.GET("/debug-config", signController.DebugConfig)
	Group.GET("/validate-telegram", signController.ValidateTelegramSession)
}