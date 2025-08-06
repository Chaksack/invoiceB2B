package services

import (
	"context"
	"encoding/json"
	"fmt"
	"invoiceB2B/internal/interfaces"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthorizationService implements the interfaces.AuthorizationService interface
type authorizationService struct {
	policyRepo repositories.PolicyRepository
	activityLogService ActivityLogService
}

// NewAuthorizationService creates a new authorization service
func NewAuthorizationService(policyRepo repositories.PolicyRepository, activityLogService ActivityLogService) interfaces.AuthorizationService {
	return &authorizationService{
		policyRepo: policyRepo,
		activityLogService: activityLogService,
	}
}

// Authorize checks if the subject is allowed to perform the action on the resource
func (s *authorizationService) Authorize(ctx context.Context, subject map[string]interface{}, action string, resource map[string]interface{}) (bool, error) {
	// Get all active policies
	policies, err := s.policyRepo.FindAll(ctx, true)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve policies: %w", err)
	}

	// Extract context information from the Fiber context if available
	var ip, userAgent, sessionID, requestID string
	if fiberCtx, ok := ctx.Value("fiber_ctx").(*fiber.Ctx); ok {
		ip = fiberCtx.IP()
		userAgent = fiberCtx.Get("User-Agent")
		sessionID = fiberCtx.Get("X-Session-ID")
		requestID = fiberCtx.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			fiberCtx.Set("X-Request-ID", requestID)
		}
	}

	// Create context attributes map
	contextAttrs := map[string]interface{}{
		"ip":         ip,
		"user_agent": userAgent,
		"session_id": sessionID,
		"request_id": requestID,
		"timestamp":  time.Now().Unix(),
	}

	// Evaluate policies in order of priority
	var matchedPolicy *models.Policy
	var decision bool
	var reason string

	for _, policy := range policies {
		// Check if policy matches
		matches, err := s.policyMatches(policy, subject, action, resource, contextAttrs)
		if err != nil {
			return false, fmt.Errorf("error evaluating policy %s: %w", policy.PolicyID, err)
		}

		if matches {
			matchedPolicy = &policy
			decision = strings.ToLower(policy.Effect) == "allow"
			if decision {
				reason = fmt.Sprintf("Allowed by policy %s", policy.PolicyID)
			} else {
				reason = fmt.Sprintf("Denied by policy %s", policy.PolicyID)
			}
			break
		}
	}

	// If no policy matched, deny by default
	if matchedPolicy == nil {
		decision = false
		reason = "No matching policy found, denied by default"
	}

	// Log the authorization decision
	err = s.LogAuthorizationDecision(ctx, subject, action, resource, decision, reason)
	if err != nil {
		// Log the error but don't fail the authorization check
		log.Printf("Failed to log authorization decision: %v", err)
	}

	return decision, nil
}

// policyMatches checks if a policy matches the given attributes
func (s *authorizationService) policyMatches(policy models.Policy, subject map[string]interface{}, action string, resource map[string]interface{}, contextAttrs map[string]interface{}) (bool, error) {
	// Check subject conditions
	if policy.SubjectConditions != "" {
		var subjectConditions map[string]interface{}
		if err := json.Unmarshal([]byte(policy.SubjectConditions), &subjectConditions); err != nil {
			return false, fmt.Errorf("invalid subject conditions: %w", err)
		}

		if !s.matchesConditions(subject, subjectConditions) {
			return false, nil
		}
	}

	// Check action
	if policy.Actions != "" {
		var actions []string
		if err := json.Unmarshal([]byte(policy.Actions), &actions); err != nil {
			return false, fmt.Errorf("invalid actions: %w", err)
		}

		actionMatched := false
		for _, a := range actions {
			if a == action || a == "*" {
				actionMatched = true
				break
			}
		}

		if !actionMatched {
			return false, nil
		}
	}

	// Check resource conditions
	if policy.ResourceConditions != "" {
		var resourceConditions map[string]interface{}
		if err := json.Unmarshal([]byte(policy.ResourceConditions), &resourceConditions); err != nil {
			return false, fmt.Errorf("invalid resource conditions: %w", err)
		}

		if !s.matchesConditions(resource, resourceConditions) {
			return false, nil
		}
	}

	// Check context conditions
	if policy.ContextConditions != "" {
		var contextConditions map[string]interface{}
		if err := json.Unmarshal([]byte(policy.ContextConditions), &contextConditions); err != nil {
			return false, fmt.Errorf("invalid context conditions: %w", err)
		}

		if !s.matchesConditions(contextAttrs, contextConditions) {
			return false, nil
		}
	}

	// All conditions matched
	return true, nil
}

