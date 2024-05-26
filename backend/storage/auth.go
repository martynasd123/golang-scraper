package storage

import (
	"errors"
	"sync"
	"time"
)

type User struct {
	Id                     *int
	Username               string
	Password               string
	DeviceIdentifier       *string
	RefreshToken           *string
	RefreshTokenValidUntil *time.Time
}

type AuthDao interface {
	CreateUser(*User) error
	GetUser(username string) (*User, error)
	UpdateUser(user *User) error
}

type InMemoryAuthDao struct {
	mu    sync.Mutex
	users map[string]User
}

func CreateAuthInMemoryDao() *InMemoryAuthDao {
	return &InMemoryAuthDao{
		users: make(map[string]User),
	}
}

func (AuthStorage *InMemoryAuthDao) CreateUser(user *User) error {
	AuthStorage.mu.Lock()
	defer AuthStorage.mu.Unlock()
	if _, exists := AuthStorage.users[user.Username]; exists {
		return errors.New("Username already exists")
	}
	AuthStorage.users[user.Username] = *user
	return nil
}

func (AuthStorage *InMemoryAuthDao) GetUser(username string) (*User, error) {
	AuthStorage.mu.Lock()
	defer AuthStorage.mu.Unlock()

	user, exists := AuthStorage.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (AuthStorage *InMemoryAuthDao) UpdateUser(user *User) error {
	AuthStorage.mu.Lock()
	defer AuthStorage.mu.Unlock()

	existingUser, exists := AuthStorage.users[user.Username]
	if !exists {
		return errors.New("user not found")
	}

	existingUser.Username = user.Username
	existingUser.Password = user.Password
	existingUser.DeviceIdentifier = user.DeviceIdentifier
	existingUser.RefreshToken = user.RefreshToken
	existingUser.RefreshTokenValidUntil = user.RefreshTokenValidUntil

	AuthStorage.users[user.Username] = existingUser

	return nil
}
