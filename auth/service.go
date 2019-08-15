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
	if err := s.repo.Create(user); err != nil {
		log.Errorf("Unable to create user [username=%s]: %v", user.Username, err)
		return err
	}
	log.Infof("Successfully created user [id=%d, username=%s]", user.ID, user.Username)
	return nil
}

func (s *authService) LoginUser(username, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		log.Errorf("Unable to get user [username=%s] from databse: %v", username, err)
		return "", err
	}

	if user == nil {
		log.Errorf("Invalid username: %s", username)
		return "", errors.New("Invalid username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Errorf("Entered password for user [id=%d, username=%s] is incorrect", user.ID, username)
		return "", err
	}

	token, err := s.getToken(user)
	if err != nil {
		log.Errorf("Unable to generate token for user [id=%d, username=%s]: %v", user.ID, user.Username, err)
		return "", err
	}

	log.Infof("Generated token for user [id=%d, username=%s]", user.ID, username)

	return token, nil
}

func (s *authService) getToken(user *User) (string, error) {

	/* Create the token */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"type":     "user",
		"exp":      time.Now().Add(168 * time.Hour).Unix(), //TODO: Change for 5 minutes this is just for example
	})

	/* Sign the token with key */
	tokenStr, err := token.SignedString([]byte(env.EnvString("JWT_KEY")))
	if err != nil {
		log.Errorf("Unable to sign token for user [id=%d, username=%s]; %v", user.ID, user.Username, err)
		return "", err
	}
	return tokenStr, nil
}
