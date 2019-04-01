package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the interface that provides authentication methods.
type AuthService interface {
	CreateUser(user *User) error
	LoginUser(username, password string) (string, error)
}

type authService struct {
	repo   AuthRepository
	jwtKey []byte
}

// NewAuthService returns a new instance of authentication service.
func NewAuthService(repo AuthRepository, jwtKey string) AuthService {
	return &authService{
		repo,
		[]byte(jwtKey),
	}
}

func (s *authService) CreateUser(user *User) error {
	return s.repo.Create(user)
}

func (s *authService) LoginUser(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("Invalid username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	token, err := s.getToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) getToken(user *User) (string, error) {

	/* Create the token */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"type":     "user",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	})

	/* Sign the token with key */
	return token.SignedString(s.jwtKey)

}
