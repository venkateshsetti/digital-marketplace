package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Item        string             `json:"item" bson:"item"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Quantity    int64              `json:"quantity" bson:"quantity"`
}

type PurchasedHistory struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ItemId         string             `json:"item_id" bson:"item_id"`
	PurchasedOn    string             `json:"purchased_on" bson:"purchased_on"`
	Quantity       string             `json:"quantity" bson:"quantity"`
	TxHash         string             `json:"tx_hash" bson:"tx_hash"`
	WalletId       string             `json:"wallet_id" bson:"wallet_id"`
	PurchasedPrice float64            `json:"purchased_price" bson:"purchased_price"`
}

type ItemUpdateRequest struct {
	Price    float64 `json:"price,omitempty" bson:"price"`
	Quantity int64   `json:"quantity,omitempty" bson:"quantity"`
}
