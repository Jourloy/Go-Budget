package impl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	budgetM "github.com/Jourloy/Go-Budget-Service/internal/storage/budgets/memory"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	userM "github.com/Jourloy/Go-Budget-Service/internal/storage/users/memory"
)

func Test_budgetService_CreateBudget(t *testing.T) {
	type args struct {
		path   string
		method string
		body   interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantCode    int
		wantErrBody string
	}{
		{
			name: `Negative (Create budget)`,
			args: args{
				path:   `/budget/`,
				method: http.MethodPost,
				body: budgets.Budget{
					Name:  `Best budget`,
					Limit: 35000,
				},
			},
			wantCode: 400,
		},
		{
			name: `Positive (Create budget)`,
			args: args{
				path:   `/budget/`,
				method: http.MethodPost,
				body: budgets.Budget{
					Name:        `Best budget`,
					Limit:       35000,
					PeriodLimit: 5000,
				},
			},
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация хранилища
			budgetsStorage := budgetM.CreateBudgetMemory()
			usersStorage := userM.CreateUserMemory()
			spendsStorage := budgetM.CreateSpendMemory()

			s := &storage.Storage{
				User:   usersStorage,
				Budget: budgetsStorage,
				Spend:  spendsStorage,
			}

			// Создание тестового пользователя
			s.User.CreateUser(&users.UserCreate{
				Username: `test_user`,
				Password: `1234`,
				Role:     `test`,
			})

			// Сервис управления бюджетом
			b := CreateBudgetService(s)

			// Body в reader и создание запроса
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(tt.args.body); err != nil {
				json.NewEncoder(&buf)
			}
			req := httptest.NewRequest(tt.args.method, tt.args.path, &buf)
			rec := httptest.NewRecorder()

			// Создание контекста
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			// Установка тестового пользователя
			c.Set(`username`, `test_user`)
			c.Set(`role`, `test`)

			// Создание бюджета
			b.CreateBudget(c)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}

func Test_budgetService_GetBudgets(t *testing.T) {
	type args struct {
		path   string
		method string
		body   interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantCode    int
		wantErrBody string
	}{
		{
			name: `Positive (Get all budgets)`,
			args: args{
				path:   `/budget/all/`,
				method: http.MethodPost,
				body: budgets.Budget{
					Name:        `Best budget`,
					Limit:       35000,
					PeriodLimit: 5000,
				},
			},
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация хранилища
			budgetsStorage := budgetM.CreateBudgetMemory()
			usersStorage := userM.CreateUserMemory()
			spendsStorage := budgetM.CreateSpendMemory()

			s := &storage.Storage{
				User:   usersStorage,
				Budget: budgetsStorage,
				Spend:  spendsStorage,
			}

			// Создание тестового пользователя
			s.User.CreateUser(&users.UserCreate{
				Username: `test_user`,
				Password: `1234`,
				Role:     `test`,
			})

			// Сервис управления бюджетом
			b := CreateBudgetService(s)

			// Body в reader и создание запроса
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(tt.args.body); err != nil {
				json.NewEncoder(&buf)
			}
			req := httptest.NewRequest(tt.args.method, tt.args.path, &buf)
			rec := httptest.NewRecorder()

			// Создание контекста
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			// Установка тестового пользователя
			c.Set(`username`, `test_user`)
			c.Set(`role`, `test`)

			// Получение бюджета
			b.GetBudgets(c)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}

func Test_budgetService_UpdateBudget(t *testing.T) {
	type args struct {
		path   string
		method string
		body   interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantCode    int
		wantErrBody string
	}{
		{
			name: `Positive (Update budget)`,
			args: args{
				path:   `/budget/testing`,
				method: http.MethodPatch,
				body: budgets.Budget{
					Name:        `Best budget`,
					Limit:       35000,
					PeriodLimit: 5000,
				},
			},
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация хранилища
			budgetsStorage := budgetM.CreateBudgetMemory()
			usersStorage := userM.CreateUserMemory()
			spendsStorage := budgetM.CreateSpendMemory()

			s := &storage.Storage{
				User:   usersStorage,
				Budget: budgetsStorage,
				Spend:  spendsStorage,
			}

			// Создание тестового пользователя
			s.User.CreateUser(&users.UserCreate{
				Username: `test_user`,
				Password: `1234`,
				Role:     `test`,
			})

			// Сервис управления бюджетом
			b := CreateBudgetService(s)

			// ПОДГОТОВКА БЮДЖЕТА //

			// Body в reader и создание запроса
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(tt.args.body); err != nil {
				json.NewEncoder(&buf)
			}
			req := httptest.NewRequest(http.MethodPost, `/budget/`, &buf)
			rec := httptest.NewRecorder()

			// Создание контекста
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			// Установка тестового пользователя
			c.Set(`username`, `test_user`)
			c.Set(`role`, `test`)

			// Создание тестового бюджета
			b.CreateBudget(c)

			// ПРОВЕРКА ОБНОВЛЕНИЯ БЮДЖЕТА //

			// Body в reader и создание запроса
			var newBuf bytes.Buffer
			if err := json.NewEncoder(&newBuf).Encode(budgets.Budget{
				Name: `The best budget`,
			}); err != nil {
				json.NewEncoder(&newBuf)
			}
			newReq := httptest.NewRequest(tt.args.method, tt.args.path, &newBuf)
			newRec := httptest.NewRecorder()

			// Создание контекста
			newC, _ := gin.CreateTestContext(newRec)
			newC.Request = newReq

			// Установка тестового пользователя
			newC.Set(`username`, `test_user`)
			newC.Set(`role`, `test`)

			// Симуляция парсинга URL от gin
			newC.AddParam(`id`, `testing`)

			b.UpdateBudget(newC)

			assert.Equal(t, tt.wantCode, newRec.Code)
		})
	}
}
