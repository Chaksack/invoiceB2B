package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeadersConfig holds the configuration for security headers
type SecurityHeadersConfig struct {
	// ContentSecurityPolicy defines the Content-Security-Policy header value
	ContentSecurityPolicy string

	// XContentTypeOptions defines the X-Content-Type-Options header value
	XContentTypeOptions string

	// XFrameOptions defines the X-Frame-Options header value
	XFrameOptions string

	// StrictTransportSecurity defines the Strict-Transport-Security header value
	StrictTransportSecurity string

	// ReferrerPolicy defines the Referrer-Policy header value
	ReferrerPolicy string

	// PermissionsPolicy defines the Permissions-Policy header value
	PermissionsPolicy string

	// XSSProtection defines the X-XSS-Protection header value
	XSSProtection string

	// CacheControl defines the Cache-Control header value
	CacheControl string
}

// DefaultSecurityHeadersConfig returns the default configuration for security headers
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		// Restrict resources to same origin, allow images, scripts, styles, and fonts from self
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' data:; script-src 'self'; style-src 'self'; font-src 'self'; connect-src 'self'",
		
		// Prevent MIME type sniffing
		XContentTypeOptions: "nosniff",
		
		// Prevent embedding in frames (clickjacking protection)
		XFrameOptions: "DENY",
		
		// Force HTTPS for 1 year, include subdomains
		StrictTransportSecurity: "max-age=31536000; includeSubDomains",
		
		// Control how much referrer information is included with requests
		ReferrerPolicy: "strict-origin-when-cross-origin",
		
		// Restrict browser features
		PermissionsPolicy: "camera=(), microphone=(), geolocation=(), interest-cohort=()",
		
		// Enable XSS filtering in browsers that support it
		XSSProtection: "1; mode=block",
		
		// Control caching of responses
		CacheControl: "no-store, max-age=0",
	}
}

// SecurityHeadersMiddleware adds security headers to all responses
type SecurityHeadersMiddleware struct {
	config SecurityHeadersConfig
}

// NewSecurityHeadersMiddleware creates a new security headers middleware with default configuration
func NewSecurityHeadersMiddleware() *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{
		config: DefaultSecurityHeadersConfig(),
	}
}

// NewSecurityHeadersMiddlewareWithConfig creates a new security headers middleware with custom configuration
func NewSecurityHeadersMiddlewareWithConfig(config SecurityHeadersConfig) *SecurityHeadersMiddleware {
	return &SecurityHeadersMiddleware{
		config: config,
	}
}

// AddSecurityHeaders adds security headers to all responses
func (shm *SecurityHeadersMiddleware) AddSecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add security headers
		if shm.config.ContentSecurityPolicy != "" {
			c.Set("Content-Security-Policy", shm.config.ContentSecurityPolicy)
		}
		
		if shm.config.XContentTypeOptions != "" {
			c.Set("X-Content-Type-Options", shm.config.XContentTypeOptions)
		}
		
		if shm.config.XFrameOptions != "" {
			c.Set("X-Frame-Options", shm.config.XFrameOptions)
		}
		
		if shm.config.StrictTransportSecurity != "" {
			c.Set("Strict-Transport-Security", shm.config.StrictTransportSecurity)
		}
		
		if shm.config.ReferrerPolicy != "" {
			c.Set("Referrer-Policy", shm.config.ReferrerPolicy)
		}
		
		if shm.config.PermissionsPolicy != "" {
			c.Set("Permissions-Policy", shm.config.PermissionsPolicy)
		}
		
		if shm.config.XSSProtection != "" {
			c.Set("X-XSS-Protection", shm.config.XSSProtection)
		}
		
		if shm.config.CacheControl != "" {
			c.Set("Cache-Control", shm.config.CacheControl)
		}
		
		// Continue to the next middleware/handler
		return c.Next()
	}
}

// AddAPISecurityHeaders adds security headers specifically for API responses
func (shm *SecurityHeadersMiddleware) AddAPISecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add API-specific security headers
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Cache-Control", "no-store, max-age=0")
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// For API endpoints, we can use a more restrictive CSP
		c.Set("Content-Security-Policy", "default-src 'none'")
		
		// Continue to the next middleware/handler
		return c.Next()
	}
}

// AddSecureDownloadHeaders adds security headers for file downloads
func (shm *SecurityHeadersMiddleware) AddSecureDownloadHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add download-specific security headers
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Content-Security-Policy", "default-src 'none'")
		c.Set("X-Download-Options", "noopen")
		
		// Set Content-Disposition to attachment to force download
		if c.Get("Content-Disposition") == "" {
			fileName := c.Query("filename", "download")
			c.Set("Content-Disposition", "attachment; filename="+fileName)
		}
		
		// Continue to the next middleware/handler
		return c.Next()
	}
}