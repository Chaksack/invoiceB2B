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

# Check if S3 bucket already exists
bucket_name="${project_name}-terraform-state"
if aws s3api head-bucket --bucket "$bucket_name" 2>/dev/null; then
  echo "S3 bucket $bucket_name already exists, skipping creation"
  s3_target=""
else
  echo "S3 bucket $bucket_name does not exist, will create it"
  s3_target="-target=aws_s3_bucket.terraform_state"
fi

# Check if DynamoDB table already exists
table_name="${project_name}-terraform-locks"
if aws dynamodb describe-table --table-name "$table_name" 2>/dev/null; then
  echo "DynamoDB table $table_name already exists, skipping creation"
  dynamodb_target=""
else
  echo "DynamoDB table $table_name does not exist, will create it"
  dynamodb_target="-target=aws_dynamodb_table.terraform_locks"
fi

# Only apply if at least one resource needs to be created
if [ -n "$s3_target" ] || [ -n "$dynamodb_target" ]; then
  terraform apply -auto-approve $s3_target $dynamodb_target
else
  echo "Both resources already exist, skipping terraform apply"
fi

echo "Backend resources created successfully!"
echo "S3 bucket: ${project_name}-terraform-state"
echo "DynamoDB table: ${project_name}-terraform-locks"

echo "You can now initialize the main Terraform configuration with the S3 backend:"
echo "terraform init"

echo "For environment-specific state files, use:"
echo "terraform init -backend-config=\"key=environments/dev/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/staging/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/prod/terraform.tfstate\""