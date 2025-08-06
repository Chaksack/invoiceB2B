package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	
	"github.com/google/uuid"
)

// Allowed file types for upload
var allowedMimeTypes = map[string]bool{
	// Images
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"image/svg+xml":   true,
	
	// Documents
	"application/pdf":                            true,
	"application/msword":                         true, // .doc
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
	"application/vnd.ms-excel":                   true, // .xls
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       true, // .xlsx
	"application/vnd.ms-powerpoint":              true, // .ppt
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
	
	// Text
	"text/csv":         true,
	"text/plain":       true,
	"text/xml":         true,
	"application/json": true,
	"application/xml":  true,
	
	// Archives (only if needed)
	"application/zip":  true,
	"application/x-7z-compressed": true,
}

// Allowed file extensions
var allowedExtensions = map[string]bool{
	// Images
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,
	
	// Documents
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	
	// Text
	".csv":  true,
	".txt":  true,
	".xml":  true,
	".json": true,
	
	// Archives (only if needed)
	".zip":  true,
	".7z":   true,
}

// High-risk file extensions that should be blocked
var highRiskExtensions = map[string]bool{
	".exe":  true,
	".dll":  true,
	".bat":  true,
	".cmd":  true,
	".sh":   true,
	".jar":  true,
	".js":   true,
	".vbs":  true,
	".ps1":  true,
	".py":   true,
	".php":  true,
	".asp":  true,
	".aspx": true,
	".jsp":  true,
	".cgi":  true,
	".com":  true,
	".scr":  true,
	".msi":  true,
}

type FileService interface {
	// Original methods (kept for backward compatibility)
	SaveFile(file *multipart.FileHeader, subDir string) (string, string, error) // returns filePath, fileName, error
	ValidateFileSize(size int64) error
	ValidateFileType(file *multipart.FileHeader) error
	GetAbsPath(relativePath string) (string, error)
	GenerateSignedURL(relativePath string, expiration time.Duration) (string, error)
	IsSensitiveFileType(fileType string) bool // checks if a file type is sensitive and should be encrypted
	DecryptFile(encryptedFilePath string) (string, error) // decrypts a file and returns the path to the decrypted file
	
	// Integrity verification methods (original)
	VerifyFileIntegrity(relativePath string) (bool, error) // verifies the integrity of a file
	CalculateChecksum(relativePath string) (string, error) // calculates the checksum of a file
	StoreChecksum(relativePath string) error // calculates and stores the checksum of a file
	ScheduleIntegrityCheck(relativePath string, interval time.Duration) error // schedules periodic integrity checks
	
	// Enhanced methods with access control
	SaveFileWithAccess(file *multipart.FileHeader, subDir string, userID string, userRole string) (string, string, error)
	GetAbsPathWithAccess(relativePath string, userID string, userRole string) (string, error)
	GenerateSignedURLWithAccess(relativePath string, expiration time.Duration, userID string, userRole string) (string, error)
	DecryptFileWithAccess(encryptedFilePath string, userID string, userRole string) (string, error)
	DeleteFile(relativePath string, userID string, userRole string) error
	
	// Enhanced integrity verification methods with access control
	VerifyFileIntegrityWithAccess(relativePath string, userID string, userRole string) (bool, error)
	CalculateChecksumWithAccess(relativePath string, userID string, userRole string) (string, error)
	StoreChecksumWithAccess(relativePath string, userID string, userRole string) error
	ScheduleIntegrityCheckWithAccess(relativePath string, interval time.Duration, userID string, userRole string) error
	
	// Access control methods
	CheckAccess(relativePath string, userID string, userRole string, accessLevel FileAccessLevel) (bool, error)
	GetAccessHistory(relativePath string, userID string, userRole string) ([]FileAccessEvent, error)
	GetUserAccessHistory(userID string, userRole string, requestingUserID string, requestingUserRole string) ([]FileAccessEvent, error)
}

type fileService struct {
	baseUploadsDir string
	maxSizeBytes   int64
	urlSecret      string // Secret for signing URLs
	virusScan      VirusScanService // Service for virus scanning
	encryption     EncryptionService // Service for file encryption
	encryptionEnabled bool // Whether encryption is enabled
	sensitiveTypes map[string]bool // Map of sensitive file types that should be encrypted
	integrity      IntegrityVerificationService // Service for file integrity verification
	integrityEnabled bool // Whether integrity verification is enabled
	verifyOnAccess bool // Whether to verify integrity on file access
	accessService  FileAccessService // Service for file access control and logging
}

