package mongo

import (
	"context"
	"digital-marketplace/config"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectToMongoDB(conf *config.AppConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.Mongo.Host+"://"+conf.Mongo.Username+":"+conf.Mongo.Password+"@cluster0.kheedti.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to MongoDB")
	return client, nil
}
