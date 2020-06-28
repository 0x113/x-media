package service

import (
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
		return err
	}

	// generate password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
	if err != nil {
		log.Errorf("Couldn't generate password for user: %v", err)
		return err
	}
	u.Password = string(hash)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	if err := s.repo.Create(u); err != nil {
		log.Errorf("Couldn't create user: %v", err)
		return err
	}

	log.Infof("Successfully create new user [username=%s]", u.Username)
	return nil
}

// GetUser calls the database layer to get user by username from the database
func (s *userService) GetUser(username string) (*models.User, error) {
	user, err := s.repo.Get(username)
	if err != nil {
		log.Errorf("Couldn't get user [username=%s]: %v", username, err)
		return nil, err
	}

	log.Infof("Successfully found user [username=%s]", username)
	return user, nil
}
