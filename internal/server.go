package internal

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/Jourloy/Go-Budget-Service/internal/handlers"
	"github.com/Jourloy/Go-Budget-Service/internal/middlewares"
	"github.com/Jourloy/Go-Budget-Service/internal/storage"
	"github.com/charmbracelet/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[server]`,
		Level:  log.DebugLevel,
	})

	Port = flag.Int(`p`, 10001, `Port of the server`)
)

// parseENV parses the environment variables.
func parseENV() {
	if env, exist := os.LookupEnv(`PORT`); exist {
		if p, err := strconv.Atoi(env); err == nil {
			Port = &p
		}
	}
}

// StartServer initializes the server and starts listening for incoming requests.
func StartServer() {
	// Initialization
	flag.Parse()
	parseENV()

	// Create router
	r := gin.New()

	// Middlewares
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{`http://localhost:10000`, `http://192.168.0.222:10000`},
		AllowMethods:     []string{`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS`},
		AllowHeaders:     []string{`Origin`, `Content-Length`, `Content-Type`, `Authorization`},
		AllowCredentials: true,
	}))
	r.Use(gin.Recovery())

	// Custom middlewares
	r.Use(loggerMiddleware())
	r.Use(middlewares.JwtMiddleware())

	// Groups
	app := r.Group(`/`)
	auth := r.Group(`/auth/`)
	budget := r.Group(`/budget/`)

	// Storage
	s := storage.CreateStorage()

	// Handlers
	handlers.RegisterAppHandler(app)
	handlers.RegisterAuthHandler(auth, s)
	handlers.RegisterBudgetHandler(budget, s)

	// Start server
	logger.Info(`Server started on port ` + strconv.Itoa(*Port))
	if err := r.Run("0.0.0.0:" + strconv.Itoa(*Port)); err != nil {
		logger.Fatal(err)
	}
}

// loggerMiddleware is a middleware function that logs information about incoming requests.
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.Info(
			`Response`,
			`method`, method,
			`path`, path,
			`status`, status,
			`latency`, latency,
		)
	}
}
