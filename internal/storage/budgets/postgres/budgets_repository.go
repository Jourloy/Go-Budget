package postgres

import (
	"database/sql"
	"errors"
	"os"

	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BudgetRepository struct {
	db *sqlx.DB
}

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Prefix: `[database-budget]`,
	Level:  log.DebugLevel,
})

func CreateBudgetRepository(db *sqlx.DB) budgets.BudgetStorage {
	return &BudgetRepository{
		db: db,
	}
}

// CreateBudget создает новый бюджет.
//
// Параметры:
// - budget: указатель на budgets.BudgetCreate структуру.
// - user: указатель на users.User струкруту.
//
// Возвращает:
// - error: ошибка при создании бюджета
func (r *BudgetRepository) CreateBudget(budget *budgets.BudgetCreate, user *users.User) error {
	b := budgets.Budget{
		ID:          uuid.NewString(),
		UserID:      user.ID,
		Name:        budget.Name,
		Limit:       budget.Limit,
		PeriodLimit: budget.PeriodLimit,
	}

	_, err := r.db.NamedExec(
		`INSERT INTO budgets (id, user_id, name, current_limit, period_limit) VALUES (:id, :user_id, :name, :current_limit, :period_limit)`,
		b,
	)

	return err
}

// GetBudgetsByUserID ищет все бюджеты пользователя.
//
// Параметры:
// - userID: ID пользователя.
//
// Возвращает:
// - []budgets.Budget: массив всех бюджетов. При ошибке вернет nil.
func (r *BudgetRepository) GetBudgetsByUserID(userID string) []budgets.Budget {
	budgets := []budgets.Budget{}

	// Получение бюджетов
	err := r.db.Select(&budgets, `SELECT * FROM budgets WHERE budgets.user_id = $1 JOIN spends ON spends.budget_id = budgets.id`, userID)

	// Если ошибка
	if err != nil {
		logger.Error(`failed to get budgets`, `err`, err)
		return nil
	}

	return budgets
}

func (r *BudgetRepository) GetBudgetByUserIDAndBudgetID(userID string, budgetID string) *budgets.Budget {
	budget := budgets.Budget{}

	err := r.db.Get(
		&budget,
		`SELECT * FROM budgets WHERE budgets.user_id = $1 AND budgets.id = $2 JOIN spends ON spends.budget_id = budgets.id`,
		userID, budgetID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Error(`budget not found`)
		return nil
	}

	if err != nil {
		logger.Error(`failed to get budget`, `err`, err)
		return nil
	}

	return &budget
}

// UpdateBudget обновляет бюджет.
//
// Параметры:
// - budget: указатель на budgets.Budget.
//
// Возвращает:
// - error: ошибка во время обновления.
func (r *BudgetRepository) UpdateBudget(budget *budgets.Budget) error {
	_, err := r.db.NamedExec(
		`UPDATE budgets SET name = :name, current_limit = :current_limit, period_limit = :period_limit, start_date = :start_date WHERE id = :id AND user_id = :user_id`,
		budget,
	)
	return err
}

// DeleteBudget удаляет бюджет.
//
// Параметры:
// - budget: указатель на budgets.Budget
//
// Возвращает:
// - error: ошибка во время обновления
func (r *BudgetRepository) DeleteBudget(budget *budgets.Budget) error {
	_, err := r.db.NamedExec(`DELETE FROM budgets WHERE id = :id AND user_id = :user_id`, budget)
	return err
}

type SpendRepository struct {
	db *sqlx.DB
}

func CreateSpendRepository(db *sqlx.DB) budgets.SpendStorage {
	return &SpendRepository{
		db: db,
	}
}

// CreateSpend создает новую операцию.
//
// Параметры:
// - spend: указатель на budgets.SpendCreate структуру.
// - budgetID: ID бюджета, в котором создана операция.
// - userID: ID пользователя, который создал операцию.
//
// Возвращает:
// - error: ошибка при создании операции.
func (r *SpendRepository) CreateSpend(spend *budgets.SpendCreate, budgetID string, userID string) error {
	s := &budgets.Spend{
		ID:          uuid.NewString(),
		UserID:      userID,
		BudgetID:    budgetID,
		Cost:        spend.Cost,
		Category:    spend.Category,
		Description: spend.Description,
		Date:        spend.Date,
		Repeat:      spend.Repeat,
	}

	_, err := r.db.NamedExec(
		`INSERT INTO spends (id, budget_id, cost, category, description, date, repeat, is_credit) VALUES (:id, :budget_id, :cost, :category, :description, :date, :repeat, :is_credit)`,
		s,
	)

	return err
}

// UpdateSpend обновляет операцию.
//
// Параметры:
// - spend: указатель на budgets.Spend.
// - budgetID: ID бюджета, в котором создана операция.
// - userID: ID пользователя.
//
// Возвращает:
// - error: ошибка во время обновления.
func (r *SpendRepository) UpdateSpend(spend *budgets.Spend, budgetID string, userID string) error {
	_, err := r.db.Exec(
		`UPDATE spends SET cost = $1, category = $2, description = $3, date = $4, repeat = $5, is_credit = $6 WHERE id = $7 AND budget_id = $8 AND user_id = $9 AND budget_id = $10`,
		spend.Cost, spend.Category, spend.Description, spend.Date, spend.Repeat, spend.IsCredit, spend.ID, spend.BudgetID, userID, budgetID)
	return err
}

// DeleteSpend удаляет операцию.
//
// Параметры:
// - spendID: ID операции.
// - budgetID: ID бюджета, в котором создана операция.
// - userID: ID пользователя, создавшего операцию.
//
// Возвращает:
// - error: ошибка во время удаления.
func (r *SpendRepository) DeleteSpend(spendID string, budgetID string, userID string) error {
	_, err := r.db.Exec(
		`DELETE FROM spends WHERE id = $1 AND budget_id = $2 AND user_id = $3`,
		spendID, budgetID, userID,
	)
	return err
}
