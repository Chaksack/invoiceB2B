package middleware

import (
	"context"
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/utils"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// RateLimiter defines the rate limiting configuration for an endpoint
type RateLimiter struct {
	// Maximum number of requests allowed in the time window
	Limit int
	// Time window in seconds
	Window int
	// Whether to use IP-based rate limiting
	IPBased bool
	// Whether to use user-based rate limiting
	UserBased bool
	// Whether to use adaptive rate limiting based on server load
	Adaptive bool
	// Base limit for adaptive rate limiting (will be adjusted based on load)
	BaseLimit int
	// Maximum limit for adaptive rate limiting
	MaxLimit int
	// Minimum limit for adaptive rate limiting
	MinLimit int
}

// RateLimitMiddleware provides rate limiting functionality
type RateLimitMiddleware struct {
	cfg *config.Config
	// Redis client for distributed rate limiting
	redisClient *redis.Client
	// Context for Redis operations
	ctx context.Context
	// Key prefix for Redis storage
	redisKeyPrefix string
	// Fallback in-memory store for rate limiting if Redis is unavailable
	ipLimits     map[string]map[int64]int
	userLimits   map[string]map[int64]int
	mu           sync.RWMutex
	// Rate limiting configuration for different endpoints
	limiters     map[string]RateLimiter
	// Whether to use Redis for rate limiting
	useRedis     bool
	// Adaptive rate limiting metrics
	requestCount     int64         // Total number of requests processed
	requestsPerMinute int64        // Requests in the last minute
	avgResponseTime  int64         // Average response time in microseconds
	lastUpdated      time.Time     // When metrics were last updated
	metricsMu        sync.RWMutex  // Mutex for synchronizing access to metrics
}

// Redis key constants
const (
	// Key prefix for IP-based rate limiting
	ipRateLimitKeyPrefix = "rate:ip:"
	// Key prefix for user-based rate limiting
	userRateLimitKeyPrefix = "rate:user:"
)

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(cfg *config.Config, redisClient *redis.Client) *RateLimitMiddleware {
	// Determine if Redis should be used
	useRedis := redisClient != nil
	
	middleware := &RateLimitMiddleware{
		cfg:              cfg,
		redisClient:      redisClient,
		ctx:              context.Background(),
		redisKeyPrefix:   "rate:",
		ipLimits:         make(map[string]map[int64]int),
		userLimits:       make(map[string]map[int64]int),
		limiters:         make(map[string]RateLimiter),
		useRedis:         useRedis,
		requestCount:     0,
		requestsPerMinute: 0,
		avgResponseTime:  0,
		lastUpdated:      time.Now(),
	}
	
	// Start a goroutine to periodically update request metrics
	go middleware.updateRequestMetrics()

	// Configure rate limiters for sensitive endpoints
	
	// Authentication endpoints
	middleware.limiters["/auth/login"] = RateLimiter{
		Limit:     10,
		Window:    60, // 10 requests per minute
		IPBased:   true,
		UserBased: false,
	}
	middleware.limiters["/auth/login/2fa/verify"] = RateLimiter{
		Limit:     5,
		Window:    60, // 5 requests per minute
		IPBased:   true,
		UserBased: false,
	}
	middleware.limiters["/auth/register"] = RateLimiter{
		Limit:     3,
		Window:    300, // 3 requests per 5 minutes
		IPBased:   true,
		UserBased: false,
	}
	middleware.limiters["/auth/refresh-token"] = RateLimiter{
		Limit:     20,
		Window:    60, // 20 requests per minute
		IPBased:   true,
		UserBased: false,
	}
	
	// KYC submission
	middleware.limiters["/user/kyc"] = RateLimiter{
		Limit:     5,
		Window:    300, // 5 requests per 5 minutes
		IPBased:   false,
		UserBased: true,
	}
	
	// Invoice uploads
	middleware.limiters["/invoices"] = RateLimiter{
		Limit:     20,
		Window:    60, // 20 requests per minute
		IPBased:   false,
		UserBased: true,
	}
	
	// Admin endpoints
	middleware.limiters["/admin/"] = RateLimiter{
		Limit:     100,
		Window:    60, // 100 requests per minute
		IPBased:   false,
		UserBased: true,
	}

	// Start a goroutine to clean up expired rate limit entries
	go middleware.cleanupExpiredEntries()

	return middleware
}

// cleanupExpiredEntries periodically removes expired rate limit entries
func (rlm *RateLimitMiddleware) cleanupExpiredEntries() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rlm.mu.Lock()
		
		now := time.Now().Unix()
		
		// Clean up IP-based limits
		for ip, windows := range rlm.ipLimits {
			for windowStart := range windows {
				if windowStart < now-3600 { // Remove entries older than 1 hour
					delete(windows, windowStart)
				}
			}
			if len(windows) == 0 {
				delete(rlm.ipLimits, ip)
			}
		}
		
		// Clean up user-based limits
		for userID, windows := range rlm.userLimits {
			for windowStart := range windows {
				if windowStart < now-3600 { // Remove entries older than 1 hour
					delete(windows, windowStart)
				}
			}
			if len(windows) == 0 {
				delete(rlm.userLimits, userID)
			}
		}
		
		rlm.mu.Unlock()
	}
}

