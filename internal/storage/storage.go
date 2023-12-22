package storage

import (
	"os"

	"github.com/Jourloy/Go-Budget-Service/internal/storage/budgets"
	budgetsRep "github.com/Jourloy/Go-Budget-Service/internal/storage/budgets/postgres"
	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	usersRep "github.com/Jourloy/Go-Budget-Service/internal/storage/users/postgres"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	User   users.UserStorage
	Budget budgets.BudgetStorage
	Spend  budgets.SpendStorage
}

var (
	// Logger for the storage package
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[database]`,
		Level:  log.DebugLevel,
	})

	DatabaseDNS string
)

// parseENV parses the environment variables.
func parseENV() {
	if env, exist := os.LookupEnv(`DATABASE_DSN`); exist {
		DatabaseDNS = env
	}
}

// CreateStorage initializes and returns a new instance of the Storage struct.
func CreateStorage() *Storage {
	// Initialization
	parseENV()

	// Create connection
	db, err := sqlx.Connect(`postgres`, DatabaseDNS)
	if err != nil {
		logger.Fatal(`Failed to connect to database`, `err`, err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Fatal(`Failed to create driver`, `err`, err)
	}

	m, err := migrate.NewWithDatabaseInstance(`file://migrations`, `postgres`, driver)
	if err != nil {
		logger.Fatal(`Failed to create migration`, `err`, err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(`Failed to migrate database`, `err`, err)
	}

	// Create repositories
	usersStorage := usersRep.CreateUserRepository(db)
	budgetsStorage := budgetsRep.CreateBudgetRepository(db)
	spendsStorage := budgetsRep.CreateSpendRepository(db)

	return &Storage{
		User:   usersStorage,
		Budget: budgetsStorage,
		Spend:  spendsStorage,
	}
}