func NewFileService(baseUploadsDir string, maxSizeBytes int64, virusScan VirusScanService, encryption EncryptionService, encryptionEnabled bool, sensitiveTypes map[string]bool, integrity IntegrityVerificationService, integrityEnabled bool, verifyOnAccess bool, accessService FileAccessService) FileService {
	// If no virus scan service is provided, create a default one with fallback scanning
	if virusScan == nil {
		// Use ScanMethod.Fallback (3) for both primary and fallback methods
		virusScan = NewVirusScanService(VirusScanConfig{
			PrimaryScanMethod:   3, // Fallback method
			FallbackScanMethod:  3, // Fallback method
			ScanTimeout:         30 * time.Second,
			LogResults:          true,
		})
	}
	
	// If no sensitive types map is provided, create an empty one
	if sensitiveTypes == nil {
		sensitiveTypes = make(map[string]bool)
	}
	
	// If no access service is provided, create a default one
	if accessService == nil {
		accessService = NewFileAccessService()
	}
	
	return &fileService{
		baseUploadsDir:    baseUploadsDir,
		maxSizeBytes:      maxSizeBytes,
		urlSecret:         uuid.NewString(), // Generate a random secret for URL signing
		virusScan:         virusScan,
		encryption:        encryption,
		encryptionEnabled: encryptionEnabled,
		sensitiveTypes:    sensitiveTypes,
		integrity:         integrity,
		integrityEnabled:  integrityEnabled,
		verifyOnAccess:    verifyOnAccess,
		accessService:     accessService,
	}
}

func (s *fileService) SaveFile(file *multipart.FileHeader, subDir string) (string, string, error) {
	// Validate file size
	if err := s.ValidateFileSize(file.Size); err != nil {
		return "", "", err
	}
	
	// Validate file type before processing
	if err := s.ValidateFileType(file); err != nil {
		return "", "", err
	}
	
	// Perform deep content inspection
	if err := s.ValidateFileContent(file); err != nil {
		return "", "", fmt.Errorf("file content validation failed: %w", err)
	}
	
	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate unique filename to prevent overwrites and sanitize
	originalFilename := filepath.Base(file.Filename)
	ext := filepath.Ext(originalFilename)
	
	// Validate file extension (double-check)
	if !allowedExtensions[strings.ToLower(ext)] {
		return "", "", fmt.Errorf("file extension %s not allowed", ext)
	}
	
	// Sanitize filename part (more robust sanitization)
	sanitizedNamePart := strings.ReplaceAll(strings.TrimSuffix(originalFilename, ext), " ", "_")
	sanitizedNamePart = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			return r
		}
		return -1 // Remove other characters
	}, sanitizedNamePart)
	if len(sanitizedNamePart) > 50 { // Truncate
		sanitizedNamePart = sanitizedNamePart[:50]
	}

	// Generate a truly unique filename with timestamp and UUID
	uniqueFilename := fmt.Sprintf("%d-%s-%s%s", time.Now().UnixNano(), uuid.NewString()[:8], sanitizedNamePart, ext)

	// Ensure subdirectory exists with secure permissions
	fullSubDirPath := filepath.Join(s.baseUploadsDir, subDir)
	if _, err := os.Stat(fullSubDirPath); os.IsNotExist(err) {
		// Create directory with restricted permissions (0750 = rwxr-x---)
		if err := os.MkdirAll(fullSubDirPath, 0750); err != nil {
			return "", "", fmt.Errorf("failed to create upload subdirectory %s: %w", subDir, err)
		}
	}

	filePath := filepath.Join(fullSubDirPath, uniqueFilename)

	// Create the destination file with restricted permissions (0640 = rw-r-----)
	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0640)
	if err != nil {
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the file content
	if _, err = io.Copy(dst, src); err != nil {
		// If copy fails, attempt to remove the partially written file
		os.Remove(filePath)
		return "", "", fmt.Errorf("failed to copy uploaded file content: %w", err)
	}

	// Sync to ensure data is written to disk
	if err = dst.Sync(); err != nil {
		os.Remove(filePath)
		return "", "", fmt.Errorf("failed to sync file to disk: %w", err)
	}

	// Calculate relative path for storage in DB
	relativePath := filepath.Join(subDir, uniqueFilename) // Path relative to baseUploadsDir
	
	// Check if the file should be encrypted
	shouldEncrypt := s.encryptionEnabled && s.encryption != nil && s.IsSensitiveFileType(subDir)
	
	// If encryption is enabled and this is a sensitive file type, encrypt it
	if shouldEncrypt {
		log.Printf("Encrypting sensitive file: %s (type: %s)", relativePath, subDir)
		
		// Close the file before encryption
		dst.Close()
		
		// Encrypt the file
		encryptedPath, err := s.encryption.EncryptFile(filePath, subDir)
		if err != nil {
			// If encryption fails, remove the original file
			os.Remove(filePath)
			return "", "", fmt.Errorf("failed to encrypt file: %w", err)
		}
		
		// Update the relative path to point to the encrypted file
		encryptedFilename := filepath.Base(encryptedPath)
		relativePath = filepath.Join(subDir, encryptedFilename)
		
		// Update the unique filename to include the encryption extension
		uniqueFilename = encryptedFilename
		
		// Update the file path to the encrypted file path
		filePath = encryptedPath
	}
	
	// Calculate and store file checksum if integrity verification is enabled
	if s.integrityEnabled && s.integrity != nil {
		// Calculate the checksum
		checksum, err := s.integrity.CalculateChecksum(filePath)
		if err != nil {
			log.Printf("Warning: Failed to calculate checksum for file %s: %v", relativePath, err)
		} else {
			// Store the checksum
			if err := s.integrity.StoreChecksum(filePath, checksum); err != nil {
				log.Printf("Warning: Failed to store checksum for file %s: %v", relativePath, err)
			} else {
				log.Printf("Stored integrity checksum for file: %s", relativePath)
			}
			
			// Schedule periodic integrity checks if configured
			if s.integrity != nil && s.verifyOnAccess {
				// Schedule a check every 24 hours (or as configured)
				if err := s.integrity.ScheduleIntegrityCheck(filePath, 24*time.Hour); err != nil {
					log.Printf("Warning: Failed to schedule integrity check for file %s: %v", relativePath, err)
				}
			}
		}
	}

	return relativePath, uniqueFilename, nil
}

