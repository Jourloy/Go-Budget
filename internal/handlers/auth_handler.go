package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/Jourloy/Go-Budget-Service/internal/auth/impl"
	"github.com/Jourloy/Go-Budget-Service/internal/storage"
)

func RegisterAuthHandler(g *gin.RouterGroup, s *storage.Storage) {
	authService := impl.CreateAuthService(s)

	g.POST(`/login/`, authService.Login)
	g.POST(`/register/`, authService.Register)
	g.POST(`/user/`, authService.GetUserData)
}
