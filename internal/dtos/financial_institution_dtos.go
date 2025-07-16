package dtos

import (
	"time"
)

// Financial Institution DTOs

// CreateFinancialInstitutionRequest is used to create a new financial institution
type CreateFinancialInstitutionRequest struct {
	Name             string   `json:"name" validate:"required,min=2"`
	Code             string   `json:"code" validate:"required,min=2"`
	Description      string   `json:"description,omitempty"`
	InterestRateMin  *float64 `json:"interestRateMin,omitempty" validate:"omitempty,gte=0,lte=1"`
	InterestRateMax  *float64 `json:"interestRateMax,omitempty" validate:"omitempty,gte=0,lte=1"`
	ProcessingFee    *float64 `json:"processingFee,omitempty" validate:"omitempty,gte=0"`
	MinInvoiceAmount *float64 `json:"minInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	MaxInvoiceAmount *float64 `json:"maxInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	TermsDays        *int     `json:"termsDays,omitempty" validate:"omitempty,gte=0"`
	IsActive         *bool    `json:"isActive,omitempty"`
}

// UpdateFinancialInstitutionRequest is used to update an existing financial institution
type UpdateFinancialInstitutionRequest struct {
	Name             *string  `json:"name,omitempty" validate:"omitempty,min=2"`
	Code             *string  `json:"code,omitempty" validate:"omitempty,min=2"`
	Description      *string  `json:"description,omitempty"`
	InterestRateMin  *float64 `json:"interestRateMin,omitempty" validate:"omitempty,gte=0,lte=1"`
	InterestRateMax  *float64 `json:"interestRateMax,omitempty" validate:"omitempty,gte=0,lte=1"`
	ProcessingFee    *float64 `json:"processingFee,omitempty" validate:"omitempty,gte=0"`
	MinInvoiceAmount *float64 `json:"minInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	MaxInvoiceAmount *float64 `json:"maxInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	TermsDays        *int     `json:"termsDays,omitempty" validate:"omitempty,gte=0"`
	IsActive         *bool    `json:"isActive,omitempty"`
}

// FinancialInstitutionResponse is used to return financial institution data
type FinancialInstitutionResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Code             string    `json:"code"`
	Description      string    `json:"description,omitempty"`
	InterestRateMin  float64   `json:"interestRateMin,omitempty"`
	InterestRateMax  float64   `json:"interestRateMax,omitempty"`
	ProcessingFee    float64   `json:"processingFee,omitempty"`
	MinInvoiceAmount float64   `json:"minInvoiceAmount,omitempty"`
	MaxInvoiceAmount float64   `json:"maxInvoiceAmount,omitempty"`
	TermsDays        int       `json:"termsDays,omitempty"`
	IsActive         bool      `json:"isActive"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// FinancialInstitutionListResponse is used to return a list of financial institutions
type FinancialInstitutionListResponse struct {
	FinancialInstitutions []FinancialInstitutionResponse `json:"financialInstitutions"`
	Total                 int64                          `json:"total"`
	Page                  int                            `json:"page"`
	PageSize              int                            `json:"pageSize"`
}

// Financial Institution Product DTOs

// CreateFinancialInstitutionProductRequest is used to create a new financial institution product
type CreateFinancialInstitutionProductRequest struct {
	FinancialInstitutionID uint     `json:"financialInstitutionId" validate:"required"`
	Name                   string   `json:"name" validate:"required,min=2"`
	Description            string   `json:"description,omitempty"`
	InterestRateMin        *float64 `json:"interestRateMin,omitempty" validate:"omitempty,gte=0,lte=1"`
	InterestRateMax        *float64 `json:"interestRateMax,omitempty" validate:"omitempty,gte=0,lte=1"`
	ProcessingFee          *float64 `json:"processingFee,omitempty" validate:"omitempty,gte=0"`
	MinInvoiceAmount       *float64 `json:"minInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	MaxInvoiceAmount       *float64 `json:"maxInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	TermsDays              *int     `json:"termsDays,omitempty" validate:"omitempty,gte=0"`
	IsActive               *bool    `json:"isActive,omitempty"`
}

// UpdateFinancialInstitutionProductRequest is used to update an existing financial institution product
type UpdateFinancialInstitutionProductRequest struct {
	Name             *string  `json:"name,omitempty" validate:"omitempty,min=2"`
	Description      *string  `json:"description,omitempty"`
	InterestRateMin  *float64 `json:"interestRateMin,omitempty" validate:"omitempty,gte=0,lte=1"`
	InterestRateMax  *float64 `json:"interestRateMax,omitempty" validate:"omitempty,gte=0,lte=1"`
	ProcessingFee    *float64 `json:"processingFee,omitempty" validate:"omitempty,gte=0"`
	MinInvoiceAmount *float64 `json:"minInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	MaxInvoiceAmount *float64 `json:"maxInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	TermsDays        *int     `json:"termsDays,omitempty" validate:"omitempty,gte=0"`
	IsActive         *bool    `json:"isActive,omitempty"`
}

