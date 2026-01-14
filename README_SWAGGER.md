# Swagger API Documentation

Dokumentasi API menggunakan Swagger/OpenAPI untuk Boilerblade API.

## Setup

Swagger sudah terintegrasi dengan aplikasi. Dokumentasi dapat diakses melalui:

```
http://localhost:3000/swagger/index.html
```

## Generate Swagger Documentation

Setelah menambahkan atau mengubah annotations di handler, generate ulang dokumentasi:

```bash
# Install swag CLI (jika belum)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
swag init -g main.go -o docs --parseDependency --parseInternal
```

## Swagger Annotations

### Main API Info

Annotations di `main.go`:

```go
// @title           Boilerblade API
// @version         1.0
// @description     This is a sample server for Boilerblade API.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:3000
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

### Handler Annotations

Contoh annotations untuk endpoint di `src/handler/user.go`:

```go
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
    // ...
}
```

### DTO Annotations

Contoh annotations untuk DTO di `src/dto/user.go`:

```go
// CreateUserRequest represents the request payload for creating a user
// @Description Request payload for creating a new user
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=3,max=100" example:"John Doe"`
    Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
    Password string `json:"password" validate:"required,min=6" example:"password123"`
}
```

## Swagger Annotations Reference

### Endpoint Annotations

- `@Summary` - Short summary of the endpoint
- `@Description` - Detailed description
- `@Tags` - Group endpoints by tag
- `@Accept` - Content types accepted (json, xml, etc.)
- `@Produce` - Content types produced
- `@Param` - Request parameters (path, query, body, header)
- `@Success` - Success response
- `@Failure` - Error response
- `@Security` - Security scheme (BearerAuth, etc.)
- `@Router` - Route path and HTTP method

### Parameter Types

- `path` - URL path parameter (e.g., `/users/{id}`)
- `query` - Query string parameter (e.g., `?limit=10`)
- `body` - Request body
- `header` - HTTP header

### Response Types

- `{object}` - JSON object
- `{array}` - JSON array
- `{string}` - String response
- `{integer}` - Integer response

## Accessing Swagger UI

1. Start the server:
   ```bash
   go run main.go
   ```

2. Open browser and navigate to:
   ```
   http://localhost:3000/swagger/index.html
   ```

3. Swagger UI akan menampilkan semua endpoints yang sudah didokumentasikan

## Testing API from Swagger UI

1. Klik pada endpoint yang ingin ditest
2. Klik "Try it out"
3. Isi parameter yang diperlukan
4. Klik "Execute"
5. Response akan ditampilkan di bawah

## Authentication

Semua endpoints memerlukan JWT authentication. Untuk test di Swagger UI:

1. Klik tombol "Authorize" di bagian atas
2. Masukkan token JWT dengan format: `Bearer <your-token>`
3. Klik "Authorize"
4. Token akan digunakan untuk semua request

## File Structure

```
docs/
├── docs.go          # Auto-generated swagger docs
├── swagger.json      # OpenAPI JSON specification
└── swagger.yaml      # OpenAPI YAML specification
```

## Adding New Endpoints

Untuk menambahkan dokumentasi untuk endpoint baru:

1. Tambahkan annotations di atas handler function
2. Tambahkan annotations di DTO jika diperlukan
3. Generate ulang swagger docs:
   ```bash
   swag init -g main.go -o docs --parseDependency --parseInternal
   ```
4. Restart server untuk melihat perubahan

## Example: Complete Endpoint Documentation

```go
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
    // Implementation
}
```

## Troubleshooting

### Swagger UI tidak muncul

1. Pastikan route swagger sudah ditambahkan di `server/rest.go`:
   ```go
   a.Get("/swagger/*", swagger.HandlerDefault)
   ```

2. Pastikan import docs di `main.go`:
   ```go
   _ "boilerblade/docs" // swagger docs
   ```

3. Pastikan swagger docs sudah di-generate:
   ```bash
   swag init -g main.go -o docs --parseDependency --parseInternal
   ```

### Annotations tidak muncul

1. Pastikan format annotations benar (harus dimulai dengan `//`)
2. Pastikan annotations berada tepat di atas function
3. Generate ulang swagger docs
4. Restart server

### Error saat generate

1. Pastikan semua dependencies terinstall:
   ```bash
   go get -u github.com/swaggo/swag/cmd/swag
   go get -u github.com/swaggo/fiber-swagger
   ```

2. Pastikan swag binary terinstall:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

## Resources

- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [swaggo/swag Documentation](https://github.com/swaggo/swag)
- [fiber-swagger Documentation](https://github.com/swaggo/fiber-swagger)
