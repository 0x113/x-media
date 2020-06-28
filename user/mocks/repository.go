package mocks

import (
	"fmt"

	"github.com/0x113/x-media/user/models"
)

// MockUserRepository representes in-memory user repository
type MockUserRepository struct {
	users map[string]*models.User
}

// NewMockUserRepository creates new instance of MockUserRepository
func NewMockUserRepository() *MockUserRepository {
	var users = map[string]*models.User{}
	users["JohnDoe"] = &models.User{
		ID:       420,
		Username: "JohnDoe",
		Password: "$2a$11$zBkkaUb7woE6Y4oGeqrzYeNlmZ.e/3IbNCfxEYtASk.YHJFYGpfzK",
		IsAdmin:  false,
	}
	return &MockUserRepository{users}
}

// Create new user in memory
func (r *MockUserRepository) Create(u *models.User) error {
	if _, ok := r.users[u.Username]; ok {
		return fmt.Errorf("Couldn't create new user: user [username=%s] already exists", u.Username)
	}

	r.users[u.Username] = u
	return nil
}

// Get user by username from memory
func (r *MockUserRepository) Get(username string) (*models.User, error) {
	if user, ok := r.users[username]; ok {
		return user, nil
	}

	return nil, fmt.Errorf("User with username: %s; doesn't exist", username)
}
