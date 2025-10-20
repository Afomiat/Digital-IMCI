// route/tree_routes.go
package route

import (
	"log"
	"net/http"

	childcontroller "github.com/Afomiat/Digital-IMCI/ruleengine/controller"
	younginfantcontroller "github.com/Afomiat/Digital-IMCI/ruleengine/controller"
	childusecase "github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
	younginfantusecase "github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
	"github.com/gin-gonic/gin"
)

func NewYoungInfantTreeRoutes(
	assessmentGroup *gin.RouterGroup,
	youngInfantUsecase *younginfantusecase.YoungInfantRuleEngineUsecase,
	youngInfantController *younginfantcontroller.YoungInfantRuleEngineController,
) {
	if youngInfantController != nil && youngInfantUsecase != nil {
		setupYoungInfantTreeRoutes(assessmentGroup, youngInfantUsecase, youngInfantController)
	} else {
		setupYoungInfantTreeRoutesUnavailable(assessmentGroup)
	}
}

func NewChildTreeRoutes(
	assessmentGroup *gin.RouterGroup,
	childUsecase *childusecase.ChildRuleEngineUsecase,
	childController *childcontroller.ChildRuleEngineController,
) {
	if childController != nil && childUsecase != nil {
		setupChildTreeRoutes(assessmentGroup, childUsecase, childController)
	} else {
		setupChildTreeRoutesUnavailable(assessmentGroup)
	}
}

func setupYoungInfantTreeRoutes(
	assessmentGroup *gin.RouterGroup,
	youngInfantUsecase *younginfantusecase.YoungInfantRuleEngineUsecase,
	youngInfantController *younginfantcontroller.YoungInfantRuleEngineController,
) {
	youngInfantGroup := assessmentGroup.Group("/young-infant")
	{
		youngInfantGroup.GET("/trees", func(c *gin.Context) {
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
				{
					"id":          "hiv_status_assessment",
					"title":       "HIV Status Assessment and Classification",
					"description": "Assess HIV status of mother and young infant",
				},
				{
					"id":          "gestation_classification",
					"title":       "Gestation and Birth Weight Classification",
					"description": "Assess gestational age and birth weight. Classify and treat according to IMCI guidelines",
				},
				{
					"id":          "developmental_assessment",
					"title":       "Developmental Milestones Assessment",
					"description": "Assess child's developmental milestones and identify delays",
				},
			}
			c.JSON(http.StatusOK, gin.H{
				"trees":     trees,
				"age_group": "young_infant",
			})
		})

		youngInfantGroup.GET("/trees/:treeId/questions", youngInfantController.GetTreeQuestions)
		youngInfantGroup.POST("/batch-process", youngInfantController.ProcessBatchAssessment)

		youngInfantGroup.GET("/tree/diarrhea", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "diarrhea_check")
		})

		youngInfantGroup.GET("/tree/jaundice", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "jaundice_check")
		})

		youngInfantGroup.GET("/tree/birth_asphyxia", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "birth_asphyxia_check")
		})

		youngInfantGroup.GET("/tree/very_severe_disease", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "very_severe_disease_check")
		})

		youngInfantGroup.GET("/tree/feeding_problem", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "feeding_problem_underweight_check")
		})

		youngInfantGroup.GET("/tree/replacement_feeding", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "replacement_feeding_check")
		})

		youngInfantGroup.GET("/tree/hiv", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "hiv_status_assessment")
		})

		youngInfantGroup.GET("/tree/gestation", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "gestation_classification")
		})

		youngInfantGroup.GET("/tree/developmental", func(c *gin.Context) {
			getYoungInfantTreeHandler(c, youngInfantUsecase, "developmental_assessment")
		})

		youngInfantGroup.POST("/:id/start-flow", youngInfantController.StartAssessmentFlow)
		youngInfantGroup.POST("/:id/answer", youngInfantController.SubmitAnswer)
	}
}

