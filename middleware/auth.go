package middleware

import (
	"boilerblade/config"
	"boilerblade/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthValidator validates JWT token from Authorization header
func AuthValidator(token string, c *fiber.Ctx) (bool, error) {
	// Get APP_KEY from context or environment
	env := c.Locals("env")
	var appKey string
	if env != nil {
		if envConfig, ok := env.(*config.Env); ok {
			appKey = envConfig.APP_KEY
		}
	}

	if appKey == "" {
		helper.LogError("JWT validation failed: APP_KEY not configured", nil, "", nil)
		return false, fiber.NewError(fiber.StatusInternalServerError, "Server configuration error")
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		helper.LogError("JWT validation failed: empty token", nil, "", nil)
		return false, fiber.NewError(fiber.StatusUnauthorized, "Missing or invalid token")
	}

	// Parse and validate JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte(appKey), nil
	})

	if err != nil {
		helper.LogError("JWT validation failed", err, "", map[string]interface{}{
			"token": token[:min(20, len(token))] + "...",
		})
		return false, fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	if !parsedToken.Valid {
		helper.LogError("JWT validation failed: invalid token", nil, "", nil)
		return false, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Extract claims and store in context
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		// Store user ID and other claims in context for use in handlers
		if userID, ok := claims["user_id"].(string); ok {
			c.Locals("user_id", userID)
		}
		if email, ok := claims["email"].(string); ok {
			c.Locals("email", email)
		}
		c.Locals("claims", claims)
	}

	helper.LogInfo("JWT validation successful", map[string]interface{}{
		"path":   c.Path(),
		"method": c.Method(),
	})

	return true, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
