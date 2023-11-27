package main

import (
	"context"
	"digital-marketplace/api"
	"digital-marketplace/config"
	"digital-marketplace/core/domain"
	"digital-marketplace/core/infrastructure/mongo"
	"digital-marketplace/core/infrastructure/redis"
	"digital-marketplace/core/utils/coingecko"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// @title Digital Marketplace
// @version 1.0
// @description Digital Market place api let users to view,list and make a purchase of items listed in inventory
// @BasePath /api/v1

func main() {

	conf := config.LoadConfig()

	client, err := mongo.ConnectToMongoDB(conf)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB %v", err.Error())
	}
	defer client.Disconnect(context.Background())

	mongo := mongo.MongoDBService(conf, client)
	coingecko := coingecko.CoingeckoService(conf)

	redisRepo, err := redis.NewRedisRepository(conf)
	if err != nil {
		log.Fatalf("Error connecting to Redis %v", err.Error())
	}
	res, err := coingecko.GetCoinList()
	if err != nil {
		log.Println("Error:-", err)
		return
	}
	// Convert struct to JSON
	jsonData, err := json.Marshal(*res)
	if err != nil {
		log.Fatal(err)
	}
	redisRepo.Client.Set(context.TODO(), "list", jsonData, 0)
	for _, data := range *res {
		redisRepo.Client.Set(context.TODO(), data.ID, data.Symbol, 0)
	}
	log.Printf("Updated the Cache with %d keys",len(*res))
	domainRepo := domain.NewDomain(conf, mongo, coingecko, redisRepo)

	router := api.NewRoutes(conf, mongo, coingecko, domainRepo,redisRepo)

	port := conf.Server.Port
	serverAddress := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", serverAddress)
	err = http.ListenAndServe(serverAddress, router)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
