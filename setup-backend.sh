#!/bin/bash
set -e

echo "Setting up Terraform backend resources..."

# Set environment if provided as an argument, default to "dev"
environment=${1:-dev}
echo "Using environment: $environment"

# Initialize and apply the bootstrap configuration
echo "Initializing bootstrap Terraform configuration..."
terraform init

echo "Creating S3 bucket and DynamoDB table for Terraform backend..."
# Extract project_name from variables.tf
project_name=$(grep -A 3 'variable "project_name"' variables.tf | grep 'default' | sed -E 's/.*"([^"]+)".*/\1/')
echo "Using project_name: $project_name"

# For dev environment, use a different bucket name prefix
if [ "$environment" = "dev" ]; then
  bucket_prefix="invoicedev"
  echo "Using dev-specific bucket prefix: $bucket_prefix"
else
  bucket_prefix="$project_name"
  echo "Using standard bucket prefix: $bucket_prefix"
fi

# Check if S3 bucket already exists
bucket_name="${bucket_prefix}-terraform-state"
if aws s3api head-bucket --bucket "$bucket_name" 2>/dev/null; then
  echo "S3 bucket $bucket_name already exists, skipping creation"
  s3_target=""
else
  echo "S3 bucket $bucket_name does not exist, will create it"
  s3_target="-target=aws_s3_bucket.terraform_state"
fi

# Check if DynamoDB table already exists
table_name="${bucket_prefix}-terraform-locks"
if aws dynamodb describe-table --table-name "$table_name" 2>/dev/null; then
  echo "DynamoDB table $table_name already exists, skipping creation"
  dynamodb_target=""
else
  echo "DynamoDB table $table_name does not exist, will create it"
  dynamodb_target="-target=aws_dynamodb_table.terraform_locks"
fi

# Only create resources if they don't exist
if [ -n "$s3_target" ] || [ -n "$dynamodb_target" ]; then
  echo "Creating resources using AWS CLI..."

  # Create S3 bucket if it doesn't exist
  if [ -n "$s3_target" ]; then
    echo "Creating S3 bucket: $bucket_name"
    aws s3api create-bucket --bucket "$bucket_name" --region us-east-1

    # Enable versioning
    aws s3api put-bucket-versioning --bucket "$bucket_name" --versioning-configuration Status=Enabled

    # Enable encryption
    aws s3api put-bucket-encryption --bucket "$bucket_name" --server-side-encryption-configuration '{"Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]}'

    # Block public access
    aws s3api put-public-access-block --bucket "$bucket_name" --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

    echo "S3 bucket created successfully"
  fi

  # Create DynamoDB table if it doesn't exist
  if [ -n "$dynamodb_target" ]; then
    echo "Creating DynamoDB table: $table_name"
    aws dynamodb create-table --table-name "$table_name" --attribute-definitions AttributeName=LockID,AttributeType=S --key-schema AttributeName=LockID,KeyType=HASH --billing-mode PAY_PER_REQUEST --region us-east-1
    echo "DynamoDB table created successfully"
  fi
else
  echo "Both resources already exist, skipping creation"
fi

echo "Backend resources created successfully!"
echo "S3 bucket: ${bucket_prefix}-terraform-state"
echo "DynamoDB table: ${bucket_prefix}-terraform-locks"

# Generate backend.tf file with the correct values
echo "Generating backend.tf file with dynamic configuration..."
./generate-backend.sh $environment

echo "You can now initialize the main Terraform configuration with the S3 backend:"
echo "terraform init"

echo "For environment-specific state files, use:"
echo "terraform init -backend-config=\"key=environments/dev/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/staging/terraform.tfstate\""
echo "terraform init -backend-config=\"key=environments/prod/terraform.tfstate\""

echo "Or regenerate the backend.tf file for a specific environment:"
echo "./generate-backend.sh dev"
echo "./generate-backend.sh staging"
echo "./generate-backend.sh prod"