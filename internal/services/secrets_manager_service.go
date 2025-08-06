package services

import (
	"encoding/json"
	"fmt"
	"invoiceB2B/internal/interfaces"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// localSecretsManagerService implements interfaces.SecretsManagerService using local file storage
// This is a placeholder implementation that can be replaced with AWS Secrets Manager later
type localSecretsManagerService struct {
	secretsDir string
	mu         sync.RWMutex
}

// NewSecretsManagerService creates a new secrets manager service
func NewSecretsManagerService(secretsDir string) (interfaces.SecretsManagerService, error) {
	// If no directory is provided, use a default
	if secretsDir == "" {
		secretsDir = "./secrets"
	}
	
	// Ensure the secrets directory exists
	if err := os.MkdirAll(secretsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create secrets directory: %w", err)
	}
	
	return &localSecretsManagerService{
		secretsDir: secretsDir,
	}, nil
}

// getSecretFilePath returns the file path for a secret
func (s *localSecretsManagerService) getSecretFilePath(secretName string) string {
	// Sanitize the secret name to ensure it's a valid filename
	sanitizedName := filepath.Base(secretName)
	return filepath.Join(s.secretsDir, sanitizedName+".json")
}

// GetSecret retrieves a secret from local storage
func (s *localSecretsManagerService) GetSecret(secretName string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
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
	var secretMap map[string]string
	if err := json.Unmarshal(data, &secretMap); err != nil {
		return nil, fmt.Errorf("failed to parse secret data: %w", err)
	}
	
	return secretMap, nil
}

// StoreSecret stores a secret in local storage
func (s *localSecretsManagerService) StoreSecret(secretName string, secretValue map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	filePath := s.getSecretFilePath(secretName)
	
	// Convert the secret to JSON
	data, err := json.MarshalIndent(secretValue, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}
	
	// Write the file with restricted permissions
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}
	
	return nil
}

// RotateSecret rotates a secret using the provided rotation function
func (s *localSecretsManagerService) RotateSecret(secretName string, rotationFunction func() (map[string]string, error)) error {
	// Check if the secret exists, but we don't need the actual value
	// In a real AWS implementation, we might pass the current secret to the rotation function
	_, err := s.GetSecret(secretName)
	if err != nil && !os.IsNotExist(err) {
		// Only return error if it's not a "not found" error
		return fmt.Errorf("failed to get current secret: %w", err)
	}
	
	// Generate new secret values
	newSecret, err := rotationFunction()
	if err != nil {
		return fmt.Errorf("failed to generate new secret values: %w", err)
	}
	
	// Store the new secret
	if err := s.StoreSecret(secretName, newSecret); err != nil {
		return fmt.Errorf("failed to store new secret: %w", err)
	}
	
	log.Printf("Successfully rotated secret: %s", secretName)
	return nil
}

// GetSecretString retrieves a specific secret value by key
func (s *localSecretsManagerService) GetSecretString(secretName, key string) (string, error) {
	secretMap, err := s.GetSecret(secretName)
	if err != nil {
		return "", err
	}
	
	value, ok := secretMap[key]
	if !ok {
		return "", fmt.Errorf("key %s not found in secret %s", key, secretName)
	}
	
	return value, nil
}

// GetSecretVersion retrieves a specific version of a secret
// This implementation doesn't support versioning, so it just returns the current secret
func (s *localSecretsManagerService) GetSecretVersion(secretName string, version int, userID string) (map[string]string, error) {
	return s.GetSecret(secretName)
}

// StoreSecretWithUser stores a secret with user information for audit logging
// This implementation ignores the user information
func (s *localSecretsManagerService) StoreSecretWithUser(secretName string, secretValue map[string]string, userID string) error {
	return s.StoreSecret(secretName, secretValue)
}

// RotateSecretWithUser rotates a secret with user information for audit logging
// This implementation ignores the user information
func (s *localSecretsManagerService) RotateSecretWithUser(secretName string, rotationFunction func() (map[string]string, error), userID string) error {
	return s.RotateSecret(secretName, rotationFunction)
}

// GetSecretStringWithUser retrieves a specific secret value by key with user information for audit logging
// This implementation ignores the user information
func (s *localSecretsManagerService) GetSecretStringWithUser(secretName, key, userID string) (string, error) {
	return s.GetSecretString(secretName, key)
}

// GetSecretHistory returns the version history of a secret
// This implementation doesn't support versioning, so it returns an empty history
func (s *localSecretsManagerService) GetSecretHistory(secretName string, userID string) ([]map[string]interface{}, error) {
	// Check if the secret exists
	_, err := s.GetSecret(secretName)
	if err != nil {
		return nil, err
	}
	
	// Return an empty history
	return []map[string]interface{}{}, nil
}

// GetAccessLog returns the access log for a secret
// This implementation doesn't support access logging, so it returns an empty log
func (s *localSecretsManagerService) GetAccessLog(secretName string, userID string) ([]map[string]interface{}, error) {
	// Check if the secret exists
	_, err := s.GetSecret(secretName)
	if err != nil {
		return nil, err
	}
	
	// Return an empty log
	return []map[string]interface{}{}, nil
}