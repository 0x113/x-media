package mysql

import (
	"database/sql"
	"time"

	"github.com/0x113/x-media/auth"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type authRepository struct {
	db     *sql.DB
	jwtKey string
}

// NewMySQLAuthRepository creates new authRepository
func NewMySQLAuthRepository(db *sql.DB, jwtKey string) auth.AuthRepository {
	return &authRepository{
		db,
		jwtKey,
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

func (r *authRepository) GenerateJWT(user *auth.User) (string, error) {
	query := `SELECT password FROM user WHERE username = ?`

	var hashedPassword string
	err := r.db.QueryRow(query, user.Username).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(5 * time.Minute) // expiration time of the token
	// generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   user.Username,
		"expires_at": expirationTime.Unix(),
	})
	tokenString, err := token.SignedString([]byte(r.jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
