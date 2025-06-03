terraform {
  backend "s3" {
    bucket         = "invoicefin-terraform-state"
    key            = "terraform.tfstate"  # Default state file path, override with -backend-config
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "invoicefin-terraform-locks"
  }
}