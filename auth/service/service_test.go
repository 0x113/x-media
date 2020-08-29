package service_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/mocks"
	"github.com/0x113/x-media/auth/models"
	"github.com/0x113/x-media/auth/service"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// AuthServiceTestSuite defenies test suite for authentication service
type AuthServiceTestSuite struct {
	suite.Suite
	httpClient  *mocks.MockClient
	authRepo    *mocks.MockAuthRepository
	authService service.AuthService
}

// SetupTest initiates new authentication service
func (suite *AuthServiceTestSuite) SetupTest() {
	// set config
	common.Config = &common.Configuration{
		AccessSecret:  "secret",
		RefreshSecret: "refresh_secret",
	}
	logrus.SetOutput(ioutil.Discard)

	suite.httpClient = &mocks.MockClient{}
	suite.authRepo = mocks.NewMockAuthRepository()
}

// TestAuthServiceTestSuite runs test suite
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) TestLogin() {
	testCases := []struct {
		name    string
		creds   *models.Credentials
		wantErr bool
		DoFunc  func(req *http.Request) (*http.Response, error)
	}{
		{
			name: "Success",
			creds: &models.Credentials{
				Username: "JohnDoe",
				Password: "test1231",
			},
			wantErr: false,
			DoFunc: func(req *http.Request) (*http.Response, error) {
				jsonStr := `{"username": "JohnDoe", "is_admin": false}`
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonStr))),
				}, nil
			},
		},
		{
			name: "Error when calling user service",
			creds: &models.Credentials{
				Username: "JohnDoe",
				Password: "test1231",
			},
			wantErr: true,
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{}, errors.New("Service is down")
			},
		},
		{
			name: "Wrong status code",
			creds: &models.Credentials{
				Username: "JohnDoe",
				Password: "incorrectPassword",
			},
			wantErr: true,
			DoFunc: func(req *http.Request) (*http.Response, error) {
				jsonStr := `{"code": 500, "message": "Invalid user credentials"}`
				return &http.Response{
					StatusCode: 500,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonStr))),
				}, nil
			},
		},
		{
			name: "Invalid response; unable to generate token",
			creds: &models.Credentials{
				Username: "JohnDoe",
				Password: "incorrectPassword",
			},
			wantErr: true,
			DoFunc: func(req *http.Request) (*http.Response, error) {
				jsonStr := `{"username": "JohnDoe"}`
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonStr))),
				}, nil
			},
		},
	}

	for _, tt := range testCases {
		// set up httpClient and auth service for the subtest
		suite.httpClient = &mocks.MockClient{tt.DoFunc}
		suite.authService = service.NewAuthService(suite.httpClient, suite.authRepo)

		suite.Run(tt.name, func() {
			token, err := suite.authService.Login(tt.creds)
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(token)
			} else {
				suite.Nil(err)
				suite.NotNil(token)
			}
		})
	}

}

func (suite *AuthServiceTestSuite) TestGenerateJWT() {
	suite.authService = service.NewAuthService(suite.httpClient, suite.authRepo)
	testCases := []struct {
		name    string
		details *models.AccessDetails
		wantErr bool
	}{
		{
			name: "Success",
			details: &models.AccessDetails{
				Username: "JohnDoe",
				IsAdmin:  new(bool),
			},
			wantErr: false,
		},
		{
			name: "Validation error",
			details: &models.AccessDetails{
				Username: "JohnDoe",
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			td, err := suite.authService.GenerateJWT(tt.details)
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(td)
			} else {
				suite.Nil(err)
				suite.NotNil(td)
			}
		})
	}
}

func (suite *AuthServiceTestSuite) TestExtractTokenMetadata() {
	suite.authService = service.NewAuthService(suite.httpClient, suite.authRepo)

	testCases := []struct {
		name          string
		details       *models.AccessDetails
		token         string
		generateToken bool
		wantErr       bool
	}{
		{
			name: "Success",
			details: &models.AccessDetails{
				Username: "JohnDoe",
				IsAdmin:  new(bool),
			},
			token:         "",
			generateToken: true,
			wantErr:       false,
		},
		{
			name: "Admin user - success",
			details: &models.AccessDetails{
				Username: "JohnDoe",
				IsAdmin:  &[]bool{true}[0], // should *bool to true; quite messy but need to be pointer for validatiote
			},
			token:         "",
			generateToken: true,
			wantErr:       false,
		},
		{
			name: "Wrong signing method",
			details: &models.AccessDetails{
				Username: "JohnDoe",
				IsAdmin:  new(bool),
			},
			token:         "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA",
			generateToken: false,
			wantErr:       true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			// generate token
			if tt.generateToken {
				token, err := suite.authService.GenerateJWT(tt.details)
				suite.Nil(err)
				tt.token = token.AccessToken
			}
			// extract data from token
			accessDetails, err := suite.authService.ExtractTokenMetadata(tt.token, common.Config.AccessSecret) // TODO: test laso refresh secret
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(accessDetails)
			} else {
				suite.Nil(err)
				suite.Equal(tt.details.Username, accessDetails.Username)
				suite.Equal(tt.details.IsAdmin, accessDetails.IsAdmin)
			}
		})
	}
}

func (suite *AuthServiceTestSuite) TestRefresh() {
	testCases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Success",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEZXRhaWxzIjp7InVzZXJuYW1lIjoiSm9obkRvZSIsImlzX2FkbWluIjpmYWxzZX0sIlV1aWQiOiJiNjZhNzIxOS1mMDdmLTQ5Y2YtODE2My0xODlkYTJmNWM4Y2MiLCJleHAiOjE1OTkzMTc3ODZ9.Yrb3n3BKP3Ol6MxVQjJAdJNCFPzOa627BBuKmlfVFfk",
			wantErr: false,
		},
		{
			name:    "No token to delete, no user with UUID from the token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEZXRhaWxzIjp7InVzZXJuYW1lIjoiSm9obkRvZSIsImlzX2FkbWluIjpmYWxzZX0sIlV1aWQiOiIyMGZhYTY3ZS1iZGEyLTQ5NTUtODUyYi0yMjI2N2EzM2VhZGQiLCJleHAiOjE1OTQ1NzcyOTJ9.W8szMrFZS0jqEi3hnX_V4GsITsVQ5bp5TN3zpQDcIes",
			wantErr: true,
		},
		{
			name:    "Wrong method",
			token:   "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA",
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			token, err := suite.authService.Refresh(tt.token)
			if tt.wantErr {
				suite.NotNil(err)
				suite.Nil(token)
			} else {
				suite.Nil(err)
				suite.NotNil(token)
			}
		})
	}
}