func setupChildTreeRoutes(
	assessmentGroup *gin.RouterGroup,
	childUsecase *childusecase.ChildRuleEngineUsecase,
	childController *childcontroller.ChildRuleEngineController,
) {
	childGroup := assessmentGroup.Group("/child")
	{
		log.Printf("🔧 Child group created, setting up endpoints...")

		childGroup.GET("/trees", func(c *gin.Context) {
			trees := []map[string]string{
				{
					"id":          "child_general_danger_signs",
					"title":       "Check for General Danger Signs",
					"description": "Assess child for general danger signs that require urgent referral",
				},
				{
					"id":          "child_cough_difficult_breathing",
					"title":       "Check for Cough or Difficult Breathing",
					"description": "Assess child for cough, breathing difficulties and classify",
				},
				{
					"id":          "child_diarrhea",
					"title":       "Check for Diarrhea",
					"description": "Assess child for diarrhea and dehydration",
				},
				{
					"id":          "child_fever",
					"title":       "Check for Fever",
					"description": "Assess child for fever, malaria risk, and measles complications",
				},
				{
					"id":          "child_ear_problem",
					"title":       "Check for Ear Problems",
					"description": "Assess for ear pain, discharge, and signs of infection",
				},
				{
					"id":          "child_anemia_check",
					"title":       "Check for Anemia",
					"description": "Assess for palmar pallor and measure hemoglobin if needed",
				},
				{
					"id":          "acute_malnutrition",
					"title":       "Assess Acute Malnutrition",
					"description": "Assess acute malnutrition in children 6 months to 5 years using WFL/H Z-score, MUAC, and oedema",
				},
				{
					"id":          "feeding_assessment",
					"title":       "Feeding Assessment",
					"description": "Assess feeding practices for children under 2 years or with Anemia/MAM",
				},
				{
					"id":          "hiv_assessment",
					"title":       "HIV Infection Classification",
					"description": "Assess HIV status of mother and child, test results, and classify HIV infection",
				},
			}
			c.JSON(http.StatusOK, gin.H{
				"trees":     trees,
				"age_group": "child",
			})
		})

		childGroup.GET("/trees/:treeId/questions", childController.GetTreeQuestions)
		childGroup.POST("/batch-process", childController.ProcessBatchAssessment)
		log.Printf("🔧 /batch-process endpoint called")

		childGroup.GET("/tree/general_danger_signs", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "child_general_danger_signs")
		})

		childGroup.GET("/tree/cough_difficult_breathing", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "child_cough_difficult_breathing")
		})
		childGroup.GET("/tree/diarrhea", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, " ")
		})
		childGroup.GET("/tree/fever", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "child_fever")
		})
		childGroup.GET("/tree/ear_problem", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "child_ear_problem")
		})
		childGroup.GET("/tree/anemia", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "child_anemia_check")
		})
		childGroup.GET("/tree/acute_malnutrition", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "acute_malnutrition")
		})
		childGroup.GET("/tree/feeding_assessment", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "feeding_assessment")
		})
		childGroup.GET("/tree/hiv_assessment", func(c *gin.Context) {
			getChildTreeHandler(c, childUsecase, "hiv_assessment")
		})

		childGroup.POST("/:id/start-flow", childController.StartAssessmentFlow)
		childGroup.POST("/:id/answer", childController.SubmitAnswer)
	}
}

func setupYoungInfantTreeRoutesUnavailable(assessmentGroup *gin.RouterGroup) {
	youngInfantGroup := assessmentGroup.Group("/young-infant")
	unavailableHandler := func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Young infant rule engine unavailable",
			"message": "Young infant rule engine failed to initialize. Check server logs.",
			"code":    "young_infant_rule_engine_unavailable",
		})
	}

	youngInfantGroup.GET("/trees", unavailableHandler)
	youngInfantGroup.GET("/tree/diarrhea", unavailableHandler)
	youngInfantGroup.GET("/tree/jaundice", unavailableHandler)
	youngInfantGroup.GET("/tree/birth_asphyxia", unavailableHandler)
	youngInfantGroup.GET("/tree/very_severe_disease", unavailableHandler)
	youngInfantGroup.GET("/tree/feeding_problem", unavailableHandler)
	youngInfantGroup.GET("/tree/replacement_feeding", unavailableHandler)
	youngInfantGroup.GET("/tree/hiv", unavailableHandler)
	youngInfantGroup.GET("/tree/gestation", unavailableHandler)
	youngInfantGroup.GET("/tree/developmental", unavailableHandler)

	youngInfantGroup.POST("/:id/start-flow", unavailableHandler)
	youngInfantGroup.POST("/:id/answer", unavailableHandler)
}

func setupChildTreeRoutesUnavailable(assessmentGroup *gin.RouterGroup) {
	childGroup := assessmentGroup.Group("/child")
	unavailableHandler := func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Child rule engine unavailable",
			"message": "Child rule engine failed to initialize. Check server logs.",
			"code":    "child_rule_engine_unavailable",
		})
	}

	childGroup.GET("/trees", unavailableHandler)
	childGroup.GET("/tree/cough", unavailableHandler)
	childGroup.GET("/tree/diarrhea", unavailableHandler)
	childGroup.GET("/tree/fever", unavailableHandler)
	childGroup.GET("/tree/malnutrition", unavailableHandler)
	childGroup.GET("/tree/acute_malnutrition", unavailableHandler)
	childGroup.GET("/tree/feeding_assessment", unavailableHandler)
	childGroup.GET("/tree/ear_infection", unavailableHandler)

	childGroup.POST("/:id/start-flow", unavailableHandler)
	childGroup.POST("/:id/answer", unavailableHandler)
}

func getYoungInfantTreeHandler(c *gin.Context, youngInfantUsecase *younginfantusecase.YoungInfantRuleEngineUsecase, treeID string) {
	tree, err := youngInfantUsecase.GetAssessmentTree(treeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get assessment tree",
			"message": err.Error(),
			"code":    "internal_error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tree":      tree,
		"age_group": "young_infant",
	})
}

func getChildTreeHandler(c *gin.Context, childUsecase *childusecase.ChildRuleEngineUsecase, treeID string) {
	tree, err := childUsecase.GetAssessmentTree(treeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get assessment tree",
			"message": err.Error(),
			"code":    "internal_error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tree":      tree,
		"age_group": "child",
	})
}
