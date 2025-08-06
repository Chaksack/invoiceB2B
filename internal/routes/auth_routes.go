package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

func SetupAuthRoutes(router fiber.Router, authHandler *handlers.AuthHandler, authMw *middleware.AuthMiddleware, csrfMw *middleware.CSRFMiddleware) {
	authGroup := router.Group("/auth")
	
	// Generate CSRF token for all auth routes
	authGroup.Use(csrfMw.GenerateToken)
	
	// Apply CSRF verification to all state-changing endpoints
	// Skip CSRF for initial login and registration to allow first-time users
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	
	// Apply CSRF verification to all other state-changing endpoints
	csrfProtected := authGroup.Group("")
	csrfProtected.Use(csrfMw.VerifyToken())
	
	csrfProtected.Post("/login/2fa/verify", authHandler.Verify2FA) // Verify OTP after login attempt
	csrfProtected.Post("/refresh-token", authHandler.RefreshToken)

	// Routes requiring authentication
	authRequired := csrfProtected.Group("")
	authRequired.Use(authMw.Protected()) // Apply JWT auth middleware

	authRequired.Post("/logout", authHandler.Logout)
	authRequired.Post("/2fa/toggle", authHandler.Enable2FA) // Enable/disable 2FA for logged-in user
}
