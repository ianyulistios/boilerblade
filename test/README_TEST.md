# User CRUD Test Suite

Test suite lengkap untuk User CRUD operations di semua layer (Repository, Usecase, Handler).

## Struktur Test

Semua test files berada di folder `test/` dengan struktur:
```
test/
├── repository/
│   └── user_test.go
├── usecase/
│   └── user_test.go
├── handler/
│   └── user_test.go
└── README_TEST.md
```

### 1. Repository Tests (`test/repository/user_test.go`)

Test untuk data access layer menggunakan mock repository:

- `TestNewUserRepository` - Test repository initialization
- `TestMockUserRepository_Create` - Test user creation
- `TestMockUserRepository_GetByID` - Test get user by ID
- `TestMockUserRepository_GetByID_NotFound` - Test error handling untuk non-existent user
- `TestMockUserRepository_GetByEmail` - Test get user by email
- `TestMockUserRepository_GetAll` - Test get all users
- `TestMockUserRepository_GetAll_WithPagination` - Test pagination
- `TestMockUserRepository_Update` - Test user update
- `TestMockUserRepository_Delete` - Test soft delete
- `TestMockUserRepository_Count` - Test count users

**Note**: Menggunakan mock repository karena SQLite memerlukan CGO yang mungkin tidak tersedia di semua environment.

### 2. Usecase Tests (`test/usecase/user_test.go`)

Test untuk business logic layer menggunakan mock repository:

- `TestNewUserUsecase` - Test usecase initialization
- `TestUserUsecase_CreateUser` - Test create user
- `TestUserUsecase_CreateUser_DuplicateEmail` - Test duplicate email validation
- `TestUserUsecase_GetUserByID` - Test get user by ID
- `TestUserUsecase_GetUserByID_NotFound` - Test error handling
- `TestUserUsecase_GetAllUsers` - Test get all users
- `TestUserUsecase_GetAllUsers_WithPagination` - Test pagination
- `TestUserUsecase_GetAllUsers_InvalidLimit` - Test pagination validation
- `TestUserUsecase_UpdateUser` - Test update user
- `TestUserUsecase_UpdateUser_NotFound` - Test error handling
- `TestUserUsecase_UpdateUser_DuplicateEmail` - Test duplicate email on update
- `TestUserUsecase_DeleteUser` - Test delete user
- `TestUserUsecase_DeleteUser_NotFound` - Test error handling

### 3. Handler Tests (`test/handler/user_test.go`)

Test untuk HTTP handler layer menggunakan mock usecase dan Fiber test utilities:

- `TestNewUserHandler` - Test handler initialization
- `TestUserHandler_CreateUser` - Test POST /users
- `TestUserHandler_CreateUser_InvalidBody` - Test invalid request body
- `TestUserHandler_CreateUser_ValidationError` - Test validation errors
- `TestUserHandler_GetUser` - Test GET /users/:id
- `TestUserHandler_GetUser_InvalidID` - Test invalid ID parameter
- `TestUserHandler_GetUser_NotFound` - Test 404 error
- `TestUserHandler_GetAllUsers` - Test GET /users
- `TestUserHandler_GetAllUsers_WithPagination` - Test pagination query params
- `TestUserHandler_UpdateUser` - Test PUT /users/:id
- `TestUserHandler_UpdateUser_NotFound` - Test 404 error
- `TestUserHandler_DeleteUser` - Test DELETE /users/:id
- `TestUserHandler_DeleteUser_NotFound` - Test 404 error
- `TestUserHandler_RegisterRoutes` - Test route registration

## Menjalankan Tests

### Run All Tests
```bash
go test ./test/... -v
```

### Run Specific Package
```bash
# Repository tests
go test ./test/repository -v

# Usecase tests
go test ./test/usecase -v

# Handler tests
go test ./test/handler -v
```

### Run Specific Test
```bash
go test ./test/handler -v -run TestUserHandler_CreateUser
```

### Run with Coverage
```bash
go test ./test/... -cover
```

### Generate Coverage Report
```bash
go test ./test/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Coverage

### Repository Layer
- ✅ Create user
- ✅ Get user by ID
- ✅ Get user by email
- ✅ Get all users with pagination
- ✅ Update user
- ✅ Soft delete user
- ✅ Count users
- ✅ Error handling (not found)

### Usecase Layer
- ✅ Create user with validation
- ✅ Duplicate email validation
- ✅ Get user by ID
- ✅ Get all users with pagination
- ✅ Pagination validation (limit, offset)
- ✅ Update user
- ✅ Duplicate email on update
- ✅ Delete user
- ✅ Error handling

### Handler Layer
- ✅ Create user endpoint
- ✅ Request body parsing
- ✅ Input validation
- ✅ Get user endpoint
- ✅ Invalid ID parameter handling
- ✅ Get all users with pagination
- ✅ Update user endpoint
- ✅ Delete user endpoint
- ✅ Error responses (400, 404, 500)
- ✅ Route registration

## Mock Implementations

### Mock User Repository
Mock repository yang mengimplementasikan `UserRepository` interface untuk testing tanpa database:

```go
type mockUserRepository struct {
    users  []*model.User
    nextID uint
}
```

### Mock User Usecase
Mock usecase yang mengimplementasikan `UserUsecase` interface untuk testing handler:

```go
type mockUserUsecase struct {
    users map[uint]*dto.UserResponse
    nextID uint
}
```

## Best Practices

1. **Isolated Tests**: Setiap test independen dan tidak bergantung pada test lain
2. **Mock Dependencies**: Menggunakan mock untuk mengisolasi layer yang ditest
3. **Error Cases**: Test mencakup error cases dan edge cases
4. **Validation**: Test validasi input dan business rules
5. **HTTP Testing**: Menggunakan Fiber test utilities untuk handler testing

## Menambahkan Test Baru

### Contoh Test Structure

```go
func TestUserHandler_NewFeature(t *testing.T) {
    // Setup
    mockUsecase := newMockUserUsecase()
    handler := NewUserHandler(mockUsecase)
    app := setupTestApp()
    app.Post("/users", handler.CreateUser)

    // Test
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

    // Verify
    if resp.StatusCode != fiber.StatusCreated {
        t.Errorf("Expected status %d, got %d", fiber.StatusCreated, resp.StatusCode)
    }
}
```

## Notes

- Repository tests menggunakan mock karena SQLite memerlukan CGO
- Handler tests menggunakan Fiber test utilities (`app.Test()`)
- Semua tests menggunakan mock untuk dependency isolation
- Test coverage mencakup success cases dan error cases
