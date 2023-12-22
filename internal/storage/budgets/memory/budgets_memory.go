package memory

import (
	"sync"

	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	"github.com/google/uuid"
)

type BudgetMemory struct {
	sync.Mutex
	Budgets map[string]budgets.Budget
}

// Это хранилище используется исключительно для тестов логики сервисов.
// Определенные методы не реализуют должную систему проверки.
//
// В качестве кэш хранилища лучше реализовать Redis.
func CreateBudgetMemory() budgets.BudgetStorage {
	budgets := make(map[string]budgets.Budget)

	return &BudgetMemory{
		Budgets: budgets,
	}
}

// CreateBudget создает новый бюджет.
//
// Параметры:
// - budget: указатель на budgets.BudgetCreate структуру.
// - user: указатель на users.User струкруту.
//
// Возвращает:
// - error: всегда nil.
func (m *BudgetMemory) CreateBudget(budget *budgets.BudgetCreate, user *users.User) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	id := uuid.NewString()
	b := budgets.Budget{
		ID:          id,
		UserID:      user.ID,
		Name:        budget.Name,
		Limit:       budget.Limit,
		PeriodLimit: budget.PeriodLimit,
	}

	m.Budgets[id] = b

	return nil
}

// GetBudgetsByUserID ищет все бюджеты пользователя.
//
// Параметры:
// - userID: ID пользователя.
//
// Возвращает:
// - []budgets.Budget: массив всех бюджетов.
func (m *BudgetMemory) GetBudgetsByUserID(userID string) []budgets.Budget {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	budgets := []budgets.Budget{}

	for _, b := range m.Budgets {
		if b.UserID == userID {
			budgets = append(budgets, b)
		}
	}

	return budgets
}

// GetBudgetByUserIDAndBudgetID находит по user id и budget id нужный бюджет.
//
// Параметры:
// - userID: ID пользователя.
// - budgetID: ID бюджета.
//
// Возвращает:
// - *budgets.Budget: указатель на budgets.Budget.
func (m *BudgetMemory) GetBudgetByUserIDAndBudgetID(userID string, budgetID string) *budgets.Budget {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for _, b := range m.Budgets {
		if b.UserID == userID && b.ID == budgetID {
			return &b
		}
	}

	return nil
}

// UpdateBudget обновляет бюджет.
//
// Параметры:
// - budget: указатель на budgets.Budget.
//
// Возвращает:
// - error: всегда nil.
func (m *BudgetMemory) UpdateBudget(budget *budgets.Budget) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Budgets[budget.ID] = *budget

	return nil
}

// DeleteBudget удаляет бюджет.
//
// Параметры:
// - budget: указатель на budgets.Budget
//
// Возвращает:
// - error: всегда nil.
func (m *BudgetMemory) DeleteBudget(budget *budgets.Budget) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	delete(m.Budgets, budget.ID)

	return nil
}

type SpendMemory struct {
	sync.Mutex
	Spends map[string]budgets.Spend
}

// Это хранилище используется исключительно для тестов логики сервисов.
// Определенные методы не реализуют должную систему проверки.
//
// В качестве кэш хранилища лучше реализовать Redis.
func CreateSpendMemory() budgets.SpendStorage {
	spends := make(map[string]budgets.Spend)

	return &SpendMemory{
		Spends: spends,
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
// - error: всегда nil.
func (m *SpendMemory) CreateSpend(spend *budgets.SpendCreate, budgetID string, userID string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	id := uuid.NewString()
	s := &budgets.Spend{
		ID:          id,
		UserID:      userID,
		BudgetID:    budgetID,
		Cost:        spend.Cost,
		Category:    spend.Category,
		Description: spend.Description,
		Date:        spend.Date,
		Repeat:      spend.Repeat,
	}

	m.Spends[id] = *s

	return nil
}

// UpdateSpend обновляет операцию.
//
// Параметры:
// - spend: указатель на budgets.Spend.
// - budgetID: ID бюджета, в котором создана операция. Здесь не используется.
// - userID: ID пользователя. Здесь не используется.
//
// Возвращает:
// - error: всегда nil.
func (m *SpendMemory) UpdateSpend(spend *budgets.Spend, userID string, budgetID string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Spends[spend.ID] = *spend

	return nil
}

// DeleteSpend удаляет операцию.
//
// Параметры:
// - spendID: ID операции.
// - budgetID: ID бюджета, в котором создана операция. Здесь не используется.
// - userID: ID пользователя, создавшего операцию. Здесь не используется.
//
// Возвращает:
// - error: всегда nil.
func (m *SpendMemory) DeleteSpend(spendID string, budgetID string, userID string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	delete(m.Spends, spendID)

	return nil
}
