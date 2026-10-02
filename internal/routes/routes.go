package routes

import (
	"net/http"

	"github.com/mehradSaeedi/beauty-store/internal/handlers"
	"gorm.io/gorm"
)

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /api/products", handlers.GetProducts(db))
	mux.HandleFunc("GET /api/products/{id}", handlers.GetProductByID(db))
	mux.HandleFunc("POST /api/products", handlers.CreateProduct(db))
	mux.HandleFunc("PUT /api/products/{id}", handlers.UpdateProduct(db))
	mux.HandleFunc("DELETE /api/products/{id}", handlers.DeleteProduct(db))
	mux.HandleFunc("POST /api/orders", handlers.CreateOrder(db))
	mux.HandleFunc("GET /api/orders", handlers.GetOrders(db))
	mux.HandleFunc("PATCH /api/orders/{id}/status", handlers.UpdateOrderStatus(db))
}
