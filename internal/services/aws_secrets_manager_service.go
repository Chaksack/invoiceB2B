package services

import (
	"context"
	"encoding/json"
	"fmt"
	"invoiceB2B/internal/interfaces"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AWSSecretsManagerConfig holds configuration for AWS Secrets Manager
type AWSSecretsManagerConfig struct {
	Region          string // AWS region
	SecretNamespace string // Namespace prefix for secrets
	TagPrefix       string // Prefix for tags
}

// awsSecretsManagerService implements interfaces.SecretsManagerService using AWS Secrets Manager
type awsSecretsManagerService struct {
	client    *secretsmanager.Client
	namespace string
	tagPrefix string
}

// NewAWSSecretsManagerService creates a new AWS Secrets Manager service
func NewAWSSecretsManagerService(cfg AWSSecretsManagerConfig) (interfaces.SecretsManagerService, error) {
	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create AWS Secrets Manager client
	client := secretsmanager.NewFromConfig(awsCfg)

	return &awsSecretsManagerService{
		client:    client,
		namespace: cfg.SecretNamespace,
		tagPrefix: cfg.TagPrefix,
	}, nil
}

// getFullSecretName returns the full secret name with namespace
func (s *awsSecretsManagerService) getFullSecretName(secretName string) string {
	if s.namespace == "" {
		return secretName
	}
	return fmt.Sprintf("%s/%s", s.namespace, secretName)
}

// GetSecret retrieves a secret from AWS Secrets Manager
func (s *awsSecretsManagerService) GetSecret(secretName string) (map[string]string, error) {
	fullSecretName := s.getFullSecretName(secretName)

	// Get the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(fullSecretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := s.client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret value: %w", err)
	}

	// Parse the secret value
	var secretMap map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretMap); err != nil {
		return nil, fmt.Errorf("failed to parse secret data: %w", err)
	}

	return secretMap, nil
}

// StoreSecret stores a secret in AWS Secrets Manager
func (s *awsSecretsManagerService) StoreSecret(secretName string, secretValue map[string]string) error {
	return s.StoreSecretWithUser(secretName, secretValue, "")
}

// StoreSecretWithUser stores a secret with user information for audit logging
func (s *awsSecretsManagerService) StoreSecretWithUser(secretName string, secretValue map[string]string, userID string) error {
	fullSecretName := s.getFullSecretName(secretName)

	// Convert the secret to JSON
	secretJSON, err := json.Marshal(secretValue)
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}

	// Check if the secret already exists
	_, err = s.client.DescribeSecret(context.TODO(), &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(fullSecretName),
	})

	if err != nil {
		// Secret doesn't exist, create it
		tags := []secretsmanager.Tag{
			{
				Key:   aws.String(fmt.Sprintf("%s/created-by", s.tagPrefix)),
				Value: aws.String(userID),
			},
			{
				Key:   aws.String(fmt.Sprintf("%s/created-at", s.tagPrefix)),
				Value: aws.String(time.Now().Format(time.RFC3339)),
			},
		}

		_, err = s.client.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
			Name:         aws.String(fullSecretName),
			SecretString: aws.String(string(secretJSON)),
			Tags:         tags,
		})
		if err != nil {
			return fmt.Errorf("failed to create secret: %w", err)
		}
	} else {
		// Secret exists, update it
		_, err = s.client.PutSecretValue(context.TODO(), &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(fullSecretName),
			SecretString: aws.String(string(secretJSON)),
		})
		if err != nil {
			return fmt.Errorf("failed to update secret: %w", err)
		}

		// Update tags
		_, err = s.client.TagResource(context.TODO(), &secretsmanager.TagResourceInput{
			SecretId: aws.String(fullSecretName),
			Tags: []secretsmanager.Tag{
				{
					Key:   aws.String(fmt.Sprintf("%s/updated-by", s.tagPrefix)),
					Value: aws.String(userID),
				},
				{
					Key:   aws.String(fmt.Sprintf("%s/updated-at", s.tagPrefix)),
					Value: aws.String(time.Now().Format(time.RFC3339)),
				},
			},
		})
		if err != nil {
			log.Printf("Warning: Failed to update tags for secret %s: %v", secretName, err)
		}
	}

	return nil
}

