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

// AdminHandler handles HTTP requests for admin-related operations.
type AdminHandler struct {
	adminService services.AdminService
	fileService  services.FileService // fileService is used for receipt uploads
	validate     *validator.Validate
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(adminService services.AdminService, fileService services.FileService, validate *validator.Validate) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		fileService:  fileService,
		validate:     validate,
	}
}

// GetAdminProfile godoc
// @Summary Get admin profile
// @Description Get the profile of the currently logged-in admin/staff member
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} dtos.StaffResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/profile/me [get]
func (h *AdminHandler) GetAdminProfile(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64) // JWT numbers are often float64
	if !ok {
		// Try string conversion if float64 fails, as some JWT libraries might store it as string
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	staffProfile, err := h.adminService.GetStaffByID(c.Context(), adminStaffID) // Assuming GetStaffByID is added to AdminService
	if err != nil {
		if strings.Contains(err.Error(), "staff not found") { // Or use errors.Is if a specific error type is defined
			return utils.HandleError(c, fiber.StatusNotFound, "Admin profile not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve admin profile.", err)
	}

	return c.Status(fiber.StatusOK).JSON(staffProfile)
}

// --- Admin User & KYC Management ---

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a paginated list of all users
// @Tags admin-users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Success 200 {object} map[string]interface{} "Returns users, total, page, pageSize"
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users [get]
func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	// TODO: Implement actual filtering based on query parameters
	// For example:
	// emailFilter := c.Query("email")
	// companyFilter := c.Query("companyName")
	// kycStatusFilter := c.Query("kycStatus")
	// Pass these filters to a modified adminService.GetAllUsers if that service method supports them.

	users, total, err := h.adminService.GetAllUsers(c.Context(), page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve users.", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users, "total": total, "page": page, "pageSize": pageSize})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get details of a specific user by their ID
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dtos.UserResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUserByID(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format.", err)
	}

	user, err := h.adminService.GetUserByID(c.Context(), uint(userID))
	if err != nil {
		// Differentiate between not found and other errors
		if strings.Contains(err.Error(), "user not found") { // Or use errors.Is if ErrUserNotFound is exported by service
			return utils.HandleError(c, fiber.StatusNotFound, "User not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve user.", err)
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// GetUserKYCDetail godoc
// @Summary Get user KYC details
// @Description Get KYC (Know Your Customer) details for a specific user
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dtos.AdminKYCDetailResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users/{id}/kyc [get]
func (h *AdminHandler) GetUserKYCDetail(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format.", err)
	}

	kycDetail, err := h.adminService.GetUserKYCDetail(c.Context(), uint(userID))
	if err != nil {
		// Differentiate between KYC not found (which might mean not submitted) and other errors
		if strings.Contains(err.Error(), "kyc record not found") { // Or use errors.Is if ErrKYCNotFound is exported
			return utils.HandleError(c, fiber.StatusNotFound, "KYC details not found for user.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve KYC details.", err)
	}
	return c.Status(fiber.StatusOK).JSON(kycDetail)
}

// ReviewKYC godoc
// @Summary Review user KYC application
// @Description Approve or reject a user's KYC (Know Your Customer) application
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dtos.AdminKYCReviewRequest true "KYC review details"
// @Success 200 {object} dtos.AdminKYCDetailResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users/{id}/kyc/review [put]
func (h *AdminHandler) ReviewKYC(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	reviewerStaffIDFloat, ok := claims["user_id"].(float64) // JWT numbers are often float64
	if !ok {
		reviewerStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(reviewerStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		reviewerStaffIDFloat = parsedID
	}
	reviewerStaffID := uint(reviewerStaffIDFloat)

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
	// Ensure req.Status is a valid models.KYCStatus if it's coming as a string
	// This validation might be better handled in the service layer or with custom validator
	if req.Status == "rejected" && (req.RejectionReason == nil || strings.TrimSpace(*req.RejectionReason) == "") {
		return utils.HandleError(c, fiber.StatusBadRequest, "Rejection reason is required when rejecting KYC.", nil)
	}

	updatedKYC, err := h.adminService.ReviewKYC(c.Context(), uint(userID), reviewerStaffID, req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update KYC status.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedKYC)
}

// --- Admin Invoice Management ---

// GetAllInvoices godoc
// @Summary Get all invoices
// @Description Get a paginated list of all invoices with optional status filter
// @Tags admin-invoices
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Param status query string false "Filter by invoice status"
// @Success 200 {object} dtos.InvoiceListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/invoices [get]
func (h *AdminHandler) GetAllInvoices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
	statusFilter := c.Query("status") // Example filter

	// The adminService.GetAllInvoices expects statusFilter as a string, not a map.
	invoices, total, err := h.adminService.GetAllInvoices(c.Context(), page, pageSize, statusFilter)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve invoices.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.InvoiceListResponse{
		Invoices: invoices, Total: total, Page: page, PageSize: pageSize,
	})
}

// GetInvoiceDetail godoc
// @Summary Get invoice details
// @Description Get detailed information for a specific invoice
// @Tags admin-invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Success 200 {object} dtos.InvoiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/invoices/{id} [get]
func (h *AdminHandler) GetInvoiceDetail(c *fiber.Ctx) error {
	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}
	invoice, err := h.adminService.GetInvoiceDetail(c.Context(), uint(invoiceID))
	if err != nil {
		if strings.Contains(err.Error(), "invoice not found") { // Or use errors.Is
			return utils.HandleError(c, fiber.StatusNotFound, "Invoice not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve invoice details.", err)
	}
	return c.Status(fiber.StatusOK).JSON(invoice)
}

// UpdateInvoiceStatus godoc
// @Summary Update invoice status
// @Description Update the status of an invoice (approve, reject, disburse, etc.)
// @Tags admin-invoices
// @Accept json
// @Produce json
// @Param id path int true "Invoice ID"
// @Param request body dtos.AdminInvoiceUpdateRequest true "Invoice status update details"
// @Success 200 {object} dtos.InvoiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/invoices/{id}/status [put]
func (h *AdminHandler) UpdateInvoiceStatus(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

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

	updatedInvoice, err := h.adminService.UpdateInvoiceStatus(c.Context(), uint(invoiceID), adminStaffID, req)
	if err != nil {
		// Provide more specific error messages based on service errors if possible
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update invoice status.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedInvoice)
}

// UploadDisbursementReceipt godoc
// @Summary Upload disbursement receipt
// @Description Upload a disbursement receipt for an invoice
// @Tags admin-invoices
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Invoice ID"
// @Param receiptFile formData file true "Receipt file to upload (PDF, PNG, JPG, JPEG)"
// @Success 200 {object} dtos.InvoiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/invoices/{id}/receipt [post]
func (h *AdminHandler) UploadDisbursementReceipt(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)
	
	// Extract user role, default to "staff" if not present
	userRole, ok := claims["role"].(string)
	if !ok {
		userRole = "staff" // Default role for admin staff if not specified in token
	}

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	file, err := c.FormFile("receiptFile")
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Receipt file is required.", err)
	}

	// Basic validation for file extension and size
	allowedExtensions := map[string]bool{".pdf": true, ".png": true, ".jpg": true, ".jpeg": true} // Added image types
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid file type. Only PDF, PNG, JPG, JPEG allowed.", nil)
	}
	
	// Use the standard ValidateFileSize method (this doesn't need access control)
	if err := h.fileService.ValidateFileSize(file.Size); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, err.Error(), err)
	}

	req := dtos.AdminUploadReceiptRequest{File: file} // File is *multipart.FileHeader
	
	// Convert adminStaffID to string for access control
	adminStaffIDStr := fmt.Sprintf("%d", adminStaffID)

	// Pass the user ID and role to the AdminService for access control
	updatedInvoice, err := h.adminService.UploadDisbursementReceipt(c.Context(), uint(invoiceID), adminStaffID, req, adminStaffIDStr, userRole)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to upload receipt.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedInvoice)
}

