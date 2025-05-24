package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

// LoadConfig loads configuration from .env file or environment variables
func LoadConfig(path string) (*Config, error) {
	currentAppEnv := os.Getenv("APP_ENV")
	if currentAppEnv == "" {
		currentAppEnv = "development"
	}

	if currentAppEnv != "production" {
		envPath := path
		if envPath == "." || envPath == "" {
			envPath = ".env"
		} else {
			envPath = path + "/.env"
		}
		err := godotenv.Overload(envPath)
		if err != nil {
			if currentAppEnv == "development" {
				fmt.Printf("Warning: .env file not found at %s. Using environment variables or defaults. Error: %v\n", envPath, err)
			}
		}

		currentAppEnv = getEnv("APP_ENV", "development")
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	accessTokenExpMinutes, _ := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_EXPIRATION_MINUTES", "15"))
	refreshTokenExpDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_EXPIRATION_DAYS", "7"))
	otpExpMinutes, _ := strconv.Atoi(getEnv("OTP_EXPIRATION_MINUTES", "5"))
	maxUploadSizeMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)

	// Construct RedisAddr from REDIS_HOST and REDIS_PORT
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	calculatedRedisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "3000"),
		AppEnv:     currentAppEnv,
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"), // Ensure these defaults match your actual fallback needs
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "invoice_db"),
		DBSslMode:  getEnv("DB_SSLMODE", "disable"), // Corrected typo: removed leading 'A'

		RedisAddr:     calculatedRedisAddr, // Use the constructed RedisAddr
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

		SMTPHost:        getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:        smtpPort,
		SMTPUser:        getEnv("SMTP_USER", ""),
		SMTPPassword:    getEnv("SMTP_PASSWORD", ""),
		SMTPSenderEmail: getEnv("SMTPSenderEmail", "Invoice App <no-reply@example.com>"), // Note: Key is SMTPSenderEmail, not SMTP_SENDER_EMAIL

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
		fmt.Printf("Warning: Could not parse REDIS_DB value '%s'. Using default 0. Error: %v\n", redisDBStr, err)
		redisDB = 0
	} else {
		redisDB = parsedRedisDB
	}
	cfg.RedisDB = redisDB

	if cfg.JWTSecret == "supersecretkey" || cfg.JWTSecret == "your_very_secret_key_for_jwt_change_this_please" || cfg.JWTSecret == "your_very_secret_key_for_jwt_change_this" {
		fmt.Println("WARNING: JWT_SECRET is set to a default/example value. Please change this for production!")
	}
	if cfg.InternalAPIKey == "default-internal-key-please-change" {
		fmt.Println("WARNING: INTERNAL_API_KEY is set to a default value. Please change this for production!")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
