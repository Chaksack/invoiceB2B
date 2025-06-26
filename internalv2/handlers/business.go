package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"invoiceB2B/internal@v2/database"
	"invoiceB2B/internal@v2/middleware"
	"invoiceB2B/internal@v2/models"
)

// BusinessHandler handles business-related requests
type BusinessHandler struct {
	db *database.DB
}

// NewBusinessHandler creates a new business handler
func NewBusinessHandler(db *database.DB) *BusinessHandler {
	return &BusinessHandler{
		db: db,
	}
}

// GetInvoices handles getting business invoices
func (h *BusinessHandler) GetInvoices(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)

	// Get query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")
	search := c.Query("search")

	// Validate and limit parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get business ID
	var businessID string
	err := h.db.QueryRow("SELECT id FROM businesses WHERE user_id = $1", user.ID).Scan(&businessID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "Business profile not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Build query with filters
	whereClause := "WHERE business_id = $1 AND deleted_at IS NULL"
	queryParams := []interface{}{businessID}
	paramIndex := 2

	if status != "" {
		whereClause += " AND status = $" + strconv.Itoa(paramIndex)
		queryParams = append(queryParams, status)
		paramIndex++
	}

	if search != "" {
		whereClause += " AND (invoice_number ILIKE $" + strconv.Itoa(paramIndex) + " OR customer_name ILIKE $" + strconv.Itoa(paramIndex) + ")"
		queryParams = append(queryParams, "%"+search+"%")
		paramIndex++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) as total FROM invoices " + whereClause
	err = h.db.QueryRow(countQuery, queryParams...).Scan(&total)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	totalPages := (total + limit - 1) / limit

	// Get invoices
	invoicesQuery := "SELECT * FROM invoices " + whereClause + " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(paramIndex) + " OFFSET $" + strconv.Itoa(paramIndex+1)
	queryParams = append(queryParams, limit, offset)

	rows, err := h.db.Query(invoicesQuery, queryParams...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(
			&invoice.ID, &invoice.BusinessID, &invoice.InvoiceNumber, &invoice.CustomerName,
			&invoice.Amount, &invoice.Currency, &invoice.Status, &invoice.IssueDate,
			&invoice.DueDate, &invoice.Description, &invoice.FileURL, &invoice.CreatedAt,
			&invoice.UpdatedAt, &invoice.DeletedAt,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Database error")
		}
		invoices = append(invoices, invoice)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Invoices retrieved successfully",
		"data":    invoices,
		"pagination": models.Pagination{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPages: totalPages,
			HasNext:   page < totalPages,
			HasPrev:   page > 1,
		},
		"timestamp": fiber.Now(),
	})
}

// GetProfile handles getting business profile
func (h *BusinessHandler) GetProfile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)

	var business models.Business
	err := h.db.QueryRow(
		"SELECT id, user_id, company_name, industry, annual_revenue, employee_count, status, created_at, updated_at FROM businesses WHERE user_id = $1",
		user.ID,
	).Scan(
		&business.ID, &business.UserID, &business.CompanyName, &business.Industry,
		&business.AnnualRevenue, &business.EmployeeCount, &business.Status,
		&business.CreatedAt, &business.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "Business profile not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Business profile retrieved successfully",
		"data":    business,
		"timestamp": fiber.Now(),
	})
}

// UpdateProfile handles updating business profile
func (h *BusinessHandler) UpdateProfile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)

	var req models.BusinessUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Update business profile
	result, err := h.db.Exec(
		"UPDATE businesses SET company_name = $1, industry = $2, annual_revenue = $3, employee_count = $4, updated_at = NOW() WHERE user_id = $5",
		req.CompanyName, req.Industry, req.AnnualRevenue, req.EmployeeCount, user.ID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Business profile not found")
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "Business profile updated successfully",
		"timestamp": fiber.Now(),
	})
} 