// updateRequestMetrics periodically updates request metrics for adaptive rate limiting
func (rlm *RateLimitMiddleware) updateRequestMetrics() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rlm.metricsMu.Lock()
		
		// Reset requests per minute counter
		atomic.StoreInt64(&rlm.requestsPerMinute, 0)
		rlm.lastUpdated = time.Now()
		
		rlm.metricsMu.Unlock()
	}
}

// recordRequest records a request and its response time for adaptive rate limiting
func (rlm *RateLimitMiddleware) recordRequest(startTime time.Time) {
	// Increment total request count
	atomic.AddInt64(&rlm.requestCount, 1)
	
	// Increment requests per minute
	atomic.AddInt64(&rlm.requestsPerMinute, 1)
	
	// Calculate response time in microseconds
	responseTime := time.Since(startTime).Microseconds()
	
	// Update average response time using exponential moving average
	rlm.metricsMu.Lock()
	defer rlm.metricsMu.Unlock()
	
	if rlm.avgResponseTime == 0 {
		rlm.avgResponseTime = responseTime
	} else {
		// Use a weight of 0.1 for the new value
		rlm.avgResponseTime = (rlm.avgResponseTime * 9 + responseTime) / 10
	}
}

// calculateLoadFactor calculates a load factor based on request metrics
// Returns a value between 0.0 and 1.0, where higher values indicate higher load
func (rlm *RateLimitMiddleware) calculateLoadFactor() float64 {
	rlm.metricsMu.RLock()
	defer rlm.metricsMu.RUnlock()
	
	// Get current requests per minute
	rpm := atomic.LoadInt64(&rlm.requestsPerMinute)
	
	// Calculate load factor based on requests per minute and average response time
	// This is a simple heuristic that can be adjusted based on application characteristics
	rpmFactor := float64(rpm) / 1000.0 // Normalize to 0-1 range (assuming 1000 rpm is high load)
	if rpmFactor > 1.0 {
		rpmFactor = 1.0
	}
	
	timeFactor := float64(rlm.avgResponseTime) / 100000.0 // Normalize to 0-1 range (assuming 100ms is high load)
	if timeFactor > 1.0 {
		timeFactor = 1.0
	}
	
	// Combine factors (giving more weight to response time)
	loadFactor := 0.4*rpmFactor + 0.6*timeFactor
	
	return loadFactor
}

// adjustRateLimit adjusts the rate limit based on the load factor
func (rlm *RateLimitMiddleware) adjustRateLimit(limiter RateLimiter) int {
	// If adaptive rate limiting is not enabled for this limiter, return the static limit
	if !limiter.Adaptive {
		return limiter.Limit
	}
	
	// Calculate current load factor
	loadFactor := rlm.calculateLoadFactor()
	
	// Adjust limit based on load factor
	// At load factor 0, use max limit
	// At load factor 1, use min limit
	// In between, interpolate linearly
	adjustedLimit := int(float64(limiter.MaxLimit) - loadFactor*float64(limiter.MaxLimit-limiter.MinLimit))
	
	// Ensure limit is within bounds
	if adjustedLimit < limiter.MinLimit {
		adjustedLimit = limiter.MinLimit
	}
	if adjustedLimit > limiter.MaxLimit {
		adjustedLimit = limiter.MaxLimit
	}
	
	return adjustedLimit
}

