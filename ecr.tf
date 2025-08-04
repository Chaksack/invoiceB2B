resource "aws_ecr_repository" "api" {
  name                 = var.ecr_repository_api_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name        = var.ecr_repository_api_name
    Project     = var.project_name
    Environment = "production"
  }
}