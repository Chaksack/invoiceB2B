package secureconfig

import (
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/interfaces"
	"invoiceB2B/internal/utils/secrets"
	"log"
	"strconv"
	"strings"
	"time"
)

// SecureConfig extends Config with secrets management capabilities
type SecureConfig struct {
	Config         *config.Config
	secretsManager interfaces.SecretsManagerService
	secretsLoaded  bool
}

// NewSecureConfig creates a new secure configuration with secrets management
func NewSecureConfig(baseConfig *config.Config, secretsManager interfaces.SecretsManagerService) *SecureConfig {
	return &SecureConfig{
		Config:         baseConfig,
		secretsManager: secretsManager,
		secretsLoaded:  false,
	}
}

// LoadSecrets loads secrets from the secrets manager
func (sc *SecureConfig) LoadSecrets() error {
	if sc.secretsLoaded {
		return nil // Already loaded
	}
	
	// Define the secret names we'll use
	dbSecretName := fmt.Sprintf("db-credentials-%s", sc.Config.AppEnv)
	jwtSecretName := fmt.Sprintf("jwt-credentials-%s", sc.Config.AppEnv)
	smtpSecretName := fmt.Sprintf("smtp-credentials-%s", sc.Config.AppEnv)
	apiSecretName := fmt.Sprintf("api-credentials-%s", sc.Config.AppEnv)
	
	// Try to load database credentials
	dbSecret, err := sc.secretsManager.GetSecret(dbSecretName)
	if err == nil {
		log.Printf("Loaded database credentials from secrets manager")
		if host, ok := dbSecret["host"]; ok {
			sc.Config.DBHost = host
		}
		if port, ok := dbSecret["port"]; ok {
			sc.Config.DBPort = port
		}
		if user, ok := dbSecret["user"]; ok {
			sc.Config.DBUser = user
		}
		if password, ok := dbSecret["password"]; ok {
			sc.Config.DBPassword = password
		}
		if name, ok := dbSecret["name"]; ok {
			sc.Config.DBName = name
		}
		if sslMode, ok := dbSecret["sslmode"]; ok {
			sc.Config.DBSslMode = sslMode
		}
		
		// Rebuild the DSN with the new values
		sc.Config.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			sc.Config.DBHost, sc.Config.DBUser, sc.Config.DBPassword, sc.Config.DBName, sc.Config.DBPort, sc.Config.DBSslMode)
	} else {
		log.Printf("Database credentials not found in secrets manager: %v. Using environment variables.", err)
	}
	
	// Try to load JWT credentials
	jwtSecret, err := sc.secretsManager.GetSecret(jwtSecretName)
	if err == nil {
		log.Printf("Loaded JWT credentials from secrets manager")
		if secret, ok := jwtSecret["secret"]; ok {
			sc.Config.JWTSecret = secret
		}
		if privateKey, ok := jwtSecret["private_key"]; ok {
			sc.Config.JWTPrivateKey = privateKey
		}
		if accessTokenExp, ok := jwtSecret["access_token_expiration_minutes"]; ok {
			if minutes, err := strconv.Atoi(accessTokenExp); err == nil {
				sc.Config.JWTAccessTokenExpirationMinutes = time.Duration(minutes) * time.Minute
			}
		}
		if refreshTokenExp, ok := jwtSecret["refresh_token_expiration_days"]; ok {
			if days, err := strconv.Atoi(refreshTokenExp); err == nil {
				sc.Config.JWTRefreshTokenExpirationDays = time.Duration(days) * 24 * time.Hour
			}
		}
	} else {
		log.Printf("JWT credentials not found in secrets manager: %v. Using environment variables.", err)
	}
	
	// Try to load SMTP credentials
	smtpSecret, err := sc.secretsManager.GetSecret(smtpSecretName)
	if err == nil {
		log.Printf("Loaded SMTP credentials from secrets manager")
		if host, ok := smtpSecret["host"]; ok {
			sc.Config.SMTPHost = host
		}
		if portStr, ok := smtpSecret["port"]; ok {
			if port, err := strconv.Atoi(portStr); err == nil {
				sc.Config.SMTPPort = port
			}
		}
		if user, ok := smtpSecret["user"]; ok {
			sc.Config.SMTPUser = user
		}
		if password, ok := smtpSecret["password"]; ok {
			sc.Config.SMTPPassword = password
		}
		if senderEmail, ok := smtpSecret["sender_email"]; ok {
			sc.Config.SMTPSenderEmail = senderEmail
		}
	} else {
		log.Printf("SMTP credentials not found in secrets manager: %v. Using environment variables.", err)
	}
	
	// Try to load API credentials
	apiSecret, err := sc.secretsManager.GetSecret(apiSecretName)
	if err == nil {
		log.Printf("Loaded API credentials from secrets manager")
		if internalApiKey, ok := apiSecret["internal_api_key"]; ok {
			sc.Config.InternalAPIKey = internalApiKey
		}
	} else {
		log.Printf("API credentials not found in secrets manager: %v. Using environment variables.", err)
	}
	
	sc.secretsLoaded = true
	return nil
}

