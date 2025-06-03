#!/bin/bash
# Exit on error, but allow for proper error handling
set -e

echo "Setting up Terraform backend resources..."

# Set environment if provided as an argument, default to "dev"
environment=${1:-dev}
echo "Using environment: $environment"

# Function to retry AWS commands with exponential backoff
function retry_aws_command {
  local max_attempts=5
  local timeout=1
  local attempt=1
  local exit_code=0

  while (( $attempt <= $max_attempts ))
  do
    echo "Attempt $attempt of $max_attempts: $@"
    "$@"
    exit_code=$?

    if [[ $exit_code -eq 0 ]]; then
      echo "Command succeeded."
      break
    fi

    echo "Command failed with exit code $exit_code. Retrying in $timeout seconds..."
    sleep $timeout
    attempt=$(( attempt + 1 ))
    timeout=$(( timeout * 2 ))
  done

  if [[ $exit_code -ne 0 ]]; then
    echo "Command '$@' failed after $max_attempts attempts"
    return $exit_code
  fi

  return 0
}

# Initialize and apply the bootstrap configuration
echo "Initializing bootstrap Terraform configuration..."
terraform init

echo "Creating S3 bucket and DynamoDB table for Terraform backend..."
# Extract project_name from variables.tf
project_name=$(grep -A 3 'variable "project_name"' variables.tf | grep 'default' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$project_name" ]; then
  echo "ERROR: Could not extract project_name from variables.tf"
  exit 1
fi
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
echo "Checking if S3 bucket $bucket_name exists..."
if retry_aws_command aws s3api head-bucket --bucket "$bucket_name" 2>/dev/null; then
  echo "S3 bucket $bucket_name already exists, skipping creation"
  s3_exists=true
else
  echo "S3 bucket $bucket_name does not exist, will create it"
  s3_exists=false
fi

# Check if DynamoDB table already exists
table_name="${bucket_prefix}-terraform-locks"
echo "Checking if DynamoDB table $table_name exists..."
if retry_aws_command aws dynamodb describe-table --table-name "$table_name" 2>/dev/null; then
  echo "DynamoDB table $table_name already exists, skipping creation"
  dynamodb_exists=true
else
  echo "DynamoDB table $table_name does not exist, will create it"
  dynamodb_exists=false
fi

# Create resources that don't exist using AWS CLI directly
# This is more reliable than using Terraform for these bootstrap resources
if [ "$s3_exists" = false ] || [ "$dynamodb_exists" = false ]; then
  echo "Creating missing resources using AWS CLI..."

  # Create S3 bucket if it doesn't exist
  if [ "$s3_exists" = false ]; then
    echo "Creating S3 bucket: $bucket_name"
    # Create the bucket with appropriate region configuration
    if [ "$AWS_REGION" = "us-east-1" ]; then
      retry_aws_command aws s3api create-bucket --bucket "$bucket_name" --region us-east-1
    else
      retry_aws_command aws s3api create-bucket --bucket "$bucket_name" --region "$AWS_REGION" --create-bucket-configuration LocationConstraint="$AWS_REGION"
    fi

    # Wait for bucket to be available
    echo "Waiting for S3 bucket to become available..."
    retry_aws_command aws s3api wait bucket-exists --bucket "$bucket_name"

    # Enable versioning
    echo "Enabling versioning on S3 bucket..."
    retry_aws_command aws s3api put-bucket-versioning --bucket "$bucket_name" --versioning-configuration Status=Enabled

    # Enable encryption
    echo "Enabling encryption on S3 bucket..."
    retry_aws_command aws s3api put-bucket-encryption --bucket "$bucket_name" --server-side-encryption-configuration '{"Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]}'

    # Block public access
    echo "Blocking public access to S3 bucket..."
    retry_aws_command aws s3api put-public-access-block --bucket "$bucket_name" --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

    echo "S3 bucket created and configured successfully"

    # Verify bucket exists after creation
    if ! retry_aws_command aws s3api head-bucket --bucket "$bucket_name" 2>/dev/null; then
      echo "ERROR: Failed to verify S3 bucket creation"
      exit 1
    fi
  fi

  # Create DynamoDB table if it doesn't exist
  if [ "$dynamodb_exists" = false ]; then
    echo "Creating DynamoDB table: $table_name"
    retry_aws_command aws dynamodb create-table \
      --table-name "$table_name" \
      --attribute-definitions AttributeName=LockID,AttributeType=S \
      --key-schema AttributeName=LockID,KeyType=HASH \
      --billing-mode PAY_PER_REQUEST \
      --region "${AWS_REGION:-us-east-1}"

    # Wait for table to be active
    echo "Waiting for DynamoDB table to become active..."
    retry_aws_command aws dynamodb wait table-exists --table-name "$table_name"

    echo "DynamoDB table created successfully"

    # Verify table exists after creation
    if ! retry_aws_command aws dynamodb describe-table --table-name "$table_name" 2>/dev/null; then
      echo "ERROR: Failed to verify DynamoDB table creation"
      exit 1
    fi
  fi

  # Set create_bootstrap_resources to true since we're creating at least one resource
  create_bootstrap_resources=true
else
  echo "Both resources already exist, skipping creation"
  # Set create_bootstrap_resources to false since both resources already exist
  create_bootstrap_resources=false
fi

# Apply Terraform configuration with the create_bootstrap_resources variable
echo "Applying Terraform configuration with create_bootstrap_resources=${create_bootstrap_resources}..."
terraform apply -auto-approve -var="create_bootstrap_resources=${create_bootstrap_resources}"

echo "Backend resources created successfully!"
echo "S3 bucket: ${bucket_prefix}-terraform-state"
echo "DynamoDB table: ${bucket_prefix}-terraform-locks"

# Generate backend.tf file with the correct values
echo "Generating backend.tf file with dynamic configuration..."
chmod +x ./generate-backend.sh
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