package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	budgetM "github.com/Jourloy/Go-Budget-Service/internal/storage/budgets/memory"
	userM "github.com/Jourloy/Go-Budget-Service/internal/storage/users/memory"
)

// Тест проверяет лишь невозможность создать, получить и модифицировать
// модель не будучи авторизированным в системе
func TestRegisterBudgetHandler(t *testing.T) {
	// Load .env
	if err := godotenv.Load(os.ExpandEnv(`../../.env`)); err != nil {
		panic(`error loading .env file`)
	}

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
			name: `Negative #1 (Create budget)`,
			args: args{
				path:   `/budget/`,
				method: http.MethodPost,
				body: budgets.Budget{
					Name: `Best budget`,
				},
			},
			wantCode: 403,
		},
		{
			name: `Negative #2 (Get budgets)`,
			args: args{
				path:   `/budget/all/`,
				method: http.MethodGet,
			},
			wantCode: 403,
		},
		{
			name: `Negative #3 (Update budget)`,
			args: args{
				path:   `/budget/bid/`,
				method: http.MethodPatch,
			},
			wantCode: 403,
		},
		{
			name: `Negative #4 (Delete budget)`,
			args: args{
				path:   `/budget/bid/`,
				method: http.MethodDelete,
			},
			wantCode: 403,
		},
		{
			name: `Negative #5 (Create spend)`,
			args: args{
				path:   `/budget/bid/spend/`,
				method: http.MethodPost,
			},
			wantCode: 403,
		},
		{
			name: `Negative #6 (Update spend)`,
			args: args{
				path:   `/budget/bid/spend/sid/`,
				method: http.MethodPatch,
			},
			wantCode: 403,
		},
		{
			name: `Negative #7 (Delete spend)`,
			args: args{
				path:   `/budget/bid/spend/sid/`,
				method: http.MethodDelete,
			},
			wantCode: 403,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			g := r.Group(`/budget/`)

			budgetsStorage := budgetM.CreateBudgetMemory()
			usersStorage := userM.CreateUserMemory()
			spendsStorage := budgetM.CreateSpendMemory()

			s := &storage.Storage{
				User:   usersStorage,
				Budget: budgetsStorage,
				Spend:  spendsStorage,
			}

			RegisterBudgetHandler(g, s)

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(tt.args.body); err != nil {
				json.NewEncoder(&buf)
			}
			req := httptest.NewRequest(tt.args.method, tt.args.path, &buf)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
			if tt.wantErrBody != `` {
				assert.Equal(t, tt.wantErrBody, strings.TrimSuffix(rec.Body.String(), `\n`))
				return
			}
		})
	}
}
