package service_test

import (
	"testing"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/service"

	"github.com/stretchr/testify/suite"
)

// AuthServiceTestSuite defenies test suite for authentication service
type AuthServiceTestSuite struct {
	suite.Suite
	authService service.AuthService
}

// SetupTest initiates new authentication service
func (suite *AuthServiceTestSuite) SetupTest() {
	suite.authService = service.NewAuthService()
}

// TestAuthServiceTestSuite runs test suite
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) TestGenerateJWT() {
	// set config
	common.Config = &common.Configuration{
		AccessSecret:  "secret",
		RefreshSecret: "refresh_secret",
	}

	td, err := suite.authService.GenerateJWT("test", false)
	suite.Nil(err)
	suite.NotNil(td)
}
