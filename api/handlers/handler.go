package handlers

import (
	"digital-marketplace/config"
	"digital-marketplace/core/domain"
	"digital-marketplace/core/infrastructure/mongo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	config    *config.AppConfig
	inventory *mongo.Inventory
}

func NewHandler(config *config.AppConfig, inventory *mongo.Inventory) *Handler {
	return &Handler{
		config:    config,
		inventory: inventory,
	}
}

func (h *Handler) InventoryItemsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := h.inventory.GetInventory()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching inventory items %v", err.Error()), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

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

func (h *Handler) UpdateInventoryItemsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		item_id := vars["item_id"]

		decoder := json.NewDecoder(r.Body)
    var req domain.ItemUpdateRequest
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
