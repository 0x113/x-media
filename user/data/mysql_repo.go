package data

import (
	"github.com/0x113/x-media/user/databases"
	"github.com/0x113/x-media/user/models"

	_ "github.com/go-sql-driver/mysql"
)

// userRepository manages the use CRUD
type userRepository struct{}

// NewMySQLUserRepository returns a new instance of UserRepository
func NewMySQLUserRepository() UserRepository {
	return &userRepository{}
}

// Create new user in the database
func (r *userRepository) Create(u *models.User) error {
	query := "INSERT INTO user (username, password, is_admin, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"

	stmt, err := databases.Database.DB.Prepare(query)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(u.Username, u.Password, u.IsAdmin, u.CreatedAt, u.UpdatedAt); err != nil {
		return err
	}

	return nil
}
