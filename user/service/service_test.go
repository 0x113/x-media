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
