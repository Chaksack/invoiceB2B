package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// VirusScanResult represents the result of a virus scan
type VirusScanResult struct {
	Clean       bool     // Whether the file is clean (no viruses detected)
	Threats     []string // List of detected threats
	ScanID      string   // Unique ID for the scan
	ScanTime    time.Time // When the scan was performed
	Description string    // Additional information about the scan
}

// VirusScanService defines the interface for virus scanning
type VirusScanService interface {
	// ScanFile scans a file for viruses
	ScanFile(filePath string) (*VirusScanResult, error)
	
	// ScanBytes scans a byte array for viruses
	ScanBytes(data []byte) (*VirusScanResult, error)
	
	// IsAvailable checks if the virus scanning service is available
	IsAvailable() bool
}

// ScanMethod represents different virus scanning methods
type ScanMethod int

const (
	// LocalClamAV uses a local ClamAV installation
	LocalClamAV ScanMethod = iota
	
	// ClamAVDaemon uses clamd daemon
	ClamAVDaemon
	
	// ExternalAPI uses an external virus scanning API
	ExternalAPI
	
	// Fallback uses a basic pattern-based scan as fallback
	Fallback
)

// VirusScanConfig holds configuration for the virus scanning service
type VirusScanConfig struct {
	// Primary scan method to use
	PrimaryScanMethod ScanMethod
	
	// Fallback scan method if primary fails
	FallbackScanMethod ScanMethod
	
	// Path to clamav binary (for LocalClamAV)
	ClamAVPath string
	
	// ClamAV daemon socket path or address (for ClamAVDaemon)
	ClamAVDaemonAddress string
	
	// External API URL (for ExternalAPI)
	ExternalAPIURL string
	
	// External API key (for ExternalAPI)
	ExternalAPIKey string
	
	// Timeout for scan operations
	ScanTimeout time.Duration
	
	// Whether to log scan results
	LogResults bool
}

// virusScanService implements VirusScanService
type virusScanService struct {
	config VirusScanConfig
}

// NewVirusScanService creates a new virus scanning service
func NewVirusScanService(config VirusScanConfig) VirusScanService {
	// Set default timeout if not specified
	if config.ScanTimeout == 0 {
		config.ScanTimeout = 30 * time.Second
	}
	
	return &virusScanService{
		config: config,
	}
}

// IsAvailable checks if the virus scanning service is available
func (s *virusScanService) IsAvailable() bool {
	switch s.config.PrimaryScanMethod {
	case LocalClamAV:
		return s.isLocalClamAVAvailable()
	case ClamAVDaemon:
		return s.isClamAVDaemonAvailable()
	case ExternalAPI:
		return s.isExternalAPIAvailable()
	case Fallback:
		return true // Fallback is always available
	default:
		return false
	}
}

// isLocalClamAVAvailable checks if local ClamAV is available
func (s *virusScanService) isLocalClamAVAvailable() bool {
	clamPath := s.config.ClamAVPath
	if clamPath == "" {
		clamPath = "clamscan" // Default command name
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, clamPath, "--version")
	err := cmd.Run()
	
	return err == nil
}

// isClamAVDaemonAvailable checks if ClamAV daemon is available
func (s *virusScanService) isClamAVDaemonAvailable() bool {
	// This is a simplified check - in a real implementation,
	// you would attempt to connect to the daemon socket
	return s.config.ClamAVDaemonAddress != ""
}

// isExternalAPIAvailable checks if the external API is available
func (s *virusScanService) isExternalAPIAvailable() bool {
	if s.config.ExternalAPIURL == "" || s.config.ExternalAPIKey == "" {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", s.config.ExternalAPIURL+"/status", nil)
	if err != nil {
		return false
	}
	
	req.Header.Set("X-API-Key", s.config.ExternalAPIKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

// ScanFile scans a file for viruses
func (s *virusScanService) ScanFile(filePath string) (*VirusScanResult, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}
	
	// Try primary scan method
	result, err := s.scanFileWithMethod(filePath, s.config.PrimaryScanMethod)
	if err == nil {
		return result, nil
	}
	
	// Log the error from primary scan
	log.Printf("Primary virus scan method failed: %v. Trying fallback method.", err)
	
	// Try fallback scan method
	if s.config.FallbackScanMethod != s.config.PrimaryScanMethod {
		result, fallbackErr := s.scanFileWithMethod(filePath, s.config.FallbackScanMethod)
		if fallbackErr == nil {
			return result, nil
		}
		
		// Both methods failed
		return nil, fmt.Errorf("all virus scan methods failed: primary: %v, fallback: %v", err, fallbackErr)
	}
	
	// No fallback or fallback is the same as primary
	return nil, err
}

// scanFileWithMethod scans a file using the specified method
func (s *virusScanService) scanFileWithMethod(filePath string, method ScanMethod) (*VirusScanResult, error) {
	switch method {
	case LocalClamAV:
		return s.scanWithLocalClamAV(filePath)
	case ClamAVDaemon:
		return s.scanWithClamAVDaemon(filePath)
	case ExternalAPI:
		return s.scanWithExternalAPI(filePath)
	case Fallback:
		// Read file content for fallback scan
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file for fallback scan: %w", err)
		}
		return s.fallbackScan(data)
	default:
		return nil, fmt.Errorf("unsupported scan method: %d", method)
	}
}

// scanWithLocalClamAV scans a file using local ClamAV installation
func (s *virusScanService) scanWithLocalClamAV(filePath string) (*VirusScanResult, error) {
	clamPath := s.config.ClamAVPath
	if clamPath == "" {
		clamPath = "clamscan" // Default command name
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ScanTimeout)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, clamPath, "--no-summary", filePath)
	output, err := cmd.CombinedOutput()
	
	result := &VirusScanResult{
		Clean:    true,
		ScanTime: time.Now(),
		ScanID:   fmt.Sprintf("clam-%d", time.Now().UnixNano()),
	}
	
	if err != nil {
		// ClamAV returns exit code 1 if virus is found
		if strings.Contains(string(output), "FOUND") {
			result.Clean = false
			
			// Extract threat names
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "FOUND") {
					parts := strings.Split(line, "FOUND")
					if len(parts) > 0 {
						threat := strings.TrimSpace(parts[len(parts)-1])
						result.Threats = append(result.Threats, threat)
					}
				}
			}
			
			result.Description = fmt.Sprintf("ClamAV found %d threats", len(result.Threats))
			return result, nil
		}
		
		// Other error
		return nil, fmt.Errorf("ClamAV scan failed: %w, output: %s", err, string(output))
	}
	
	result.Description = "ClamAV scan completed, no threats found"
	return result, nil
}

