package main

import (
	"context"
	"fmt"
	"invoiceB2B/internal/dtos"
	"net/http"
	"os"
	"path/filepath" // Ensure filepath is imported
	"strings"
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
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
)

type NuxtProjectConfig struct {
	Name     string
	URLPath  string
	DistPath string
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if _, err := os.Stat(cfg.UploadsDir); os.IsNotExist(err) {
		log.Infof("Uploads directory %s does not exist. Creating...", cfg.UploadsDir)
		if err := os.MkdirAll(cfg.UploadsDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create uploads directory: %v", err)
		}
		subDirs := []string{"invoices", "receipts", "kyc"}
		for _, dir := range subDirs {
			path := fmt.Sprintf("%s/%s", cfg.UploadsDir, dir)
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				log.Fatalf("Failed to create subdirectory %s: %v", path, err)
			}
		}
	}

	customValidator := utils.NewCustomValidator()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Info("Database connected successfully.")

	err = db.AutoMigrate(
		&models.User{}, &models.Staff{}, &models.KYCDetail{},
		&models.Invoice{}, &models.Transaction{}, &models.ActivityLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Info("Database migrations completed.")

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Info("Redis connected successfully.")

	notificationService, err := services.NewNotificationService(cfg)
	if err != nil {
		log.Infof("Warning: Failed to initialize Notification Service (RabbitMQ): %v. Some event notifications might not work.", err)
	} else {
		log.Info("Notification Service (RabbitMQ) initialized.")
		if notificationService != nil { // Ensure service is not nil before deferring Close
			defer notificationService.Close()
		}
	}

	jwtService := services.NewJWTService(cfg)
	emailService := services.NewEmailService(cfg)
	otpService := services.NewOTPService(rdb, cfg.OTPExpirationMinutes)
	fileService := services.NewFileService(cfg.UploadsDir, cfg.MaxUploadSizeMB*1024*1024)

	userRepo := repositories.NewUserRepository(db)
	kycRepo := repositories.NewKYCRepository(db)
	invoiceRepo := repositories.NewInvoiceRepository(db)
	staffRepo := repositories.NewStaffRepository(db)
	activityLogRepo := repositories.NewActivityLogRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	activityLogSvc := services.NewActivityLogService(activityLogRepo)
	authService := services.NewAuthService(userRepo, staffRepo, kycRepo, jwtService, emailService, otpService, notificationService, activityLogSvc, cfg)
	userService := services.NewUserService(userRepo, kycRepo, activityLogSvc)
	invoiceService := services.NewInvoiceService(invoiceRepo, userRepo, transactionRepo, fileService, notificationService, activityLogSvc, emailService, cfg)
	internalService := services.NewInternalService(invoiceRepo, activityLogSvc)

	// Initialize AdminService (pass all dependencies, including staffRepo)
	adminService := services.NewAdminService(
		userRepo,
		kycRepo,
		staffRepo,
		invoiceRepo,
		transactionRepo,
		activityLogSvc,
		emailService,
		nil,
		fileService,
		cfg,
	)

	// --- Create Superadmin User (if not exists) ---
	createSuperAdminIfNotExists(adminService, cfg)

	authHandler := handlers.NewAuthHandler(authService, customValidator.Validator)
	userHandler := handlers.NewUserHandler(userService, customValidator.Validator)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService, fileService, customValidator.Validator)
	adminHandler := handlers.NewAdminHandler(adminService, fileService, customValidator.Validator)
	internalHandler := handlers.NewInternalHandler(internalService, customValidator.Validator)

	app := fiber.New(fiber.Config{
		ErrorHandler: utils.GlobalErrorHandler,
	})

	nuxtProjects := []NuxtProjectConfig{
		{Name: "Dashboard", URLPath: "/", DistPath: "./client/dist"},
	}
	setupNuxtFrontendServers(app, nuxtProjects)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5000,http://localhost:3000,http://localhost:3001",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(flogger.New(flogger.Config{
		Format:     "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}))

	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	adminMiddleware := middleware.NewAdminMiddleware(staffRepo)
	internalApiMiddleware := middleware.NewInternalAPIMiddleware(cfg.InternalAPIKey)

	apiV1 := app.Group("/api/v1")
	routes.SetupAuthRoutes(apiV1, authHandler, authMiddleware)
	routes.SetupUserRoutes(apiV1, userHandler, authMiddleware)
	routes.SetupInvoiceRoutes(apiV1, invoiceHandler, authMiddleware, adminMiddleware)
	routes.SetupAdminRoutes(apiV1, adminHandler, authMiddleware, adminMiddleware)
	routes.SetupInternalRoutes(apiV1, internalHandler, internalApiMiddleware)

	app.Get("/api/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "Welcome to the Invoice Financing API!",
			"version":   "1.2.0",
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
	log.Infof("Starting server on port %s in %s mode", port, cfg.AppEnv)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func setupNuxtFrontendServers(app *fiber.App, projects []NuxtProjectConfig) {
	for _, project := range projects {
		relativeDistPath := project.DistPath // e.g., "./client/dist"

		absDistPath, err := filepath.Abs(relativeDistPath)
		if err != nil {
			log.Errorf("Nuxt ('%s'): Error getting absolute path for DistPath '%s': %v", project.Name, relativeDistPath, err)
			continue
		}
		log.Infof("Nuxt ('%s'): project.URLPath='%s', relativeDistPath='%s', absoluteDistPathForStat='%s'",
			project.Name, project.URLPath, relativeDistPath, absDistPath)

		if _, err := os.Stat(absDistPath); os.IsNotExist(err) {
			log.Warnf("Nuxt ('%s'): Distribution DIRECTORY NOT FOUND at absolute path: '%s' (from relative '%s'). This Nuxt app will not be served.",
				project.Name, absDistPath, relativeDistPath)
			continue
		} else if err != nil {
			log.Warnf("Nuxt ('%s'): Error stating distribution DIRECTORY at absolute path: '%s': %v. This Nuxt app will not be served.",
				project.Name, absDistPath, err)
			continue
		}
		log.Infof("Nuxt ('%s'): Distribution DIRECTORY found at absolute path: '%s'", project.Name, absDistPath)

		// Specific handler for "/" if project.URLPath is "/"
		if project.URLPath == "/" {
			app.Get("/", func(c *fiber.Ctx) error {
				// Construct absolute path to index.html
				absIndexPath := filepath.Join(absDistPath, "index.html")
				log.Infof("Nuxt ('%s') - GET / handler: Checking for index.html at absolute path '%s'", project.Name, absIndexPath)

				if _, err := os.Stat(absIndexPath); os.IsNotExist(err) {
					log.Errorf("Nuxt ('%s') - GET / handler: index.html NOT FOUND by os.Stat at absolute path '%s'", project.Name, absIndexPath)
					return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("index.html not found by debug handler (abs path: %s)", absIndexPath))
				} else if err != nil {
					log.Errorf("Nuxt ('%s') - GET / handler: Error stating index.html at absolute path '%s': %v", project.Name, absIndexPath, err)
					return c.Status(fiber.StatusInternalServerError).SendString("Error checking for index.html")
				}
				log.Infof("Nuxt ('%s') - GET / handler: Attempting to serve index.html from absolute path '%s'", project.Name, absIndexPath)
				return c.SendFile(absIndexPath)
			})
		}

		// General static file serving
		app.Use(project.URLPath, filesystem.New(filesystem.Config{
			Root:         http.Dir(absDistPath), // http.Dir needs a valid filesystem path
			Browse:       false,
			Index:        "index.html",                             // Relative to Root
			NotFoundFile: filepath.Join(absDistPath, "index.html"), // Absolute path for robustness
		}))
		log.Infof("Nuxt ('%s'): Serving static files from URLPath '%s' using filesystem root '%s'",
			project.Name, project.URLPath, absDistPath)

	}
}

