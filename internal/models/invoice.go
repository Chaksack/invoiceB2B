package models

import (
	"gorm.io/gorm"
	"time"
)

type InvoiceStatus string

const (
	InvoicePendingReview    InvoiceStatus = "pending_review"
	InvoiceApproved         InvoiceStatus = "approved"
	InvoiceRejected         InvoiceStatus = "rejected"
	InvoiceDisbursed        InvoiceStatus = "disbursed"
	InvoiceRepaymentPending InvoiceStatus = "repayment_pending"
	InvoiceRepaid           InvoiceStatus = "repaid"
)

type Invoice struct {
	gorm.Model
	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	InvoiceNumber     string `gorm:"type:varchar(100);null"`
	IssuerName        string `gorm:"type:varchar(255);null"`
	IssuerBankAccount string `gorm:"type:varchar(100);null"`
	IssuerBankName    string `gorm:"type:varchar(100);null"`
	DebtorName        string `gorm:"type:varchar(255);null"`

	Amount   float64    `gorm:"type:decimal(18,2);null"`
	Currency string     `gorm:"type:varchar(3);null"`
	DueDate  *time.Time `gorm:"null"`

	Status           InvoiceStatus `gorm:"type:varchar(30);default:'pending_review';not null"`
	OriginalFilePath string        `gorm:"type:varchar(500);null"`
	JSONData         string        `gorm:"type:jsonb;null"`

	UploadedAt time.Time `gorm:"autoCreateTime"`

	ApprovedByID *uint  `gorm:"null"`
	ApprovedBy   *Staff `gorm:"foreignKey:ApprovedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ApprovedAt   *time.Time

	DisbursedByID *uint  `gorm:"null"`
	DisbursedBy   *Staff `gorm:"foreignKey:DisbursedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DisbursedAt   *time.Time

	FinancingFeePercentage  *float64 `gorm:"type:decimal(5,2);null"`
	FinancedAmount          *float64 `gorm:"type:decimal(18,2);null"`
	DisbursementReceiptPath *string  `gorm:"type:varchar(500);null"`
	ProcessingError         *string  `gorm:"type:text;null"`

	Transactions []Transaction `gorm:"foreignKey:InvoiceID"`
}

type TransactionType string

const (
	TransactionDisbursement TransactionType = "disbursement"
	TransactionRepayment    TransactionType = "repayment"
)

type Transaction struct {
	gorm.Model
	InvoiceID uint    `gorm:"not null;index"`
	Invoice   Invoice `gorm:"foreignKey:InvoiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Type            TransactionType `gorm:"type:varchar(20);not null"`
	Amount          float64         `gorm:"type:decimal(18,2);not null"`
	TransactionDate time.Time       `gorm:"not null"`
	ReferenceID     *string         `gorm:"type:varchar(100);null"`
}
