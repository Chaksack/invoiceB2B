package application

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"invoiceB2B/terraform_test/utils"
)

func TestECSConfiguration(t *testing.T) {
	// Setup Fiber app for testing
	app := utils.SetupFiberTestApp()

	// Create a test endpoint to validate ECS configuration
	app.Get("/test/ecs", func(c *fiber.Ctx) error {
		// Get the project root directory
		projectRoot, err := os.Getwd()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get working directory",
			})
		}
		// Go up two levels to get to the project root
		projectRoot = filepath.Dir(filepath.Dir(projectRoot))

		// Create a test configuration
		config := utils.NewTerraformTestConfig(
			projectRoot,
			"ecs.tf",
			"variables.tf",
		)

		// Initialize Terraform
		if err := utils.RunTerraformInit(t, config); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to initialize Terraform",
			})
		}

		// Run Terraform plan
		plan, err := utils.RunTerraformPlan(t, config)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to run Terraform plan",
			})
		}

		// Validate ECS cluster configuration
		utils.ValidateTerraformResource(t, plan, "aws_ecs_cluster", "main", map[string]interface{}{
			"name":                 "${var.ecs_cluster_name}",
			"setting.0.name":       "containerInsights",
			"setting.0.value":      "enabled",
			"tags.Name":            "${var.ecs_cluster_name}",
			"tags.Project":         "${var.project_name}",
			"tags.Environment":     "production",
		})

		// Validate CloudWatch Log Group configuration
		utils.ValidateTerraformResource(t, plan, "aws_cloudwatch_log_group", "ecs_logs", map[string]interface{}{
			"name":              "/ecs/${var.project_name}",
			"retention_in_days": 30,
			"tags.Name":         "${var.project_name}-ecs-logs",
			"tags.Project":      "${var.project_name}",
			"tags.Environment":  "production",
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "ECS configuration validated successfully",
		})
	})

	// Test the endpoint
	resp, err := app.Test(httptest.NewRequest("GET", "/test/ecs", nil))
	require.NoError(t, err)

	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestECSResourceCreation(t *testing.T) {
	// Create a mock AWS service
	mockAWS := utils.NewMockAWSService()

	// Create mock ECS cluster
	ecsCluster := utils.MockECSCluster{
		ID:   "ecs-cluster-12345",
		Name: "invoice-b2b-cluster",
		Tags: map[string]string{
			"Name":        "invoice-b2b-cluster",
			"Project":     "invoice",
			"Environment": "production",
		},
	}

	// Add ECS cluster to mock AWS service
	mockAWS.AddResource("aws_ecs_cluster", ecsCluster.ID, map[string]interface{}{
		"id":   ecsCluster.ID,
		"name": ecsCluster.Name,
		"setting": []map[string]interface{}{
			{
				"name":  "containerInsights",
				"value": "enabled",
			},
		},
		"tags": ecsCluster.Tags,
	})

	// Create mock CloudWatch Log Group
	logGroup := map[string]interface{}{
		"id":                "log-group-12345",
		"name":              "/ecs/invoice",
		"retention_in_days": 30,
		"tags": map[string]string{
			"Name":        "invoice-ecs-logs",
			"Project":     "invoice",
			"Environment": "production",
		},
	}

	// Add CloudWatch Log Group to mock AWS service
	mockAWS.AddResource("aws_cloudwatch_log_group", logGroup["id"].(string), logGroup)

	// Test retrieving ECS cluster
	ecsClusterResource, err := mockAWS.GetResource("aws_ecs_cluster", ecsCluster.ID)
	require.NoError(t, err)
	assert.Equal(t, ecsCluster.Name, ecsClusterResource["name"])
	assert.Equal(t, ecsCluster.Tags, ecsClusterResource["tags"])

	// Test retrieving CloudWatch Log Group
	logGroupResource, err := mockAWS.GetResource("aws_cloudwatch_log_group", logGroup["id"].(string))
	require.NoError(t, err)
	assert.Equal(t, logGroup["name"], logGroupResource["name"])
	assert.Equal(t, logGroup["retention_in_days"], logGroupResource["retention_in_days"])
	assert.Equal(t, logGroup["tags"], logGroupResource["tags"])
}