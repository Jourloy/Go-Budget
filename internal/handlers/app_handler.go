package handlers

import (
	"github.com/Jourloy/Go-Budget-Service/internal/app"
	"github.com/gin-gonic/gin"
)

func RegisterAppHandler(g *gin.RouterGroup) {
	appService := app.CreateAppService()

	// Проверка работоспособности
	g.GET(`/live`, appService.LiveCheck)
}
