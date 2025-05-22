package dtos

type DisbursementRequest struct {
	InvoiceID         uint    `json:"invoiceId" validate:"required"`
	Amount            float64 `json:"amount" validate:"required,gt=0"`
	Currency          string  `json:"currency" validate:"required,len=3"`
	BankAccountNumber string  `json:"bankAccountNumber" validate:"required"`
	BankName          string  `json:"bankName" validate:"required"`
	RecipientName     string  `json:"recipientName" validate:"required"`
}

type DisbursementResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type RepaymentRequest struct {
	InvoiceID          uint    `json:"invoiceId" validate:"required"`
	UserID             uint    `json:"userId" validate:"required"`
	Amount             float64 `json:"amount" validate:"required,gt=0"`
	Currency           string  `json:"currency" validate:"required,len=3"`
	PaymentMethodToken string  `json:"paymentMethodToken" validate:"required"`
}

type RepaymentResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type PaymentStatusResponse struct {
	TransactionID string  `json:"transactionId"`
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	Amount        float64 `json:"amount,omitempty"`
	Currency      string  `json:"currency,omitempty"`
}
