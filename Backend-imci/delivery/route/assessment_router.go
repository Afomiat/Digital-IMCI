// route/assessment_routes.go
package route

import (
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/repository"
	"github.com/Afomiat/Digital-IMCI/usecase"
	ruleenginecontroller "github.com/Afomiat/Digital-IMCI/ruleengine/controller"
	ruleengineengine "github.com/Afomiat/Digital-IMCI/ruleengine/engine"
	ruleengineusecase "github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewAssessmentRouter(
	env *config.Env,
	timeout time.Duration,
	db *pgxpool.Pool,
	group *gin.RouterGroup,
) {
	assessmentRepo := repository.NewAssessmentRepo(db)
	patientRepo := repository.NewPatientRepo(db)
	medicalProfessionalAnswerRepo := repository.NewMedicalProfessionalAnswerRepo(db)
	clinicalFindingsRepo := repository.NewClinicalFindingsRepo(db)
	classificationRepo := repository.NewClassificationRepo(db)
	treatmentPlanRepo := repository.NewTreatmentPlanRepo(db)
	counselingRepo := repository.NewCounselingRepo(db)

	assessmentUsecase := usecase.NewAssessmentUsecase(assessmentRepo, patientRepo, timeout)
	
	var ruleEngine *ruleengineengine.RuleEngine
	var ruleEngineErr error
	
	ruleEngine, ruleEngineErr = ruleengineengine.NewRuleEngine()
	if ruleEngineErr != nil {
		log.Printf("🚨 Rule engine initialization failed: %v", ruleEngineErr)
		log.Printf("⚠️  Assessment creation will work, but IMCI flow will be disabled")
	} else {
		log.Printf("✅ Rule engine initialized successfully")
	}

	var ruleEngineController *ruleenginecontroller.RuleEngineController
	var ruleEngineUsecase *ruleengineusecase.RuleEngineUsecase
	
	if ruleEngine != nil {
		ruleEngineUsecase = ruleengineusecase.NewRuleEngineUsecase(
			ruleEngine,
			assessmentRepo,
			medicalProfessionalAnswerRepo,
			clinicalFindingsRepo,
			classificationRepo,
			treatmentPlanRepo,
			counselingRepo,
			timeout,
		)
		ruleEngineController = ruleenginecontroller.NewRuleEngineController(ruleEngineUsecase)
	}

	assessmentController := controller.NewAssessmentController(assessmentUsecase)

	assessmentGroup := group.Group("/assessments")
	{
		assessmentGroup.POST("", assessmentController.CreateAssessment)
		assessmentGroup.GET("/:id", assessmentController.GetAssessment)
		assessmentGroup.GET("", assessmentController.ListAssessments) 
		assessmentGroup.PUT("/:id", assessmentController.UpdateAssessment) 
		assessmentGroup.DELETE("/:id", assessmentController.DeleteAssessment) 
		
		NewTreeRoutes(assessmentGroup, ruleEngineUsecase, ruleEngineController)
	}
}