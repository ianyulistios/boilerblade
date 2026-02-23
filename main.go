package main

import (
	"boilerblade/config"
	"boilerblade/server"
	"boilerblade/src/migration"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"

	_ "boilerblade/docs" // swagger docs
)

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

func main() {
	// Load environment variables
	env := &config.Env{}
	if err := envconfig.Process("", env); err != nil {
		log.Fatal("Failed to load environment variables:", err)
	}

	// Initialize app with all configuration and connections
	app, err := server.NewApp(env)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// Run database migrations (Goose)
	if app.Config.Database != nil {
		if err := migration.RunMigrations(app.Config.Database); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Database migration completed")
	}

	// Get server mode from environment (http, amqp, or both)
	serverMode := strings.ToLower(env.SERVER_MODE)
	if serverMode == "" {
		serverMode = "both" // Default to both
	}

	switch serverMode {
	case "http":
		// Start HTTP server only
		log.Println("Starting HTTP server only...")
		app.ServeHTTP()

	case "amqp":
		// Start AMQP consumers only
		log.Println("Starting AMQP consumers only...")
		if err := app.AMQPServe(); err != nil {
			log.Fatal("Failed to start AMQP consumers:", err)
		}

	case "both":
		// Start both HTTP server and AMQP consumers
		log.Println("Starting HTTP server and AMQP consumers...")

		// Start HTTP server in goroutine
		go func() {
			app.ServeHTTP()
		}()

		// Start AMQP consumers in goroutine
		go func() {
			if err := app.AMQPServe(); err != nil {
				log.Fatal("Failed to start AMQP consumers:", err)
			}
		}()

		// Wait for shutdown signal
		app.WaitForShutdown()

	default:
		log.Fatalf("Invalid SERVER_MODE: %s. Valid options: http, amqp, both", serverMode)
	}
}
