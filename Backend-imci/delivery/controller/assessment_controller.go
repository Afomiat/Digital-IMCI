// delivery/controller/assessment_controller.go
package controller

import (
	"net/http"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssessmentController struct {
	AssessmentUsecase domain.AssessmentUsecase
}

func NewAssessmentController(assessmentUsecase domain.AssessmentUsecase) *AssessmentController {
	return &AssessmentController{
		AssessmentUsecase: assessmentUsecase,
	}
}

func (ac *AssessmentController) CreateAssessment(c *gin.Context) {
	var request domain.CreateAssessmentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
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

	mpID := medicalProfessionalID.(uuid.UUID)// ??

	assessment, err := ac.AssessmentUsecase.CreateAssessment(c.Request.Context(), &request, mpID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case domain.ErrPatientNotFound:
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		case domain.ErrInvalidWeight, domain.ErrInvalidAgeForAssessment:
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to create assessment",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Assessment created successfully",
		"assessment": assessment,
	})
}

func (ac *AssessmentController) GetAssessment(c *gin.Context) {
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

	assessment, err := ac.AssessmentUsecase.GetAssessment(c.Request.Context(), assessmentID, mpID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		if err == domain.ErrAssessmentNotFound {
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to get assessment",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assessment": assessment,
	})
}

// delivery/controller/assessment_controller.go
// Add these methods to your existing controller:

func (ac *AssessmentController) ListAssessments(c *gin.Context) {
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

	// Optional patient ID filter
	patientIDStr := c.Query("patient_id")
	var patientID uuid.UUID
	var err error
	
	if patientIDStr != "" {
		patientID, err = uuid.Parse(patientIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Invalid patient ID",
				Message: "Patient ID must be a valid UUID",
				Code:    "validation_error",
			})
			return
		}
	}

	var assessments []*domain.Assessment
	if patientID != uuid.Nil {
		assessments, err = ac.AssessmentUsecase.GetAssessmentsByPatient(c.Request.Context(), patientID, mpID)
	} else {
		// You might want to add a method to get all assessments for the medical professional
		// For now, return empty or implement as needed
		assessments = []*domain.Assessment{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get assessments",
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assessments": assessments,
		"count":       len(assessments),
	})
}

func (ac *AssessmentController) UpdateAssessment(c *gin.Context) {
	assessmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid assessment ID",
			Message: "Assessment ID must be a valid UUID",
			Code:    "validation_error",
		})
		return
	}

	var request domain.UpdateAssessmentRequest
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

	// Get existing assessment
	assessment, err := ac.AssessmentUsecase.GetAssessment(c.Request.Context(), assessmentID, mpID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Assessment not found",
			Message: err.Error(),
			Code:    "not_found",
		})
		return
	}

	// Update fields
	// Add your update logic here based on request

	if err := ac.AssessmentUsecase.UpdateAssessment(c.Request.Context(), assessment); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update assessment",
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Assessment updated successfully",
		"assessment": assessment,
	})
}

func (ac *AssessmentController) DeleteAssessment(c *gin.Context) {
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

	if err := ac.AssessmentUsecase.DeleteAssessment(c.Request.Context(), assessmentID, mpID); err != nil {
		if err == domain.ErrAssessmentNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Assessment not found",
				Message: err.Error(),
				Code:    "not_found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete assessment",
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assessment deleted successfully",
	})
}