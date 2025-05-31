package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

// SetupAdminRoutes configures the routes for admin-specific operations.
func SetupAdminRoutes(
	router fiber.Router,
	adminHandler *handlers.AdminHandler,
	authMw *middleware.AuthMiddleware,
	adminMw *middleware.AdminMiddleware,
) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(authMw.Protected())
	adminGroup.Use(adminMw.AdminRequired())

	// --- Admin Profile ---
	adminGroup.Get("/profile/me", adminHandler.GetAdminProfile)

	// --- Admin User & KYC Management ---
	adminUsersGroup := adminGroup.Group("/users")
	adminUsersGroup.Get("", adminHandler.GetAllUsers)
	adminUsersGroup.Get("/:id", adminHandler.GetUserByID)
	adminUsersGroup.Get("/:id/kyc", adminHandler.GetUserKYCDetail)
	adminUsersGroup.Put("/:id/kyc/review", adminHandler.ReviewKYC)
	adminUsersGroup.Get("/:id/activity-logs", adminHandler.GetUserActivityLogs)

	// --- Admin Invoice Management ---
	adminInvoicesGroup := adminGroup.Group("/invoices")
	adminInvoicesGroup.Get("", adminHandler.GetAllInvoices)
	adminInvoicesGroup.Get("/:id", adminHandler.GetInvoiceDetail)
	adminInvoicesGroup.Put("/:id/status", adminHandler.UpdateInvoiceStatus)
	adminInvoicesGroup.Post("/:id/receipt", adminHandler.UploadDisbursementReceipt)
	adminInvoicesGroup.Get("/:id/download-pdf", adminHandler.DownloadInvoicePDF)

	// --- Admin Staff Management ---
	adminStaffGroup := adminGroup.Group("/staff")
	adminStaffGroup.Post("", adminHandler.CreateStaff)
	adminStaffGroup.Get("", adminHandler.GetAllStaff)
	adminStaffGroup.Put("/:id", adminHandler.UpdateStaff)
	adminStaffGroup.Delete("/:id", adminHandler.DeleteStaff)

	// --- Admin Activity Logs & Analytics ---
	adminGroup.Get("/activity-logs", adminHandler.GetActivityLogs)
	// Dashboard analytics
	adminGroup.Get("/dashboard/analytics", adminHandler.GetAdminDashboardAnalytics)
}