// RotateSecret rotates a secret using the provided rotation function
func (s *awsSecretsManagerService) RotateSecret(secretName string, rotationFunction func() (map[string]string, error)) error {
	return s.RotateSecretWithUser(secretName, rotationFunction, "")
}

// RotateSecretWithUser rotates a secret using the provided rotation function with user information
func (s *awsSecretsManagerService) RotateSecretWithUser(secretName string, rotationFunction func() (map[string]string, error), userID string) error {
	// Get the current secret
	currentSecret, err := s.GetSecret(secretName)
	if err != nil {
		// If the secret doesn't exist, create it
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			newSecret, err := rotationFunction()
			if err != nil {
				return fmt.Errorf("failed to generate new secret values: %w", err)
			}
			return s.StoreSecretWithUser(secretName, newSecret, userID)
		}
		return fmt.Errorf("failed to get current secret: %w", err)
	}

	// Generate new secret values
	newSecret, err := rotationFunction()
	if err != nil {
		return fmt.Errorf("failed to generate new secret values: %w", err)
	}

	// Merge any keys that weren't rotated
	for k, v := range currentSecret {
		if _, exists := newSecret[k]; !exists {
			newSecret[k] = v
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
func (s *awsSecretsManagerService) GetSecretString(secretName, key string) (string, error) {
	return s.GetSecretStringWithUser(secretName, key, "")
}

// GetSecretStringWithUser retrieves a specific secret value by key with user information
func (s *awsSecretsManagerService) GetSecretStringWithUser(secretName, key, userID string) (string, error) {
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

// GetSecretVersion retrieves a specific version of a secret
func (s *awsSecretsManagerService) GetSecretVersion(secretName string, version int, userID string) (map[string]string, error) {
	fullSecretName := s.getFullSecretName(secretName)

	// Get the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(fullSecretName),
	}

	// If version is specified, use it
	if version > 0 {
		// List versions to find the version ID
		versionsInput := &secretsmanager.ListSecretVersionIdsInput{
			SecretId: aws.String(fullSecretName),
		}

		versionsOutput, err := s.client.ListSecretVersionIds(context.TODO(), versionsInput)
		if err != nil {
			return nil, fmt.Errorf("failed to list secret versions: %w", err)
		}

		// Find the version ID for the specified version
		if version > len(versionsOutput.Versions) {
			return nil, fmt.Errorf("version %d not found for secret %s", version, secretName)
		}

		// Sort versions by creation date
		// Note: This is a simplified approach. In a real implementation, you would need to
		// sort the versions by creation date and find the correct version ID.
		input.VersionId = versionsOutput.Versions[version-1].VersionId
	}

	result, err := s.client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret value: %w", err)
	}

	// Parse the secret value
	var secretMap map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretMap); err != nil {
		return nil, fmt.Errorf("failed to parse secret data: %w", err)
	}

	// Log access if user ID is provided
	if userID != "" {
		_, err = s.client.TagResource(context.TODO(), &secretsmanager.TagResourceInput{
			SecretId: aws.String(fullSecretName),
			Tags: []secretsmanager.Tag{
				{
					Key:   aws.String(fmt.Sprintf("%s/accessed-by", s.tagPrefix)),
					Value: aws.String(userID),
				},
				{
					Key:   aws.String(fmt.Sprintf("%s/accessed-at", s.tagPrefix)),
					Value: aws.String(time.Now().Format(time.RFC3339)),
				},
			},
		})
		if err != nil {
			log.Printf("Warning: Failed to update access tags for secret %s: %v", secretName, err)
		}
	}

	return secretMap, nil
}

