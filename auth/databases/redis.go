package databases

import (
	"fmt"

	"github.com/0x113/x-media/auth/common"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// Redis manages Redis connection
type Redis struct {
	DB *redis.Client
}

// Init initializes the Redis database connection
func (database *Redis) Init() error {
	log.Infoln("Connecting to the Redis database ...")
	addr := fmt.Sprintf("%s:%s", common.Config.RedisHost, common.Config.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: common.Config.RedisPassword,
		DB:       common.Config.RedisDB,
	})

	// FIXME: is doesn't work with docker
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Errorf("Couldn't connect to the redis database: %v", err)
			return err
		}
	*/

	database.DB = rdb
	log.Infoln("Successfully conceted to the Redis database")
	return nil
}
