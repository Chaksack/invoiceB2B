package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"invoiceB2B/internal/models"
	"log"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *models.Invoice) error
	Update(ctx context.Context, invoice *models.Invoice) error
	FindByID(ctx context.Context, id uint) (*models.Invoice, error)
	FindByIDWithRelations(ctx context.Context, id uint) (*models.Invoice, error) // e.g. Preload User, Staff
	FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Invoice, int64, error)
	FindAll(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.Invoice, int64, error)
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *models.Invoice) error {
	if err := r.db.WithContext(ctx).Create(invoice).Error; err != nil {
		log.Printf("Error creating invoice in DB: %v", err)
		return err
	}
	return nil
}

func (r *invoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
	if err := r.db.WithContext(ctx).Save(invoice).Error; err != nil {
		log.Printf("Error updating invoice %d in DB: %v", invoice.ID, err)
		return err
	}
	return nil
}

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
	// Add more filters as needed (e.g., date range, invoice_number like)

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
		log.Printf("Error fetching all invoices with filters: %v", err)
		return nil, 0, err
	}
	return invoices, total, nil
}
