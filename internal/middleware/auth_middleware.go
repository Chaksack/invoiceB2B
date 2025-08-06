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

		// Token validation now includes:
		// 1. Signature verification with RSA public key
		// 2. Expiration time validation
		// 3. Not before time validation
		// 4. Audience validation
		// 5. Token revocation check
		// 6. Token type validation (access vs refresh)
		
		// Extract user ID and role from claims for convenience
		userID, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)
		
		// Store these in context for easy access in handlers
		c.Locals("user_id", userID)
		c.Locals("role", role)
		
		// Create a jwt.Token object to store in locals
		// Use RS256 since we're now using asymmetric signing
		tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenWithClaims.Raw = tokenStr
		
		// Store the full token object for handlers that need access to all claims
		c.Locals("user", tokenWithClaims)

		return c.Next()
	}
}

// Logout handler helper to revoke the current token
func (am *AuthMiddleware) RevokeCurrentToken(c *fiber.Ctx) error {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Token not found in context", nil)
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Invalid token claims", nil)
	}
	
	// Get the token ID (jti) from claims
	tokenID, ok := claims["jti"].(string)
	if !ok {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Token ID not found in claims", nil)
	}
	
	// Revoke the token
	if err := am.jwtService.RevokeToken(tokenID); err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to revoke token", err)
	}
	
	return nil
}