// ValidateFileType performs comprehensive validation of file types
func (s *fileService) ValidateFileType(file *multipart.FileHeader) error {
	// Check file extension
	originalExt := filepath.Ext(file.Filename)
	ext := strings.ToLower(originalExt)
	
	// First, explicitly check for high-risk extensions
	if highRiskExtensions[ext] {
		return fmt.Errorf("high-risk file extension %s is not allowed", ext)
	}
	
	// Then check if the extension is in the allowed list
	if !allowedExtensions[ext] {
		return fmt.Errorf("file extension %s not allowed", ext)
	}
	
	// Open the file to check its content type
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for type validation: %w", err)
	}
	defer src.Close()
	
	// Read the first 4KB for content type detection and magic number validation
	// This is more than the standard 512 bytes to better detect file types
	buffer := make([]byte, 4096)
	bytesRead, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file content for type detection: %w", err)
	}
	
	// Trim buffer to actual bytes read
	buffer = buffer[:bytesRead]
	
	// Reset the file pointer to the beginning
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	// Detect content type
	contentType := http.DetectContentType(buffer)
	
	// Check if the content type is in the allowed list
	if !allowedMimeTypes[contentType] {
		return fmt.Errorf("detected content type %s not allowed", contentType)
	}
	
	// Verify that the content type matches the file extension
	if !s.validateContentTypeMatchesExtension(contentType, ext) {
		return fmt.Errorf("content type %s does not match file extension %s", contentType, ext)
	}
	
	// Perform additional validation based on file type
	switch {
	case strings.HasPrefix(contentType, "image/"):
		if err := s.validateImageFile(buffer, contentType, ext); err != nil {
			return err
		}
	case strings.HasPrefix(contentType, "application/pdf"):
		if err := s.validatePDFFile(buffer); err != nil {
			return err
		}
	case strings.HasPrefix(contentType, "application/msword") || 
	     strings.HasPrefix(contentType, "application/vnd.openxmlformats-officedocument"):
		if err := s.validateOfficeFile(buffer, ext); err != nil {
			return err
		}
	case strings.HasPrefix(contentType, "application/zip") || 
	     strings.HasPrefix(contentType, "application/x-7z-compressed"):
		if err := s.validateArchiveFile(buffer, ext); err != nil {
			return err
		}
	}
	
	// Check for embedded malicious content
	if s.detectMaliciousContent(buffer) {
		return fmt.Errorf("potential malicious content detected in file")
	}
	
	return nil
}

