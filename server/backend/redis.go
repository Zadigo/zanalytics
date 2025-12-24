package backend

import (
	"log"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient initializes and returns a new Redis client
func NewRedisClient(addr, password string, db int) *redis.Client {
	options, err := redis.ParseURL("redis://:@localhost:6379/0")

	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
		return nil
	}

	client := redis.NewClient(options)
	return client
}
