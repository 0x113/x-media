package service

import (
	"fmt"
	"time"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
)

// AuthService decribes authentication service
type AuthService interface {
	GenerateJWT(accessDetails *models.AccessDetails) (*models.TokenDetails, error)
	ExtractTokenMetadata(tokenString string) (*models.AccessDetails, error)
}

type authService struct{}

// NewAuthService creates new instance of authentication service
func NewAuthService() AuthService {
	return &authService{}
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

	if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
		return &models.AccessDetails{
			Username: claims.Details.Username,
			IsAdmin:  claims.Details.IsAdmin,
		}, nil
	}

	return nil, err
}