// matchesConditions checks if attributes match the given conditions
func (s *authorizationService) matchesConditions(attrs map[string]interface{}, conditions map[string]interface{}) bool {
	for key, condition := range conditions {
		// Get the attribute value
		attrValue, exists := attrs[key]
		if !exists {
			return false
		}

		// Handle different condition types
		switch cond := condition.(type) {
		case map[string]interface{}:
			// Complex condition with operators
			for op, value := range cond {
				switch op {
				case "eq":
					if attrValue != value {
						return false
					}
				case "ne":
					if attrValue == value {
						return false
					}
				case "gt":
					// Convert to comparable types
					if !s.compareGreaterThan(attrValue, value) {
						return false
					}
				case "lt":
					// Convert to comparable types
					if !s.compareLessThan(attrValue, value) {
						return false
					}
				case "in":
					// Check if value is in array
					if !s.isInArray(attrValue, value) {
						return false
					}
				case "contains":
					// Check if string contains substring
					if !s.containsString(attrValue, value) {
						return false
					}
				}
			}
		default:
			// Simple equality check
			if attrValue != condition {
				return false
			}
		}
	}

	return true
}

// Helper functions for comparisons
func (s *authorizationService) compareGreaterThan(a, b interface{}) bool {
	// Convert to float64 for numeric comparison
	aFloat, aOk := a.(float64)
	bFloat, bOk := b.(float64)

	if aOk && bOk {
		return aFloat > bFloat
	}

	// Try string comparison
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)

	if aOk && bOk {
		return aStr > bStr
	}

	return false
}

func (s *authorizationService) compareLessThan(a, b interface{}) bool {
	// Convert to float64 for numeric comparison
	aFloat, aOk := a.(float64)
	bFloat, bOk := b.(float64)

	if aOk && bOk {
		return aFloat < bFloat
	}

	// Try string comparison
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)

	if aOk && bOk {
		return aStr < bStr
	}

	return false
}

func (s *authorizationService) isInArray(item, array interface{}) bool {
	arr, ok := array.([]interface{})
	if !ok {
		return false
	}

	for _, v := range arr {
		if v == item {
			return true
		}
	}

	return false
}

func (s *authorizationService) containsString(str, substr interface{}) bool {
	strVal, strOk := str.(string)
	substrVal, substrOk := substr.(string)

	if strOk && substrOk {
		return strings.Contains(strVal, substrVal)
	}

	return false
}

// AddPolicy adds a new authorization policy
func (s *authorizationService) AddPolicy(policyID string, policyData map[string]interface{}) error {
	// Create a new policy
	policy := models.Policy{
		PolicyID:    policyID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Active:      true,
	}
	
	// Manually assign fields from the map
	if name, ok := policyData["name"].(string); ok {
		policy.Name = name
	}
	
	if description, ok := policyData["description"].(string); ok {
		policy.Description = description
	}
	
	if effect, ok := policyData["effect"].(string); ok {
		policy.Effect = effect
	}
	
	if priority, ok := policyData["priority"].(float64); ok {
		policy.Priority = int(priority)
	}
	
	if createdBy, ok := policyData["created_by"].(float64); ok {
		policy.CreatedBy = uint(createdBy)
	}

	// Convert complex fields to JSON strings
	if subjectConditions, ok := policyData["subject_conditions"].(map[string]interface{}); ok {
		jsonData, err := json.Marshal(subjectConditions)
		if err != nil {
			return fmt.Errorf("failed to marshal subject conditions: %w", err)
		}
		policy.SubjectConditions = string(jsonData)
	}

	if actions, ok := policyData["actions"].([]interface{}); ok {
		jsonData, err := json.Marshal(actions)
		if err != nil {
			return fmt.Errorf("failed to marshal actions: %w", err)
		}
		policy.Actions = string(jsonData)
	}

	if resourceConditions, ok := policyData["resource_conditions"].(map[string]interface{}); ok {
		jsonData, err := json.Marshal(resourceConditions)
		if err != nil {
			return fmt.Errorf("failed to marshal resource conditions: %w", err)
		}
		policy.ResourceConditions = string(jsonData)
	}

	if contextConditions, ok := policyData["context_conditions"].(map[string]interface{}); ok {
		jsonData, err := json.Marshal(contextConditions)
		if err != nil {
			return fmt.Errorf("failed to marshal context conditions: %w", err)
		}
		policy.ContextConditions = string(jsonData)
	}

	// Save to database
	_, err := s.policyRepo.Create(context.Background(), &policy)
	if err != nil {
		return fmt.Errorf("failed to create policy: %w", err)
	}

	return nil
}

// RemovePolicy removes an authorization policy
func (s *authorizationService) RemovePolicy(policyID string) error {
	policy, err := s.policyRepo.FindByPolicyID(context.Background(), policyID)
	if err != nil {
		return fmt.Errorf("failed to find policy: %w", err)
	}

	return s.policyRepo.Delete(context.Background(), policy.ID)
}

