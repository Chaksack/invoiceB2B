package routes

import (
	"github.com/gofiber/fiber/v2"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
)

func SetupInvoiceRoutes(router fiber.Router, invoiceHandler *handlers.InvoiceHandler, authMw *middleware.AuthMiddleware, adminMw *middleware.AdminMiddleware) {
	userInvoiceGroup := router.Group("/invoices")
	userInvoiceGroup.Use(authMw.Protected())

	userInvoiceGroup.Post("", invoiceHandler.UploadInvoice)
	userInvoiceGroup.Get("", invoiceHandler.GetUserInvoices)
	userInvoiceGroup.Get("/:id", invoiceHandler.GetInvoiceByID)
	userInvoiceGroup.Get("/:id/viewreceipt", invoiceHandler.ViewReceipt)
	userInvoiceGroup.Get("/:id/receipt", invoiceHandler.DownloadReceipt)
}
