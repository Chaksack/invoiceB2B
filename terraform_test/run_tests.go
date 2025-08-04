package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Parse command line flags
	var (
		testDir   = flag.String("dir", "", "Directory to run tests from (default: all)")
		testName  = flag.String("test", "", "Specific test to run (e.g., TestVPCConfiguration)")
		verbose   = flag.Bool("v", false, "Verbose output")
		coverage  = flag.Bool("coverage", false, "Generate coverage report")
		coverFile = flag.String("coverfile", "coverage.out", "Coverage file name")
	)
	flag.Parse()

	// Get the project root directory
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Determine which directories to test
	var testDirs []string
	if *testDir != "" {
		// Test a specific directory
		testDirs = []string{*testDir}
	} else {
		// Test all directories
		testDirs = []string{
			"infrastructure",
			"application",
			"monitoring",
		}
	}

	// Build the test command
	var testCmd []string
	testCmd = append(testCmd, "go", "test")

	// Add verbose flag if requested
	if *verbose {
		testCmd = append(testCmd, "-v")
	}

	// Add coverage flags if requested
	if *coverage {
		testCmd = append(testCmd, "-coverprofile="+*coverFile)
	}

	// Add specific test if requested
	if *testName != "" {
		testCmd = append(testCmd, "-run="+*testName)
	}

	// Run tests in each directory
	var failedTests []string
	for _, dir := range testDirs {
		testPath := filepath.Join(projectRoot, dir)
		fmt.Printf("Running tests in %s...\n", testPath)

		// Create the command
		cmd := exec.Command(testCmd[0], testCmd[1:]...)
		cmd.Dir = testPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		startTime := time.Now()
		err := cmd.Run()
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("Tests in %s failed after %s: %v\n", dir, duration, err)
			failedTests = append(failedTests, dir)
		} else {
			fmt.Printf("Tests in %s passed in %s\n", dir, duration)
		}
	}

	// Generate coverage report if requested
	if *coverage {
		fmt.Println("Generating coverage report...")
		cmd := exec.Command("go", "tool", "cover", "-html="+*coverFile)
		cmd.Dir = projectRoot
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error generating coverage report: %v\n", err)
		}
	}

	// Print summary
	fmt.Println("\nTest Summary:")
	if len(failedTests) > 0 {
		fmt.Printf("Failed tests in directories: %s\n", strings.Join(failedTests, ", "))
		os.Exit(1)
	} else {
		fmt.Println("All tests passed!")
	}
}