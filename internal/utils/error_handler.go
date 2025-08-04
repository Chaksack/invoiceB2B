package utils

import (
	"errors"
	"fmt"
	"invoiceB2B/internal/services"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrorCode represents a unique error code for API errors
type ErrorCode string

// Define standard error codes
const (
	ErrCodeUnknown             ErrorCode = "UNKNOWN_ERROR"
	ErrCodeInvalidInput        ErrorCode = "INVALID_INPUT"
	ErrCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrCodeConflict            ErrorCode = "CONFLICT"
	ErrCodeInternalServer      ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError       ErrorCode = "DATABASE_ERROR"
	ErrCodeValidationFailed    ErrorCode = "VALIDATION_FAILED"
	ErrCodeInvalidCredentials  ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeUserNotFound        ErrorCode = "USER_NOT_FOUND"
	ErrCodeEmailExists         ErrorCode = "EMAIL_EXISTS"
	ErrCodeOTPInvalid          ErrorCode = "OTP_INVALID"
	ErrCode2FANotEnabled       ErrorCode = "2FA_NOT_ENABLED"
	ErrCodeAccountNotActive    ErrorCode = "ACCOUNT_NOT_ACTIVE"
	ErrCodeKYCNotApproved      ErrorCode = "KYC_NOT_APPROVED"
	ErrCodeRefreshTokenInvalid ErrorCode = "REFRESH_TOKEN_INVALID"
	ErrCodeTokenBlacklisted    ErrorCode = "TOKEN_BLACKLISTED"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status    string      `json:"status"`
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// getErrorCode maps errors to error codes
func getErrorCode(err error) (ErrorCode, int) {
	// Default values
	code := ErrCodeUnknown
	statusCode := fiber.StatusInternalServerError

	// Check for specific error types
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		code = ErrCodeNotFound
		statusCode = fiber.StatusNotFound
	case errors.Is(err, gorm.ErrInvalidTransaction) ||
		errors.Is(err, gorm.ErrNotImplemented) ||
		errors.Is(err, gorm.ErrMissingWhereClause) ||
		errors.Is(err, gorm.ErrUnsupportedRelation):
		code = ErrCodeDatabaseError
		statusCode = fiber.StatusInternalServerError
	case errors.Is(err, services.ErrInvalidCredentials):
		code = ErrCodeInvalidCredentials
		statusCode = fiber.StatusUnauthorized
	case errors.Is(err, services.ErrUserNotFound):
		code = ErrCodeUserNotFound
		statusCode = fiber.StatusNotFound
	case errors.Is(err, services.ErrEmailExists):
		code = ErrCodeEmailExists
		statusCode = fiber.StatusConflict
	case errors.Is(err, services.ErrOTPInvalidOrExpired):
		code = ErrCodeOTPInvalid
		statusCode = fiber.StatusUnauthorized
	case errors.Is(err, services.Err2FANotEnabled):
		code = ErrCode2FANotEnabled
		statusCode = fiber.StatusBadRequest
	case errors.Is(err, services.ErrAccountNotActive):
		code = ErrCodeAccountNotActive
		statusCode = fiber.StatusForbidden
	case errors.Is(err, services.ErrKYCNotApproved):
		code = ErrCodeKYCNotApproved
		statusCode = fiber.StatusForbidden
	case errors.Is(err, services.ErrRefreshTokenInvalid):
		code = ErrCodeRefreshTokenInvalid
		statusCode = fiber.StatusUnauthorized
	case errors.Is(err, services.ErrTokenBlacklisted):
		code = ErrCodeTokenBlacklisted
		statusCode = fiber.StatusUnauthorized
	}

	// Check for Fiber's own error type
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		statusCode = fiberError.Code
		// Map HTTP status codes to error codes
		switch statusCode {
		case fiber.StatusBadRequest:
			code = ErrCodeInvalidInput
		case fiber.StatusUnauthorized:
			code = ErrCodeUnauthorized
		case fiber.StatusForbidden:
			code = ErrCodeForbidden
		case fiber.StatusNotFound:
			code = ErrCodeNotFound
		case fiber.StatusConflict:
			code = ErrCodeConflict
		}
	}

	return code, statusCode
}

