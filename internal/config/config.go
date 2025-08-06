package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"log"
)

// Config holds all configuration for the application
type Config struct {
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string
	DSN        string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	RabbitMQURL                            string
	RabbitMQEventExchangeName              string
	RabbitMQUserRegisteredRoutingKey       string
	RabbitMQInvoiceUploadedRoutingKey      string
	RabbitMQInvoiceStatusUpdatedRoutingKey string
	RabbitMQKYCStatusUpdatedRoutingKey     string

	JWTSecret                       string
	JWTPrivateKey                   string // RSA private key for JWT signing (PEM format)
	JWTAccessTokenExpirationMinutes time.Duration
	JWTRefreshTokenExpirationDays   time.Duration

	// CSRF Protection Configuration
	CSRFSecret            string        // Secret key for CSRF token signing (separate from JWT)
	CSRFTokenExpiration   time.Duration // CSRF token expiration time
	CSRFCookieName        string        // Name of the cookie to store CSRF token
	CSRFHeaderName        string        // Name of the header to send CSRF token
	CSRFCookieSecure      bool          // Whether to set Secure flag on CSRF cookie
	CSRFCookieSameSite    string        // SameSite attribute for CSRF cookie (Strict, Lax, None)
	CSRFCookieHTTPOnly    bool          // Whether to set HttpOnly flag on CSRF cookie
	CSRFPerSessionTokens  bool          // Whether to use per-session CSRF tokens
	CSRFExemptPaths       []string      // Paths exempt from CSRF protection

	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPassword    string
	SMTPSenderEmail string

	OTPExpirationMinutes time.Duration
	UploadsDir           string
	MaxUploadSizeMB      int64
	InternalAPIKey       string
	
	// Virus Scanning Configuration
	ClamAVEnabled       bool   // Whether to use ClamAV for virus scanning
	ClamAVPath          string // Path to ClamAV binary
	VirusScanAPIEnabled bool   // Whether to use external virus scanning API
	VirusScanAPIURL     string // URL of the external virus scanning API
	VirusScanAPIKey     string // API key for the external virus scanning API
	
	// File Encryption Configuration
	EncryptionEnabled   bool              // Whether to encrypt sensitive files
	EncryptionMasterKey string            // Master key for file encryption (hex-encoded)
	EncryptionTempDir   string            // Directory for temporary decrypted files
	EncryptionFileExt   string            // File extension for encrypted files
	SensitiveFileTypes  map[string]bool   // Map of file types that should be encrypted
	TypeEncryptionKeys  map[string]string // Map of file types to their encryption keys (hex-encoded)
	
	// File Integrity Verification Configuration
	IntegrityVerificationEnabled bool   // Whether to verify file integrity
	IntegrityAlgorithm          string // Hash algorithm to use for integrity verification (sha256, sha512, hmac-sha256)
	IntegrityHMACKey            string // Secret key for HMAC (if using HMAC algorithm)
	IntegrityChecksumDir        string // Directory for storing checksums
	IntegrityVerifyOnAccess     bool   // Whether to verify integrity on file access
	IntegrityScheduleChecks     bool   // Whether to schedule periodic integrity checks
	IntegrityCheckInterval      int    // Interval for periodic integrity checks (in hours)
}

