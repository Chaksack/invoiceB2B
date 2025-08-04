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

func TestVPCConfiguration(t *testing.T) {
	// Setup Fiber app for testing
	app := utils.SetupFiberTestApp()

	// Create a test endpoint to validate VPC configuration
	app.Get("/test/vpc", func(c *fiber.Ctx) error {
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
			"vpc.tf",
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

		// Validate VPC configuration
		utils.ValidateTerraformResource(t, plan, "aws_vpc", "main", map[string]interface{}{
			"cidr_block":           "10.0.0.0/16",
			"enable_dns_support":   true,
			"enable_dns_hostnames": true,
			"tags.Name":            "${var.project_name}-vpc",
			"tags.Project":         "${var.project_name}",
			"tags.Environment":     "production",
		})

		// Validate Internet Gateway configuration
		utils.ValidateTerraformResource(t, plan, "aws_internet_gateway", "main", map[string]interface{}{
			"vpc_id":           "${aws_vpc.main.id}",
			"tags.Name":        "${var.project_name}-igw",
			"tags.Project":     "${var.project_name}",
			"tags.Environment": "production",
		})

		// Validate Public Subnet configuration
		utils.ValidateTerraformResource(t, plan, "aws_subnet", "public", map[string]interface{}{
			"vpc_id":                  "${aws_vpc.main.id}",
			"map_public_ip_on_launch": true,
			"tags.Tier":               "Public",
		})

		// Validate Private Subnet configuration
		utils.ValidateTerraformResource(t, plan, "aws_subnet", "private", map[string]interface{}{
			"vpc_id":                  "${aws_vpc.main.id}",
			"map_public_ip_on_launch": false,
			"tags.Tier":               "Private",
		})

		// Validate NAT Gateway configuration
		utils.ValidateTerraformResource(t, plan, "aws_nat_gateway", "main", map[string]interface{}{
			"subnet_id":       "${aws_subnet.public[count.index].id}",
			"allocation_id":   "${aws_eip.nat[count.index].id}",
			"tags.Name":       "${var.project_name}-nat-gw-${count.index + 1}",
			"tags.Project":    "${var.project_name}",
			"tags.Environment": "production",
		})

		// Validate Public Route Table configuration
		utils.ValidateTerraformResource(t, plan, "aws_route_table", "public", map[string]interface{}{
			"vpc_id":           "${aws_vpc.main.id}",
			"tags.Name":        "${var.project_name}-public-rt",
			"tags.Project":     "${var.project_name}",
			"tags.Environment": "production",
		})

		// Validate Private Route Table configuration
		utils.ValidateTerraformResource(t, plan, "aws_route_table", "private", map[string]interface{}{
			"vpc_id":           "${aws_vpc.main.id}",
			"tags.Name":        "${var.project_name}-private-rt-${count.index + 1}",
			"tags.Project":     "${var.project_name}",
			"tags.Environment": "production",
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "VPC configuration validated successfully",
		})
	})

	// Test the endpoint
	resp, err := app.Test(httptest.NewRequest("GET", "/test/vpc", nil))
	require.NoError(t, err)

	// Assert response
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestVPCResourceCreation(t *testing.T) {
	// Create a mock AWS service
	mockAWS := utils.NewMockAWSService()

	// Create mock VPC
	vpc := utils.MockVPC{
		ID:                "vpc-12345",
		CIDRBlock:         "10.0.0.0/16",
		EnableDNSSupport:  true,
		EnableDNSHostnames: true,
		Tags: map[string]string{
			"Name":        "invoice-vpc",
			"Project":     "invoice",
			"Environment": "production",
		},
	}

	// Add VPC to mock AWS service
	mockAWS.AddResource("aws_vpc", vpc.ID, map[string]interface{}{
		"id":                 vpc.ID,
		"cidr_block":         vpc.CIDRBlock,
		"enable_dns_support": vpc.EnableDNSSupport,
		"enable_dns_hostnames": vpc.EnableDNSHostnames,
		"tags":               vpc.Tags,
	})

	// Create mock Internet Gateway
	igw := utils.MockInternetGateway{
		ID:    "igw-12345",
		VPCID: vpc.ID,
		Tags: map[string]string{
			"Name":        "invoice-igw",
			"Project":     "invoice",
			"Environment": "production",
		},
	}

	// Add Internet Gateway to mock AWS service
	mockAWS.AddResource("aws_internet_gateway", igw.ID, map[string]interface{}{
		"id":     igw.ID,
		"vpc_id": igw.VPCID,
		"tags":   igw.Tags,
	})

	// Create mock Public Subnet
	publicSubnet := utils.MockSubnet{
		ID:                  "subnet-public-12345",
		VPCID:               vpc.ID,
		CIDRBlock:           "10.0.1.0/24",
		AvailabilityZone:    "us-east-1a",
		MapPublicIPOnLaunch: true,
		Tags: map[string]string{
			"Name":        "invoice-public-subnet-1",
			"Project":     "invoice",
			"Environment": "production",
			"Tier":        "Public",
		},
	}

	// Add Public Subnet to mock AWS service
	mockAWS.AddResource("aws_subnet", publicSubnet.ID, map[string]interface{}{
		"id":                    publicSubnet.ID,
		"vpc_id":                publicSubnet.VPCID,
		"cidr_block":            publicSubnet.CIDRBlock,
		"availability_zone":     publicSubnet.AvailabilityZone,
		"map_public_ip_on_launch": publicSubnet.MapPublicIPOnLaunch,
		"tags":                  publicSubnet.Tags,
	})

	// Create mock Private Subnet
	privateSubnet := utils.MockSubnet{
		ID:                  "subnet-private-12345",
		VPCID:               vpc.ID,
		CIDRBlock:           "10.0.101.0/24",
		AvailabilityZone:    "us-east-1a",
		MapPublicIPOnLaunch: false,
		Tags: map[string]string{
			"Name":        "invoice-private-subnet-1",
			"Project":     "invoice",
			"Environment": "production",
			"Tier":        "Private",
		},
	}

	// Add Private Subnet to mock AWS service
	mockAWS.AddResource("aws_subnet", privateSubnet.ID, map[string]interface{}{
		"id":                    privateSubnet.ID,
		"vpc_id":                privateSubnet.VPCID,
		"cidr_block":            privateSubnet.CIDRBlock,
		"availability_zone":     privateSubnet.AvailabilityZone,
		"map_public_ip_on_launch": privateSubnet.MapPublicIPOnLaunch,
		"tags":                  privateSubnet.Tags,
	})

	// Test retrieving VPC
	vpcResource, err := mockAWS.GetResource("aws_vpc", vpc.ID)
	require.NoError(t, err)
	assert.Equal(t, vpc.CIDRBlock, vpcResource["cidr_block"])
	assert.Equal(t, vpc.EnableDNSSupport, vpcResource["enable_dns_support"])
	assert.Equal(t, vpc.EnableDNSHostnames, vpcResource["enable_dns_hostnames"])
	assert.Equal(t, vpc.Tags, vpcResource["tags"])

	// Test retrieving Internet Gateway
	igwResource, err := mockAWS.GetResource("aws_internet_gateway", igw.ID)
	require.NoError(t, err)
	assert.Equal(t, igw.VPCID, igwResource["vpc_id"])
	assert.Equal(t, igw.Tags, igwResource["tags"])

	// Test retrieving Public Subnet
	publicSubnetResource, err := mockAWS.GetResource("aws_subnet", publicSubnet.ID)
	require.NoError(t, err)
	assert.Equal(t, publicSubnet.VPCID, publicSubnetResource["vpc_id"])
	assert.Equal(t, publicSubnet.CIDRBlock, publicSubnetResource["cidr_block"])
	assert.Equal(t, publicSubnet.AvailabilityZone, publicSubnetResource["availability_zone"])
	assert.Equal(t, publicSubnet.MapPublicIPOnLaunch, publicSubnetResource["map_public_ip_on_launch"])
	assert.Equal(t, publicSubnet.Tags, publicSubnetResource["tags"])

	// Test retrieving Private Subnet
	privateSubnetResource, err := mockAWS.GetResource("aws_subnet", privateSubnet.ID)
	require.NoError(t, err)
	assert.Equal(t, privateSubnet.VPCID, privateSubnetResource["vpc_id"])
	assert.Equal(t, privateSubnet.CIDRBlock, privateSubnetResource["cidr_block"])
	assert.Equal(t, privateSubnet.AvailabilityZone, privateSubnetResource["availability_zone"])
	assert.Equal(t, privateSubnet.MapPublicIPOnLaunch, privateSubnetResource["map_public_ip_on_launch"])
	assert.Equal(t, privateSubnet.Tags, privateSubnetResource["tags"])
}