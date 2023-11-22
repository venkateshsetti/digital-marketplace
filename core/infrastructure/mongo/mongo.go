package mongo

import (
	"context"
	"digital-marketplace/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

func ConnectToMongoDB(conf *config.AppConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.Mongo.Host+"://"+conf.Mongo.Username+":"+conf.Mongo.Password+"@cluster1.clc0pny.mongodb.net/?retryWrites=true&w=majority"))
	log.Println("Successfully connected to MongoDB")
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return client, nil
}
