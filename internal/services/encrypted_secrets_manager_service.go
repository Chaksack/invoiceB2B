package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"invoiceB2B/internal/interfaces"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SecretVersion represents a versioned secret
type SecretVersion struct {
	Version    int               `json:"version"`
	Data       map[string]string `json:"data"`
	CreatedAt  time.Time         `json:"created_at"`
	RotatedAt  *time.Time        `json:"rotated_at,omitempty"`
	RotatedBy  string            `json:"rotated_by,omitempty"`
	AccessedAt *time.Time        `json:"accessed_at,omitempty"`
}

// EncryptedSecret represents an encrypted secret with version history
type EncryptedSecret struct {
	CurrentVersion int                      `json:"current_version"`
	Versions       map[int][]byte           `json:"versions"` // Encrypted version data
	Metadata       map[string]string        `json:"metadata"` // Unencrypted metadata
	AccessLog      []map[string]interface{} `json:"access_log"`
}

// encryptedSecretsManagerService implements interfaces.SecretsManagerService with encryption
type encryptedSecretsManagerService struct {
	secretsDir string
	masterKey  []byte
	mu         sync.RWMutex
}

// NewEncryptedSecretsManagerService creates a new encrypted secrets manager service
func NewEncryptedSecretsManagerService(secretsDir string, masterKeyHex string) (interfaces.SecretsManagerService, error) {
	// If no directory is provided, use a default
	if secretsDir == "" {
		secretsDir = "./secrets"
	}
	
	// Ensure the secrets directory exists
	if err := os.MkdirAll(secretsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create secrets directory: %w", err)
	}
	
	// Derive master key from provided key or generate a new one
	var masterKey []byte
	if masterKeyHex == "" {
		// Generate a new random master key
		masterKey = make([]byte, 32) // 256 bits
		if _, err := io.ReadFull(rand.Reader, masterKey); err != nil {
			return nil, fmt.Errorf("failed to generate master key: %w", err)
		}
		
		// Save the master key to a file for future use
		masterKeyPath := filepath.Join(secretsDir, "master.key")
		if err := os.WriteFile(masterKeyPath, []byte(base64.StdEncoding.EncodeToString(masterKey)), 0600); err != nil {
			return nil, fmt.Errorf("failed to save master key: %w", err)
		}
	} else {
		// Use the provided master key
		var err error
		masterKey, err = base64.StdEncoding.DecodeString(masterKeyHex)
		if err != nil {
			// If not base64, derive a key from the string
			masterKey = deriveSecretKey([]byte(masterKeyHex), nil, 32)
		}
	}
	
	return &encryptedSecretsManagerService{
		secretsDir: secretsDir,
		masterKey:  masterKey,
	}, nil
}

// getSecretFilePath returns the file path for a secret
func (s *encryptedSecretsManagerService) getSecretFilePath(secretName string) string {
	// Sanitize the secret name to ensure it's a valid filename
	sanitizedName := filepath.Base(secretName)
	return filepath.Join(s.secretsDir, sanitizedName+".enc.json")
}

// deriveSecretKey derives a key from a password and salt
func deriveSecretKey(password, salt []byte, keyLen int) []byte {
	if salt == nil {
		salt = []byte("default-salt-for-secrets-manager")
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

// encrypt encrypts data using AES-GCM
func (s *encryptedSecretsManagerService) encrypt(data []byte, secretName string) ([]byte, error) {
	// Derive a key for this specific secret
	key := deriveSecretKey(s.masterKey, []byte(secretName), 32)
	
	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to create nonce: %w", err)
	}
	
	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func (s *encryptedSecretsManagerService) decrypt(ciphertext []byte, secretName string) ([]byte, error) {
	// Derive a key for this specific secret
	key := deriveSecretKey(s.masterKey, []byte(secretName), 32)
	
	// Create cipher block
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
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	
	return plaintext, nil
}

// loadEncryptedSecret loads an encrypted secret from disk
func (s *encryptedSecretsManagerService) loadEncryptedSecret(secretName string) (*EncryptedSecret, error) {
	filePath := s.getSecretFilePath(secretName)
	
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("secret %s not found", secretName)
	}
	
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret file: %w", err)
	}
	
	// Parse the JSON
	var encryptedSecret EncryptedSecret
	if err := json.Unmarshal(data, &encryptedSecret); err != nil {
		return nil, fmt.Errorf("failed to parse secret data: %w", err)
	}
	
	return &encryptedSecret, nil
}

