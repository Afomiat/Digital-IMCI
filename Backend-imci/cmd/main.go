// main.go
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

	// Initialize repositories - they already return the interface type
	medicalProfessionalRepo := repository.NewMedicalProfessionalRepo(db)
	otpRepo := repository.NewOtpRepository(db)

	// Initialize SMS service
	var smsService service.SMSService
	if env.UseMockSMS {
		smsService = service.NewMockSMSService()
		log.Println("Using Mock SMS Service for development")
	} else {
		smsService = service.NewSMSService(
			env.TwilioAccountSID,
			env.TwilioAuthToken,
			env.TwilioFromNumber,
		)
		log.Println("Using Twilio SMS Service")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	route.Setup(env, timeout, db, r, medicalProfessionalRepo, otpRepo, smsService)

	if err := r.Run(env.LocalServerPort); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}