// scanWithClamAVDaemon scans a file using ClamAV daemon
func (s *virusScanService) scanWithClamAVDaemon(filePath string) (*VirusScanResult, error) {
	// This is a placeholder for ClamAV daemon integration
	// In a real implementation, you would connect to the daemon socket
	// and send the SCAN command with the file path
	
	return nil, fmt.Errorf("ClamAV daemon scanning not implemented")
}

// scanWithExternalAPI scans a file using an external API
func (s *virusScanService) scanWithExternalAPI(filePath string) (*VirusScanResult, error) {
	if s.config.ExternalAPIURL == "" || s.config.ExternalAPIKey == "" {
		return nil, fmt.Errorf("external API URL or key not configured")
	}
	
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for API scan: %w", err)
	}
	defer file.Close()
	
	// Create a pipe for streaming the file
	pr, pw := io.Pipe()
	defer pr.Close()
	
	// Create multipart writer
	go func() {
		defer pw.Close()
		
		// Copy file to pipe
		if _, err := io.Copy(pw, file); err != nil {
			log.Printf("Error copying file to pipe: %v", err)
		}
	}()
	
	// Create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ScanTimeout)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.ExternalAPIURL+"/scan", pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create API request: %w", err)
	}
	
	req.Header.Set("X-API-Key", s.config.ExternalAPIKey)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Filename", filepath.Base(filePath))
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API scan request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}
	
	// This is a simplified response parsing
	// In a real implementation, you would parse a structured JSON response
	result := &VirusScanResult{
		Clean:       !strings.Contains(string(body), "threat"),
		ScanTime:    time.Now(),
		ScanID:      fmt.Sprintf("api-%d", time.Now().UnixNano()),
		Description: fmt.Sprintf("API scan completed with status %d", resp.StatusCode),
	}
	
	if !result.Clean {
		result.Threats = []string{"API detected threat"}
	}
	
	return result, nil
}

// ScanBytes scans a byte array for viruses
func (s *virusScanService) ScanBytes(data []byte) (*VirusScanResult, error) {
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", "virus-scan-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up
	
	// Write data to the temporary file
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	
	// Close the file before scanning
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}
	
	// Scan the temporary file
	return s.ScanFile(tempFile.Name())
}

// fallbackScan implements a basic pattern-based scan as fallback
func (s *virusScanService) fallbackScan(data []byte) (*VirusScanResult, error) {
	result := &VirusScanResult{
		Clean:       true,
		ScanTime:    time.Now(),
		ScanID:      fmt.Sprintf("fallback-%d", time.Now().UnixNano()),
		Description: "Fallback scan completed",
	}
	
	// Check for EICAR test virus signature
	eicarPattern := []byte("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*")
	if bytes.Contains(data, eicarPattern) {
		result.Clean = false
		result.Threats = append(result.Threats, "EICAR-Test-File")
	}
	
	// Check for common malicious patterns
	maliciousPatterns := []struct {
		name    string
		pattern []byte
	}{
		{"Suspicious-Shell-Code", []byte("\\x90\\x90\\x90\\x90")}, // NOP sled
		{"Suspicious-Exec-Command", []byte("exec(")},
		{"Suspicious-System-Command", []byte("system(")},
		{"Suspicious-Shell-Exec", []byte("shell_exec(")},
		{"Suspicious-Eval", []byte("eval(")},
		{"Suspicious-Script-Tag", []byte("<script")},
	}
	
	for _, p := range maliciousPatterns {
		if bytes.Contains(data, p.pattern) {
			result.Clean = false
			result.Threats = append(result.Threats, p.name)
		}
	}
	
	if !result.Clean {
		result.Description = fmt.Sprintf("Fallback scan found %d potential threats", len(result.Threats))
	}
	
	return result, nil
}