// saveEncryptedSecret saves an encrypted secret to disk
func (s *encryptedSecretsManagerService) saveEncryptedSecret(secretName string, encryptedSecret *EncryptedSecret) error {
	filePath := s.getSecretFilePath(secretName)
	
	// Convert the secret to JSON
	data, err := json.MarshalIndent(encryptedSecret, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}
	
	// Write the file with restricted permissions
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}
	
	return nil
}

// logAccess logs access to a secret
func (s *encryptedSecretsManagerService) logAccess(encryptedSecret *EncryptedSecret, operation string, userID string) {
	// Create a log entry
	logEntry := map[string]interface{}{
		"operation":  operation,
		"timestamp":  time.Now().Format(time.RFC3339),
		"user_id":    userID,
		"version":    encryptedSecret.CurrentVersion,
	}
	
	// Add the log entry to the access log
	encryptedSecret.AccessLog = append(encryptedSecret.AccessLog, logEntry)
	
	// Limit the size of the access log
	if len(encryptedSecret.AccessLog) > 100 {
		encryptedSecret.AccessLog = encryptedSecret.AccessLog[len(encryptedSecret.AccessLog)-100:]
	}
}

// GetSecret retrieves a secret
func (s *encryptedSecretsManagerService) GetSecret(secretName string) (map[string]string, error) {
	return s.GetSecretVersion(secretName, -1, "")
}

// GetSecretVersion retrieves a specific version of a secret
func (s *encryptedSecretsManagerService) GetSecretVersion(secretName string, version int, userID string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Load the encrypted secret
	encryptedSecret, err := s.loadEncryptedSecret(secretName)
	if err != nil {
		return nil, err
	}
	
	// If version is -1, use the current version
	if version == -1 {
		version = encryptedSecret.CurrentVersion
	}
	
	// Check if the requested version exists
	encryptedData, ok := encryptedSecret.Versions[version]
	if !ok {
		return nil, fmt.Errorf("version %d of secret %s not found", version, secretName)
	}
	
	// Decrypt the data
	decryptedData, err := s.decrypt(encryptedData, secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}
	
	// Parse the JSON
	var secretVersion SecretVersion
	if err := json.Unmarshal(decryptedData, &secretVersion); err != nil {
		return nil, fmt.Errorf("failed to parse secret data: %w", err)
	}
	
	// Update access time
	now := time.Now()
	secretVersion.AccessedAt = &now
	
	// Re-encrypt and save the updated version
	updatedData, err := json.Marshal(secretVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated secret data: %w", err)
	}
	
	encryptedUpdatedData, err := s.encrypt(updatedData, secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt updated secret data: %w", err)
	}
	
	encryptedSecret.Versions[version] = encryptedUpdatedData
	
	// Log the access
	s.logAccess(encryptedSecret, "get", userID)
	
	// Save the updated secret
	if err := s.saveEncryptedSecret(secretName, encryptedSecret); err != nil {
		log.Printf("Warning: Failed to save updated access time for secret %s: %v", secretName, err)
	}
	
	return secretVersion.Data, nil
}

// StoreSecret stores a secret
func (s *encryptedSecretsManagerService) StoreSecret(secretName string, secretValue map[string]string) error {
	return s.StoreSecretWithUser(secretName, secretValue, "")
}

// StoreSecretWithUser stores a secret with user information
func (s *encryptedSecretsManagerService) StoreSecretWithUser(secretName string, secretValue map[string]string, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Try to load existing secret
	var encryptedSecret *EncryptedSecret
	var err error
	
	encryptedSecret, err = s.loadEncryptedSecret(secretName)
	if err != nil {
		// If the secret doesn't exist, create a new one
		if os.IsNotExist(err) {
			encryptedSecret = &EncryptedSecret{
				CurrentVersion: 1,
				Versions:       make(map[int][]byte),
				Metadata:       make(map[string]string),
				AccessLog:      make([]map[string]interface{}, 0),
			}
		} else {
			return err
		}
	}
	
	// Create a new version
	newVersion := encryptedSecret.CurrentVersion + 1
	
	// Create the secret version
	secretVersion := SecretVersion{
		Version:   newVersion,
		Data:      secretValue,
		CreatedAt: time.Now(),
	}
	
	// Convert to JSON
	versionData, err := json.Marshal(secretVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}
	
	// Encrypt the data
	encryptedData, err := s.encrypt(versionData, secretName)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret data: %w", err)
	}
	
	// Update the encrypted secret
	encryptedSecret.Versions[newVersion] = encryptedData
	encryptedSecret.CurrentVersion = newVersion
	
	// Log the access
	s.logAccess(encryptedSecret, "store", userID)
	
	// Save the updated secret
	if err := s.saveEncryptedSecret(secretName, encryptedSecret); err != nil {
		return fmt.Errorf("failed to save secret: %w", err)
	}
	
	return nil
}

