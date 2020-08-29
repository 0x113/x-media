package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/mocks"
	"github.com/0x113/x-media/auth/models"
	"github.com/0x113/x-media/auth/service"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
)

// AuthHandlerTestSuite defines the test suite
type AuthHandlerTestSuite struct {
	suite.Suite
	httpClient  *mocks.MockClient
	authRepo    *mocks.MockAuthRepository
	authService service.AuthService
}

// SetupTest initiates new authentication service and sets the config
func (suite *AuthHandlerTestSuite) SetupTest() {
	// set config
	common.Config = &common.Configuration{
		AccessSecret:  "secret",
		RefreshSecret: "refresh_secret",
	}
	logrus.SetOutput(ioutil.Discard)
	suite.authRepo = mocks.NewMockAuthRepository()
}

// TestAuthHandlerTestSuite runs the test suite
func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (suite *AuthHandlerTestSuite) TestGenerateToken() {
	e := echo.New()

	testCases := []struct {
		name               string
		json               string
		expectedStatusCode int
		wantErr            bool
		DoFunc             func(req *http.Request) (*http.Response, error)
	}{
		{
			name:               "Success",
			json:               `{"username": "JohnDoe", "password": "test1231"}`,
			expectedStatusCode: 200,
			wantErr:            false,
			DoFunc: func(req *http.Request) (*http.Response, error) {
				jsonStr := `{"username": "JohnDoe", "is_admin": false}`
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonStr))),
				}, nil
			},
		},
		{
			name:               "Invalid json",
			json:               `{"username": "JohnDoe", "is_cool": "no"}`,
			expectedStatusCode: 500,
			wantErr:            true,
			DoFunc: func(req *http.Request) (*http.Response, error) {
				jsonStr := `{"code": 500, "message": "Invalid user credentials"}`
				return &http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonStr))),
				}, nil
			},
		},
		{
			name:               "Binding error",
			json:               ``,
			expectedStatusCode: 400,
			wantErr:            true,
			DoFunc:             nil,
		},
	}

	for _, tt := range testCases {
		// set up httpClient, auth service and handler
		suite.httpClient = &mocks.MockClient{tt.DoFunc}
		suite.authService = service.NewAuthService(suite.httpClient, suite.authRepo)
		h := authHandler{suite.authService}

		// run the subtest
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
	suite.httpClient = &mocks.MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(nil),
			}, nil
		},
	}
	suite.authService = service.NewAuthService(suite.httpClient, suite.authRepo)

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
			expectedStatusCode: 400,
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

func (suite *AuthHandlerTestSuite) TestRefreshToken() {
	suite.authService = service.NewAuthService(suite.httpClient, suite.authRepo)
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
			json:               `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEZXRhaWxzIjp7InVzZXJuYW1lIjoiSm9obkRvZSIsImlzX2FkbWluIjpmYWxzZX0sIlV1aWQiOiJiNjZhNzIxOS1mMDdmLTQ5Y2YtODE2My0xODlkYTJmNWM4Y2MiLCJleHAiOjE1OTkzMTc3ODZ9.Yrb3n3BKP3Ol6MxVQjJAdJNCFPzOa627BBuKmlfVFfk"}`,
			expectedStatusCode: 200,
			wantErr:            false,
		},
		{
			name:               "Empty request body",
			json:               ``,
			expectedStatusCode: 400,
			wantErr:            true,
		},
		{
			name:               "Invalid token type",
			json:               `{"token": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA"}`,
			expectedStatusCode: 400,
			wantErr:            true,
		},
		{
			name:               "Not token",
			json:               `{"token": 1}`,
			expectedStatusCode: 500,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/token/refresh", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.RefreshToken(c)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}

}
