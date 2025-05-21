package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"invoiceB2B/internal/config"
	"invoiceB2B/internal/database"
	"invoiceB2B/internal/handlers"
	"invoiceB2B/internal/middleware"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"invoiceB2B/internal/routes"
	"invoiceB2B/internal/services"
	"invoiceB2B/internal/utils"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize Validator
	customValidator := utils.NewCustomValidator()

	// Initialize Database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully.")

	// Auto-migrate schema
	log.Println("Running database migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Staff{},
		&models.KYCDetail{},
		&models.Invoice{},
		&models.Transaction{},
		&models.ActivityLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations completed.")

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully.")

	// Initialize Notification Service (RabbitMQ)
	notificationService, err := services.NewNotificationService(cfg)
	if err != nil {
		// Depending on how critical RabbitMQ is at startup,
		// you might log a warning and continue, or terminate.
		log.Printf("Warning: Failed to initialize Notification Service (RabbitMQ): %v. Some event notifications might not work.", err)
	} else {
		log.Println("Notification Service (RabbitMQ) initialized.")
		defer notificationService.Close()
	}

	// Initialize Services
	jwtService := services.NewJWTService(cfg)
	emailService := services.NewEmailService(cfg)
	otpService := services.NewOTPService(rdb, cfg.OTPExpirationMinutes)

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(db)
	kycRepo := repositories.NewKYCRepository(db)
	// staffRepo := repositories.NewStaffRepository(db)
	// activityLogRepo := repositories.NewActivityLogRepository(db)

	authService := services.NewAuthService(userRepo, kycRepo, jwtService, emailService, otpService, notificationService, cfg)
	userService := services.NewUserService(userRepo, kycRepo)
	// paymentService := services.NewPaymentService() // Placeholder

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService, customValidator.Validator)
	userHandler := handlers.NewUserHandler(userService, customValidator.Validator)
	// paymentHandler := handlers.NewPaymentHandler(paymentService, customValidator.Validator)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.GlobalErrorHandler,
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}))

	// JWT Middleware instance
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup Routes
	apiGroup := app.Group("/api/v1")
	routes.SetupAuthRoutes(apiGroup, authHandler, authMiddleware)
	routes.SetupUserRoutes(apiGroup, userHandler, authMiddleware)
	// routes.SetupPaymentRoutes(apiGroup, paymentHandler, authMiddleware)
	// routes.SetupAdminRoutes(apiGroup, ...)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "Welcome to the Invoice Financing API!",
			"version":   "1.2.0", // Incremented version for RabbitMQ integration
			"timestamp": time.Now(),
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Route '%s' not found on method %s.", c.Path(), c.Method()),
		})
	})

	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting server on port %s in %s mode", port, cfg.AppEnv)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
