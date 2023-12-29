package middlewares

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtLogger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[middleware-jwt]`,
		Level:  log.DebugLevel,
	})

	Secret string
)

// parseENV парсит переменные окружения.
func parseENV() {
	env, exist := os.LookupEnv(`SECRET`)
	if !exist {
		jwtLogger.Fatal(`Error loading SECRET from .env file`)
	}
	Secret = env
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		parseENV()

		// Получает токены из кук
		a, err := c.Cookie(`access_token`)
		if err != nil {
			c.Next()
		}

		r, err := c.Cookie(`refresh_token`)
		if err != nil {
			c.Next()
		}

		// Проверяет наличие токенов
		if a == `` || r == `` {
			c.Next()
		}

		// Декодирование и верификация
		claims := jwt.MapClaims{}
		if _, err := jwt.ParseWithClaims(a, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(Secret), nil
		}); err != nil {
			c.Next()
		}

		// Проверка ника и роли
		var username string
		var role string

		for key, val := range claims {
			if key == `username` {
				username = val.(string)
			}
			if key == `role` {
				role = val.(string)
			}
		}

		if username == `` || role == `` {
			c.Next()
		}

		// Установка данных в контекст
		c.Set(`username`, username)
		c.Set(`role`, role)

		c.Next()
	}
}
