terraform {
  backend "s3" {
    bucket         = "invoicefin-terraform-state"
    key            = "terraform.tfstate"  # Default state file path, override with -backend-config
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "${var.project_name}-terraform-locks"
  }
}