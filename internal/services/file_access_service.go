package services

import (
	"log"
	"path/filepath"
	"strings"
	"time"
)

// FileAccessLevel defines the level of access a user has to a file
type FileAccessLevel int

const (
	// NoAccess means the user has no access to the file
	NoAccess FileAccessLevel = iota
	// ReadAccess means the user can read the file
	ReadAccess
	// WriteAccess means the user can read and write to the file
	WriteAccess
	// FullAccess means the user has full control over the file
	FullAccess
)

// FileAccessEvent represents a file access event for logging
type FileAccessEvent struct {
	UserID       string
	UserRole     string
	FilePath     string
	AccessType   string // "read", "write", "delete", etc.
	AccessTime   time.Time
	AccessResult bool   // true if access was granted, false if denied
	IPAddress    string
	UserAgent    string
	Reason       string // reason for access denial, if applicable
}

// FileAccessService defines the interface for file access control
type FileAccessService interface {
	// CheckAccess checks if a user has the specified access level to a file
	CheckAccess(userID string, userRole string, filePath string, requiredAccess FileAccessLevel) (bool, error)
	
	// LogAccess logs a file access event
	LogAccess(event FileAccessEvent) error
	
	// GetAccessHistory returns the access history for a file
	GetAccessHistory(filePath string) ([]FileAccessEvent, error)
	
	// GetUserAccessHistory returns the access history for a user
	GetUserAccessHistory(userID string) ([]FileAccessEvent, error)
}

// fileAccessService implements FileAccessService
type fileAccessService struct {
	// Map of file paths to user IDs and their access levels
	// In a real implementation, this would be stored in a database
	accessMap map[string]map[string]FileAccessLevel
	
	// Access history for audit trails
	// In a real implementation, this would be stored in a database
	accessHistory []FileAccessEvent
	
	// Role-based access control rules
	// Map of roles to file path patterns and access levels
	roleRules map[string]map[string]FileAccessLevel
}

// NewFileAccessService creates a new file access service
func NewFileAccessService() FileAccessService {
	service := &fileAccessService{
		accessMap:     make(map[string]map[string]FileAccessLevel),
		accessHistory: make([]FileAccessEvent, 0),
		roleRules:     make(map[string]map[string]FileAccessLevel),
	}
	
	// Initialize default role-based access rules
	service.initializeDefaultRules()
	
	return service
}

// initializeDefaultRules sets up default access rules based on roles
func (s *fileAccessService) initializeDefaultRules() {
	// Admin role has full access to all files
	s.roleRules["admin"] = map[string]FileAccessLevel{
		"*": FullAccess,
	}
	
	// Staff role has read access to all files, write access to some directories
	s.roleRules["staff"] = map[string]FileAccessLevel{
		"*":                ReadAccess,
		"invoices/*":       WriteAccess,
		"receipts/*":       WriteAccess,
		"reports/*":        WriteAccess,
	}
	
	// User role has limited access
	s.roleRules["user"] = map[string]FileAccessLevel{
		"*":                NoAccess,
		"invoices/*":       ReadAccess,  // Can read all invoices
		"invoices/own/*":   WriteAccess, // Can write to own invoices
		"receipts/own/*":   WriteAccess, // Can write to own receipts
	}
	
	// Guest role has very limited access
	s.roleRules["guest"] = map[string]FileAccessLevel{
		"*":                NoAccess,
		"public/*":         ReadAccess,
	}
}

// CheckAccess checks if a user has the specified access level to a file
func (s *fileAccessService) CheckAccess(userID string, userRole string, filePath string, requiredAccess FileAccessLevel) (bool, error) {
	// Normalize file path
	filePath = filepath.Clean(filePath)
	
	// First check user-specific permissions (highest priority)
	if userMap, exists := s.accessMap[filePath]; exists {
		if accessLevel, exists := userMap[userID]; exists {
			return accessLevel >= requiredAccess, nil
		}
	}
	
	// If no user-specific permissions, check role-based permissions
	if roleRules, exists := s.roleRules[userRole]; exists {
		// Check for exact path match
		if accessLevel, exists := roleRules[filePath]; exists {
			return accessLevel >= requiredAccess, nil
		}
		
		// Check for pattern matches (e.g., "invoices/*")
		highestAccess := NoAccess
		for pattern, accessLevel := range roleRules {
			if pattern == "*" || matchesPattern(filePath, pattern) {
				if accessLevel > highestAccess {
					highestAccess = accessLevel
				}
			}
		}
		
		// Special handling for "own" directories
		if strings.Contains(filePath, "/own/") {
			// Extract the owner ID from the path (assuming format like "invoices/own/USER_ID/...")
			parts := strings.Split(filePath, "/")
			if len(parts) >= 3 {
				ownIndex := -1
				for i, part := range parts {
					if part == "own" && i+1 < len(parts) {
						ownIndex = i + 1
						break
					}
				}
				
				if ownIndex >= 0 && parts[ownIndex] != userID {
					// User is trying to access another user's files
					return false, nil
				}
			}
		}
		
		return highestAccess >= requiredAccess, nil
	}
	
	// Default deny if no matching rules
	return false, nil
}

