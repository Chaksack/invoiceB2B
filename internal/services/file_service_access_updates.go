package services

import (
	"fmt"
	"mime/multipart"
	"time"
)

// SaveFileWithAccess is an enhanced version of SaveFile that includes access control
func (s *fileService) SaveFileWithAccess(file *multipart.FileHeader, subDir string, userID string, userRole string) (string, string, error) {
	// Check if the user has write access to the directory
	hasAccess, err := s.CheckAccess(subDir, userID, userRole, WriteAccess)
	if err != nil {
		return "", "", fmt.Errorf("failed to check access for file upload: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     subDir,
			AccessType:   "upload",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return "", "", fmt.Errorf("access denied: insufficient permissions to upload file to this directory")
	}
	
	// Call the original SaveFile method
	relativePath, fileName, err := s.SaveFile(file, subDir)
	if err != nil {
		// Log the failure
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     subDir,
			AccessType:   "upload",
			AccessResult: false,
			Reason:       fmt.Sprintf("file upload failed: %v", err),
		})
		
		return "", "", err
	}
	
	// Log the successful upload
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "upload",
		AccessResult: true,
	})
	
	return relativePath, fileName, nil
}

// GetAbsPathWithAccess is an enhanced version of GetAbsPath that includes access control
func (s *fileService) GetAbsPathWithAccess(relativePath string, userID string, userRole string) (string, error) {
	// Check if the user has read access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, ReadAccess)
	if err != nil {
		return "", fmt.Errorf("failed to check access for file path: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "read_path",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return "", fmt.Errorf("access denied: insufficient permissions to access this file")
	}
	
	// Call the original GetAbsPath method
	absPath, err := s.GetAbsPath(relativePath)
	if err != nil {
		return "", err
	}
	
	// Log the successful access
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "read_path",
		AccessResult: true,
	})
	
	return absPath, nil
}

// GenerateSignedURLWithAccess is an enhanced version of GenerateSignedURL that includes access control
func (s *fileService) GenerateSignedURLWithAccess(relativePath string, expiration time.Duration, userID string, userRole string) (string, error) {
	// Check if the user has read access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, ReadAccess)
	if err != nil {
		return "", fmt.Errorf("failed to check access for signed URL: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "generate_url",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return "", fmt.Errorf("access denied: insufficient permissions to generate signed URL for this file")
	}
	
	// Call the original GenerateSignedURL method
	signedURL, err := s.GenerateSignedURL(relativePath, expiration)
	if err != nil {
		return "", err
	}
	
	// Log the successful access
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "generate_url",
		AccessResult: true,
	})
	
	return signedURL, nil
}

// DecryptFileWithAccess is an enhanced version of DecryptFile that includes access control
func (s *fileService) DecryptFileWithAccess(encryptedFilePath string, userID string, userRole string) (string, error) {
	// Check if the user has read access to the file
	hasAccess, err := s.CheckAccess(encryptedFilePath, userID, userRole, ReadAccess)
	if err != nil {
		return "", fmt.Errorf("failed to check access for file decryption: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     encryptedFilePath,
			AccessType:   "decrypt",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return "", fmt.Errorf("access denied: insufficient permissions to decrypt this file")
	}
	
	// Call the original DecryptFile method
	decryptedPath, err := s.DecryptFile(encryptedFilePath)
	if err != nil {
		// Log the failure
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     encryptedFilePath,
			AccessType:   "decrypt",
			AccessResult: false,
			Reason:       fmt.Sprintf("file decryption failed: %v", err),
		})
		
		return "", err
	}
	
	// Log the successful decryption
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     encryptedFilePath,
		AccessType:   "decrypt",
		AccessResult: true,
	})
	
	return decryptedPath, nil
}

