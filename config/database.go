package config

import (
	"boilerblade/helper"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (e *Env) InitDatabase() *gorm.DB {
	// Configure GORM logger based on mode
	var logLevel logger.LogLevel
	if e.MODE == "production" {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	// Generate DSN based on database type
	dsn, dbType := e.getDSN()

	var db *gorm.DB
	var err error

	// Open connection based on database type
	switch strings.ToLower(e.DB_TYPE) {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
	case "postgres", "postgresql":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
	default:
		helper.LogError("Unsupported database type", fmt.Errorf("database type %s is not supported", e.DB_TYPE), e.DB_TYPE, map[string]interface{}{
			"db_type":   e.DB_TYPE,
			"supported": []string{"mysql", "postgres", "postgresql"},
		})
		return nil
	}

	if err != nil {
		helper.LogError("Database connection failed", err, e.DB_HOST, map[string]interface{}{
			"db_type":  dbType,
			"host":     e.DB_HOST,
			"port":     e.DB_PORT,
			"user":     e.DB_USER,
			"password": "***",
			"dbname":   e.DB_NAME,
		})
		return nil
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		helper.LogError("Database connection pool setup failed", err, e.DB_HOST, nil)
		return nil
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(e.DB_MAX_OPEN_CONNS)
	sqlDB.SetMaxIdleConns(e.DB_MAX_IDLE_CONNS)
	sqlDB.SetConnMaxLifetime(time.Duration(e.DB_MAX_LIFETIME_CONNS) * time.Second)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		helper.LogError("Database ping failed", err, e.DB_HOST, map[string]interface{}{
			"host": e.DB_HOST,
			"port": e.DB_PORT,
		})
		return nil
	}

	helper.LogInfo("Database connection initialized", map[string]interface{}{
		"db_type":            dbType,
		"host":               e.DB_HOST,
		"port":               e.DB_PORT,
		"user":               e.DB_USER,
		"dbname":             e.DB_NAME,
		"max_open_conns":     e.DB_MAX_OPEN_CONNS,
		"max_idle_conns":     e.DB_MAX_IDLE_CONNS,
		"max_lifetime_conns": e.DB_MAX_LIFETIME_CONNS,
	})

	return db
}

// getDSN generates the Data Source Name based on database type
func (e *Env) getDSN() (string, string) {
	dbType := strings.ToLower(e.DB_TYPE)

	switch dbType {
	case "mysql":
		// MySQL DSN format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			e.DB_USER, e.DB_PASSWORD, e.DB_HOST, e.DB_PORT, e.DB_NAME)
		return dsn, "mysql"
	case "postgres", "postgresql":
		// PostgreSQL DSN format: host=host port=port user=user password=password dbname=dbname sslmode=disable TimeZone=UTC
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			e.DB_HOST, e.DB_PORT, e.DB_USER, e.DB_PASSWORD, e.DB_NAME)
		return dsn, "postgres"
	default:
		// Default to PostgreSQL
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			e.DB_HOST, e.DB_PORT, e.DB_USER, e.DB_PASSWORD, e.DB_NAME)
		return dsn, "postgres"
	}
}
