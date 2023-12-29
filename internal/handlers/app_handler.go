package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/Jourloy/Go-Budget-Service/internal/app"
)

func RegisterAppHandler(g *gin.RouterGroup) {
	appService := app.CreateAppService()

	// Проверка работоспособности
	g.GET(`/live`, appService.LiveCheck)
}
