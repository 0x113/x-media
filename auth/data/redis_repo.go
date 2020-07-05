package data

import (
	"context"
	"time"

	"github.com/0x113/x-media/auth/databases"
	"github.com/0x113/x-media/auth/models"
)

// authRepository manages the authentication CRUD
type authRepository struct{}

// NewRedisAuthRepository returns a new instace of the authentication repository
func NewRedisAuthRepository() AuthRepository {
	return &authRepository{}
}

// Save stores the token details in the Redis database
func (r *authRepository) Save(username string, token *models.TokenDetails) error {
	// convert Unix to UTC
	atExpires := time.Unix(token.AtExpires, 0)
	rtExpires := time.Unix(token.RtExpires, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := databases.Database.DB.Set(ctx, token.AccessUuid, username, atExpires.Sub(time.Now())).Err(); err != nil {
		return err
	}

	if err := databases.Database.DB.Set(ctx, token.RefreshUuid, username, rtExpires.Sub(time.Now())).Err(); err != nil {
		return err
	}

	return nil
}
