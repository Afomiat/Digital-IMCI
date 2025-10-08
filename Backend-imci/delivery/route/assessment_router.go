package route

import (
	"log"
	"net/http"
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
	// Initialize all repositories
	assessmentRepo := repository.NewAssessmentRepo(db)
	patientRepo := repository.NewPatientRepo(db)
	medicalProfessionalAnswerRepo := repository.NewMedicalProfessionalAnswerRepo(db)
	clinicalFindingsRepo := repository.NewClinicalFindingsRepo(db)
	classificationRepo := repository.NewClassificationRepo(db)
	treatmentPlanRepo := repository.NewTreatmentPlanRepo(db)
	counselingRepo := repository.NewCounselingRepo(db)

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
		
		if ruleEngineController != nil && ruleEngineUsecase != nil {
			// New endpoint to get assessment tree
			assessmentGroup.GET("/tree/birth_asphyxia", func(c *gin.Context) {
				tree, err := ruleEngineUsecase.GetAssessmentTree("birth_asphyxia_check")
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to get assessment tree",
						"message": err.Error(),
						"code":    "internal_error",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"tree": tree,
				})
			})
			
			assessmentGroup.POST("/:id/start-flow", ruleEngineController.StartAssessmentFlow)
			assessmentGroup.POST("/:id/answer", ruleEngineController.SubmitAnswer)
		} else {
			// Fallback endpoints when rule engine is unavailable
			assessmentGroup.POST("/:id/start-flow", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error": "IMCI rule engine unavailable",
					"message": "Rule engine failed to initialize. Check server logs.",
					"code": "rule_engine_unavailable",
				})
			})
			assessmentGroup.POST("/:id/answer", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error": "IMCI rule engine unavailable",
					"message": "Rule engine failed to initialize. Check server logs.",
					"code": "rule_engine_unavailable",
				})
			})
		}
	}
}