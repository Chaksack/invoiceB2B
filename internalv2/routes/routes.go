package routes

import (
	"github.com/gofiber/fiber/v2"

	"invoiceB2B/internal@v2/config"
	"invoiceB2B/internal@v2/database"
	"invoiceB2B/internal@v2/handlers"
	"invoiceB2B/internal@v2/middleware"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, db *database.DB, cfg *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	businessHandler := handlers.NewBusinessHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(cfg)

	// API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/profile", authMiddleware, authHandler.Profile)

	// Business routes (require authentication)
	business := api.Group("/business", authMiddleware, middleware.RoleMiddleware("business"))
	business.Get("/invoices", businessHandler.GetInvoices)
	business.Get("/profile", businessHandler.GetProfile)
	business.Put("/profile", businessHandler.UpdateProfile)

	// Admin routes (require admin authentication)
	admin := api.Group("/admin", authMiddleware, middleware.RoleMiddleware("admin"))
	admin.Get("/businesses", adminHandler.GetBusinesses)
	admin.Get("/dashboard-summary", adminHandler.GetDashboardSummary)
	admin.Put("/businesses/:id/status", adminHandler.UpdateBusinessStatus)
} 