// createSuperAdminIfNotExists function
func createSuperAdminIfNotExists(adminService services.AdminService, cfg *config.Config) {
	superAdminEmail := os.Getenv("SUPERADMIN_EMAIL")
	superAdminPassword := os.Getenv("SUPERADMIN_PASSWORD")
	superAdminFirstName := os.Getenv("SUPERADMIN_FIRSTNAME")
	superAdminLastName := os.Getenv("SUPERADMIN_LASTNAME")
	superAdminRole := "admin"

	if superAdminEmail == "" || superAdminPassword == "" {
		log.Info("Superadmin credentials (SUPERADMIN_EMAIL, SUPERADMIN_PASSWORD) not found in environment variables. Skipping superadmin creation.")
		return
	}
	if superAdminFirstName == "" {
		superAdminFirstName = "Super" // Default value
	}
	if superAdminLastName == "" {
		superAdminLastName = "Admin" // Default value
	}

	// The CreateStaff service method already checks if the email exists.
	// It returns an error if the staff member already exists.
	superAdminReq := dtos.CreateStaffRequest{
		Email:     superAdminEmail,
		Password:  superAdminPassword, // This will be hashed by CreateStaff
		FirstName: superAdminFirstName,
		LastName:  superAdminLastName,
		Role:      superAdminRole,
	}

	log.Infof("Attempting to create/verify superadmin: %s", superAdminEmail)
	_, err := adminService.CreateStaff(context.Background(), superAdminReq)
	if err != nil {
		if strings.Contains(err.Error(), "staff with this email already exists") {
			log.Infof("Superadmin with email %s already exists. No action taken.", superAdminEmail)
		} else {
			log.Infof("Failed to create superadmin %s: %v", superAdminEmail, err)
		}
	} else {
		log.Infof("Superadmin %s created successfully with role %s.", superAdminEmail, superAdminRole)
	}
}
