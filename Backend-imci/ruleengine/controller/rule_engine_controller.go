package controller

import (
	"net/http"

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

// FIXED: Remove AssessmentID from request body since it comes from URL
type StartAssessmentFlowRequest struct {
	TreeID string `json:"tree_id" binding:"required"`
}

func (rc *RuleEngineController) StartAssessmentFlow(c *gin.Context) {
	var req StartAssessmentFlowRequest

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

	response, err := rc.ruleEngineUsecase.StartAssessmentFlow(c.Request.Context(), usecase.StartFlowRequest{
		AssessmentID: assessmentID, // Get from URL parameter
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

type SubmitAnswerRequest struct {
	NodeID string      `json:"node_id" binding:"required"`
	Answer interface{} `json:"answer" binding:"required"`
}

func (rc *RuleEngineController) SubmitAnswer(c *gin.Context) {
	var req SubmitAnswerRequest

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

	response, err := rc.ruleEngineUsecase.SubmitAnswer(c.Request.Context(), usecase.SubmitAnswerRequest{
		AssessmentID: assessmentID, // Get from URL parameter
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

// ErrorResponse matches your existing controller structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}