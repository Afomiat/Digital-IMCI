// route/assessment_routes.go
package route

import (
	"log"
	"time"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/controller"
	"github.com/Afomiat/Digital-IMCI/repository"
	"github.com/Afomiat/Digital-IMCI/usecase"
	younginfantcontroller "github.com/Afomiat/Digital-IMCI/ruleengine/controller"
	childcontroller "github.com/Afomiat/Digital-IMCI/ruleengine/controller"
	"github.com/Afomiat/Digital-IMCI/ruleengine/engine"
	younginfantusecase "github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
	childusecase "github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
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
	
	var youngInfantController *younginfantcontroller.YoungInfantRuleEngineController
	var youngInfantUsecase *younginfantusecase.YoungInfantRuleEngineUsecase
	var childController *childcontroller.ChildRuleEngineController
	var childUsecase *childusecase.ChildRuleEngineUsecase
		youngInfantEngine, err := engine.NewYoungInfantRuleEngine()
	if err != nil {
		log.Printf("⚠️  Young infant rule engine initialization failed: %v", err)
	} else {
		log.Printf("✅ Young infant rule engine initialized successfully")
		youngInfantUsecase = younginfantusecase.NewYoungInfantRuleEngineUsecase(
			youngInfantEngine,
			assessmentRepo,
			medicalProfessionalAnswerRepo,
			clinicalFindingsRepo,
			classificationRepo,
			treatmentPlanRepo,
			counselingRepo,
			timeout,
		)
		youngInfantController = younginfantcontroller.NewYoungInfantRuleEngineController(youngInfantUsecase)
		log.Printf("✅ Young infant rule engine use case initialized successfully")
	}

	childEngine, err := engine.NewChildRuleEngine()
	if err != nil {
		log.Printf("⚠️  Child rule engine initialization failed: %v", err)
	} else {
		log.Printf("✅ Child rule engine initialized successfully")
		childUsecase = childusecase.NewChildRuleEngineUsecase(
			childEngine,
			assessmentRepo,
			medicalProfessionalAnswerRepo,
			clinicalFindingsRepo,
			classificationRepo,
			treatmentPlanRepo,
			counselingRepo,
			timeout,
		)
		childController = childcontroller.NewChildRuleEngineController(childUsecase)
		log.Printf("✅ Child rule engine use case initialized successfully")
	}

	assessmentController := controller.NewAssessmentController(assessmentUsecase)

	assessmentGroup := group.Group("/assessments")
	{
		assessmentGroup.POST("", assessmentController.CreateAssessment)
		assessmentGroup.GET("/:id", assessmentController.GetAssessment)
		assessmentGroup.GET("", assessmentController.ListAssessments) 
		assessmentGroup.PUT("/:id", assessmentController.UpdateAssessment) 
		assessmentGroup.DELETE("/:id", assessmentController.DeleteAssessment) 
		
		NewYoungInfantTreeRoutes(assessmentGroup, youngInfantUsecase, youngInfantController)
		NewChildTreeRoutes(assessmentGroup, childUsecase, childController)
	}
}