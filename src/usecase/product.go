package usecase

import (
	"boilerblade/src/dto"
	"boilerblade/src/model"
	"boilerblade/src/repository"
	"errors"
	"math"
)

// ProductUsecase defines the interface for product business logic
type ProductUsecase interface {
	CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProductByID(id uint) (*dto.ProductResponse, error)
	GetAllProducts(limit, offset int) (*dto.ProductListResponse, error)
	UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error
}

// productUsecase implements ProductUsecase interface
type productUsecase struct {
	productRepo repository.ProductRepository
}

// NewProductUsecase creates a new product usecase instance
func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (uc *productUsecase) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Create product model
	product := &model.Product{
		// TODO: Map fields from request to model
	}

	// Save to database
	if err := uc.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Return response DTO
	return &dto.ProductResponse{
		ID: product.ID,
		// TODO: Map other fields
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetProductByID retrieves a product by ID
func (uc *productUsecase) GetProductByID(id uint) (*dto.ProductResponse, error) {
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return &dto.ProductResponse{
		ID: product.ID,
		// TODO: Map other fields
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetAllProducts retrieves all products with pagination
func (uc *productUsecase) GetAllProducts(limit, offset int) (*dto.ProductListResponse, error) {
	// Validate pagination
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Get products and total count
	products, err := uc.productRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.productRepo.Count()
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dto.ProductResponse{
			ID: product.ID,
			// TODO: Map other fields
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.ProductListResponse{
		Products:      productResponses,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}, nil
}

// UpdateProduct updates an existing product
func (uc *productUsecase) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Get existing product
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// TODO: Update fields if provided

	// Save updates
	if err := uc.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Return response DTO
	return &dto.ProductResponse{
		ID: product.ID,
		// TODO: Map other fields
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// DeleteProduct deletes a product
func (uc *productUsecase) DeleteProduct(id uint) error {
	// Check if product exists
	_, err := uc.productRepo.GetByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	// Delete product
	return uc.productRepo.Delete(id)
}
