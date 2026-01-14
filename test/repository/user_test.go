package repository_test

import (
	"boilerblade/src/model"
	"testing"

	"gorm.io/gorm"
)

// mockUserRepository is a mock implementation for testing without database
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
	m.nextID++
	m.users = append(m.users, user)
	return nil
}

func (m *mockUserRepository) GetByID(id uint) (*model.User, error) {
	for _, user := range m.users {
		if user.ID == id && !user.DeletedAt.Valid {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email && !user.DeletedAt.Valid {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) GetAll(limit, offset int) ([]model.User, error) {
	start := offset
	if start > len(m.users) {
		start = len(m.users)
	}
	end := start + limit
	if end > len(m.users) {
		end = len(m.users)
	}
	if start >= end {
		return []model.User{}, nil
	}
	users := make([]model.User, 0)
	for i := start; i < end; i++ {
		users = append(users, *m.users[i])
	}
	return users, nil
}

func (m *mockUserRepository) Update(user *model.User) error {
	for i, u := range m.users {
		if u.ID == user.ID {
			m.users[i] = user
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Delete(id uint) error {
	for i, user := range m.users {
		if user.ID == id {
			user.DeletedAt = gorm.DeletedAt{Valid: true}
			m.users[i] = user
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Count() (int64, error) {
	count := int64(0)
	for _, user := range m.users {
		if !user.DeletedAt.Valid {
			count++
		}
	}
	return count, nil
}

func TestNewUserRepository(t *testing.T) {
	repo := newMockUserRepository()
	if repo == nil {
		t.Error("newMockUserRepository returned nil")
	}
}

func TestMockUserRepository_Create(t *testing.T) {
	repo := newMockUserRepository()

	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == 0 {
		t.Error("User ID should be set after creation")
	}
}

func TestMockUserRepository_GetByID(t *testing.T) {
	repo := newMockUserRepository()

	// Create a user first
	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	repo.Create(user)

	// Get by ID
	retrieved, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Retrieved user should not be nil")
	}

	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, retrieved.ID)
	}
}

func TestMockUserRepository_GetByID_NotFound(t *testing.T) {
	repo := newMockUserRepository()

	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestMockUserRepository_GetByEmail(t *testing.T) {
	repo := newMockUserRepository()

	// Create a user first
	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	repo.Create(user)

	// Get by email
	retrieved, err := repo.GetByEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Retrieved user should not be nil")
	}

	if retrieved.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", retrieved.Email)
	}
}

func TestMockUserRepository_GetAll(t *testing.T) {
	repo := newMockUserRepository()

	// Create multiple users
	for i := 0; i < 5; i++ {
		user := &model.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		repo.Create(user)
	}

	// Get all with pagination
	users, err := repo.GetAll(10, 0)
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}
}

func TestMockUserRepository_GetAll_WithPagination(t *testing.T) {
	repo := newMockUserRepository()

	// Create 10 users
	for i := 0; i < 10; i++ {
		user := &model.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		repo.Create(user)
	}

	// Get first page
	users, err := repo.GetAll(5, 0)
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}

	// Get second page
	users, err = repo.GetAll(5, 5)
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}
}

func TestMockUserRepository_Update(t *testing.T) {
	repo := newMockUserRepository()

	// Create a user
	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	repo.Create(user)

	// Update user
	user.Name = "Updated User"
	user.Email = "updated@example.com"

	err := repo.Update(user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	updated, err := repo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if updated.Name != "Updated User" {
		t.Errorf("Expected name 'Updated User', got '%s'", updated.Name)
	}
}

func TestMockUserRepository_Delete(t *testing.T) {
	repo := newMockUserRepository()

	// Create a user
	user := &model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	repo.Create(user)

	// Delete user
	err := repo.Delete(user.ID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion (soft delete)
	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Error("User should be deleted (soft delete)")
	}
}

func TestMockUserRepository_Count(t *testing.T) {
	repo := newMockUserRepository()

	// Create multiple users
	for i := 0; i < 5; i++ {
		user := &model.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		repo.Create(user)
	}

	count, err := repo.Count()
	if err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}

	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}
}
