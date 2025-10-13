package controller

import (
	"net/http"

	"github.com/Afomiat/Digital-IMCI/ruleengine/domain"
	"github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RuleEngineController struct {
	ruleEngineUsecase *usecase.RuleEngineUsecase
}

func NewRuleEngineController(ruleEngineUsecase *usecase.RuleEngineUsecase) *RuleEngineController {
	return &RuleEngineController{
		ruleEngineUsecase: ruleEngineUsecase,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (rc *RuleEngineController) StartAssessmentFlow(c *gin.Context) {
	var req struct {
		TreeID string `json:"tree_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "validation_error",
		})
		return
	}

	assessmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid assessment ID",
			Message: "Assessment ID must be a valid UUID",
			Code:    "validation_error",
		})
		return
	}

	medicalProfessionalID, exists := c.Get("medical_professional_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Medical professional ID not found",
			Code:    "unauthorized",
		})
		return
	}

	mpID := medicalProfessionalID.(uuid.UUID)

	response, err := rc.ruleEngineUsecase.StartAssessmentFlow(c.Request.Context(), domain.StartFlowRequest{
		AssessmentID: assessmentID, 
		TreeID:       req.TreeID,
	}, mpID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to start assessment flow",
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assessment flow started successfully",
		"data":    response,
	})
}

func (rc *RuleEngineController) SubmitAnswer(c *gin.Context) {
	var req struct {
		NodeID string      `json:"node_id" binding:"required"`
		Answer interface{} `json:"answer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "validation_error",
		})
		return
	}

	assessmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid assessment ID",
			Message: "Assessment ID must be a valid UUID",
			Code:    "validation_error",
		})
		return
	}

	medicalProfessionalID, exists := c.Get("medical_professional_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Medical professional ID not found",
			Code:    "unauthorized",
		})
		return
	}

	mpID := medicalProfessionalID.(uuid.UUID)

	response, err := rc.ruleEngineUsecase.SubmitAnswer(c.Request.Context(), domain.SubmitAnswerRequest{
		AssessmentID: assessmentID, 
		NodeID:       req.NodeID,
		Answer:       req.Answer,
	}, mpID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to submit answer",
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Answer submitted successfully",
		"data":    response,
	})
}

// Batch Processing Endpoints
func (ctrl *RuleEngineController) ProcessBatchAssessment(c *gin.Context) {
    var req domain.BatchProcessRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request body",
            "message": err.Error(),
        })
        return
    }

    medicalProfessionalIDInterface, exists := c.Get("medical_professional_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Medical professional ID not found in context",
        })
        return
    }

    medicalProfessionalID, ok := medicalProfessionalIDInterface.(uuid.UUID)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Medical professional ID has wrong type in context",
        })
        return
    }

    response, err := ctrl.ruleEngineUsecase.ProcessBatchAssessment(
        c.Request.Context(), 
        req, 
        medicalProfessionalID,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to process batch assessment",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, response)
}

func (ctrl *RuleEngineController) GetTreeQuestions(c *gin.Context) {
	treeID := c.Param("treeId")
	if treeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tree ID is required",
		})
		return
	}

	tree, err := ctrl.ruleEngineUsecase.GetTreeQuestions(treeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get tree questions",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tree_id":      tree.AssessmentID,
		"title":        tree.Title,
		"instructions": tree.Instructions,
		"questions":    tree.QuestionsFlow,
	})
}