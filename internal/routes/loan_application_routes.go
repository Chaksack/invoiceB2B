package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

// SetupLoanApplicationRoutes configures centralized loan application routes for users
func SetupLoanApplicationRoutes(
	router fiber.Router, 
	loanAppHandler *handlers.LoanApplicationHandler, 
	authMw *middleware.AuthMiddleware, 
	csrfMw *middleware.CSRFMiddleware,
) {
	// All loan application routes require authentication
	loanGroup := router.Group("/loan-applications")
	loanGroup.Use(authMw.Protected())
	loanGroup.Use(csrfMw.GenerateToken)

	// Read-only operations (GET requests)
	loanReadGroup := loanGroup.Group("")
	loanReadGroup.Get("", loanAppHandler.GetUserLoanApplications)       // Get user's loan applications
	loanReadGroup.Get("/:id", loanAppHandler.GetLoanApplication)        // Get specific loan application
	loanReadGroup.Get("/stats", loanAppHandler.GetLoanApplicationStats) // Get user's loan statistics

	// State-changing operations (POST, PUT, DELETE) with CSRF protection
	loanWriteGroup := loanGroup.Group("")
	loanWriteGroup.Use(csrfMw.VerifyToken())

	// Core loan application endpoints
	loanWriteGroup.Post("", loanAppHandler.CreateLoanApplication)           // Create new loan application
	loanWriteGroup.Post("/manual", loanAppHandler.CreateManualLoanApplication) // Create manual loan with all data

	// KYB (Know Your Business) operations
	kybGroup := loanWriteGroup.Group("/:id/kyb")
	kybGroup.Post("", loanAppHandler.SubmitKYBInformation) // Submit KYB information

	// Financial statement operations - Unified endpoint
	financialGroup := loanWriteGroup.Group("/:id/financial-statements")
	financialGroup.Post("", loanAppHandler.SubmitFinancialStatementUnified) // Unified financial statement submission (data + file)

	// Document upload operations
	documentsGroup := loanWriteGroup.Group("/:id/documents")
	documentsGroup.Post("/upload", loanAppHandler.UploadDocument) // Upload supporting documents
}