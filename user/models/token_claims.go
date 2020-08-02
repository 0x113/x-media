package models

// TokenClaims defines details like username and is_admin, which are needed
// to generate access token
type TokenClaims struct {
	Username string `json:"username" example:"TheBill"`
	IsAdmin  bool   `json:"is_admin" example:"false"`
}
