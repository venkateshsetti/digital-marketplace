package handlers

import (
	"digital-marketplace/config"
	core "digital-marketplace/core"
	"digital-marketplace/core/infrastructure/mongo"
	"digital-marketplace/core/utils/coingecko"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	config    *config.AppConfig
	inventory *mongo.Inventory
	coingecko *coingecko.Coingecko
}

func NewHandler(config *config.AppConfig, inventory *mongo.Inventory, coingecko *coingecko.Coingecko) *Handler {
	return &Handler{
		config:    config,
		inventory: inventory,
		coingecko: coingecko,
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
		result, err := h.coingecko.GetCoinList()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching Tokens List %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// @Summary Get a list of items
// @Description Get a list of items from the server
// @Tags         Items
// @Param        token  query string  true "coingecko id for user specified token (use ref list_coins api)"
// @Accept       json
// @Produce      json
// @Router /items/{token} [get]
func (h *Handler) InventoryItemsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]
		result, err := h.inventory.GetInventory(token)
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
// @Router /purchasedHistory/{wallet_id} [get]
func (h *Handler) PreviouslyPurchasedItemsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		wallet_id := vars["wallet_id"]
		result, err := h.inventory.PurchasedHistory(wallet_id)
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
// @Router /updateInventory/{item_id} [get]
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

		result, err := h.inventory.UpdateInventoryItems(item_id, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating inventory of given item id %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// @Summary Get a list of items
// @Description Get a list of items from the server
// @Produce json
func (h *Handler) ExecuteOrderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
