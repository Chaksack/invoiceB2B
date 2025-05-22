package handlers

import (
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

type AdminHandler struct {
	adminService services.AdminService
	fileService  services.FileService
	validate     *validator.Validate
}

func NewAdminHandler(adminService services.AdminService, fileService services.FileService, validate *validator.Validate) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		fileService:  fileService,
		validate:     validate,
	}
}

// --- Admin User & KYC Management ---

func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	// Add filters like email, companyName, kycStatus from c.Query()

	users, total, err := h.adminService.GetAllUsers(c.Context(), page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve users.", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users, "total": total, "page": page, "pageSize": pageSize})
}

func (h *AdminHandler) GetUserKYCDetail(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format.", err)
	}

	kycDetail, err := h.adminService.GetUserKYCDetail(c.Context(), uint(userID))
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "KYC details not found for user.", err)
	}
	return c.Status(fiber.StatusOK).JSON(kycDetail)
}

func (h *AdminHandler) ReviewKYC(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims) // Assuming admin is a staff member
	reviewerStaffIDStr := claims["user_id"].(string)               // Staff ID from token
	reviewerStaffID, _ := strconv.ParseUint(reviewerStaffIDStr, 10, 64)

	userIDStr := c.Params("id") // User whose KYC is being reviewed
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format.", err)
	}

	var req dtos.AdminKYCReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}
	if req.Status == "rejected" && (req.RejectionReason == nil || *req.RejectionReason == "") {
		return utils.HandleError(c, fiber.StatusBadRequest, "Rejection reason is required when rejecting KYC.", nil)
	}

	updatedKYC, err := h.adminService.ReviewKYC(c.Context(), uint(userID), uint(reviewerStaffID), req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update KYC status.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedKYC)
}

// --- Admin Invoice Management ---
func (h *AdminHandler) GetAllInvoices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	statusFilter := c.Query("status") // Example filter

	invoices, total, err := h.adminService.GetAllInvoices(c.Context(), page, pageSize, statusFilter)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve invoices.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.InvoiceListResponse{
		Invoices: invoices, Total: total, Page: page, PageSize: pageSize,
	})
}

func (h *AdminHandler) GetInvoiceDetail(c *fiber.Ctx) error {
	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}
	invoice, err := h.adminService.GetInvoiceDetail(c.Context(), uint(invoiceID))
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Invoice not found.", err)
	}
	return c.Status(fiber.StatusOK).JSON(invoice)
}

func (h *AdminHandler) UpdateInvoiceStatus(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	adminStaffIDStr := claims["user_id"].(string)
	adminStaffID, _ := strconv.ParseUint(adminStaffIDStr, 10, 64)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	var req dtos.AdminInvoiceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	updatedInvoice, err := h.adminService.UpdateInvoiceStatus(c.Context(), uint(invoiceID), uint(adminStaffID), req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update invoice status.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedInvoice)
}

func (h *AdminHandler) UploadDisbursementReceipt(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	adminStaffIDStr := claims["user_id"].(string)
	adminStaffID, _ := strconv.ParseUint(adminStaffIDStr, 10, 64)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	file, err := c.FormFile("receiptFile")
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Receipt file is required.", err)
	}
	// Basic validation
	allowedExtensions := map[string]bool{".pdf": true, ".csv": true}
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[strings.ToLower(ext)] {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid file type. Only PDF, CSV allowed.", nil)
	}
	if err := h.fileService.ValidateFileSize(file.Size); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, err.Error(), err)
	}

	req := dtos.AdminUploadReceiptRequest{File: file}

	updatedInvoice, err := h.adminService.UploadDisbursementReceipt(c.Context(), uint(invoiceID), uint(adminStaffID), req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to upload receipt.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedInvoice)
}

// --- Admin Staff Management ---
func (h *AdminHandler) CreateStaff(c *fiber.Ctx) error {
	var req dtos.CreateStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}
	staff, err := h.adminService.CreateStaff(c.Context(), req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to create staff.", err)
	}
	return c.Status(fiber.StatusCreated).JSON(staff)
}

func (h *AdminHandler) GetAllStaff(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	staffList, total, err := h.adminService.GetAllStaff(c.Context(), page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve staff.", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"staff": staffList, "total": total, "page": page, "pageSize": pageSize})
}

func (h *AdminHandler) UpdateStaff(c *fiber.Ctx) error {
	staffIDStr := c.Params("id")
	staffID, err := strconv.ParseUint(staffIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid staff ID format.", err)
	}
	var req dtos.UpdateStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}
	updatedStaff, err := h.adminService.UpdateStaff(c.Context(), uint(staffID), req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update staff.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedStaff)
}

func (h *AdminHandler) DeleteStaff(c *fiber.Ctx) error {
	staffIDStr := c.Params("id")
	staffID, err := strconv.ParseUint(staffIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid staff ID format.", err)
	}
	if err := h.adminService.DeleteStaff(c.Context(), uint(staffID)); err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to delete staff.", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Admin Activity Logs ---
func (h *AdminHandler) GetActivityLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))
	// Add filters: c.Query("user_id"), c.Query("staff_id"), c.Query("action")

	logs, total, err := h.adminService.GetActivityLogs(c.Context(), page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve activity logs.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.ActivityLogListResponse{
		Logs: logs, Total: total, Page: page, PageSize: pageSize,
	})
}
