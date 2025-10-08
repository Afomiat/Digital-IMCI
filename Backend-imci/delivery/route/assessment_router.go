// delivery/route/assessment_router.go
package route

import (
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/repository"
	"github.com/Afomiat/Digital-IMCI/usecase"
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
	// Initialize all repositories
	assessmentRepo := repository.NewAssessmentRepo(db)
	patientRepo := repository.NewPatientRepo(db)
	medicalProfessionalAnswerRepo := repository.NewMedicalProfessionalAnswerRepo(db)
	clinicalFindingsRepo := repository.NewClinicalFindingsRepo(db)
	classificationRepo := repository.NewClassificationRepo(db)
	treatmentPlanRepo := repository.NewTreatmentPlanRepo(db)
	counselingRepo := repository.NewCounselingRepo(db) // Add this

	// Initialize usecases
	assessmentUsecase := usecase.NewAssessmentUsecase(assessmentRepo, patientRepo, timeout)
	
	// Initialize rule engine with proper error handling
	var ruleEngine *ruleengineengine.RuleEngine
	var ruleEngineErr error
	
	ruleEngine, ruleEngineErr = ruleengineengine.NewRuleEngine()
	if ruleEngineErr != nil {
		log.Printf("🚨 Rule engine initialization failed: %v", ruleEngineErr)
		log.Printf("⚠️  Assessment creation will work, but IMCI flow will be disabled")
	} else {
		log.Printf("✅ Rule engine initialized successfully")
	}

	var ruleEngineController *controller.RuleEngineController
	if ruleEngine != nil {
		ruleEngineUsecase := ruleengineusecase.NewRuleEngineUsecase(
			ruleEngine,
			assessmentRepo,
			medicalProfessionalAnswerRepo,
			clinicalFindingsRepo,
			classificationRepo,
			treatmentPlanRepo,
			counselingRepo, // Add this
			timeout,
		)
		ruleEngineController = controller.NewRuleEngineController(ruleEngineUsecase)
	}

	assessmentController := controller.NewAssessmentController(assessmentUsecase)

	assessmentGroup := group.Group("/assessments")
	{
		// Basic assessment operations
		assessmentGroup.POST("", assessmentController.CreateAssessment)
		assessmentGroup.GET("/:id", assessmentController.GetAssessment)
		assessmentGroup.GET("", assessmentController.ListAssessments) // Add this
		assessmentGroup.PUT("/:id", assessmentController.UpdateAssessment) // Add this
		assessmentGroup.DELETE("/:id", assessmentController.DeleteAssessment) // Add this
		
		// Rule engine endpoints (only if rule engine is available)
		if ruleEngineController != nil {
			assessmentGroup.POST("/:id/start-flow", ruleEngineController.StartAssessmentFlow)
			assessmentGroup.POST("/:id/answer", ruleEngineController.SubmitAnswer)
			// assessmentGroup.GET("/:id/status", ruleEngineController.GetAssessmentStatus)
		} else {
			// Provide informative error for rule engine endpoints
			assessmentGroup.POST("/:id/start-flow", func(c *gin.Context) {
				c.JSON(503, gin.H{
					"error": "IMCI rule engine unavailable",
					"message": "Rule engine failed to initialize. Check server logs.",
					"code": "rule_engine_unavailable",
				})
			})
		}
	}
}