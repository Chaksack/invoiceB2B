package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TerraformTestConfig holds configuration for Terraform tests
type TerraformTestConfig struct {
	ProjectRoot string
	TFFilePath  string
	VarFiles    []string
}

// NewTerraformTestConfig creates a new TerraformTestConfig
func NewTerraformTestConfig(projectRoot, tfFilePath string, varFiles ...string) *TerraformTestConfig {
	return &TerraformTestConfig{
		ProjectRoot: projectRoot,
		TFFilePath:  tfFilePath,
		VarFiles:    varFiles,
	}
}

// TerraformResource represents a Terraform resource
type TerraformResource struct {
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Instances []map[string]interface{} `json:"instances"`
}

// TerraformPlan represents a Terraform plan
type TerraformPlan struct {
	ResourceChanges []struct {
		Address      string `json:"address"`
		Mode         string `json:"mode"`
		Type         string `json:"type"`
		Name         string `json:"name"`
		ProviderName string `json:"provider_name"`
		Change       struct {
			Actions []string               `json:"actions"`
			Before  interface{}            `json:"before"`
			After   map[string]interface{} `json:"after"`
		} `json:"change"`
	} `json:"resource_changes"`
}

// RunTerraformInit initializes Terraform in the specified directory
func RunTerraformInit(t *testing.T, config *TerraformTestConfig) error {
	cmd := exec.Command("terraform", "init", "-backend=false")
	cmd.Dir = config.ProjectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Terraform init failed: %s", output)
		return fmt.Errorf("terraform init failed: %w", err)
	}
	return nil
}

// RunTerraformPlan runs terraform plan and returns the plan output
func RunTerraformPlan(t *testing.T, config *TerraformTestConfig) (*TerraformPlan, error) {
	// Create a temporary directory for the plan file
	tmpDir, err := ioutil.TempDir("", "terraform-test")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	planFile := filepath.Join(tmpDir, "tfplan.json")

	// Build the terraform plan command
	args := []string{
		"plan",
		"-no-color",
		fmt.Sprintf("-target=%s", config.TFFilePath),
		"-out=tfplan",
	}

	// Add var files if provided
	for _, varFile := range config.VarFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", varFile))
	}

	// Run terraform plan
	cmd := exec.Command("terraform", args...)
	cmd.Dir = config.ProjectRoot
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Terraform plan failed: %s", output)
		return nil, fmt.Errorf("terraform plan failed: %w", err)
	}

	// Convert the plan to JSON
	cmd = exec.Command("terraform", "show", "-json", "tfplan")
	cmd.Dir = config.ProjectRoot
	output, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("terraform show failed: %w", err)
	}

	// Write the JSON output to a file
	if err := ioutil.WriteFile(planFile, output, 0644); err != nil {
		return nil, fmt.Errorf("failed to write plan file: %w", err)
	}

	// Parse the JSON plan
	planData, err := ioutil.ReadFile(planFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	var plan TerraformPlan
	if err := json.Unmarshal(planData, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w", err)
	}

	return &plan, nil
}

// ValidateTerraformResource validates a Terraform resource against expected values
func ValidateTerraformResource(t *testing.T, plan *TerraformPlan, resourceType, resourceName string, expectedValues map[string]interface{}) {
	t.Helper()

	// Find the resource in the plan
	var resourceChange *struct {
		Address      string `json:"address"`
		Mode         string `json:"mode"`
		Type         string `json:"type"`
		Name         string `json:"name"`
		ProviderName string `json:"provider_name"`
		Change       struct {
			Actions []string               `json:"actions"`
			Before  interface{}            `json:"before"`
			After   map[string]interface{} `json:"after"`
		} `json:"change"`
	}

	for i, rc := range plan.ResourceChanges {
		if rc.Type == resourceType && rc.Name == resourceName {
			resourceChange = &plan.ResourceChanges[i]
			break
		}
	}

	assert.NotNil(t, resourceChange, "Resource %s.%s not found in plan", resourceType, resourceName)
	if resourceChange == nil {
		return
	}

	// Validate expected values
	for key, expectedValue := range expectedValues {
		// Handle nested keys with dot notation (e.g., "tags.Name")
		keys := strings.Split(key, ".")
		actualValue := getNestedValue(resourceChange.Change.After, keys)
		
		assert.Equal(t, expectedValue, actualValue, "Resource %s.%s has incorrect value for %s", resourceType, resourceName, key)
	}
}

// getNestedValue retrieves a nested value from a map using a slice of keys
func getNestedValue(data map[string]interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		return nil
	}

	if len(keys) == 1 {
		return data[keys[0]]
	}

	if nestedData, ok := data[keys[0]].(map[string]interface{}); ok {
		return getNestedValue(nestedData, keys[1:])
	}

	return nil
}

// SetupFiberTestApp creates a new Fiber app for testing
func SetupFiberTestApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	return app
}