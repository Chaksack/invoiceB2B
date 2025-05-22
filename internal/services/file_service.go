package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileService interface {
	SaveFile(file *multipart.FileHeader, subDir string) (string, string, error) // returns filePath, fileName, error
	ValidateFileSize(size int64) error
	GetAbsPath(relativePath string) (string, error)
}

type fileService struct {
	baseUploadsDir string
	maxSizeBytes   int64
}

func NewFileService(baseUploadsDir string, maxSizeBytes int64) FileService {
	return &fileService{
		baseUploadsDir: baseUploadsDir,
		maxSizeBytes:   maxSizeBytes,
	}
}

func (s *fileService) SaveFile(file *multipart.FileHeader, subDir string) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate unique filename to prevent overwrites and sanitize
	originalFilename := filepath.Base(file.Filename)
	ext := filepath.Ext(originalFilename)
	// Sanitize filename part (simple sanitization, more robust might be needed)
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

	uniqueFilename := fmt.Sprintf("%d-%s-%s%s", time.Now().UnixNano(), uuid.NewString()[:8], sanitizedNamePart, ext)

	// Ensure subdirectory exists
	fullSubDirPath := filepath.Join(s.baseUploadsDir, subDir)
	if _, err := os.Stat(fullSubDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullSubDirPath, os.ModePerm); err != nil {
			return "", "", fmt.Errorf("failed to create upload subdirectory %s: %w", subDir, err)
		}
	}

	filePath := filepath.Join(fullSubDirPath, uniqueFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", "", fmt.Errorf("failed to copy uploaded file content: %w", err)
	}

	// Return relative path for storage in DB, and the generated filename for display/download
	relativePath := filepath.Join(subDir, uniqueFilename) // Path relative to baseUploadsDir
	return relativePath, uniqueFilename, nil
}

func (s *fileService) ValidateFileSize(size int64) error {
	if size > s.maxSizeBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed %d bytes", size, s.maxSizeBytes)
	}
	return nil
}

func (s *fileService) GetAbsPath(relativePath string) (string, error) {
	absPath := filepath.Join(s.baseUploadsDir, relativePath)
	// Optional: Check if file exists, etc.
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found at path: %s", absPath)
	}
	return absPath, nil
}
