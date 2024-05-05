package internal

import (
	"fmt"

	"gorm.io/gorm"
)

type DbProduct struct {
	gorm.Model
	ID          uint64
	Name        string
	Sku         string
	Description string
	Price       float32
	Image       string
}

func (DbProduct) TableName() string {
	return "catalog_products"
}

const (
	ErrorId = 0
)

// Create a new DbProduct
func CreateProduct(db *gorm.DB, product *DbProduct) (uint64, error) {
	result := db.Create(product)
	if result.Error != nil {
		return ErrorId, fmt.Errorf("failed to create a product: %w", result.Error)
	}
	return product.ID, nil
}

// Read a DbProduct by ID
func GetProductByID(db *gorm.DB, id uint64) (*DbProduct, error) {
	var product DbProduct
	result := db.First(&product, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get a product %d: %w", id, result.Error)
	}
	return &product, nil
}

// Update a DbProduct
func UpdateProduct(db *gorm.DB, product *DbProduct) error {
	result := db.Save(product)
	if result.Error != nil {
		return fmt.Errorf("failed to update a product %d: %w", product.ID, result.Error)
	}
	return nil
}

// Delete a DbProduct by ID
func DeleteProductByID(db *gorm.DB, id uint64) error {
	result := db.Delete(&DbProduct{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delet a product %d: %w", id, result.Error)
	}
	return nil
}

// Get all DbProducts
func GetAllProducts(db *gorm.DB) ([]*DbProduct, error) {
	var products []*DbProduct
	result := db.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get products: %w", result.Error)
	}
	return products, nil
}
