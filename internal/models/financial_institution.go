package models

import (
	"gorm.io/gorm"
	"time"
)

// FinancialInstitution represents a financial institution that can provide financing for invoices
type FinancialInstitution struct {
	gorm.Model
	Name              string  `gorm:"type:varchar(255);not null"`
	Code              string  `gorm:"type:varchar(50);uniqueIndex;not null"`
	Description       string  `gorm:"type:text;null"`
	InterestRateMin   float64 `gorm:"type:decimal(5,4);null"`
	InterestRateMax   float64 `gorm:"type:decimal(5,4);null"`
	ProcessingFee     float64 `gorm:"type:decimal(10,2);null"`
	MinInvoiceAmount  float64 `gorm:"type:decimal(15,2);null"`
	MaxInvoiceAmount  float64 `gorm:"type:decimal(15,2);null"`
	TermsDays         int     `gorm:"null"`
	IsActive          bool    `gorm:"default:true"`
	
	// Relationships
	Invoices []Invoice `gorm:"foreignKey:FinancialInstitutionID"`
}

// FinancialInstitutionProduct represents a product offered by a financial institution
type FinancialInstitutionProduct struct {
	gorm.Model
	FinancialInstitutionID uint                 `gorm:"not null;index"`
	FinancialInstitution   FinancialInstitution `gorm:"foreignKey:FinancialInstitutionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name                   string               `gorm:"type:varchar(255);not null"`
	Description            string               `gorm:"type:text;null"`
	InterestRateMin        float64              `gorm:"type:decimal(5,4);null"`
	InterestRateMax        float64              `gorm:"type:decimal(5,4);null"`
	ProcessingFee          float64              `gorm:"type:decimal(10,2);null"`
	MinInvoiceAmount       float64              `gorm:"type:decimal(15,2);null"`
	MaxInvoiceAmount       float64              `gorm:"type:decimal(15,2);null"`
	TermsDays              int                  `gorm:"null"`
	IsActive               bool                 `gorm:"default:true"`
}

// FinancialInstitutionTerm represents specific terms offered by a financial institution
type FinancialInstitutionTerm struct {
	gorm.Model
	FinancialInstitutionID      uint                 `gorm:"not null;index"`
	FinancialInstitution        FinancialInstitution `gorm:"foreignKey:FinancialInstitutionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FinancialInstitutionProductID *uint                `gorm:"null;index"`
	FinancialInstitutionProduct  *FinancialInstitutionProduct `gorm:"foreignKey:FinancialInstitutionProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name                        string               `gorm:"type:varchar(255);not null"`
	Description                 string               `gorm:"type:text;null"`
	InterestRate                float64              `gorm:"type:decimal(5,4);not null"`
	ProcessingFee               float64              `gorm:"type:decimal(10,2);null"`
	TermDays                    int                  `gorm:"not null"`
	MinInvoiceAmount            float64              `gorm:"type:decimal(15,2);null"`
	MaxInvoiceAmount            float64              `gorm:"type:decimal(15,2);null"`
	IsActive                    bool                 `gorm:"default:true"`
	ValidFrom                   *time.Time           `gorm:"null"`
	ValidUntil                  *time.Time           `gorm:"null"`
}