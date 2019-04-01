package auth

import (
	"encoding/json"
	"net/http"
)

// AuthHandler does something, but I don't know how to write it TODO: make this comment more clear
type AuthHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GenerateJWT(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService AuthService
}

// NewAuthHandler returns new instance of authentication handler.
func NewAuthHandler(authService AuthService) AuthHandler {
	return &authHandler{
		authService,
	}
}

func (h *authHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	err := h.authService.CreateUser(&user)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to create user"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully Created!"})
}

func (h *authHandler) GenerateJWT(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	token, err := h.authService.LoginUser(user.Username, user.Password)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to generate jwt"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
