package migration

import (
	"boilerblade/helper"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// RunMigrations runs all pending Goose SQL migrations using the same database
// connection as the given GORM DB. Dialect is inferred from the GORM dialector
// (postgres, mysql supported).
func RunMigrations(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get underlying *sql.DB: %w", err)
	}

	dialect, err := gooseDialectFromGORM(db)
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(dialect); err != nil {
		helper.LogError("Failed to set Goose dialect", err, dialect, map[string]interface{}{
			"source": "RunMigrations",
		})
		return fmt.Errorf("set goose dialect: %w", err)
	}

	helper.LogInfo("Starting Goose migrations", map[string]interface{}{
		"source":  "RunMigrations",
		"dialect": dialect,
	})

	ctx := context.Background()
	if err := goose.UpContext(ctx, sqlDB, "migrations"); err != nil {
		helper.LogError("Goose migration failed", err, "", map[string]interface{}{
			"source": "RunMigrations",
		})
		return fmt.Errorf("goose up: %w", err)
	}

	helper.LogInfo("Goose migrations completed successfully", map[string]interface{}{
		"source": "RunMigrations",
	})
	return nil
}

// gooseDialectFromGORM returns the Goose dialect name from the GORM DB dialector.
func gooseDialectFromGORM(db *gorm.DB) (string, error) {
	name := db.Dialector.Name()
	switch strings.ToLower(name) {
	case "postgres", "postgresql":
		return "postgres", nil
	case "mysql":
		return "mysql", nil
	default:
		return "", fmt.Errorf("unsupported database dialect for Goose: %s", name)
	}
}

// RunMigrationsWithDB runs Goose migrations using a raw *sql.DB and dialect.
// Use this when you have *sql.DB and dialect ("postgres" or "mysql") without GORM.
func RunMigrationsWithDB(sqlDB *sql.DB, dialect string) error {
	if sqlDB == nil {
		return fmt.Errorf("database connection is nil")
	}
	dialect = strings.ToLower(dialect)
	if dialect != "postgres" && dialect != "mysql" {
		return fmt.Errorf("unsupported dialect for Goose: %s", dialect)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	ctx := context.Background()
	return goose.UpContext(ctx, sqlDB, "migrations")
}
