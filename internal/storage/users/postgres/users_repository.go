package postgres

import (
	"database/sql"
	"errors"
	"os"

	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Prefix: `[database-user]`,
	Level:  log.DebugLevel,
})

type UserRepository struct {
	db *sqlx.DB
}

func CreateUserRepository(db *sqlx.DB) users.UserStorage {
	return &UserRepository{
		db: db,
	}
}

// CreateUser создает нового пользователя.
//
// Параметры:
// - user: указатель на users.UserCreate структуру.
//
// Возвращает:
// - error: ошибка при создании пользователя.
func (r *UserRepository) CreateUser(user *users.UserCreate) error {
	// Проверка на существование
	if u, err := r.GetUserByUsername(user.Username); u != nil {
		return errors.New(`user already exists`)
	} else if err != nil {
		return err
	}

	tokens := make([]string, 0)

	u := &users.User{
		ID:            uuid.NewString(),
		Username:      user.Username,
		Role:          user.Role,
		RefreshTokens: tokens,
	}

	// Хеширование пароля
	hash := hashString(user.Password)
	u.Password = hash

	// Создание
	_, err := r.db.NamedExec(`INSERT INTO users (id, username, password, role) VALUES (:id, :username, :password, :role)`, u)
	return err
}

// GetUserByUsername ищет пользователя по нику.
//
// Параметры:
// - username: ник пользователя.
//
// Возвращает:
// - *users.User: указатель на users.User структуру. Возвращает nil если пользователь не найден.
// - error: ошибка во время получения пользователя.
func (r *UserRepository) GetUserByUsername(username string) (*users.User, error) {
	user := users.User{}

	// Получаем пользователя
	err := r.db.Get(&user, `SELECT * FROM users WHERE username = $1`, username)

	// Если пользователь не найден
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	// Если получили ошибку
	if err != nil {
		logger.Error(`failed to get user`, `err`, err)
		return nil, err
	}

	return &users.User{
		ID:            user.ID,
		Username:      user.Username,
		Role:          user.Role,
		Password:      user.Password,
		RefreshTokens: user.RefreshTokens,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

// UpdateUser обновляет пользователя.
//
// Параметры:
// - user: указатель на users.User.
//
// Возвращает:
// - error: ошибка во время обновления.
func (r *UserRepository) UpdateUser(user *users.User) error {
	_, err := r.db.NamedExec(`UPDATE users SET username = :username, password = :password, role = :role, refresh_tokens = :refresh_tokens WHERE id = :id`, user)
	return err
}

// DeleteUser удаляет пользователя.
//
// Параметры:
// - id: ID пользователя.
//
// Возвращает:
// - error: ошибка во время удаления.
func (r *UserRepository) DeleteUser(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}

// hashString генерирует хэш на основе вводной строки.
//
// Параметры:
// - s: вводная строка.
//
// Возвращает:
// - string: хэшированная строка.
func hashString(s string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(bytes)
}