// DownloadInvoicePDF godoc
// @Summary Download invoice PDF
// @Description Download an invoice as a PDF document
// @Tags admin-invoices
// @Accept json
// @Produce application/json
// @Param id path int true "Invoice ID"
// @Success 200 {object} services.InvoicePDFResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/invoices/{id}/download-pdf [get]
func (h *AdminHandler) DownloadInvoicePDF(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	invoiceIDStr := c.Params("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid invoice ID format.", err)
	}

	pdfResponse, err := h.adminService.DownloadInvoicePDF(c.Context(), uint(invoiceID), adminStaffID)
	if err != nil {
		if strings.Contains(err.Error(), "failed to generate PDF") { // Or use errors.Is
			return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate invoice PDF.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Error preparing invoice PDF for download.", err)
	}

	// Assuming pdfResponse.FilePath is an absolute path to the file on the server
	// or a relative path that fileService can resolve to an absolute one.
	// For direct download, you'd typically use c.SendFile() or c.Download()
	// c.Download() sets Content-Disposition for a nice filename.
	// If pdfResponse.FilePath is relative, construct absolute path:
	// absPath, err := h.fileService.GetAbsPath(pdfResponse.FilePath)
	// if err != nil {
	//    return utils.HandleError(c, fiber.StatusInternalServerError, "PDF file path error.", err)
	// }
	// return c.Download(absPath, pdfResponse.FileName)

	// For now, returning JSON with path info, client can construct download URL or request separately
	return c.Status(fiber.StatusOK).JSON(pdfResponse)
}

