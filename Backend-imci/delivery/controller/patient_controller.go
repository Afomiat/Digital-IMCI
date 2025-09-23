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



func (pc *PatientController) CreatePatient(c *gin.Context) {
	var request domain.CreatePatientRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	patient := &domain.Patient{
		Name:        request.Name,
		DateOfBirth: dateOfBirth,
		Gender:      request.Gender,
		IsOffline:   request.IsOffline,
	}

	err = pc.PatientUsecase.CreatePatient(c.Request.Context(), patient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	patient, err := pc.PatientUsecase.GetPatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patient": patient,
	})
}

func (pc *PatientController) GetAllPatients(c *gin.Context) {
	patients, err := pc.PatientUsecase.GetAllPatients(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"count":    len(patients),
	})
}

func (pc *PatientController) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	var request domain.UpdatePatientRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dateOfBirth, err := time.Parse("2006-01-02", request.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	patient := &domain.Patient{
		ID:          id,
		Name:        request.Name,
		DateOfBirth: dateOfBirth,
		Gender:      request.Gender,
		IsOffline:   request.IsOffline,
	}

	err = pc.PatientUsecase.UpdatePatient(c.Request.Context(), patient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	err = pc.PatientUsecase.DeletePatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient deleted successfully",
	})
}