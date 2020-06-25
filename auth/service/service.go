package service

import (
	"time"

	"github.com/0x113/x-media/auth/common"
	"github.com/0x113/x-media/auth/models"

	"github.com/dgrijalva/jwt-go"
)

// AuthService decribes authentication service
type AuthService interface {
	GenerateJWT(username string, isAdmin bool) (*models.TokenDetails, error)
}

type authService struct{}

// NewAuthService creates new instance of authentication service
func NewAuthService() AuthService {
	return &authService{}
}

// GenerateJWT generates new token from provided data
func (s *authService) GenerateJWT(username string, isAdmin bool) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	var err error
	// access token
	atClaims := &models.TokenClaims{
		username,
		isAdmin,
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
		username,
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(common.Config.RefreshSecret))

	return td, err
}
