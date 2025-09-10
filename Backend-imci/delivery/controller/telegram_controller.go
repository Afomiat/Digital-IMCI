package controller

import (
	"net/http"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/gin-gonic/gin"
)

type TelegramController struct {
	TelegramService domain.TelegramService
}

func NewTelegramController(telegramService domain.TelegramService) *TelegramController {
	return &TelegramController{
		TelegramService: telegramService,
	}
}

func (tc *TelegramController) GetStartLink(ctx *gin.Context) {
	startLink := tc.TelegramService.GetStartLink()
	
	ctx.JSON(http.StatusOK, gin.H{
		"start_link": startLink,
		"message":    "Use this link to start the Telegram bot without typing /start",
	})
}