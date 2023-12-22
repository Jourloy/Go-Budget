package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	budgetM "github.com/Jourloy/Go-Budget-Service/internal/storage/budgets/memory"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	userM "github.com/Jourloy/Go-Budget-Service/internal/storage/users/memory"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAuthHandler(t *testing.T) {
	// Load .env
	if err := godotenv.Load(os.ExpandEnv(`../../.env`)); err != nil {
		panic(`error loading .env file`)
	}

	type args struct {
		path   string
		method string
		body   users.User
	}
	tests := []struct {
		name        string
		args        args
		wantCode    int
		wantErrBody string
	}{
		{
			name: `Negative #1 (Without body)`,
			args: args{
				path:   `/auth/login/`,
				method: http.MethodPost,
			},
			wantCode: 400,
		},
		{
			name: `Negative #2 (Only one field in body)`,
			args: args{
				path:   `/auth/login/`,
				method: http.MethodPost,
				body: users.User{
					Username: `antoxa`,
				},
			},
			wantCode: 400,
		},
		{
			name: `Positive`,
			args: args{
				path:   `/auth/login/`,
				method: http.MethodPost,
				body: users.User{
					Username: `igoryan`,
					Password: `qwerty`,
				},
			},
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			g := r.Group(`/auth/`)

			budgetsStorage := budgetM.CreateBudgetMemory()
			usersStorage := userM.CreateUserMemory()
			spendsStorage := budgetM.CreateSpendMemory()

			s := &storage.Storage{
				User:   usersStorage,
				Budget: budgetsStorage,
				Spend:  spendsStorage,
			}

			RegisterAuthHandler(g, s)

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