// InitializeSecrets creates initial secrets in the secrets manager based on current config
func (sc *SecureConfig) InitializeSecrets() error {
	// Only initialize secrets if they don't exist yet
	
	// Database credentials
	dbSecretName := fmt.Sprintf("db-credentials-%s", sc.Config.AppEnv)
	_, err := sc.secretsManager.GetSecret(dbSecretName)
	if err != nil && strings.Contains(err.Error(), "not found") {
		dbSecret := map[string]string{
			"host":     sc.Config.DBHost,
			"port":     sc.Config.DBPort,
			"user":     sc.Config.DBUser,
			"password": sc.Config.DBPassword,
			"name":     sc.Config.DBName,
			"sslmode":  sc.Config.DBSslMode,
		}
		if err := sc.secretsManager.StoreSecret(dbSecretName, dbSecret); err != nil {
			log.Printf("Warning: Failed to store database credentials in secrets manager: %v", err)
		} else {
			log.Printf("Initialized database credentials in secrets manager")
		}
	}
	
	// JWT credentials
	jwtSecretName := fmt.Sprintf("jwt-credentials-%s", sc.Config.AppEnv)
	_, err = sc.secretsManager.GetSecret(jwtSecretName)
	if err != nil && strings.Contains(err.Error(), "not found") {
		jwtSecret := map[string]string{
			"secret":                          sc.Config.JWTSecret,
			"private_key":                     sc.Config.JWTPrivateKey,
			"access_token_expiration_minutes": strconv.Itoa(int(sc.Config.JWTAccessTokenExpirationMinutes.Minutes())),
			"refresh_token_expiration_days":   strconv.Itoa(int(sc.Config.JWTRefreshTokenExpirationDays.Hours() / 24)),
		}
		if err := sc.secretsManager.StoreSecret(jwtSecretName, jwtSecret); err != nil {
			log.Printf("Warning: Failed to store JWT credentials in secrets manager: %v", err)
		} else {
			log.Printf("Initialized JWT credentials in secrets manager")
		}
	}
	
	// SMTP credentials
	smtpSecretName := fmt.Sprintf("smtp-credentials-%s", sc.Config.AppEnv)
	_, err = sc.secretsManager.GetSecret(smtpSecretName)
	if err != nil && strings.Contains(err.Error(), "not found") {
		smtpSecret := map[string]string{
			"host":         sc.Config.SMTPHost,
			"port":         strconv.Itoa(sc.Config.SMTPPort),
			"user":         sc.Config.SMTPUser,
			"password":     sc.Config.SMTPPassword,
			"sender_email": sc.Config.SMTPSenderEmail,
		}
		if err := sc.secretsManager.StoreSecret(smtpSecretName, smtpSecret); err != nil {
			log.Printf("Warning: Failed to store SMTP credentials in secrets manager: %v", err)
		} else {
			log.Printf("Initialized SMTP credentials in secrets manager")
		}
	}
	
	// API credentials
	apiSecretName := fmt.Sprintf("api-credentials-%s", sc.Config.AppEnv)
	_, err = sc.secretsManager.GetSecret(apiSecretName)
	if err != nil && strings.Contains(err.Error(), "not found") {
		apiSecret := map[string]string{
			"internal_api_key": sc.Config.InternalAPIKey,
		}
		if err := sc.secretsManager.StoreSecret(apiSecretName, apiSecret); err != nil {
			log.Printf("Warning: Failed to store API credentials in secrets manager: %v", err)
		} else {
			log.Printf("Initialized API credentials in secrets manager")
		}
	}
	
	return nil
}

// SetupCredentialRotation configures automatic rotation of credentials
func (sc *SecureConfig) SetupCredentialRotation() {
	// Database password rotation (every 30 days)
	dbSecretName := fmt.Sprintf("db-credentials-%s", sc.Config.AppEnv)
	dbRotationFunc := func() (map[string]string, error) {
		// Get current secret
		currentSecret, err := sc.secretsManager.GetSecret(dbSecretName)
		if err != nil {
			return nil, fmt.Errorf("failed to get current DB secret: %w", err)
		}
		
		// Generate new password
		newPassword, err := secrets.GenerateSecurePassword(24)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new DB password: %w", err)
		}
		
		// In a real implementation, this would update the database user's password
		
		// Update the secret with the new password
		currentSecret["password"] = newPassword
		currentSecret["rotated_at"] = time.Now().Format(time.RFC3339)
		
		return currentSecret, nil
	}
	secrets.ScheduleRotation(sc.secretsManager, dbSecretName, dbRotationFunc, 30*24*time.Hour) // 30 days
	
	// JWT secret rotation (every 90 days)
	jwtSecretName := fmt.Sprintf("jwt-credentials-%s", sc.Config.AppEnv)
	jwtRotationFunc := func() (map[string]string, error) {
		// Get current secret
		currentSecret, err := sc.secretsManager.GetSecret(jwtSecretName)
		if err != nil {
			return nil, fmt.Errorf("failed to get current JWT secret: %w", err)
		}
		
		// Generate new JWT secret
		newSecret, err := secrets.GenerateSecureToken(32)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new JWT secret: %w", err)
		}
		
		// Update the secret
		currentSecret["secret"] = newSecret
		currentSecret["rotated_at"] = time.Now().Format(time.RFC3339)
		
		return currentSecret, nil
	}
	secrets.ScheduleRotation(sc.secretsManager, jwtSecretName, jwtRotationFunc, 90*24*time.Hour) // 90 days
	
	// API key rotation (every 60 days)
	apiSecretName := fmt.Sprintf("api-credentials-%s", sc.Config.AppEnv)
	apiRotationFunc := func() (map[string]string, error) {
		// Get current secret
		currentSecret, err := sc.secretsManager.GetSecret(apiSecretName)
		if err != nil {
			return nil, fmt.Errorf("failed to get current API secret: %w", err)
		}
		
		// Generate new API key
		newApiKey, err := secrets.GenerateSecureToken(32)
		if err != nil {
			return nil, fmt.Errorf("failed to generate new API key: %w", err)
		}
		
		// Update the secret
		currentSecret["internal_api_key"] = newApiKey
		currentSecret["rotated_at"] = time.Now().Format(time.RFC3339)
		
		return currentSecret, nil
	}
	secrets.ScheduleRotation(sc.secretsManager, apiSecretName, apiRotationFunc, 60*24*time.Hour) // 60 days
}