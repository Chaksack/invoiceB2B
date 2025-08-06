package services

import (
	"context"
	"encoding/json"
	"fmt"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuditLogService provides comprehensive security audit logging
type AuditLogService interface {
	// LogSecurityEvent logs a security-related event
	LogSecurityEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, details map[string]interface{}) error
	
	// LogAuthenticationEvent logs an authentication event (login, logout, 2FA, etc.)
	LogAuthenticationEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, success bool, details map[string]interface{}) error
	
	// LogAccessEvent logs an access control event (permission check, resource access, etc.)
	LogAccessEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, resource string, resourceID string, action string, allowed bool) error
	
	// LogDataEvent logs a data-related event (read, write, delete, etc.)
	LogDataEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, dataType string, dataID string, action string, details map[string]interface{}) error
	
	// LogAdminEvent logs an administrative action
	LogAdminEvent(ctx context.Context, eventType string, staffID uint, targetType string, targetID string, action string, details map[string]interface{}) error
	
 // GetAuditLogs retrieves audit logs with filtering and pagination
	GetAuditLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ActivityLog, int64, error)
}

type auditLogService struct {
	activityLogRepo repositories.ActivityLogRepository
}

// NewAuditLogService creates a new audit log service
func NewAuditLogService(activityLogRepo repositories.ActivityLogRepository) AuditLogService {
	return &auditLogService{
		activityLogRepo: activityLogRepo,
	}
}

// LogSecurityEvent logs a security-related event
func (s *auditLogService) LogSecurityEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, details map[string]interface{}) error {
	// Extract context information
	var ip, userAgent, requestID, sessionID string
	if fiberCtx, ok := ctx.Value("fiber_ctx").(*fiber.Ctx); ok {
		ip = fiberCtx.IP()
		userAgent = fiberCtx.Get("User-Agent")
		requestID = fiberCtx.Get("X-Request-ID")
		sessionID = fiberCtx.Get("X-Session-ID")
		
		// Generate request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
			fiberCtx.Set("X-Request-ID", requestID)
		}
	}
	
	// Add context information to details
	if details == nil {
		details = make(map[string]interface{})
	}
	
	details["ip"] = ip
	details["user_agent"] = userAgent
	details["request_id"] = requestID
	details["session_id"] = sessionID
	details["timestamp"] = time.Now().Format(time.RFC3339)
	
	// Convert details to JSON
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}
	
 // Create activity log entry
	// Add description to details
	details["description"] = fmt.Sprintf("Security event: %s", eventType)
	details["user_agent"] = userAgent
	
	// Update details JSON
	detailsJSON, err = json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal updated details: %w", err)
	}
	
	// Create IP address pointer
	var ipPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	
	activityLog := &models.ActivityLog{
		UserID:     userID,
		StaffID:    staffID,
		Action:     fmt.Sprintf("SECURITY_%s", eventType),
		Details:    string(detailsJSON),
		IPAddress:  ipPtr,
	}
	
	// Save to database
	return s.activityLogRepo.Create(ctx, activityLog)
}

// LogAuthenticationEvent logs an authentication event
func (s *auditLogService) LogAuthenticationEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, success bool, details map[string]interface{}) error {
	// Extract context information
	var ip, userAgent, requestID, sessionID string
	if fiberCtx, ok := ctx.Value("fiber_ctx").(*fiber.Ctx); ok {
		ip = fiberCtx.IP()
		userAgent = fiberCtx.Get("User-Agent")
		requestID = fiberCtx.Get("X-Request-ID")
		sessionID = fiberCtx.Get("X-Session-ID")
		
		// Generate request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
			fiberCtx.Set("X-Request-ID", requestID)
		}
	}
	
	// Add context information to details
	if details == nil {
		details = make(map[string]interface{})
	}
	
	details["ip"] = ip
	details["user_agent"] = userAgent
	details["request_id"] = requestID
	details["session_id"] = sessionID
	details["timestamp"] = time.Now().Format(time.RFC3339)
	details["success"] = success
	
	// Convert details to JSON
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}
	
 // Create activity log entry
	result := "SUCCESS"
	if !success {
		result = "FAILURE"
	}
	
	// Add description to details
	details["description"] = fmt.Sprintf("Authentication event: %s (%s)", eventType, result)
	
	// Update details JSON
	detailsJSON, err = json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal updated details: %w", err)
	}
	
	// Create IP address pointer
	var ipPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	
	activityLog := &models.ActivityLog{
		UserID:     userID,
		StaffID:    staffID,
		Action:     fmt.Sprintf("AUTH_%s_%s", eventType, result),
		Details:    string(detailsJSON),
		IPAddress:  ipPtr,
	}
	
	// Save to database
	return s.activityLogRepo.Create(ctx, activityLog)
}

