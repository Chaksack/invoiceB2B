package routes

import (
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInternalRoutes(router fiber.Router, internalHandler *handlers.InternalHandler, internalApiMw *middleware.InternalAPIMiddleware) {
	internalGroup := router.Group("/internal")
	internalGroup.Use(internalApiMw.RequireAPIKey())

	internalGroup.Put("/invoices/:id/processed-data", internalHandler.UpdateInvoiceWithProcessedData)
}
