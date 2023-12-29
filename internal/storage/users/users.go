package users

import "github.com/lib/pq"

type User struct {
	ID            string         `db:"id"`
	Username      string         `db:"username"`
	Password      string         `db:"password"`
	Role          string         `db:"role"`
	RefreshTokens pq.StringArray `db:"refresh_tokens"`
	CreatedAt     string         `db:"created_at"`
	UpdatedAt     string         `db:"updated_at"`
}

type UserCreate struct {
	Username string
	Password string
	Role     string
}

type UserStorage interface {
	CreateUser(user *UserCreate) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
}
