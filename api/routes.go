package api

import (
	"digital-marketplace/api/handlers"
	"digital-marketplace/config"
	"digital-marketplace/core/domain"
	"digital-marketplace/core/infrastructure/mongo"
	"digital-marketplace/core/infrastructure/redis"
	"digital-marketplace/core/utils/coingecko"

	_ "digital-marketplace/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	// swagger embed files
	"github.com/gorilla/mux"
)

func NewRoutes(config *config.AppConfig, mongo *mongo.MongoDB, coingecko *coingecko.Coingecko, domain *domain.Domain,	redisClient *redis.RedisRepository	) *mux.Router {
	webHandler := handlers.NewHandler(config, mongo, coingecko, domain,redisClient)
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/list_coins", webHandler.GetCoinsListHandler()).Methods("GET")
	r.HandleFunc("/api/v1/items", webHandler.InventoryItemsHandler()).Queries("token","{token}").Methods("GET")
	r.HandleFunc("/api/v1/purchasedHistory", webHandler.PreviouslyPurchasedItemsHandler()).Queries("user_id","{user_id}").Methods("GET")
	r.HandleFunc("/api/v1/updateInventory/{item_id}", webHandler.UpdateInventoryItemsHandler()).Methods("PATCH")
	r.HandleFunc("/api/v1/execute_order", webHandler.ExecuteOrderHandler()).Methods("POST")

	// Swagger UI route
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL to API definition
	))
	return r
}