// --- Admin Staff Management ---

// CreateStaff godoc
// @Summary Create staff member
// @Description Create a new staff member (admin, reviewer, etc.)
// @Tags admin-staff
// @Accept json
// @Produce json
// @Param request body dtos.CreateStaffRequest true "Staff creation details"
// @Success 201 {object} dtos.StaffResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse "Staff with this email already exists"
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/staff [post]
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
		if strings.Contains(err.Error(), "already exists") {
			return utils.HandleError(c, fiber.StatusConflict, "Staff with this email already exists.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to create staff.", err)
	}
	return c.Status(fiber.StatusCreated).JSON(staff)
}

// GetAllStaff godoc
// @Summary Get all staff members
// @Description Get a paginated list of all staff members
// @Tags admin-staff
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Success 200 {object} map[string]interface{} "Returns staff, total, page, pageSize"
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/staff [get]
func (h *AdminHandler) GetAllStaff(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	staffList, total, err := h.adminService.GetAllStaff(c.Context(), page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve staff.", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"staff": staffList, "total": total, "page": page, "pageSize": pageSize})
}

// UpdateStaff godoc
// @Summary Update staff member
// @Description Update details of an existing staff member
// @Tags admin-staff
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Param request body dtos.UpdateStaffRequest true "Staff update details"
// @Success 200 {object} dtos.StaffResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/staff/{id} [put]
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
	// Note: UpdateStaffRequest might have all fields as pointers to distinguish
	// between empty value and not provided. Validation might need to be adjusted.
	// if errs := h.validate.Struct(req); errs != nil {
	//    return utils.HandleValidationError(c, errs)
	// }
	updatedStaff, err := h.adminService.UpdateStaff(c.Context(), uint(staffID), req)
	if err != nil {
		if strings.Contains(err.Error(), "staff not found") { // Or use errors.Is
			return utils.HandleError(c, fiber.StatusNotFound, "Staff member not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update staff.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedStaff)
}

// DeleteStaff godoc
// @Summary Delete staff member
// @Description Delete an existing staff member
// @Tags admin-staff
// @Accept json
// @Produce json
// @Param id path int true "Staff ID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/staff/{id} [delete]
func (h *AdminHandler) DeleteStaff(c *fiber.Ctx) error {
	staffIDStr := c.Params("id")
	staffID, err := strconv.ParseUint(staffIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid staff ID format.", err)
	}
	if err := h.adminService.DeleteStaff(c.Context(), uint(staffID)); err != nil {
		if strings.Contains(err.Error(), "not found") { // Or use errors.Is
			return utils.HandleError(c, fiber.StatusNotFound, "Staff member not found for deletion.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to delete staff.", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Admin Activity Logs & Analytics ---

// GetActivityLogs godoc
// @Summary Get activity logs
// @Description Get a paginated list of all activity logs with optional filters
// @Tags admin-logs
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 20)"
// @Param user_id query string false "Filter by user ID"
// @Param staff_id query string false "Filter by staff ID"
// @Param action query string false "Filter by action type"
// @Success 200 {object} dtos.ActivityLogListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/activity-logs [get]
func (h *AdminHandler) GetActivityLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	filters := make(map[string]string)
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if staffID := c.Query("staff_id"); staffID != "" {
		filters["staff_id"] = staffID
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	// Add more filters as needed, e.g., date range

	logs, total, err := h.adminService.GetActivityLogs(c.Context(), page, pageSize, filters)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve activity logs.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.ActivityLogListResponse{
		Logs: logs, Total: total, Page: page, PageSize: pageSize,
	})
}

// GetUserActivityLogs godoc
// @Summary Get user activity logs
// @Description Get activity logs specific to a user
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 20)"
// @Param action query string false "Filter logs by action type"
// @Success 200 {object} dtos.ActivityLogListResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/users/{id}/activity-logs [get]
func (h *AdminHandler) GetUserActivityLogs(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid user ID format.", err)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	filters := make(map[string]string)
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	// user_id is passed as a direct argument to the service, not in the filters map for this specific function.
	// However, the service interface `GetUserActivityLogs` expects filters map[string]string,
	// so we pass the filters map which might contain other filters like 'action'.
	// The service implementation will handle the userID separately and merge with other filters if needed.

	logs, total, err := h.adminService.GetUserActivityLogs(c.Context(), uint(userID), page, pageSize, filters)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "User not found, cannot retrieve activity logs.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve user activity logs.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.ActivityLogListResponse{
		Logs: logs, Total: total, Page: page, PageSize: pageSize,
	})
}

// GetAdminDashboardAnalytics godoc
// @Summary Get dashboard analytics
// @Description Get aggregated data for the admin dashboard (users, invoices, transactions)
// @Tags admin-analytics
// @Accept json
// @Produce json
// @Success 200 {object} services.AdminDashboardAnalytics
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/dashboard/analytics [get]
func (h *AdminHandler) GetAdminDashboardAnalytics(c *fiber.Ctx) error {
	analytics, err := h.adminService.GetAdminDashboardAnalytics(c.Context())
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve dashboard analytics.", err)
	}
	return c.Status(fiber.StatusOK).JSON(analytics)
}

// --- Financial Institution Management ---

// CreateFinancialInstitution godoc
// @Summary Create financial institution
// @Description Create a new financial institution
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param request body dtos.CreateFinancialInstitutionRequest true "Financial institution creation details"
// @Success 201 {object} dtos.FinancialInstitutionResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse "Financial institution with this code already exists"
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions [post]
func (h *AdminHandler) CreateFinancialInstitution(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	var req dtos.CreateFinancialInstitutionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	fi, err := h.adminService.CreateFinancialInstitution(c.Context(), adminStaffID, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return utils.HandleError(c, fiber.StatusConflict, "Financial institution with this code already exists.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to create financial institution.", err)
	}
	return c.Status(fiber.StatusCreated).JSON(fi)
}

// GetAllFinancialInstitutions godoc
// @Summary Get all financial institutions
// @Description Get a paginated list of all financial institutions with optional filters
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Param name query string false "Filter by name"
// @Param code query string false "Filter by code"
// @Param is_active query string false "Filter by active status (true/false)"
// @Success 200 {object} dtos.FinancialInstitutionListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions [get]
func (h *AdminHandler) GetAllFinancialInstitutions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	filters := make(map[string]string)
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if code := c.Query("code"); code != "" {
		filters["code"] = code
	}
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive
	}

	financialInstitutions, total, err := h.adminService.GetAllFinancialInstitutions(c.Context(), page, pageSize, filters)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve financial institutions.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.FinancialInstitutionListResponse{
		FinancialInstitutions: financialInstitutions, Total: total, Page: page, PageSize: pageSize,
	})
}

