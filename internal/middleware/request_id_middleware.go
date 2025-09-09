package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates unique request IDs for API traceability
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware creates a new request ID middleware
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

// GenerateRequestID generates a unique request ID for each API request
func (m *RequestIDMiddleware) GenerateRequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID already exists (from load balancer or client)
		requestID := c.Get("X-Request-ID")
		
		// Generate new request ID if none exists
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Set request ID in context for handlers to use
		c.Set("X-Request-ID", requestID)
		
		// Continue to next middleware
		return c.Next()
	}
}