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
	loanAppHandler *handlers.LoanApplicationHandler,
	reportingHandler *handlers.ReportingHandler,
	authMw *middleware.AuthMiddleware,
	adminMw *middleware.AdminMiddleware,
	csrfMw *middleware.CSRFMiddleware,
) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(authMw.Protected())
	adminGroup.Use(adminMw.AdminRequired())
	
	// Generate CSRF token for all admin routes
	adminGroup.Use(csrfMw.GenerateToken)
	
	// Create a group for read-only operations (GET)
	adminReadGroup := adminGroup.Group("")
	
	// Create a group for state-changing operations (POST, PUT, DELETE)
	// that require CSRF verification
	adminWriteGroup := adminGroup.Group("")
	adminWriteGroup.Use(csrfMw.VerifyToken())

	// --- Admin Profile ---
	adminReadGroup.Get("/profile/me", adminHandler.GetAdminProfile)

	// --- Admin User & KYC Management ---
	// Read operations
	adminUsersReadGroup := adminReadGroup.Group("/users")
	adminUsersReadGroup.Get("", adminHandler.GetAllUsers)
	adminUsersReadGroup.Get("/:id", adminHandler.GetUserByID)
	adminUsersReadGroup.Get("/:id/kyc", adminHandler.GetUserKYCDetail)
	adminUsersReadGroup.Get("/:id/activity-logs", adminHandler.GetUserActivityLogs)
	
	// Write operations
	adminUsersWriteGroup := adminWriteGroup.Group("/users")
	adminUsersWriteGroup.Put("/:id/kyc/review", adminHandler.ReviewKYC)

	// --- Admin Invoice Management ---
	// Read operations
	adminInvoicesReadGroup := adminReadGroup.Group("/invoices")
	adminInvoicesReadGroup.Get("", adminHandler.GetAllInvoices)
	adminInvoicesReadGroup.Get("/:id", adminHandler.GetInvoiceDetail)
	adminInvoicesReadGroup.Get("/:id/download-pdf", adminHandler.DownloadInvoicePDF)
	
	// Write operations
	adminInvoicesWriteGroup := adminWriteGroup.Group("/invoices")
	adminInvoicesWriteGroup.Put("/:id/status", adminHandler.UpdateInvoiceStatus)
	adminInvoicesWriteGroup.Post("/:id/receipt", adminHandler.UploadDisbursementReceipt)

	// --- Admin Staff Management ---
	// Read operations
	adminStaffReadGroup := adminReadGroup.Group("/staff")
	adminStaffReadGroup.Get("", adminHandler.GetAllStaff)
	
	// Write operations
	adminStaffWriteGroup := adminWriteGroup.Group("/staff")
	adminStaffWriteGroup.Post("", adminHandler.CreateStaff)
	adminStaffWriteGroup.Put("/:id", adminHandler.UpdateStaff)
	adminStaffWriteGroup.Delete("/:id", adminHandler.DeleteStaff)

	// --- Admin Activity Logs & Analytics ---
	adminReadGroup.Get("/activity-logs", adminHandler.GetActivityLogs)
	// Dashboard analytics
	adminReadGroup.Get("/dashboard/analytics", adminHandler.GetAdminDashboardAnalytics)

	// --- Admin Financial Institution Management ---
	// Read operations
	adminFIReadGroup := adminReadGroup.Group("/financial-institutions")
	adminFIReadGroup.Get("", adminHandler.GetAllFinancialInstitutions)
	adminFIReadGroup.Get("/:id", adminHandler.GetFinancialInstitutionByID)
	adminFIReadGroup.Get("/products/:id", adminHandler.GetFinancialInstitutionProductByID)
	adminFIReadGroup.Get("/:fiId/products", adminHandler.GetFinancialInstitutionProducts)
	adminFIReadGroup.Get("/terms/:id", adminHandler.GetFinancialInstitutionTermByID)
	adminFIReadGroup.Get("/:fiId/terms", adminHandler.GetFinancialInstitutionTerms)
	adminFIReadGroup.Get("/products/:productId/terms", adminHandler.GetFinancialInstitutionTermsByProduct)
	
	// Write operations
	adminFIWriteGroup := adminWriteGroup.Group("/financial-institutions")
	adminFIWriteGroup.Post("", adminHandler.CreateFinancialInstitution)
	adminFIWriteGroup.Put("/:id", adminHandler.UpdateFinancialInstitution)
	adminFIWriteGroup.Delete("/:id", adminHandler.DeleteFinancialInstitution)
	adminFIWriteGroup.Post("/products", adminHandler.CreateFinancialInstitutionProduct)
	adminFIWriteGroup.Put("/products/:id", adminHandler.UpdateFinancialInstitutionProduct)
	adminFIWriteGroup.Delete("/products/:id", adminHandler.DeleteFinancialInstitutionProduct)
	adminFIWriteGroup.Post("/terms", adminHandler.CreateFinancialInstitutionTerm)
	adminFIWriteGroup.Put("/terms/:id", adminHandler.UpdateFinancialInstitutionTerm)
	adminFIWriteGroup.Delete("/terms/:id", adminHandler.DeleteFinancialInstitutionTerm)

	// --- Admin Loan Application Management ---
	// Read operations
	adminLoanReadGroup := adminReadGroup.Group("/loan-applications")
	adminLoanReadGroup.Get("", loanAppHandler.GetAllLoanApplications)
	adminLoanReadGroup.Get("/:id", loanAppHandler.GetLoanApplicationDetail)
	adminLoanReadGroup.Get("/:id/kyb", loanAppHandler.GetLoanApplicationKYB)
	adminLoanReadGroup.Get("/stats", loanAppHandler.GetLoanApplicationStats)
	
	// Write operations
	adminLoanWriteGroup := adminWriteGroup.Group("/loan-applications")
	adminLoanWriteGroup.Post("/:id/review", loanAppHandler.AdminReviewLoanApplication)
	adminLoanWriteGroup.Put("/:id/kyb/review", loanAppHandler.ReviewKYBInformation)
	adminLoanWriteGroup.Post("/:id/send-to-financial-institution", loanAppHandler.SendToFinancialInstitution)

	// --- Admin Reporting ---
	// Read operations (reports are read-only by nature)
	adminReportsReadGroup := adminReadGroup.Group("/reports")
	adminReportsReadGroup.Get("/loan-applications", reportingHandler.GenerateLoanApplicationReport)
	adminReportsReadGroup.Get("/invoices", reportingHandler.GenerateInvoiceReport)
	adminReportsReadGroup.Get("/users", reportingHandler.GenerateUserReport)
	adminReportsReadGroup.Get("/activity", reportingHandler.GenerateActivityReport)
	adminReportsReadGroup.Get("/financial-institutions", reportingHandler.GenerateFinancialInstitutionReport)
	adminReportsReadGroup.Get("/overview", reportingHandler.GenerateOverviewReport)
	
	// Write operations (for generating and exporting reports)
	adminReportsWriteGroup := adminWriteGroup.Group("/reports")
	adminReportsWriteGroup.Post("/generate", reportingHandler.GenerateReport)
	adminReportsWriteGroup.Post("/export", reportingHandler.ExportReport)
}
