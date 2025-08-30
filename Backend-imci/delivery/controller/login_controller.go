// delivery/controller/login_controller.go
package controller

import (
    "net/http"

    "github.com/Afomiat/Digital-IMCI/domain"
    "github.com/gin-gonic/gin"
)

type LoginController struct {
    LoginUsecase domain.LoginUsecase
}

func NewLoginController(loginUsecase domain.LoginUsecase) *LoginController {
    return &LoginController{
        LoginUsecase: loginUsecase,
    }
}

func (lc *LoginController) Login(c *gin.Context) {
    var request domain.LoginRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := lc.LoginUsecase.Login(c.Request.Context(), &request)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "data":    response,
    })
}

func (lc *LoginController) RefreshToken(c *gin.Context) {
    // Implement refresh token endpoint
    c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

func (lc *LoginController) Logout(c *gin.Context) {
    // Implement logout endpoint
    c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}