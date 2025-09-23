package route

import (
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/usecase"
	"github.com/Afomiat/Digital-IMCI/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPatientRouter(
	env *config.Env,
	timeout time.Duration,
	db *pgxpool.Pool,
	group *gin.RouterGroup,
) {
	patientRepo := repository.NewPatientRepo(db)
	patientUsecase := usecase.NewPatientUsecase(patientRepo, timeout)
	patientController := controller.NewPatientController(patientUsecase)

	patientGroup := group.Group("/patients")
	{
		patientGroup.POST("", patientController.CreatePatient)
		patientGroup.GET("", patientController.GetAllPatients)
		patientGroup.GET("/:id", patientController.GetPatient)
		patientGroup.PUT("/:id", patientController.UpdatePatient)
		patientGroup.DELETE("/:id", patientController.DeletePatient)
	}
}