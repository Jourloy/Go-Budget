package handlers

import (
	"github.com/Jourloy/Go-Budget-Service/internal/auth/impl"
	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	"github.com/gin-gonic/gin"
)

func RegisterAuthHandler(g *gin.RouterGroup, s *storage.Storage) {
	authService := impl.CreateAuthService(s)

	g.POST(`/login/`, authService.Login)
	g.POST(`/register/`, authService.Register)
	g.POST(`/user/`, authService.GetUserData)
}
