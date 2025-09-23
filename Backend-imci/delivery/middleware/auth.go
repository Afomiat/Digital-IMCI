// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/Afomiat/Digital-IMCI/config"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	env           *config.Env
	blacklistRepo domain.TokenBlacklistRepository
}

func NewAuthMiddleware(env *config.Env, blacklistRepo domain.TokenBlacklistRepository) *AuthMiddleware {
	return &AuthMiddleware{
		env:           env,
		blacklistRepo: blacklistRepo,
	}
}

func (am *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/v1/login" || c.Request.URL.Path == "/api/v1/signup" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		if am.blacklistRepo != nil {
			blacklisted, err := am.blacklistRepo.IsTokenBlacklisted(c.Request.Context(), tokenString)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}
			if blacklisted {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
				c.Abort()
				return
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(am.env.AccessTokenSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("userID", claims["id"])
		c.Set("phone", claims["phone"])
		c.Set("role", claims["role"])
		c.Set("token", tokenString)

		c.Next()
	}
}