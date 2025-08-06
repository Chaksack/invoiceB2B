package interfaces

import "context"

// AuthorizationService defines the interface for authorization services
type AuthorizationService interface {
	// Authorize checks if the subject is allowed to perform the action on the resource
	Authorize(ctx context.Context, subject map[string]interface{}, action string, resource map[string]interface{}) (bool, error)
	
	// AddPolicy adds a new authorization policy
	AddPolicy(policyID string, policy map[string]interface{}) error
	
	// RemovePolicy removes an authorization policy
	RemovePolicy(policyID string) error
	
	// GetPolicy retrieves an authorization policy
	GetPolicy(policyID string) (map[string]interface{}, error)
	
	// LogAuthorizationDecision logs an authorization decision for audit purposes
	LogAuthorizationDecision(ctx context.Context, subject map[string]interface{}, action string, resource map[string]interface{}, decision bool, reason string) error
}