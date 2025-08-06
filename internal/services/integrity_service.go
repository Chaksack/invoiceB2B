package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// IntegrityAlgorithm represents the hash algorithm used for integrity verification
type IntegrityAlgorithm string

const (
	// SHA256 algorithm
	SHA256 IntegrityAlgorithm = "sha256"
	// SHA512 algorithm
	SHA512 IntegrityAlgorithm = "sha512"
	// HMAC-SHA256 algorithm
	HMACSHA256 IntegrityAlgorithm = "hmac-sha256"
)

// IntegrityVerificationService defines the interface for file integrity verification
type IntegrityVerificationService interface {
	// CalculateChecksum calculates the checksum of a file
	CalculateChecksum(filePath string) (string, error)
	
	// VerifyChecksum verifies the checksum of a file
	VerifyChecksum(filePath, expectedChecksum string) (bool, error)
	
	// StoreChecksum stores the checksum of a file
	StoreChecksum(filePath, checksum string) error
	
	// GetStoredChecksum retrieves the stored checksum of a file
	GetStoredChecksum(filePath string) (string, error)
	
	// VerifyFileIntegrity verifies the integrity of a file using the stored checksum
	VerifyFileIntegrity(filePath string) (bool, error)
	
	// ScheduleIntegrityCheck schedules periodic integrity checks for a file
	ScheduleIntegrityCheck(filePath string, interval time.Duration) error
	
	// CancelIntegrityCheck cancels scheduled integrity checks for a file
	CancelIntegrityCheck(filePath string) error
}

// IntegrityConfig holds configuration for the integrity verification service
type IntegrityConfig struct {
	// Algorithm used for integrity verification
	Algorithm IntegrityAlgorithm
	
	// Secret key for HMAC (if using HMAC algorithm)
	HMACKey string
	
	// Directory for storing checksums
	ChecksumDir string
	
	// Whether to verify integrity on file access
	VerifyOnAccess bool
	
	// Whether to schedule periodic integrity checks
	ScheduleChecks bool
	
	// Interval for periodic integrity checks
	CheckInterval time.Duration
}

// integrityService implements IntegrityVerificationService
type integrityService struct {
	config IntegrityConfig
	
	// Map of scheduled integrity checks
	scheduledChecks     map[string]*time.Ticker
	scheduledChecksMutex sync.Mutex
}

// NewIntegrityVerificationService creates a new integrity verification service
func NewIntegrityVerificationService(config IntegrityConfig) (IntegrityVerificationService, error) {
	// Set default algorithm if not specified
	if config.Algorithm == "" {
		config.Algorithm = SHA256
	}
	
	// Set default check interval if not specified
	if config.CheckInterval == 0 {
		config.CheckInterval = 24 * time.Hour // Default to daily checks
	}
	
	// Set default checksum directory if not specified
	if config.ChecksumDir == "" {
		config.ChecksumDir = ".checksums"
	}
	
	// Create checksum directory if it doesn't exist
	if _, err := os.Stat(config.ChecksumDir); os.IsNotExist(err) {
		if err := os.MkdirAll(config.ChecksumDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create checksum directory: %w", err)
		}
	}
	
	return &integrityService{
		config:              config,
		scheduledChecks:     make(map[string]*time.Ticker),
		scheduledChecksMutex: sync.Mutex{},
	}, nil
}

// CalculateChecksum calculates the checksum of a file
func (s *integrityService) CalculateChecksum(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file for checksum calculation: %w", err)
	}
	defer file.Close()
	
	// Create the appropriate hash function
	var h hash.Hash
	switch s.config.Algorithm {
	case SHA256:
		h = sha256.New()
	case SHA512:
		h = sha512.New()
	case HMACSHA256:
		h = hmac.New(sha256.New, []byte(s.config.HMACKey))
	default:
		return "", fmt.Errorf("unsupported hash algorithm: %s", s.config.Algorithm)
	}
	
	// Calculate the hash
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("failed to calculate file hash: %w", err)
	}
	
	// Return the hex-encoded hash
	return hex.EncodeToString(h.Sum(nil)), nil
}

