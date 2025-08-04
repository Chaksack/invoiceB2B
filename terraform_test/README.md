# Terraform Unit Tests

This directory contains unit tests for the Terraform configurations in the project. The tests use Golang and the Gofiber framework to validate the Terraform configurations.

## Test Structure

The tests are organized into the following directories:

- `infrastructure/`: Tests for infrastructure resources (VPC, security groups, etc.)
- `application/`: Tests for application resources (ECS, RDS, etc.)
- `monitoring/`: Tests for monitoring resources (Prometheus, Grafana, etc.)
- `utils/`: Utility functions for testing Terraform configurations

## Prerequisites

Before running the tests, make sure you have the following installed:

- Go 1.16 or later
- Terraform 1.0 or later
- The following Go packages:
  - github.com/gofiber/fiber/v2
  - github.com/stretchr/testify

You can install the required Go packages with:

```bash
go get github.com/gofiber/fiber/v2
go get github.com/stretchr/testify
```

## Running the Tests

You can run the tests using the provided `run_tests.go` script:

```bash
# Run all tests
go run run_tests.go

# Run tests in a specific directory
go run run_tests.go -dir infrastructure

# Run a specific test
go run run_tests.go -test TestVPCConfiguration

# Run tests with verbose output
go run run_tests.go -v

# Generate a coverage report
go run run_tests.go -coverage
```

### Command Line Options

The `run_tests.go` script supports the following command line options:

- `-dir`: Directory to run tests from (default: all)
- `-test`: Specific test to run (e.g., TestVPCConfiguration)
- `-v`: Verbose output
- `-coverage`: Generate coverage report
- `-coverfile`: Coverage file name (default: coverage.out)

## Test Types

The tests are divided into two main types:

1. **Configuration Tests**: These tests validate the Terraform configuration by running `terraform init` and `terraform plan`, then checking that the resources in the plan have the expected configuration.

2. **Resource Creation Tests**: These tests use a mock AWS service to simulate the creation of AWS resources and validate that the resources are created with the expected attributes.

## Adding New Tests

To add a new test for a Terraform file:

1. Identify the appropriate directory for the test (infrastructure, application, or monitoring).
2. Create a new test file with the naming convention `<resource>_test.go`.
3. Implement the configuration test and resource creation test for the Terraform file.
4. Run the tests to ensure they pass.

## Test Utilities

The `utils/` directory contains utility functions for testing Terraform configurations:

- `test_utils.go`: Contains functions for running Terraform commands and validating Terraform resources.
- `mock_aws.go`: Contains a mock implementation of AWS services for testing resource creation.

## Example Test

Here's an example of a configuration test for a VPC:

```go
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
        // Go up one level to get to the project root
        projectRoot = filepath.Dir(projectRoot)

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
```

## Troubleshooting

If you encounter any issues running the tests, check the following:

- Make sure you have the required Go packages installed.
- Make sure Terraform is installed and in your PATH.
- Make sure you're running the tests from the correct directory.
- Check the error messages for clues about what might be wrong.

If you're still having issues, try running the tests with the `-v` flag for more detailed output.

### Known Limitations

#### Configuration Tests

The configuration tests (TestVPCConfiguration, TestSecurityGroupsConfiguration, etc.) attempt to run actual Terraform commands (init, plan) to validate the Terraform configurations. These tests may fail in environments where:

- Terraform is not installed
- The tests don't have permission to run Terraform commands
- The Terraform commands take too long to execute (timeout errors)

If you encounter timeout errors when running the configuration tests, you can increase the timeout by modifying the test code or focus on the resource creation tests instead.

#### Resource Creation Tests

The resource creation tests (TestVPCResourceCreation, TestSecurityGroupsResourceCreation, etc.) use a mock AWS service to simulate the creation of AWS resources. These tests don't require Terraform to be installed and should pass in any environment.

If you're primarily interested in validating the structure and attributes of the AWS resources defined in the Terraform configurations, the resource creation tests should be sufficient.