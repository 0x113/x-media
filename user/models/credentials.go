package models

// Credentials defines credentials which will be use to authenticate user
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
