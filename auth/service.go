package auth

import (
	"errors"
	"time"

	"github.com/0x113/x-media/env"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the interface that provides authentication methods.
type AuthService interface {
	CreateUser(user *User) error
	LoginUser(username, password string) (string, error)
}

type authService struct {
	repo AuthRepository
}

// NewAuthService returns a new instance of authentication service.
func NewAuthService(repo AuthRepository) AuthService {
	return &authService{
		repo,
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
		log.WithFields(log.Fields{"user": username, "error": err.Error()}).Error("Unable to generate token")
		return "", err
	}

	log.Infof("Generated token for user %s", username)

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
	return token.SignedString([]byte(env.EnvString("jwt_key")))

}
