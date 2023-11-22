package routes

import (
	"digital-marketplace/api/handlers"
	"digital-marketplace/config"
	"digital-marketplace/core/infrastructure/mongo"
	"github.com/gorilla/mux"
)

func NewRoutes(config *config.AppConfig, inventory *mongo.Inventory) *mux.Router {
	webHandler := handlers.NewHandler(config, inventory)
	r := mux.NewRouter()

	r.HandleFunc("/items", webHandler.InventoryItemsHandler()).Methods("GET")
	r.HandleFunc("/purchasedHistory/{wallet_id}", webHandler.PreviouslyPurchasedItemsHandler()).Methods("GET")
	r.HandleFunc("/updateInventory/{item_id}", webHandler.UpdateInventoryItemsHandler()).Methods("PATCH")
	return r
}
