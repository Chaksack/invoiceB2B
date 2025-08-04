# Development Environment Configuration

This directory contains Terraform configuration specific to the development environment.

## Files

- `terraform.tfvars`: Contains environment-specific variable values for the development environment.

## Usage

The GitHub workflow automatically uses these configurations when deploying to the development environment. The workflow is triggered when:

1. Code is pushed to the `staging` branch
2. A workflow is manually triggered with the environment set to `dev`

## Local Development

To use these configurations locally, run:

```bash
terraform init -backend-config="key=environments/dev/terraform.tfstate"
terraform plan -var-file=environments/dev/terraform.tfvars
terraform apply -var-file=environments/dev/terraform.tfvars
```

## Environment-Specific Values

The development environment uses:
- Different resource naming with `-dev` suffixes
- Development-specific database names
- Separate credentials for services like RabbitMQ
- The same network CIDR ranges as staging, as they are not deployed simultaneously

## Notes

- The development environment is intended for testing and development purposes only.
- It may have fewer resources or smaller instance sizes than production.
- Sensitive data like passwords should be managed through AWS Secrets Manager in a real deployment.