// GlobalErrorHandler provides a centralized way to handle errors.
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	// Generate a request ID for tracking
	requestID := uuid.New().String()

	// Get error code and status code
	errorCode, statusCode := getErrorCode(err)

	// Default message
	message := "An unexpected error occurred. Please try again later."

	// Check for specific error messages
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		message = fiberError.Message
	} else if err != nil {
		// Use the error's message for known error types
		switch errorCode {
		case ErrCodeNotFound, ErrCodeInvalidCredentials, ErrCodeUserNotFound, 
			 ErrCodeEmailExists, ErrCodeOTPInvalid, ErrCode2FANotEnabled, 
			 ErrCodeAccountNotActive, ErrCodeKYCNotApproved, ErrCodeRefreshTokenInvalid, 
			 ErrCodeTokenBlacklisted:
			message = err.Error()
		case ErrCodeDatabaseError:
			message = "A database operation failed due to an internal issue."
		}
	}

	// Additional details for development environment
	var details interface{}
	if os.Getenv("APP_ENV") == "development" {
		details = map[string]interface{}{
			"error": err.Error(),
			"stack": getStackTrace(2), // Skip this function and the caller
		}
	}

	// Create the error response
	response := ErrorResponse{
		Status:    "error",
		Code:      errorCode,
		Message:   message,
		Details:   details,
		RequestID: requestID,
		Timestamp: time.Now(),
	}

	// Log the error with context
	logErrorWithContext(c, err, statusCode, requestID, message)

	// Return the error response
	return c.Status(statusCode).JSON(response)
}

// HandleError is a utility for handlers to return structured errors.
func HandleError(c *fiber.Ctx, statusCode int, msg string, originalError error) error {
	// If originalError is a known service error, return it directly
	// This allows the GlobalErrorHandler to properly identify the error type
	if originalError != nil {
		switch {
		case errors.Is(originalError, services.ErrInvalidCredentials),
			 errors.Is(originalError, services.ErrUserNotFound),
			 errors.Is(originalError, services.ErrEmailExists),
			 errors.Is(originalError, services.ErrOTPInvalidOrExpired),
			 errors.Is(originalError, services.Err2FANotEnabled),
			 errors.Is(originalError, services.ErrAccountNotActive),
			 errors.Is(originalError, services.ErrKYCNotApproved),
			 errors.Is(originalError, services.ErrRefreshTokenInvalid),
			 errors.Is(originalError, services.ErrTokenBlacklisted):
			return originalError
		}
	}

	// Log the error
	if originalError != nil {
		log.Printf("Handler Error: %s - Original: %v - Path: %s", msg, originalError, c.Path())
	} else {
		log.Printf("Handler Info/Error: %s - Path: %s", msg, c.Path())
	}

	// Create a new fiber.Error with the status code and message
	return fiber.NewError(statusCode, msg)
}

// HandleValidationError is a utility for handling validation errors
func HandleValidationError(c *fiber.Ctx, err error) error {
	// Extract validation errors and format them
	validationErrors := extractValidationErrors(err)

	// Create a detailed error response
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Status:    "error",
		Code:      ErrCodeValidationFailed,
		Message:   "Validation failed. Please check your input.",
		Details:   validationErrors,
		RequestID: uuid.New().String(),
		Timestamp: time.Now(),
	})
}

// Helper functions

// getStackTrace returns a simplified stack trace
func getStackTrace(skip int) []string {
	var stack []string
	for i := skip; i < skip+5; i++ { // Capture 5 frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}
	return stack
}

// logErrorWithContext logs an error with additional context
func logErrorWithContext(c *fiber.Ctx, err error, statusCode int, requestID, message string) {
	logEntry := fmt.Sprintf(
		"[%s] Error: %s - Status: %d - RequestID: %s - Path: %s - Method: %s - IP: %s",
		time.Now().Format(time.RFC3339),
		message,
		statusCode,
		requestID,
		c.Path(),
		c.Method(),
		c.IP(),
	)

	if err != nil {
		logEntry += fmt.Sprintf(" - Original Error: %v", err)
	}

	if statusCode >= 500 {
		log.Printf("SERVER ERROR: %s", logEntry)
	} else {
		log.Printf("CLIENT ERROR: %s", logEntry)
	}
}

// extractValidationErrors extracts validation errors from the validator
func extractValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	// Check if it's a validator.ValidationErrors type
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			// Get a user-friendly error message based on the validation tag
			fieldName := e.Field()
			// Convert first letter to lowercase for JSON field name convention
			if len(fieldName) > 0 {
				fieldName = strings.ToLower(fieldName[:1]) + fieldName[1:]
			}

			// Add the formatted error message
			errors[fieldName] = formatErrorMessage(e)
		}
		return errors
	}

	// If it's not a validation error or we can't extract details,
	// just add the generic error message
	errors["_error"] = err.Error()

	return errors
}
