package services

import (
	"fmt"
	"os"
	"path/filepath"
)

// CheckAccess checks if a user has the specified access level to a file
func (s *fileService) CheckAccess(relativePath string, userID string, userRole string, accessLevel FileAccessLevel) (bool, error) {
	// Normalize the relative path
	relativePath = filepath.Clean(relativePath)
	
	// Use the access service to check access
	hasAccess, err := s.accessService.CheckAccess(userID, userRole, relativePath, accessLevel)
	if err != nil {
		return false, fmt.Errorf("failed to check access: %w", err)
	}
	
	return hasAccess, nil
}

// DeleteFile deletes a file with access control
func (s *fileService) DeleteFile(relativePath string, userID string, userRole string) error {
	// Check if the user has write access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, WriteAccess)
	if err != nil {
		return fmt.Errorf("failed to check access for file deletion: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "delete",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return fmt.Errorf("access denied: insufficient permissions to delete file")
	}
	
	// Get the absolute path
	absPath, err := s.GetAbsPath(relativePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	
	// Delete the file
	if err := os.Remove(absPath); err != nil {
		// Log the failure
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "delete",
			AccessResult: false,
			Reason:       fmt.Sprintf("file deletion failed: %v", err),
		})
		
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	// Log the successful deletion
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "delete",
		AccessResult: true,
	})
	
	return nil
}

// GetAccessHistory returns the access history for a file
func (s *fileService) GetAccessHistory(relativePath string, userID string, userRole string) ([]FileAccessEvent, error) {
	// Check if the user has read access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, ReadAccess)
	if err != nil {
		return nil, fmt.Errorf("failed to check access for access history: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "get_history",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return nil, fmt.Errorf("access denied: insufficient permissions to view file access history")
	}
	
	// Get the access history
	history, err := s.accessService.GetAccessHistory(relativePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get access history: %w", err)
	}
	
	// Log the successful access
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "get_history",
		AccessResult: true,
	})
	
	return history, nil
}

// GetUserAccessHistory returns the access history for a user
func (s *fileService) GetUserAccessHistory(userID string, userRole string, requestingUserID string, requestingUserRole string) ([]FileAccessEvent, error) {
	// Only allow users to view their own access history, or admins to view any user's history
	if userID != requestingUserID && requestingUserRole != "admin" {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       requestingUserID,
			UserRole:     requestingUserRole,
			FilePath:     "",
			AccessType:   "get_user_history",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return nil, fmt.Errorf("access denied: insufficient permissions to view user access history")
	}
	
	// Get the user access history
	history, err := s.accessService.GetUserAccessHistory(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user access history: %w", err)
	}
	
	// Log the successful access
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       requestingUserID,
		UserRole:     requestingUserRole,
		FilePath:     "",
		AccessType:   "get_user_history",
		AccessResult: true,
	})
	
	return history, nil
}