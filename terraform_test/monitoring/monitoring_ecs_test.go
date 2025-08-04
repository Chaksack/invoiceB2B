package monitoring

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

func TestMonitoringECSConfiguration(t *testing.T) {
	// Setup Fiber app for testing
	app := utils.SetupFiberTestApp()

	// Create a test endpoint to validate monitoring ECS configuration
	app.Get("/test/monitoring-ecs", func(c *fiber.Ctx) error {
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
			"monitoring_ecs.tf",
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

		// Validate EFS File System for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_efs_file_system", "prometheus_data", map[string]interface{}{
			"creation_token": "${var.project_name}-prometheus-data-efs",
			"encrypted":      true,
			"tags.Name":      "${var.project_name}-prometheus-data",
			"tags.Project":   "${var.project_name}",
		})

		// Validate EFS Access Point for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_efs_access_point", "prometheus_data", map[string]interface{}{
			"file_system_id":           "${aws_efs_file_system.prometheus_data.id}",
			"posix_user.0.gid":         65534,
			"posix_user.0.uid":         65534,
			"root_directory.0.path":    "/prometheus-data",
			"tags.Name":                "${var.project_name}-prometheus-ap",
			"tags.Project":             "${var.project_name}",
		})

		// Validate ECR Repository for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_ecr_repository", "prometheus", map[string]interface{}{
			"name":                 "${var.project_name}-prometheus",
			"image_tag_mutability": "MUTABLE",
			"tags.Name":            "${var.project_name}-prometheus-ecr",
			"tags.Project":         "${var.project_name}",
		})

		// Validate Security Group Rule for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_security_group_rule", "prometheus_ingress", map[string]interface{}{
			"security_group_id":        "${aws_security_group.ecs_tasks.id}",
			"type":                     "ingress",
			"from_port":                9090,
			"to_port":                  9090,
			"protocol":                 "tcp",
			"source_security_group_id": "${aws_security_group.alb.id}",
			"description":              "Allow traffic from ALB to Prometheus",
		})

		// Validate ALB Target Group for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_lb_target_group", "prometheus", map[string]interface{}{
			"name":        "${var.project_name}-prometheus-tg",
			"port":        9090,
			"protocol":    "HTTP",
			"vpc_id":      "${aws_vpc.main.id}",
			"target_type": "ip",
			"tags.Name":   "${var.project_name}-prometheus-tg",
			"tags.Project": "${var.project_name}",
		})

		// Validate ALB Listener Rule for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_lb_listener_rule", "prometheus_rule", map[string]interface{}{
			"listener_arn":                "${aws_lb_listener.https.arn}",
			"priority":                    130,
			"action.0.type":               "forward",
			"action.0.target_group_arn":   "${aws_lb_target_group.prometheus.arn}",
			"condition.0.path_pattern.0.values.0": "/prometheus/*",
		})

		// Validate ECS Task Definition for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_ecs_task_definition", "prometheus", map[string]interface{}{
			"family":                   "${var.project_name}-prometheus",
			"network_mode":             "awsvpc",
			"requires_compatibilities.0": "FARGATE",
			"cpu":                      "1024",
			"memory":                   "2048",
			"execution_role_arn":       "${aws_iam_role.ecs_execution_role.arn}",
			"task_role_arn":            "${aws_iam_role.ecs_task_role.arn}",
			"tags.Name":                "${var.project_name}-prometheus-task",
			"tags.Project":             "${var.project_name}",
			"tags.Environment":         "production",
		})

		// Validate ECS Service for Prometheus
		utils.ValidateTerraformResource(t, plan, "aws_ecs_service", "prometheus", map[string]interface{}{
			"name":            "${var.project_name}-prometheus",
			"cluster":         "${aws_ecs_cluster.main.id}",
			"task_definition": "${aws_ecs_task_definition.prometheus.arn}",
			"desired_count":   1,
			"launch_type":     "FARGATE",
			"tags.Name":       "${var.project_name}-prometheus-service",
			"tags.Project":    "${var.project_name}",
			"tags.Environment": "production",
		})

		// Validate IAM Policy for EFS access
		utils.ValidateTerraformResource(t, plan, "aws_iam_policy", "monitoring_efs_access", map[string]interface{}{
			"name":        "${var.project_name}-monitoring-efs-access",
			"description": "Allow ECS tasks to access EFS for monitoring services",
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Monitoring ECS configuration validated successfully",
		})
	})

	// Test the endpoint
	resp, err := app.Test(httptest.NewRequest("GET", "/test/monitoring-ecs", nil))
	require.NoError(t, err)

	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestMonitoringECSResourceCreation(t *testing.T) {
	// Create a mock AWS service
	mockAWS := utils.NewMockAWSService()

	// Create mock EFS File System for Prometheus
	prometheusEFS := map[string]interface{}{
		"id":             "fs-prometheus-12345",
		"creation_token": "invoice-prometheus-data-efs",
		"encrypted":      true,
		"tags": map[string]string{
			"Name":    "invoice-prometheus-data",
			"Project": "invoice",
		},
	}

	// Add EFS File System to mock AWS service
	mockAWS.AddResource("aws_efs_file_system", prometheusEFS["id"].(string), prometheusEFS)

	// Create mock EFS Access Point for Prometheus
	prometheusAP := map[string]interface{}{
		"id":             "fsap-prometheus-12345",
		"file_system_id": prometheusEFS["id"].(string),
		"posix_user": map[string]interface{}{
			"gid": 65534,
			"uid": 65534,
		},
		"root_directory": map[string]interface{}{
			"path": "/prometheus-data",
			"creation_info": map[string]interface{}{
				"owner_gid":   65534,
				"owner_uid":   65534,
				"permissions": "755",
			},
		},
		"tags": map[string]string{
			"Name":    "invoice-prometheus-ap",
			"Project": "invoice",
		},
	}

	// Add EFS Access Point to mock AWS service
	mockAWS.AddResource("aws_efs_access_point", prometheusAP["id"].(string), prometheusAP)

	// Create mock ECR Repository for Prometheus
	prometheusECR := map[string]interface{}{
		"id":                   "ecr-prometheus-12345",
		"name":                 "invoice-prometheus",
		"image_tag_mutability": "MUTABLE",
		"tags": map[string]string{
			"Name":    "invoice-prometheus-ecr",
			"Project": "invoice",
		},
	}

	// Add ECR Repository to mock AWS service
	mockAWS.AddResource("aws_ecr_repository", prometheusECR["id"].(string), prometheusECR)

	// Create mock ECS Task Definition for Prometheus
	prometheusTaskDef := map[string]interface{}{
		"id":                      "task-def-prometheus-12345",
		"family":                  "invoice-prometheus",
		"network_mode":            "awsvpc",
		"requires_compatibilities": []string{"FARGATE"},
		"cpu":                     "1024",
		"memory":                  "2048",
		"tags": map[string]string{
			"Name":        "invoice-prometheus-task",
			"Project":     "invoice",
			"Environment": "production",
		},
	}

	// Add ECS Task Definition to mock AWS service
	mockAWS.AddResource("aws_ecs_task_definition", prometheusTaskDef["id"].(string), prometheusTaskDef)

	// Create mock ECS Service for Prometheus
	prometheusService := map[string]interface{}{
		"id":              "service-prometheus-12345",
		"name":            "invoice-prometheus",
		"cluster":         "ecs-cluster-12345",
		"task_definition": prometheusTaskDef["id"].(string),
		"desired_count":   1,
		"launch_type":     "FARGATE",
		"tags": map[string]string{
			"Name":        "invoice-prometheus-service",
			"Project":     "invoice",
			"Environment": "production",
		},
	}

	// Add ECS Service to mock AWS service
	mockAWS.AddResource("aws_ecs_service", prometheusService["id"].(string), prometheusService)

	// Test retrieving EFS File System
	efsResource, err := mockAWS.GetResource("aws_efs_file_system", prometheusEFS["id"].(string))
	require.NoError(t, err)
	assert.Equal(t, prometheusEFS["creation_token"], efsResource["creation_token"])
	assert.Equal(t, prometheusEFS["encrypted"], efsResource["encrypted"])
	assert.Equal(t, prometheusEFS["tags"], efsResource["tags"])

	// Test retrieving EFS Access Point
	apResource, err := mockAWS.GetResource("aws_efs_access_point", prometheusAP["id"].(string))
	require.NoError(t, err)
	assert.Equal(t, prometheusAP["file_system_id"], apResource["file_system_id"])
	assert.Equal(t, prometheusAP["tags"], apResource["tags"])

	// Test retrieving ECR Repository
	ecrResource, err := mockAWS.GetResource("aws_ecr_repository", prometheusECR["id"].(string))
	require.NoError(t, err)
	assert.Equal(t, prometheusECR["name"], ecrResource["name"])
	assert.Equal(t, prometheusECR["image_tag_mutability"], ecrResource["image_tag_mutability"])
	assert.Equal(t, prometheusECR["tags"], ecrResource["tags"])

	// Test retrieving ECS Task Definition
	taskDefResource, err := mockAWS.GetResource("aws_ecs_task_definition", prometheusTaskDef["id"].(string))
	require.NoError(t, err)
	assert.Equal(t, prometheusTaskDef["family"], taskDefResource["family"])
	assert.Equal(t, prometheusTaskDef["network_mode"], taskDefResource["network_mode"])
	assert.Equal(t, prometheusTaskDef["cpu"], taskDefResource["cpu"])
	assert.Equal(t, prometheusTaskDef["memory"], taskDefResource["memory"])
	assert.Equal(t, prometheusTaskDef["tags"], taskDefResource["tags"])

	// Test retrieving ECS Service
	serviceResource, err := mockAWS.GetResource("aws_ecs_service", prometheusService["id"].(string))
	require.NoError(t, err)
	assert.Equal(t, prometheusService["name"], serviceResource["name"])
	assert.Equal(t, prometheusService["cluster"], serviceResource["cluster"])
	assert.Equal(t, prometheusService["task_definition"], serviceResource["task_definition"])
	assert.Equal(t, prometheusService["desired_count"], serviceResource["desired_count"])
	assert.Equal(t, prometheusService["launch_type"], serviceResource["launch_type"])
	assert.Equal(t, prometheusService["tags"], serviceResource["tags"])
}