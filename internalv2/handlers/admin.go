package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"invoiceB2B/internal@v2/database"
	"invoiceB2B/internal@v2/middleware"
	"invoiceB2B/internal@v2/models"
)

// AdminHandler handles admin-related requests
type AdminHandler struct {
	db *database.DB
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(db *database.DB) *AdminHandler {
	return &AdminHandler{
		db: db,
	}
}

// GetBusinesses handles getting all businesses
func (h *AdminHandler) GetBusinesses(c *fiber.Ctx) error {
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

	// Build query with filters
	whereClause := "WHERE 1=1"
	queryParams := []interface{}{}
	paramIndex := 1

	if status != "" {
		whereClause += " AND b.status = $" + strconv.Itoa(paramIndex)
		queryParams = append(queryParams, status)
		paramIndex++
	}

	if search != "" {
		whereClause += " AND b.company_name ILIKE $" + strconv.Itoa(paramIndex)
		queryParams = append(queryParams, "%"+search+"%")
		paramIndex++
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) as total FROM businesses b " + whereClause
	err := h.db.QueryRow(countQuery, queryParams...).Scan(&total)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	totalPages := (total + limit - 1) / limit

	// Get businesses with user info
	businessesQuery := `
		SELECT 
			b.*,
			u.email,
			COUNT(i.id) as total_invoices,
			COALESCE(SUM(CASE WHEN i.status = 'funded' THEN i.amount ELSE 0 END), 0) as total_funded
		FROM businesses b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN invoices i ON b.id = i.business_id
		` + whereClause + `
		GROUP BY b.id, u.email
		ORDER BY b.created_at DESC
		LIMIT $` + strconv.Itoa(paramIndex) + ` OFFSET $` + strconv.Itoa(paramIndex+1)

	queryParams = append(queryParams, limit, offset)

	rows, err := h.db.Query(businessesQuery, queryParams...)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}
	defer rows.Close()

	var businesses []models.BusinessWithUser
	for rows.Next() {
		var business models.BusinessWithUser
		err := rows.Scan(
			&business.ID, &business.UserID, &business.CompanyName, &business.Industry,
			&business.AnnualRevenue, &business.EmployeeCount, &business.Status,
			&business.CreatedAt, &business.UpdatedAt, &business.Email,
			&business.TotalInvoices, &business.TotalFunded,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Database error")
		}
		businesses = append(businesses, business)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Businesses retrieved successfully",
		"data":    businesses,
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

// GetDashboardSummary handles getting admin dashboard summary
func (h *AdminHandler) GetDashboardSummary(c *fiber.Ctx) error {
	// Get total businesses
	var totalBusinesses int
	err := h.db.QueryRow("SELECT COUNT(*) FROM businesses").Scan(&totalBusinesses)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Get pending businesses
	var pendingBusinesses int
	err = h.db.QueryRow("SELECT COUNT(*) FROM businesses WHERE status = 'pending'").Scan(&pendingBusinesses)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Get total invoices
	var totalInvoices int
	err = h.db.QueryRow("SELECT COUNT(*) FROM invoices WHERE deleted_at IS NULL").Scan(&totalInvoices)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Get total funded amount
	var totalFunded float64
	err = h.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM invoices WHERE status = 'funded' AND deleted_at IS NULL").Scan(&totalFunded)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Get recent invoices
	rows, err := h.db.Query(`
		SELECT i.id, i.invoice_number, i.customer_name, i.amount, i.status, i.created_at, b.company_name
		FROM invoices i
		JOIN businesses b ON i.business_id = b.id
		WHERE i.deleted_at IS NULL
		ORDER BY i.created_at DESC
		LIMIT 5
	`)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}
	defer rows.Close()

	type RecentInvoice struct {
		ID            string  `json:"id"`
		InvoiceNumber string  `json:"invoice_number"`
		CustomerName  string  `json:"customer_name"`
		Amount        float64 `json:"amount"`
		Status        string  `json:"status"`
		CreatedAt     string  `json:"created_at"`
		CompanyName   string  `json:"company_name"`
	}

	var recentInvoices []RecentInvoice
	for rows.Next() {
		var invoice RecentInvoice
		err := rows.Scan(
			&invoice.ID, &invoice.InvoiceNumber, &invoice.CustomerName,
			&invoice.Amount, &invoice.Status, &invoice.CreatedAt, &invoice.CompanyName,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Database error")
		}
		recentInvoices = append(recentInvoices, invoice)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Dashboard summary retrieved successfully",
		"data": fiber.Map{
			"total_businesses":     totalBusinesses,
			"pending_businesses":   pendingBusinesses,
			"total_invoices":       totalInvoices,
			"total_funded_amount":  totalFunded,
			"recent_invoices":      recentInvoices,
		},
		"timestamp": fiber.Now(),
	})
}

// UpdateBusinessStatus handles updating business status
func (h *AdminHandler) UpdateBusinessStatus(c *fiber.Ctx) error {
	businessID := c.Params("id")
	if businessID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Business ID is required")
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=approved rejected suspended"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Update business status
	result, err := h.db.Exec(
		"UPDATE businesses SET status = $1, updated_at = NOW() WHERE id = $2",
		req.Status, businessID,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Business not found")
	}

	// If approved, also approve the user
	if req.Status == "approved" {
		_, err = h.db.Exec(
			"UPDATE users SET is_approved = true WHERE id = (SELECT user_id FROM businesses WHERE id = $1)",
			businessID,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Database error")
		}
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "Business status updated successfully",
		"timestamp": fiber.Now(),
	})
} 