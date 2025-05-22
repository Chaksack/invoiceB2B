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
}

// LoadConfig loads configuration from .env file or environment variables
func LoadConfig(path string) (*Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(path + "/.env")
		if err != nil && os.Getenv("APP_ENV") == "development" {
			fmt.Printf("Warning: .env file not found at %s/.env. Using environment variables or defaults.\n", path)
		}
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	accessTokenExpMinutes, _ := strconv.Atoi(getEnv("JWT_ACCESS_TOKEN_EXPIRATION_MINUTES", "15"))
	refreshTokenExpDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_TOKEN_EXPIRATION_DAYS", "7"))
	otpExpMinutes, _ := strconv.Atoi(getEnv("OTP_EXPIRATION_MINUTES", "5"))
	maxUploadSizeMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)

	cfg := &Config{
		AppPort:                                getEnv("APP_PORT", "3000"),
		AppEnv:                                 getEnv("APP_ENV", "development"),
		DBHost:                                 getEnv("DB_HOST", "localhost"),
		DBPort:                                 getEnv("DB_PORT", "5432"),
		DBUser:                                 getEnv("DB_USER", "user"),
		DBPassword:                             getEnv("DB_PASSWORD", "password"),
		DBName:                                 getEnv("DB_NAME", "invoice_db"),
		DBSslMode:                              getEnv("DB_SSLMODE", "disable"),
		RedisAddr:                              getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:                          getEnv("REDIS_PASSWORD", ""),
		RabbitMQURL:                            getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQEventExchangeName:              getEnv("RABBITMQ_EVENT_EXCHANGE_NAME", "invoice_events_exchange"), // Updated exchange name
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
		SMTPSenderEmail: getEnv("SMTPSenderEmail", "Invoice App <no-reply@example.com>"),

		OTPExpirationMinutes: time.Duration(otpExpMinutes) * time.Minute,
		UploadsDir:           getEnv("UPLOADS_DIR", "./uploads"),
		MaxUploadSizeMB:      maxUploadSizeMB,
	}

	cfg.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSslMode)

	redisDBStr := getEnv("REDIS_DB", "0")
	var redisDB int
	_, err := fmt.Sscan(redisDBStr, &redisDB)
	if err != nil {
		fmt.Printf("Warning: Could not parse REDIS_DB value '%s'. Using default 0. Error: %v\n", redisDBStr, err)
		redisDB = 0
	}
	cfg.RedisDB = redisDB

	if cfg.JWTSecret == "supersecretkey" || cfg.JWTSecret == "your_very_secret_key_for_jwt_change_this_please" {
		fmt.Println("WARNING: JWT_SECRET is set to a default/example value. Please change this in your .env file for production!")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
