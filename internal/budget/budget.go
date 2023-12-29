package budget

import "github.com/gin-gonic/gin"

type BudgetService interface {
	CreateBudget(c *gin.Context)
	GetBudgets(c *gin.Context)
	UpdateBudget(c *gin.Context)
	DeleteBudget(c *gin.Context)
	ChangeDays(c *gin.Context)
}
