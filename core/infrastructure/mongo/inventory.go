package mongo

import (
	"context"
	"digital-marketplace/config"
	core "digital-marketplace/core"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	config *config.AppConfig
	client *mongo.Client
}

func MongoDBService(config *config.AppConfig, client *mongo.Client) *MongoDB {
	return &MongoDB{
		config: config,
		client: client,
	}
}

func (i *MongoDB) GetInventory() ([]core.InventoryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var inventoryItems []core.InventoryItem
	configData := i.config.Mongo

	db := i.client.Database(configData.Database)
	collection := db.Collection("items")

	cursor, err := collection.Find(ctx, bson.D{})
	log.Println("Fetching items...")
	if err != nil {
		log.Println("Error fetching data: ", err)
		return nil, err
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			log.Println("Error closing cursor:", err)
		}
	}()

	for cursor.Next(ctx) {
		var result core.InventoryItem
		if err = cursor.Decode(&result); err != nil {
			log.Println("Error decoding data: ", err)
			return nil, err
		}
		inventoryItems = append(inventoryItems, result)
	}

	if err = cursor.Err(); err != nil {
		log.Println("Error during cursor iteration: ", err)
		return nil, err
	}

	return inventoryItems, nil
}

func (i *MongoDB) PurchasedHistory(userID string) ([]core.PurchasedHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var purchasedHistory []core.PurchasedHistory
	configData := i.config.Mongo

	db := i.client.Database(configData.Database)
	collection := db.Collection("purchase_history")

	filter := bson.M{"user_id":userID}
    
	cursor, err := collection.Find(ctx, filter)
	log.Println("Fetching history...",cursor.Current,filter)
	if err != nil {
		log.Println("Error fetching data: ", err)
		return nil, err
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			log.Println("Error closing cursor:", err)
			return
		}
	}()

	for cursor.Next(ctx) {
		log.Println("started...")
		var result core.PurchasedHistory
		if err = cursor.Decode(&result); err != nil {
			log.Println("Error decoding data: ", err)
			return nil, err
		}
		log.Println(result,"hhhh")
		purchasedHistory = append(purchasedHistory, result)
	}

	if err = cursor.Err(); err != nil {
		log.Println("Error during cursor iteration: ", err)
		return nil, err
	}

	return purchasedHistory, nil
}

func (i *MongoDB) UpdateInventoryItems(item_id string, req core.ItemUpdateRequest) (core.InventoryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configData := i.config.Mongo

	var update primitive.M

	db := i.client.Database(configData.Database)
	collection := db.Collection("items")

	objId, _ := primitive.ObjectIDFromHex(item_id)

	filter := bson.D{{Key: "_id", Value: objId}}

	if req.Quantity == 0 {
		update = bson.M{
			"$set": bson.M{"price": req.Price},
		}
	} else if req.Price == 0.0 {
		update = bson.M{
			"$set": bson.M{"quantity": req.Quantity},
		}
	} else {
		update = bson.M{
			"$set": bson.M{"quantity": req.Quantity, "price": req.Price},
		}
	}

	var resp core.InventoryItem

	result := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	log.Println("Updating item...")
	if result.Err() != nil {
		return core.InventoryItem{}, result.Err()
	}
	if err := result.Decode(&resp); err != nil {
		log.Fatal("Error decoding data:", err)
		return core.InventoryItem{}, err
	}

	return resp, nil
}
