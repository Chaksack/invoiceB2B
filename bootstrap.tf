# Bootstrap Terraform resources
# This file creates the S3 bucket and DynamoDB table needed for the Terraform backend

# Use the project_name variable from variables.tf
# No need to redefine it here as it's already defined in variables.tf

resource "aws_s3_bucket" "terraform_state" {
  bucket = "${var.project_name}-terraform-state"

  # Prevent accidental deletion of this S3 bucket
  lifecycle {
    prevent_destroy = true
    # Ignore errors related to bucket already existing
    ignore_changes = [bucket]
    # Prevent errors when bucket already exists
    create_before_destroy = false
  }

  tags = {
    Name        = "${var.project_name}-terraform-state"
    Environment = "All"
    Project     = var.project_name
  }
}

resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  versioning_configuration {
    status = "Enabled"
  }

  # Ignore errors if bucket already exists
  lifecycle {
    ignore_changes = [bucket]
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }

  # Ignore errors if bucket already exists
  lifecycle {
    ignore_changes = [bucket]
  }
}

resource "aws_s3_bucket_public_access_block" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

  # Ignore errors if bucket already exists
  lifecycle {
    ignore_changes = [bucket]
  }
}

resource "aws_dynamodb_table" "terraform_locks" {
  name         = "${var.project_name}-terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  lifecycle {
    # Ignore errors related to table already existing
    ignore_changes = [name]
    # Prevent errors when table already exists
    create_before_destroy = false
  }

  tags = {
    Name        = "${var.project_name}-terraform-locks"
    Environment = "All"
    Project     = var.project_name
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