package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	JWTAccessTokenExpirationMinutes time.Duration
	JWTRefreshTokenExpirationDays   time.Duration

	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPassword    string
	SMTPSenderEmail string

	OTPExpirationMinutes time.Duration
	UploadsDir           string
	MaxUploadSizeMB      int64
	InternalAPIKey       string
}

func LoadConfig(path string) (*Config, error) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
		log.Println("Warning: APP_ENV not set, defaulting to 'development'.")
	}

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

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	calculatedRedisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

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
		JWTAccessTokenExpirationMinutes: time.Duration(accessTokenExpMinutes) * time.Minute,
		JWTRefreshTokenExpirationDays:   time.Duration(refreshTokenExpDays) * 24 * time.Hour,

		SMTPHost:        getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:        smtpPort,
		SMTPUser:        getEnv("SMTP_USER", "your_email@example.com"),
		SMTPPassword:    getEnv("SMTP_PASSWORD", "your_app_password"),
		SMTPSenderEmail: getEnv("SMTPSenderEmail", "Profundr <no-reply@profundr.io>"),

		OTPExpirationMinutes: time.Duration(otpExpMinutes) * time.Minute,
		UploadsDir:           getEnv("UPLOADS_DIR", "./uploads"),
		MaxUploadSizeMB:      maxUploadSizeMB,
		InternalAPIKey:       getEnv("INTERNAL_API_KEY", "default-internal-key-please-change"),
	}

	cfg.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslMode)

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
		return value
	}
	log.Printf("Environment variable %s not set, using default value: %s", key, defaultValue)
	return defaultValue
}
