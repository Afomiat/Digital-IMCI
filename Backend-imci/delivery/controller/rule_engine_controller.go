// delivery/controller/rule_engine_controller.go
package controller

import (
	"net/http"

	ruleengineusecase "github.com/Afomiat/Digital-IMCI/ruleengine/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RuleEngineController struct {
	ruleEngineUsecase *ruleengineusecase.RuleEngineUsecase
}

func NewRuleEngineController(ruleEngineUsecase *ruleengineusecase.RuleEngineUsecase) *RuleEngineController {
	return &RuleEngineController{
		ruleEngineUsecase: ruleEngineUsecase,
	}
}

// FIXED: Proper request types
type StartAssessmentFlowRequest struct {
	AssessmentID uuid.UUID `json:"assessment_id" binding:"required"`
}

type SubmitAnswerRequest struct {
	AssessmentID uuid.UUID   `json:"assessment_id" binding:"required"`
	NodeID       string      `json:"node_id" binding:"required"`
	Answer       interface{} `json:"answer" binding:"required"`
}

func (rc *RuleEngineController) StartAssessmentFlow(c *gin.Context) {
	var request StartAssessmentFlowRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "validation_error",
		})
		return
	}

	// Get medical professional ID from auth middleware
	medicalProfessionalID, exists := c.Get("medical_professional_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Medical professional ID not found",
			Code:    "unauthorized",
		})
		return
	}

	mpID := medicalProfessionalID.(uuid.UUID)// was it hashed ??

	result, err := rc.ruleEngineUsecase.StartAssessmentFlow(c.Request.Context(), request.AssessmentID, mpID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to start assessment flow",
			Message: err.Error(),
			Code:    "rule_engine_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assessment flow started",
		"data":    result,
	})
}

func (rc *RuleEngineController) SubmitAnswer(c *gin.Context) {
	var request SubmitAnswerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "validation_error",
		})
		return
	}

	// Get medical professional ID from auth middleware
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

	result, err := rc.ruleEngineUsecase.SubmitAnswer(c.Request.Context(), request.AssessmentID, mpID, request.NodeID, request.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to process answer",
			Message: err.Error(),
			Code:    "rule_engine_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Answer processed successfully",
		"data":    result,
	})
}