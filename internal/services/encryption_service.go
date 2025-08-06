package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// EncryptionService defines the interface for file encryption and decryption
type EncryptionService interface {
	// EncryptFile encrypts a file and returns the path to the encrypted file
	EncryptFile(filePath string, sensitiveType string) (string, error)
	
	// DecryptFile decrypts a file and returns the path to the decrypted file
	DecryptFile(encryptedFilePath string) (string, error)
	
	// IsEncrypted checks if a file is encrypted
	IsEncrypted(filePath string) bool
	
	// GetEncryptionKey returns the encryption key for a specific sensitive type
	GetEncryptionKey(sensitiveType string) ([]byte, error)
}

// encryptionService implements EncryptionService
type encryptionService struct {
	// Master encryption key used to derive type-specific keys
	masterKey []byte
	
	// Map of sensitive file types to their encryption keys
	typeKeys map[string][]byte
	
	// Directory for temporary decrypted files
	tempDir string
	
	// File extension for encrypted files
	encryptedExt string
}

// EncryptionConfig holds configuration for the encryption service
type EncryptionConfig struct {
	// Master key for encryption (hex-encoded)
	MasterKey string
	
	// Directory for temporary decrypted files
	TempDir string
	
	// File extension for encrypted files
	EncryptedExt string
	
	// Map of sensitive file types to their encryption keys (hex-encoded)
	TypeKeys map[string]string
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(config EncryptionConfig) (EncryptionService, error) {
	// Decode master key from hex
	masterKey, err := hex.DecodeString(config.MasterKey)
	if err != nil {
		// If master key is not valid hex, derive it from the string
		masterKey = deriveKey([]byte(config.MasterKey), nil, 32)
	}
	
	// Initialize type keys map
	typeKeys := make(map[string][]byte)
	
	// Process type-specific keys
	for fileType, keyHex := range config.TypeKeys {
		// Decode key from hex
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			// If key is not valid hex, derive it from the master key and file type
			key = deriveKey(masterKey, []byte(fileType), 32)
		}
		typeKeys[fileType] = key
	}
	
	// Set default encrypted extension if not provided
	encryptedExt := config.EncryptedExt
	if encryptedExt == "" {
		encryptedExt = ".enc"
	}
	
	// Set default temp directory if not provided
	tempDir := config.TempDir
	if tempDir == "" {
		tempDir = os.TempDir()
	}
	
	// Create temp directory if it doesn't exist
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tempDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
	}
	
	return &encryptionService{
		masterKey:    masterKey,
		typeKeys:     typeKeys,
		tempDir:      tempDir,
		encryptedExt: encryptedExt,
	}, nil
}

// EncryptFile encrypts a file and returns the path to the encrypted file
func (s *encryptionService) EncryptFile(filePath string, sensitiveType string) (string, error) {
	// Get encryption key for the sensitive type
	key, err := s.GetEncryptionKey(sensitiveType)
	if err != nil {
		return "", err
	}
	
	// Read the file
	plaintext, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to create nonce: %w", err)
	}
	
	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	
	// Create encrypted file path
	encryptedFilePath := filePath + s.encryptedExt
	
	// Write encrypted data to file
	if err := ioutil.WriteFile(encryptedFilePath, ciphertext, 0600); err != nil {
		return "", fmt.Errorf("failed to write encrypted file: %w", err)
	}
	
	// Remove original file
	if err := os.Remove(filePath); err != nil {
		// If we can't remove the original file, remove the encrypted file to avoid duplicates
		os.Remove(encryptedFilePath)
		return "", fmt.Errorf("failed to remove original file: %w", err)
	}
	
	return encryptedFilePath, nil
}

// DecryptFile decrypts a file and returns the path to the decrypted file
func (s *encryptionService) DecryptFile(encryptedFilePath string) (string, error) {
	// Check if file is encrypted
	if !s.IsEncrypted(encryptedFilePath) {
		return "", fmt.Errorf("file is not encrypted: %s", encryptedFilePath)
	}
	
	// Read the encrypted file
	ciphertext, err := ioutil.ReadFile(encryptedFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read encrypted file: %w", err)
	}
	
	// Try each encryption key until one works
	var plaintext []byte
	var decryptionErr error
	
	// First try with type-specific keys
	for _, key := range s.typeKeys {
		plaintext, decryptionErr = s.decryptWithKey(ciphertext, key)
		if decryptionErr == nil {
			break
		}
	}
	
	// If type-specific keys didn't work, try with master key
	if decryptionErr != nil {
		plaintext, decryptionErr = s.decryptWithKey(ciphertext, s.masterKey)
		if decryptionErr != nil {
			return "", fmt.Errorf("failed to decrypt file: %w", decryptionErr)
		}
	}
	
	// Create decrypted file path (in temp directory)
	baseFileName := filepath.Base(encryptedFilePath)
	baseFileName = strings.TrimSuffix(baseFileName, s.encryptedExt)
	decryptedFilePath := filepath.Join(s.tempDir, baseFileName)
	
	// Write decrypted data to file
	if err := ioutil.WriteFile(decryptedFilePath, plaintext, 0600); err != nil {
		return "", fmt.Errorf("failed to write decrypted file: %w", err)
	}
	
	return decryptedFilePath, nil
}

// decryptWithKey attempts to decrypt data with a specific key
func (s *encryptionService) decryptWithKey(ciphertext, key []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Check if ciphertext is long enough
	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	
	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	
	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	
	return plaintext, nil
}

// IsEncrypted checks if a file is encrypted
func (s *encryptionService) IsEncrypted(filePath string) bool {
	return strings.HasSuffix(filePath, s.encryptedExt)
}

// GetEncryptionKey returns the encryption key for a specific sensitive type
func (s *encryptionService) GetEncryptionKey(sensitiveType string) ([]byte, error) {
	// Check if we have a key for this type
	if key, ok := s.typeKeys[sensitiveType]; ok {
		return key, nil
	}
	
	// If no type-specific key, derive one from the master key and type
	key := deriveKey(s.masterKey, []byte(sensitiveType), 32)
	
	// Store the derived key for future use
	s.typeKeys[sensitiveType] = key
	
	return key, nil
}

// deriveKey derives a key from a password and salt using PBKDF2-like approach
// This is a simplified version for demonstration purposes
func deriveKey(password, salt []byte, keyLen int) []byte {
	if salt == nil {
		salt = []byte("default-salt-for-encryption-service")
	}
	
	// Create a SHA-256 hash of password and salt
	hash := sha256.New()
	hash.Write(password)
	hash.Write(salt)
	key := hash.Sum(nil)
	
	// If key is not long enough, keep hashing
	for len(key) < keyLen {
		hash.Reset()
		hash.Write(key)
		key = append(key, hash.Sum(nil)...)
	}
	
	// Truncate key to desired length
	return key[:keyLen]
}

// GenerateRandomKey generates a random encryption key
func GenerateRandomKey(keyLen int) ([]byte, error) {
	key := make([]byte, keyLen)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

// EncodeKeyToHex encodes a key to hex string
func EncodeKeyToHex(key []byte) string {
	return hex.EncodeToString(key)
}

// DecodeKeyFromHex decodes a key from hex string
func DecodeKeyFromHex(hexKey string) ([]byte, error) {
	return hex.DecodeString(hexKey)
}