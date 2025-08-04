package handlers

import (
	"fmt"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type InvoiceHandler struct {
	invoiceService services.InvoiceService
	fileService    services.FileService
	validate       *validator.Validate
}

func NewInvoiceHandler(invoiceService services.InvoiceService, fileService services.FileService, validate *validator.Validate) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
		fileService:    fileService,
		validate:       validate,
	}
}

// UploadInvoice godoc
// @Summary Upload a new invoice
// @Description Upload a new invoice file for processing
// @Tags invoices
// @Accept multipart/form-data
// @Produce json
// @Param invoiceFile formData file true "Invoice file to upload (PDF, CSV, JPEG, JPG, PNG)"
// @Success 201 {object} dtos.InvoiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /invoices [post]
func (h *InvoiceHandler) UploadInvoice(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	file, err := c.FormFile("invoiceFile")
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invoice file is required.", err)
	}

	// Updated allowed extensions
	allowedExtensions := map[string]bool{
		".pdf": true, ".csv": true,
		".jpeg": true, ".jpg": true, ".png": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid file type. Allowed: PDF, CSV, JPEG, JPG, PNG.", nil)
	}
	if err := h.fileService.ValidateFileSize(file.Size); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, err.Error(), err)
	}

	req := dtos.InvoiceUploadRequest{File: file}

	invoiceResponse, err := h.invoiceService.CreateInvoice(c.Context(), uint(userID), req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to upload invoice.", err)
	}

	return c.Status(fiber.StatusCreated).JSON(invoiceResponse)
}

// GetUserInvoices godoc
// @Summary Get user invoices
// @Description Get a paginated list of invoices for the authenticated user
// @Tags invoices
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Success 200 {object} dtos.InvoiceListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /invoices [get]
func (h *InvoiceHandler) GetUserInvoices(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	invoices, total, err := h.invoiceService.GetUserInvoices(c.Context(), uint(userID), page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve invoices.", err)
	}

	return c.Status(fiber.StatusOK).JSON(dtos.InvoiceListResponse{
		Invoices: invoices,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetInvoiceByID godoc
// @Summary Get invoice by ID
// @Description Get details of a specific invoice by its ID
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} dtos.InvoiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /invoices/{id} [get]
func (h *InvoiceHandler) GetInvoiceByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	invoice, err := h.invoiceService.GetInvoiceByIDForUser(c.Context(), uint(invoiceID), uint(userID))
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Invoice not found or access denied.", err)
	}
	return c.Status(fiber.StatusOK).JSON(invoice)
}

// ViewReceipt godoc
// @Summary View invoice receipt
// @Description View the receipt for a specific invoice
// @Tags invoices
// @Accept json
// @Produce image/jpeg,image/png,application/pdf
// @Param id path int true "Invoice ID"
// @Success 200 {file} file "Receipt file"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /invoices/{id}/viewreceipt [get]
func (h *InvoiceHandler) ViewReceipt(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	filePath, fileName, err := h.invoiceService.GetReceiptPathForUser(c.Context(), uint(invoiceID), uint(userID))
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Receipt not found or access denied.", err)
	}

	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	return c.SendFile(filePath)
}

// DownloadReceipt godoc
// @Summary Download invoice receipt
// @Description Download the receipt for a specific invoice
// @Tags invoices
// @Accept json
// @Produce application/octet-stream
// @Param id path int true "Invoice ID"
// @Success 200 {file} file "Receipt file"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /invoices/{id}/receipt [get]
func (h *InvoiceHandler) DownloadReceipt(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	filePath, fileName, err := h.invoiceService.GetReceiptPathForUser(c.Context(), uint(invoiceID), uint(userID))
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Receipt not found or access denied.", err)
	}

	return c.Download(filePath, fileName)
}

// SelectFinancialInstitution godoc
// @Summary Select financial institution for invoice
// @Description Select a financial institution to finance a specific invoice
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param request body dtos.SelectFinancialInstitutionRequest true "Financial institution selection details"
// @Success 200 {object} dtos.InvoiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /invoices/{id}/select-financial-institution [post]
func (h *InvoiceHandler) SelectFinancialInstitution(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	var req dtos.SelectFinancialInstitutionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Validation failed.", err)
	}

	invoice, err := h.invoiceService.SelectFinancialInstitution(c.Context(), uint(invoiceID), uint(userID), req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to select financial institution.", err)
	}

	return c.Status(fiber.StatusOK).JSON(invoice)
}

// GetSuggestedFinancialInstitutions godoc
// @Summary Get suggested financial institutions for invoice
// @Description Get a list of suggested financial institutions for a specific invoice based on compatibility
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} dtos.SuggestedFinancialInstitutionsResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /invoices/{id}/suggested-financial-institutions [get]
func (h *InvoiceHandler) GetSuggestedFinancialInstitutions(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	suggestions, err := h.invoiceService.GetSuggestedFinancialInstitutions(c.Context(), uint(invoiceID), uint(userID))
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to get suggested financial institutions.", err)
	}

	return c.Status(fiber.StatusOK).JSON(suggestions)
}
