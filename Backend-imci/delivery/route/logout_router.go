package route

import (
	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/delivery/middleware"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
)

func NewLogoutRouter(
	env *config.Env,
	Group *gin.RouterGroup,
	blacklistRepo domain.TokenBlacklistRepository,
) {
	authController := middleware.NewAuthController(blacklistRepo)
	
	Group.POST("/logout", authController.Logout)
}