// validateContentTypeMatchesExtension checks if the detected content type matches the file extension
func (s *fileService) validateContentTypeMatchesExtension(contentType, ext string) bool {
	// Map of expected content types for each extension
	expectedContentTypes := map[string][]string{
		".jpg":  {"image/jpeg"},
		".jpeg": {"image/jpeg"},
		".png":  {"image/png"},
		".gif":  {"image/gif"},
		".webp": {"image/webp"},
		".svg":  {"image/svg+xml"},
		".pdf":  {"application/pdf"},
		".doc":  {"application/msword"},
		".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		".xls":  {"application/vnd.ms-excel"},
		".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		".ppt":  {"application/vnd.ms-powerpoint"},
		".pptx": {"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		".csv":  {"text/csv", "text/plain", "application/csv"},
		".txt":  {"text/plain"},
		".xml":  {"text/xml", "application/xml"},
		".json": {"application/json", "text/plain"},
		".zip":  {"application/zip"},
		".7z":   {"application/x-7z-compressed"},
	}
	
	// If we don't have expected content types for this extension, allow it
	expected, exists := expectedContentTypes[ext]
	if !exists {
		return true
	}
	
	// Check if the detected content type is in the list of expected content types
	for _, expectedType := range expected {
		if strings.HasPrefix(contentType, expectedType) {
			return true
		}
	}
	
	return false
}

// validateImageFile performs additional validation for image files
func (s *fileService) validateImageFile(buffer []byte, contentType, ext string) error {
	// Check for specific image types
	switch contentType {
	case "image/jpeg":
		// JPEG files start with FF D8
		if len(buffer) < 2 || buffer[0] != 0xFF || buffer[1] != 0xD8 {
			return fmt.Errorf("invalid JPEG file header")
		}
	case "image/png":
		// PNG files start with 89 50 4E 47 0D 0A 1A 0A
		pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		if len(buffer) < 8 {
			return fmt.Errorf("file too small to be a valid PNG")
		}
		for i := 0; i < 8; i++ {
			if buffer[i] != pngSignature[i] {
				return fmt.Errorf("invalid PNG file signature")
			}
		}
	case "image/gif":
		// GIF files start with GIF87a or GIF89a
		if len(buffer) < 6 {
			return fmt.Errorf("file too small to be a valid GIF")
		}
		if !bytes.Equal(buffer[0:3], []byte("GIF")) {
			return fmt.Errorf("invalid GIF file signature")
		}
		if !bytes.Equal(buffer[3:6], []byte("87a")) && !bytes.Equal(buffer[3:6], []byte("89a")) {
			return fmt.Errorf("invalid GIF version")
		}
	case "image/svg+xml":
		// SVG files should contain an <svg tag
		if !bytes.Contains(buffer, []byte("<svg")) {
			return fmt.Errorf("invalid SVG file content")
		}
		// Check for potentially malicious script tags
		if bytes.Contains(buffer, []byte("<script")) {
			return fmt.Errorf("SVG file contains potentially malicious script tags")
		}
	}
	
	return nil
}

// validatePDFFile performs additional validation for PDF files
func (s *fileService) validatePDFFile(buffer []byte) error {
	// PDF files start with %PDF-
	if len(buffer) < 5 || !bytes.Equal(buffer[0:5], []byte("%PDF-")) {
		return fmt.Errorf("invalid PDF file header")
	}
	
	// Check for potentially malicious JavaScript in PDF
	if bytes.Contains(buffer, []byte("/JS")) || bytes.Contains(buffer, []byte("/JavaScript")) {
		return fmt.Errorf("PDF file contains potentially malicious JavaScript")
	}
	
	// Check for potentially malicious actions in PDF
	if bytes.Contains(buffer, []byte("/Launch")) || bytes.Contains(buffer, []byte("/URI")) {
		// This is a simplistic check - in a real implementation, you'd want to parse the PDF
		// and analyze the actions more thoroughly
		log.Printf("Warning: PDF file contains potentially risky actions (Launch or URI)")
	}
	
	return nil
}

// validateOfficeFile performs additional validation for Microsoft Office files
func (s *fileService) validateOfficeFile(buffer []byte, ext string) error {
	// Office Open XML files (.docx, .xlsx, .pptx) are ZIP files
	if ext == ".docx" || ext == ".xlsx" || ext == ".pptx" {
		// Check for ZIP file signature
		if len(buffer) < 4 || !bytes.Equal(buffer[0:4], []byte("PK\x03\x04")) {
			return fmt.Errorf("invalid Office Open XML file signature")
		}
		
		// In a real implementation, you'd want to extract and scan the XML content
		// for potentially malicious macros, external links, etc.
	}
	
	// Legacy Office files (.doc, .xls, .ppt) use the Compound File Binary Format
	if ext == ".doc" || ext == ".xls" || ext == ".ppt" {
		// Check for CFBF signature (D0 CF 11 E0 A1 B1 1A E1)
		cfbfSignature := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
		if len(buffer) < 8 {
			return fmt.Errorf("file too small to be a valid Office file")
		}
		for i := 0; i < 8; i++ {
			if buffer[i] != cfbfSignature[i] {
				return fmt.Errorf("invalid Office file signature")
			}
		}
		
		// Check for potentially malicious macros
		if bytes.Contains(buffer, []byte("VBA")) || bytes.Contains(buffer, []byte("Macro")) {
			log.Printf("Warning: Office file may contain macros")
		}
	}
	
	return nil
}

// validateArchiveFile performs additional validation for archive files
func (s *fileService) validateArchiveFile(buffer []byte, ext string) error {
	// ZIP files start with PK\x03\x04
	if ext == ".zip" {
		if len(buffer) < 4 || !bytes.Equal(buffer[0:4], []byte("PK\x03\x04")) {
			return fmt.Errorf("invalid ZIP file signature")
		}
	}
	
	// 7Z files start with 7z\xBC\xAF\x27\x1C
	if ext == ".7z" {
		if len(buffer) < 6 || !bytes.Equal(buffer[0:6], []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}) {
			return fmt.Errorf("invalid 7Z file signature")
		}
	}
	
	// In a real implementation, you'd want to scan the archive contents
	// for potentially malicious files
	
	return nil
}

