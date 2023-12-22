package impl

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"time"

	"github.com/Jourloy/Go-Budget-Service/internal/budget"
	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type budgetService struct {
	storage *storage.Storage
}

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[budget]`,
		Level:  log.DebugLevel,
	})
)

func CreateBudgetService(storage *storage.Storage) budget.BudgetService {
	return &budgetService{
		storage: storage,
	}
}

type BudgetCreateData struct {
	Name        string `json:"name"`
	Limit       int    `json:"limit"`
	PeriodLimit int    `json:"periodLimit"`
}

func (s *budgetService) CreateBudget(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	// Check body
	var body BudgetCreateData
	if ok := s.parseBody(c, &body); !ok {
		return
	}

	if body.Name == `` || body.Limit == 0 || body.PeriodLimit == 0 {
		logger.Error(`body are invalid`)
		c.String(400, `body are invalid`)
		return
	}

	// Create budget
	err := s.storage.Budget.CreateBudget(&budgets.BudgetCreate{
		Name:        body.Name,
		Limit:       body.Limit,
		PeriodLimit: body.PeriodLimit,
	}, user)

	if err != nil {
		logger.Error(`failed to create budget`, `err`, err)
		c.String(500, `failed to create budget`)
		return
	}

	c.String(200, `ok`)
}

type SpendResponse struct {
	ID          string  `json:"id"`
	Cost        int     `json:"cost"`
	Category    string  `json:"category"`
	IsCredit    bool    `json:"isCredit"`
	Description *string `json:"description"`
	Date        *string `json:"date"`
	Repeat      *string `json:"repeat"`
	CreatedAT   string  `json:"createdAt"`
	UpdatedAT   string  `json:"updatedAt"`
}

type BudgetResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Limit       int             `json:"limit"`
	PeriodLimit int             `json:"periodLimit"`
	StartDate   string          `json:"startDate"`
	CreatedAT   string          `json:"createdAt"`
	UpdatedAT   string          `json:"updatedAt"`
	DaysPassed  int             `json:"daysPassed"`
	DaysLeft    int             `json:"daysLeft"`
	TodayLimit  int             `json:"todayLimit"`
	MonthIncome int             `json:"monthIncome"`
	MonthSpend  int             `json:"monthSpend"`
	WeekLimit   int             `json:"weekLimit"`
	Credit      int             `json:"credit"`
	Spends      []SpendResponse `json:"spends"`
	TodayBudget int             `json:"todayBudget"`
}

func (s *budgetService) GetBudgets(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	// Get budgets
	budgets := s.storage.Budget.GetBudgetsByUserID(user.ID)

	// Sort budgets
	sort.SliceStable(budgets, func(i, j int) bool {
		return budgets[i].ID < budgets[j].ID
	})

	budgetsResponse := []BudgetResponse{}
	for _, budget := range budgets {
		budgetsResponse = append(budgetsResponse, *s.calculateBudget(&budget))
	}

	c.JSON(200, budgetsResponse)
}

func (s *budgetService) calculateBudget(budget *budgets.Budget) *BudgetResponse {
	// Calculate days passed
	startDay, err := time.Parse(`2006-01-02T15:04:05.9Z`, budget.StartDate)
	if err != nil {
		logger.Error(`failed to parse start date`, `err`, err)
		return nil
	}
	today := time.Now()
	daysPassed := int(today.Sub(startDay).Hours()/24) + 1

	// Calculate today limit
	spendsCost := 0
	monthSpend := 0
	monthIncome := 0
	credit := 0

	spendsResponse := []SpendResponse{}
	for _, spend := range budget.Spends {
		spendsResponse = append(spendsResponse, SpendResponse{
			ID:          spend.ID,
			Cost:        spend.Cost,
			Category:    spend.Category,
			Description: spend.Description,
			Date:        spend.Date,
			Repeat:      spend.Repeat,
			IsCredit:    spend.IsCredit,
			CreatedAT:   spend.CreatedAt,
			UpdatedAT:   spend.UpdatedAt,
		})

		// Skip planned spend
		if spend.Date != nil {
			continue
		}

		// Skip credit
		if spend.IsCredit {
			credit += spend.Cost
			continue
		}

		// Calculate spends
		spendsCost += spend.Cost

		// Parse date
		spendDate, err := time.Parse(`2006-01-02T15:04:05.9Z`, spend.CreatedAt)
		if err != nil {
			logger.Error(`failed to parse spend date`, `err`, err)
			return nil
		}

		// Skip old spends
		if spendDate.Month() != today.Month() {
			continue
		}

		// Calculate month stats
		if spend.Cost < 0 {
			monthSpend += spend.Cost
		} else {
			monthIncome += spend.Cost
		}
	}

	todayLimit := budget.PeriodLimit*daysPassed + spendsCost
	todayBudget := budget.Limit + spendsCost + credit
	weekLimit := todayLimit + budget.PeriodLimit*(7-s.getWeekDay(today))

	// Calculate days left
	daysLeft := budget.Limit/budget.PeriodLimit - daysPassed

	return &BudgetResponse{
		ID:          budget.ID,
		Name:        budget.Name,
		Limit:       budget.Limit,
		PeriodLimit: budget.PeriodLimit,
		StartDate:   budget.StartDate,
		CreatedAT:   budget.CreatedAt,
		UpdatedAT:   budget.UpdatedAt,
		DaysPassed:  daysPassed,
		DaysLeft:    daysLeft,
		TodayLimit:  todayLimit,
		Spends:      spendsResponse,
		MonthSpend:  monthSpend,
		MonthIncome: monthIncome,
		WeekLimit:   weekLimit,
		TodayBudget: todayBudget,
		Credit:      credit,
	}
}

func (s *budgetService) getWeekDay(day time.Time) int {
	switch day.Weekday() {
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	case time.Sunday:
		return 7
	}

	return 0
}

type BudgetUpdateData struct {
	Name        string `json:"name,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	PeriodLimit int    `json:"periodLimit,omitempty"`
	StartDate   string `json:"startDate,omitempty"`
}

