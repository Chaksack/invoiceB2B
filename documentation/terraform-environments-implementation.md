# Terraform Environments Implementation

## Overview

This document describes the implementation of separate development and production environments for Terraform in the InvoiceB2B project. The implementation allows for environment-specific configurations while maintaining a consistent deployment process through GitHub Actions.

## Changes Made

1. **Created Environment Directory Structure**
   - Created `environments/dev/` and `environments/prod/` directories to store environment-specific configurations

2. **Created Environment-Specific Configuration Files**
   - Added `terraform.tfvars` files for each environment with appropriate values
   - Created README.md files in each environment directory explaining usage

3. **Updated GitHub Workflow**
   - Modified Terraform commands to use environment-specific tfvars files
   - Ensured the workflow uses the correct environment based on the branch:
     - `main` branch → `prod` environment
     - `staging` branch → `dev` environment

## Environment-Specific Configurations

### Development Environment (`environments/dev/`)

The development environment configuration includes:
- Environment name set to "dev"
- Resource naming with "-dev" suffixes
- Development-specific database names
- Separate credentials for services like RabbitMQ

### Production Environment (`environments/prod/`)

The production environment configuration includes:
- Environment name set to "prod"
- Resource naming with "-prod" suffixes
- Production-specific database names
- More secure credentials for services
- Different network CIDR ranges to allow both environments to exist simultaneously if needed

## How It Works

1. The GitHub workflow determines the environment based on the branch:
   ```yaml
   ENVIRONMENT: ${{ github.event.inputs.environment || (github.ref == 'refs/heads/main' && 'prod') || (github.ref == 'refs/heads/staging' && 'dev') || 'staging' }}
   ```

2. Terraform is initialized with an environment-specific state file:
   ```yaml
   terraform init -reconfigure -backend-config="key=environments/${{ env.ENVIRONMENT }}/terraform.tfstate"
   ```

3. Terraform commands use the environment-specific tfvars file:
   ```yaml
   terraform plan -var-file=environments/${{ env.ENVIRONMENT }}/terraform.tfvars -out=tfplan
   ```

4. This ensures that each environment has:
   - Its own Terraform state file
   - Its own configuration values
   - Isolated resources that don't conflict with other environments

## Benefits

1. **Environment Isolation**: Each environment has its own resources with unique names, preventing conflicts
2. **Configuration Management**: Environment-specific variables are stored in separate files, making them easier to manage
3. **Consistent Deployment Process**: The same workflow is used for all environments, reducing the risk of deployment errors
4. **Branch-Based Deployments**: Deployments to different environments are automatically triggered by pushing to the appropriate branch

## Local Development

To use these configurations locally, run:

```bash
# For development environment
terraform init -backend-config="key=environments/dev/terraform.tfstate"
terraform plan -var-file=environments/dev/terraform.tfvars
terraform apply -var-file=environments/dev/terraform.tfvars

# For production environment
terraform init -backend-config="key=environments/prod/terraform.tfstate"
terraform plan -var-file=environments/prod/terraform.tfvars
terraform apply -var-file=environments/prod/terraform.tfvars
```

## Future Improvements

1. Consider adding a staging environment with its own configuration
2. Implement environment-specific resource sizing (e.g., smaller instances for development)
3. Use AWS Secrets Manager for sensitive values instead of storing them in tfvars files
4. Add environment-specific monitoring and alerting configurations