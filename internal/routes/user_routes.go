package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

func SetupUserRoutes(router fiber.Router, userHandler *handlers.UserHandler, authMw *middleware.AuthMiddleware, csrfMw *middleware.CSRFMiddleware) {
	userGroup := router.Group("/user")
	userGroup.Use(authMw.Protected()) // All user routes require authentication
	
	// Generate CSRF token for all user routes
	userGroup.Use(csrfMw.GenerateToken)
	
	// Create a group for read-only operations (GET)
	userReadGroup := userGroup.Group("")
	userReadGroup.Get("/profile", userHandler.GetUserProfile)
	
	// Create a group for state-changing operations (PUT, POST)
	// that require CSRF verification
	userWriteGroup := userGroup.Group("")
	userWriteGroup.Use(csrfMw.VerifyToken())
	userWriteGroup.Put("/profile", userHandler.UpdateUserProfile)
	
	// KYC operations
	kycReadGroup := userReadGroup.Group("/kyc")
	kycReadGroup.Get("", userHandler.GetKYCStatus)
	
	kycWriteGroup := userWriteGroup.Group("/kyc")
	kycWriteGroup.Post("", userHandler.SubmitKYC)
}
