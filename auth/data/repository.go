package data

import "github.com/0x113/x-media/auth/models"

// AuthRepository manages the operations on the database for
// the authetication service
type AuthRepository interface {
	Save(username string, token *models.TokenDetails) error
	Delete(uuid string) error
}
