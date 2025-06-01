# Bootstrap Terraform resources
# This file creates the S3 bucket and DynamoDB table needed for the Terraform backend
# Run this with: terraform init && terraform apply
# After running this, you can use the main Terraform configuration with the S3 backend

provider "aws" {
  region = "us-east-1"  # Use the same region as in your backend configuration
  alias  = "bootstrap"
}

resource "aws_s3_bucket" "terraform_state" {
  provider = aws.bootstrap
  bucket   = "invoiceb2b-terraform-state"

  # Prevent accidental deletion of this S3 bucket
  lifecycle {
    prevent_destroy = true
  }

  tags = {
    Name        = "Terraform State"
    Environment = "All"
    Project     = "invoiceb2b"
  }
}

resource "aws_s3_bucket_versioning" "terraform_state" {
  provider = aws.bootstrap
  bucket   = aws_s3_bucket.terraform_state.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  provider = aws.bootstrap
  bucket   = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "terraform_state" {
  provider = aws.bootstrap
  bucket   = aws_s3_bucket.terraform_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_dynamodb_table" "terraform_locks" {
  provider     = aws.bootstrap
  name         = "invoiceb2b-terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name        = "Terraform Locks"
    Environment = "All"
    Project     = "invoiceb2b"
  }
}

output "s3_bucket_name" {
  value       = aws_s3_bucket.terraform_state.id
  description = "The name of the S3 bucket"
}

output "dynamodb_table_name" {
  value       = aws_dynamodb_table.terraform_locks.id
  description = "The name of the DynamoDB table"
}