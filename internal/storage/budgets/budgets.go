package budgets

import "github.com/Jourloy/Go-Budget-Service/internal/storage/users"

type Budget struct {
	ID          string `db:"id" json:"id"`
	UserID      string `db:"user_id" json:"user_id"`
	Name        string `db:"name" json:"name"`
	Limit       int    `db:"current_limit" json:"limit"`
	PeriodLimit int    `db:"period_limit" json:"periodLimit"`
	StartDate   string `db:"start_date" json:"startDate"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"updatedAt"`
	Spends      []Spend
}

type BudgetCreate struct {
	Name        string
	Limit       int
	PeriodLimit int
	StartDate   string
}

type BudgetStorage interface {
	CreateBudget(budget *BudgetCreate, user *users.User) error
	GetBudgetsByUserID(userID string) []Budget
	GetBudgetByUserIDAndBudgetID(userID string, budgetID string) *Budget
	UpdateBudget(budget *Budget) error
	DeleteBudget(budget *Budget) error
}

type Spend struct {
	ID          string  `db:"id"`
	UserID      string  `db:"user_id"`
	BudgetID    string  `db:"budget_id"`
	Cost        int     `db:"cost"`
	Category    string  `db:"category"`
	IsCredit    bool    `db:"is_credit"`
	Description *string `db:"description"`
	Date        *string `db:"date"`
	Repeat      *string `db:"repeat"`
	CreatedAt   string  `db:"created_at"`
	UpdatedAt   string  `db:"updated_at"`
}

type SpendCreate struct {
	Cost        int
	Category    string
	IsCredit    bool
	Description *string
	Date        *string
	Repeat      *string
}

type SpendStorage interface {
	CreateSpend(spend *SpendCreate, budgetID string, userID string) error
	UpdateSpend(spend *Spend, userID string, budgetID string) error
	DeleteSpend(spendID string, budgetID string, userID string) error
}
