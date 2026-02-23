# Boilerblade CLI Generator

CLI tool untuk generate boilerplate code untuk Model, Repository, Usecase, Handler, dan DTO, serta membuat project baru.

## Installation

### Option 1: Install Binary (Recommended)

**Install via Go:**
```bash
go install ./cmd/cli
```

Setelah install, pastikan `$GOPATH/bin` ada di PATH, lalu gunakan:
```bash
boilerblade new my-api
boilerblade make all -name=Product
```

**Build Binary:**
```bash
# Using Makefile
make build
make install

# Or using build script
./build.sh        # Linux/Mac
build.bat         # Windows

# Or manually
go build -o bin/boilerblade ./cmd/cli
```

Lihat [README_BUILD.md](README_BUILD.md) untuk detail instalasi.

### Option 2: Run Directly (Development)

Tidak perlu install, langsung jalankan dengan `go run`:

```bash
go run cmd/cli/main.go new my-api
go run cmd/cli/main.go make model -name=Product
```

## Usage

### Create New Project

```bash
boilerblade new <project-name>
```

Contoh:
```bash
boilerblade new my-api
```

Ini akan membuat project Boilerblade baru di direktori `my-api` dengan semua file dan konfigurasi yang diperlukan.

### Generate Code (Make Command)

**Using Binary:**
```bash
boilerblade make <resource> -name=<EntityName> [-fields=<fields>]
```

**Using Go Run:**
```bash
go run cmd/cli/main.go make <resource> -name=<EntityName> [-fields=<fields>]
```

### Backward Compatibility (Old Format)

Format lama masih didukung untuk backward compatibility:
```bash
boilerblade -layer=<layer> -name=<EntityName> [-fields=<fields>]
```

### Layers

- `model` - Generate model file
- `repository` - Generate repository file
- `usecase` - Generate usecase file
- `handler` - Generate handler file
- `dto` - Generate DTO file
- `consumer` - Generate RabbitMQ consumer from name and optional title (general-purpose)
- `migration` - Create Goose SQL migration (postgres + mysql placeholder files)
- `all` - Generate all layers (model, dto, repository, usecase, handler)

### Options

- `-layer` (required) - Layer yang ingin di-generate
- `-name` (required) - Nama entity (e.g., Product, Order, Category)
- `-fields` (optional) - Fields untuk model (format: `Name:Type:Tag,Field2:Type2:Tag2`)
- `-help` - Show help message

### Field Format

Format: `FieldName:FieldType:GormTag`

Contoh:
- `Name:string:required` - Field Name dengan type string dan tag "required"
- `Price:float64:required` - Field Price dengan type float64 dan tag "required"
- `Stock:int:required` - Field Stock dengan type int dan tag "required"
- `Description:string` - Field Description dengan type string tanpa tag khusus

## Examples

### Create New Project

```bash
boilerblade new my-api
```

### Generate Model Only

```bash
boilerblade make model -name=Product -fields="Name:string:required,Price:float64:required,Stock:int:required"
```

### Generate Repository Only

```bash
boilerblade make repository -name=Product
```

### Generate All Layers

```bash
boilerblade make all -name=Order -fields="OrderNumber:string:required,Total:float64:required,Status:string:required"
```

### Generate Handler Only

```bash
boilerblade make handler -name=Product
```

### Generate DTO Only

```bash
boilerblade make dto -name=Product -fields="Name:string:required,Price:float64:required"
```

### Generate RabbitMQ Consumer (general-purpose)

Generates a RabbitMQ consumer from **name** and optional **title** only (no entity/fields). Use any name you want (e.g. OrderEvents, payment, inventory).

```bash
boilerblade make consumer -name=OrderEvents
boilerblade make consumer -name=payment -title="Payment"
boilerblade make consumer -name=order_events -title="Order Events"
```

Generated files:
- `constants/amqp_<identifier>.go` – exchange name, created/updated queue names, routing keys, intervals (identifier is snake_case from name, e.g. order_events)
- `src/consumer/<identifier>.go` – consumer struct, `ProcessCreated`, `ProcessUpdated`, generic message handlers (payload as `map[string]interface{}`); add your logic in `handleCreatedMessage` / `handleUpdatedMessage`

