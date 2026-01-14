# CRUD User Implementation

Contoh implementasi CRUD User dengan struktur Clean Architecture:
**Handler → DTO → Usecase → Model → Repository**

## Struktur File

```
src/
├── model/
│   └── user.go              # Entity/Model User
├── repository/
│   └── user_repository.go   # Data access layer
├── dto/
│   └── user_dto.go          # Data Transfer Objects
├── usecase/
│   └── user_usecase.go      # Business logic layer
└── handler/
    └── user_handler.go      # HTTP handlers
```

## Flow Request

```
HTTP Request
    ↓
Handler (user_handler.go)
    ↓ Parse & Validate
DTO (user_dto.go)
    ↓
Usecase (user_usecase.go)
    ↓ Business Logic
Repository (user_repository.go)
    ↓
Model (user.go)
    ↓
Database (GORM)
```

## API Endpoints

### 1. Create User
```http
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

### 2. Get User by ID
```http
GET /api/v1/users/:id
```

### 3. Get All Users (with pagination)
```http
GET /api/v1/users?limit=10&offset=0
```

### 4. Update User
```http
PUT /api/v1/users/:id
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

### 5. Delete User
```http
DELETE /api/v1/users/:id
```

## Setup Database Migration

Tambahkan di `main.go` atau buat migration script:

```go
import "boilerblade/src/migration"

// After app initialization
if err := migration.CreateUsersTable(app.Config.Database); err != nil {
    log.Fatal("Failed to migrate database:", err)
}
```

## Testing dengan cURL

### Create User
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get All Users
```bash
curl http://localhost:3000/api/v1/users?limit=10&offset=0
```

### Get User by ID
```bash
curl http://localhost:3000/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:3000/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:3000/api/v1/users/1
```

## Catatan Penting

1. **Password Hashing**: Saat ini password disimpan plain text. Untuk production, gunakan bcrypt atau argon2.
2. **JWT Authentication**: Routes sudah dilindungi dengan JWT middleware.
3. **Validation**: Menggunakan `go-playground/validator` untuk validasi request.
4. **Soft Delete**: User menggunakan soft delete (tidak benar-benar dihapus dari database).
