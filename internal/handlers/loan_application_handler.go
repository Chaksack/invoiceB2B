package handlers

import (
	"strconv"
	"time"

	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

type LoanApplicationHandler struct {
	loanAppService services.LoanApplicationService
	fileService    services.FileService
	validate       *validator.Validate
}

func NewLoanApplicationHandler(loanAppService services.LoanApplicationService, fileService services.FileService, validate *validator.Validate) *LoanApplicationHandler {
	return &LoanApplicationHandler{
		loanAppService: loanAppService,
		fileService:    fileService,
		validate:       validate,
	}
}

// @Summary Create a new loan application
// @Description Create a new loan application with basic details
// @Tags Loan Applications
// @Accept json
// @Produce json
// @Param request body dtos.CreateLoanApplicationRequest true "Loan application request"
// @Success 201 {object} dtos.LoanApplicationResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications [post]
func (h *LoanApplicationHandler) CreateLoanApplication(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	var request dtos.CreateLoanApplicationRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Create loan application
	loanApp, err := h.loanAppService.CreateLoanApplication(userID, &request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to create loan application", err)
	}

	response := dtos.ToLoanApplicationResponse(loanApp)
	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Loan application created successfully",
		Data:    response,
	})
}

// @Summary Create a manual loan application
// @Description Create a loan application with all KYB and financial data provided manually
// @Tags Loan Applications
// @Accept json
// @Produce json
// @Param request body dtos.ManualLoanInputRequest true "Manual loan application request"
// @Success 201 {object} dtos.LoanApplicationResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications/manual [post]
func (h *LoanApplicationHandler) CreateManualLoanApplication(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	var request dtos.ManualLoanInputRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Create manual loan application
	loanApp, err := h.loanAppService.CreateManualLoanApplication(userID, &request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to create manual loan application", err)
	}

	response := dtos.ToLoanApplicationResponse(loanApp)
	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Manual loan application created successfully",
		Data:    response,
	})
}

// @Summary Get loan application by ID
// @Description Retrieve a specific loan application by its ID
// @Tags Loan Applications
// @Accept json
// @Produce json
// @Param id path int true "Loan Application ID"
// @Success 200 {object} dtos.LoanApplicationResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications/{id} [get]
func (h *LoanApplicationHandler) GetLoanApplication(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid loan application ID", err)
	}

	loanApp, err := h.loanAppService.GetLoanApplication(uint(id), userID)
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Loan application not found", err)
	}

	response := dtos.ToLoanApplicationResponse(loanApp)
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Loan application retrieved successfully",
		Data:    response,
	})
}

// @Summary Get user's loan applications
// @Description Retrieve all loan applications for the authenticated user with pagination
// @Tags Loan Applications
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dtos.LoanApplicationListResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications [get]
func (h *LoanApplicationHandler) GetUserLoanApplications(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	applications, totalCount, err := h.loanAppService.GetLoanApplicationsByUser(userID, page, pageSize)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve loan applications", err)
	}

	var responses []dtos.LoanApplicationResponse
	for _, app := range applications {
		responses = append(responses, dtos.ToLoanApplicationResponse(&app))
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	response := dtos.LoanApplicationListResponse{
		Applications: responses,
		TotalCount:   totalCount,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Loan applications retrieved successfully",
		Data:    response,
	})
}

// @Summary Submit KYB information
// @Description Submit Know Your Business information for a loan application
// @Tags Loan Applications
// @Accept json
// @Produce json
// @Param id path int true "Loan Application ID"
// @Param request body dtos.SubmitKYBInformationRequest true "KYB information request"
// @Success 200 {object} dtos.KYBInformationResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications/{id}/kyb [post]
func (h *LoanApplicationHandler) SubmitKYBInformation(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	idStr := c.Params("id")
	loanAppID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid loan application ID", err)
	}

	// Verify loan application belongs to user
	_, err = h.loanAppService.GetLoanApplication(uint(loanAppID), userID)
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Loan application not found", err)
	}

	var request dtos.SubmitKYBInformationRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Submit KYB information
	kybInfo, err := h.loanAppService.SubmitKYBInformation(uint(loanAppID), &request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to submit KYB information", err)
	}

	response := dtos.ToKYBInformationResponse(kybInfo)
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "KYB information submitted successfully",
		Data:    response,
	})
}

// @Summary Submit financial statement
// @Description Submit financial statement data for a loan application
// @Tags Loan Applications
// @Accept json
// @Produce json
// @Param id path int true "Loan Application ID"
// @Param request body dtos.SubmitFinancialStatementRequest true "Financial statement request"
// @Success 200 {object} dtos.FinancialStatementResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications/{id}/financial-statements [post]
func (h *LoanApplicationHandler) SubmitFinancialStatement(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	idStr := c.Params("id")
	loanAppID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid loan application ID", err)
	}

	// Verify loan application belongs to user
	_, err = h.loanAppService.GetLoanApplication(uint(loanAppID), userID)
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Loan application not found", err)
	}

	var request dtos.SubmitFinancialStatementRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Submit financial statement
	statement, err := h.loanAppService.SubmitFinancialStatement(uint(loanAppID), &request, "")
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to submit financial statement", err)
	}

	response := dtos.ToFinancialStatementResponse(statement)
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Financial statement submitted successfully",
		Data:    response,
	})
}

