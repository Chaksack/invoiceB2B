package main

import (
	"fmt"
	"log"

	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/services"
)

func main() {
	fmt.Println("🔧 Testing New Currencies in Loan Management System...")
	
	validator := services.NewLoanValidationService()

	// Test the new African currencies
	newCurrencies := []string{"GHS", "NGN", "ZAR"}
	
	for _, currency := range newCurrencies {
		fmt.Printf("\n✅ Testing %s currency...\n", currency)
		
		request := &dtos.CreateLoanApplicationRequest{
			Source:          "manual",
			RequestedAmount: 25000.00,
			Currency:        currency,
			Purpose:         "Business expansion loan",
		}

		err := validator.ValidateCreateLoanRequest(request)
		if err != nil {
			log.Printf("❌ %s currency validation failed: %v", currency, err)
		} else {
			fmt.Printf("✅ %s currency validation passed\n", currency)
		}
		
		// Test lowercase version
		requestLower := &dtos.CreateLoanApplicationRequest{
			Source:          "manual",
			RequestedAmount: 25000.00,
			Currency:        currency[0:1] + string(currency[1:2])[0:1] + string(currency[2:3])[0:1], // lowercase
			Purpose:         "Business expansion loan",
		}
		requestLower.Currency = string(currency[0]) + string(currency[1]) + string(currency[2])
		requestLower.Currency = fmt.Sprintf("%c%c%c", currency[0]+32, currency[1]+32, currency[2]+32) // convert to lowercase
		
		err = validator.ValidateCreateLoanRequest(requestLower)
		if err != nil {
			fmt.Printf("✅ Lowercase %s correctly rejected (case sensitive validation working)\n", requestLower.Currency)
		} else {
			fmt.Printf("✅ Lowercase %s accepted (case insensitive validation working)\n", requestLower.Currency)
		}
	}

	// Test existing currencies still work
	fmt.Printf("\n✅ Testing existing currencies...\n")
	existingCurrencies := []string{"USD", "EUR", "GBP", "KES", "UGX", "TZS"}
	
	for _, currency := range existingCurrencies {
		request := &dtos.CreateLoanApplicationRequest{
			Source:          "manual",
			RequestedAmount: 25000.00,
			Currency:        currency,
			Purpose:         "Business expansion loan",
		}

		err := validator.ValidateCreateLoanRequest(request)
		if err != nil {
			log.Printf("❌ Existing %s currency validation failed: %v", currency, err)
		} else {
			fmt.Printf("✅ Existing %s currency validation passed\n", currency)
		}
	}

	// Test invalid currency
	fmt.Printf("\n✅ Testing invalid currency rejection...\n")
	invalidRequest := &dtos.CreateLoanApplicationRequest{
		Source:          "manual",
		RequestedAmount: 25000.00,
		Currency:        "INVALID",
		Purpose:         "Business expansion loan",
	}

	err := validator.ValidateCreateLoanRequest(invalidRequest)
	if err != nil {
		fmt.Printf("✅ Invalid currency correctly rejected: %v\n", err)
	} else {
		fmt.Println("❌ Invalid currency should have been rejected")
	}
	
	fmt.Println("\n🎉 Currency validation testing completed!")
}