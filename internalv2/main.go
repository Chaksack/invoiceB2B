package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"invoiceB2B/internal@v2/config"
	"invoiceB2B/internal@v2/database"
	"invoiceB2B/internal@v2/handlers"
	"invoiceB2B/internal@v2/middleware"
	"invoiceB2B/internal@v2/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		AppName:      "InvoiceB2B Internal API v2",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "Internal API v2 is healthy",
			"version":   "2.0.0",
			"timestamp": fiber.Now(),
		})
	})

	// API index endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "InvoiceB2B Internal API v2",
			"version":   "2.0.0",
			"timestamp": fiber.Now(),
			"endpoints": fiber.Map{
				"auth":     "/api/auth",
				"business": "/api/business",
				"admin":    "/api/admin",
				"health":   "/health",
			},
		})
	})

	// Setup routes
	routes.SetupRoutes(app, db, cfg)

	// Start server
	port := os.Getenv("INTERNAL_API_PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("Starting Internal API v2 server on port %s", port)
	log.Fatal(app.Listen(":" + port))
} 