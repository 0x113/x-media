package handler

// NOTE: it's not used in the handler implementation, it's used for docs only
// userPayload represents payload which should be send while creating user
type userPayload struct {
	Username string `json:"username" example:"TheBill"`
	Password string `json:"password" example:"SuperSecretAndStrongPassword123#!"`
}

type userCreateResponse struct {
	Message string `json:"message" example:"Successfully created new user"`
}

type userValidatePayload struct {
	Username string `json:"username" example:"TheBill"`
	Password string `json:"password" example:"hashedPasswordHere"`
}
