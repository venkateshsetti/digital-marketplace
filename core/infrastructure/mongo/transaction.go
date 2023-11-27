package mongo

import (
	"context"
	"crypto/sha256"
	core "digital-marketplace/core"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentService simulates an external payment service
type PaymentService struct{}

// GenerateRandomTxHash generates a random transaction hash for a mocked crypto transaction.
func GenerateRandomTxHash() string {
	// Use the current time as a seed for randomness
	rand.Seed(time.Now().UnixNano())

	// Generate a random string (you might use your own logic here)
	randomString := generateRandomString(16)

	// Create a unique identifier by hashing the random string
	txHash := hashString(randomString)

	return txHash
}

// generateRandomString generates a random string of a given length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// hashString hashes a string using SHA-256 and returns the hex-encoded result.
func hashString(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// ProcessPayment simulates the payment process
func (ps *PaymentService) ProcessPayment(amount float64) (bool, string, error) {
	// Simulate payment success for demonstration purposes
	txHash := GenerateRandomTxHash()
	return true, txHash, nil
}

var PurchaseMutex sync.Mutex

// PurchaseItem attempts to purchase an item from the inventory
func (m *MongoDB) PurchaseItem(userID string, itemID string, quantity int32) (*core.ExecuteOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Acquire the mutex to ensure atomicity
	PurchaseMutex.Lock()
	defer PurchaseMutex.Unlock()
	// Check item availability and lock it
	var item core.InventoryItem
	db := m.client.Database(m.config.Mongo.Database)
	collection := db.Collection("items")
	purchaseCollection := db.Collection("purchase_history")
	objId, _ := primitive.ObjectIDFromHex(itemID)
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objId, "quantity": bson.M{"$gte": quantity}, "locked": false},
		bson.M{"$set": bson.M{"locked": true}},options,
	).Decode(&item)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Item %s is not available for purchase", objId.Hex())
	}

	// Simulate payment initiation
	paymentService := &PaymentService{}
	paymentAmount := item.Price * float64(quantity) // Use the item's price for payment
	paymentSuccess, txHash, paymentErr := paymentService.ProcessPayment(paymentAmount)

	// Unlock the item if payment fails
	defer func() {
		if !paymentSuccess {
			_, _ = collection.UpdateOne(
				ctx,
				bson.M{"_id": objId},
				bson.M{"$set": bson.M{"locked": false}},
			)
		}
	}()

	if paymentErr != nil || !paymentSuccess {
		return nil, fmt.Errorf("Payment for item %s failed", objId.Hex())
	}

	// Update the inventory count
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objId, "quantity": bson.M{"$gte": quantity}},
		bson.M{"$inc": bson.M{"quantity": -quantity}},
	)
	_, _ = collection.UpdateOne(
		ctx,
		bson.M{"_id": objId},
		bson.M{"$set": bson.M{"locked": false}},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to update inventory for item %s", objId.Hex())
	}

	record := core.PurchasedHistory{
		ItemId:         objId.String(),
		PurchasedOn:    time.Now().String(),
		Quantity:       quantity,
		TxHash:         txHash,
		UserID:         userID,
		PurchasedPrice: paymentAmount,
	}

	// Insert the document into the collection
	_, err = purchaseCollection.InsertOne(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert  purchased record %v", err)

	}
	fmt.Printf("User %s successfully purchased item %s\n", userID, objId.Hex())
	return &core.ExecuteOrderResponse{TxID: txHash, Message: "Order Placed Sucessfully"}, nil
}
