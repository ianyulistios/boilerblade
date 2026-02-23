# Boilerblade

A production-ready Go boilerplate project following Clean Architecture principles, designed to accelerate backend development with built-in code generation, authentication, database migrations, and support for both HTTP REST APIs and AMQP message queues.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture & Design Patterns](#architecture--design-patterns)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Commands](#commands)
- [API Documentation](#api-documentation)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Additional Documentation](#additional-documentation)

## 🎯 Overview

Boilerblade is a comprehensive Go boilerplate that provides:

- **Clean Architecture** implementation with clear separation of concerns
- **Code Generator CLI** for rapid development of CRUD operations
- **JWT Authentication** middleware
- **Database Migrations** via Goose (SQL migrations, versioned)
- **Dual Server Mode**: HTTP REST API and/or AMQP message queue consumers
- **Multiple Database Support**: PostgreSQL, MySQL, SQLite
- **Redis Integration** for caching and session management
- **Swagger/OpenAPI** documentation
- **Comprehensive Testing** structure

## ✨ Features

### Core Features

- 🏗️ **Clean Architecture** - Layered architecture (Handler → DTO → Usecase → Repository → Model)
- 🔐 **JWT Authentication** - Secure token-based authentication middleware
- 🗄️ **Database Support** - PostgreSQL, MySQL, and SQLite via GORM
- 📦 **Code Generator** - CLI tool to generate models, repositories, usecases, handlers, and DTOs
- 🔄 **Goose Migrations** - Versioned SQL migrations (PostgreSQL & MySQL)
- 📡 **Dual Server Mode** - Run HTTP server, AMQP consumers, or both simultaneously
- 📚 **Swagger Documentation** - Auto-generated API documentation
- 🧪 **Test Structure** - Organized test files for handlers, usecases, and repositories
- ⚡ **High Performance** - Built on Fiber framework for fast HTTP handling
- 🔌 **AMQP Support** - RabbitMQ message queue integration for async processing

### Technology Stack

- **Framework**: [Fiber v2](https://github.com/gofiber/fiber) - Fast HTTP framework
- **ORM**: [GORM](https://gorm.io/) - Database ORM
- **Database**: PostgreSQL, MySQL, SQLite
- **Cache**: [Redis](https://redis.io/) - In-memory data store
- **Message Queue**: [RabbitMQ](https://www.rabbitmq.com/) via AMQP
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: [Swagger/OpenAPI](https://swagger.io/)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)

## 🏛️ Architecture & Design Patterns

### Clean Architecture

Boilerblade follows **Clean Architecture** principles with clear layer separation:

```
┌─────────────────────────────────────────┐
│         HTTP Request/Response            │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Handler Layer                   │  ← HTTP handlers, request parsing
│    (src/handler/)                       │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         DTO Layer                       │  ← Data Transfer Objects
│    (src/dto/)                           │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Usecase Layer                   │  ← Business logic
│    (src/usecase/)                       │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Repository Layer                │  ← Data access abstraction
│    (src/repository/)                    │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Model Layer                     │  ← Domain entities
│    (src/model/)                         │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Database (GORM)                 │
└─────────────────────────────────────────┘
```

### Design Patterns Used

1. **Repository Pattern** - Abstracts data access logic
2. **Dependency Injection** - Loose coupling between layers
3. **Factory Pattern** - Connection initialization (Database, Redis, AMQP)
4. **Goose** - Versioned SQL migrations with Up/Down support
5. **Middleware Pattern** - Authentication and request processing
6. **Strategy Pattern** - Multiple database drivers support

### Request Flow

```
HTTP Request
    ↓
Middleware (Auth, CORS, Logger, Recover)
    ↓
Handler (Parse request, validate)
    ↓
DTO (Transform request to domain objects)
    ↓
Usecase (Business logic, validation)
    ↓
Repository (Data access)
    ↓
Model (Domain entity)
    ↓
Database (GORM)
```

## 📁 Project Structure

```
boilerblade/
├── bin/                          # Compiled binaries
│   └── boilerblade.exe          # CLI generator binary
│
├── cmd/                          # Command-line applications
│   └── generate/                 # Code generator CLI
│       └── main.go
│
├── config/                       # Configuration management
│   ├── amqp/                     # AMQP connection management
│   ├── database.go               # Database configuration
│   ├── env.go                    # Environment variables
│   ├── init.go                   # Configuration initialization
│   └── redis.go                  # Redis configuration
│
├── constants/                    # Application constants
│   └── amqp.go
│
├── docs/                         # Swagger documentation (auto-generated)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── helper/                       # Helper utilities
│   └── log.go                    # Logging utilities
│
├── internal/                     # Internal packages
│   └── generator/                # Code generator implementation
│       ├── generator.go
│       └── generator_test.go
│
├── middleware/                   # HTTP middleware
│   └── auth.go                   # JWT authentication middleware
│
├── server/                       # Server setup
│   ├── app.go                    # Application initialization
│   ├── rest.go                   # HTTP routes setup
│   └── amqp.go                   # AMQP consumers setup
│
├── src/                          # Application source code
│   ├── consumer/                 # AMQP message consumers
│   │   └── user.go
│   ├── dto/                      # Data Transfer Objects
│   │   └── user.go
│   ├── handler/                  # HTTP handlers
│   │   └── user.go
│   ├── migration/                # Database migrations
│   │   ├── migration.go
│   │   └── README.md
│   ├── model/                    # Domain models/entities
│   │   ├── user.go
│   │   └── product.go
│   ├── repository/               # Data access layer
│   │   └── user.go
│   └── usecase/                  # Business logic layer
│       └── user.go
│
├── test/                         # Test files
│   ├── handler/
│   ├── repository/
│   ├── usecase/
│   └── README_TEST.md
│
├── .env                          # Environment variables (create from .env.example)
├── .env.example                  # Environment template (app, DB, Redis, AMQP, Goose)
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── main.go                       # Application entry point
├── Makefile                      # Build automation
├── build.sh                      # Linux/Mac build script
├── build.bat                     # Windows build script
├── install.ps1                   # Windows global installer (C:\boilerblade\bin + PATH)
├── install.sh                    # macOS/Linux global installer (~/.local/bin or /usr/local/bin)
├── installer/                    # Native installers: .msi (Windows), .deb (Linux), .pkg (macOS)
│
└── README files:
    ├── README.md                 # This file (main documentation)
    ├── README_BUILD.md           # Build and installation guide
    ├── README_CLI.md             # CLI generator usage
    ├── README_CRUD_USER.md       # CRUD implementation example
    └── README_SWAGGER.md         # Swagger documentation guide
```

## 🚀 Getting Started

### Prerequisites

- **Go 1.24+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL/MySQL** (optional) - For database
- **Redis** (optional) - For caching
- **RabbitMQ** (optional) - For message queue

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd boilerblade
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your credentials (DB, Redis, AMQP, APP_KEY, etc.)
   ```

4. **Build the CLI tool** (optional)
   ```bash
   # Using Makefile (Linux/Mac)
   make build
   
   # Using build script
   ./build.sh        # Linux/Mac
   build.bat         # Windows
   
   # Or manually
   go build -o bin/boilerblade ./cmd/cli
   ```

   For detailed installation instructions, see [README_BUILD.md](README_BUILD.md).

   **Install globally (like Composer):** Run `boilerblade` from CMD, Git Bash, or Terminal from any directory. See [README_INSTALL.md](README_INSTALL.md) for Windows (`install.ps1`) and macOS/Linux (`install.sh`).

5. **Run the application**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:3000` (default port).

### Docker Setup (Alternative)

You can also run the entire stack using Docker Compose:

```bash
# Start all services (app, PostgreSQL, Redis, RabbitMQ)
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down
```

For detailed Docker setup instructions, see [README_DOCKER.md](README_DOCKER.md).

## ⚙️ Configuration

Configuration is managed through environment variables. Copy `.env.example` to `.env` and customize:

### Server Configuration

```env
MODE=development                    # development, production
FIBER_PORT=3000                     # HTTP server port
FIBER_APP_NAME=boilerblade          # Application name
APP_KEY=your-secret-key-here        # JWT secret key (change in production!)
SERVER_MODE=both                    # http, amqp, or both
```

### Database Configuration

```env
ENABLE_DB=true                      # Enable/disable database
DB_TYPE=postgres                    # postgres, mysql, sqlite
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=boilerblade
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=10
DB_MAX_LIFETIME_CONNS=10
```

### Redis Configuration

```env
ENABLE_REDIS=true                   # Enable/disable Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### AMQP Configuration

```env
ENABLE_AMQP=true                    # Enable/disable AMQP
AMQP_HOST=localhost
AMQP_PORT=5672
AMQP_USER=guest
AMQP_PASSWORD=guest
```

### Connection Flags

You can disable specific connections by setting flags to `false`:
- `ENABLE_DB=false` - Disable database connection
- `ENABLE_REDIS=false` - Disable Redis connection
- `ENABLE_AMQP=false` - Disable AMQP connection

## 🛠️ Commands

### Running the Application

```bash
# Development mode
go run main.go

# Build and run
go build -o boilerblade.exe
./boilerblade.exe
```

### CLI Commands

Boilerblade provides a powerful CLI tool similar to Laravel's Artisan:

#### Create New Project

```bash
# Create a new Boilerblade project
boilerblade new my-api
```

This will scaffold a complete Boilerblade project in the `my-api` directory with all necessary files and configurations.

#### Generate Code (Make Command)

The `make` command helps you quickly scaffold CRUD operations:

```bash
# Generate all layers for an entity
boilerblade make all -name=Product -fields="Name:string:required,Price:float64:required"

# Generate specific layer
boilerblade make model -name=Order
boilerblade make repository -name=Order
boilerblade make usecase -name=Order
boilerblade make handler -name=Order
boilerblade make dto -name=Order
```

**Available Resources:**
- `model` - Domain model/entity
- `repository` - Data access layer
- `usecase` - Business logic layer
- `handler` - HTTP handlers
- `dto` - Data Transfer Objects
- `all` - Generate all layers at once

**Field Format:**
```
FieldName:FieldType:GormTag
```

**Examples:**
```bash
# Simple entity
boilerblade make all -name=Category -fields="Name:string:required"

# Complex entity
boilerblade make all -name=Product -fields="Name:string:required,Price:float64:required,Stock:int:required,Description:string"
```

#### Other Commands

```bash
# Show help
boilerblade help

# Show version
boilerblade version
```

**Backward Compatibility:**
The old flag-based format still works for backward compatibility:
```bash
boilerblade -layer=all -name=Product -fields="Name:string:required"
```

For detailed CLI usage, see [README_CLI.md](README_CLI.md).

### Build Commands

```bash
# Build binary
make build

# Install to GOPATH/bin
make install

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

## 📚 API Documentation

### Swagger UI

Once the server is running, access Swagger documentation at:

```
http://localhost:3000/swagger/index.html
```

### Generate Swagger Docs

After adding or modifying API endpoints:

```bash
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g main.go -o docs --parseDependency --parseInternal
```

For detailed Swagger setup, see [README_SWAGGER.md](README_SWAGGER.md).

### API Endpoints

All API endpoints are prefixed with `/api/v1` and require JWT authentication.

**Example User Endpoints:**
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users` - Get all users (with pagination)
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

For complete CRUD implementation example, see [README_CRUD_USER.md](README_CRUD_USER.md).

### Authentication

All endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 🔄 Development Workflow

### Creating a New Feature

1. **Generate code using CLI**
   ```bash
   boilerblade make all -name=Product -fields="Name:string:required,Price:float64:required"
   ```

2. **Add a Goose migration** (if the feature adds new tables)
   Add SQL file(s) under `src/migration/migrations/`, e.g. `00003_create_orders_table.postgres.sql` and `00003_create_orders_table.mysql.sql`. See [src/migration/README.md](src/migration/README.md).

3. **Register routes**
   Edit `server/rest.go`:
   ```go
   // Initialize dependencies
   productRepo := repository.NewProductRepository(a.Config.Database)
   productUsecase := usecase.NewProductUsecase(productRepo)
   productHandler := handler.NewProductHandler(productUsecase)
   
   // Register routes
   productHandler.RegisterRoutes(apiV1Group)
   ```

4. **Add Swagger annotations**
   Add annotations to handler methods (see [README_SWAGGER.md](README_SWAGGER.md))

5. **Generate Swagger docs**
   ```bash
   swag init -g main.go -o docs --parseDependency --parseInternal
   ```

6. **Test the endpoints**
   - Use Swagger UI: `http://localhost:3000/swagger/index.html`
   - Or use cURL/Postman

### Server Modes

The application supports three server modes:

1. **HTTP Only** (`SERVER_MODE=http`)
   - Runs only HTTP REST API server

2. **AMQP Only** (`SERVER_MODE=amqp`)
   - Runs only AMQP message queue consumers

3. **Both** (`SERVER_MODE=both` - default)
   - Runs both HTTP server and AMQP consumers concurrently

## 🧪 Testing

The project includes a structured testing setup:

```
test/
├── handler/          # HTTP handler tests
├── repository/       # Repository/data access tests
├── usecase/          # Business logic tests
└── README_TEST.md    # Testing documentation
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./test/handler/...
```

For testing guidelines, see [test/README_TEST.md](test/README_TEST.md).

## 📖 Additional Documentation

- **[README_BUILD.md](README_BUILD.md)** - Detailed build and installation instructions
- **[README_INSTALL.md](README_INSTALL.md)** - Global installer (script: install.ps1 / install.sh; native: .msi, .deb, .pkg — see [installer/README.md](installer/README.md))
- **[README_CLI.md](README_CLI.md)** - CLI code generator usage guide
- **[README_CRUD_USER.md](README_CRUD_USER.md)** - Complete CRUD implementation example
- **[README_SWAGGER.md](README_SWAGGER.md)** - Swagger/OpenAPI documentation guide
- **[README_DOCKER.md](README_DOCKER.md)** - Docker and Docker Compose setup guide
- **[src/migration/README.md](src/migration/README.md)** - Database migration system documentation
- **[test/README_TEST.md](test/README_TEST.md)** - Testing guidelines and examples

## 🎯 Key Concepts

### Clean Architecture Layers

1. **Handler** - HTTP request/response handling, validation
2. **DTO** - Data transformation between layers
3. **Usecase** - Business logic, orchestration
4. **Repository** - Data access abstraction
5. **Model** - Domain entities

### Dependency Flow

Dependencies flow inward:
- Handler depends on Usecase
- Usecase depends on Repository
- Repository depends on Model
- Model has no dependencies

### Benefits

- **Testability** - Each layer can be tested independently
- **Maintainability** - Clear separation of concerns
- **Flexibility** - Easy to swap implementations
- **Scalability** - Easy to add new features

## 🔒 Security

- JWT-based authentication
- Password hashing (implement bcrypt/argon2 for production)
- CORS configuration
- Input validation
- SQL injection protection (via GORM)
- Error handling without exposing internals

## 🚧 Roadmap

- [ ] Add more authentication providers (OAuth2, etc.)
- [ ] GraphQL support
- [ ] WebSocket support
- [ ] Rate limiting middleware
- [ ] Request/response logging
- [ ] Health check endpoints
- [ ] Metrics and monitoring
- [x] Docker and Docker Compose setup
- [ ] Kubernetes deployment examples

## 📝 License

[Add your license here]

## 🤝 Contributing

[Add contribution guidelines here]

## 📧 Support

[Add support/contact information here]

---

**Built with ❤️ using Go and Clean Architecture principles**
