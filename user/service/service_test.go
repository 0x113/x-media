package service_test

import (
	"io/ioutil"
	"testing"

	"github.com/0x113/x-media/user/mocks"
	"github.com/0x113/x-media/user/models"
	"github.com/0x113/x-media/user/service"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// UserServiceTestSuite defines test suite for UserService
type UserServiceTestSuite struct {
	suite.Suite
	userRepo    *mocks.MockUserRepository
	userService service.UserService
}

// SetupTest initiates mocked database and new user service
func (suite *UserServiceTestSuite) SetupTest() {
	suite.userRepo = mocks.NewMockUserRepository()
	suite.userService = service.NewUserService(suite.userRepo)
	logrus.SetOutput(ioutil.Discard)
}

// TestUserServiceTestSuite runs test suite
func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestCreateUser() {
	testCases := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "Success",
			user: &models.User{
				Username: "test123",
				Password: "strongpassword",
			},
			wantErr: false,
		},
		{
			name: "Duplicate username",
			user: &models.User{
				Username: "test123",
				Password: "strongpassword",
			},
			wantErr: true,
		},
		{
			name: "Empty username",
			user: &models.User{
				Username: "",
				Password: "strongpassword",
			},
			wantErr: true,
		},
		{
			name: "Provided IsAdmin field",
			user: &models.User{
				Username: "adminUser",
				Password: "adminPass",
				IsAdmin:  true,
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			err := suite.userService.CreateUser(tt.user)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *UserServiceTestSuite) TestGetUser() {
	testCases := []struct {
		name         string
		username     string
		expectedUser *models.User
		wantErr      bool
	}{
		{
			name:     "Success",
			username: "JohnDoe", // this user exists in a mocked database
			expectedUser: &models.User{
				ID:       420,
				Username: "JohnDoe",
				Password: "$2a$11$zBkkaUb7woE6Y4oGeqrzYeNlmZ.e/3IbNCfxEYtASk.YHJFYGpfzK",
				IsAdmin:  false,
			},
			wantErr: false,
		},
		{
			name:         "Non-existent user",
			username:     "NotJohnDoe",
			expectedUser: nil,
			wantErr:      true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			user, err := suite.userService.GetUser(tt.username)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}

			suite.Equal(tt.expectedUser, user)
		})
	}
}

func (suite *UserServiceTestSuite) TestValidateUser() {
	testCases := []struct {
		name           string
		creds          *models.Credentials
		expectedClaims *models.TokenClaims
		wantErr        bool
	}{
		{
			name: "Success",
			creds: &models.Credentials{
				Username: "JohnDoe",
				Password: "test1231",
			},
			expectedClaims: &models.TokenClaims{
				Username: "JohnDoe",
				IsAdmin:  false,
			},
			wantErr: false,
		},
		{
			name: "Wrong password",
			creds: &models.Credentials{
				Username: "JohnDoe",
				Password: "ThatNotJohnDoesPassword",
			},
			expectedClaims: nil,
			wantErr:        true,
		},
		{
			name: "Non-existent user",
			creds: &models.Credentials{
				Username: "JanKowalski",
				Password: "SuperCoolPassword",
			},
			expectedClaims: nil,
			wantErr:        true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			claims, err := suite.userService.ValidateUser(tt.creds)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedClaims, claims)
		})
	}
}
