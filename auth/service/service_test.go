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

func (suite *AuthServiceTestSuite) TestExtractTokenMetadata() {
	// set config
	common.Config = &common.Configuration{
		AccessSecret:  "secret",
		RefreshSecret: "refresh_secret",
	}

	testCases := []struct {
		name          string
		username      string
		token         string
		generateToken bool
		isAdmin       bool
		wantErr       bool
	}{
		{
			name:          "Success",
			username:      "JohnDoe",
			token:         "",
			generateToken: true,
			isAdmin:       false,
			wantErr:       false,
		},
		{
			name:          "Admin user - success",
			username:      "JohnDoe",
			token:         "",
			generateToken: true,
			isAdmin:       true,
			wantErr:       false,
		},
		{
			name:          "Wrong signing method",
			username:      "",
			token:         "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA",
			generateToken: false,
			isAdmin:       false,
			wantErr:       true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			// generate token
			if tt.generateToken {
				token, err := suite.authService.GenerateJWT(tt.username, tt.isAdmin)
				tt.token = token.AccessToken
				suite.Nil(err)
			}
			// extract data from token
			accessDetails, err := suite.authService.ExtractTokenMetadata(tt.token)
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(accessDetails)
			} else {
				suite.Nil(err)
				suite.Equal(tt.username, accessDetails.Username)
				suite.Equal(tt.isAdmin, accessDetails.IsAdmin)
			}
		})
	}
}
