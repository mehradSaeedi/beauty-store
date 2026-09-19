package repository

import (
	"github.com/mehradSaeedi/beauty-store/internal/models"
	"gorm.io/gorm"
)

func GetAllProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product

	result := db.Find(&products)

	return products, result.Error
}

func GetProductByID(db *gorm.DB, id uint) (models.Product, error) {
	var product models.Product

	result := db.First(&product, id)

	return product, result.Error
}

func CreateProduct(db *gorm.DB, product *models.Product) error {
	result := db.Create(product)
	return result.Error
}

func UpdateProduct(db *gorm.DB, id uint, product *models.Product) error {
	var existingProduct models.Product

	result := db.First(&existingProduct, id)
	if result.Error != nil {
		return result.Error
	}

	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Stock = product.Stock

	result = db.Save(&existingProduct)

	return result.Error
}

func DeleteProduct(db *gorm.DB, id uint) error {
	var product models.Product

	result := db.First(&product, id)
	if result.Error != nil {
		return result.Error
	}

	// This does a soft delete becuase our model embeds gorm.Model
	result = db.Delete(&product)

	return result.Error

}
