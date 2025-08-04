package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ImprovedExtractValidationErrors extracts validation errors from the validator
// and formats them in a more user-friendly way
func ImprovedExtractValidationErrors(err error) map[string]string {
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

			// Add the formatted error message using our custom formatter
			errors[fieldName] = improvedFormatErrorMessage(e)
		}
		return errors
	}

	// If it's not a validation error or we can't extract details,
	// just add the generic error message
	errors["_error"] = err.Error()

	return errors
}

// improvedFormatErrorMessage formats validation errors into user-friendly messages
func improvedFormatErrorMessage(fe validator.FieldError) string {
	// Customize error messages based on tag
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", strings.ToLower(fe.Field()))
	case "email":
		return fmt.Sprintf("The %s field must be a valid email address.", strings.ToLower(fe.Field()))
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters long.", strings.ToLower(fe.Field()), fe.Param())
	case "max":
		return fmt.Sprintf("The %s field may not be greater than %s characters.", strings.ToLower(fe.Field()), fe.Param())
	case "len":
		return fmt.Sprintf("The %s field must be %s characters long.", strings.ToLower(fe.Field()), fe.Param())
	case "numeric":
		return fmt.Sprintf("The %s field must be a number.", strings.ToLower(fe.Field()))
	case "gt":
		return fmt.Sprintf("The %s field must be greater than %s.", strings.ToLower(fe.Field()), fe.Param())
	case "lt":
		return fmt.Sprintf("The %s field must be less than %s.", strings.ToLower(fe.Field()), fe.Param())
	case "gte":
		return fmt.Sprintf("The %s field must be greater than or equal to %s.", strings.ToLower(fe.Field()), fe.Param())
	case "lte":
		return fmt.Sprintf("The %s field must be less than or equal to %s.", strings.ToLower(fe.Field()), fe.Param())
	case "oneof":
		return fmt.Sprintf("The %s field must be one of: %s.", strings.ToLower(fe.Field()), fe.Param())
	case "url":
		return fmt.Sprintf("The %s field must be a valid URL.", strings.ToLower(fe.Field()))
	case "uuid":
		return fmt.Sprintf("The %s field must be a valid UUID.", strings.ToLower(fe.Field()))
	case "json":
		return fmt.Sprintf("The %s field must be a valid JSON string.", strings.ToLower(fe.Field()))
	default:
		return fmt.Sprintf("Validation failed on field '%s' with tag '%s'.", fe.Field(), fe.Tag())
	}
}

// UpdateHandleValidationError is an improved version of HandleValidationError
// that uses ImprovedExtractValidationErrors
func UpdateHandleValidationError(c *fiber.Ctx, err error) error {
	// Extract validation errors and format them
	validationErrors := ImprovedExtractValidationErrors(err)

	// Generate a request ID for tracking
	requestID := uuid.New().String()

	// Create a detailed error response
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Status:    "error",
		Code:      ErrCodeValidationFailed,
		Message:   "Validation failed. Please check your input.",
		Details:   validationErrors,
		RequestID: requestID,
		Timestamp: time.Now(),
	})
}
