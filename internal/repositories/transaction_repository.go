package repositories

import (
	"context"
	"invoiceB2B/internal/models"
	"log"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	FindByInvoiceID(ctx context.Context, invoiceID uint) ([]models.Transaction, error)
	// Add other methods as needed, e.g., FindByID
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		log.Printf("Error creating transaction in DB: %v", err)
		return err
	}
	return nil
}

func (r *transactionRepository) FindByInvoiceID(ctx context.Context, invoiceID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).Order("created_at DESC").Find(&transactions).Error; err != nil {
		log.Printf("Error finding transactions by invoice ID %d: %v", invoiceID, err)
		return nil, err
	}
	return transactions, nil
}
