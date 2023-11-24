package mongo

import (
	"context"
	core "digital-marketplace/core"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PaymentService simulates an external payment service
type PaymentService struct{}

// ProcessPayment simulates the payment process
func (ps *PaymentService) ProcessPayment(amount float64) (bool, error) {
	// Simulate payment success for demonstration purposes
	return true, nil
}

var PurchaseMutex sync.Mutex

// PurchaseItem attempts to purchase an item from the inventory
func PurchaseItem(ctx context.Context, inventoryCollection *mongo.Collection, userID int, itemID primitive.ObjectID, purchaseChannel chan struct{}) error {
	// Acquire the mutex to ensure atomicity
	PurchaseMutex.Lock()
	defer PurchaseMutex.Unlock()
	// Check item availability and lock it
	var item core.InventoryItem
	err := inventoryCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": itemID, "quantity": bson.M{"$gt": 0}, "locked": false},
		bson.M{"$set": bson.M{"locked": true}},
	).Decode(&item)

	if err != nil {
		return log.Printf("Item %s is not available for purchase", itemID.Hex())
	}

	// Simulate payment initiation
	paymentService := &PaymentService{}
	paymentAmount := item.Price // Use the item's price for payment
	paymentSuccess, paymentErr := paymentService.ProcessPayment(paymentAmount)

	// Unlock the item if payment fails
	defer func() {
		if !paymentSuccess {
			_, _ = inventoryCollection.UpdateOne(
				ctx,
				bson.M{"_id": itemID},
				bson.M{"$set": bson.M{"locked": false}},
			)
		}
	}()

	if paymentErr != nil || !paymentSuccess {
		return fmt.Errorf("Payment for item %s failed", itemID.Hex())
	}

	// Update the inventory count
	_, err = inventoryCollection.UpdateOne(
		ctx,
		bson.M{"_id": itemID, "quantity": bson.M{"$gt": 0}},
		bson.M{"$inc": bson.M{"quantity": -1}},
	)

	if err != nil {
		return fmt.Errorf("Failed to update inventory for item %s", itemID.Hex())
	}

	fmt.Printf("User %d successfully purchased item %s\n", userID, itemID.Hex())
	return nil
}
