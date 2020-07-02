package handler

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0x113/x-media/user/mocks"
	"github.com/0x113/x-media/user/service"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

// UserHandlerTestSuite defines the test suite
type UserHandlerTestSuite struct {
	suite.Suite
	userRepo    *mocks.MockUserRepository
	userService service.UserService
}

// SetupTest inititates mocked database and new user service
func (suite *UserHandlerTestSuite) SetupTest() {
	suite.userRepo = mocks.NewMockUserRepository()
	suite.userService = service.NewUserService(suite.userRepo)
	logrus.SetOutput(ioutil.Discard)
}

// TestUserHandlerTestSuite runs the test suite
func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (suite *UserHandlerTestSuite) TestCreateUser() {
	e := echo.New()
	h := &userHandler{suite.userService}
	testCases := []struct {
		name               string
		json               string
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Success",
			json:               `{"username": "test", "password": "strongpassword"}`,
			expectedStatusCode: 201,
			wantErr:            false,
		},
		{
			name:               "Decoding error",
			json:               ``,
			expectedStatusCode: 400,
			wantErr:            true,
		},
		{
			name:               "Existing user",
			json:               `{"username": "JohnDoe", "password": "strong"}`,
			expectedStatusCode: 500,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/create", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h.CreateUser(c)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedStatusCode, rec.Code)
		})
	}
}

func (suite *UserHandlerTestSuite) TestValidateUser() {
	e := echo.New()
	h := &userHandler{suite.userService}
	testCases := []struct {
		name               string
		json               string
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "Success",
			json:               `{"username": "JohnDoe", "password": "test1231"}`,
			expectedStatusCode: 200,
			wantErr:            false,
		},
		{
			name:               "Bad request",
			json:               `{"username: "test123"}`,
			expectedStatusCode: 400,
			wantErr:            true,
		},
		{
			name:               "Wrong password",
			json:               `{"username": "JohnDoe", "password": "strong"}`,
			expectedStatusCode: 500,
			wantErr:            true,
		},
	}

	for _, tt := range testCases {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/validate", strings.NewReader(tt.json))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h.ValidateUser(c)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
			suite.Equal(tt.expectedStatusCode, rec.Code)
		})
	}
}
