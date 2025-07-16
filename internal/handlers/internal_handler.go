package handlers

import (
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InternalHandler struct {
	internalService services.InternalService
	validate        *validator.Validate
}

func NewInternalHandler(internalService services.InternalService, validate *validator.Validate) *InternalHandler {
	return &InternalHandler{
		internalService: internalService,
		validate:        validate,
	}
}

// UpdateInvoiceWithProcessedData godoc
// @Summary Update invoice with processed data
// @Description Update an invoice with data processed by n8n or other internal services
// @Tags internal
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param request body dtos.UpdateInvoiceProcessedDataRequest true "Processed invoice data"
// @Success 200 {object} dtos.InvoiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /internal/invoices/{id}/processed-data [put]
func (h *InternalHandler) UpdateInvoiceWithProcessedData(c *fiber.Ctx) error {
	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	var req dtos.UpdateInvoiceProcessedDataRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body for processed data.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	// The IP address here would be of the n8n service if it's making the call
	ipAddress := c.IP()

	updatedInvoice, err := h.internalService.UpdateInvoiceWithProcessedData(c.Context(), uint(invoiceID), req, ipAddress)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update invoice with processed data.", err)
	}

	return c.Status(fiber.StatusOK).JSON(updatedInvoice)
}
