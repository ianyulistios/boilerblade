package server

import (
	"boilerblade/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	*fiber.App
	Config *config.AppConfig
}

// NewApp creates a new App instance with initialized configuration
func NewApp(env *config.Env) (*App, error) {
	// Initialize connections based on provided env
	cfg, err := config.InitializeWithEnv(env)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config: cfg,
	}

	return app, nil
}

// WaitForShutdown waits for interrupt signal to gracefully shutdown
func (a *App) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")
}