// LogAccessEvent logs an access control event
func (s *auditLogService) LogAccessEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, resource string, resourceID string, action string, allowed bool) error {
	// Extract context information
	var ip, userAgent, requestID, sessionID string
	if fiberCtx, ok := ctx.Value("fiber_ctx").(*fiber.Ctx); ok {
		ip = fiberCtx.IP()
		userAgent = fiberCtx.Get("User-Agent")
		requestID = fiberCtx.Get("X-Request-ID")
		sessionID = fiberCtx.Get("X-Session-ID")
		
		// Generate request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
			fiberCtx.Set("X-Request-ID", requestID)
		}
	}
	
	// Create details
	details := map[string]interface{}{
		"ip":          ip,
		"user_agent":  userAgent,
		"request_id":  requestID,
		"session_id":  sessionID,
		"timestamp":   time.Now().Format(time.RFC3339),
		"resource":    resource,
		"resource_id": resourceID,
		"action":      action,
		"allowed":     allowed,
	}
	
	// Convert details to JSON
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}
	
	// Create activity log entry
	result := "ALLOWED"
	if !allowed {
		result = "DENIED"
	}
	
	// Add description to details
	details["description"] = fmt.Sprintf("Access event: %s %s %s (%s)", action, resource, resourceID, result)
	
	// Update details JSON
	detailsJSON, err = json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal updated details: %w", err)
	}
	
	// Create IP address pointer
	var ipPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	
	activityLog := &models.ActivityLog{
		UserID:     userID,
		StaffID:    staffID,
		Action:     fmt.Sprintf("ACCESS_%s_%s", eventType, result),
		Details:    string(detailsJSON),
		IPAddress:  ipPtr,
	}
	
	// Save to database
	return s.activityLogRepo.Create(ctx, activityLog)
}

// LogDataEvent logs a data-related event
func (s *auditLogService) LogDataEvent(ctx context.Context, eventType string, userID *uint, staffID *uint, dataType string, dataID string, action string, details map[string]interface{}) error {
	// Extract context information
	var ip, userAgent, requestID, sessionID string
	if fiberCtx, ok := ctx.Value("fiber_ctx").(*fiber.Ctx); ok {
		ip = fiberCtx.IP()
		userAgent = fiberCtx.Get("User-Agent")
		requestID = fiberCtx.Get("X-Request-ID")
		sessionID = fiberCtx.Get("X-Session-ID")
		
		// Generate request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
			fiberCtx.Set("X-Request-ID", requestID)
		}
	}
	
	// Add context information to details
	if details == nil {
		details = make(map[string]interface{})
	}
	
	details["ip"] = ip
	details["user_agent"] = userAgent
	details["request_id"] = requestID
	details["session_id"] = sessionID
	details["timestamp"] = time.Now().Format(time.RFC3339)
	details["data_type"] = dataType
	details["data_id"] = dataID
	details["action"] = action
	
	// Convert details to JSON
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}
	
	// Add description to details
	details["description"] = fmt.Sprintf("Data event: %s %s %s", action, dataType, dataID)
	
	// Update details JSON
	detailsJSON, err = json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal updated details: %w", err)
	}
	
	// Create IP address pointer
	var ipPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	
	// Create activity log entry
	activityLog := &models.ActivityLog{
		UserID:     userID,
		StaffID:    staffID,
		Action:     fmt.Sprintf("DATA_%s_%s", action, eventType),
		Details:    string(detailsJSON),
		IPAddress:  ipPtr,
	}
	
	// Save to database
	return s.activityLogRepo.Create(ctx, activityLog)
}

// LogAdminEvent logs an administrative action
func (s *auditLogService) LogAdminEvent(ctx context.Context, eventType string, staffID uint, targetType string, targetID string, action string, details map[string]interface{}) error {
	// Extract context information
	var ip, userAgent, requestID, sessionID string
	if fiberCtx, ok := ctx.Value("fiber_ctx").(*fiber.Ctx); ok {
		ip = fiberCtx.IP()
		userAgent = fiberCtx.Get("User-Agent")
		requestID = fiberCtx.Get("X-Request-ID")
		sessionID = fiberCtx.Get("X-Session-ID")
		
		// Generate request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
			fiberCtx.Set("X-Request-ID", requestID)
		}
	}
	
	// Add context information to details
	if details == nil {
		details = make(map[string]interface{})
	}
	
	details["ip"] = ip
	details["user_agent"] = userAgent
	details["request_id"] = requestID
	details["session_id"] = sessionID
	details["timestamp"] = time.Now().Format(time.RFC3339)
	details["target_type"] = targetType
	details["target_id"] = targetID
	details["action"] = action
	
	// Convert details to JSON
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}
	
	// Add description to details
	details["description"] = fmt.Sprintf("Admin event: %s %s %s", action, targetType, targetID)
	
	// Update details JSON
	detailsJSON, err = json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal updated details: %w", err)
	}
	
	// Create IP address pointer
	var ipPtr *string
	if ip != "" {
		ipPtr = &ip
	}
	
	// Create activity log entry
	staffIDPtr := &staffID
	activityLog := &models.ActivityLog{
		StaffID:    staffIDPtr,
		Action:     fmt.Sprintf("ADMIN_%s_%s", action, eventType),
		Details:    string(detailsJSON),
		IPAddress:  ipPtr,
	}
	
	// Save to database
	return s.activityLogRepo.Create(ctx, activityLog)
}

// GetAuditLogs retrieves audit logs with filtering and pagination
func (s *auditLogService) GetAuditLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ActivityLog, int64, error) {
	return s.activityLogRepo.FindAll(ctx, page, pageSize, filters)
}