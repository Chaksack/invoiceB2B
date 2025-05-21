package dtos

import "github.com/google/uuid"

// Payment DTOs

// DisbursementRequest is used when the admin initiates a fund transfer to the manufacturer.
type DisbursementRequest struct {
	InvoiceID         uuid.UUID `json:"invoiceId" validate:"required"`
	Amount            float64   `json:"amount" validate:"required,gt=0"`
	Currency          string    `json:"currency" validate:"required,len=3"`
	BankAccountNumber string    `json:"bankAccountNumber" validate:"required"` // Manufacturer's bank account
	BankName          string    `json:"bankName" validate:"required"`
	RecipientName     string    `json:"recipientName" validate:"required"`
	// Add other necessary fields for your payment gateway
}

type DisbursementResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"` // e.g., "SUCCESS", "PENDING", "FAILED"
	Message       string `json:"message"`
	// GatewaySpecificResponse map[string]interface{} `json:"gatewaySpecificResponse,omitempty"`
}

// RepaymentRequest is used when the user makes a repayment.
type RepaymentRequest struct {
	InvoiceID          uuid.UUID `json:"invoiceId" validate:"required"`
	UserID             uuid.UUID `json:"userId" validate:"required"` // To identify the payer
	Amount             float64   `json:"amount" validate:"required,gt=0"`
	Currency           string    `json:"currency" validate:"required,len=3"`
	PaymentMethodToken string    `json:"paymentMethodToken" validate:"required"` // Token from client-side payment integration (e.g., Stripe token)
	// Add other necessary fields
}

type RepaymentResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"` // e.g., "SUCCESS", "PENDING", "FAILED"
	Message       string `json:"message"`
}

// PaymentStatusResponse is used to check the status of any transaction.
type PaymentStatusResponse struct {
	TransactionID string  `json:"transactionId"`
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	Amount        float64 `json:"amount,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	// Add other relevant fields like timestamps, etc.
}
