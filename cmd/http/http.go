package main

import (
	"os"

	"github.com/Jourloy/Go-Budget-Service/internal"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[cmd]`,
	})
)

func main() {
	// Load .env
	if err := godotenv.Load(`.env`); err != nil {
		logger.Fatal(`Error loading .env file`)
	}

	// Start server
	internal.StartServer()
}
