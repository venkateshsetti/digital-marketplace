package api

import (
	"digital-marketplace/api/handlers"
	"digital-marketplace/config"
	"digital-marketplace/core/infrastructure/mongo"
	"digital-marketplace/core/utils/coingecko"

	_ "digital-marketplace/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	// swagger embed files
	"github.com/gorilla/mux"
)

func NewRoutes(config *config.AppConfig, inventory *mongo.Inventory, coingecko *coingecko.Coingecko) *mux.Router {
	webHandler := handlers.NewHandler(config, inventory, coingecko)
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/list_coins", webHandler.GetCoinsListHandler()).Methods("GET")
	r.HandleFunc("/api/v1/items/{token}", webHandler.InventoryItemsHandler()).Methods("GET")
	r.HandleFunc("/api/v1/purchasedHistory/{wallet_id}", webHandler.PreviouslyPurchasedItemsHandler()).Methods("GET")
	r.HandleFunc("/api/v1/updateInventory/{item_id}", webHandler.UpdateInventoryItemsHandler()).Methods("PATCH")
	r.HandleFunc("/api/v1/execute_order", webHandler.InventoryItemsHandler()).Methods("POST")

	// Swagger UI route
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL to API definition
	))
	return r
}