// VerifyFileIntegrityWithAccess is an enhanced version of VerifyFileIntegrity that includes access control
func (s *fileService) VerifyFileIntegrityWithAccess(relativePath string, userID string, userRole string) (bool, error) {
	// Check if the user has read access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, ReadAccess)
	if err != nil {
		return false, fmt.Errorf("failed to check access for integrity verification: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "verify_integrity",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return false, fmt.Errorf("access denied: insufficient permissions to verify integrity of this file")
	}
	
	// Call the original VerifyFileIntegrity method
	isValid, err := s.VerifyFileIntegrity(relativePath)
	if err != nil {
		// Log the failure
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "verify_integrity",
			AccessResult: false,
			Reason:       fmt.Sprintf("integrity verification failed: %v", err),
		})
		
		return false, err
	}
	
	// Log the successful verification
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "verify_integrity",
		AccessResult: true,
	})
	
	return isValid, nil
}

// CalculateChecksumWithAccess is an enhanced version of CalculateChecksum that includes access control
func (s *fileService) CalculateChecksumWithAccess(relativePath string, userID string, userRole string) (string, error) {
	// Check if the user has read access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, ReadAccess)
	if err != nil {
		return "", fmt.Errorf("failed to check access for checksum calculation: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "calculate_checksum",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return "", fmt.Errorf("access denied: insufficient permissions to calculate checksum for this file")
	}
	
	// Call the original CalculateChecksum method
	checksum, err := s.CalculateChecksum(relativePath)
	if err != nil {
		// Log the failure
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "calculate_checksum",
			AccessResult: false,
			Reason:       fmt.Sprintf("checksum calculation failed: %v", err),
		})
		
		return "", err
	}
	
	// Log the successful calculation
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "calculate_checksum",
		AccessResult: true,
	})
	
	return checksum, nil
}

// StoreChecksumWithAccess is an enhanced version of StoreChecksum that includes access control
func (s *fileService) StoreChecksumWithAccess(relativePath string, userID string, userRole string) error {
	// Check if the user has write access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, WriteAccess)
	if err != nil {
		return fmt.Errorf("failed to check access for checksum storage: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "store_checksum",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return fmt.Errorf("access denied: insufficient permissions to store checksum for this file")
	}
	
	// Call the original StoreChecksum method
	err = s.StoreChecksum(relativePath)
	if err != nil {
		// Log the failure
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "store_checksum",
			AccessResult: false,
			Reason:       fmt.Sprintf("checksum storage failed: %v", err),
		})
		
		return err
	}
	
	// Log the successful storage
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "store_checksum",
		AccessResult: true,
	})
	
	return nil
}

// ScheduleIntegrityCheckWithAccess is an enhanced version of ScheduleIntegrityCheck that includes access control
func (s *fileService) ScheduleIntegrityCheckWithAccess(relativePath string, interval time.Duration, userID string, userRole string) error {
	// Check if the user has write access to the file
	hasAccess, err := s.CheckAccess(relativePath, userID, userRole, WriteAccess)
	if err != nil {
		return fmt.Errorf("failed to check access for scheduling integrity check: %w", err)
	}
	
	if !hasAccess {
		// Log the access denial
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "schedule_integrity_check",
			AccessResult: false,
			Reason:       "insufficient permissions",
		})
		
		return fmt.Errorf("access denied: insufficient permissions to schedule integrity check for this file")
	}
	
	// Call the original ScheduleIntegrityCheck method
	err = s.ScheduleIntegrityCheck(relativePath, interval)
	if err != nil {
		// Log the failure
		s.accessService.LogAccess(FileAccessEvent{
			UserID:       userID,
			UserRole:     userRole,
			FilePath:     relativePath,
			AccessType:   "schedule_integrity_check",
			AccessResult: false,
			Reason:       fmt.Sprintf("scheduling integrity check failed: %v", err),
		})
		
		return err
	}
	
	// Log the successful scheduling
	s.accessService.LogAccess(FileAccessEvent{
		UserID:       userID,
		UserRole:     userRole,
		FilePath:     relativePath,
		AccessType:   "schedule_integrity_check",
		AccessResult: true,
	})
	
	return nil
}