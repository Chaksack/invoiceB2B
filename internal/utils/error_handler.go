package utils

import (
	"errors"
	"invoiceB2B/internal/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GlobalErrorHandler provides a centralized way to handle errors.
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "An unexpected error occurred. Please try again later."
	var additionalDetails interface{}

	// Check for Fiber's own error type
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		code = fiberError.Code
		message = fiberError.Message
	} else {
		// Handle GORM errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = fiber.StatusNotFound
			message = "The requested resource was not found."
		} else if errors.Is(err, gorm.ErrInvalidTransaction) ||
			errors.Is(err, gorm.ErrNotImplemented) ||
			errors.Is(err, gorm.ErrMissingWhereClause) ||
			errors.Is(err, gorm.ErrUnsupportedRelation) {
			code = fiber.StatusInternalServerError
			message = "A database operation failed due to an internal issue."
			log.Printf("GORM Error: %v - Path: %s", err, c.Path())
		} else {
			switch {
			case errors.Is(err, services.ErrInvalidCredentials):
				code = fiber.StatusUnauthorized
				message = err.Error()
			case errors.Is(err, services.ErrUserNotFound):
				code = fiber.StatusNotFound
				message = err.Error()
			case errors.Is(err, services.ErrEmailExists):
				code = fiber.StatusConflict
				message = err.Error()
			case errors.Is(err, services.ErrOTPInvalidOrExpired):
				code = fiber.StatusUnauthorized
				message = err.Error()
			case errors.Is(err, services.Err2FANotEnabled):
				code = fiber.StatusBadRequest // Or Unauthorized, depending on context
				message = err.Error()
			case errors.Is(err, services.ErrAccountNotActive):
				code = fiber.StatusForbidden
				message = err.Error()
			case errors.Is(err, services.ErrKYCNotApproved):
				code = fiber.StatusForbidden
				message = err.Error()
			case errors.Is(err, services.ErrRefreshTokenInvalid):
				code = fiber.StatusUnauthorized
				message = err.Error()
			case errors.Is(err, services.ErrTokenBlacklisted):
				code = fiber.StatusUnauthorized
				message = err.Error()
			default:
				// Fallback for generic errors not caught above
				log.Printf("Unhandled Error: %v - Path: %s - Method: %s", err, c.Path(), c.Method())
			}
		}
	}

	// For development, you might want to expose more error details
	if os.Getenv("APP_ENV") == "development" {
		additionalDetails = err.Error() // Show original error message in dev
	}

	// Log the error with more context
	if code >= 500 {
		log.Printf("Server Error (%d): %s - Path: %s - Method: %s - Error: %v", code, message, c.Path(), c.Method(), err)
	} else {
		log.Printf("Client Error (%d): %s - Path: %s - Method: %s - Error: %v", code, message, c.Path(), c.Method(), err)
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"details": additionalDetails,
	})
}

// HandleError is a utility for handlers to return structured errors.
// This function now directly calls fiber.NewError to ensure the GlobalErrorHandler is triggered.
func HandleError(c *fiber.Ctx, statusCode int, msg string, originalError error) error {
	// Log the original error for internal tracking if it exists
	fullMessage := msg
	if originalError != nil {
		log.Printf("Handler Error: %s - Original: %v - Path: %s", msg, originalError, c.Path())
		// Optionally append original error to message for dev, or let GlobalErrorHandler handle it
		// if os.Getenv("APP_ENV") == "development" {
		// 	fullMessage = fmt.Sprintf("%s: %v", msg, originalError)
		// }
	} else {
		log.Printf("Handler Info/Error: %s - Path: %s", msg, c.Path())
	}

	// Return a new fiber.Error. The GlobalErrorHandler will catch this.
	// Pass the originalError if you want GlobalErrorHandler to inspect it further.
	// For now, we are creating a new fiber.Error with the statusCode and msg.
	// If originalError is one of our custom service errors, it's better to return originalError directly.
	// This function is more for creating new errors within handlers.

	// If originalError is already a well-defined error our global handler knows,
	// it might be better to return that directly from the service layer.
	// This HandleError is for when a handler itself identifies an issue.
	return fiber.NewError(statusCode, fullMessage)
}
