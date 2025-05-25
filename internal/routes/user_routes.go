package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

func SetupUserRoutes(router fiber.Router, userHandler *handlers.UserHandler, authMw *middleware.AuthMiddleware) {
	userGroup := router.Group("/user")
	userGroup.Use(authMw.Protected()) // All user routes require authentication

	userGroup.Get("/profile", userHandler.GetUserProfile)
	userGroup.Put("/profile", userHandler.UpdateUserProfile)

	kycGroup := userGroup.Group("/kyc")
	kycGroup.Post("", userHandler.SubmitKYC)
	kycGroup.Get("", userHandler.GetKYCStatus)
}
