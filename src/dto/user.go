package dto

// CreateUserRequest represents the request payload for creating a user
// @Description Request payload for creating a new user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100" example:"John Doe"`       // User's full name
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`  // User's email address
	Password string `json:"password" validate:"required,min=6" example:"password123"`        // User's password (min 6 characters)
}

// UpdateUserRequest represents the request payload for updating a user
// @Description Request payload for updating an existing user
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty,min=3,max=100" example:"John Doe Updated"`       // User's full name (optional)
	Email    string `json:"email" validate:"omitempty,email" example:"john.doe.updated@example.com"` // User's email address (optional)
	Password string `json:"password" validate:"omitempty,min=6" example:"newpassword123"`             // User's password (optional, min 6 characters)
}

// UserResponse represents the user response data
// @Description User response data
type UserResponse struct {
	ID        uint   `json:"id" example:"1"`                                    // User ID
	Name      string `json:"name" example:"John Doe"`                           // User's full name
	Email     string `json:"email" example:"john.doe@example.com"`             // User's email address
	CreatedAt string `json:"created_at" example:"2024-01-01 00:00:00"`          // User creation timestamp
	UpdatedAt string `json:"updated_at" example:"2024-01-01 00:00:00"`          // User last update timestamp
}

// UserListResponse represents paginated user list response
// @Description Paginated list of users
type UserListResponse struct {
	Users      []UserResponse `json:"users"`       // List of users
	Total      int64          `json:"total" example:"100"`        // Total number of users
	Limit      int            `json:"limit" example:"10"`        // Number of items per page
	Offset     int            `json:"offset" example:"0"`         // Offset for pagination
	TotalPages int            `json:"total_pages" example:"10"`   // Total number of pages
}
