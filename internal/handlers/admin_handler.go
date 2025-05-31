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

// GetAdminProfile retrieves the profile of the currently logged-in admin/staff member.
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

// GetAllUsers retrieves a paginated list of all users.
// Query parameters: page, pageSize, email, companyName, kycStatus (conceptual)
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

// GetUserByID retrieves a specific user by their ID.
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

// GetUserKYCDetail retrieves KYC details for a specific user.
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

// ReviewKYC allows an admin to approve or reject a user's KYC application.
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

// GetAllInvoices retrieves a paginated list of all invoices, with optional filters.
// Query parameters: page, pageSize, status
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

// GetInvoiceDetail retrieves details for a specific invoice.
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

// UpdateInvoiceStatus allows an admin to update the status of an invoice.
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

// UploadDisbursementReceipt allows an admin to upload a disbursement receipt for an invoice.
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
	if err := h.fileService.ValidateFileSize(file.Size); err != nil { // Assuming fileService has ValidateFileSize
		return utils.HandleError(c, fiber.StatusBadRequest, err.Error(), err)
	}

	req := dtos.AdminUploadReceiptRequest{File: file} // File is *multipart.FileHeader

	updatedInvoice, err := h.adminService.UploadDisbursementReceipt(c.Context(), uint(invoiceID), adminStaffID, req)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to upload receipt.", err)
	}
	return c.Status(fiber.StatusOK).JSON(updatedInvoice)
}

// DownloadInvoicePDF allows an admin to download an invoice as a PDF.
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

// CreateStaff allows an admin to create a new staff member.
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

// GetAllStaff retrieves a paginated list of all staff members.
func (h *AdminHandler) GetAllStaff(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	staffList, total, err := h.adminService.GetAllStaff(c.Context(), page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve staff.", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"staff": staffList, "total": total, "page": page, "pageSize": pageSize})
}

// UpdateStaff allows an admin to update details of an existing staff member.
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

// DeleteStaff allows an admin to delete a staff member.
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

// GetActivityLogs retrieves a paginated list of all activity logs, with optional filters.
// Query parameters: page, pageSize, user_id, staff_id, action
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

// GetUserActivityLogs retrieves activity logs specific to a user.
// Path parameter: id (userID)
// Query parameters: page, pageSize, action
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

// GetAdminDashboardAnalytics retrieves aggregated data for the admin dashboard.
func (h *AdminHandler) GetAdminDashboardAnalytics(c *fiber.Ctx) error {
	analytics, err := h.adminService.GetAdminDashboardAnalytics(c.Context())
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve dashboard analytics.", err)
	}
	return c.Status(fiber.StatusOK).JSON(analytics)
}
