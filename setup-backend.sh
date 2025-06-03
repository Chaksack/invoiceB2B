#!/bin/bash
set -e

echo "Setting up Terraform backend resources..."

# Initialize and apply the bootstrap configuration
echo "Initializing bootstrap Terraform configuration..."
terraform init

echo "Creating S3 bucket and DynamoDB table for Terraform backend..."
# Extract project_name from variables.tf
project_name=$(grep -A 3 'variable "project_name"' variables.tf | grep 'default' | sed -E 's/.*"([^"]+)".*/\1/')
echo "Using project_name: $project_name"
terraform apply -auto-approve -target=aws_s3_bucket.terraform_state -target=aws_dynamodb_table.terraform_locks

echo "Backend resources created successfully!"
echo "S3 bucket: ${project_name}-terraform-state"
echo "DynamoDB table: ${project_name}-terraform-locks"

echo "You can now initialize the main Terraform configuration with the S3 backend:"
echo "terraform init"

echo "For environment-specific state files, use:"
echo "terraform init -backend-config=\"key=environments/dev/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/staging/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/prod/terraform.tfstate\""