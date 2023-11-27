package core

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Item        string             `json:"item" bson:"item"`
	Description string             `json:"description" bson:"description"`
	Price       float64           `json:"price" bson:"price"`
	Quantity    int32              `json:"quantity" bson:"quantity"`
	Locked      bool               `json:"-" bson:"locked"`
}

type PurchasedHistory struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ItemId         string             `json:"item_id" bson:"item_id"`
	PurchasedOn    string             `json:"purchased_on" bson:"purchased_on"`
	Quantity       int32              `json:"quantity" bson:"quantity"`
	TxHash         string             `json:"tx_hash" bson:"tx_hash"`
	UserID       string             `json:"user_id" bson:"user_id"`
	PurchasedPrice float64            `json:"purchased_price" bson:"purchased_price"`
}

type ItemUpdateRequest struct {
	Price    float64 `json:"price,omitempty" bson:"price"`
	Quantity int64   `json:"quantity,omitempty" bson:"quantity"`
}

type ExecuteOrderRequest struct {
	ItemID   string `json:"item_id"`
	UserID   string `json:"user_id"`
	Quantity int32  `json:"quantity`
}

type ExecuteOrderResponse struct {
	TxID    string `json:"tx_id"`
	Message string `json:"message"`
}
