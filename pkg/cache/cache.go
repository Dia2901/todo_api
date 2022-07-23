package cache

import (
	"gojek/web-server-gin/pkg/config"
	"gojek/web-server-gin/pkg/handleError"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func SetupRedis() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: config.RedisPass,
		DB:       1,
	})

	_, err := redisClient.Ping().Result()
	handleError.Check(err)

	RedisClient = redisClient
}
