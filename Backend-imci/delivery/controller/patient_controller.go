// controller/patient_controller.go
package controller

import (
	"net/http"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PatientController struct {
	PatientUsecase domain.PatientUsecase
}

func NewPatientController(patientUsecase domain.PatientUsecase) *PatientController {
	return &PatientController{
		PatientUsecase: patientUsecase,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (pc *PatientController) CreatePatient(c *gin.Context) {
	var request domain.CreatePatientRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "validation_error",
		})
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid date format",
			Message: "Use YYYY-MM-DD format",
			Code:    "validation_error",
		})
		return
	}

	patient := &domain.Patient{
		Name:        request.Name,
		DateOfBirth: dateOfBirth,
		Gender:      domain.Gender(request.Gender),
		IsOffline:   request.IsOffline,
	}

	err = pc.PatientUsecase.CreatePatient(c.Request.Context(), patient)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case domain.ErrNameRequired, domain.ErrInvalidDateOfBirth, domain.ErrInvalidGender:
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to create patient",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient created successfully",
		"patient": patient,
	})
}

func (pc *PatientController) GetPatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid patient ID",
			Message: "Please provide a valid UUID",
			Code:    "validation_error",
		})
		return
	}

	patient, err := pc.PatientUsecase.GetPatient(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		if err == domain.ErrPatientNotFound {
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to get patient",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patient": patient,
	})
}

func (pc *PatientController) GetAllPatients(c *gin.Context) {
	var pagination domain.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid pagination parameters",
			Message: err.Error(),
			Code:    "validation_error",
		})
		return
	}

	patients, totalCount, err := pc.PatientUsecase.GetAllPatients(c.Request.Context(), pagination.Page, pagination.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get patients",
			Message: err.Error(),
			Code:    "internal_error",
		})
		return
	}

	totalPages := (totalCount + pagination.PerPage - 1) / pagination.PerPage

	c.JSON(http.StatusOK, domain.PaginatedPatientsResponse{
		Patients:   patients,
		TotalCount: totalCount,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: totalPages,
	})
}

func (pc *PatientController) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid patient ID",
			Message: "Please provide a valid UUID",
			Code:    "validation_error",
		})
		return
	}

	var request domain.UpdatePatientRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "validation_error",
		})
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid date format",
			Message: "Use YYYY-MM-DD format",
			Code:    "validation_error",
		})
		return
	}

	patient := &domain.Patient{
		ID:          id,
		Name:        request.Name,
		DateOfBirth: dateOfBirth,
		Gender:      domain.Gender(request.Gender),
		IsOffline:   request.IsOffline,
	}

	err = pc.PatientUsecase.UpdatePatient(c.Request.Context(), patient)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		switch err {
		case domain.ErrPatientNotFound:
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		case domain.ErrNameRequired, domain.ErrInvalidDateOfBirth, domain.ErrInvalidGender:
			statusCode = http.StatusBadRequest
			errorCode = "validation_error"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to update patient",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient updated successfully",
		"patient": patient,
	})
}

func (pc *PatientController) DeletePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid patient ID",
			Message: "Please provide a valid UUID",
			Code:    "validation_error",
		})
		return
	}

	err = pc.PatientUsecase.DeletePatient(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "internal_error"

		if err == domain.ErrPatientNotFound {
			statusCode = http.StatusNotFound
			errorCode = "not_found"
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   "Failed to delete patient",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient deleted successfully",
	})
}