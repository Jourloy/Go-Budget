package memory

import (
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
)

type UserMemory struct {
	sync.Mutex
	Users map[string]users.User
}

// Это хранилище используется исключительно для тестов логики сервисов.
// Определенные методы не реализуют должную систему проверки.
//
// В качестве кэш хранилища лучше реализовать Redis
func CreateUserMemory() users.UserStorage {
	users := make(map[string]users.User)

	return &UserMemory{
		Users: users,
	}
}

// CreateUser создает нового пользователя.
//
// Параметры:
// - user: указатель на users.UserCreate структуру.
//
// Возвращает:
// - error: всегда nil.
func (m *UserMemory) CreateUser(user *users.UserCreate) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	tokens := make([]string, 0)
	id := uuid.NewString()
	u := &users.User{
		ID:            id,
		Username:      user.Username,
		Role:          user.Role,
		RefreshTokens: tokens,
	}

	// Хеширование пароля
	hash := hashString(user.Password)
	u.Password = hash

	m.Users[id] = *u
	return nil
}

// GetUserByUsername ищет пользователя по нику.
//
// Параметры:
// - username: ник пользователя.
//
// Возвращает:
// - *users.User: указатель на users.User структуру. Возвращает nil если пользователь не найден.
// - error: всегда nil.
func (m *UserMemory) GetUserByUsername(username string) (*users.User, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for _, u := range m.Users {
		if u.Username == username {
			return &u, nil
		}
	}

	return nil, nil
}

// UpdateUser обновляет пользователя.
//
// Параметры:
// - user: указатель на users.User.
//
// Возвращает:
// - error: всегда nil.
func (m *UserMemory) UpdateUser(user *users.User) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.Users[user.ID] = *user

	return nil
}

// DeleteUser удаляет пользователя.
//
// Параметры:
// - id: ID пользователя.
//
// Возвращает:
// - error: всегда nil.
func (m *UserMemory) DeleteUser(id string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	delete(m.Users, id)

	return nil
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
