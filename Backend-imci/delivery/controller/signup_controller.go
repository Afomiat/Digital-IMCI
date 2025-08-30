// delivery/controller/signup_controller.go
package controller

import (
	"net/http"
	"sync"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
)

// Temporary storage for signup data
var tempSignupData = struct {
	sync.RWMutex
	data map[string]domain.SignupForm
}{
	data: make(map[string]domain.SignupForm),
}

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	env           *config.Env
}

func NewSignupController(signupUsecase domain.SignupUsecase, env *config.Env) *SignupController {
	return &SignupController{
		SignupUsecase: signupUsecase,
		env:           env,
	}
}

func (sc *SignupController) Signup(ctx *gin.Context) {
	var form domain.SignupForm

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if phone already exists
	existingProfessional, _ := sc.SignupUsecase.GetMedicalProfessionalByPhone(ctx, form.Phone)
	if existingProfessional != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already registered!"})
		return
	}

	// Store the signup data temporarily
	tempSignupData.Lock()
	tempSignupData.data[form.Phone] = form
	tempSignupData.Unlock()

	// Create a temporary professional object for OTP
	professional := &domain.MedicalProfessional{
		FullName: form.FullName,
		Phone:    form.Phone,
		Role:     form.Role,
	}

	err := sc.SignupUsecase.SendOtp(ctx, professional)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OTP sending failed: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (sc *SignupController) Verify(ctx *gin.Context) {
	var otp domain.VerifyOtp

	if err := ctx.ShouldBindJSON(&otp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify OTP
	_, err := sc.SignupUsecase.VerifyOtp(ctx, &otp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the stored signup data
	tempSignupData.RLock()
	form, exists := tempSignupData.data[otp.Phone]
	tempSignupData.RUnlock()

	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Signup data not found. Please sign up again."})
		return
	}

	// Register the user automatically after successful verification
	professionalID, err := sc.SignupUsecase.RegisterMedicalProfessional(ctx, &form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Registration failed: " + err.Error()})
		return
	}

	// Clean up temporary data
	tempSignupData.Lock()
	delete(tempSignupData.data, otp.Phone)
	tempSignupData.Unlock()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Registration completed successfully",
		"user_id": professionalID,
	})
}