// FinancialInstitutionProductResponse is used to return financial institution product data
type FinancialInstitutionProductResponse struct {
	ID                     uint      `json:"id"`
	FinancialInstitutionID uint      `json:"financialInstitutionId"`
	Name                   string    `json:"name"`
	Description            string    `json:"description,omitempty"`
	InterestRateMin        float64   `json:"interestRateMin,omitempty"`
	InterestRateMax        float64   `json:"interestRateMax,omitempty"`
	ProcessingFee          float64   `json:"processingFee,omitempty"`
	MinInvoiceAmount       float64   `json:"minInvoiceAmount,omitempty"`
	MaxInvoiceAmount       float64   `json:"maxInvoiceAmount,omitempty"`
	TermsDays              int       `json:"termsDays,omitempty"`
	IsActive               bool      `json:"isActive"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

// FinancialInstitutionProductListResponse is used to return a list of financial institution products
type FinancialInstitutionProductListResponse struct {
	Products []FinancialInstitutionProductResponse `json:"products"`
	Total    int64                                 `json:"total"`
	Page     int                                   `json:"page"`
	PageSize int                                   `json:"pageSize"`
}

// Financial Institution Term DTOs

// CreateFinancialInstitutionTermRequest is used to create a new financial institution term
type CreateFinancialInstitutionTermRequest struct {
	FinancialInstitutionID      uint      `json:"financialInstitutionId" validate:"required"`
	FinancialInstitutionProductID *uint     `json:"financialInstitutionProductId,omitempty"`
	Name                        string    `json:"name" validate:"required,min=2"`
	Description                 string    `json:"description,omitempty"`
	InterestRate                float64   `json:"interestRate" validate:"required,gte=0,lte=1"`
	ProcessingFee               *float64  `json:"processingFee,omitempty" validate:"omitempty,gte=0"`
	TermDays                    int       `json:"termDays" validate:"required,gte=0"`
	MinInvoiceAmount            *float64  `json:"minInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	MaxInvoiceAmount            *float64  `json:"maxInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	IsActive                    *bool     `json:"isActive,omitempty"`
	ValidFrom                   *string   `json:"validFrom,omitempty" validate:"omitempty,datetime=2006-01-02"`
	ValidUntil                  *string   `json:"validUntil,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// UpdateFinancialInstitutionTermRequest is used to update an existing financial institution term
type UpdateFinancialInstitutionTermRequest struct {
	Name             *string  `json:"name,omitempty" validate:"omitempty,min=2"`
	Description      *string  `json:"description,omitempty"`
	InterestRate     *float64 `json:"interestRate,omitempty" validate:"omitempty,gte=0,lte=1"`
	ProcessingFee    *float64 `json:"processingFee,omitempty" validate:"omitempty,gte=0"`
	TermDays         *int     `json:"termDays,omitempty" validate:"omitempty,gte=0"`
	MinInvoiceAmount *float64 `json:"minInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	MaxInvoiceAmount *float64 `json:"maxInvoiceAmount,omitempty" validate:"omitempty,gte=0"`
	IsActive         *bool    `json:"isActive,omitempty"`
	ValidFrom        *string  `json:"validFrom,omitempty" validate:"omitempty,datetime=2006-01-02"`
	ValidUntil       *string  `json:"validUntil,omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// FinancialInstitutionTermResponse is used to return financial institution term data
type FinancialInstitutionTermResponse struct {
	ID                          uint      `json:"id"`
	FinancialInstitutionID      uint      `json:"financialInstitutionId"`
	FinancialInstitutionProductID *uint     `json:"financialInstitutionProductId,omitempty"`
	Name                        string    `json:"name"`
	Description                 string    `json:"description,omitempty"`
	InterestRate                float64   `json:"interestRate"`
	ProcessingFee               float64   `json:"processingFee,omitempty"`
	TermDays                    int       `json:"termDays"`
	MinInvoiceAmount            float64   `json:"minInvoiceAmount,omitempty"`
	MaxInvoiceAmount            float64   `json:"maxInvoiceAmount,omitempty"`
	IsActive                    bool      `json:"isActive"`
	ValidFrom                   *time.Time `json:"validFrom,omitempty"`
	ValidUntil                  *time.Time `json:"validUntil,omitempty"`
	CreatedAt                   time.Time  `json:"createdAt"`
	UpdatedAt                   time.Time  `json:"updatedAt"`
}

// FinancialInstitutionTermListResponse is used to return a list of financial institution terms
type FinancialInstitutionTermListResponse struct {
	Terms    []FinancialInstitutionTermResponse `json:"terms"`
	Total    int64                              `json:"total"`
	Page     int                                `json:"page"`
	PageSize int                                `json:"pageSize"`
}
