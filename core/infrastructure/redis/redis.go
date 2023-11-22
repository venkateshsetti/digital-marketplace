package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository() (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	return &RedisRepository{
		client: client,
	}, nil
}