// GetPolicy retrieves an authorization policy
func (s *authorizationService) GetPolicy(policyID string) (map[string]interface{}, error) {
	policy, err := s.policyRepo.FindByPolicyID(context.Background(), policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to find policy: %w", err)
	}

	// Convert to map
	result := map[string]interface{}{
		"id":          policy.ID,
		"policy_id":   policy.PolicyID,
		"name":        policy.Name,
		"description": policy.Description,
		"effect":      policy.Effect,
		"priority":    policy.Priority,
		"created_at":  policy.CreatedAt,
		"updated_at":  policy.UpdatedAt,
		"created_by":  policy.CreatedBy,
		"active":      policy.Active,
	}

	// Parse JSON fields
	if policy.SubjectConditions != "" {
		var subjectConditions map[string]interface{}
		if err := json.Unmarshal([]byte(policy.SubjectConditions), &subjectConditions); err == nil {
			result["subject_conditions"] = subjectConditions
		}
	}

	if policy.Actions != "" {
		var actions []string
		if err := json.Unmarshal([]byte(policy.Actions), &actions); err == nil {
			result["actions"] = actions
		}
	}

	if policy.ResourceConditions != "" {
		var resourceConditions map[string]interface{}
		if err := json.Unmarshal([]byte(policy.ResourceConditions), &resourceConditions); err == nil {
			result["resource_conditions"] = resourceConditions
		}
	}

	if policy.ContextConditions != "" {
		var contextConditions map[string]interface{}
		if err := json.Unmarshal([]byte(policy.ContextConditions), &contextConditions); err == nil {
			result["context_conditions"] = contextConditions
		}
	}

	return result, nil
}

// LogAuthorizationDecision logs an authorization decision for audit purposes
func (s *authorizationService) LogAuthorizationDecision(ctx context.Context, subject map[string]interface{}, action string, resource map[string]interface{}, decision bool, reason string) error {
	// Extract subject information
	subjectID := ""
	subjectType := ""
	
	if id, ok := subject["id"].(string); ok {
		subjectID = id
	} else if id, ok := subject["id"].(float64); ok {
		subjectID = fmt.Sprintf("%d", int(id))
	}
	
	if typ, ok := subject["type"].(string); ok {
		subjectType = typ
	}
	
	// Extract resource information
	resourceID := ""
	resourceType := ""
	
	if id, ok := resource["id"].(string); ok {
		resourceID = id
	} else if id, ok := resource["id"].(float64); ok {
		resourceID = fmt.Sprintf("%d", int(id))
	}
	
	if typ, ok := resource["type"].(string); ok {
		resourceType = typ
	}
	
	// Convert maps to JSON for storage
	subjectJSON, _ := json.Marshal(subject)
	resourceJSON, _ := json.Marshal(resource)
	
	// Extract context information
	var ip, userAgent, sessionID, requestID, policyID string
	if fiberCtx, ok := ctx.Value("fiber_ctx").(*fiber.Ctx); ok {
		ip = fiberCtx.IP()
		userAgent = fiberCtx.Get("User-Agent")
		sessionID = fiberCtx.Get("X-Session-ID")
		requestID = fiberCtx.Get("X-Request-ID")
	}
	
	// Extract policy ID from reason if available
	if strings.Contains(reason, "policy") {
		parts := strings.Split(reason, " ")
		for i, part := range parts {
			if part == "policy" && i < len(parts)-1 {
				policyID = parts[i+1]
				break
			}
		}
	}
	
	// Create log entry
	authLog := &models.AuthorizationLog{
		Timestamp:    time.Now(),
		SubjectID:    subjectID,
		SubjectType:  subjectType,
		SubjectData:  string(subjectJSON),
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ResourceData: string(resourceJSON),
		Decision:     decision,
		Reason:       reason,
		IP:           ip,
		UserAgent:    userAgent,
		SessionID:    sessionID,
		RequestID:    requestID,
		PolicyID:     policyID,
	}
	
	// Log to database
	err := s.policyRepo.LogAuthorizationDecision(ctx, authLog)
	if err != nil {
		return fmt.Errorf("failed to log authorization decision: %w", err)
	}
	
	// Also log to activity log service
	var userID *uint
	if subjectType == "user" && subjectID != "" {
		if id, err := strconv.ParseUint(subjectID, 10, 32); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}
	
	var staffID *uint
	if subjectType == "staff" && subjectID != "" {
		if id, err := strconv.ParseUint(subjectID, 10, 32); err == nil {
			sid := uint(id)
			staffID = &sid
		}
	}
	
	actionType := "AUTHORIZATION"
	if decision {
		actionType += "_ALLOWED"
	} else {
		actionType += "_DENIED"
	}
	
	details := fmt.Sprintf("Action: %s, Resource: %s, Reason: %s", action, resourceType, reason)
	
	_ = s.activityLogService.LogActivity(ctx, staffID, userID, actionType, details, "")
	
	return nil
}