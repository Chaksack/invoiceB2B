package handlers

import (
	"strconv"
	"time"

	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

type ReportingHandler struct {
	reportingService services.ReportingService
	validate         *validator.Validate
}

func NewReportingHandler(reportingService services.ReportingService, validate *validator.Validate) *ReportingHandler {
	return &ReportingHandler{
		reportingService: reportingService,
		validate:         validate,
	}
}

// @Summary Generate a report
// @Description Generate a comprehensive report based on the specified type and filters
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param request body dtos.ReportRequest true "Report generation request"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/generate [post]
func (h *ReportingHandler) GenerateReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	var request dtos.ReportRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Generate report
	reportData, err := h.reportingService.GenerateReport(&request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Report generated successfully",
		Data:    reportData,
	})
}

// @Summary Generate loan application report
// @Description Generate a comprehensive loan application report with statistics and trends
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Filter by status"
// @Param user_id query int false "Filter by user ID"
// @Param group_by query string false "Group by period" Enums(daily, weekly, monthly, yearly)
// @Success 200 {object} dtos.LoanApplicationReportResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/loan-applications [get]
func (h *ReportingHandler) GenerateLoanApplicationReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	request, err := h.parseReportRequest(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request parameters", err)
	}
	request.ReportType = "loan_applications"

	// Generate loan application report
	reportData, err := h.reportingService.GenerateLoanApplicationReport(request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate loan application report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Loan application report generated successfully",
		Data:    reportData,
	})
}

// @Summary Generate invoice report
// @Description Generate a comprehensive invoice report with statistics and trends
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Filter by status"
// @Param user_id query int false "Filter by user ID"
// @Param group_by query string false "Group by period" Enums(daily, weekly, monthly, yearly)
// @Success 200 {object} dtos.InvoiceReportResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/invoices [get]
func (h *ReportingHandler) GenerateInvoiceReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	request, err := h.parseReportRequest(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request parameters", err)
	}
	request.ReportType = "invoices"

	// Generate invoice report
	reportData, err := h.reportingService.GenerateInvoiceReport(request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate invoice report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Invoice report generated successfully",
		Data:    reportData,
	})
}

// @Summary Generate user report
// @Description Generate a comprehensive user report with statistics and trends
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param group_by query string false "Group by period" Enums(daily, weekly, monthly, yearly)
// @Success 200 {object} dtos.UserReportResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/users [get]
func (h *ReportingHandler) GenerateUserReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	request, err := h.parseReportRequest(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request parameters", err)
	}
	request.ReportType = "users"

	// Generate user report
	reportData, err := h.reportingService.GenerateUserReport(request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate user report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "User report generated successfully",
		Data:    reportData,
	})
}

// @Summary Generate activity report
// @Description Generate a comprehensive activity report with statistics and trends
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param group_by query string false "Group by period" Enums(daily, weekly, monthly, yearly)
// @Success 200 {object} dtos.ActivityReportResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/activity [get]
func (h *ReportingHandler) GenerateActivityReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	request, err := h.parseReportRequest(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request parameters", err)
	}
	request.ReportType = "activity"

	// Generate activity report
	reportData, err := h.reportingService.GenerateActivityReport(request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate activity report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Activity report generated successfully",
		Data:    reportData,
	})
}

// @Summary Generate financial institution report
// @Description Generate a comprehensive financial institution report with statistics
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dtos.FinancialInstitutionReportResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/financial-institutions [get]
func (h *ReportingHandler) GenerateFinancialInstitutionReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	request, err := h.parseReportRequest(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request parameters", err)
	}
	request.ReportType = "financial_institutions"

	// Generate financial institution report
	reportData, err := h.reportingService.GenerateFinancialInstitutionReport(request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate financial institution report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Financial institution report generated successfully",
		Data:    reportData,
	})
}

// @Summary Generate overview report
// @Description Generate a comprehensive overview report with all system statistics
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} dtos.OverviewReportResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/overview [get]
func (h *ReportingHandler) GenerateOverviewReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	request, err := h.parseReportRequest(c)
	if err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request parameters", err)
	}
	request.ReportType = "overview"

	// Generate overview report
	reportData, err := h.reportingService.GenerateOverviewReport(request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate overview report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Overview report generated successfully",
		Data:    reportData,
	})
}

// @Summary Export report
// @Description Export a report in the specified format (JSON, CSV, PDF)
// @Tags Admin - Reports
// @Accept json
// @Produce json
// @Param request body dtos.ReportRequest true "Report export request"
// @Success 200 {object} dtos.ExportReportResponse
// @Failure 400 {object} utils.SwaggerErrorResponse
// @Failure 401 {object} utils.SwaggerErrorResponse
// @Failure 403 {object} utils.SwaggerErrorResponse
// @Failure 500 {object} utils.SwaggerErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/reports/export [post]
func (h *ReportingHandler) ExportReport(c *fiber.Ctx) error {
	// Check admin role
	userRole, ok := c.Locals("user_role").(string)
	if !ok || userRole != "admin" {
		return utils.HandleError(c, fiber.StatusForbidden, "Admin access required", nil)
	}

	var request dtos.ReportRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(&request); err != nil {
		return utils.HandleValidationError(c, err)
	}

	// Generate report data first
	reportData, err := h.reportingService.GenerateReport(&request)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate report data", err)
	}

	// Export report
	exportResponse, err := h.reportingService.ExportReport(&request, reportData)
	if err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to export report", err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse{
		Status:  "success",
		Message: "Report exported successfully",
		Data:    exportResponse,
	})
}

// parseReportRequest parses query parameters into a ReportRequest
func (h *ReportingHandler) parseReportRequest(c *fiber.Ctx) (*dtos.ReportRequest, error) {
	request := &dtos.ReportRequest{}

	// Parse start date
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, err
		}
		request.StartDate = &startDate
	}

	// Parse end date
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, err
		}
		request.EndDate = &endDate
	}

	// Parse other parameters
	request.Status = c.Query("status")
	request.GroupBy = c.Query("group_by")

	// Parse user ID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return nil, err
		}
		userIDUint := uint(userID)
		request.UserID = &userIDUint
	}

	return request, nil
}