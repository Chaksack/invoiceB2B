package infrastructure

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

func TestSecurityGroupsConfiguration(t *testing.T) {
	// Setup Fiber app for testing
	app := utils.SetupFiberTestApp()

	// Create a test endpoint to validate security groups configuration
	app.Get("/test/security-groups", func(c *fiber.Ctx) error {
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
			"security_groups.tf",
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

		// Validate ALB security group configuration
		utils.ValidateTerraformResource(t, plan, "aws_security_group", "alb", map[string]interface{}{
			"name":        "${var.project_name}-alb-sg",
			"description": "Security group for ALB",
			"vpc_id":      "${aws_vpc.main.id}",
			"tags.Name":   "${var.project_name}-alb-sg",
			"tags.Project": "${var.project_name}",
		})

		// Validate ECS tasks security group configuration
		utils.ValidateTerraformResource(t, plan, "aws_security_group", "ecs_tasks", map[string]interface{}{
			"name":        "${var.project_name}-ecs-tasks-sg",
			"description": "Security group for ECS tasks",
			"vpc_id":      "${aws_vpc.main.id}",
			"tags.Name":   "${var.project_name}-ecs-tasks-sg",
			"tags.Project": "${var.project_name}",
		})

		// Validate RDS security group configuration
		utils.ValidateTerraformResource(t, plan, "aws_security_group", "rds", map[string]interface{}{
			"name":        "${var.project_name}-rds-sg",
			"description": "Security group for RDS instance",
			"vpc_id":      "${aws_vpc.main.id}",
			"tags.Name":   "${var.project_name}-rds-sg",
			"tags.Project": "${var.project_name}",
		})

		// Validate ElastiCache security group configuration
		utils.ValidateTerraformResource(t, plan, "aws_security_group", "elasticache", map[string]interface{}{
			"name":        "${var.project_name}-elasticache-sg",
			"description": "Security group for ElastiCache Redis",
			"vpc_id":      "${aws_vpc.main.id}",
			"tags.Name":   "${var.project_name}-elasticache-sg",
			"tags.Project": "${var.project_name}",
		})

		// Validate RabbitMQ security group configuration
		utils.ValidateTerraformResource(t, plan, "aws_security_group", "rabbitmq", map[string]interface{}{
			"name":        "${var.project_name}-rabbitmq-sg",
			"description": "Security group for Amazon MQ (RabbitMQ)",
			"vpc_id":      "${aws_vpc.main.id}",
			"tags.Name":   "${var.project_name}-rabbitmq-sg",
			"tags.Project": "${var.project_name}",
		})

		// Validate EFS security group configuration
		utils.ValidateTerraformResource(t, plan, "aws_security_group", "efs", map[string]interface{}{
			"name":        "${var.project_name}-efs-sg",
			"description": "Security group for EFS mount targets",
			"vpc_id":      "${aws_vpc.main.id}",
			"tags.Name":   "${var.project_name}-efs-sg",
			"tags.Project": "${var.project_name}",
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Security groups configuration validated successfully",
		})
	})

	// Test the endpoint
	resp, err := app.Test(httptest.NewRequest("GET", "/test/security-groups", nil))
	require.NoError(t, err)

	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestSecurityGroupsResourceCreation(t *testing.T) {
	// Create a mock AWS service
	mockAWS := utils.NewMockAWSService()

	// Create mock VPC for reference
	vpc := utils.MockVPC{
		ID:       "vpc-12345",
		CIDRBlock: "10.0.0.0/16",
		Tags: map[string]string{
			"Name": "invoice-vpc",
		},
	}

	// Add VPC to mock AWS service
	mockAWS.AddResource("aws_vpc", vpc.ID, map[string]interface{}{
		"id":         vpc.ID,
		"cidr_block": vpc.CIDRBlock,
		"tags":       vpc.Tags,
	})

	// Create mock ALB security group
	albSG := utils.MockSecurityGroup{
		ID:          "sg-alb-12345",
		VPCID:       vpc.ID,
		Name:        "invoice-alb-sg",
		Description: "Security group for ALB",
		IngressRules: []utils.MockSecurityGroupRule{
			{
				Protocol:    "tcp",
				FromPort:    80,
				ToPort:      80,
				CIDRBlocks:  []string{"0.0.0.0/0"},
				Description: "HTTP",
			},
			{
				Protocol:    "tcp",
				FromPort:    443,
				ToPort:      443,
				CIDRBlocks:  []string{"0.0.0.0/0"},
				Description: "HTTPS",
			},
		},
		EgressRules: []utils.MockSecurityGroupRule{
			{
				Protocol:    "-1",
				FromPort:    0,
				ToPort:      0,
				CIDRBlocks:  []string{"0.0.0.0/0"},
				Description: "All outbound traffic",
			},
		},
		Tags: map[string]string{
			"Name":    "invoice-alb-sg",
			"Project": "invoice",
		},
	}

	// Add ALB security group to mock AWS service
	mockAWS.AddResource("aws_security_group", albSG.ID, map[string]interface{}{
		"id":          albSG.ID,
		"vpc_id":      albSG.VPCID,
		"name":        albSG.Name,
		"description": albSG.Description,
		"ingress":     albSG.IngressRules,
		"egress":      albSG.EgressRules,
		"tags":        albSG.Tags,
	})

	// Create mock ECS tasks security group
	ecsTasksSG := utils.MockSecurityGroup{
		ID:          "sg-ecs-tasks-12345",
		VPCID:       vpc.ID,
		Name:        "invoice-ecs-tasks-sg",
		Description: "Security group for ECS tasks",
		IngressRules: []utils.MockSecurityGroupRule{
			{
				Protocol:    "tcp",
				FromPort:    3000, // app_port
				ToPort:      3000,
				Description: "API port",
			},
			{
				Protocol:    "tcp",
				FromPort:    5678,
				ToPort:      5678,
				Description: "N8N port",
			},
			{
				Protocol:    "tcp",
				FromPort:    9000,
				ToPort:      9000,
				Description: "SonarQube port",
			},
		},
		EgressRules: []utils.MockSecurityGroupRule{
			{
				Protocol:    "-1",
				FromPort:    0,
				ToPort:      0,
				CIDRBlocks:  []string{"0.0.0.0/0"},
				Description: "All outbound traffic",
			},
		},
		Tags: map[string]string{
			"Name":    "invoice-ecs-tasks-sg",
			"Project": "invoice",
		},
	}

	// Add ECS tasks security group to mock AWS service
	mockAWS.AddResource("aws_security_group", ecsTasksSG.ID, map[string]interface{}{
		"id":          ecsTasksSG.ID,
		"vpc_id":      ecsTasksSG.VPCID,
		"name":        ecsTasksSG.Name,
		"description": ecsTasksSG.Description,
		"ingress":     ecsTasksSG.IngressRules,
		"egress":      ecsTasksSG.EgressRules,
		"tags":        ecsTasksSG.Tags,
	})

	// Test retrieving ALB security group
	albSGResource, err := mockAWS.GetResource("aws_security_group", albSG.ID)
	require.NoError(t, err)
	assert.Equal(t, albSG.VPCID, albSGResource["vpc_id"])
	assert.Equal(t, albSG.Name, albSGResource["name"])
	assert.Equal(t, albSG.Description, albSGResource["description"])
	assert.Equal(t, albSG.Tags, albSGResource["tags"])

	// Test retrieving ECS tasks security group
	ecsTasksSGResource, err := mockAWS.GetResource("aws_security_group", ecsTasksSG.ID)
	require.NoError(t, err)
	assert.Equal(t, ecsTasksSG.VPCID, ecsTasksSGResource["vpc_id"])
	assert.Equal(t, ecsTasksSG.Name, ecsTasksSGResource["name"])
	assert.Equal(t, ecsTasksSG.Description, ecsTasksSGResource["description"])
	assert.Equal(t, ecsTasksSG.Tags, ecsTasksSGResource["tags"])
}