// matchesPattern checks if a file path matches a pattern
func matchesPattern(filePath, pattern string) bool {
	// Handle wildcard at the end
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(filePath, prefix+"/")
	}
	
	// Handle wildcard at the beginning
	if strings.HasPrefix(pattern, "*/") {
		suffix := strings.TrimPrefix(pattern, "*/")
		return strings.HasSuffix(filePath, "/"+suffix)
	}
	
	// Handle exact match
	return filePath == pattern
}

// LogAccess logs a file access event
func (s *fileAccessService) LogAccess(event FileAccessEvent) error {
	// Set access time if not provided
	if event.AccessTime.IsZero() {
		event.AccessTime = time.Now()
	}
	
	// Log the event
	log.Printf("File access: user=%s role=%s file=%s type=%s result=%v ip=%s reason=%s",
		event.UserID, event.UserRole, event.FilePath, event.AccessType, event.AccessResult,
		event.IPAddress, event.Reason)
	
	// Store the event in the access history
	s.accessHistory = append(s.accessHistory, event)
	
	// In a real implementation, this would be stored in a database
	
	return nil
}

// GetAccessHistory returns the access history for a file
func (s *fileAccessService) GetAccessHistory(filePath string) ([]FileAccessEvent, error) {
	// Normalize file path
	filePath = filepath.Clean(filePath)
	
	// Filter access history by file path
	var history []FileAccessEvent
	for _, event := range s.accessHistory {
		if event.FilePath == filePath {
			history = append(history, event)
		}
	}
	
	return history, nil
}

// GetUserAccessHistory returns the access history for a user
func (s *fileAccessService) GetUserAccessHistory(userID string) ([]FileAccessEvent, error) {
	// Filter access history by user ID
	var history []FileAccessEvent
	for _, event := range s.accessHistory {
		if event.UserID == userID {
			history = append(history, event)
		}
	}
	
	return history, nil
}

// SetUserAccess sets the access level for a specific user to a specific file
func (s *fileAccessService) SetUserAccess(userID string, filePath string, accessLevel FileAccessLevel) error {
	// Normalize file path
	filePath = filepath.Clean(filePath)
	
	// Initialize user map if it doesn't exist
	if _, exists := s.accessMap[filePath]; !exists {
		s.accessMap[filePath] = make(map[string]FileAccessLevel)
	}
	
	// Set access level
	s.accessMap[filePath][userID] = accessLevel
	
	return nil
}

// RemoveUserAccess removes user-specific access for a file
func (s *fileAccessService) RemoveUserAccess(userID string, filePath string) error {
	// Normalize file path
	filePath = filepath.Clean(filePath)
	
	// Check if user map exists
	if userMap, exists := s.accessMap[filePath]; exists {
		// Remove user access
		delete(userMap, userID)
		
		// Remove file entry if no users left
		if len(userMap) == 0 {
			delete(s.accessMap, filePath)
		}
	}
	
	return nil
}

// AddRoleRule adds a role-based access rule
func (s *fileAccessService) AddRoleRule(role string, pattern string, accessLevel FileAccessLevel) error {
	// Initialize role rules if they don't exist
	if _, exists := s.roleRules[role]; !exists {
		s.roleRules[role] = make(map[string]FileAccessLevel)
	}
	
	// Set access level for pattern
	s.roleRules[role][pattern] = accessLevel
	
	return nil
}

// RemoveRoleRule removes a role-based access rule
func (s *fileAccessService) RemoveRoleRule(role string, pattern string) error {
	// Check if role rules exist
	if roleRules, exists := s.roleRules[role]; exists {
		// Remove pattern rule
		delete(roleRules, pattern)
		
		// Remove role entry if no patterns left
		if len(roleRules) == 0 {
			delete(s.roleRules, role)
		}
	}
	
	return nil
}