// ValidateFileContent performs deep inspection of file content based on file type
func (s *fileService) ValidateFileContent(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for content validation: %w", err)
	}
	defer src.Close()
	
	// Read the entire file content
	content, err := ioutil.ReadAll(src)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}
	
	// Reset file pointer
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}
	
	// Detect content type
	contentType := http.DetectContentType(content[:min(512, len(content))])
	
	// Validate based on content type
	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		return s.validateJPEG(content)
	case strings.HasPrefix(contentType, "image/png"):
		return s.validatePNG(content)
	case strings.HasPrefix(contentType, "application/pdf"):
		return s.validatePDF(content)
	case strings.HasPrefix(contentType, "text/csv"):
		return s.validateCSV(content)
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
}

// validateJPEG validates JPEG file structure
func (s *fileService) validateJPEG(content []byte) error {
	// Check for JPEG magic numbers
	if len(content) < 2 {
		return fmt.Errorf("file too small to be a valid JPEG")
	}
	
	// JPEG files start with FF D8
	if content[0] != 0xFF || content[1] != 0xD8 {
		return fmt.Errorf("invalid JPEG file header")
	}
	
	// JPEG files end with FF D9
	if len(content) < 4 || content[len(content)-2] != 0xFF || content[len(content)-1] != 0xD9 {
		return fmt.Errorf("invalid JPEG file footer")
	}
	
	// Check for malicious content
	if s.detectMaliciousContent(content) {
		return fmt.Errorf("potential malicious content detected in JPEG file")
	}
	
	return nil
}

// validatePNG validates PNG file structure
func (s *fileService) validatePNG(content []byte) error {
	// Check for PNG magic numbers
	if len(content) < 8 {
		return fmt.Errorf("file too small to be a valid PNG")
	}
	
	// PNG files start with 89 50 4E 47 0D 0A 1A 0A
	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i := 0; i < 8; i++ {
		if content[i] != pngSignature[i] {
			return fmt.Errorf("invalid PNG file signature")
		}
	}
	
	// Check for malicious content
	if s.detectMaliciousContent(content) {
		return fmt.Errorf("potential malicious content detected in PNG file")
	}
	
	return nil
}

// validatePDF validates PDF file structure
func (s *fileService) validatePDF(content []byte) error {
	// Check for PDF magic numbers
	if len(content) < 5 {
		return fmt.Errorf("file too small to be a valid PDF")
	}
	
	// PDF files start with %PDF-
	if !bytes.HasPrefix(content, []byte("%PDF-")) {
		return fmt.Errorf("invalid PDF file header")
	}
	
	// PDF files should end with %%EOF
	if !bytes.Contains(content[len(content)-1024:], []byte("%%EOF")) {
		return fmt.Errorf("invalid PDF file footer")
	}
	
	// Check for potentially dangerous PDF features
	dangerousPatterns := []string{
		"/JS", "/JavaScript", "/AA", "/OpenAction", "/Launch",
		"/RichMedia", "/SubmitForm", "/ImportData",
	}
	
	for _, pattern := range dangerousPatterns {
		if bytes.Contains(content, []byte(pattern)) {
			return fmt.Errorf("potentially dangerous PDF feature detected: %s", pattern)
		}
	}
	
	// Check for malicious content
	if s.detectMaliciousContent(content) {
		return fmt.Errorf("potential malicious content detected in PDF file")
	}
	
	return nil
}

// validateCSV validates CSV file structure
func (s *fileService) validateCSV(content []byte) error {
	// Basic CSV validation - check for proper line endings and commas
	lines := bytes.Split(content, []byte("\n"))
	
	if len(lines) < 1 {
		return fmt.Errorf("empty CSV file")
	}
	
	// Check for formula injection attacks in CSV
	formulaPatterns := []string{
		"=cmd", "=powershell", "=shell", "=exec",
		"@cmd", "@powershell", "@shell", "@exec",
		"+cmd", "+powershell", "+shell", "+exec",
		"-cmd", "-powershell", "-shell", "-exec",
		"|cmd", "|powershell", "|shell", "|exec",
	}
	
	for _, line := range lines {
		for _, pattern := range formulaPatterns {
			if bytes.Contains(bytes.ToLower(line), []byte(pattern)) {
				return fmt.Errorf("potential formula injection detected in CSV file")
			}
		}
	}
	
	// Check for malicious content
	if s.detectMaliciousContent(content) {
		return fmt.Errorf("potential malicious content detected in CSV file")
	}
	
	return nil
}

