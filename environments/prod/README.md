# Production Environment Configuration

This directory contains Terraform configuration specific to the production environment.

## Files

- `terraform.tfvars`: Contains environment-specific variable values for the production environment.

## Usage

The GitHub workflow automatically uses these configurations when deploying to the production environment. The workflow is triggered when:

1. Code is pushed to the `main` branch
2. A workflow is manually triggered with the environment set to `prod`

## Manual Approval

Deployments to the production environment require manual approval in the GitHub workflow. This adds an extra layer of safety for production deployments.

## Local Development

To use these configurations locally, run:

```bash
terraform init -backend-config="key=environments/prod/terraform.tfstate"
terraform plan -var-file=environments/prod/terraform.tfvars
terraform apply -var-file=environments/prod/terraform.tfvars
```

## Environment-Specific Values

The production environment uses:
- Different resource naming with `-prod` suffixes
- Production-specific database names
- Separate, more secure credentials for services like RabbitMQ
- Different network CIDR ranges from development to allow both environments to exist simultaneously if needed

## Notes

- The production environment is intended for live, customer-facing deployments.
- It should have appropriate scaling, monitoring, and backup configurations.
- All sensitive data should be managed through AWS Secrets Manager in a real deployment.
- Always test changes in the development environment before deploying to production.