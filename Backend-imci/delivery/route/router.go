package route

import (
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/middleware"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/repository"
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
	telegramService domain.TelegramService,
	whatsappService domain.WhatsAppService,
	telegramRepo domain.TelegramRepository,
	passwordResetRepo domain.PasswordResetRepository,
) {
	// Initialize Redis blacklist (optional - can be nil if not using Redis)
	var blacklistRepo domain.TokenBlacklistRepository
	if env.RedisURL != "" {
		redisRepo, err := repository.NewRedisTokenBlacklist(env.RedisURL)
		if err != nil {
			log.Printf("Warning: Redis blacklist not available: %v", err)
		} else {
			blacklistRepo = redisRepo
			defer redisRepo.Close()
		}
	}

	// Create auth middleware with blacklist support
	authMiddleware := middleware.NewAuthMiddleware(env, blacklistRepo).Handler()
	
	// Public routes
	public := r.Group("/api/v1")
	NewSignUpRouter(env, timeout, db, public, medicalProfessionalRepo, otpRepo, telegramService, whatsappService, telegramRepo)
	NewLoginRouter(env, timeout, db, public, medicalProfessionalRepo)
	
	// Add password reset routes
	NewPasswordResetRouter(env, timeout, db, public, medicalProfessionalRepo, telegramRepo, telegramService, whatsappService, passwordResetRepo)
	
	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(authMiddleware)
	
	// Add patient routes under protected group
	NewPatientRouter(env, timeout, db, protected)
	
	// Logout routes
	NewLogoutRouter(env, protected, blacklistRepo)
}