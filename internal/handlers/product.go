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

func GetProducts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := repository.GetAllProducts(db)
		if err != nil {
			http.Error(w, "Failed to get products", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func GetProductByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.ParseUint(idStr, 10, 0)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		product, err := repository.GetProductByID(db, uint(id))

		if errors.Is(err, repository.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Failed to get product", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

func CreateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product

		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Invalid JSON request", http.StatusBadRequest)
			return
		}

		err = repository.CreateProduct(db, &product)

		if errors.Is(err, repository.ErrInvalidPrice) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(product)

	}
}

func UpdateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.ParseUint(idStr, 10, 0)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var product models.Product

		err = json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		product.ID = uint(id)

		err = repository.UpdateProduct(db, uint(id), &product)

		if errors.Is(err, repository.ErrInvalidPrice) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, repository.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(product)

	}
}

func DeleteProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.ParseUint(idStr, 10, 0)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		err = repository.DeleteProduct(db, uint(id))

		if errors.Is(err, repository.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}