// VerifyChecksum verifies the checksum of a file
func (s *integrityService) VerifyChecksum(filePath, expectedChecksum string) (bool, error) {
	// Calculate the current checksum
	checksum, err := s.CalculateChecksum(filePath)
	if err != nil {
		return false, err
	}
	
	// Compare checksums
	return checksum == expectedChecksum, nil
}

// getChecksumFilePath returns the path to the checksum file for a given file
func (s *integrityService) getChecksumFilePath(filePath string) string {
	// Use the absolute path to avoid collisions
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		// Fall back to the original path if we can't get the absolute path
		absPath = filePath
	}
	
	// Replace path separators with underscores
	safePath := strings.ReplaceAll(absPath, string(os.PathSeparator), "_")
	
	// Create the checksum file path
	return filepath.Join(s.config.ChecksumDir, safePath+".checksum")
}

// StoreChecksum stores the checksum of a file
func (s *integrityService) StoreChecksum(filePath, checksum string) error {
	// Get the checksum file path
	checksumFilePath := s.getChecksumFilePath(filePath)
	
	// Create the checksum file
	if err := ioutil.WriteFile(checksumFilePath, []byte(checksum), 0600); err != nil {
		return fmt.Errorf("failed to write checksum file: %w", err)
	}
	
	return nil
}

// GetStoredChecksum retrieves the stored checksum of a file
func (s *integrityService) GetStoredChecksum(filePath string) (string, error) {
	// Get the checksum file path
	checksumFilePath := s.getChecksumFilePath(filePath)
	
	// Read the checksum file
	data, err := ioutil.ReadFile(checksumFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read checksum file: %w", err)
	}
	
	return string(data), nil
}

// VerifyFileIntegrity verifies the integrity of a file using the stored checksum
func (s *integrityService) VerifyFileIntegrity(filePath string) (bool, error) {
	// Get the stored checksum
	storedChecksum, err := s.GetStoredChecksum(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to get stored checksum: %w", err)
	}
	
	// Verify the checksum
	return s.VerifyChecksum(filePath, storedChecksum)
}

// ScheduleIntegrityCheck schedules periodic integrity checks for a file
func (s *integrityService) ScheduleIntegrityCheck(filePath string, interval time.Duration) error {
	s.scheduledChecksMutex.Lock()
	defer s.scheduledChecksMutex.Unlock()
	
	// Cancel any existing scheduled check for this file
	if ticker, exists := s.scheduledChecks[filePath]; exists {
		ticker.Stop()
	}
	
	// Create a new ticker
	ticker := time.NewTicker(interval)
	s.scheduledChecks[filePath] = ticker
	
	// Start the integrity check goroutine
	go func() {
		for range ticker.C {
			// Verify the file integrity
			ok, err := s.VerifyFileIntegrity(filePath)
			if err != nil {
				// Log the error
				fmt.Printf("Error verifying integrity of %s: %v\n", filePath, err)
				continue
			}
			
			if !ok {
				// Log the integrity violation
				fmt.Printf("Integrity violation detected for %s\n", filePath)
				
				// TODO: Implement additional actions for integrity violations
				// e.g., send alerts, restore from backup, etc.
			}
		}
	}()
	
	return nil
}

// CancelIntegrityCheck cancels scheduled integrity checks for a file
func (s *integrityService) CancelIntegrityCheck(filePath string) error {
	s.scheduledChecksMutex.Lock()
	defer s.scheduledChecksMutex.Unlock()
	
	// Cancel the scheduled check
	if ticker, exists := s.scheduledChecks[filePath]; exists {
		ticker.Stop()
		delete(s.scheduledChecks, filePath)
	}
	
	return nil
}

// CalculateAndStoreChecksum calculates and stores the checksum of a file
func (s *integrityService) CalculateAndStoreChecksum(filePath string) error {
	// Calculate the checksum
	checksum, err := s.CalculateChecksum(filePath)
	if err != nil {
		return err
	}
	
	// Store the checksum
	return s.StoreChecksum(filePath, checksum)
}