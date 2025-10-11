// route/tree_routes.go
package route

import (
	"net/http"

	ruleenginecontroller "github.com/Afomiat/Digital-IMCI/ruleengine/controller"
	ruleengineusecase "github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
	"github.com/gin-gonic/gin"
)

func NewTreeRoutes(
	assessmentGroup *gin.RouterGroup,
	ruleEngineUsecase *ruleengineusecase.RuleEngineUsecase,
	ruleEngineController *ruleenginecontroller.RuleEngineController,
) {
	if ruleEngineController != nil && ruleEngineUsecase != nil {
		setupTreeRoutes(assessmentGroup, ruleEngineUsecase, ruleEngineController)
	} else {
		setupTreeRoutesUnavailable(assessmentGroup)
	}
}

func setupTreeRoutes(
	assessmentGroup *gin.RouterGroup,
	ruleEngineUsecase *ruleengineusecase.RuleEngineUsecase,
	ruleEngineController *ruleenginecontroller.RuleEngineController,
) {
	assessmentGroup.GET("/trees", func(c *gin.Context) {
		trees := []map[string]string{
			{
				"id":          "birth_asphyxia_check",
				"title":       "Check for Birth Asphyxia",
				"description": "Assess newborn for birth asphyxia and provide immediate resuscitation if needed",
			},
			{
				"id":          "very_severe_disease_check", 
				"title":       "Check for Very Severe Disease",
				"description": "Assess young infants (0-2 months) for very severe disease and local bacterial infection",
			},
			{
				"id":          "jaundice_check",
				"title":       "Check for Jaundice in Young Infant", 
				"description": "Assess young infants (0-2 months) for jaundice and classify severity",
			},
			{
				"id":          "diarrhea_check",
				"title":       "Check for Diarrhea and Dehydration",
				"description": "Assess young infants for diarrhea and classify dehydration severity",
			},
			{
				"id":          "feeding_problem_underweight_check", 
				"title":       "Assess Feeding Problems and Underweight",
				"description": "Assess infant feeding practices and nutritional status",
			},
			{
				"id":          "replacement_feeding_check",
				"title":       "Assess Replacement Feeding for HIV-Positive Mothers", 
				"description": "Assess feeding practices for infants not receiving breast milk",
			},
		}
		c.JSON(http.StatusOK, gin.H{
			"trees": trees,
		})
	})

	assessmentGroup.GET("/tree/diarrhea", func(c *gin.Context) {
		getTreeHandler(c, ruleEngineUsecase, "diarrhea_check")
	})

	assessmentGroup.GET("/tree/jaundice", func(c *gin.Context) {
		getTreeHandler(c, ruleEngineUsecase, "jaundice_check")
	})

	assessmentGroup.GET("/tree/birth_asphyxia", func(c *gin.Context) {
		getTreeHandler(c, ruleEngineUsecase, "birth_asphyxia_check")
	})

	assessmentGroup.GET("/tree/very_severe_disease", func(c *gin.Context) {
		getTreeHandler(c, ruleEngineUsecase, "very_severe_disease_check")
	})

	assessmentGroup.GET("/tree/feeding_problem", func(c *gin.Context) { 
		getTreeHandler(c, ruleEngineUsecase, "feeding_problem_underweight_check")
	})
	assessmentGroup.GET("/tree/replacement_feeding", func(c *gin.Context) {
		getTreeHandler(c, ruleEngineUsecase, "replacement_feeding_check")
	})

	assessmentGroup.POST("/:id/start-flow", ruleEngineController.StartAssessmentFlow)
	assessmentGroup.POST("/:id/answer", ruleEngineController.SubmitAnswer)
}

func setupTreeRoutesUnavailable(assessmentGroup *gin.RouterGroup) {
	unavailableHandler := func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "IMCI rule engine unavailable",
			"message": "Rule engine failed to initialize. Check server logs.",
			"code": "rule_engine_unavailable",
		})
	}

	assessmentGroup.GET("/trees", unavailableHandler)
	assessmentGroup.GET("/tree/diarrhea", unavailableHandler)
	assessmentGroup.GET("/tree/jaundice", unavailableHandler)
	assessmentGroup.GET("/tree/birth_asphyxia", unavailableHandler)
	assessmentGroup.GET("/tree/very_severe_disease", unavailableHandler)
	assessmentGroup.GET("/tree/feeding_problem", unavailableHandler) 
	assessmentGroup.GET("/tree/replacement_feeding", unavailableHandler)
	assessmentGroup.POST("/:id/start-flow", unavailableHandler)
	assessmentGroup.POST("/:id/answer", unavailableHandler)
}

func getTreeHandler(c *gin.Context, ruleEngineUsecase *ruleengineusecase.RuleEngineUsecase, treeID string) {
	tree, err := ruleEngineUsecase.GetAssessmentTree(treeID)
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
}