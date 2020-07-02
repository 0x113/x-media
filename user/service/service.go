package service

import (
	"fmt"
	"time"

	"github.com/0x113/x-media/user/data"
	"github.com/0x113/x-media/user/models"

	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// UserService describes user service
type UserService interface {
	CreateUser(u *models.User) error
	ValidateUser(creds *models.Credentials) (*models.TokenClaims, error)
	GetUser(username string) (*models.User, error)
}

type userService struct {
	repo data.UserRepository
}

// NewUserService creates new instance of UserService
func NewUserService(repo data.UserRepository) UserService {
	return &userService{repo}
}

// CreateUser calls the database layer to create new user in the database
func (s *userService) CreateUser(u *models.User) error {
	// validate user
	validation := validator.New()
	if err := validation.Struct(u); err != nil {
		log.Errorf("Couldn't validate user: %v", err)
		return fmt.Errorf("Couldn't validate provided user data. Only two fields must be provided: username and password. Username must be at least 2 characters long and max 32 characters long. Password should be at least 8 characters long.") // NOTE: quite messy, probably should be better documented
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
	if err != nil {
		log.Errorf("Couldn't hash password for user: %v", err)
		return fmt.Errorf("Couldn't hash password")
	}
	u.Password = string(hash)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	if err := s.repo.Create(u); err != nil {
		log.Errorf("Couldn't create user: %v", err)
		return fmt.Errorf("Couldn't create new user: %v", err)
	}

	log.Infof("Successfully create new user [username=%s]", u.Username)
	return nil
}

// ValidateUser checks if provided credentials match with the data in the database
func (s *userService) ValidateUser(creds *models.Credentials) (*models.TokenClaims, error) {
	user, err := s.GetUser(creds.Username)
	if err != nil {
		return nil, err // no need to log, 'cause GetUser method does it
	}
	// compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Errorf("Wrong password for user [username=%s]: %v", creds.Username, err)
		return nil, fmt.Errorf("Invalid user credentials")
	}

	return &models.TokenClaims{user.Username, user.IsAdmin}, nil
}

// GetUser calls the database layer to get user by username from the database
func (s *userService) GetUser(username string) (*models.User, error) {
	user, err := s.repo.Get(username)
	if err != nil {
		log.Errorf("Couldn't get user [username=%s]: %v", username, err)
		return nil, fmt.Errorf("Couldn't get the user from the database: %v", err)
	}

	log.Infof("Successfully found user [username=%s]", username)
	return user, nil
}
