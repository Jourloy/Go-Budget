package impl

import (
	"encoding/json"
	"io"
	"os"

	"github.com/Jourloy/Go-Budget-Service/internal/spend"
	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type spendService struct {
	storage *storage.Storage
}

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[spend]`,
		Level:  log.DebugLevel,
	})
)

func CreateSpendService(storage *storage.Storage) spend.SpendsService {
	return &spendService{
		storage: storage,
	}
}

type SpendCreateData struct {
	Cost        int     `json:"cost"`
	BudgetID    string  `json:"budgetId"`
	Category    string  `json:"category"`
	IsCredit    bool    `json:"isCredit"`
	Description *string `json:"description"`
	Date        *string `json:"date"`
	Repeat      *string `json:"repeat"`
}

func (s *spendService) CreateSpend(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	// Get IDs
	budgetID := c.Param(`bid`)

	// Check body
	var body SpendCreateData
	if ok := s.parseBody(c, &body); !ok {
		return
	}

	// Get budget
	budget := s.storage.Budget.GetBudgetByUserIDAndBudgetID(user.ID, budgetID)
	if budget == nil {
		logger.Error(`budget not found`)
		c.String(404, `budget not found`)
		return
	}

	logger.Debug(`debug`, `date`, body.Date)

	// Create spend
	if err := s.storage.Spend.CreateSpend(&budgets.SpendCreate{
		Cost:        body.Cost,
		Category:    body.Category,
		Description: body.Description,
		Date:        body.Date,
		Repeat:      body.Repeat,
		IsCredit:    body.IsCredit,
	}, budget.ID, user.ID); err != nil {
		logger.Error(`failed to create spend`, `err`, err)
		c.String(500, `failed to create spend`)
		return
	}

	c.String(200, `ok`)
}

func (s *spendService) UpdateSpend(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	// Get IDs
	budgetID := c.Param(`bid`)
	spendID := c.Param(`sid`)

	if budgetID == `` || spendID == `` {
		logger.Error(`budget id or spend id not found`)
		c.String(400, `budget id or spend id not found`)
		return
	}

	// Parse body
	var body budgets.Spend
	if ok := s.parseBody(c, &body); !ok {
		return
	}

	if err := s.storage.Spend.UpdateSpend(&body, budgetID, user.ID); err != nil {
		logger.Error(`failed to update spend`, `err`, err)
		c.String(500, `failed to update spend`)
		return
	}

	c.String(200, `ok`)
}

func (s *spendService) DeleteSpend(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	// Get IDs
	budgetID := c.Param(`bid`)
	spendID := c.Param(`sid`)

	if budgetID == `` || spendID == `` {
		logger.Error(`budget id or spend id not found`)
		c.String(400, `budget id or spend id not found`)
		return
	}

	if err := s.storage.Spend.DeleteSpend(spendID, budgetID, user.ID); err != nil {
		logger.Error(`failed to delete spend`, `err`, err)
		c.String(400, `failed to delete spend`)
		return
	}

	c.String(200, `ok`)
}

func (s *spendService) check(c *gin.Context) (*users.User, bool) {
	// Check user
	username, exist := c.Get(`username`)
	if !exist {
		logger.Error(`failed to get username`)
		c.String(403, `failed to get username`)
		return nil, false
	}

	// Get user
	user, err := s.storage.User.GetUserByUsername(username.(string))
	if user == nil || err != nil {
		logger.Error(`failed to get user`)
		c.String(500, `failed to get user`)
		return nil, false
	}

	return user, true
}

func (s *spendService) parseBody(c *gin.Context, body interface{}) bool {
	// Check body
	if c.Request.Body == nil {
		logger.Error(`body not found`)
		c.String(400, `body not found`)
		return false
	}

	// Read body
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(`failed to read body`, `err`, err)
		c.String(500, `failed to read body`)
		return false
	}

	// Unmarshal
	if err := json.Unmarshal(b, &body); err != nil {
		logger.Error(`failed to unmarshal body`, `err`, err)
		c.String(400, `failed to unmarshal body`)
		return false
	}

	return true
}
