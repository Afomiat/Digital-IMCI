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

func NewSignUpRouter(
	env *config.Env,
	timeout time.Duration,
	db *pgxpool.Pool,
	group *gin.RouterGroup,
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
) {
	otpRepo := repository.NewOtpRepository(db)
	telegramRepo := repository.NewTelegramRepository(db)
	
	var telegramService domain.TelegramService
	if env.TelegramBotToken != "" {
		telegramSvc, err := service.GetTelegramService( 
			env.TelegramBotToken, 
			telegramRepo, 
			otpRepo,
		)
		if err != nil {
			log.Printf("Warning: Telegram service not available: %v", err)
		} else {
			telegramService = telegramSvc
			log.Printf("Telegram Bot Service created successfully (polling not started)")
		}
	}

	var whatsappService domain.WhatsAppService
	if env.MetaWhatsAppAccessToken != "" && env.MetaWhatsAppPhoneNumberID != "" {
		whatsappService = service.NewMetaWhatsAppService(
			env.MetaWhatsAppAccessToken,
			env.MetaWhatsAppPhoneNumberID,
		)
		log.Println("Meta WhatsApp Service initialized successfully for signup")
	}

	signUsecase := usecase.NewSignupUsecase(
		medicalProfessionalRepo, 
		otpRepo, 
		telegramService, 
		whatsappService,
		timeout, 
		env,
	)
	
	signController := controller.NewSignupController(signUsecase, telegramRepo, env, telegramService)
	
	group.POST("/signup", signController.Signup)
	group.POST("/verify", signController.Verify)
	group.GET("/debug-config", signController.DebugConfig)
	group.GET("/validate-telegram", signController.ValidateTelegramSession)
	
	group.POST("/start-telegram-bot", signController.StartTelegramBot)
	group.POST("/stop-telegram-bot", signController.StopTelegramBot)
	
	telegramController := controller.NewTelegramController(telegramService)
	group.GET("/telegram/start-link", telegramController.GetStartLink)
	group.GET("/telegram/signup-qr", telegramController.GenerateSignupQR)
}