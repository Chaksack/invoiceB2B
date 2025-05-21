package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
)

// SetupAuthRoutes configures authentication related routes
func SetupAuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler) {
	authGroup := router.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	// authGroup.Post("/login", authHandler.Login) // To be implemented
	// authGroup.Post("/logout", authHandler.Logout) // To be implemented, requires auth middleware
	// authGroup.Post("/refresh-token", authHandler.RefreshToken) // To be implemented
	// authGroup.Post("/2fa/setup", authHandler.Setup2FA) // To be implemented, requires auth middleware
	// authGroup.Post("/2fa/verify", authHandler.Verify2FA) // To be implemented
}

// TODO: Add other route setup files (user_routes.go, admin_routes.go, invoice_routes.go)