// @Summary Upload financial statement file
// @Description Upload a financial statement file for a loan application
// @Tags Loan Applications
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Loan Application ID"
// @Param file formData file true "Financial statement file"
// @Param type formData string true "Statement type (bank or mobile_money)"
// @Param provider_name formData string true "Bank or mobile money provider name"
// @Param account_number formData string true "Account number"
// @Success 200 {object} dtos.DocumentUploadResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications/{id}/financial-statements/upload [post]
func (h *LoanApplicationHandler) UploadFinancialStatementFile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	idStr := c.Params("id")
	loanAppID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid loan application ID", err)
	}

	// Verify loan application belongs to user
	_, err = h.loanAppService.GetLoanApplication(uint(loanAppID), userID)
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Loan application not found", err)
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "File is required", err)
	}

	// Validate file
	if err := h.fileService.ValidateFileType(file); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid file type", err)
	}

	if err := h.fileService.ValidateFileSize(file.Size); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "File size exceeds limit", err)
	}

	// Save file
	fileName, filePath, err := h.fileService.SaveFile(file, "financial-statements")
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to save file", err)
	}

	response := dtos.DocumentUploadResponse{
		FileName:   fileName,
		FilePath:   filePath,
		FileSize:   file.Size,
		UploadedAt: time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Financial statement file uploaded successfully",
		Data:    response,
	})
}

// @Summary Upload invoice or contract document
// @Description Upload an invoice or contract document for loan application processing
// @Tags Loan Applications
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Loan Application ID"
// @Param file formData file true "Invoice or contract file"
// @Param source formData string true "Document source (invoice or contract)"
// @Success 200 {object} dtos.DocumentUploadResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/loan-applications/{id}/documents [post]
func (h *LoanApplicationHandler) UploadDocument(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "User not authenticated", nil)
	}

	idStr := c.Params("id")
	loanAppID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid loan application ID", err)
	}

	// Verify loan application belongs to user
	_, err = h.loanAppService.GetLoanApplication(uint(loanAppID), userID)
	if err != nil {
		return utils.HandleError(c, fiber.StatusNotFound, "Loan application not found", err)
	}

	// Get source parameter
	source := c.FormValue("source")
	if source != "invoice" && source != "contract" {
		return utils.HandleError(c, fiber.StatusBadRequest, "Source must be 'invoice' or 'contract'", nil)
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "File is required", err)
	}

	// Validate file
	if err := h.fileService.ValidateFileType(file); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid file type", err)
	}

	if err := h.fileService.ValidateFileSize(file.Size); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "File size exceeds limit", err)
	}

	// Save file
	fileName, filePath, err := h.fileService.SaveFile(file, "loan-documents")
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to save file", err)
	}

	// Process document upload
	err = h.loanAppService.ProcessDocumentUpload(uint(loanAppID), filePath, models.LoanApplicationSource(source))
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to process document", err)
	}

	response := dtos.DocumentUploadResponse{
		FileName:   fileName,
		FilePath:   filePath,
		FileSize:   file.Size,
		UploadedAt: time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Document uploaded and processing initiated",
		Data:    response,
	})
}

// Admin endpoints

// @Summary Get loan application statistics (Admin)
// @Description Retrieve loan application statistics for admin dashboard
// @Tags Admin - Loan Applications
// @Accept json
// @Produce json
// @Success 200 {object} dtos.LoanApplicationStatsResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/loan-applications/stats [get]
func (h *LoanApplicationHandler) GetLoanApplicationStats(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	stats, err := h.loanAppService.GetLoanApplicationStats()
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to get statistics", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Statistics retrieved successfully",
		Data:    stats,
	})
}

// @Summary Admin review loan application
// @Description Admin review and approve/reject a loan application
// @Tags Admin - Loan Applications
// @Accept json
// @Produce json
// @Param id path int true "Loan Application ID"
// @Param request body dtos.AdminReviewRequest true "Admin review request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/loan-applications/{id}/review [post]
func (h *LoanApplicationHandler) AdminReviewLoanApplication(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	adminID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Admin not authenticated", nil)
	}

	idStr := c.Params("id")
	loanAppID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid loan application ID", err)
	}

	var request dtos.AdminReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Process admin review
	err = h.loanAppService.AdminReview(uint(loanAppID), adminID, &request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to process admin review", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Loan application reviewed successfully",
		Data:    nil,
	})
}

// @Summary Get all loan applications (Admin)
// @Description Retrieve all loan applications with pagination for admin
// @Tags Admin - Loan Applications
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} dtos.LoanApplicationListResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/loan-applications [get]
func (h *LoanApplicationHandler) GetAllLoanApplications(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get all applications (admin can see all)
	applications, totalCount, err := h.loanAppService.GetLoanApplicationsByUser(0, page, pageSize) // 0 means all users for admin
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to retrieve loan applications", err)
	}

	var responses []dtos.LoanApplicationResponse
	for _, app := range applications {
		responses = append(responses, dtos.ToLoanApplicationResponse(&app))
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	response := dtos.LoanApplicationListResponse{
		Applications: responses,
		TotalCount:   totalCount,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Loan applications retrieved successfully",
		Data:    response,
	})
}

// @Summary Send loan application to financial institution
// @Description Admin sends loan application data (both raw and processed) to a financial institution
// @Tags Admin - Loan Applications
// @Accept json
// @Produce json
// @Param id path int true "Loan Application ID"
// @Param request body dtos.SendToFinancialInstitutionRequest true "Send to financial institution request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 404 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/loan-applications/{id}/send-to-financial-institution [post]
func (h *LoanApplicationHandler) SendToFinancialInstitution(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	adminID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.HandleError(c, fiber.StatusUnauthorized, "Admin not authenticated", nil)
	}

	idStr := c.Params("id")
	loanAppID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid loan application ID", err)
	}

	var request dtos.SendToFinancialInstitutionRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Send loan application to financial institution
	err = h.loanAppService.SendToFinancialInstitution(uint(loanAppID), adminID, &request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to send loan application to financial institution", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Loan application sent to financial institution successfully",
		Data:    nil,
	})
}