package middleware

import (
	"context"
	"fmt"
	"invoiceB2B/internal/interfaces"
	"invoiceB2B/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AuthorizationMiddleware provides attribute-based access control
type AuthorizationMiddleware struct {
	authzService interfaces.AuthorizationService
}

// NewAuthorizationMiddleware creates a new authorization middleware
func NewAuthorizationMiddleware(authzService interfaces.AuthorizationService) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		authzService: authzService,
	}
}

// RequirePermission checks if the user has permission to perform the action on the resource
func (am *AuthorizationMiddleware) RequirePermission(action string, resourceType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract subject information from JWT token
		token, ok := c.Locals("user").(*jwt.Token)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Authentication required", nil)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
		}

		// Create subject attributes
		subject := map[string]interface{}{
			"id":   claims["user_id"],
			"type": claims["role"],
		}

		// Add all claims to subject attributes
		for key, value := range claims {
			subject[key] = value
		}

		// Create resource attributes
		resource := map[string]interface{}{
			"type": resourceType,
		}

		// Add resource ID if available in path params
		if id := c.Params("id"); id != "" {
			resource["id"] = id
		}

		// Add additional resource attributes based on the resource type
		am.addResourceAttributes(c, resource, resourceType)

		// Store the Fiber context in the context for access in the authorization service
		ctx := context.WithValue(c.Context(), "fiber_ctx", c)

		// Check if the action is allowed
		allowed, err := am.authzService.Authorize(ctx, subject, action, resource)
		if err != nil {
			return utils.HandleError(c, fiber.StatusInternalServerError, "Authorization check failed", err)
		}

		if !allowed {
			return utils.HandleError(c, fiber.StatusForbidden, "Permission denied", nil)
		}

		// Store authorization context for handlers
		c.Locals("authz_subject", subject)
		c.Locals("authz_action", action)
		c.Locals("authz_resource", resource)

		return c.Next()
	}
}

// CheckOwnership is a specialized middleware that ensures a user can only access their own resources
func (am *AuthorizationMiddleware) CheckOwnership(resourceType string, idParam string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user ID from JWT token
		token, ok := c.Locals("user").(*jwt.Token)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Authentication required", nil)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid user ID in token", nil)
		}

		// Get resource ID from path params
		resourceID := c.Params(idParam)
		if resourceID == "" {
			return utils.HandleError(c, fiber.StatusBadRequest, "Resource ID not provided", nil)
		}

		// Create subject and resource for authorization check
		subject := map[string]interface{}{
			"id":   userID,
			"type": claims["role"],
		}

		resource := map[string]interface{}{
			"id":    resourceID,
			"type":  resourceType,
			"owner": userID, // Assume the resource has an owner field
		}

		// Add additional resource attributes
		am.addResourceAttributes(c, resource, resourceType)

		// Store the Fiber context in the context for access in the authorization service
		ctx := context.WithValue(c.Context(), "fiber_ctx", c)

		// Check if the user is allowed to access this resource
		allowed, err := am.authzService.Authorize(ctx, subject, "access", resource)
		if err != nil {
			return utils.HandleError(c, fiber.StatusInternalServerError, "Authorization check failed", err)
		}

		if !allowed {
			return utils.HandleError(c, fiber.StatusForbidden, "You do not have permission to access this resource", nil)
		}

		return c.Next()
	}
}

// addResourceAttributes adds resource-specific attributes based on the resource type
func (am *AuthorizationMiddleware) addResourceAttributes(c *fiber.Ctx, resource map[string]interface{}, resourceType string) {
	switch resourceType {
	case "invoice":
		// Add invoice-specific attributes
		if id := c.Params("id"); id != "" {
			// In a real implementation, you would fetch the invoice from the database
			// and add relevant attributes like status, amount, owner, etc.
			resource["id"] = id
		}
	case "user":
		// Add user-specific attributes
		if id := c.Params("id"); id != "" {
			// In a real implementation, you would fetch the user from the database
			// and add relevant attributes like role, status, etc.
			resource["id"] = id
		}
	case "kyc":
		// Add KYC-specific attributes
		if id := c.Params("id"); id != "" {
			// In a real implementation, you would fetch the KYC details from the database
			// and add relevant attributes like status, verification level, etc.
			resource["id"] = id
		}
	}

	// Add query parameters as resource attributes with "query:" prefix
	c.QueryParser(nil) // Ensure query params are parsed
	for key, value := range c.Queries() {
		resource["query:"+key] = value
	}

	// For POST/PUT requests, add body parameters as resource attributes with "body:" prefix
	if c.Method() == "POST" || c.Method() == "PUT" {
		var body map[string]interface{}
		if err := c.BodyParser(&body); err == nil {
			for key, value := range body {
				// Skip sensitive fields
				if !isSensitiveField(key) {
					resource["body:"+key] = value
				}
			}
			// Reset body for handlers to read it again
			c.Request().ResetBody()
		}
	}
}

// isSensitiveField checks if a field name indicates sensitive data that shouldn't be included in authorization logs
func isSensitiveField(field string) bool {
	field = strings.ToLower(field)
	sensitiveFields := []string{
		"password", "secret", "token", "key", "auth", "credential", "credit_card",
		"card", "ssn", "social", "tax", "bank", "account", "routing", "pin", "cvv",
	}

	for _, sensitive := range sensitiveFields {
		if strings.Contains(field, sensitive) {
			return true
		}
	}

	return false
}

// ResourceBasedAuth is a factory function that creates middleware for specific resource types and actions
func (am *AuthorizationMiddleware) ResourceBasedAuth(resourceType string, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract subject information from JWT token
		token, ok := c.Locals("user").(*jwt.Token)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Authentication required", nil)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
		}

		// Create subject attributes
		subject := map[string]interface{}{
			"id":   claims["user_id"],
			"type": claims["role"],
		}

		// Add all claims to subject attributes
		for key, value := range claims {
			subject[key] = value
		}

		// Create resource attributes
		resource := map[string]interface{}{
			"type": resourceType,
		}

		// Add resource ID if available in path params
		if id := c.Params("id"); id != "" {
			resource["id"] = id
		}

		// Add additional resource attributes
		am.addResourceAttributes(c, resource, resourceType)

		// Store the Fiber context in the context for access in the authorization service
		ctx := context.WithValue(c.Context(), "fiber_ctx", c)

		// Check if the action is allowed
		allowed, err := am.authzService.Authorize(ctx, subject, action, resource)
		if err != nil {
			return utils.HandleError(c, fiber.StatusInternalServerError, fmt.Sprintf("Authorization check failed: %v", err), nil)
		}

		if !allowed {
			return utils.HandleError(c, fiber.StatusForbidden, fmt.Sprintf("You do not have permission to %s this %s", action, resourceType), nil)
		}

		return c.Next()
	}
}