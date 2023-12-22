package memory

import (
	"testing"

	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	"github.com/stretchr/testify/assert"
)

func TestBudgetMemory_CreateBudget(t *testing.T) {
	type args struct {
		budget *budgets.BudgetCreate
		user   *users.User
	}
	tests := []struct {
		name    string
		m       *BudgetMemory
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				budget: &budgets.BudgetCreate{
					Name:        `budget`,
					Limit:       10000,
					PeriodLimit: 1000,
				},
				user: &users.User{
					ID: `testing`,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateBudgetMemory()

			// Создание бюджета
			if err := m.CreateBudget(tt.args.budget, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("BudgetMemory.CreateBudget() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка его наличия
			b := m.GetBudgetsByUserID(`testing`)
			assert.Equal(t, 1, len(b))
		})
	}
}

func TestBudgetMemory_GetBudgetsByUserID(t *testing.T) {
	type args struct {
		budget *budgets.BudgetCreate
		user   *users.User
	}
	tests := []struct {
		name    string
		m       *BudgetMemory
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				budget: &budgets.BudgetCreate{
					Name:        `budget`,
					Limit:       10000,
					PeriodLimit: 1000,
				},
				user: &users.User{
					ID: `testing`,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateBudgetMemory()

			// Создание бюджета
			if err := m.CreateBudget(tt.args.budget, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("BudgetMemory.CreateBudget() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка его наличия
			b := m.GetBudgetsByUserID(`testing`)
			assert.Equal(t, 1, len(b))

			// Проверка если ID неверный
			newB := m.GetBudgetsByUserID(`test`)
			assert.Equal(t, 0, len(newB))
		})
	}
}

func TestBudgetMemory_UpdateBudget(t *testing.T) {
	type args struct {
		budget *budgets.BudgetCreate
		user   *users.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				budget: &budgets.BudgetCreate{
					Name:        `budget`,
					Limit:       10000,
					PeriodLimit: 1000,
				},
				user: &users.User{
					ID: `testing`,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateBudgetMemory()

			// Создание бюджета
			if err := m.CreateBudget(tt.args.budget, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("BudgetMemory.CreateBudget() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка имени
			b := m.GetBudgetsByUserID(`testing`)
			assert.Equal(t, b[0].Name, tt.args.budget.Name)

			// Обновление имени
			b[0].Name = `The best budget`
			m.UpdateBudget(&b[0])

			// Проверка обновления
			bNew := m.GetBudgetsByUserID(`testing`)
			assert.Equal(t, bNew[0].Name, `The best budget`)
		})
	}
}

func TestBudgetMemory_DeleteBudget(t *testing.T) {
	type args struct {
		budget *budgets.BudgetCreate
		user   *users.User
	}
	tests := []struct {
		name    string
		m       *BudgetMemory
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				budget: &budgets.BudgetCreate{
					Name:        `budget`,
					Limit:       10000,
					PeriodLimit: 1000,
				},
				user: &users.User{
					ID: `testing`,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateBudgetMemory()

			// Создание бюджета
			if err := m.CreateBudget(tt.args.budget, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("BudgetMemory.CreateBudget() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка его наличия
			b := m.GetBudgetsByUserID(`testing`)
			assert.Equal(t, 1, len(b))

			// Удаление
			m.DeleteBudget(&b[0])

			// Проверка удаления
			newB := m.GetBudgetsByUserID(`testing`)
			assert.Equal(t, 0, len(newB))
		})
	}
}

func TestBudgetMemory_CreateSpendt(t *testing.T) {
	type args struct {
		spend    *budgets.SpendCreate
		budgetID string
		userID   string
	}
	tests := []struct {
		name    string
		m       *BudgetMemory
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				spend: &budgets.SpendCreate{
					Cost:     100,
					Category: `food`,
				},
				budgetID: `test_bid`,
				userID:   `test_uid`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateSpendMemory()

			// Создание операции
			if err := m.CreateSpend(tt.args.spend, tt.args.budgetID, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("SpendMemory.CreateSpend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
