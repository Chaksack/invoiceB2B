package middleware

import (
	"invoiceB2B/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

type InternalAPIMiddleware struct {
	apiKey string
}

func NewInternalAPIMiddleware(apiKey string) *InternalAPIMiddleware {
	if apiKey == "" || apiKey == "default-internal-key-please-change" {
		log.Printf("WARN: INTERNAL_API_KEY is not set securely. Please configure it properly for production.")
	}
	return &InternalAPIMiddleware{apiKey: apiKey}
}

func (m *InternalAPIMiddleware) RequireAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		providedKey := c.Get("X-Internal-API-Key")

		if m.apiKey == "" || m.apiKey == "default-internal-key-please-change" {
			log.Printf("ERROR: Internal API Key is not configured or is default, blocking request for security.")
			return utils.HandleError(c, fiber.StatusServiceUnavailable, "Internal API endpoint not available due to misconfiguration.", nil)
		}

		if providedKey == "" {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Missing X-Internal-API-Key header.", nil)
		}

		if providedKey != m.apiKey {
			return utils.HandleError(c, fiber.StatusForbidden, "Invalid X-Internal-API-Key.", nil)
		}
		return c.Next()
	}
}
