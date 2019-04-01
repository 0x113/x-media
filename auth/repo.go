package auth

// AuthRepository represtents authentication repository.
type AuthRepository interface {
	Create(user *User) error
	GenerateJWT(user *User) (string, error)
}
