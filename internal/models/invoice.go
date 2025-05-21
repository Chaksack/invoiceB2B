package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// InvoiceStatus type for invoice status
type InvoiceStatus string

const (
	InvoicePendingReview    InvoiceStatus = "pending_review"
	InvoiceApproved         InvoiceStatus = "approved"
	InvoiceRejected         InvoiceStatus = "rejected"
	InvoiceDisbursed        InvoiceStatus = "disbursed"
	InvoiceRepaymentPending InvoiceStatus = "repayment_pending"
	InvoiceRepaid           InvoiceStatus = "repaid"
)

// Invoice represents an invoice uploaded by a user
type Invoice struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	InvoiceNumber     string `gorm:"type:varchar(100);not null"`
	IssuerName        string `gorm:"type:varchar(255);null"`
	IssuerBankAccount string `gorm:"type:varchar(100);null"`
	IssuerBankName    string `gorm:"type:varchar(100);null"`
	DebtorName        string `gorm:"type:varchar(255);null"`

	Amount   float64 `gorm:"type:decimal(18,2);not null"`
	Currency string  `gorm:"type:varchar(3);not null"`
	DueDate  *time.Time

	Status           InvoiceStatus `gorm:"type:varchar(30);default:'pending_review';not null"`
	OriginalFilePath string        `gorm:"type:varchar(500);null"`
	JSONData         string        `gorm:"type:jsonb;null"`

	UploadedAt time.Time `gorm:"autoCreateTime"`

	ApprovedByID *uuid.UUID `gorm:"type:uuid;null"`
	ApprovedBy   *Staff     `gorm:"foreignKey:ApprovedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ApprovedAt   *time.Time

	DisbursedByID *uuid.UUID `gorm:"type:uuid;null"`
	DisbursedBy   *Staff     `gorm:"foreignKey:DisbursedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DisbursedAt   *time.Time

	FinancingFeePercentage  *float64 `gorm:"type:decimal(5,2);null"`
	FinancedAmount          *float64 `gorm:"type:decimal(18,2);null"`
	DisbursementReceiptPath *string  `gorm:"type:varchar(500);null"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Transactions []Transaction `gorm:"foreignKey:InvoiceID"`
}

func (inv *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	return
}

type TransactionType string

const (
	TransactionDisbursement TransactionType = "disbursement"
	TransactionRepayment    TransactionType = "repayment"
)

type Transaction struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	InvoiceID uuid.UUID `gorm:"type:uuid;not null;index"`
	Invoice   Invoice   `gorm:"foreignKey:InvoiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Type            TransactionType `gorm:"type:varchar(20);not null"`
	Amount          float64         `gorm:"type:decimal(18,2);not null"`
	TransactionDate time.Time       `gorm:"not null"`
	ReferenceID     *string         `gorm:"type:varchar(100);null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
