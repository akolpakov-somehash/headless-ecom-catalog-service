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

type DbWrapper interface {
	Create(interface{}) *gorm.DB
	First(interface{}, ...interface{}) *gorm.DB
	Save(interface{}) *gorm.DB
	Delete(interface{}, ...interface{}) *gorm.DB
	Find(interface{}, ...interface{}) *gorm.DB
}

type ProductServiceInterface interface {
	CreateProduct(product *DbProduct) (uint64, error)
	GetProductByID(id uint64) (*DbProduct, error)
	UpdateProduct(product *DbProduct) error
	DeleteProductByID(id uint64) error
	GetAllProducts() ([]*DbProduct, error)
}

type ProductService struct {
	DB DbWrapper
}

// Create a new DbProduct
func (p *ProductService) CreateProduct(product *DbProduct) (uint64, error) {
	result := p.DB.Create(product)
	if result.Error != nil {
		return ErrorId, fmt.Errorf("failed to create a product: %w", result.Error)
	}
	return product.ID, nil
}

// Read a DbProduct by ID
func (p *ProductService) GetProductByID(id uint64) (*DbProduct, error) {
	product := DbProduct{}
	result := p.DB.First(&product, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get a product %d: %w", id, result.Error)
	}
	return &product, nil
}

// Update a DbProduct
func (p *ProductService) UpdateProduct(product *DbProduct) error {
	result := p.DB.Save(product)
	if result.Error != nil {
		return fmt.Errorf("failed to update a product %d: %w", product.ID, result.Error)
	}
	return nil
}

// Delete a DbProduct by ID
func (p *ProductService) DeleteProductByID(id uint64) error {
	result := p.DB.Delete(&DbProduct{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete a product %d: %w", id, result.Error)
	}
	return nil
}

// Get all DbProducts
func (p *ProductService) GetAllProducts() ([]*DbProduct, error) {
	var products []*DbProduct
	result := p.DB.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get products: %w", result.Error)
	}
	return products, nil
}
