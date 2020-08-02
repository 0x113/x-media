package models

import "github.com/dgrijalva/jwt-go"

// TokenDetails defines token details
type TokenDetails struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEZXRhaWxzIjp7InVzZXJuYW1lIjoiSm9obkRvZSIsImlzX2FkbWluIjpmYWxzZX0sIlV1aWQiOiJmMTk0YWZkYy1iNTA1LTRjMmYtYTc1NC02ZTQ0NjA5YzZlODAiLCJleHAiOjE1OTQ1NzUwMzB9.h9YpZNRkriaBvi3c1kt9Rm6NyWAfKDI2a2y2gQRCOOU"`
	AccessUuid   string `json:"-"`
	RefreshUuid  string `json:"-"`
	AtExpires    int64  `json:"-"`
	RtExpires    int64  `json:"-"`
}

// TokenClaims defines custom token claims
type TokenClaims struct {
	Details *AccessDetails
	Uuid    string
	jwt.StandardClaims
}

// AccessDetails defines access details e.g. isAdmin, username
type AccessDetails struct {
	Username string `json:"username" validate:"required"`
	IsAdmin  *bool  `json:"is_admin" validate:"required"`
}

// UuidAccessDetails defines extented AccessDetails model with token uuid NOTE: should be named better
type UuidAccessDetails struct {
	*AccessDetails
	Uuid string
}

// TokenString defines the models which will be used to validate the token
type TokenString struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEZXRhaWxzIjp7InVzZXJuYW1lIjoiSm9obkRvZSIsImlzX2FkbWluIjpmYWxzZX0sIlV1aWQiOiJmMTk0YWZkYy1iNTA1LTRjMmYtYTc1NC02ZTQ0NjA5YzZlODAiLCJleHAiOjE1OTQ1NzUwMzB9.h9YpZNRkriaBvi3c1kt9Rm6NyWAfKDI2a2y2gQRCOOU"`
}
