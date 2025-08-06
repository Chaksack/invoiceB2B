package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"invoiceB2B/internal/interfaces"
	"log"
	"time"
)

// ScheduleRotation schedules automatic rotation of a secret
func ScheduleRotation(secretsService interfaces.SecretsManagerService, secretName string, rotationFunction func() (map[string]string, error), interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := secretsService.RotateSecret(secretName, rotationFunction); err != nil {
				log.Printf("Failed to rotate secret %s: %v", secretName, err)
			}
		}
	}()
}

// GenerateSecurePassword generates a cryptographically secure random password
func GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 16 // Minimum password length for security
	}
	
	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	// Convert to base64 to get printable characters
	password := base64.StdEncoding.EncodeToString(randomBytes)
	
	// Trim to desired length
	if len(password) > length {
		password = password[:length]
	}
	
	return password, nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	if length < 16 {
		length = 32 // Minimum token length for security
	}
	
	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	// Convert to base64 to get printable characters
	token := base64.StdEncoding.EncodeToString(randomBytes)
	
	// Trim to desired length
	if len(token) > length {
		token = token[:length]
	}
	
	return token, nil
}

// CreateDatabaseRotationFunction creates a function for rotating database credentials
func CreateDatabaseRotationFunction(passwordLength int) func() (map[string]string, error) {
	return func() (map[string]string, error) {
		// Generate a new secure password
		newPassword, err := GenerateSecurePassword(passwordLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new database password: %w", err)
		}
		
		// In a real implementation, this would update the database user's password
		// For now, we just return the new credentials
		return map[string]string{
			"password": newPassword,
			"rotated_at": time.Now().Format(time.RFC3339),
		}, nil
	}
}

// CreateJWTSecretRotationFunction creates a function for rotating JWT secrets
func CreateJWTSecretRotationFunction(secretLength int) func() (map[string]string, error) {
	return func() (map[string]string, error) {
		// Generate a new secure JWT secret
		newSecret, err := GenerateSecureToken(secretLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new JWT secret: %w", err)
		}
		
		return map[string]string{
			"secret": newSecret,
			"rotated_at": time.Now().Format(time.RFC3339),
		}, nil
	}
}

// CreateAPIKeyRotationFunction creates a function for rotating API keys
func CreateAPIKeyRotationFunction(keyLength int) func() (map[string]string, error) {
	return func() (map[string]string, error) {
		// Generate a new secure API key
		newAPIKey, err := GenerateSecureToken(keyLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new API key: %w", err)
		}
		
		return map[string]string{
			"api_key": newAPIKey,
			"rotated_at": time.Now().Format(time.RFC3339),
		}, nil
	}
}