Register the consumer in `server/amqp.go` and add your business logic in the handler TODOs.

### Generate Goose Migration

Creates a new Goose migration with placeholder Up/Down SQL for both PostgreSQL and MySQL.

```bash
boilerblade make migration -name=add_orders_table
boilerblade make migration -name=add_index_to_users
```

Generated files (in `src/migration/migrations/`):
- `<timestamp>_<name>.postgres.sql`
- `<timestamp>_<name>.mysql.sql`

Edit the files to add your SQL, then run the app (migrations run on startup) or use the Goose CLI.

### Show Help

```bash
boilerblade help
```

### Show Version

```bash
boilerblade version
```

## Generated Files

### Model (`src/model/<entity>.go`)
- Struct dengan GORM tags
- TableName() method
- Standard fields: ID, CreatedAt, UpdatedAt, DeletedAt

### Repository (`src/repository/<entity>.go`)
- Interface dengan CRUD methods
- Implementation dengan GORM
- Methods: Create, GetByID, GetAll, Update, Delete, Count

### Usecase (`src/usecase/<entity>.go`)
- Interface dengan business logic methods
- Implementation dengan repository dependency
- Methods: Create, GetByID, GetAll, Update, Delete
- Pagination support

### Handler (`src/handler/<entity>.go`)
- HTTP handlers untuk Fiber
- RegisterRoutes() method
- CRUD endpoints: GET, POST, PUT, DELETE
- Error handling dengan helper

### Consumer (`consumer` – RabbitMQ, general-purpose)
- **-name** (required): consumer name (e.g. OrderEvents, payment, order_events). Normalized to identifier (snake_case) and struct name (PascalCase).
- **-title** (optional): human-readable title (e.g. "Order Events"); used in comments and logs.
- **constants/amqp_&lt;identifier&gt;.go** – Exchange, queue names, routing keys, intervals.
- **src/consumer/&lt;identifier&gt;.go** – Consumer struct (no usecase/DTO), `ProcessCreated`, `ProcessUpdated`, handlers with generic payload; add your logic and register in `server/amqp.go`.

### DTO (`src/dto/<entity>.go`)
- CreateRequest struct
- UpdateRequest struct
- Response struct
- ListResponse struct dengan pagination
- ToResponse() converter function

## After Generation

Setelah generate, Anda perlu:

1. **Update Model Fields** - Sesuaikan fields di model sesuai kebutuhan
2. **Update DTO Fields** - Sesuaikan fields di DTO sesuai kebutuhan
3. **Update Usecase Logic** - Implement mapping dari DTO ke Model
4. **Register Routes** - Tambahkan handler ke `server/rest.go`:
   ```go
   productRepo := repository.NewProductRepository(a.Config.Database)
   productUsecase := usecase.NewProductUsecase(productRepo)
   productHandler := handler.NewProductHandler(productUsecase)
   productHandler.RegisterRoutes(apiV1Group)
   ```
5. **Add Goose migration** (if new tables) - Buat file SQL di `src/migration/migrations/`, misalnya `00003_create_products_table.postgres.sql` dan `.mysql.sql`. Lihat [src/migration/README.md](src/migration/README.md).

## Notes

- File yang sudah ada tidak akan di-overwrite (akan error jika file sudah ada)
- Nama entity akan otomatis di-convert ke format yang benar (e.g., "product" -> "Product")
- Route name akan otomatis plural (e.g., "product" -> "products")
- Table name akan otomatis plural snake_case (e.g., "Product" -> "products")

## Troubleshooting

### File Already Exists
Jika file sudah ada, generator akan error. Hapus file yang sudah ada atau rename terlebih dahulu.

### Invalid Field Format
Pastikan format field sesuai: `Name:Type:Tag` dengan comma sebagai separator.

### Missing Dependencies
Pastikan semua dependencies sudah terinstall dengan `go mod tidy`.
