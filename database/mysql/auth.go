package mysql

import (
	"database/sql"

	"github.com/0x113/x-media/auth"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type authRepository struct {
	db *sql.DB
}

// NewMySQLAuthRepository creates new authRepository
func NewMySQLAuthRepository(db *sql.DB) auth.AuthRepository {
	return &authRepository{
		db,
	}
}

func (r *authRepository) Create(user *auth.User) error {
	query := "INSERT INTO user (username, password) VALUE (?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(user.Username, hashedPassword)
	if err != nil {
		return err
	}

	newID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	log.Infof("Created user with id %d", newID)
	return nil
}
