package usecase

import (
	"boilerblade/src/dto"
	"boilerblade/src/model"
	"boilerblade/src/repository"
	"errors"
	"math"
)

// UserUsecase defines the interface for user business logic
type UserUsecase interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	GetAllUsers(limit, offset int) (*dto.UserListResponse, error)
	UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
}

// userUsecase implements UserUsecase interface
type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new user usecase instance
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (uc *userUsecase) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existingUser, _ := uc.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create user model
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // In production, hash the password
	}

	// Save to database
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Return response DTO
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetUserByID retrieves a user by ID
func (uc *userUsecase) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetAllUsers retrieves all users with pagination
func (uc *userUsecase) GetAllUsers(limit, offset int) (*dto.UserListResponse, error) {
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

	// Get users and total count
	users, err := uc.userRepo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.userRepo.Count()
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser updates an existing user
func (uc *userUsecase) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email already exists
		existingUser, _ := uc.userRepo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password // In production, hash the password
	}

	// Save updates
	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Return response DTO
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// DeleteUser deletes a user
func (uc *userUsecase) DeleteUser(id uint) error {
	// Check if user exists
	_, err := uc.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	// Delete user
	return uc.userRepo.Delete(id)
}
