package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

func SetupInvoiceRoutes(router fiber.Router, invoiceHandler *handlers.InvoiceHandler, authMw *middleware.AuthMiddleware, adminMw *middleware.AdminMiddleware, csrfMw *middleware.CSRFMiddleware) {
	userInvoiceGroup := router.Group("/invoices")
	
	// Apply authentication middleware
	userInvoiceGroup.Use(authMw.Protected())
	
	// Generate CSRF token for all invoice routes
	userInvoiceGroup.Use(csrfMw.GenerateToken)
	
	// Set up GET routes (no CSRF verification needed)
	userInvoiceGroup.Get("", invoiceHandler.GetUserInvoices)
	userInvoiceGroup.Get("/:id", invoiceHandler.GetInvoiceByID)
	userInvoiceGroup.Get("/:id/viewreceipt", invoiceHandler.ViewReceipt)
	userInvoiceGroup.Get("/:id/receipt", invoiceHandler.DownloadReceipt)
	userInvoiceGroup.Get("/:id/suggested-financial-institutions", invoiceHandler.GetSuggestedFinancialInstitutions)
	
	// Set up POST routes with CSRF verification
	csrfProtected := userInvoiceGroup.Group("")
	csrfProtected.Use(csrfMw.VerifyToken())
	
	csrfProtected.Post("", invoiceHandler.UploadInvoice)
	csrfProtected.Post("/:id/select-financial-institution", invoiceHandler.SelectFinancialInstitution)
}
