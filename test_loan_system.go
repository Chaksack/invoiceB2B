package main

import (
	"fmt"
	"log"

	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/services"
)

func main() {
	fmt.Println("🔧 Testing Loan Management System Validation...")
	
	// Test validation service only
	testValidationService()
	
	fmt.Println("\n🎉 Loan Management System validation testing completed!")
	fmt.Println("📋 Repository and integration tests require database setup")
}


func testValidationService() {
	fmt.Println("\n✅ Testing Validation Service...")
	
	validator := services.NewLoanValidationService()

	// Test valid loan application request
	validRequest := &dtos.CreateLoanApplicationRequest{
		Source:          "manual",
		RequestedAmount: 25000.00,
		Currency:        "USD",
		Purpose:         "Business expansion loan",
	}

	err := validator.ValidateCreateLoanRequest(validRequest)
	if err != nil {
		log.Printf("❌ Valid request validation failed: %v", err)
	} else {
		fmt.Println("✅ Valid request validation passed")
	}

	// Test invalid loan amount (too small)
	invalidRequest := &dtos.CreateLoanApplicationRequest{
		Source:          "manual",
		RequestedAmount: 500.00, // Below minimum
		Currency:        "USD",
		Purpose:         "Test loan",
	}

	err = validator.ValidateCreateLoanRequest(invalidRequest)
	if err != nil {
		fmt.Printf("✅ Invalid amount validation correctly rejected: %v\n", err)
	} else {
		fmt.Println("❌ Invalid amount validation should have failed")
	}

	// Test invalid currency
	invalidCurrencyRequest := &dtos.CreateLoanApplicationRequest{
		Source:          "manual",
		RequestedAmount: 25000.00,
		Currency:        "INVALID", // Unsupported currency
		Purpose:         "Test loan",
	}

	err = validator.ValidateCreateLoanRequest(invalidCurrencyRequest)
	if err != nil {
		fmt.Printf("✅ Invalid currency validation correctly rejected: %v\n", err)
	} else {
		fmt.Println("❌ Invalid currency validation should have failed")
	}

	// Test status transition validation
	err = validator.ValidateStatusTransition(models.LoanApplicationPending, models.LoanApplicationUnderReview)
	if err != nil {
		log.Printf("❌ Valid status transition failed: %v", err)
	} else {
		fmt.Println("✅ Valid status transition validation passed")
	}

	err = validator.ValidateStatusTransition(models.LoanApplicationCompleted, models.LoanApplicationPending)
	if err != nil {
		fmt.Printf("✅ Invalid status transition correctly rejected: %v\n", err)
	} else {
		fmt.Println("❌ Invalid status transition should have failed")
	}

	// Test KYB validation
	validKYB := &dtos.SubmitKYBInformationRequest{
		BusinessName:            "Test Business Ltd",
		BusinessRegistrationNo:  "REG123456",
		BusinessType:            "Limited Company",
		IndustryType:            "Technology",
		BusinessAddress:         "123 Business Street, City",
		TaxIdentificationNumber: "TAX123456",
		YearsInOperation:        5,
		BusinessDescription:     "A technology consulting business",
	}

	err = validator.ValidateKYBInformation(validKYB)
	if err != nil {
		log.Printf("❌ Valid KYB validation failed: %v", err)
	} else {
		fmt.Println("✅ Valid KYB validation passed")
	}

	invalidKYB := &dtos.SubmitKYBInformationRequest{
		BusinessName:            "", // Required field missing
		BusinessRegistrationNo:  "REG123456",
		BusinessType:            "Limited Company",
		IndustryType:            "Technology",
		BusinessAddress:         "123 Business Street",
		TaxIdentificationNumber: "TAX123456",
		YearsInOperation:        0, // Below minimum
		BusinessDescription:     "Test business",
	}

	err = validator.ValidateKYBInformation(invalidKYB)
	if err != nil {
		fmt.Printf("✅ Invalid KYB validation correctly rejected: %v\n", err)
	} else {
		fmt.Println("❌ Invalid KYB validation should have failed")
	}
}
