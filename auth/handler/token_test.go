package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/models"
	"github.com/0x113/x-media/auth/service"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
)

// AuthHandlerTestSuite defines the test suite
type AuthHandlerTestSuite struct {
	suite.Suite
	authService service.AuthService
}

// SetupTest initiates new authentication service and sets the config
func (suite *AuthHandlerTestSuite) SetupTest() {
	// set config
	common.Config = &common.Configuration{
		AccessSecret:  "secret",
		RefreshSecret: "refresh_secret",
	}
	suite.authService = service.NewAuthService()
}

// TestAuthHandlerTestSuite runs the test suite
func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (suite *AuthHandlerTestSuite) TestGenerateToken() {
	e := echo.New()
	h := authHandler{suite.authService}
	testCases := []struct {
		name               string
		json               string
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Success",
			json:               `{"username": "JohnDoe", "is_admin": false}`,
			expectedStatusCode: 200,
			wantErr:            false,
		},
		{
			name:               "Invalid json",
			json:               `{"username": "JohnDoe", "is_cool": "no"}`,
			expectedStatusCode: 500,
			wantErr:            true,
		},
		{
			name:               "Binding error",
			json:               ``,
			expectedStatusCode: 422,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token/generate", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.GenerateToken(c)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedStatusCode, rec.Code)
		})
	}
}

func (suite *AuthHandlerTestSuite) TestGetTokenMetadata() {
	e := echo.New()
	h := authHandler{suite.authService}
	testCases := []struct {
		name               string
		json               string
		expectedStatusCode int
		generateToken      bool
		wantErr            bool
	}{
		{
			name:               "Success",
			json:               ``,
			expectedStatusCode: 200,
			generateToken:      true,
			wantErr:            false,
		},
		{
			name:               "Binding error",
			json:               ``,
			expectedStatusCode: 422,
			wantErr:            true,
		},
		{
			name:               "Expired token",
			json:               `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjhmYWFlMmE2LTQ3MGEtNGQ2NC05ZGFmLTk2ZDFlYjJlMTVkMCIsImV4cCI6MTU5MjgzNjk5MywidXNlcl9pZCI6IjEyMyJ9.z88nbWhEamEjZbOBqz8cxYgrFWvbvvs2PJ1OjhStFu4"}`,
			expectedStatusCode: 500,
			wantErr:            true,
		},
		{
			name:               "Not a token",
			json:               `{"token": "it's definitely not a token"}`,
			expectedStatusCode: 500,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			// generate token, it's always valid so metadata needs to be validated
			if tt.generateToken {
				token, err := suite.authService.GenerateJWT(&models.AccessDetails{"JohnDoe", new(bool)})
				suite.Nil(err)
				tt.json = fmt.Sprintf(`{"token": "%s"}`, token.AccessToken)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token/validate", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.GetTokenMetadata(c)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedStatusCode, rec.Code)
		})
	}
}
