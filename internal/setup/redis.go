package setup

import (
	"os"
	"strconv"

	"github.com/go-redis/redis/v7"
)

// Redis sets up a connection to a Redis server specified by the .env config.
func Redis() *redis.Client {

	db, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		LogCommon(err).Fatal("REDIS_DATABASE type conversion")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	return client
}
