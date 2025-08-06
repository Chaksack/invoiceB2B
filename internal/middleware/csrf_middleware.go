package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// CSRFMiddleware provides protection against Cross-Site Request Forgery attacks
type CSRFMiddleware struct {
	cfg *config.Config
	// tokenStore stores session-specific CSRF tokens when per-session tokens are enabled
	// In a production environment, this would be replaced with Redis
	tokenStore map[string]map[string]time.Time // map[sessionID]map[tokenID]expiryTime
}

// NewCSRFMiddleware creates a new CSRF middleware instance
func NewCSRFMiddleware(cfg *config.Config) *CSRFMiddleware {
	middleware := &CSRFMiddleware{
		cfg:        cfg,
		tokenStore: make(map[string]map[string]time.Time),
	}
	
	// Start a goroutine to clean up expired tokens
	go middleware.cleanupExpiredTokens()
	
	return middleware
}

// cleanupExpiredTokens periodically removes expired CSRF tokens
func (cm *CSRFMiddleware) cleanupExpiredTokens() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		
		// Clean up expired tokens
		for sessionID, tokens := range cm.tokenStore {
			for tokenID, expiry := range tokens {
				if now.After(expiry) {
					delete(tokens, tokenID)
				}
			}
			
			// Remove empty session entries
			if len(tokens) == 0 {
				delete(cm.tokenStore, sessionID)
			}
		}
	}
}

// GenerateToken creates a new CSRF token
func (cm *CSRFMiddleware) GenerateToken(c *fiber.Ctx) error {
	// Get user token from context if available
	var userID string
	var sessionID string
	
	if token, ok := c.Locals("user").(*jwt.Token); ok {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if id, ok := claims["user_id"].(string); ok {
				userID = id
				// Use user ID as session ID for authenticated users
				sessionID = userID
			}
		}
	}
	
	// If no user ID (not authenticated), use IP address as session ID
	if sessionID == "" {
		sessionID = c.IP()
	}
	
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return utils.HandleError(c, fiber.StatusInternalServerError, "Failed to generate CSRF token", err)
	}
	
	// Create base64 encoded token
	csrfToken := base64.StdEncoding.EncodeToString(tokenBytes)
	
	// Generate a unique token ID
	tokenID := uuid.New().String()
	
	// Store token in token store if per-session tokens are enabled
	if cm.cfg.CSRFPerSessionTokens {
		// Initialize session entry if it doesn't exist
		if _, exists := cm.tokenStore[sessionID]; !exists {
			cm.tokenStore[sessionID] = make(map[string]time.Time)
		}
		
		// Store token with expiration time
		cm.tokenStore[sessionID][tokenID] = time.Now().Add(cm.cfg.CSRFTokenExpiration)
	}
	
	// Create HMAC signature using the token, user ID, and token ID
	signature := cm.createSignature(csrfToken, userID, tokenID)
	
	// Set CSRF token in cookie with configured security settings
	cookie := fiber.Cookie{
		Name:     cm.cfg.CSRFCookieName,
		Value:    fmt.Sprintf("%s.%s.%s", csrfToken, tokenID, signature),
		Path:     "/",
		Expires:  time.Now().Add(cm.cfg.CSRFTokenExpiration),
		HTTPOnly: cm.cfg.CSRFCookieHTTPOnly,
		SameSite: cm.cfg.CSRFCookieSameSite,
		Secure:   cm.cfg.CSRFCookieSecure,
	}
	c.Cookie(&cookie)
	
	// Also set the token in response header for JavaScript to read
	c.Set(cm.cfg.CSRFHeaderName, csrfToken)
	
	return c.Next()
}

// VerifyToken validates the CSRF token in the request
func (cm *CSRFMiddleware) VerifyToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip for GET, HEAD, OPTIONS requests as they should be safe
		method := c.Method()
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return c.Next()
		}
		
		// Check if the path is exempt from CSRF protection
		path := c.Path()
		for _, exemptPath := range cm.cfg.CSRFExemptPaths {
			if strings.HasPrefix(path, exemptPath) {
				return c.Next()
			}
		}
		
		// Get user token from context if available
		var userID string
		var sessionID string
		
		if token, ok := c.Locals("user").(*jwt.Token); ok {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if id, ok := claims["user_id"].(string); ok {
					userID = id
					// Use user ID as session ID for authenticated users
					sessionID = userID
				}
			}
		}
		
		// If no user ID (not authenticated), use IP address as session ID
		if sessionID == "" {
			sessionID = c.IP()
		}
		
		// Get token from header
		headerToken := c.Get(cm.cfg.CSRFHeaderName)
		
		// Get token from cookie
		cookieToken := c.Cookies(cm.cfg.CSRFCookieName)
		if cookieToken == "" {
			return utils.HandleError(c, fiber.StatusForbidden, "CSRF token cookie missing", nil)
		}
		
		// Split cookie token into token, token ID, and signature
		parts := strings.Split(cookieToken, ".")
		if len(parts) != 3 {
			return utils.HandleError(c, fiber.StatusForbidden, "Invalid CSRF token format", nil)
		}
		
		cookieTokenValue := parts[0]
		tokenID := parts[1]
		cookieSignature := parts[2]
		
		// Verify that header token matches cookie token
		if headerToken != cookieTokenValue {
			return utils.HandleError(c, fiber.StatusForbidden, "CSRF token mismatch", nil)
		}
		
		// If per-session tokens are enabled, verify that the token exists in the token store
		if cm.cfg.CSRFPerSessionTokens {
			if sessionTokens, exists := cm.tokenStore[sessionID]; !exists {
				return utils.HandleError(c, fiber.StatusForbidden, "Invalid CSRF token session", nil)
			} else if expiry, exists := sessionTokens[tokenID]; !exists {
				return utils.HandleError(c, fiber.StatusForbidden, "Invalid CSRF token ID", nil)
			} else if time.Now().After(expiry) {
				// Remove expired token
				delete(sessionTokens, tokenID)
				return utils.HandleError(c, fiber.StatusForbidden, "Expired CSRF token", nil)
			}
		}
		
		// Verify signature
		expectedSignature := cm.createSignature(cookieTokenValue, userID, tokenID)
		if !hmac.Equal([]byte(cookieSignature), []byte(expectedSignature)) {
			return utils.HandleError(c, fiber.StatusForbidden, "Invalid CSRF token signature", nil)
		}
		
		return c.Next()
	}
}

// createSignature creates an HMAC signature for the CSRF token
func (cm *CSRFMiddleware) createSignature(token, userID, tokenID string) string {
	// Use a separate secret key for CSRF tokens
	h := hmac.New(sha256.New, []byte(cm.cfg.CSRFSecret))
	h.Write([]byte(token))
	h.Write([]byte(userID)) // Bind token to user if authenticated
	h.Write([]byte(tokenID)) // Include token ID in signature
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GetTokenStore returns the token store for testing purposes
func (cm *CSRFMiddleware) GetTokenStore() map[string]map[string]time.Time {
	return cm.tokenStore
}