// RotateSecret rotates a secret using the provided rotation function
func (s *encryptedSecretsManagerService) RotateSecret(secretName string, rotationFunction func() (map[string]string, error)) error {
	return s.RotateSecretWithUser(secretName, rotationFunction, "")
}

// RotateSecretWithUser rotates a secret using the provided rotation function with user information
func (s *encryptedSecretsManagerService) RotateSecretWithUser(secretName string, rotationFunction func() (map[string]string, error), userID string) error {
	// Load the current secret
	currentSecret, err := s.GetSecret(secretName)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to get current secret: %w", err)
	}
	
	// Generate new secret values
	newSecret, err := rotationFunction()
	if err != nil {
		return fmt.Errorf("failed to generate new secret values: %w", err)
	}
	
	// If the current secret exists, merge any keys that weren't rotated
	if currentSecret != nil {
		for k, v := range currentSecret {
			if _, exists := newSecret[k]; !exists {
				newSecret[k] = v
			}
		}
	}
	
	// Store the new secret
	if err := s.StoreSecretWithUser(secretName, newSecret, userID); err != nil {
		return fmt.Errorf("failed to store new secret: %w", err)
	}
	
	log.Printf("Successfully rotated secret: %s", secretName)
	return nil
}

// GetSecretString retrieves a specific secret value by key
func (s *encryptedSecretsManagerService) GetSecretString(secretName, key string) (string, error) {
	return s.GetSecretStringWithUser(secretName, key, "")
}

// GetSecretStringWithUser retrieves a specific secret value by key with user information
func (s *encryptedSecretsManagerService) GetSecretStringWithUser(secretName, key, userID string) (string, error) {
	secretMap, err := s.GetSecretVersion(secretName, -1, userID)
	if err != nil {
		return "", err
	}
	
	value, ok := secretMap[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s", key, secretName)
	}
	
	return value, nil
}

// GetSecretHistory returns the version history of a secret
func (s *encryptedSecretsManagerService) GetSecretHistory(secretName string, userID string) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Load the encrypted secret
	encryptedSecret, err := s.loadEncryptedSecret(secretName)
	if err != nil {
		return nil, err
	}
	
	// Create a history list
	history := make([]map[string]interface{}, 0, len(encryptedSecret.Versions))
	
	// Add each version to the history
	for version := range encryptedSecret.Versions {
		// Decrypt the version data
		encryptedData := encryptedSecret.Versions[version]
		decryptedData, err := s.decrypt(encryptedData, secretName)
		if err != nil {
			log.Printf("Warning: Failed to decrypt version %d of secret %s: %v", version, secretName, err)
			continue
		}
		
		// Parse the JSON
		var secretVersion SecretVersion
		if err := json.Unmarshal(decryptedData, &secretVersion); err != nil {
			log.Printf("Warning: Failed to parse version %d of secret %s: %v", version, secretName, err)
			continue
		}
		
		// Add to history
		historyEntry := map[string]interface{}{
			"version":     secretVersion.Version,
			"created_at":  secretVersion.CreatedAt,
			"rotated_at":  secretVersion.RotatedAt,
			"rotated_by":  secretVersion.RotatedBy,
			"accessed_at": secretVersion.AccessedAt,
		}
		
		history = append(history, historyEntry)
	}
	
	// Log the access
	s.logAccess(encryptedSecret, "get_history", userID)
	
	// Save the updated secret
	if err := s.saveEncryptedSecret(secretName, encryptedSecret); err != nil {
		log.Printf("Warning: Failed to save access log for secret %s: %v", secretName, err)
	}
	
	return history, nil
}

// GetAccessLog returns the access log for a secret
func (s *encryptedSecretsManagerService) GetAccessLog(secretName string, userID string) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Load the encrypted secret
	encryptedSecret, err := s.loadEncryptedSecret(secretName)
	if err != nil {
		return nil, err
	}
	
	// Log the access
	s.logAccess(encryptedSecret, "get_access_log", userID)
	
	// Save the updated secret
	if err := s.saveEncryptedSecret(secretName, encryptedSecret); err != nil {
		log.Printf("Warning: Failed to save access log for secret %s: %v", secretName, err)
	}
	
	return encryptedSecret.AccessLog, nil
}