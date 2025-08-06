package interfaces

// SecretsManagerService provides methods to interact with a secrets management system
type SecretsManagerService interface {
	// GetSecret retrieves a secret by name
	GetSecret(secretName string) (map[string]string, error)
	
	// StoreSecret stores a secret
	StoreSecret(secretName string, secretValue map[string]string) error
	
	// RotateSecret rotates a secret using the provided rotation function
	RotateSecret(secretName string, rotationFunction func() (map[string]string, error)) error
	
	// GetSecretString retrieves a specific secret value by key
	GetSecretString(secretName, key string) (string, error)
	
	// Enhanced methods for versioning, audit logging, and user tracking
	
	// GetSecretVersion retrieves a specific version of a secret
	// If version is -1, returns the current version
	GetSecretVersion(secretName string, version int, userID string) (map[string]string, error)
	
	// StoreSecretWithUser stores a secret with user information for audit logging
	StoreSecretWithUser(secretName string, secretValue map[string]string, userID string) error
	
	// RotateSecretWithUser rotates a secret with user information for audit logging
	RotateSecretWithUser(secretName string, rotationFunction func() (map[string]string, error), userID string) error
	
	// GetSecretStringWithUser retrieves a specific secret value by key with user information for audit logging
	GetSecretStringWithUser(secretName, key, userID string) (string, error)
	
	// GetSecretHistory returns the version history of a secret
	GetSecretHistory(secretName string, userID string) ([]map[string]interface{}, error)
	
	// GetAccessLog returns the access log for a secret
	GetAccessLog(secretName string, userID string) ([]map[string]interface{}, error)
}