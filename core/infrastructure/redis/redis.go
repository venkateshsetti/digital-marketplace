package redis

import (
	"context"
	"digital-marketplace/config"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	Client *redis.Client
	config    *config.AppConfig
}

func NewRedisRepository(config    *config.AppConfig	) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + fmt.Sprint(config.Redis.Port),
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	log.Println("Successfully connected to redis")
	return &RedisRepository{
		Client: client,
	}, nil
}

func (r *RedisRepository) GetTokenData(key string) (string, error) {
	res := r.Client.Get(context.TODO(), key)
	return res.Result()

}
