package backend

import (
	"log"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient initializes and returns a new Redis client
func NewRedisClient(config *ServerBackendConfig) (*redis.Client, error) {
	options, err := redis.ParseURL("redis://:@localhost:6379/0")

	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
		return nil, err
	}

	client := redis.NewClient(options)
	return client, nil
}
