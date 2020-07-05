package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/data"
	"github.com/0x113/x-media/auth/httpclient"
	"github.com/0x113/x-media/auth/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
)

// AuthService decribes authentication service
type AuthService interface {
	Login(creds *models.Credentials) (*models.TokenDetails, error)
	GenerateJWT(accessDetails *models.AccessDetails) (*models.TokenDetails, error)
	ExtractTokenMetadata(tokenString string) (*models.AccessDetails, error)
}

type authService struct {
	httpClient httpclient.HTTPClient
	repo       data.AuthRepository
}

// NewAuthService creates new instance of authentication service
func NewAuthService(httpClient httpclient.HTTPClient, repo data.AuthRepository) AuthService {
	return &authService{httpClient, repo}
}

// Login calls the user service to check if provided credentials are correct
// and generates authentication token
func (s *authService) Login(creds *models.Credentials) (*models.TokenDetails, error) {
	// convert credentials to json
	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		log.Errorf("Couldn't convert credentials to json: %v", err)
		return nil, fmt.Errorf("Couldn't convert credentials to the json")

	}
	// call the user service to check is provided data is correct
	req, err := http.NewRequest(http.MethodPost, "http://xmedia-user-svc:8002/api/v1/user/validate", bytes.NewBuffer(jsonCreds))
	if err != nil {
		log.Errorf("Couldn't prepare request: %v", err)
		return nil, fmt.Errorf("Couldn't prepare the request")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Errorf("Couldn't to execute request: %v", err)
		return nil, fmt.Errorf("Couldn't connect to the user service")
	}

	if res.StatusCode != http.StatusOK {
		log.Errorf("Expected status code: %d, got: %d", http.StatusOK, res.StatusCode)
		errMsg := new(models.Error)
		if err := json.NewDecoder(res.Body).Decode(errMsg); err != nil {
			return nil, fmt.Errorf("Couldn't decode the response from the user service")
		}
		return nil, fmt.Errorf(errMsg.Message) // return the error message from the user service
	}
	defer res.Body.Close()

	// decode the response
	accessDetails := new(models.AccessDetails)
	if err := json.NewDecoder(res.Body).Decode(accessDetails); err != nil {
		log.Errorf("Couldn't decode the response: %v", err)
		return nil, fmt.Errorf("Couldn't decode the response from the user service")
	}

	// generate the access and refresh token
	token, err := s.GenerateJWT(accessDetails)
	if err != nil {
		return nil, err // no need to log, 'cause GenerateJWT does it
	}

	// save tokens to the redis database
	if err := s.repo.Save(accessDetails.Username, token); err != nil {
		log.Errorf("Unable to save the access and refresh token to the database: %v", err)
		return nil, fmt.Errorf("Unable to save the access and refresh token to the database")
	}

	return token, nil
}

// GenerateJWT generates new token from provided data
func (s *authService) GenerateJWT(accessDetails *models.AccessDetails) (*models.TokenDetails, error) {
	// validation
	validate := validator.New()
	if err := validate.Struct(accessDetails); err != nil {
		log.Errorf("Couldn't validate access details for generating JWT: %v", err)
		return nil, fmt.Errorf("Provided credentials are invalid")
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
		log.Errorf("Couldn't sign the authentication token: %v", err)
		return nil, fmt.Errorf("Couldn't generate the authentication token")
	}
	td.AtExpires = time.Now().Add(15 * time.Minute).Unix()
	td.AccessUuid = uuid.NewV4().String()

	// refresh token
	rtClaims := &models.TokenClaims{
		accessDetails,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(common.Config.RefreshSecret))
	if err != nil {
		log.Errorf("Couldn't sign the refresh token: %v", err)
		return nil, fmt.Errorf("Couldn't generate the authentication and refresh token")
	}
	td.RtExpires = time.Now().Add(7 * 24 * time.Hour).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	return td, nil
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
		log.Errorf("Unable to extract the metadata from token; token is nil")
		return nil, fmt.Errorf("Token is nil")

	}
	// return claims if token is valid and token claims are same as models.TokenClaims
	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return &models.AccessDetails{
			Username: claims.Details.Username,
			IsAdmin:  claims.Details.IsAdmin,
		}, nil
	}

	log.Errorf("Couldn't parse the token: %v", err)
	return nil, fmt.Errorf("Couldn't parse provided token")
}
