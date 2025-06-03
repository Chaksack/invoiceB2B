#!/bin/bash
set -e

echo "Setting up Terraform backend resources..."

# Navigate to the bootstrap directory
cd bootstrap

# Initialize and apply the bootstrap configuration
echo "Initializing bootstrap Terraform configuration..."
terraform init

echo "Creating S3 bucket and DynamoDB table for Terraform backend..."
terraform apply -auto-approve

# Navigate back to the main directory
cd ..

echo "Backend resources created successfully!"
echo "S3 bucket: invoiceapp-terraform-state"
echo "DynamoDB table: invoice-terraform-locks"

echo "You can now initialize the main Terraform configuration with the S3 backend:"
echo "terraform init"

echo "For environment-specific state files, use:"
echo "terraform init -backend-config=\"key=environments/dev/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/staging/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/prod/terraform.tfstate\""