package services

import (
	"context"
	"errors"
	"fmt"
	"invoiceB2B/internal/dtos"
	"log"

	"github.com/google/uuid" // Keep for simulated external transaction IDs
)

type PaymentService interface {
	InitiateDisbursement(ctx context.Context, req dtos.DisbursementRequest) (*dtos.DisbursementResponse, error)
	ProcessRepayment(ctx context.Context, req dtos.RepaymentRequest) (*dtos.RepaymentResponse, error)
	CheckPaymentStatus(ctx context.Context, transactionID string) (*dtos.PaymentStatusResponse, error)
}

type paymentService struct {
}

func NewPaymentService() PaymentService {
	return &paymentService{}
}

func (s *paymentService) InitiateDisbursement(ctx context.Context, req dtos.DisbursementRequest) (*dtos.DisbursementResponse, error) {
	log.Printf("Attempting to initiate disbursement for Invoice ID %d, Amount: %.2f %s to Account: %s",
		req.InvoiceID, req.Amount, req.Currency, req.BankAccountNumber)

	if req.Amount <= 0 {
		return nil, errors.New("disbursement amount must be positive")
	}

	simulatedTransactionID := "DISB_" + uuid.New().String()
	log.Printf("Simulated disbursement successful. Transaction ID: %s", simulatedTransactionID)

	return &dtos.DisbursementResponse{
		TransactionID: simulatedTransactionID,
		Status:        "SUCCESS",
		Message:       "Disbursement initiated successfully (simulated).",
	}, nil
}

func (s *paymentService) ProcessRepayment(ctx context.Context, req dtos.RepaymentRequest) (*dtos.RepaymentResponse, error) {
	log.Printf("Attempting to process repayment for Invoice ID %d, Amount: %.2f %s from User ID: %d",
		req.InvoiceID, req.Amount, req.Currency, req.UserID)

	if req.Amount <= 0 {
		return nil, errors.New("repayment amount must be positive")
	}
	if req.PaymentMethodToken == "INVALID_TOKEN_FOR_SIMULATION" {
		return &dtos.RepaymentResponse{
			TransactionID: "",
			Status:        "FAILED",
			Message:       "Repayment failed due to invalid payment method (simulated).",
		}, nil
	}

	simulatedTransactionID := "REPAY_" + uuid.New().String()
	log.Printf("Simulated repayment successful. Transaction ID: %s", simulatedTransactionID)

	return &dtos.RepaymentResponse{
		TransactionID: simulatedTransactionID,
		Status:        "SUCCESS",
		Message:       "Repayment processed successfully (simulated).",
	}, nil
}

func (s *paymentService) CheckPaymentStatus(ctx context.Context, transactionID string) (*dtos.PaymentStatusResponse, error) {
	log.Printf("Checking payment status for Transaction ID: %s", transactionID)

	if len(transactionID) < 5 {
		return nil, fmt.Errorf("invalid transaction ID format for simulation")
	}

	status := "PENDING"
	message := "Payment status is pending (simulated)."

	if transactionID[:5] == "DISB_" {
		if len(transactionID)%2 == 0 {
			status = "SUCCESS"
			message = "Disbursement completed successfully (simulated)."
		}
	} else if transactionID[:6] == "REPAY_" {
		status = "SUCCESS"
		message = "Repayment confirmed (simulated)."
	}

	return &dtos.PaymentStatusResponse{
		TransactionID: transactionID,
		Status:        status,
		Message:       message,
		Amount:        100.00,
		Currency:      "USD",
	}, nil
}
