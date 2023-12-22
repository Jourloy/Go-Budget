package auth

import "github.com/gin-gonic/gin"

type AuthService interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetUserData(c *gin.Context)
}
