package models

import "github.com/dgrijalva/jwt-go"

// TokenDetails defines token details
type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenClaims defines custom token claims
type TokenClaims struct {
	Username string
	IsAdmin  bool
	jwt.StandardClaims
}
