variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1" # Change to your desired region
}

variable "project_name" {
  description = "A short name for the project, used for naming resources"
  type        = string
  default     = "invoiceb2b"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "List of CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "List of CIDR blocks for private subnets (for ECS tasks, RDS, ElastiCache)"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24"]
}

variable "availability_zones" {
  description = "List of Availability Zones to use"
  type        = list(string)
  # Ensure these are valid for your chosen region
  default = ["us-east-1a", "us-east-1b"]
}

variable "github_oidc_provider_url" {
  description = "URL of the GitHub OIDC provider (e.g., token.actions.githubusercontent.com)"
  type        = string
}

variable "github_actions_role_name" {
  description = "Name for the IAM role assumed by GitHub Actions for ECR/general access"
  type        = string
  default     = "GitHubActionECRAccessRole"
}

variable "github_actions_ecs_deploy_role_name" {
  description = "Name for the IAM role assumed by GitHub Actions for ECS deployments"
  type        = string
  default     = "GitHubActionECSDeployRole"
}

variable "github_repository" {
  description = "GitHub repository (owner/repo) for OIDC trust policy"
  type        = string
  # Example: "your-github-username/your-repo-name"
}

variable "ecr_repository_api_name" {
  description = "Name for the ECR repository for the API service"
  type        = string
  default     = "invoice-api" # Corresponds to secrets.ECR_REPOSITORY_API
}

variable "ecs_cluster_name" {
  description = "Name for the ECS cluster"
  type        = string
  default     = "invoice-b2b-cluster" # Corresponds to secrets.ECS_CLUSTER_NAME
}

variable "app_port" {
  description = "Port the API application listens on inside the container"
  type        = number
  default     = 3000
}

variable "db_username" {
  description = "Username for the RDS database"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Password for the RDS database"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Name for the main application database in RDS"
  type        = string
  default     = "invoice_db"
}

variable "sonarqube_db_name" {
  description = "Name for the SonarQube database in RDS"
  type        = string
  default     = "sonarqube_db"
}

// Add more variables as needed for N8N, SonarQube, other secrets, etc.
variable "n8n_encryption_key" {
  description = "N8N encryption key"
  type        = string
  sensitive   = true
}

variable "internal_api_key" {
  description = "Internal API Key for service-to-service communication"
  type        = string
  sensitive   = true
}

variable "rabbitmq_user" {
  description = "RabbitMQ default username"
  type        = string
  default     = "guest"
  sensitive   = true
}

variable "rabbitmq_password" {
  description = "RabbitMQ default password"
  type        = string
  default     = "guest"
  sensitive   = true
}
