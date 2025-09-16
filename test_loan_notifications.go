package main

import (
	"fmt"
	"log"
	"time"

	"invoiceB2B/internal/config"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"invoiceB2B/internal/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("=== Testing Loan Creation Notifications ===")
	
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize services
	emailService := services.NewEmailService(cfg)
	notificationService, err := services.NewNotificationService(cfg)
	if err != nil {
		log.Printf("Failed to initialize notification service: %v", err)
	}

	// Initialize repositories and services
	loanRepo := repositories.NewLoanApplicationRepository(db)
	
	// For testing, we'll create a mock loan validation service
	validationService := &MockLoanValidationService{}
	
	// Initialize loan application service
	loanService := services.NewLoanApplicationService(
		db,
		loanRepo,
		validationService,
		nil, // n8nService not needed for this test
		nil, // fileService not needed for this test
		notificationService,
		emailService,
	)

	// Test 1: Create a test user if not exists
	fmt.Println("\n1. Setting up test user...")
	testUser := &models.User{
		Email:       "test@example.com",
		CompanyName: "Test Company Ltd",
		FirstName:   "John",
		LastName:    "Doe",
		IsActive:    true,
	}

	// Check if test user exists
	var existingUser models.User
	result := db.Where("email = ?", testUser.Email).First(&existingUser)
	if result.Error == gorm.ErrRecordNotFound {
		// Create test user
		if err := db.Create(testUser).Error; err != nil {
			log.Fatalf("Failed to create test user: %v", err)
		}
		fmt.Printf("Created test user with ID: %d\n", testUser.ID)
	} else {
		testUser = &existingUser
		fmt.Printf("Using existing test user with ID: %d\n", testUser.ID)
	}

	// Test 2: Create a loan application to trigger notifications
	fmt.Println("\n2. Creating test loan application...")
	loanRequest := &dtos.CreateLoanApplicationRequest{
		Source:          "manual",
		RequestedAmount: 50000.00,
		Currency:        "USD",
		Purpose:         "Working capital for business expansion",
	}

	loanApp, err := loanService.CreateLoanApplication(testUser.ID, loanRequest)
	if err != nil {
		log.Fatalf("Failed to create loan application: %v", err)
	}

	fmt.Printf("✅ Loan application created successfully!\n")
	fmt.Printf("   - Application ID: %s\n", loanApp.ApplicationReference)
	fmt.Printf("   - Amount: %.2f %s\n", loanApp.RequestedAmount, loanApp.Currency)
	fmt.Printf("   - Status: %s\n", loanApp.Status)
	fmt.Printf("   - User: %s (%s)\n", testUser.Email, testUser.CompanyName)

	// Test 3: Test manual notification service calls
	fmt.Println("\n3. Testing notification services directly...")
	
	// Test email service
	fmt.Println("Testing email service...")
	testEmailSubject := "Test Loan Application Notification"
	testEmailBody := fmt.Sprintf(`
	<html>
	<body>
		<h2>Test Notification</h2>
		<p>This is a test notification for loan application %s</p>
		<p><strong>Amount:</strong> %.2f %s</p>
		<p><strong>User:</strong> %s</p>
	</body>
	</html>`, loanApp.ApplicationReference, loanApp.RequestedAmount, loanApp.Currency, testUser.Email)

	if err := emailService.SendEmail("admin@profundr.io", testEmailSubject, testEmailBody); err != nil {
		fmt.Printf("❌ Email test failed: %v\n", err)
	} else {
		fmt.Printf("✅ Email notification sent successfully\n")
	}

	// Test Slack service
	if notificationService != nil {
		fmt.Println("Testing Slack notification service...")
		loanAmount := fmt.Sprintf("%.2f %s", loanApp.RequestedAmount, loanApp.Currency)
		createdAt := loanApp.CreatedAt.Format("2006-01-02 15:04:05")
		if err := notificationService.SendLoanCreationAlert(
			testUser.Email,
			testUser.CompanyName,
			loanAmount,
			loanApp.ApplicationReference,
			loanApp.Purpose,
			string(loanApp.Source),
			string(loanApp.Status),
			createdAt,
		); err != nil {
			fmt.Printf("❌ Slack notification failed: %v\n", err)
		} else {
			fmt.Printf("✅ Slack notification sent successfully\n")
		}

		// Test generic Slack message
		testMessage := fmt.Sprintf("Test notification: New loan application %s created for %s", loanApp.ApplicationReference, testUser.Email)
		if err := notificationService.SendSlackNotification(testMessage); err != nil {
			fmt.Printf("❌ Generic Slack message failed: %v\n", err)
		} else {
			fmt.Printf("✅ Generic Slack message sent successfully\n")
		}
	} else {
		fmt.Printf("⚠️ Notification service not available (RabbitMQ connection issue)\n")
	}

	// Wait a moment for async notifications to process
	fmt.Println("\n4. Waiting for async notifications to process...")
	time.Sleep(5 * time.Second)

	fmt.Println("\n=== Test Complete ===")
	fmt.Println("Check your email and Slack channel for notifications!")
	fmt.Println("\nConfiguration used:")
	fmt.Printf("- SMTP Host: %s\n", cfg.SMTPHost)
	fmt.Printf("- Slack Webhook URL configured: %t\n", cfg.SlackWebhookURL != "")
	if cfg.SlackWebhookURL == "" {
		fmt.Println("⚠️ Note: SLACK_WEBHOOK_URL is not configured. Update your .env file with a valid webhook URL.")
	}
}

// MockLoanValidationService for testing
type MockLoanValidationService struct{}

func (m *MockLoanValidationService) ValidateCreateLoanRequest(request *dtos.CreateLoanApplicationRequest) error {
	// Simple validation for testing
	if request.RequestedAmount <= 0 {
		return fmt.Errorf("requested amount must be greater than 0")
	}
	if request.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	return nil
}

func (m *MockLoanValidationService) ValidateManualLoanRequest(request *dtos.ManualLoanInputRequest) error {
	// Simple validation for testing
	if request.RequestedAmount <= 0 {
		return fmt.Errorf("requested amount must be greater than 0")
	}
	if request.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	return nil
}

func (m *MockLoanValidationService) ValidateStatusTransition(currentStatus models.LoanApplicationStatus, newStatus models.LoanApplicationStatus) error {
	// For testing, allow all transitions
	return nil
}

func (m *MockLoanValidationService) ValidateKYBInformation(request *dtos.SubmitKYBInformationRequest) error {
	// For testing, skip validation
	return nil
}

func (m *MockLoanValidationService) ValidateFinancialStatement(request *dtos.SubmitFinancialStatementRequest) error {
	// For testing, skip validation
	return nil
}