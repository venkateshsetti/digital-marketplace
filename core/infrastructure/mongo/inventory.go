package mongo

import (
	"context"
	"digital-marketplace/config"
	"digital-marketplace/core/domain"
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

func (i *Inventory) GetInventory() ([]domain.InventoryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var inventoryItems []domain.InventoryItem
	configData := i.config.Mongo

	db := i.client.Database(configData.Database)
	collection := db.Collection("inventory")

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
		var result domain.InventoryItem
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

func (i *Inventory) PurchasedHistory(wallet_id string) ([]domain.PurchasedHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var purchasedHistory []domain.PurchasedHistory
	configData := i.config.Mongo

	db := i.client.Database(configData.Database)
	collection := db.Collection("purchase_history")

	filter := bson.D{{"wallet_id", wallet_id}}

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
		var result domain.PurchasedHistory
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

func (i *Inventory) UpdateInventoryItems(item_id string, req domain.ItemUpdateRequest) (domain.InventoryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configData := i.config.Mongo

	var update primitive.M

	db := i.client.Database(configData.Database)
	collection := db.Collection("inventory")

	objId, _ := primitive.ObjectIDFromHex(item_id)

	filter := bson.D{{"_id", objId}}

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

	var resp domain.InventoryItem

	result := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	log.Println("Updating item...")
	if result.Err() != nil {
		return domain.InventoryItem{}, result.Err()
	}
	if err := result.Decode(&resp); err != nil {
		log.Fatal("Error decoding data:", err)
		return domain.InventoryItem{}, err
	}

	return resp, nil
}
