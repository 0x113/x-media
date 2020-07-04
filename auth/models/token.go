package models

import "github.com/dgrijalva/jwt-go"

// TokenDetails defines token details
type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

// TokenClaims defines custom token claims
type TokenClaims struct {
	Details *AccessDetails
	jwt.StandardClaims
}

// AccessDetails defines access details e.g. isAdmin, username
type AccessDetails struct {
	Username string `json:"username" validate:"required"`
	IsAdmin  *bool  `json:"is_admin" validate:"required"`
}

// TokenString defines the models which will be used to validate the token
type TokenString struct {
	Token string `json:"token"`
}
