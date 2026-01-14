package dto

import "boilerblade/src/model"

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	Name string `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Name string `json:"name" validate:"omitempty,required"`
	Price float64 `json:"price" validate:"omitempty,required"`
}

// ProductResponse represents the response payload for a single product
type ProductResponse struct {
	ID        uint   `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ProductListResponse represents the response payload for a list of products with pagination
type ProductListResponse struct {
	Products      []ProductResponse `json:"products"`
	Total      int64          `json:"total"`
	Limit      int            `json:"limit"`
	Offset     int            `json:"offset"`
	TotalPages int            `json:"total_pages"`
}

// ToProductResponse converts a model.Product to ProductResponse DTO
func ToProductResponse(product *model.Product) ProductResponse {
	return ProductResponse{
		ID: product.ID,
		// TODO: Map other fields
		CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