// detectMaliciousContent performs comprehensive detection of malicious content in files
// Uses the VirusScanService if available, otherwise falls back to pattern matching
func (s *fileService) detectMaliciousContent(content []byte) bool {
	// Use VirusScanService if available
	if s.virusScan != nil && s.virusScan.IsAvailable() {
		// Scan the content using the virus scan service
		result, err := s.virusScan.ScanBytes(content)
		if err == nil {
			// If scan was successful, return whether the file is clean or not
			return !result.Clean
		}
		// If scan failed, log the error and fall back to pattern matching
		log.Printf("Virus scan failed: %v. Falling back to pattern matching.", err)
	}
	
	// Fall back to pattern matching if virus scan service is not available or failed
	
	// Convert content to string for string-based pattern matching
	contentStr := string(content)
	
	// 1. Check for known virus test signatures
	knownSignatures := [][]byte{
		// EICAR test virus signature
		[]byte("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*"),
		// Other known malware signatures could be added here
	}
	
	for _, signature := range knownSignatures {
		if bytes.Contains(content, signature) {
			log.Printf("Detected known malware signature")
			return true
		}
	}
	
	// 2. Check for executable code and shell commands
	executablePatterns := []string{
		// Shellcode markers
		"\\x90\\x90\\x90\\x90", // NOP sled
		"\\x31\\xc0\\x50\\x68", // Common shellcode start
		"\\x48\\x31\\xc0",      // Common x64 shellcode start
		"\\xeb\\x3e\\x5b\\x31", // Common JMP/CALL/POP technique
		
		// Command execution
		"cmd.exe", "powershell.exe", "bash", "/bin/sh", "/bin/bash",
		"cmd /c", "cmd /k", "cmd.exe /c", "powershell -e", "powershell -enc",
		
		// Function calls that can execute code
		"eval(", "setTimeout(", "setInterval(", "Function(", "constructor",
		"system(", "exec(", "shell_exec(", "passthru(", "proc_open(",
		"popen(", "assert(", "create_function(", "call_user_func",
		"include(", "include_once(", "require(", "require_once(",
		
		// Reverse shells
		"socket(", "connect(", "bind(", "listen(", "accept(",
		"fsockopen(", "pfsockopen(",
		
		// Base64 encoded PHP tags (common in web shells)
		"PD9waHA", "PHBocA", "ZXZhbCg", "c3lzdGVt", "ZXhlYyg", "cGFzc3RocnU",
		
		// Common web shell names and patterns
		"c99.php", "r57.php", "shell.php", "wso.php", "b374k",
		"weevely", "china chopper", "webshell", "backdoor",
	}
	
	for _, pattern := range executablePatterns {
		if strings.Contains(contentStr, pattern) {
			log.Printf("Detected executable code pattern: %s", pattern)
			return true
		}
	}
	
	// 3. Check for common exploit patterns using regex
	exploitPatterns := []*regexp.Regexp{
		// Script injection
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		regexp.MustCompile(`(?i)<script[^>]*src\s*=`),
		regexp.MustCompile(`(?i)javascript:`),
		
		// Event handlers
		regexp.MustCompile(`(?i)\s+on\w+\s*=`), // Matches onload=, onerror=, etc.
		
		// iframes
		regexp.MustCompile(`(?i)<iframe[^>]*src\s*=`),
		
		// Data URIs
		regexp.MustCompile(`(?i)data:text/html`),
		regexp.MustCompile(`(?i)data:application/javascript`),
		
		// SQL Injection patterns
		regexp.MustCompile(`(?i)'\s*OR\s*'.*?'\s*=\s*'`),
		regexp.MustCompile(`(?i)'\s*OR\s*1\s*=\s*1`),
		regexp.MustCompile(`(?i)'\s*;\s*DROP\s+TABLE`),
		
		// XML external entity (XXE) attacks
		regexp.MustCompile(`(?i)<!ENTITY\s+\w+\s+SYSTEM`),
		
		// Command injection
		regexp.MustCompile(`(?i)[;&|]\s*(?:ls|dir|cat|type|more|wget|curl)`),
		
		// PHP code injection
		regexp.MustCompile(`(?i)<\?php`),
		
		// Obfuscated JavaScript
		regexp.MustCompile(`(?i)eval\s*\(\s*(?:atob|base64_decode|String\.fromCharCode)`),
		regexp.MustCompile(`(?i)(?:unescape|decodeURIComponent|fromCharCode)\s*\(`),
		
		// Obfuscated PowerShell
		regexp.MustCompile(`(?i)-(?:enc|encodedcommand)\s+[A-Za-z0-9+/=]{20,}`),
	}
	
	for _, pattern := range exploitPatterns {
		if pattern.Match(content) {
			log.Printf("Detected exploit pattern: %s", pattern.String())
			return true
		}
	}
	
	// 4. File type-specific checks
	
	// PDF specific checks
	if bytes.HasPrefix(content, []byte("%PDF-")) {
		pdfExploitPatterns := []string{
			"/JS ", "/JavaScript ", // JavaScript in PDF
			"/Launch", // Can launch external applications
			"/OpenAction", // Automatically executed when PDF is opened
			"/AA ", // Additional actions
			"/JBIG2Decode", // Associated with exploits
			"/RichMedia", // Can embed Flash
			"/GoTo", // Can redirect to other documents
			"/GoToR", // Can redirect to remote documents
			"/GoToE", // Can embed files
			"/SubmitForm", // Can submit forms to external URLs
			"/XFA", // XML Forms Architecture, can contain JavaScript
		}
		
		for _, pattern := range pdfExploitPatterns {
			if strings.Contains(contentStr, pattern) {
				log.Printf("Detected PDF exploit pattern: %s", pattern)
				return true
			}
		}
	}
	
	// Office document specific checks
	if bytes.HasPrefix(content, []byte("PK")) || // Office Open XML
	   bytes.HasPrefix(content, []byte{0xD0, 0xCF, 0x11, 0xE0}) { // Compound File Binary Format
		officeExploitPatterns := []string{
			"VBA", "vbaProject", // VBA macros
			"\\w:autoOpen", "\\w:AutoExec", "Auto_Open", "AutoOpen", // Auto-executing macros
			"ActiveX", "CLSID", // ActiveX controls
			"powershell", "cmd.exe", "cscript", "wscript", // Shell commands
			"http://", "https://", "ftp://", // External links
			"DDEAUTO", "DDEEXEC", // DDE exploits
			"INCLUDEPICTURE", "INCLUDETEXT", // Can include external content
			"IMPORTDATA", // Can import external data
		}
		
		for _, pattern := range officeExploitPatterns {
			if strings.Contains(contentStr, pattern) {
				log.Printf("Detected Office document exploit pattern: %s", pattern)
				return true
			}
		}
	}
	
	// 5. Check for entropy (potential encryption/obfuscation)
	// This is a simple implementation - a more sophisticated one would use Shannon entropy
	if len(content) > 256 {
		// Take a sample of the content
		sample := content[:256]
		uniqueChars := make(map[byte]bool)
		for _, b := range sample {
			uniqueChars[b] = true
		}
		
		// If more than 80% of possible byte values are used, it might be encrypted/compressed/obfuscated
		if len(uniqueChars) > 204 { // 80% of 256
			// This is just a heuristic - would need further analysis
			log.Printf("Warning: High entropy detected, possible encryption or obfuscation")
			// We don't return true here as it's just a warning, not a definitive detection
		}
	}
	
	// 6. Check for polyglot files (files that are valid as multiple types)
	// For example, a file that is both a valid image and contains executable code
	if len(content) > 4 {
		// Check if file starts with image signature but contains code
		isImage := bytes.HasPrefix(content, []byte{0xFF, 0xD8, 0xFF}) || // JPEG
		           bytes.HasPrefix(content, []byte{0x89, 0x50, 0x4E, 0x47}) || // PNG
		           bytes.HasPrefix(content, []byte("GIF8")) // GIF
		
		if isImage {
			// Check for code after image data
			codePatterns := []string{"<?php", "<script", "eval(", "function(", "exec(", "shell_exec("}
			for _, pattern := range codePatterns {
				if strings.Contains(contentStr, pattern) {
					log.Printf("Detected code in image file: %s", pattern)
					return true
				}
			}
		}
	}
	
	return false
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// scanForViruses scans a file for viruses using the VirusScanService
// This replaces the old simulateVirusScan method with real virus scanning
func (s *fileService) scanForViruses(filePath string) (bool, []string, error) {
	// Use VirusScanService if available
	if s.virusScan != nil && s.virusScan.IsAvailable() {
		// Scan the file using the virus scan service
		result, err := s.virusScan.ScanFile(filePath)
		if err != nil {
			return false, nil, fmt.Errorf("virus scan failed: %w", err)
		}
		
		// Return the scan result
		return !result.Clean, result.Threats, nil
	}
	
	// If virus scan service is not available, read the file and use pattern matching
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read file for virus scan: %w", err)
	}
	
	// Use pattern matching as fallback
	isMalicious := s.detectMaliciousContent(content)
	
	if isMalicious {
		return true, []string{"Suspicious content detected by pattern matching"}, nil
	}
	
	return false, nil, nil
}

