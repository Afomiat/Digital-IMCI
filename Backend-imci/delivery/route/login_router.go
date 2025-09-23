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

func NewLoginRouter(
	env *config.Env,
	timeout time.Duration,
	db *pgxpool.Pool,
	group *gin.RouterGroup,
	medicalProfessionalRepo domain.MedicalProfessionalRepository,
) {
	// Login doesn't need OTP, Telegram, or WhatsApp services
	loginUsecase := usecase.NewLoginUsecase(medicalProfessionalRepo, timeout, env)
	loginController := controller.NewLoginController(loginUsecase)

	group.POST("/login", loginController.Login)
	group.POST("/refresh-token", loginController.RefreshToken)
}