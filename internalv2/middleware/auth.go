package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"invoiceB2B/internal@v2/config"
	"invoiceB2B/internal@v2/models"
)

// Claims represents JWT claims
type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates authentication middleware
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Authorization header required",
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token claims",
			})
		}

		// Set user in context
		c.Locals("user", &models.User{
			ID:    claims.ID,
			Email: claims.Email,
			Role:  claims.Role,
		})

		return c.Next()
	}
}

// RoleMiddleware creates role-based authorization middleware
func RoleMiddleware(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		if user.Role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Insufficient permissions",
			})
		}
		return c.Next()
	}
}

// GetUserFromContext extracts user from context
func GetUserFromContext(c *fiber.Ctx) *models.User {
	user := c.Locals("user").(*models.User)
	return user
} 