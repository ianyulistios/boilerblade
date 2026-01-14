package handler_test

import (
	"boilerblade/src/dto"
	"boilerblade/src/handler"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// mockUserUsecase is a mock implementation of UserUsecase for testing
type mockUserUsecase struct {
	users  map[uint]*dto.UserResponse
	nextID uint
}

func newMockUserUsecase() *mockUserUsecase {
	return &mockUserUsecase{
		users:  make(map[uint]*dto.UserResponse),
		nextID: 1,
	}
}

func (m *mockUserUsecase) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check for duplicate email
	for _, user := range m.users {
		if user.Email == req.Email {
			return nil, errors.New("email already exists")
		}
	}

	user := &dto.UserResponse{
		ID:        m.nextID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: "2024-01-01 00:00:00",
		UpdatedAt: "2024-01-01 00:00:00",
	}
	m.users[m.nextID] = user
	m.nextID++
	return user, nil
}

func (m *mockUserUsecase) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserUsecase) GetAllUsers(limit, offset int) (*dto.UserListResponse, error) {
	users := make([]dto.UserResponse, 0)
	for _, user := range m.users {
		users = append(users, *user)
	}

	start := offset
	if start > len(users) {
		start = len(users)
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	if start < end {
		users = users[start:end]
	} else {
		users = []dto.UserResponse{}
	}

	return &dto.UserListResponse{
		Users:      users,
		Total:      int64(len(m.users)),
		Limit:      limit,
		Offset:     offset,
		TotalPages: (len(m.users) + limit - 1) / limit,
	}, nil
}

func (m *mockUserUsecase) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check for duplicate email
		for id2, user2 := range m.users {
			if user2.Email == req.Email && id2 != id {
				return nil, errors.New("email already exists")
			}
		}
		user.Email = req.Email
	}

	return user, nil
}

func (m *mockUserUsecase) DeleteUser(id uint) error {
	_, ok := m.users[id]
	if !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestNewUserHandler(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)

	if userHandler == nil {
		t.Error("NewUserHandler returned nil")
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Post("/users", userHandler.CreateUser)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("Expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
	}
}

func TestUserHandler_CreateUser_InvalidBody(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Post("/users", userHandler.CreateUser)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestUserHandler_CreateUser_ValidationError(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Post("/users", userHandler.CreateUser)

	// Invalid request (missing required fields)
	reqBody := dto.CreateUserRequest{
		Name:  "",              // Missing name
		Email: "invalid-email", // Invalid email
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Get("/users/:id", userHandler.GetUser)

	// Create a user first
	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUsecase.CreateUser(&reqBody)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	// Verify response body
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	if response["data"] == nil {
		t.Error("Response should contain 'data' field")
	}
}

func TestUserHandler_GetUser_InvalidID(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Get("/users/:id", userHandler.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Get("/users/:id", userHandler.GetUser)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Get("/users", userHandler.GetAllUsers)

	// Create some users
	for i := 0; i < 3; i++ {
		reqBody := dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		mockUsecase.CreateUser(&reqBody)
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestUserHandler_GetAllUsers_WithPagination(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Get("/users", userHandler.GetAllUsers)

	// Create some users
	for i := 0; i < 5; i++ {
		reqBody := dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		mockUsecase.CreateUser(&reqBody)
	}

	req := httptest.NewRequest(http.MethodGet, "/users?limit=2&offset=0", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Put("/users/:id", userHandler.UpdateUser)

	// Create a user first
	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUsecase.CreateUser(&reqBody)

	// Update user
	updateReq := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestUserHandler_UpdateUser_NotFound(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Put("/users/:id", userHandler.UpdateUser)

	updateReq := dto.UpdateUserRequest{
		Name: "Updated User",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Delete("/users/:id", userHandler.DeleteUser)

	// Create a user first
	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUsecase.CreateUser(&reqBody)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Handler returns 200 with message, not 204
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	app.Delete("/users/:id", userHandler.DeleteUser)

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}
}

func TestUserHandler_RegisterRoutes(t *testing.T) {
	mockUsecase := newMockUserUsecase()
	userHandler := handler.NewUserHandler(mockUsecase)
	app := setupTestApp()
	apiGroup := app.Group("/api/v1")

	userHandler.RegisterRoutes(apiGroup)

	// Test that routes are registered by making requests
	testCases := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/api/v1/users", http.StatusOK},
		{http.MethodGet, "/api/v1/users/1", http.StatusNotFound},    // No user exists yet
		{http.MethodPost, "/api/v1/users", http.StatusBadRequest},   // Missing body
		{http.MethodPut, "/api/v1/users/1", http.StatusNotFound},    // No user exists
		{http.MethodDelete, "/api/v1/users/1", http.StatusNotFound}, // No user exists
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		if tc.method == http.MethodPost || tc.method == http.MethodPut {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request to %s %s: %v", tc.method, tc.path, err)
		}
		// Just verify route exists (status may vary based on implementation)
		if resp.StatusCode == http.StatusNotFound && tc.path != "/api/v1/users/1" {
			// Route not found means route wasn't registered
			t.Errorf("Route %s %s not registered", tc.method, tc.path)
		}
	}
}
