# Backend configuration for storing Terraform state in S3
# This allows for team collaboration, state locking, and versioning
#
# To use with different environments, initialize with:
# terraform init -backend-config="key=environments/dev/terraform.tfstate"
# terraform init -backend-config="key=environments/staging/terraform.tfstate"
# terraform init -backend-config="key=environments/prod/terraform.tfstate"
# Temporarily commented out S3 backend configuration until the S3 bucket is created
# To enable this backend, first run the setup-backend.sh script to create the required resources
# Then uncomment this configuration and run: terraform init
# Backend configuration for storing Terraform state in S3
# This allows for team collaboration, state locking, and versioning
#
# Note: The bucket and dynamodb_table names follow the pattern:
# - bucket: "${var.project_name}-terraform-state"
# - dynamodb_table: "${var.project_name}-terraform-locks"
#
# However, variables cannot be used directly in the backend configuration.
# The values here must match the values in bootstrap.tf.
# terraform {
#   backend "s3" {
#     bucket         = "invoicefin-terraform-state"
#     key            = "terraform.tfstate"  # Default state file path, override with -backend-config
#     region         = "us-east-1"
#     encrypt        = true
#     dynamodb_table = "invoicefin-terraform-locks"
#   }
# }