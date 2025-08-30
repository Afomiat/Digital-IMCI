package route

import (
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/service"
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
	smsService service.SMSService,
) {
	signUsecase := usecase.NewSignupUsecase(medicalProfessionalRepo, otpRepo, smsService, timeout, env)
	signController := controller.NewSignupController(signUsecase, env)

	Group.POST("/signup", signController.Signup)
	Group.POST("/verify", signController.Verify)

}