// GetSecretHistory returns the version history of a secret
func (s *awsSecretsManagerService) GetSecretHistory(secretName string, userID string) ([]map[string]interface{}, error) {
	fullSecretName := s.getFullSecretName(secretName)

	// List versions
	input := &secretsmanager.ListSecretVersionIdsInput{
		SecretId:          aws.String(fullSecretName),
		IncludeDeprecated: aws.Bool(true),
	}

	result, err := s.client.ListSecretVersionIds(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to list secret versions: %w", err)
	}

	// Create a history list
	history := make([]map[string]interface{}, 0, len(result.Versions))

	// Add each version to the history
	for _, version := range result.Versions {
		historyEntry := map[string]interface{}{
			"version_id":   *version.VersionId,
			"created_date": version.CreatedDate.Format(time.RFC3339),
		}

		if version.LastAccessedDate != nil {
			historyEntry["last_accessed_date"] = version.LastAccessedDate.Format(time.RFC3339)
		}

		history = append(history, historyEntry)
	}

	// Log access if user ID is provided
	if userID != "" {
		_, err = s.client.TagResource(context.TODO(), &secretsmanager.TagResourceInput{
			SecretId: aws.String(fullSecretName),
			Tags: []secretsmanager.Tag{
				{
					Key:   aws.String(fmt.Sprintf("%s/history-accessed-by", s.tagPrefix)),
					Value: aws.String(userID),
				},
				{
					Key:   aws.String(fmt.Sprintf("%s/history-accessed-at", s.tagPrefix)),
					Value: aws.String(time.Now().Format(time.RFC3339)),
				},
			},
		})
		if err != nil {
			log.Printf("Warning: Failed to update history access tags for secret %s: %v", secretName, err)
		}
	}

	return history, nil
}

// GetAccessLog returns the access log for a secret
func (s *awsSecretsManagerService) GetAccessLog(secretName string, userID string) ([]map[string]interface{}, error) {
	fullSecretName := s.getFullSecretName(secretName)

	// Describe the secret to get tags
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(fullSecretName),
	}

	result, err := s.client.DescribeSecret(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe secret: %w", err)
	}

	// Create an access log from tags
	accessLog := make([]map[string]interface{}, 0)

	// Extract access information from tags
	for _, tag := range result.Tags {
		if strings.HasPrefix(*tag.Key, fmt.Sprintf("%s/accessed-", s.tagPrefix)) {
			// Parse the tag key to extract the timestamp
			parts := strings.Split(*tag.Key, "/")
			if len(parts) < 2 {
				continue
			}

			operation := parts[len(parts)-1]
			if operation == "accessed-by" {
				// Find the corresponding timestamp
				timestamp := ""
				for _, t := range result.Tags {
					if *t.Key == fmt.Sprintf("%s/accessed-at", s.tagPrefix) {
						timestamp = *t.Value
						break
					}
				}

				accessLog = append(accessLog, map[string]interface{}{
					"operation":  "get",
					"user_id":    *tag.Value,
					"timestamp":  timestamp,
					"version_id": "current", // AWS Secrets Manager doesn't track which version was accessed
				})
			}
		}
	}

	// Log this access if user ID is provided
	if userID != "" {
		_, err = s.client.TagResource(context.TODO(), &secretsmanager.TagResourceInput{
			SecretId: aws.String(fullSecretName),
			Tags: []secretsmanager.Tag{
				{
					Key:   aws.String(fmt.Sprintf("%s/log-accessed-by", s.tagPrefix)),
					Value: aws.String(userID),
				},
				{
					Key:   aws.String(fmt.Sprintf("%s/log-accessed-at", s.tagPrefix)),
					Value: aws.String(time.Now().Format(time.RFC3339)),
				},
			},
		})
		if err != nil {
			log.Printf("Warning: Failed to update log access tags for secret %s: %v", secretName, err)
		}
	}

	return accessLog, nil
}