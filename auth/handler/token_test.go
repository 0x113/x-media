package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0x113/x-media/auth/common"
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
