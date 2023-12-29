package handlers

import (
	"github.com/gin-gonic/gin"

	budgetS "github.com/Jourloy/Go-Budget-Service/internal/budget/impl"
	spendS "github.com/Jourloy/Go-Budget-Service/internal/spend/impl"
	"github.com/Jourloy/Go-Budget-Service/internal/storage"
)

func RegisterBudgetHandler(g *gin.RouterGroup, s *storage.Storage) {
	budgetService := budgetS.CreateBudgetService(s)
	spendService := spendS.CreateSpendService(s)

	g.POST(`/`, budgetService.CreateBudget)
	g.GET(`/all/`, budgetService.GetBudgets)
	g.PATCH(`/:bid/`, budgetService.UpdateBudget)
	g.DELETE(`/:bid/`, budgetService.DeleteBudget)
	g.POST(`/:bid/days/:mod/`, budgetService.ChangeDays)

	g.POST(`/:bid/spend/`, spendService.CreateSpend)
	g.PATCH(`/:bid/spend/:sid/`, spendService.UpdateSpend)
	g.DELETE(`:bid/spend/:sid/`, spendService.DeleteSpend)
}
