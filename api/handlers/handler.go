package handlers

import (
	"context"
	"digital-marketplace/config"
	"digital-marketplace/core"
	"digital-marketplace/core/domain"
	"digital-marketplace/core/infrastructure/mongo"
	"digital-marketplace/core/infrastructure/redis"
	"digital-marketplace/core/utils/coingecko"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	config    *config.AppConfig
	mongo     *mongo.MongoDB
	coingecko *coingecko.Coingecko
	domain    *domain.Domain
	redisClient *redis.RedisRepository
}

func NewHandler(config *config.AppConfig, mongo *mongo.MongoDB, coingecko *coingecko.Coingecko, domain *domain.Domain,	redisClient *redis.RedisRepository) *Handler {
	return &Handler{
		config:    config,
		mongo:     mongo,
		coingecko: coingecko,
		domain:    domain,
		redisClient: redisClient,
	}
}

// @Summary Get a list of tokens
// @Description Get a list of tokens to get the price of specified token for items api
// @Tags         Coins
// @Accept       json
// @Produce      json
// @Router  /list_coins [get]
func (h *Handler) GetCoinsListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data coingecko.CoinsList
		result := h.redisClient.Client.Get(context.TODO(),"list")
		res,err := result.Result()
		if err != nil{
			log.Println("Error: ",err)
		}
		json.Unmarshal([]byte(res),&data)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

// @Summary Get a list of items
// @Description Get a list of items from the server
// @Tags         Items
// @Param        token  query string  true "coingecko id for user specified token (use ref list_coins api)"
// @Accept       json
// @Produce      json
// @Router /items [get]
func (h *Handler) InventoryItemsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println(vars)
		token := vars["token"]
		result, err := h.domain.GetInventory(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching inventory items %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// @Summary Review the items
// @Description Get a list of items purchased by given wallet id
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        user_id  query string  true "user id for user to review his purchase history"
// @Router /purchasedHistory [get]
func (h *Handler) PreviouslyPurchasedItemsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		wallet_id := vars["user_id"]
		result, err := h.mongo.PurchasedHistory(wallet_id)
		log.Println(result,err)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching purchased histroy of given wallet id %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// @Summary Update the quantity and/or price of an item
// @Description  Update the quantity and/or price of an item by an admin
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        item_id  path string  true "item id for admin to update the details"
// @Param        data    body core.ItemUpdateRequest true "To update price and/or quantity "
// @Router /updateInventory/{item_id} [patch]
func (h *Handler) UpdateInventoryItemsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		item_id := vars["item_id"]

		decoder := json.NewDecoder(r.Body)
		var req core.ItemUpdateRequest
		err := decoder.Decode(&req)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(req)

		result, err := h.mongo.UpdateInventoryItems(item_id, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating inventory of given item id %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// @Summary Execute Order
// @Description Creates a record after validating the transaction against the given order deta
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param  data  body  core.ExecuteOrderRequest true "details of purchase order"
// @Router /execute_order [post]
func (h *Handler) ExecuteOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData core.ExecuteOrderRequest

		// Decode the JSON request body into the struct
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			// Handle the error, e.g., return a bad request response
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if requestData.Quantity < 1{
			http.Error(w,"Quantity cant be less than 1",http.StatusNotAcceptable)
			return
		}
		result, err := h.domain.ExecuteOrder(requestData.UserID, requestData.ItemID, requestData.Quantity)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error purchaseing item of given item id %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
