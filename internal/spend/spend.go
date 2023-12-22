package spend

import "github.com/gin-gonic/gin"

type SpendsService interface {
	CreateSpend(c *gin.Context)
	UpdateSpend(c *gin.Context)
	DeleteSpend(c *gin.Context)
}
