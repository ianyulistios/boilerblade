package server

import (
	"boilerblade/middleware"
	"boilerblade/src/handler"
	"boilerblade/src/repository"
	"boilerblade/src/usecase"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/keyauth/v2"
	"github.com/gofiber/swagger"
)

// App and Connection types are defined in app.go

func (a *App) Routes() {
	// Swagger documentation route (before authentication)
	a.Get("/swagger/*", swagger.HandlerDefault)

	apiV1Group := a.Group("/api/v1")

	apiV1Group.Use(recover.New())
	apiV1Group.Use(logger.New())

	// Allow all origins, methods, and headers
	apiV1Group.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: fmt.Sprintf("%s,%s,%s,%s", fiber.MethodPut, fiber.MethodPost, fiber.MethodGet, fiber.MethodDelete),
	}))

	// Store env in context for middleware access
	apiV1Group.Use(func(c *fiber.Ctx) error {
		c.Locals("env", a.Config.Env)
		return c.Next()
	})

	// JWT Authentication middleware
	apiV1Group.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:Authorization",
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			return middleware.AuthValidator(key, c)
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
		},
	}))

	// Initialize dependencies
	userRepo := repository.NewUserRepository(a.Config.Database)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// Register handler routes
	userHandler.RegisterRoutes(apiV1Group)
}

func (a *App) ServeHTTP() {
	// Initialize Fiber app
	a.App = fiber.New(fiber.Config{
		AppName: a.Config.Env.FIBER_APP_NAME,
	})

	//Init Routes
	a.Routes()

	// Get port from environment config
	port := a.Config.Env.FIBER_PORT
	if port == "" {
		port = "3000" // Default port
	}

	listenerPort := fmt.Sprintf(":%s", port)
	log.Fatal(a.Listen(listenerPort))
}
