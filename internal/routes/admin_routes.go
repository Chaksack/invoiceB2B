package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

func SetupAdminRoutes(
	router fiber.Router,
	adminHandler *handlers.AdminHandler,
	authMw *middleware.AuthMiddleware,
	adminMw *middleware.AdminMiddleware,
) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(authMw.Protected())      // First, ensure user is authenticated
	adminGroup.Use(adminMw.AdminRequired()) // Then, ensure user is an admin

	// Admin User & KYC Management
	adminUsersGroup := adminGroup.Group("/users")
	adminUsersGroup.Get("", adminHandler.GetAllUsers)
	adminUsersGroup.Get("/:id/kyc", adminHandler.GetUserKYCDetail)
	adminUsersGroup.Put("/:id/kyc/review", adminHandler.ReviewKYC) // Combines approve/reject

	// Admin Invoice Management
	adminInvoicesGroup := adminGroup.Group("/invoices")
	adminInvoicesGroup.Get("", adminHandler.GetAllInvoices)
	adminInvoicesGroup.Get("/:id", adminHandler.GetInvoiceDetail)
	adminInvoicesGroup.Put("/:id/status", adminHandler.UpdateInvoiceStatus) // For approve, reject, disburse
	adminInvoicesGroup.Post("/:id/receipt", adminHandler.UploadDisbursementReceipt)

	// Admin Staff Management
	adminStaffGroup := adminGroup.Group("/staff")
	adminStaffGroup.Post("", adminHandler.CreateStaff)
	adminStaffGroup.Get("", adminHandler.GetAllStaff)
	adminStaffGroup.Put("/:id", adminHandler.UpdateStaff)
	adminStaffGroup.Delete("/:id", adminHandler.DeleteStaff)

	// Admin Activity Logs
	adminGroup.Get("/activity-logs", adminHandler.GetActivityLogs)
}
