package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

func SetupAuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler, authMw *middleware.AuthMiddleware) {
	authGroup := router.Group("/auth")

	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/login/2fa/verify", authHandler.Verify2FA) // Verify OTP after login attempt

	authGroup.Post("/refresh-token", authHandler.RefreshToken)

	// Routes requiring authentication
	authRequired := authGroup.Group("")
	authRequired.Use(authMw.Protected()) // Apply JWT auth middleware

	authRequired.Post("/logout", authHandler.Logout)
	authRequired.Post("/2fa/toggle", authHandler.Enable2FA) // Enable/disable 2FA for logged-in user
}
