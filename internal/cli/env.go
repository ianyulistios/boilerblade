package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

// envExampleContent is the default .env.example template written when generating a project or when missing.
const envExampleContent = `# =============================================================================
# Boilerblade Environment - Copy to .env and fill in your values
# =============================================================================

# --- App ---
MODE=development
FIBER_PORT=3000
FIBER_APP_NAME=boilerblade
APP_KEY=your-secret-key-for-jwt-min-32-chars
SERVER_MODE=both

# --- Connection flags (true/false) ---
ENABLE_DB=true
ENABLE_REDIS=true
ENABLE_AMQP=true

# --- Database (used by app and Goose migrations) ---
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=boilerblade
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=10
DB_MAX_LIFETIME_CONNS=10

# --- Goose CLI (optional; for running goose from command line) ---
# GOOSE_DRIVER=postgres
# GOOSE_DBSTRING=host=localhost port=5432 user=postgres password=postgres dbname=boilerblade sslmode=disable
# GOOSE_MIGRATION_DIR=./src/migration/migrations
# GOOSE_TABLE=goose_db_version

# --- Redis ---
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# --- AMQP (e.g. RabbitMQ) ---
AMQP_HOST=localhost
AMQP_PORT=5672
AMQP_USER=guest
AMQP_PASSWORD=guest
`

// EnsureEnvExample creates .env.example in dir if it does not exist.
// Used when generating a project so the new project has the env template.
func EnsureEnvExample(dir string) error {
	p := filepath.Join(dir, ".env.example")
	if _, err := os.Stat(p); err == nil {
		return nil
	}
	return os.WriteFile(p, []byte(envExampleContent), 0644)
}

// EnsureEnvFromExample creates .env from .env.example in dir if .env does not exist.
func EnsureEnvFromExample(dir string) (created bool, err error) {
	envPath := filepath.Join(dir, ".env")
	examplePath := filepath.Join(dir, ".env.example")
	if _, err := os.Stat(envPath); err == nil {
		return false, nil
	}
	content, err := os.ReadFile(examplePath)
	if err != nil {
		return false, fmt.Errorf("read .env.example: %w", err)
	}
	if err := os.WriteFile(envPath, content, 0644); err != nil {
		return false, err
	}
	return true, nil
}
