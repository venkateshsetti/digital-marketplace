package redis

import (
	"context"
	"digital-marketplace/core/utils/coingecko"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client    *redis.Client
	coingecko *coingecko.Coingecko
}

func NewRedisRepository(coingecko *coingecko.Coingecko) (*RedisRepository, error) {
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
		client:    client,
		coingecko: coingecko,
	}, nil
}

func (r *RedisRepository) SetData() {
	res, err := r.coingecko.GetCoinList()
	if err != nil {
		log.Println("Error:-", err)
		return
	}
	for _, data := range *res {
		r.client.Set(context.TODO(), data.ID, data.Symbol, 0)
	}
	return
}

func (r *RedisRepository) GetTokenData(key string) (string, error) {
	res := r.client.Get(context.TODO(), key)
	return res.Result()

}