// simulateVirusScan is kept for backward compatibility
// It now uses the real virus scanning service if available
func (s *fileService) simulateVirusScan(fileContent []byte) bool {
	// Use VirusScanService if available
	if s.virusScan != nil && s.virusScan.IsAvailable() {
		result, err := s.virusScan.ScanBytes(fileContent)
		if err == nil {
			return !result.Clean
		}
	}
	
	// Fall back to pattern matching
	return s.detectMaliciousContent(fileContent)
}

func (s *fileService) ValidateFileSize(size int64) error {
	if size > s.maxSizeBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed %d bytes", size, s.maxSizeBytes)
	}
	return nil
}

func (s *fileService) GetAbsPath(relativePath string) (string, error) {
	absPath := filepath.Join(s.baseUploadsDir, relativePath)
	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found at path: %s", absPath)
	}
	return absPath, nil
}

// GenerateSignedURL creates a signed URL with a short expiration time
func (s *fileService) GenerateSignedURL(relativePath string, expiration time.Duration) (string, error) {
	// Check if the file exists
	_, err := s.GetAbsPath(relativePath)
	if err != nil {
		return "", err
	}
	
	// Calculate expiration timestamp
	expiresAt := time.Now().Add(expiration).Unix()
	
	// Create the signature
	h := hmac.New(sha256.New, []byte(s.urlSecret))
	h.Write([]byte(relativePath))
	h.Write([]byte(fmt.Sprintf("%d", expiresAt)))
	signature := hex.EncodeToString(h.Sum(nil))
	
	// Construct the signed URL
	// In a real implementation, this would be a proper URL with host, etc.
	signedURL := fmt.Sprintf("/api/v1/files/download?path=%s&expires=%d&signature=%s", 
		relativePath, expiresAt, signature)
	
	return signedURL, nil
}

