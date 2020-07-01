package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/httpclient"
	"github.com/0x113/x-media/auth/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)

// AuthService decribes authentication service
type AuthService interface {
	Login(creds *models.Credentials) (*models.TokenDetails, error)
	GenerateJWT(accessDetails *models.AccessDetails) (*models.TokenDetails, error)
	ExtractTokenMetadata(tokenString string) (*models.AccessDetails, error)
}

type authService struct {
	httpClient httpclient.HTTPClient
}

// NewAuthService creates new instance of authentication service
func NewAuthService(httpClient httpclient.HTTPClient) AuthService {
	return &authService{httpClient}
}

// Login calls the user service to check if provided credentials are correct
// and generates authentication token
func (s *authService) Login(creds *models.Credentials) (*models.TokenDetails, error) {
	// convert credentials to json
	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		log.Errorf("Couldn't convert credentials to json: %v", err)
		return nil, err

	}
	// call the user service to check is provided data is correct
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1/api/v1/user/validate", bytes.NewBuffer(jsonCreds))
	if err != nil {
		log.Errorf("Couldn't prepare request: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Host = "usersvc"

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Errorf("Couldn't to execute request: %v", err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Errorf("Expected status code: %d, got: %d", http.StatusOK, res.StatusCode)
		return nil, errors.New("Wrong status code")
	}
	defer res.Body.Close()

	// decode the response
	accessDetails := new(models.AccessDetails)
	if err := json.NewDecoder(res.Body).Decode(accessDetails); err != nil {
		log.Errorf("Couldn't decode the response: %v", err)
		return nil, err
	}

	token, err := s.GenerateJWT(accessDetails)
	if err != nil {
		return nil, err // no need to log, 'cause GenerateJWT does it
	}

	return token, nil
}

// GenerateJWT generates new token from provided data
func (s *authService) GenerateJWT(accessDetails *models.AccessDetails) (*models.TokenDetails, error) {
	// validation
	validate := validator.New()
	if err := validate.Struct(accessDetails); err != nil {
		return nil, err
	}

	td := &models.TokenDetails{}
	var err error
	// access token
	atClaims := &models.TokenClaims{
		accessDetails,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(common.Config.AccessSecret))
	if err != nil {
		return nil, err
	}

	// refresh token
	rtClaims := &models.TokenClaims{
		accessDetails,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(common.Config.RefreshSecret))

	return td, err
}

// ExtractTokenMetadata extracts data from provided JSON Web Token
func (s *authService) ExtractTokenMetadata(tokenString string) (*models.AccessDetails, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.Config.AccessSecret), nil
	})

	// make sure that token is not nil
	if token == nil {
		return nil, err

	}
	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return &models.AccessDetails{
			Username: claims.Details.Username,
			IsAdmin:  claims.Details.IsAdmin,
		}, nil
	}

	return nil, err
}
