package models

import (
	"time"
)

// Policy represents an authorization policy for attribute-based access control
type Policy struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PolicyID    string    `gorm:"uniqueIndex" json:"policy_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Effect      string    `json:"effect"` // "allow" or "deny"
	
	// Subject conditions (who)
	SubjectConditions string `gorm:"type:text" json:"subject_conditions"`
	
	// Action conditions (what)
	Actions string `gorm:"type:text" json:"actions"`
	
	// Resource conditions (which)
	ResourceConditions string `gorm:"type:text" json:"resource_conditions"`
	
	// Context conditions (when/where)
	ContextConditions string `gorm:"type:text" json:"context_conditions"`
	
	// Priority determines the order of evaluation (higher numbers evaluated first)
	Priority int `json:"priority"`
	
	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint      `json:"created_by"`
	Active    bool      `json:"active"`
}

// AuthorizationLog represents a log entry for an authorization decision
type AuthorizationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Timestamp time.Time `json:"timestamp"`
	
	// Who
	SubjectID   string `json:"subject_id"`
	SubjectType string `json:"subject_type"` // "user", "staff", "system", etc.
	SubjectData string `gorm:"type:text" json:"subject_data"`
	
	// What
	Action string `json:"action"`
	
	// Which
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	ResourceData string `gorm:"type:text" json:"resource_data"`
	
	// Result
	Decision bool   `json:"decision"` // true = allowed, false = denied
	Reason   string `json:"reason"`
	
	// Context
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	SessionID string    `json:"session_id"`
	RequestID string    `json:"request_id"`
	PolicyID  string    `json:"policy_id"`
}