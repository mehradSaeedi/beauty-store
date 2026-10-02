package repository

import (
	"errors"

	"github.com/mehradSaeedi/beauty-store/internal/models"
	"gorm.io/gorm"
)

var (
	ErrInvalidQuantity = errors.New("Quantity must be greater than zero")
	ErrInvalidPrice    = errors.New("Price must be greater than zero")
	ErrNotEnoughStock  = errors.New("Not enough stock")
	ErrProductNotFound = errors.New("Product not found")
	ErrOrderNotFound   = errors.New("Order not found")
	ErrInvalidStatus   = errors.New("Invalid order status")
)

func GetAllProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product

	result := db.Find(&products)

	return products, result.Error
}

func GetProductByID(db *gorm.DB, id uint) (models.Product, error) {
	var product models.Product

	result := db.First(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return product, ErrProductNotFound
	}

	return product, result.Error
}

func CreateProduct(db *gorm.DB, product *models.Product) error {

	if product.Price <= 0 {
		return ErrInvalidPrice
	}

	result := db.Create(product)
	return result.Error
}

func UpdateProduct(db *gorm.DB, id uint, product *models.Product) error {
	var existingProduct models.Product

	result := db.First(&existingProduct, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrProductNotFound
	}

	if result.Error != nil {
		return result.Error
	}

	if product.Price <= 0 {
		return ErrInvalidPrice
	}

	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Stock = product.Stock
	existingProduct.ImageURL = product.ImageURL

	result = db.Save(&existingProduct)

	return result.Error
}

func DeleteProduct(db *gorm.DB, id uint) error {
	var product models.Product

	result := db.First(&product, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrProductNotFound
	}

	if result.Error != nil {
		return result.Error
	}

	// This does a soft delete becuase our model embeds gorm.Model
	result = db.Delete(&product)
	return result.Error

}

func CreateOrder(
	db *gorm.DB,
	customerName string,
	customerPhone string,
	customerAddress string,
	items []models.OrderItem,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var totalPrice int64

		order := models.Order{
			CustomerName:    customerName,
			CustomerPhone:   customerPhone,
			CustomerAddress: customerAddress,
			Status:          "pending",
		}

		for i := range items {
			item := &items[i]

			if item.Quantity == 0 {
				return ErrInvalidQuantity
			}

			var product models.Product

			result := tx.First(&product, item.ProductID)

			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ErrProductNotFound
			}

			if result.Error != nil {
				return result.Error
			}

			result = tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrNotEnoughStock
			}

			item.UnitPrice = product.Price
			totalPrice += product.Price * int64(item.Quantity)
		}

		order.TotalPrice = totalPrice

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID

			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func GetOrders(db *gorm.DB) ([]models.Order, error) {
	var orders []models.Order

	// Preload tells GORM: "when you fetch orders, also fetch their associated OrderItemsa and Product"
	result := db.
		Preload("Items").
		Preload("Items.Product").
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func UpdateOrderStatus(db *gorm.DB, orderID uint, status string) error {
	switch status {
	case "pending", "processing", "shipped", "cancelled", "completed":
		// valid
	default:
		return ErrInvalidStatus
	}

	result := db.
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}