func (s *budgetService) UpdateBudget(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	// Check body
	var body BudgetUpdateData
	if ok := s.parseBody(c, &body); !ok {
		return
	}

	// Parse budget id
	budgetID := c.Param(`id`)
	if budgetID == `` {
		logger.Error(`budget id is required`)
		c.String(400, `budget id is required`)
		return
	}

	// Update budget
	err := s.storage.Budget.UpdateBudget(&budgets.Budget{
		ID:          budgetID,
		UserID:      user.ID,
		Name:        body.Name,
		Limit:       body.Limit,
		PeriodLimit: body.PeriodLimit,
		StartDate:   body.StartDate,
	})
	if err != nil {
		logger.Error(`failed to update budget`, `err`, err)
		c.String(400, `failed to update budget`)
		return
	}

	c.String(200, `ok`)
}

func (s *budgetService) ChangeDays(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	budgetID := c.Param(`bid`)

	// Get budgets
	budget := s.storage.Budget.GetBudgetByUserIDAndBudgetID(user.ID, budgetID)
	if budget == nil {
		logger.Error(`budget not found`)
		c.String(404, `budget not found`)
		return
	}

	startDate, err := time.Parse(`2006-01-02T15:04:05.9Z`, budget.StartDate)
	if err != nil {
		logger.Error(`failed to parse start date`, `err`, err)
		c.String(500, `failed to parse start date`)
		return
	}

	t := 1 * 24 * time.Hour
	if c.Param(`mod`) == `remove` {
		t = -24 * time.Hour
	}

	newStartDate := startDate.Add(t)
	budget.StartDate = newStartDate.Format(`2006-01-02T15:04:05.9Z`)

	// Update budget
	err = s.storage.Budget.UpdateBudget(budget)
	if err != nil {
		logger.Error(`failed to update budget`, `err`, err)
		c.String(400, `failed to update budget`)
		return
	}

	c.String(200, `ok`)
}

func (s *budgetService) DeleteBudget(c *gin.Context) {
	// Check user
	user, ok := s.check(c)
	if !ok {
		return
	}

	// Parse budget id
	budgetID := c.Param(`id`)
	if budgetID == `` {
		logger.Error(`budget id is required`)
		c.String(400, `budget id is required`)
		return
	}

	// Get budget
	budget := s.storage.Budget.GetBudgetByUserIDAndBudgetID(user.ID, budgetID)
	if budget == nil {
		logger.Error(`budget not found`)
		c.String(404, `budget not found`)
		return
	}

	// Delete budget
	err := s.storage.Budget.DeleteBudget(budget)
	if err != nil {
		logger.Error(`failed to delete budget`, `err`, err)
		c.String(500, `failed to delete budget`)
		return
	}

	c.String(200, `ok`)
}

func (s *budgetService) check(c *gin.Context) (*users.User, bool) {
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

func (s *budgetService) parseBody(c *gin.Context, body interface{}) bool {
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
		c.String(500, `failed to unmarshal body`)
		return false
	}

	return true
}
