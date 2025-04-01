package configs

import (
	"backend-developer-assignment/platform/cache"
	"fmt"
	"os"
	"strconv"
)

// RedisConnection creates a new Redis client
func RedisConnection() *cache.RedisClient {
	// Get Redis connection details from environment
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		var err error
		redisDB, err = strconv.Atoi(dbStr)
		if err != nil {
			redisDB = 0
		}
	}

	// Create Redis address
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Create and return Redis client
	return cache.NewRedisClient(redisAddr, redisPassword, redisDB)
}