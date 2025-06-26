package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles errors in the application
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Log the error
	log.Printf("Error: %v", err)

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"success":   false,
		"message":   message,
		"timestamp": fiber.Now(),
	})
} 