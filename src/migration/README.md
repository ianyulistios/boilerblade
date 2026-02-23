# Goose Migration System

Migrations are managed with [Goose](https://github.com/pressly/goose). SQL migration files are stored in `src/migration/migrations/` and are embedded into the binary. Dialect is inferred from the GORM database (PostgreSQL or MySQL).

## How It Works

- **SQL migrations**: Each file is named `NNNNN_description.dialect.sql` (e.g. `00001_create_users_table.postgres.sql`, `00001_create_users_table.mysql.sql`). Goose runs only the files matching the current database dialect.
- **Up/Down**: Each file has `-- +goose Up` and `-- +goose Down` sections. On startup, `RunMigrations` runs all pending **Up** migrations.
- **Versioning**: Goose tracks applied migrations in the `goose_db_version` table.

## Adding a New Migration

### Option 1: CLI (recommended)

From the project root:

```bash
boilerblade make migration -name=add_orders_table
```

This creates two files in `src/migration/migrations/`: `<timestamp>_add_orders_table.postgres.sql` and `<timestamp>_add_orders_table.mysql.sql` with placeholder `-- +goose Up` and `-- +goose Down` sections. Edit them to add your SQL.

### Option 2: Manual

1. Add one or two new SQL files under `src/migration/migrations/`:
   - For both dialects: `00003_create_orders_table.postgres.sql` and `00003_create_orders_table.mysql.sql`
   - Or a single dialect-neutral file if your SQL is compatible: `00003_create_orders_table.sql`

2. Use this format:

```sql
-- +goose Up
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    -- ...
);

-- +goose Down
DROP TABLE IF EXISTS orders;
```

3. Rebuild and run the app; migrations run automatically on startup.

## API

### `RunMigrations(db *gorm.DB) error`

Runs all pending Goose migrations using the same connection as the given GORM DB. Dialect is taken from GORM (postgres or mysql).

```go
if err := migration.RunMigrations(app.Config.Database); err != nil {
    log.Fatal("Failed to migrate:", err)
}
```

### `RunMigrationsWithDB(sqlDB *sql.DB, dialect string) error`

Runs migrations with a raw `*sql.DB` and dialect (`"postgres"` or `"mysql"`). Use when you are not using GORM.

```go
if err := migration.RunMigrationsWithDB(sqlDB, "postgres"); err != nil {
    log.Fatal("Failed to migrate:", err)
}
```

## Usage in main.go

```go
import "boilerblade/src/migration"

func main() {
    app, _ := server.NewApp(env)

    if app.Config.Database != nil {
        if err := migration.RunMigrations(app.Config.Database); err != nil {
            log.Fatal("Failed to migrate database:", err)
        }
        log.Println("Database migration completed")
    }
}
```

## CLI (optional)

You can also run migrations via the Goose CLI:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
GOOSE_DRIVER=postgres GOOSE_DBSTRING="host=localhost port=5432 user=postgres password=postgres dbname=boilerblade sslmode=disable" GOOSE_MIGRATION_DIR=./src/migration/migrations goose up
```

## Current Migrations

- `00001_create_users_table` – users table (PostgreSQL + MySQL)
- `00002_create_products_table` – products table (PostgreSQL + MySQL)

## MySQL note

For MySQL, the DSN should include `multiStatements=true` when a single migration file runs multiple statements (e.g. `CREATE TABLE` plus `CREATE INDEX`). The project’s database config adds this where needed.
