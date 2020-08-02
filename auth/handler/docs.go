package handler

// NOTE: it's not used in the handler implementation, it's used for docs only
type generateTokenPayload struct {
	Username string `json:"username" example:"TheBill"`
	Password string `json:"password" example:"SuperSecretAndStrongPassword123#!"`
}
