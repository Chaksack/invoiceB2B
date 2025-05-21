package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{Validator: validator.New()}
}

// FormatValidationError formats validation errors into a readable structure.
func FormatValidationError(err error) []map[string]string {
	var errors []map[string]string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors = append(errors, map[string]string{
				"field":   fieldErr.Field(),
				"tag":     fieldErr.Tag(),
				"value":   fmt.Sprintf("%v", fieldErr.Value()),
				"message": formatErrorMessage(fieldErr),
			})
		}
	}
	return errors
}

func formatErrorMessage(fe validator.FieldError) string {
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
	default:
		return fmt.Sprintf("Validation failed on field '%s' with tag '%s'.", fe.Field(), fe.Tag())
	}
}

// HandleValidationError is a utility function for handlers to return validation errors.
func HandleValidationError(c *fiber.Ctx, errs error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  "error",
		"message": "Validation failed",
		"errors":  FormatValidationError(errs),
	})
}
