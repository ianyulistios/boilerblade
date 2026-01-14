package repository

import (
	"boilerblade/src/model"

	"gorm.io/gorm"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(id uint) (*model.Product, error)
	GetAll(limit, offset int) ([]model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
	Count() (int64, error)
}

// productRepository implements ProductRepository interface
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// Create creates a new product
func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products with pagination
func (r *productRepository) GetAll(limit, offset int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Limit(limit).Offset(offset).Find(&products).Error
	return products, err
}

// Update updates an existing product
func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// Delete soft deletes a product
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// Count returns the total number of products
func (r *productRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Product{}).Count(&count).Error
	return count, err
}
