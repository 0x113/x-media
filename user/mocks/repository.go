package mocks

import (
	"fmt"

	"github.com/0x113/x-media/user/models"
)

// MockUserRepository representes in-memory user repository
type MockUserRepository struct {
	users map[int]*models.User
}

// NewMockUserRepository creates new instance of MockUserRepository
func NewMockUserRepository() *MockUserRepository {
	var users = map[int]*models.User{}
	users[420] = &models.User{
		ID:       420,
		Username: "JohnDoe",
		Password: "ThisShouldSuperStrongAndSecret",
		IsAdmin:  false,
	}
	return &MockUserRepository{users}
}

// Create new user in memory
func (r *MockUserRepository) Create(u *models.User) error {
	if _, ok := r.users[u.ID]; ok {
		return fmt.Errorf("Couldn't create new user: user [id=%d] already exists", u.ID)
	}

	r.users[u.ID] = u
	return nil
}
