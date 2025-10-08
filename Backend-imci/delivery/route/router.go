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
) {
	medicalProfessionalRepo := repository.NewMedicalProfessionalRepo(db)
	
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

	authMiddleware := middleware.NewAuthMiddleware(env, blacklistRepo).Handler()
	
	public := r.Group("/api/v1")
	protected := r.Group("/api/v1")
	protected.Use(authMiddleware)
	
	NewSignUpRouter(env, timeout, db, public, medicalProfessionalRepo)
	NewLoginRouter(env, timeout, db, public, medicalProfessionalRepo)
	NewPasswordResetRouter(env, timeout, db, public, medicalProfessionalRepo)
	NewPatientRouter(env, timeout, db, protected)
	NewLogoutRouter(env, protected, blacklistRepo)
	NewAssessmentRouter(env, timeout, db, protected)

}