// getRealIP extracts the real client IP from request headers
func (rlm *RateLimitMiddleware) getRealIP(c *fiber.Ctx) string {
	// Check for X-Forwarded-For header
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs (client, proxy1, proxy2, ...)
		// The first one is the original client IP
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			// Trim whitespace and return the first IP
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check for X-Real-IP header
	if xrip := c.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	
	// Fall back to the direct connection IP
	return c.IP()
}

// RateLimit returns a middleware function that applies rate limiting
func (rlm *RateLimitMiddleware) RateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record start time for adaptive rate limiting
		startTime := time.Now()
		
		path := c.Path()
		
		// Find the most specific rate limiter for the path
		var limiter RateLimiter
		var found bool
		
		// Check for exact path match
		limiter, found = rlm.limiters[path]
		
		// If not found, check for prefix match
		if !found {
			for pattern, l := range rlm.limiters {
				if len(pattern) > 0 && pattern[len(pattern)-1] == '/' && len(path) >= len(pattern) && path[:len(pattern)] == pattern {
					limiter = l
					found = true
					break
				}
			}
		}
		
		// If no rate limiter is configured for this path, skip rate limiting
		if !found {
			// Record metrics even for non-rate-limited paths
			go rlm.recordRequest(startTime)
			return c.Next()
		}
		
		// Adjust rate limit if adaptive rate limiting is enabled
		effectiveLimit := limiter.Limit
		if limiter.Adaptive {
			effectiveLimit = rlm.adjustRateLimit(limiter)
		}
		
		// Apply IP-based rate limiting
		if limiter.IPBased {
			// Get the real client IP, handling proxies
			ip := rlm.getRealIP(c)
			
			// Apply rate limiting with the real IP
			if !rlm.checkRateLimit(ip, effectiveLimit, limiter.Window, true) {
				// Record metrics even for rejected requests
				go rlm.recordRequest(startTime)
				return utils.HandleError(c, fiber.StatusTooManyRequests, 
					fmt.Sprintf("Rate limit exceeded. Try again in %d seconds", limiter.Window), nil)
			}
		}
		
		// Apply user-based rate limiting
		if limiter.UserBased {
			userID := c.Locals("user_id")
			if userID != nil {
				userIDStr := fmt.Sprintf("%v", userID)
				if !rlm.checkRateLimit(userIDStr, effectiveLimit, limiter.Window, false) {
					// Record metrics even for rejected requests
					go rlm.recordRequest(startTime)
					return utils.HandleError(c, fiber.StatusTooManyRequests, 
						fmt.Sprintf("Rate limit exceeded. Try again in %d seconds", limiter.Window), nil)
				}
			}
		}
		
		// Set rate limit headers
		c.Set("X-RateLimit-Limit", strconv.Itoa(effectiveLimit))
		c.Set("X-RateLimit-Window", strconv.Itoa(limiter.Window))
		
		// Create a wrapper handler to record metrics after the request is processed
		err := c.Next()
		
		// Record metrics after the request is processed
		go rlm.recordRequest(startTime)
		
		return err
	}
}

// checkRateLimit checks if the request is within the rate limit
func (rlm *RateLimitMiddleware) checkRateLimit(key string, limit int, window int, isIP bool) bool {
	// Use Redis if available, otherwise fall back to in-memory
	if rlm.useRedis {
		return rlm.checkRateLimitRedis(key, limit, window, isIP)
	} else {
		return rlm.checkRateLimitMemory(key, limit, window, isIP)
	}
}

