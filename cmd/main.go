package main

import (
	"context"
	"digital-marketplace/api"
	"digital-marketplace/config"
	"digital-marketplace/core/domain"
	"digital-marketplace/core/infrastructure/mongo"
	"digital-marketplace/core/infrastructure/redis"
	"digital-marketplace/core/utils/coingecko"
	"fmt"
	"log"
	"net/http"
)

// @title Digital Marketplace
// @version 1.0
// @description Digital Market place api let users to view,list and make a purchase of items listed in inventory
// @host localhost:8080
// @BasePath /api/v1

func main() {

	conf := config.LoadConfig()

	client, err := mongo.ConnectToMongoDB(conf)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB %v", err.Error())
	}
	defer client.Disconnect(context.Background())

	inventory := mongo.InventoryService(conf, client)
	coingecko := coingecko.CoingeckoService(conf)

	redisRepo, err := redis.NewRedisRepository()
	if err != nil {
		log.Fatalf("Error connecting to Redis %v", err.Error())
	}
	res, err := coingecko.GetCoinList()
	if err != nil {
		log.Println("Error:-", err)
		return
	}
	for _, data := range *res {
		redisRepo.Client.Set(context.TODO(), data.ID, data.Symbol, 0)
	}
	domainRepo := domain.NewDomain(conf, inventory, coingecko)

	router := api.NewRoutes(conf, inventory, coingecko, domainRepo)

	port := conf.Server.Port
	serverAddress := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", serverAddress)
	err = http.ListenAndServe(serverAddress, router)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
