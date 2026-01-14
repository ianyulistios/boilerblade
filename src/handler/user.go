package handler

import (
	"boilerblade/helper"
	"boilerblade/src/dto"
	"boilerblade/src/usecase"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userUsecase usecase.UserUsecase
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		validator:   validator.New(),
	}
}

// RegisterRoutes registers all user routes to the router group
func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/users", h.GetAllUsers)
	router.Get("/users/:id", h.GetUser)
	router.Post("/users", h.CreateUser)
	router.Put("/users/:id", h.UpdateUser)
	router.Delete("/users/:id", h.DeleteUser)
}

// CreateUser handles POST /users
// @Summary      Create a new user
// @Description  Create a new user with name, email, and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserRequest  true  "User data"
// @Success      201   {object}  map[string]interface{}  "User created successfully"
// @Failure      400   {object}  map[string]interface{}  "Invalid request body or validation failed"
// @Failure      500   {object}  map[string]interface{}  "Internal server error"
// @Security     BearerAuth
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		helper.LogError("Failed to parse request body", err, c.Path(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		helper.LogError("Validation failed", err, c.Path(), req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Call usecase
	user, err := h.userUsecase.CreateUser(&req)
	if err != nil {
		helper.LogError("Failed to create user", err, c.Path(), req)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	helper.LogInfo("User created successfully", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}

// GetUser handles GET /users/:id
// @Summary      Get user by ID
// @Description  Get a user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}  "User data"
// @Failure      400  {object}  map[string]interface{}  "Invalid user ID"
// @Failure      404  {object}  map[string]interface{}  "User not found"
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Get ID from params
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Call usecase
	user, err := h.userUsecase.GetUserByID(uint(id))
	if err != nil {
		helper.LogError("Failed to get user", err, c.Path(), map[string]interface{}{"id": id})
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": user,
	})
}

// GetAllUsers handles GET /users
// @Summary      Get all users
// @Description  Get all users with pagination support
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Limit number of results (default: 10, max: 100)"
// @Param        offset  query     int  false  "Offset for pagination (default: 0)"
// @Success      200     {object}  map[string]interface{}  "List of users"
// @Failure      500     {object}  map[string]interface{}  "Internal server error"
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Get pagination parameters
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Call usecase
	users, err := h.userUsecase.GetAllUsers(limit, offset)
	if err != nil {
		helper.LogError("Failed to get users", err, c.Path(), nil)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": users,
	})
}

// UpdateUser handles PUT /users/:id
// @Summary      Update user
// @Description  Update an existing user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int                   true  "User ID"
// @Param        user  body      dto.UpdateUserRequest  true  "User data to update"
// @Success      200   {object}  map[string]interface{}  "User updated successfully"
// @Failure      400   {object}  map[string]interface{}  "Invalid request body or validation failed"
// @Failure      404   {object}  map[string]interface{}  "User not found"
// @Failure      500   {object}  map[string]interface{}  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Get ID from params
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req dto.UpdateUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		helper.LogError("Failed to parse request body", err, c.Path(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		helper.LogError("Validation failed", err, c.Path(), req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Call usecase
	user, err := h.userUsecase.UpdateUser(uint(id), &req)
	if err != nil {
		helper.LogError("Failed to update user", err, c.Path(), map[string]interface{}{"id": id})
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	helper.LogInfo("User updated successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser handles DELETE /users/:id
// @Summary      Delete user
// @Description  Delete a user by ID (soft delete)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}  "User deleted successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid user ID"
// @Failure      404  {object}  map[string]interface{}  "User not found"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Get ID from params
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Call usecase
	if err := h.userUsecase.DeleteUser(uint(id)); err != nil {
		helper.LogError("Failed to delete user", err, c.Path(), map[string]interface{}{"id": id})
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	helper.LogInfo("User deleted successfully", map[string]interface{}{
		"user_id": id,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
