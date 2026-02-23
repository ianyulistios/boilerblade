package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const (
	migrationDir = "src/migration/migrations"
)

// MigrationGen generates a new Goose SQL migration (postgres + mysql placeholder files).
type MigrationGen struct {
	Name      string // user-provided name (e.g. add_orders_table)
	Version   string // timestamp version (YYYYMMDDHHMMSS)
	Normalized string // snake_case name for filename
}

// NewMigrationGen creates a migration generator from the given name.
// Name is normalized to snake_case (e.g. "Add Orders Table" -> "add_orders_table").
func NewMigrationGen(name string) *MigrationGen {
	name = strings.TrimSpace(name)
	normalized := toSnakeCase(name)
	if normalized == "" {
		normalized = "migration"
	}
	return &MigrationGen{
		Name:       name,
		Version:    time.Now().Format("20060102150405"),
		Normalized: normalized,
	}
}

func toSnakeCase(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '-' {
			b.WriteByte('_')
		} else if (unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_') && unicode.IsLower(r) {
			b.WriteRune(r)
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	result := b.String()
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return strings.Trim(result, "_")
}

// Generate creates the migration files (postgres and mysql) with placeholder Up/Down.
func (m *MigrationGen) Generate() error {
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations dir: %w", err)
	}

	baseName := m.Version + "_" + m.Normalized

	postgresPath := filepath.Join(migrationDir, baseName+".postgres.sql")
	mysqlPath := filepath.Join(migrationDir, baseName+".mysql.sql")

	for path, content := range map[string]string{
		postgresPath: m.postgresContent(),
		mysqlPath:    m.mysqlContent(),
	} {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file already exists: %s", path)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}

	return nil
}

func (m *MigrationGen) postgresContent() string {
	return `-- +goose Up
-- TODO: add your PostgreSQL migration SQL here
-- Example:
-- CREATE TABLE example (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
-- );

-- +goose Down
-- TODO: add your rollback SQL here (e.g. DROP TABLE IF EXISTS example;
`
}

func (m *MigrationGen) mysqlContent() string {
	return `-- +goose Up
-- TODO: add your MySQL migration SQL here
-- Example:
-- CREATE TABLE example (
--     id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
--     updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
-- TODO: add your rollback SQL here (e.g. DROP TABLE IF EXISTS example;
`
}

// GeneratedFiles returns the paths of the files that would be created (for messaging).
func (m *MigrationGen) GeneratedFiles() (postgres, mysql string) {
	baseName := m.Version + "_" + m.Normalized
	return filepath.Join(migrationDir, baseName+".postgres.sql"),
		filepath.Join(migrationDir, baseName+".mysql.sql")
}
