package handler

import (
	"boilerblade/helper"
	"boilerblade/src/dto"
	"boilerblade/src/usecase"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	productUsecase usecase.ProductUsecase
	validator                   *validator.Validate
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		validator:                   validator.New(),
	}
}

// RegisterRoutes registers all product routes to the router group
func (h *ProductHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/products", h.GetAllProducts)
	router.Get("/products/:id", h.GetProduct)
	router.Post("/products", h.CreateProduct)
	router.Put("/products/:id", h.UpdateProduct)
	router.Delete("/products/:id", h.DeleteProduct)
}

// CreateProduct handles the creation of a new product
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		helper.LogError("Failed to parse request body", err, c.Path(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	productResponse, err := h.productUsecase.CreateProduct(&req)
	if err != nil {
		helper.LogError("Failed to create product", err, c.Path(), req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(productResponse)
}

// GetProduct handles retrieving a product by ID
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		helper.LogError("Invalid product ID parameter", err, c.Path(), map[string]interface{}{"id_param": c.Params("id")})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	productResponse, err := h.productUsecase.GetProductByID(uint(id))
	if err != nil {
		helper.LogError("Failed to get product", err, c.Path(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(productResponse)
}

// GetAllProducts handles retrieving all products with pagination
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		offset = 0
	}

	productListResponse, err := h.productUsecase.GetAllProducts(limit, offset)
	if err != nil {
		helper.LogError("Failed to get products", err, c.Path(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(productListResponse)
}

// UpdateProduct handles updating an existing product
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		helper.LogError("Invalid product ID parameter", err, c.Path(), map[string]interface{}{"id_param": c.Params("id")})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		helper.LogError("Failed to parse request body", err, c.Path(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	productResponse, err := h.productUsecase.UpdateProduct(uint(id), &req)
	if err != nil {
		helper.LogError("Failed to update product", err, c.Path(), req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(productResponse)
}

// DeleteProduct handles soft deleting a product
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		helper.LogError("Invalid product ID parameter", err, c.Path(), map[string]interface{}{"id_param": c.Params("id")})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	if err := h.productUsecase.DeleteProduct(uint(id)); err != nil {
		helper.LogError("Failed to delete product", err, c.Path(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