// IsSensitiveFileType checks if a file type is sensitive and should be encrypted
func (s *fileService) IsSensitiveFileType(fileType string) bool {
	// If encryption is not enabled, no file types are considered sensitive
	if !s.encryptionEnabled || s.encryption == nil {
		return false
	}
	
	// Check if the file type is in the sensitive types map
	return s.sensitiveTypes[fileType]
}

// DecryptFile decrypts a file and returns the path to the decrypted file
func (s *fileService) DecryptFile(encryptedFilePath string) (string, error) {
	// If encryption is not enabled or encryption service is not available, return an error
	if !s.encryptionEnabled || s.encryption == nil {
		return "", fmt.Errorf("encryption is not enabled or encryption service is not available")
	}
	
	// Check if the file exists
	absPath, err := s.GetAbsPath(encryptedFilePath)
	if err != nil {
		return "", err
	}
	
	// Check if the file is encrypted
	if !s.encryption.IsEncrypted(absPath) {
		return "", fmt.Errorf("file is not encrypted: %s", encryptedFilePath)
	}
	
	// Decrypt the file
	decryptedPath, err := s.encryption.DecryptFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt file: %w", err)
	}
	
	return decryptedPath, nil
}

// VerifyFileIntegrity verifies the integrity of a file
func (s *fileService) VerifyFileIntegrity(relativePath string) (bool, error) {
	// If integrity verification is not enabled or integrity service is not available, return an error
	if !s.integrityEnabled || s.integrity == nil {
		return false, fmt.Errorf("integrity verification is not enabled or integrity service is not available")
	}
	
	// Get the absolute path
	absPath, err := s.GetAbsPath(relativePath)
	if err != nil {
		return false, err
	}
	
	// Verify the file integrity
	return s.integrity.VerifyFileIntegrity(absPath)
}

// CalculateChecksum calculates the checksum of a file
func (s *fileService) CalculateChecksum(relativePath string) (string, error) {
	// If integrity verification is not enabled or integrity service is not available, return an error
	if !s.integrityEnabled || s.integrity == nil {
		return "", fmt.Errorf("integrity verification is not enabled or integrity service is not available")
	}
	
	// Get the absolute path
	absPath, err := s.GetAbsPath(relativePath)
	if err != nil {
		return "", err
	}
	
	// Calculate the checksum
	return s.integrity.CalculateChecksum(absPath)
}

// StoreChecksum calculates and stores the checksum of a file
func (s *fileService) StoreChecksum(relativePath string) error {
	// If integrity verification is not enabled or integrity service is not available, return an error
	if !s.integrityEnabled || s.integrity == nil {
		return fmt.Errorf("integrity verification is not enabled or integrity service is not available")
	}
	
	// Get the absolute path
	absPath, err := s.GetAbsPath(relativePath)
	if err != nil {
		return err
	}
	
	// Calculate the checksum
	checksum, err := s.integrity.CalculateChecksum(absPath)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}
	
	// Store the checksum
	return s.integrity.StoreChecksum(absPath, checksum)
}

// ScheduleIntegrityCheck schedules periodic integrity checks for a file
func (s *fileService) ScheduleIntegrityCheck(relativePath string, interval time.Duration) error {
	// If integrity verification is not enabled or integrity service is not available, return an error
	if !s.integrityEnabled || s.integrity == nil {
		return fmt.Errorf("integrity verification is not enabled or integrity service is not available")
	}
	
	// Get the absolute path
	absPath, err := s.GetAbsPath(relativePath)
	if err != nil {
		return err
	}
	
	// Schedule the integrity check
	return s.integrity.ScheduleIntegrityCheck(absPath, interval)
}
