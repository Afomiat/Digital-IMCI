// controller/password_reset_controller.go
package controller

import (
	"net/http"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
)

type PasswordResetController struct {
	PasswordResetUsecase domain.PasswordResetUsecase
}

func NewPasswordResetController(passwordResetUsecase domain.PasswordResetUsecase) *PasswordResetController {
	return &PasswordResetController{
		PasswordResetUsecase: passwordResetUsecase,
	}
}



func (c *PasswordResetController) ForgotPassword(ctx *gin.Context) {
	var req domain.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.PasswordResetUsecase.InitiatePasswordReset(ctx, req.Phone, req.UseWhatsApp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	method := "whatsapp"
	if !req.UseWhatsApp {
		method = "telegram"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset OTP has been sent",
		"phone":   req.Phone,
		"method":  method,
	})
}

func (c *PasswordResetController) VerifyResetOTP(ctx *gin.Context) {
	var req domain.VerifyResetOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValid, err := c.PasswordResetUsecase.VerifyPasswordResetOTP(ctx, req.Phone, req.OTP)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "OTP verified successfully",
		"phone":    req.Phone,
		"verified": true,
	})
}

func (c *PasswordResetController) ResetPassword(ctx *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	isValid, err := c.PasswordResetUsecase.VerifyPasswordResetOTP(ctx, req.Phone, req.OTP)
	if err != nil || !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	if err := c.PasswordResetUsecase.ResetPassword(ctx, req.Phone, req.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
		"phone":   req.Phone,
	})
}