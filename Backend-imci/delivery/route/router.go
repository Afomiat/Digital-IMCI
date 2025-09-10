package route

import (
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(
	env *config.Env,
	timeout time.Duration,
	db *pgxpool.Pool,
	r *gin.Engine,
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
	otpRepo domain.OtpRepository,
	telegramService domain.TelegramService, // Only Telegram, no SMS
) {
	PublicRout := r.Group("")
	
	// Create Telegram controller
	telegramController := controller.NewTelegramController(telegramService)
	
	// Setup signup routes
	NewSignUpRouter(env, timeout, db, PublicRout, medicalProfessionalRepo, otpRepo, telegramService)
	
	// Setup login routes  
	NewLoginRouter(env, timeout, db, PublicRout, medicalProfessionalRepo)
	
	// ✅ Add Telegram utility routes
	PublicRout.GET("/telegram/start-link", telegramController.GetStartLink)
}