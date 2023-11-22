package main

import (
	"context"
	"digital-marketplace/api"
	"digital-marketplace/config"
	"digital-marketplace/core/infrastructure/mongo"
	"fmt"
	"log"
	"net/http"
)

func main() {

	conf := config.LoadConfig()

	client, err := mongo.ConnectToMongoDB(conf)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB %v", err.Error())
	}
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	inventory := mongo.InventoryService(conf, client)

	//redisRepo, _ := redis.NewRedisRepository()

	router := routes.NewRoutes(conf, inventory)

	port := conf.Server.Port
	serverAddress := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", serverAddress)
	err = http.ListenAndServe(serverAddress, router)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
