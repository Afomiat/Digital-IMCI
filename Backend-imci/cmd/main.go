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

	// Initialize Telegram service - REAL SERVICE ONLY
	if env.TelegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	telegramService, err := service.NewTelegramBotService(env.TelegramBotToken, telegramRepo)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	log.Printf("Telegram Bot Service initialized successfully")

	// ✅ REMOVED SMS service entirely - we're only using Telegram now

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Pass only telegramService (no SMS service)
	route.Setup(env, timeout, db, r, medicalProfessionalRepo, otpRepo, telegramService)

	if err := r.Run(env.LocalServerPort); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}