// GetFinancialInstitutionByID godoc
// @Summary Get financial institution by ID
// @Description Get details of a specific financial institution by its ID
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Financial Institution ID"
// @Success 200 {object} dtos.FinancialInstitutionResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/{id} [get]
func (h *AdminHandler) GetFinancialInstitutionByID(c *fiber.Ctx) error {
	fiIDStr := c.Params("id")
	fiID, err := strconv.ParseUint(fiIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid financial institution ID format.", err)
	}

	fi, err := h.adminService.GetFinancialInstitutionByID(c.Context(), uint(fiID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve financial institution.", err)
	}
	return c.Status(fiber.StatusOK).JSON(fi)
}

// UpdateFinancialInstitution godoc
// @Summary Update financial institution
// @Description Update an existing financial institution
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Financial Institution ID"
// @Param request body dtos.UpdateFinancialInstitutionRequest true "Financial institution update details"
// @Success 200 {object} dtos.FinancialInstitutionResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse "Financial institution with this code already exists"
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/{id} [put]
func (h *AdminHandler) UpdateFinancialInstitution(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	fiIDStr := c.Params("id")
	fiID, err := strconv.ParseUint(fiIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid financial institution ID format.", err)
	}

	var req dtos.UpdateFinancialInstitutionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	fi, err := h.adminService.UpdateFinancialInstitution(c.Context(), uint(fiID), adminStaffID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution not found.", err)
		}
		if strings.Contains(err.Error(), "already exists") {
			return utils.HandleError(c, fiber.StatusConflict, "Financial institution with this code already exists.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update financial institution.", err)
	}
	return c.Status(fiber.StatusOK).JSON(fi)
}

// DeleteFinancialInstitution godoc
// @Summary Delete financial institution
// @Description Delete an existing financial institution
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Financial Institution ID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/{id} [delete]
func (h *AdminHandler) DeleteFinancialInstitution(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	fiIDStr := c.Params("id")
	fiID, err := strconv.ParseUint(fiIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid financial institution ID format.", err)
	}

	if err := h.adminService.DeleteFinancialInstitution(c.Context(), uint(fiID), adminStaffID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to delete financial institution.", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Financial Institution Product Management ---

// CreateFinancialInstitutionProduct godoc
// @Summary Create financial institution product
// @Description Create a new product for a financial institution
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param request body dtos.CreateFinancialInstitutionProductRequest true "Financial institution product creation details"
// @Success 201 {object} dtos.FinancialInstitutionProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse "Financial institution not found"
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/products [post]
func (h *AdminHandler) CreateFinancialInstitutionProduct(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	var req dtos.CreateFinancialInstitutionProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	product, err := h.adminService.CreateFinancialInstitutionProduct(c.Context(), adminStaffID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to create financial institution product.", err)
	}
	return c.Status(fiber.StatusCreated).JSON(product)
}

// GetFinancialInstitutionProductByID godoc
// @Summary Get financial institution product by ID
// @Description Get details of a specific financial institution product by its ID
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dtos.FinancialInstitutionProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/products/{id} [get]
func (h *AdminHandler) GetFinancialInstitutionProductByID(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid product ID format.", err)
	}

	product, err := h.adminService.GetFinancialInstitutionProductByID(c.Context(), uint(productID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution product not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve financial institution product.", err)
	}
	return c.Status(fiber.StatusOK).JSON(product)
}

// GetFinancialInstitutionProducts godoc
// @Summary Get financial institution products
// @Description Get a paginated list of products for a specific financial institution
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param fiId path int true "Financial Institution ID"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Success 200 {object} dtos.FinancialInstitutionProductListResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/{fiId}/products [get]
func (h *AdminHandler) GetFinancialInstitutionProducts(c *fiber.Ctx) error {
	fiIDStr := c.Params("fiId")
	fiID, err := strconv.ParseUint(fiIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid financial institution ID format.", err)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	products, total, err := h.adminService.GetFinancialInstitutionProducts(c.Context(), uint(fiID), page, pageSize)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve financial institution products.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.FinancialInstitutionProductListResponse{
		Products: products, Total: total, Page: page, PageSize: pageSize,
	})
}

// UpdateFinancialInstitutionProduct godoc
// @Summary Update financial institution product
// @Description Update an existing financial institution product
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body dtos.UpdateFinancialInstitutionProductRequest true "Product update details"
// @Success 200 {object} dtos.FinancialInstitutionProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/products/{id} [put]
func (h *AdminHandler) UpdateFinancialInstitutionProduct(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	productIDStr := c.Params("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid product ID format.", err)
	}

	var req dtos.UpdateFinancialInstitutionProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	product, err := h.adminService.UpdateFinancialInstitutionProduct(c.Context(), uint(productID), adminStaffID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution product not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update financial institution product.", err)
	}
	return c.Status(fiber.StatusOK).JSON(product)
}

// DeleteFinancialInstitutionProduct godoc
// @Summary Delete financial institution product
// @Description Delete a financial institution product
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/products/{id} [delete]
func (h *AdminHandler) DeleteFinancialInstitutionProduct(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	productIDStr := c.Params("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid product ID format.", err)
	}

	if err := h.adminService.DeleteFinancialInstitutionProduct(c.Context(), uint(productID), adminStaffID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution product not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to delete financial institution product.", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Financial Institution Term Management ---

// CreateFinancialInstitutionTerm creates a new financial institution term.
func (h *AdminHandler) CreateFinancialInstitutionTerm(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	var req dtos.CreateFinancialInstitutionTermRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	term, err := h.adminService.CreateFinancialInstitutionTerm(c.Context(), adminStaffID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution or product not found.", err)
		}
		if strings.Contains(err.Error(), "invalid") {
			return utils.HandleError(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to create financial institution term.", err)
	}
	return c.Status(fiber.StatusCreated).JSON(term)
}

// GetFinancialInstitutionTermByID retrieves a specific financial institution term by its ID.
func (h *AdminHandler) GetFinancialInstitutionTermByID(c *fiber.Ctx) error {
	termIDStr := c.Params("id")
	termID, err := strconv.ParseUint(termIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid term ID format.", err)
	}

	term, err := h.adminService.GetFinancialInstitutionTermByID(c.Context(), uint(termID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution term not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve financial institution term.", err)
	}
	return c.Status(fiber.StatusOK).JSON(term)
}

// GetFinancialInstitutionTerms godoc
// @Summary Get financial institution terms
// @Description Get a paginated list of terms for a specific financial institution
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param fiId path int true "Financial Institution ID"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Success 200 {object} dtos.FinancialInstitutionTermListResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/{fiId}/terms [get]
func (h *AdminHandler) GetFinancialInstitutionTerms(c *fiber.Ctx) error {
	fiIDStr := c.Params("fiId")
	fiID, err := strconv.ParseUint(fiIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid financial institution ID format.", err)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	terms, total, err := h.adminService.GetFinancialInstitutionTerms(c.Context(), uint(fiID), page, pageSize)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve financial institution terms.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.FinancialInstitutionTermListResponse{
		Terms: terms, Total: total, Page: page, PageSize: pageSize,
	})
}

// GetFinancialInstitutionTermsByProduct godoc
// @Summary Get terms by product
// @Description Get a paginated list of terms for a specific financial institution product
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Success 200 {object} dtos.FinancialInstitutionTermListResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/products/{productId}/terms [get]
func (h *AdminHandler) GetFinancialInstitutionTermsByProduct(c *fiber.Ctx) error {
	productIDStr := c.Params("productId")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid product ID format.", err)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	terms, total, err := h.adminService.GetFinancialInstitutionTermsByProduct(c.Context(), uint(productID), page, pageSize)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution product not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve financial institution terms by product.", err)
	}
	return c.Status(fiber.StatusOK).JSON(dtos.FinancialInstitutionTermListResponse{
		Terms: terms, Total: total, Page: page, PageSize: pageSize,
	})
}

// UpdateFinancialInstitutionTerm godoc
// @Summary Update financial institution term
// @Description Update an existing financial institution term
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Term ID"
// @Param request body dtos.UpdateFinancialInstitutionTermRequest true "Term update details"
// @Success 200 {object} dtos.FinancialInstitutionTermResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/terms/{id} [put]
func (h *AdminHandler) UpdateFinancialInstitutionTerm(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	termIDStr := c.Params("id")
	termID, err := strconv.ParseUint(termIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid term ID format.", err)
	}

	var req dtos.UpdateFinancialInstitutionTermRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body.", err)
	}
	if errs := h.validate.Struct(req); errs != nil {
		return utils.HandleValidationError(c, errs)
	}

	term, err := h.adminService.UpdateFinancialInstitutionTerm(c.Context(), uint(termID), adminStaffID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution term not found.", err)
		}
		if strings.Contains(err.Error(), "invalid") {
			return utils.HandleError(c, fiber.StatusBadRequest, err.Error(), err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to update financial institution term.", err)
	}
	return c.Status(fiber.StatusOK).JSON(term)
}

// DeleteFinancialInstitutionTerm godoc
// @Summary Delete financial institution term
// @Description Delete a financial institution term
// @Tags admin-financial-institutions
// @Accept json
// @Produce json
// @Param id path int true "Term ID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/financial-institutions/terms/{id} [delete]
func (h *AdminHandler) DeleteFinancialInstitutionTerm(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims.", nil)
	}
	adminStaffIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		adminStaffIDStr, okStr := claims["user_id"].(string)
		if !okStr {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID type in token.", nil)
		}
		parsedID, err := strconv.ParseFloat(adminStaffIDStr, 64)
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid staff ID format in token string.", err)
		}
		adminStaffIDFloat = parsedID
	}
	adminStaffID := uint(adminStaffIDFloat)

	termIDStr := c.Params("id")
	termID, err := strconv.ParseUint(termIDStr, 10, 64)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid term ID format.", err)
	}

	if err := h.adminService.DeleteFinancialInstitutionTerm(c.Context(), uint(termID), adminStaffID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return utils.HandleError(c, fiber.StatusNotFound, "Financial institution term not found.", err)
		}
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to delete financial institution term.", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
