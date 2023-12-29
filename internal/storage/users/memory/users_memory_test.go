package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Jourloy/Go-Budget-Service/internal/storage/users"
)

func TestUserMemory_CreateUser(t *testing.T) {
	type args struct {
		user *users.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				user: &users.UserCreate{
					Username: `test`,
					Password: `1q2w`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateUserMemory()

			if err := m.CreateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserMemory.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка имени
			u, e := m.GetUserByUsername(`test`)
			assert.True(t, e == nil)
			assert.Equal(t, tt.args.user.Username, u.Username)
			assert.NotEqual(t, tt.args.user.Password, u.Password) // Пароль должен быть хэширован
		})
	}
}

func TestUserMemory_UpdateUser(t *testing.T) {
	type args struct {
		user *users.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				user: &users.UserCreate{
					Username: `test`,
					Password: `1q2w`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateUserMemory()

			if err := m.CreateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserMemory.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка имени
			u, e := m.GetUserByUsername(`test`)
			assert.True(t, e == nil)

			u.Role = `admin`

			m.UpdateUser(u)

			// Проверка изменения
			newU, newE := m.GetUserByUsername(`test`)
			assert.True(t, newE == nil)
			assert.Equal(t, `admin`, newU.Role)
		})
	}
}

func TestUserMemory_DeleteUser(t *testing.T) {
	type args struct {
		user *users.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: `Positive`,
			args: args{
				user: &users.UserCreate{
					Username: `test`,
					Password: `1q2w`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// in-memory хранилище
			m := CreateUserMemory()

			if err := m.CreateUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserMemory.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Проверка имени
			u, e := m.GetUserByUsername(`test`)
			assert.True(t, e == nil)

			m.DeleteUser(u.ID)

			// Проверка удаления
			newU, newE := m.GetUserByUsername(`test`)
			assert.True(t, newE == nil)
			assert.Nil(t, newU)
		})
	}
}
