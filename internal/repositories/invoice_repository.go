package repositories

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"invoiceB2B/internal/models"
	"log"
	"time" // Added for CountOverdue
)

// InvoiceRepository interface defines the operations for invoice data management.
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *models.Invoice) error
	Update(ctx context.Context, invoice *models.Invoice) error
	FindByID(ctx context.Context, id uint) (*models.Invoice, error)
	FindByIDWithRelations(ctx context.Context, id uint) (*models.Invoice, error)
	FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Invoice, int64, error)
	// Renamed from FindAll to FindAllWithRelations to match service layer usage and make preloading explicit.
	FindAll(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.Invoice, int64, error)
	// Added methods for analytics
	CountByStatus(ctx context.Context, status models.InvoiceStatus, filters map[string]interface{}) (int64, error)
	CountOverdue(ctx context.Context) (int64, error)
	SumAmountByStatus(ctx context.Context, statuses []models.InvoiceStatus, amountField string) (float64, error)
}

type invoiceRepository struct {
	db *gorm.DB
}

// NewInvoiceRepository creates a new instance of InvoiceRepository.
func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

// Create persists a new invoice record to the database.
func (r *invoiceRepository) Create(ctx context.Context, invoice *models.Invoice) error {
	if err := r.db.WithContext(ctx).Create(invoice).Error; err != nil {
		log.Printf("Error creating invoice in DB: %v", err)
		return err
	}
	return nil
}

// Update saves changes to an existing invoice record in the database.
func (r *invoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
	if err := r.db.WithContext(ctx).Save(invoice).Error; err != nil {
		log.Printf("Error updating invoice %d in DB: %v", invoice.ID, err)
		return err
	}
	return nil
}

// FindByID retrieves a single invoice by its ID.
func (r *invoiceRepository) FindByID(ctx context.Context, id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.WithContext(ctx).First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		log.Printf("Error finding invoice by ID %d in DB: %v", id, err)
		return nil, err
	}
	return &invoice, nil
}

// FindByIDWithRelations retrieves a single invoice by its ID, preloading specified related entities.
func (r *invoiceRepository) FindByIDWithRelations(ctx context.Context, id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("ApprovedBy").
		Preload("DisbursedBy").
		First(&invoice, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		log.Printf("Error finding invoice by ID %d with relations in DB: %v", id, err)
		return nil, err
	}
	return &invoice, nil
}

// FindByUserID retrieves a paginated list of invoices for a specific user.
func (r *invoiceRepository) FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Invoice{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting invoices for user %d: %v", userID, err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&invoices).Error; err != nil {
		log.Printf("Error fetching invoices for user %d: %v", userID, err)
		return nil, 0, err
	}
	return invoices, total, nil
}

// FindAllWithRelations retrieves a paginated list of all invoices, applying filters and preloading relations.
// Renamed from FindAll.
func (r *invoiceRepository) FindAll(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Invoice{})

	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if userID, ok := filters["user_id"]; ok && userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting all invoices with filters: %v", err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Preload("User").
		Preload("ApprovedBy").
		Preload("DisbursedBy").
		Offset(offset).Limit(pageSize).Find(&invoices).Error
	if err != nil {
		log.Printf("Error fetching all invoices with filters and relations: %v", err)
		return nil, 0, err
	}
	return invoices, total, nil
}

// CountByStatus counts invoices based on their status.
// The filters map can be used for additional criteria (e.g., user_id).
// If status is an empty string, it counts all invoices matching other filters.
func (r *invoiceRepository) CountByStatus(ctx context.Context, status models.InvoiceStatus, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Invoice{})

	if status != "" { // If status is provided, filter by it
		query = query.Where("status = ?", status)
	}

	if filters != nil { // Apply additional filters if any
		for key, value := range filters {
			if vStr, ok := value.(string); ok && vStr == "" { // Skip empty string filters
				continue
			}
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		log.Printf("Error counting invoices by status '%s' with filters %v: %v", status, filters, err)
		return 0, err
	}
	return count, nil
}

// CountOverdue counts invoices that are past their due date and have status 'disbursed'.
func (r *invoiceRepository) CountOverdue(ctx context.Context) (int64, error) {
	var count int64
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.Invoice{}).
		Where("status = ?", models.InvoiceDisbursed). // Assuming models.InvoiceDisbursed is defined
		Where("due_date < ?", now).
		Count(&count).Error
	if err != nil {
		log.Printf("Error counting overdue invoices: %v", err)
		return 0, err
	}
	return count, nil
}

// SumAmountByStatus calculates the sum of a specified amount field for invoices matching given statuses.
// statuses: A slice of InvoiceStatus to filter by. If empty, no status filter is applied.
// amountField: The database column name of the amount to sum (e.g., "amount", "financed_amount").
func (r *invoiceRepository) SumAmountByStatus(ctx context.Context, statuses []models.InvoiceStatus, amountField string) (float64, error) {
	var totalAmount struct {
		Total float64
	}

	if amountField != "amount" && amountField != "financed_amount" {
		log.Printf("Error: Invalid amount field '%s' for summation.", amountField)
		return 0, errors.New("invalid amount field for summation")
	}

	query := r.db.WithContext(ctx).Model(&models.Invoice{})

	if len(statuses) > 0 {
		query = query.Where("status IN (?)", statuses)
	}

	selectQuery := fmt.Sprintf("COALESCE(SUM(%s), 0) as total", amountField)
	if err := query.Select(selectQuery).Scan(&totalAmount).Error; err != nil {
		log.Printf("Error summing '%s' for statuses %v: %v", amountField, statuses, err)
		return 0, err
	}
	return totalAmount.Total, nil
}
