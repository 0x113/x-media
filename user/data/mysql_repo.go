package data

import (
	"database/sql"
	"fmt"

	"github.com/0x113/x-media/user/databases"
	"github.com/0x113/x-media/user/models"

	"github.com/go-sql-driver/mysql"
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
	defer stmt.Close()

	if _, err := stmt.Exec(u.Username, u.Password, u.IsAdmin, u.CreatedAt, u.UpdatedAt); err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return fmt.Errorf("User %s already exists in the database", u.Username)
		}
		return err
	}

	return nil
}

// Get user by username from the database
func (r *userRepository) Get(username string) (*models.User, error) {
	query := "SELECT * FROM user WHERE username = ?"

	var user models.User
	if err := databases.Database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("There is no user %s in the database", username)
		}
		return nil, err
	}
	return &user, nil
}
