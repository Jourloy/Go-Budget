package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type appService struct{}

var startTime time.Time

// init задает время старта сервера
func init() {
	startTime = time.Now()
}

func CreateAppService() *appService {
	return &appService{}
}

// LiveCheck возвращает текущий статус сервера
func (s *appService) LiveCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		`status`: `OK`,
		`uptime`: time.Since(startTime) / time.Second,
	})
}
