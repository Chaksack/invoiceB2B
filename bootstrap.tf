# Bootstrap Terraform resources
# This file creates the S3 bucket and DynamoDB table needed for the Terraform backend

# Use the project_name variable from variables.tf
# No need to redefine it here as it's already defined in variables.tf

locals {
  # For dev environment, use a different bucket name prefix
  # If the environment variable BUCKET_PREFIX is set, use that instead
  # This allows for overriding the bucket prefix at runtime
  bucket_prefix = var.environment == "dev" ? (coalesce(var.bucket_prefix_override, "invoicedev", "invoiceapidev")) : var.project_name

  # Check if resources exist with the invoiceapidev prefix (for backward compatibility)
  # These will be computed after the data sources are evaluated
  use_api_prefix_for_s3 = (!var.create_bootstrap_resources && var.environment == "dev" &&
  local.bucket_prefix != "invoiceapidev" &&
  length(data.aws_s3_bucket.existing_terraform_state_api) > 0)

  use_api_prefix_for_dynamodb = (!var.create_bootstrap_resources && var.environment == "dev" &&
  local.bucket_prefix != "invoiceapidev" &&
  length(data.aws_dynamodb_table.existing_terraform_locks_api) > 0)
}

# Check if the S3 bucket already exists with the current prefix
data "aws_s3_bucket" "existing_terraform_state" {
  count  = var.create_bootstrap_resources ? 0 : 1
  bucket = "${local.bucket_prefix}-terraform-state"
}

# Check if the S3 bucket already exists with the invoiceapidev prefix (for backward compatibility)
data "aws_s3_bucket" "existing_terraform_state_api" {
  count  = (!var.create_bootstrap_resources && var.environment == "dev" && local.bucket_prefix != "invoiceapidev") ? 1 : 0
  bucket = "invoiceapidev-terraform-state"
}


resource "aws_s3_bucket" "terraform_state" {
  count  = var.create_bootstrap_resources ? 1 : 0
  bucket = local.use_api_prefix_for_s3 ? "invoiceapidev-terraform-state" : "${local.bucket_prefix}-terraform-state"

  # Prevent accidental deletion of this S3 bucket
  lifecycle {
    prevent_destroy = true
    # Ignore errors related to bucket already existing
    ignore_changes = [bucket, id]
    # Prevent errors when bucket already exists
    create_before_destroy = false
  }

  tags = {
    Name        = "${local.bucket_prefix}-terraform-state"
    Environment = "All"
    Project     = var.project_name
  }
}

resource "aws_s3_bucket_versioning" "terraform_state" {
  count  = var.create_bootstrap_resources ? 1 : 0
  bucket = aws_s3_bucket.terraform_state[0].id

  versioning_configuration {
    status = "Enabled"
  }

  # Ignore errors if bucket already exists
  lifecycle {
    ignore_changes = [bucket]
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  count  = var.create_bootstrap_resources ? 1 : 0
  bucket = aws_s3_bucket.terraform_state[0].id

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
  count  = var.create_bootstrap_resources ? 1 : 0
  bucket = aws_s3_bucket.terraform_state[0].id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true

  # Ignore errors if bucket already exists
  lifecycle {
    ignore_changes = [bucket]
  }
}

# Check if the DynamoDB table already exists with the current prefix
data "aws_dynamodb_table" "existing_terraform_locks" {
  count = var.create_bootstrap_resources ? 0 : 1
  name  = "${local.bucket_prefix}-terraform-locks"
}

# Check if the DynamoDB table already exists with the invoiceapidev prefix (for backward compatibility)
data "aws_dynamodb_table" "existing_terraform_locks_api" {
  count = (!var.create_bootstrap_resources && var.environment == "dev" && local.bucket_prefix != "invoiceapidev") ? 1 : 0
  name  = "invoiceapidev-terraform-locks"
}


resource "aws_dynamodb_table" "terraform_locks" {
  count        = var.create_bootstrap_resources ? 1 : 0
  name         = local.use_api_prefix_for_dynamodb ? "invoiceapidev-terraform-locks" : "${local.bucket_prefix}-terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  lifecycle {
    # Ignore errors related to table already existing
    ignore_changes = [name, id]
    # Prevent errors when table already exists
    create_before_destroy = false
  }

  tags = {
    Name        = "${local.bucket_prefix}-terraform-locks"
    Environment = "All"
    Project     = var.project_name
  }
}

output "s3_bucket_name" {
  value = var.create_bootstrap_resources ? aws_s3_bucket.terraform_state[0].id : (
  try(data.aws_s3_bucket.existing_terraform_state[0].id,
    try(data.aws_s3_bucket.existing_terraform_state_api[0].id, "${local.bucket_prefix}-terraform-state")
  )
  )
  description = "The name of the S3 bucket"
}

output "dynamodb_table_name" {
  value = var.create_bootstrap_resources ? aws_dynamodb_table.terraform_locks[0].id : (
  try(data.aws_dynamodb_table.existing_terraform_locks[0].name,
    try(data.aws_dynamodb_table.existing_terraform_locks_api[0].name, "${local.bucket_prefix}-terraform-locks")
  )
  )
  description = "The name of the DynamoDB table"
}