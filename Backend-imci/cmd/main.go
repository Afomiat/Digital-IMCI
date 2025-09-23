package main

import (
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/route"
	"github.com/Afomiat/Digital-IMCI/repository"
	"github.com/Afomiat/Digital-IMCI/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	env := config.NewEnv()

	db := config.ConnectPostgres(env)
	timeout := time.Duration(env.ContextTimeout) * time.Second

	// Initialize repositories
	medicalProfessionalRepo := repository.NewMedicalProfessionalRepo(db)
	otpRepo := repository.NewOtpRepository(db)
	telegramRepo := repository.NewTelegramRepository(db)
	
	// Initialize password reset repository
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// Initialize Telegram service with OTP repo
	if env.TelegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	telegramService, err := service.NewTelegramBotService(
		env.TelegramBotToken, 
		telegramRepo, 
		otpRepo,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	log.Printf("Telegram Bot Service initialized successfully")

	// Initialize Meta WhatsApp service
	if env.MetaWhatsAppAccessToken == "" || env.MetaWhatsAppPhoneNumberID == "" {
		log.Fatal("Meta WhatsApp credentials are required")
	}

	whatsappService := service.NewMetaWhatsAppService(
		env.MetaWhatsAppAccessToken,
		env.MetaWhatsAppPhoneNumberID,
	)
	log.Println("Meta WhatsApp Service initialized successfully")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Pass all services to route setup including passwordResetRepo
	route.Setup(env, timeout, db, r, medicalProfessionalRepo, otpRepo, telegramService, whatsappService, telegramRepo, passwordResetRepo)

	if err := r.Run(env.LocalServerPort); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}