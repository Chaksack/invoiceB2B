package handlers

import (
	"fmt"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"
	"path/filepath"
	"strconv"

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

func (h *InvoiceHandler) UploadInvoice(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	file, err := c.FormFile("invoiceFile")
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invoice file is required.", err)
	}

	allowedExtensions := map[string]bool{".pdf": true, ".csv": true}
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid file type. Only PDF and CSV are allowed.", nil)
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
