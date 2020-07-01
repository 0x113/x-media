package data

import "github.com/0x113/x-media/user/models"

// UserRepository contains all methods for operation on User model
type UserRepository interface {
	Create(u *models.User) error
	Get(username string) (*models.User, error)
}
