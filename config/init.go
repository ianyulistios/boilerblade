package config

import (
	"boilerblade/config/amqp"
	"boilerblade/helper"
	"errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrAMQPNotInitialized = errors.New("AMQP connection is not initialized")
)

// AppConfig holds all application configuration and connections
type AppConfig struct {
	Env      *Env
	Database *gorm.DB
	Redis    *redis.Client
	AMQP     amqp.IAMQPConnection
}

// ConnectionOptions defines which connections to initialize
type ConnectionOptions struct {
	EnableDB    bool
	EnableRedis bool
	EnableAMQP  bool
}

// DefaultConnectionOptions returns default connection options (all enabled)
func DefaultConnectionOptions() *ConnectionOptions {
	return &ConnectionOptions{
		EnableDB:    true,
		EnableRedis: true,
		EnableAMQP:  true,
	}
}

// Initialize loads environment variables and initializes connections based on Env flags
func Initialize() (*AppConfig, error) {
	// Load environment variables
	env := &Env{}
	if err := envconfig.Process("", env); err != nil {
		helper.LogError("Failed to load environment variables", err, "", nil)
		return nil, err
	}

	// Use environment flags to determine which connections to initialize
	options := &ConnectionOptions{
		EnableDB:    env.ENABLE_DB,
		EnableRedis: env.ENABLE_REDIS,
		EnableAMQP:  env.ENABLE_AMQP,
	}

	return InitializeWithOptions(env, options)
}

// InitializeWithEnv initializes connections based on provided env
func InitializeWithEnv(env *Env) (*AppConfig, error) {
	// Use environment flags to determine which connections to initialize
	options := &ConnectionOptions{
		EnableDB:    env.ENABLE_DB,
		EnableRedis: env.ENABLE_REDIS,
		EnableAMQP:  env.ENABLE_AMQP,
	}

	return InitializeWithOptions(env, options)
}

// InitializeWithOptions loads environment variables and initializes connections based on provided options
func InitializeWithOptions(env *Env, options *ConnectionOptions) (*AppConfig, error) {
	if options == nil {
		options = DefaultConnectionOptions()
	}

	cfg := &AppConfig{
		Env: env,
	}

	// Initialize Database if enabled
	if options.EnableDB {
		cfg.Database = env.InitDatabase()
		helper.LogInfo("Database connection initialization attempted", map[string]interface{}{
			"enabled": true,
			"ready":   cfg.Database != nil,
		})
	} else {
		helper.LogInfo("Database connection disabled", map[string]interface{}{
			"enabled": false,
		})
	}

	// Initialize Redis if enabled
	if options.EnableRedis {
		cfg.Redis = env.InitRedis()
		helper.LogInfo("Redis connection initialization attempted", map[string]interface{}{
			"enabled": true,
			"ready":   cfg.Redis != nil,
		})
	} else {
		helper.LogInfo("Redis connection disabled", map[string]interface{}{
			"enabled": false,
		})
	}

	// Initialize AMQP if enabled
	var amqpConn *amqp.IAMQPConnection
	if options.EnableAMQP {
		amqpConn = env.InitAMQP()
		if amqpConn != nil {
			cfg.AMQP = *amqpConn
		}
		helper.LogInfo("AMQP connection initialization attempted", map[string]interface{}{
			"enabled": true,
			"ready":   amqpConn != nil,
		})
	} else {
		helper.LogInfo("AMQP connection disabled", map[string]interface{}{
			"enabled": false,
		})
	}

	// Log initialization summary
	helper.LogInfo("Application configuration initialized", map[string]interface{}{
		"mode":          env.MODE,
		"app_name":      env.FIBER_APP_NAME,
		"port":          env.FIBER_PORT,
		"db_type":       env.DB_TYPE,
		"db_enabled":    options.EnableDB,
		"db_ready":      cfg.Database != nil,
		"redis_enabled": options.EnableRedis,
		"redis_ready":   cfg.Redis != nil,
		"amqp_enabled":  options.EnableAMQP,
		"amqp_ready":    amqpConn != nil,
	})

	return cfg, nil
}

// EnsureAMQP ensures AMQP connection is initialized
// If AMQP was disabled via ENABLE_AMQP=false, it will be force-enabled
// This method is useful when AMQP is needed but was not initialized during startup
func (cfg *AppConfig) EnsureAMQP() error {
	if cfg.AMQP == nil {
		// Force enable AMQP if it was disabled
		if !cfg.Env.ENABLE_AMQP {
			helper.LogInfo("ENABLE_AMQP was false, forcing AMQP initialization", map[string]interface{}{
				"source": "AppConfig.EnsureAMQP",
			})
			cfg.Env.ENABLE_AMQP = true
		}

		// Initialize AMQP connection
		amqpConn := cfg.Env.InitAMQP()
		if amqpConn == nil {
			helper.LogError("Failed to initialize AMQP connection", ErrAMQPNotInitialized, "", map[string]interface{}{
				"source": "AppConfig.EnsureAMQP",
			})
			return ErrAMQPNotInitialized
		}

		cfg.AMQP = *amqpConn
		helper.LogInfo("AMQP connection initialized", map[string]interface{}{
			"source": "AppConfig.EnsureAMQP",
		})
	}

	return nil
}
