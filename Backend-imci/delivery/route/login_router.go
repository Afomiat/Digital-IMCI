// delivery/route/login_router.go
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
    Group *gin.RouterGroup,
    medicalProfessionalRepo domain.MedicalProfessionalRepository,
) {
    loginUsecase := usecase.NewLoginUsecase(medicalProfessionalRepo, timeout, env)
    loginController := controller.NewLoginController(loginUsecase)

    Group.POST("/login", loginController.Login)
    Group.POST("/refresh-token", loginController.RefreshToken)
    Group.POST("/logout", loginController.Logout)
}