// checkRateLimitRedis implements a sliding window rate limit algorithm using Redis
func (rlm *RateLimitMiddleware) checkRateLimitRedis(key string, limit int, window int, isIP bool) bool {
	// Create a Redis key with the appropriate prefix
	var redisKey string
	if isIP {
		redisKey = ipRateLimitKeyPrefix + key
	} else {
		redisKey = userRateLimitKeyPrefix + key
	}
	
	now := time.Now().UnixNano() / int64(time.Millisecond)
	windowMs := int64(window * 1000) // Convert seconds to milliseconds
	
	// Remove timestamps older than the current window
	clearBefore := now - windowMs
	
	// Use Redis pipeline to execute multiple commands atomically
	pipe := rlm.redisClient.Pipeline()
	
	// Add the current timestamp to the sorted set
	pipe.ZAdd(rlm.ctx, redisKey, &redis.Z{
		Score:  float64(now),
		Member: now,
	})
	
	// Remove timestamps outside the current window
	pipe.ZRemRangeByScore(rlm.ctx, redisKey, "0", fmt.Sprintf("%d", clearBefore))
	
	// Count the number of requests in the current window
	countCmd := pipe.ZCount(rlm.ctx, redisKey, fmt.Sprintf("%d", clearBefore), "+inf")
	
	// Set expiration on the key to clean up automatically
	pipe.Expire(rlm.ctx, redisKey, time.Duration(window*2)*time.Second)
	
	// Execute the pipeline
	_, err := pipe.Exec(rlm.ctx)
	if err != nil {
		// If Redis fails, fall back to in-memory rate limiting
		return rlm.checkRateLimitMemory(key, limit, window, isIP)
	}
	
	// Get the count result
	count, err := countCmd.Result()
	if err != nil {
		// If Redis fails, fall back to in-memory rate limiting
		return rlm.checkRateLimitMemory(key, limit, window, isIP)
	}
	
	// Check if the count is within the limit
	return count <= int64(limit)
}

// checkRateLimitMemory implements a fixed window rate limit algorithm using in-memory storage
func (rlm *RateLimitMiddleware) checkRateLimitMemory(key string, limit int, window int, isIP bool) bool {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()
	
	now := time.Now().Unix()
	windowStart := now / int64(window) * int64(window)
	
	var limits map[string]map[int64]int
	if isIP {
		limits = rlm.ipLimits
	} else {
		limits = rlm.userLimits
	}
	
	if _, exists := limits[key]; !exists {
		limits[key] = make(map[int64]int)
	}
	
	if _, exists := limits[key][windowStart]; !exists {
		limits[key][windowStart] = 0
	}
	
	limits[key][windowStart]++
	
	return limits[key][windowStart] <= limit
}

// GetRateLimitInfo returns the current rate limit information for a key
func (rlm *RateLimitMiddleware) GetRateLimitInfo(key string, window int, isIP bool) (int, int) {
	// Use Redis if available, otherwise fall back to in-memory
	if rlm.useRedis {
		return rlm.getRateLimitInfoRedis(key, window, isIP)
	} else {
		return rlm.getRateLimitInfoMemory(key, window, isIP)
	}
}

// getRateLimitInfoRedis gets rate limit information from Redis
func (rlm *RateLimitMiddleware) getRateLimitInfoRedis(key string, window int, isIP bool) (int, int) {
	// Create a Redis key with the appropriate prefix
	var redisKey string
	if isIP {
		redisKey = ipRateLimitKeyPrefix + key
	} else {
		redisKey = userRateLimitKeyPrefix + key
	}
	
	now := time.Now().UnixNano() / int64(time.Millisecond)
	windowMs := int64(window * 1000) // Convert seconds to milliseconds
	clearBefore := now - windowMs
	
	// Count the number of requests in the current window
	count, err := rlm.redisClient.ZCount(rlm.ctx, redisKey, fmt.Sprintf("%d", clearBefore), "+inf").Result()
	if err != nil {
		// If Redis fails, fall back to in-memory rate limiting
		return rlm.getRateLimitInfoMemory(key, window, isIP)
	}
	
	// Calculate remaining time in the window (in seconds)
	remainingTime := window - int((now/1000)%int64(window))
	
	return int(count), remainingTime
}

// getRateLimitInfoMemory gets rate limit information from in-memory storage
func (rlm *RateLimitMiddleware) getRateLimitInfoMemory(key string, window int, isIP bool) (int, int) {
	rlm.mu.RLock()
	defer rlm.mu.RUnlock()
	
	now := time.Now().Unix()
	windowStart := now / int64(window) * int64(window)
	
	var limits map[string]map[int64]int
	if isIP {
		limits = rlm.ipLimits
	} else {
		limits = rlm.userLimits
	}
	
	if _, exists := limits[key]; !exists {
		return 0, window - int(now-windowStart)
	}
	
	if count, exists := limits[key][windowStart]; exists {
		return count, window - int(now-windowStart)
	}
	
	return 0, window
}