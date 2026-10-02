package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/mehradSaeedi/beauty-store/internal/models"
	"github.com/mehradSaeedi/beauty-store/internal/repository"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerAddress string `json:"customer_address"`

	Items []struct {
		ProductID uint `json:"id"`
		Quantity  uint `json:"quantity"`
	} `json:"items"`
}

func CreateOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateOrderRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(request.Items) == 0 {
			http.Error(w, "Order must contain at least one element", http.StatusBadRequest)
			return
		}

		if request.CustomerName == "" ||
			request.CustomerPhone == "" ||
			request.CustomerAddress == "" {
			http.Error(w, "Missing customer information", http.StatusBadRequest)
			return
		}

		items := make([]models.OrderItem, 0, len(request.Items))

		for _, item := range request.Items {
			items = append(items, models.OrderItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		err = repository.CreateOrder(
			db,
			request.CustomerName,
			request.CustomerPhone,
			request.CustomerAddress,
			items,
		)
		if err != nil {
			if errors.Is(err, repository.ErrNotEnoughStock) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}

			if errors.Is(err, repository.ErrProductNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			if errors.Is(err, repository.ErrInvalidQuantity) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.Error(w, "Could not create order", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Order created successfully ",
		})
	}
}

func GetOrders(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := repository.GetOrders(db)
		if err != nil {
			http.Error(w, "Could not get orders", http.StatusInternalServerError)
			log.Println("GetOrders error:", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
	}
}

func UpdateOrderStatus(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		var request struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		err = repository.UpdateOrderStatus(db, uint(id), request.Status)
		if errors.Is(err, repository.ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if errors.Is(err, repository.ErrInvalidStatus) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Println("UpdateOrderRequest error:", err.Error())
			http.Error(w, "Could not update order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