func LoadConfig(path string) (*Config, error) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
		log.Println("Warning: APP_ENV not set, defaulting to 'development'.")
	}

	log.Printf("Loading configuration with APP_ENV=%s", appEnv)

	if appEnv != "production" {
		envPath := path
		if envPath == "." || envPath == "" {
			envPath = ".env"
		} else {
			envPath = filepath.Join(path, ".env")
		}

		err := godotenv.Load(envPath)
		if err != nil {
			if appEnv == "development" {
				log.Printf("Warning: .env file not found or error loading from %s. Relying on OS environment variables or defaults. Error: %v\n", envPath, err)
			}
		} else {
			log.Printf("Loaded environment variables from %s (for non-production)\n", envPath)
		}
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "465"))
	accessTokenExpMinutes, _ := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_EXPIRATION_MINUTES", "15"))
	refreshTokenExpDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_EXPIRATION_DAYS", "7"))
	otpExpMinutes, _ := strconv.Atoi(getEnv("OTP_EXPIRATION_MINUTES", "5"))
	maxUploadSizeMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)
	csrfTokenExpMinutes, _ := strconv.Atoi(getEnv("CSRF_TOKEN_EXPIRATION_MINUTES", "30"))
	
	// Parse CSRF cookie secure flag
	csrfCookieSecure := true
	if getEnv("CSRF_COOKIE_SECURE", "true") == "false" {
		csrfCookieSecure = false
	}
	
	// Parse CSRF cookie HttpOnly flag
	csrfCookieHTTPOnly := true
	if getEnv("CSRF_COOKIE_HTTP_ONLY", "true") == "false" {
		csrfCookieHTTPOnly = false
	}
	
	// Parse CSRF per-session tokens flag
	csrfPerSessionTokens := true
	if getEnv("CSRF_PER_SESSION_TOKENS", "true") == "false" {
		csrfPerSessionTokens = false
	}
	
	// Parse CSRF exempt paths
	csrfExemptPathsStr := getEnv("CSRF_EXEMPT_PATHS", "/api/v1/auth/login,/api/v1/auth/register,/api/v1/internal")
	csrfExemptPaths := []string{}
	if csrfExemptPathsStr != "" {
		csrfExemptPaths = strings.Split(csrfExemptPathsStr, ",")
	}

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	calculatedRedisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Parse ClamAV enabled flag
	clamAVEnabled := false
	if getEnv("CLAMAV_ENABLED", "false") == "true" {
		clamAVEnabled = true
	}
	
	// Parse virus scan API enabled flag
	virusScanAPIEnabled := false
	if getEnv("VIRUS_SCAN_API_ENABLED", "false") == "true" {
		virusScanAPIEnabled = true
	}
	
	// Parse encryption enabled flag
	encryptionEnabled := false
	if getEnv("ENCRYPTION_ENABLED", "false") == "true" {
		encryptionEnabled = true
	}
	
	// Initialize sensitive file types map
	sensitiveFileTypes := make(map[string]bool)
	sensitiveTypesStr := getEnv("SENSITIVE_FILE_TYPES", "kyc,financial,personal")
	for _, fileType := range strings.Split(sensitiveTypesStr, ",") {
		sensitiveFileTypes[strings.TrimSpace(fileType)] = true
	}
	
	// Initialize type encryption keys map
	typeEncryptionKeys := make(map[string]string)
	// Parse type-specific encryption keys if provided
	typeKeysStr := getEnv("TYPE_ENCRYPTION_KEYS", "")
	if typeKeysStr != "" {
		for _, pair := range strings.Split(typeKeysStr, ",") {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) == 2 {
				fileType := strings.TrimSpace(parts[0])
				key := strings.TrimSpace(parts[1])
				typeEncryptionKeys[fileType] = key
			}
		}
	}
	
	// Parse integrity verification enabled flag
	integrityVerificationEnabled := false
	if getEnv("INTEGRITY_VERIFICATION_ENABLED", "false") == "true" {
		integrityVerificationEnabled = true
	}
	
	// Parse integrity verify on access flag
	integrityVerifyOnAccess := false
	if getEnv("INTEGRITY_VERIFY_ON_ACCESS", "false") == "true" {
		integrityVerifyOnAccess = true
	}
	
	// Parse integrity schedule checks flag
	integrityScheduleChecks := false
	if getEnv("INTEGRITY_SCHEDULE_CHECKS", "false") == "true" {
		integrityScheduleChecks = true
	}
	
	// Parse integrity check interval
	integrityCheckInterval, _ := strconv.Atoi(getEnv("INTEGRITY_CHECK_INTERVAL", "24"))
	
	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "3000"),
		AppEnv:     appEnv, // Use the appEnv determined at the start from OS
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "invoice_db"),
		DBSslMode:  getEnv("DB_SSLMODE", "disable"),

		RedisAddr:     calculatedRedisAddr,
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		RabbitMQURL:                            getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQEventExchangeName:              getEnv("RABBITMQ_EVENT_EXCHANGE_NAME", "invoice_events_exchange"),
		RabbitMQUserRegisteredRoutingKey:       getEnv("RABBITMQ_USER_REGISTERED_ROUTING_KEY", "user.registered"),
		RabbitMQInvoiceUploadedRoutingKey:      getEnv("RABBITMQ_INVOICE_UPLOADED_ROUTING_KEY", "invoice.uploaded"),
		RabbitMQInvoiceStatusUpdatedRoutingKey: getEnv("RABBITMQ_INVOICE_STATUS_UPDATED_ROUTING_KEY", "invoice.status.updated"),
		RabbitMQKYCStatusUpdatedRoutingKey:     getEnv("RABBITMQ_KYC_STATUS_UPDATED_ROUTING_KEY", "kyc.status.updated"),

		JWTSecret:                       getEnv("JWT_SECRET", "supersecretkey"),
		JWTPrivateKey:                   getEnv("JWT_PRIVATE_KEY", ""), // Empty default, will generate a key if not provided
		JWTAccessTokenExpirationMinutes: time.Duration(accessTokenExpMinutes) * time.Minute,
		JWTRefreshTokenExpirationDays:   time.Duration(refreshTokenExpDays) * 24 * time.Hour,

		// CSRF Configuration
		CSRFSecret:           getEnv("CSRF_SECRET", "csrf-secret-key-change-in-production"),
		CSRFTokenExpiration:  time.Duration(csrfTokenExpMinutes) * time.Minute,
		CSRFCookieName:       getEnv("CSRF_COOKIE_NAME", "X-CSRF-Token"),
		CSRFHeaderName:       getEnv("CSRF_HEADER_NAME", "X-CSRF-Token"),
		CSRFCookieSecure:     csrfCookieSecure,
		CSRFCookieSameSite:   getEnv("CSRF_COOKIE_SAME_SITE", "Strict"),
		CSRFCookieHTTPOnly:   csrfCookieHTTPOnly,
		CSRFPerSessionTokens: csrfPerSessionTokens,
		CSRFExemptPaths:      csrfExemptPaths,

		SMTPHost:        getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:        smtpPort,
		SMTPUser:        getEnv("SMTP_USER", "your_email@example.com"),
		SMTPPassword:    getEnv("SMTP_PASSWORD", "your_app_password"),
		SMTPSenderEmail: getEnv("SMTPSenderEmail", "Profundr <no-reply@profundr.io>"),

		OTPExpirationMinutes: time.Duration(otpExpMinutes) * time.Minute,
		UploadsDir:           getEnv("UPLOADS_DIR", "./uploads"),
		MaxUploadSizeMB:      maxUploadSizeMB,
		InternalAPIKey:       getEnv("INTERNAL_API_KEY", "default-internal-key-please-change"),
		
		// Virus Scanning Configuration
		ClamAVEnabled:       clamAVEnabled,
		ClamAVPath:          getEnv("CLAMAV_PATH", "clamscan"),
		VirusScanAPIEnabled: virusScanAPIEnabled,
		VirusScanAPIURL:     getEnv("VIRUS_SCAN_API_URL", ""),
		VirusScanAPIKey:     getEnv("VIRUS_SCAN_API_KEY", ""),
		
		// File Encryption Configuration
		EncryptionEnabled:   encryptionEnabled,
		EncryptionMasterKey: getEnv("ENCRYPTION_MASTER_KEY", "default-encryption-key-change-in-production"),
		EncryptionTempDir:   getEnv("ENCRYPTION_TEMP_DIR", os.TempDir()),
		EncryptionFileExt:   getEnv("ENCRYPTION_FILE_EXT", ".enc"),
		SensitiveFileTypes:  sensitiveFileTypes,
		TypeEncryptionKeys:  typeEncryptionKeys,
		
		// File Integrity Verification Configuration
		IntegrityVerificationEnabled: integrityVerificationEnabled,
		IntegrityAlgorithm:          getEnv("INTEGRITY_ALGORITHM", "sha256"),
		IntegrityHMACKey:            getEnv("INTEGRITY_HMAC_KEY", "default-hmac-key-change-in-production"),
		IntegrityChecksumDir:        getEnv("INTEGRITY_CHECKSUM_DIR", ".checksums"),
		IntegrityVerifyOnAccess:     integrityVerifyOnAccess,
		IntegrityScheduleChecks:     integrityScheduleChecks,
		IntegrityCheckInterval:      integrityCheckInterval,
	}

	cfg.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslMode)

	// Log the DSN with password masked for debugging
	maskedDSN := fmt.Sprintf("host=%s user=%s password=*** dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBSslMode)
	log.Printf("Database DSN: %s", maskedDSN)

	redisDBStr := getEnv("REDIS_DB", "0")
	var redisDB int
	parsedRedisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Printf("Warning: Could not parse REDIS_DB value '%s'. Using default 0. Error: %v\n", redisDBStr, err)
		redisDB = 0 // Default value
	} else {
		redisDB = parsedRedisDB
	}
	cfg.RedisDB = redisDB

	// Production warnings
	if cfg.AppEnv == "production" {
		if cfg.JWTSecret == "supersecretkey" || cfg.JWTSecret == "your_very_secret_key_for_jwt_change_this_please" || cfg.JWTSecret == "your_very_secret_key_for_jwt_change_this" {
			log.Println("CRITICAL WARNING: JWT_SECRET is set to a default/example value in a production environment. This is insecure. Please set a strong, unique JWT_SECRET environment variable.")
		}
		if cfg.InternalAPIKey == "default-internal-key-please-change" {
			log.Println("CRITICAL WARNING: INTERNAL_API_KEY is set to a default value in a production environment. This is insecure. Please set a strong, unique INTERNAL_API_KEY environment variable.")
		}
	}
	// Log a few key values to help debug in ECS logs
	log.Printf("Config Loaded: APP_ENV=%s, DB_HOST=%s, REDIS_ADDR=%s", cfg.AppEnv, cfg.DBHost, cfg.RedisAddr)

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Environment variable %s found with value: %s", key, value)
		return value
	}
	log.Printf("Environment variable %s not set, using default value: %s", key, defaultValue)
	return defaultValue
}
