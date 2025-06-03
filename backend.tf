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
terraform {
  backend "s3" {
    bucket         = "invoiceapp-terraform-state"
    key            = "terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "invoice-terraform-locks"
  }
}