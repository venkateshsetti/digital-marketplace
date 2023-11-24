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

type Inventory struct {
	config *config.AppConfig
	client *mongo.Client
}

func InventoryService(config *config.AppConfig, client *mongo.Client) *Inventory {
	return &Inventory{
		config: config,
		client: client,
	}
}

func (i *Inventory) GetInventory() ([]core.InventoryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var inventoryItems []core.InventoryItem
	configData := i.config.Mongo

	db := i.client.Database(configData.Database)
	collection := db.Collection("items")

	cursor, err := collection.Find(ctx, bson.D{})
	log.Println("Fetching items...")
	if err != nil {
		log.Fatal("Error fetching data: ", err)
		return nil, err
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			log.Fatal("Error closing cursor:", err)
		}
	}()

	for cursor.Next(ctx) {
		var result core.InventoryItem
		if err = cursor.Decode(&result); err != nil {
			log.Fatal("Error decoding data: ", err)
			return nil, err
		}
		inventoryItems = append(inventoryItems, result)
	}

	if err = cursor.Err(); err != nil {
		log.Fatal("Error during cursor iteration: ", err)
		return nil, err
	}

	return inventoryItems, nil
}

func (i *Inventory) PurchasedHistory(walletID string) ([]core.PurchasedHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var purchasedHistory []core.PurchasedHistory
	configData := i.config.Mongo

	db := i.client.Database(configData.Database)
	collection := db.Collection("purchase_history")

	filter := bson.D{{Key: "wallet_id", Value: walletID}}

	cursor, err := collection.Find(ctx, filter)
	log.Println("Fetching history...")
	if err != nil {
		log.Fatal("Error fetching data: ", err)
		return nil, err
	}
	defer func() {
		if err = cursor.Close(ctx); err != nil {
			log.Fatal("Error closing cursor:", err)
		}
	}()

	for cursor.Next(ctx) {
		var result core.PurchasedHistory
		if err = cursor.Decode(&result); err != nil {
			log.Println("err", err)
			log.Fatal("Error decoding data: ", err)
			return nil, err
		}
		purchasedHistory = append(purchasedHistory, result)
	}

	if err = cursor.Err(); err != nil {
		log.Fatal("Error during cursor iteration: ", err)
		return nil, err
	}

	return purchasedHistory, nil
}

func (i *Inventory) UpdateInventoryItems(item_id string, req core.ItemUpdateRequest) (core.InventoryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configData := i.config.Mongo

	var update primitive.M

	db := i.client.Database(configData.Database)
	collection := db.Collection("inventory")

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
