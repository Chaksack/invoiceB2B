package middleware

import (
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	jwtService services.JWTService
}

func NewAuthMiddleware(jwtService services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (am *AuthMiddleware) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Missing Authorization Header", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid Authorization Header format", nil)
		}

		tokenStr := parts[1]
		claims, err := am.jwtService.ValidateToken(tokenStr, false) // false for access token
		if err != nil {
			return utils.HandleError(c, fiber.StatusUnauthorized, "Invalid or expired token", err)
		}

		// Check if token is blacklisted (for logout)
		// This requires otpService or a dedicated blacklist service.
		// For simplicity, let's assume otpService has IsTokenBlacklisted.
		// This dependency might need to be passed to AuthMiddleware or handled differently.
		// isBlacklisted, _ := am.otpService.IsTokenBlacklisted(c.Context(), tokenStr)
		// if isBlacklisted {
		// 	return utils.HandleError(c, fiber.StatusUnauthorized, "Token has been invalidated", nil)
		// }

		// Store claims in context for handlers to use
		// Convert claims (jwt.MapClaims) to a more usable struct if needed, or pass as is.
		// For simplicity, passing the raw token object which contains claims.
		// Handlers will need to cast c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)

		// Create a new token object to store in locals, as the parsed one might not be ideal.
		// It's better to store the parsed claims directly or a custom user identity struct.
		// For now, storing the validated token object.
		// A common practice is to parse into your custom Claims struct.

		// Create a jwt.Token object to store in locals. This is a bit of a workaround
		// as Fiber doesn't have a direct way to pass arbitrary structs easily without type assertion.
		// The claims are already validated.
		tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Recreate a token object with validated claims
		tokenWithClaims.Raw = tokenStr                                       // Store the raw token string as well if needed later (e.g. for blacklist)

		c.Locals("user", tokenWithClaims) // Store the validated token object (which includes claims)
		// c.Locals("user_id", claims["user_id"]) // Example of storing specific claim

		return c.Next()
	}
}
