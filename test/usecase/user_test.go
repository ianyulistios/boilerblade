package usecase_test

import (
	"boilerblade/src/dto"
	"boilerblade/src/model"
	"boilerblade/src/usecase"
	"testing"
	"time"

	"gorm.io/gorm"
)

// mockUserRepository is a mock implementation of UserRepository for testing
type mockUserRepository struct {
	users  []*model.User
	nextID uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make([]*model.User, 0),
		nextID: 1,
	}
}

func (m *mockUserRepository) Create(user *model.User) error {
	user.ID = m.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.nextID++
	m.users = append(m.users, user)
	return nil
}

func (m *mockUserRepository) GetByID(id uint) (*model.User, error) {
	for _, user := range m.users {
		if user.ID == id && user.DeletedAt.Time.IsZero() {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email && user.DeletedAt.Time.IsZero() {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetAll(limit, offset int) ([]model.User, error) {
	activeUsers := make([]model.User, 0)
	for _, user := range m.users {
		if user.DeletedAt.Time.IsZero() {
			activeUsers = append(activeUsers, *user)
		}
	}

	start := offset
	if start > len(activeUsers) {
		start = len(activeUsers)
	}

	end := start + limit
	if end > len(activeUsers) {
		end = len(activeUsers)
	}

	if start >= end {
		return []model.User{}, nil
	}

	return activeUsers[start:end], nil
}

func (m *mockUserRepository) Update(user *model.User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			user.UpdatedAt = time.Now()
			m.users[i] = user
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Delete(id uint) error {
	user, err := m.GetByID(id)
	if err != nil {
		return err
	}
	user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return nil
}

func (m *mockUserRepository) Count() (int64, error) {
	count := int64(0)
	for _, user := range m.users {
		if user.DeletedAt.Time.IsZero() {
			count++
		}
	}
	return count, nil
}

func TestNewUserUsecase(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	if uc == nil {
		t.Error("NewUserUsecase returned nil")
	}
}

func TestUserUsecase_CreateUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	req := &dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := uc.CreateUser(req)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if resp == nil {
		t.Fatal("Response should not be nil")
	}

	if resp.ID == 0 {
		t.Error("User ID should be set")
	}

	if resp.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", resp.Name)
	}

	if resp.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", resp.Email)
	}
}

func TestUserUsecase_CreateUser_DuplicateEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Create first user
	req1 := &dto.CreateUserRequest{
		Name:     "Test User 1",
		Email:    "test@example.com",
		Password: "password123",
	}
	uc.CreateUser(req1)

	// Try to create user with same email
	req2 := &dto.CreateUserRequest{
		Name:     "Test User 2",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err := uc.CreateUser(req2)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}

	if err.Error() != "email already exists" {
		t.Errorf("Expected 'email already exists' error, got '%s'", err.Error())
	}
}

func TestUserUsecase_GetUserByID(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Create a user first
	req := &dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, _ := uc.CreateUser(req)

	// Get user by ID
	resp, err := uc.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if resp == nil {
		t.Fatal("Response should not be nil")
	}

	if resp.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, resp.ID)
	}
}

func TestUserUsecase_GetUserByID_NotFound(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	_, err := uc.GetUserByID(999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got '%s'", err.Error())
	}
}

func TestUserUsecase_GetAllUsers(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Create multiple users with unique emails
	for i := 0; i < 5; i++ {
		req := &dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		// Create user directly in mock to avoid email validation
		user := &model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}
		mockRepo.Create(user)
	}

	resp, err := uc.GetAllUsers(10, 0)
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if resp == nil {
		t.Fatal("Response should not be nil")
	}

	if len(resp.Users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(resp.Users))
	}

	if resp.Total != 5 {
		t.Errorf("Expected total 5, got %d", resp.Total)
	}
}

func TestUserUsecase_GetAllUsers_WithPagination(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Create 10 users with unique emails by modifying email
	for i := 0; i < 10; i++ {
		user := &model.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		// Create user directly in mock - mock will allow duplicates for testing
		// We'll modify the mock to allow this
		user.ID = mockRepo.nextID
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		mockRepo.nextID++
		mockRepo.users = append(mockRepo.users, user)
	}

	// Get first page
	resp, err := uc.GetAllUsers(5, 0)
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	if len(resp.Users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(resp.Users))
	}

	if resp.TotalPages != 2 {
		t.Errorf("Expected 2 total pages, got %d", resp.TotalPages)
	}
}

func TestUserUsecase_GetAllUsers_InvalidLimit(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Test with invalid limit (should default to 10)
	resp, err := uc.GetAllUsers(-1, 0)
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	if resp.Limit != 10 {
		t.Errorf("Expected limit to be defaulted to 10, got %d", resp.Limit)
	}

	// Test with limit > 100 (should be capped at 100)
	resp, err = uc.GetAllUsers(200, 0)
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	if resp.Limit != 100 {
		t.Errorf("Expected limit to be capped at 100, got %d", resp.Limit)
	}
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Create a user first
	req := &dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, _ := uc.CreateUser(req)

	// Update user
	updateReq := &dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	resp, err := uc.UpdateUser(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	if resp.Name != "Updated User" {
		t.Errorf("Expected name 'Updated User', got '%s'", resp.Name)
	}

	if resp.Email != "updated@example.com" {
		t.Errorf("Expected email 'updated@example.com', got '%s'", resp.Email)
	}
}

func TestUserUsecase_UpdateUser_NotFound(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	updateReq := &dto.UpdateUserRequest{
		Name: "Updated User",
	}

	_, err := uc.UpdateUser(999, updateReq)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got '%s'", err.Error())
	}
}

func TestUserUsecase_UpdateUser_DuplicateEmail(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Create two users
	req1 := &dto.CreateUserRequest{
		Name:     "Test User 1",
		Email:    "test1@example.com",
		Password: "password123",
	}
	req2 := &dto.CreateUserRequest{
		Name:     "Test User 2",
		Email:    "test2@example.com",
		Password: "password123",
	}
	user1, _ := uc.CreateUser(req1)
	uc.CreateUser(req2)

	// Try to update user2 with user1's email
	updateReq := &dto.UpdateUserRequest{
		Email: "test1@example.com",
	}

	_, err := uc.UpdateUser(user1.ID+1, updateReq)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}

	if err.Error() != "email already exists" {
		t.Errorf("Expected 'email already exists' error, got '%s'", err.Error())
	}
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// Create a user
	req := &dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, _ := uc.CreateUser(req)

	// Delete user
	err := uc.DeleteUser(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion
	_, err = uc.GetUserByID(created.ID)
	if err == nil {
		t.Error("User should be deleted")
	}
}

func TestUserUsecase_DeleteUser_NotFound(t *testing.T) {
	mockRepo := newMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	err := uc.DeleteUser(999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got '%s'", err.Error())
	}
}
