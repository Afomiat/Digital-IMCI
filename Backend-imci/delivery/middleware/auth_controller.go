// middleware/auth_controller.go
package middleware

import (
	"net/http"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	blacklistRepo domain.TokenBlacklistRepository
}

func NewAuthController(blacklistRepo domain.TokenBlacklistRepository) *AuthController {
	return &AuthController{
		blacklistRepo: blacklistRepo,
	}
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	phone, _ := ctx.Get("phone")
	role, _ := ctx.Get("role")
	token, exists := ctx.Get("token")
	
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token not found in context"})
		return
	}

	tokenString, ok := token.(string)
	if !ok || tokenString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	if ac.blacklistRepo != nil {
		expiration := 24 * time.Hour 
		err := ac.blacklistRepo.BlacklistToken(ctx.Request.Context(), tokenString, expiration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to invalidate token",
				"details": "Please try again or contact support if the issue persists",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Logged out successfully",
		"details":      "Your session has been terminated and access token invalidated.",
		"user_id":      userID,
		"phone":        phone,
		"role":         role,
		"logged_out":   true,
		"token_invalidated": true,
		"timestamp":    time.Now().Format(time.RFC3339),
	})
}

