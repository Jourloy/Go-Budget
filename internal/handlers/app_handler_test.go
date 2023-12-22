package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAppHandler(t *testing.T) {
	// Load .env
	if err := godotenv.Load(os.ExpandEnv(`../../.env`)); err != nil {
		panic(`error loading .env file`)
	}

	type args struct {
		path   string
		method string
	}
	tests := []struct {
		name        string
		args        args
		wantCode    int
		wantErrBody string
	}{
		{
			name: `Positive`,
			args: args{
				path:   `/live`,
				method: http.MethodGet,
			},
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			g := r.Group(`/`)

			RegisterAppHandler(g)

			req := httptest.NewRequest(tt.args.method, tt.args.path, nil)
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
