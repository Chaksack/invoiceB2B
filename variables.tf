variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-east-1" # Change to your desired region
}

variable "project_name" {
  description = "A short name for the project, used for naming resources"
  type        = string
  default     = "smeloan"
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


variable "ecr_repository_api_name" {
  description = "Name for the ECR repository for the API service"
  type        = string
  default     = "smeloan-api"
}

variable "ecs_cluster_name" {
  description = "Name for the ECS cluster"
  type        = string
  default     = "smeloan-cluster" # Corresponds to secrets.ECS_CLUSTER_NAME, uses project_name prefix
}

variable "app_port" {
  description = "Port the API application listens on inside the container"
  type        = number
  default     = 3000
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

variable "n8n_db_name" {
  description = "Name for the n8n database in RDS"
  type        = string
  default     = "n8n_db"
}

// Add more variables as needed for N8N, SonarQube, other secrets, etc.



variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "staging"
}

variable "create_bootstrap_resources" {
  description = "Whether to create the bootstrap resources (S3 bucket and DynamoDB table)"
  type        = bool
  default     = false
}

variable "bucket_prefix_override" {
  description = "Override for the bucket prefix used for S3 and DynamoDB resources (useful for migration or special environments)"
  type        = string
  default     = null
}

variable "bastion_public_key" {
  description = "Public SSH key for accessing the bastion host. If not provided, a new key pair will be generated."
  type        = string
  default     = ""
}

variable "bastion_ami" {
  description = "AMI ID for the bastion host (Amazon Linux 2023)"
  type        = string
  default     = "ami-0230bd60aa48260c6"  # Amazon Linux 2023 AMI (us-east-1) - update as needed
}

variable "n8n_generic_timezone" {
  description = "Timezone for the N8N container"
  type        = string
  default     = "UTC"
}

variable "n8n_webhook_url" {
  description = "The public-facing base URL for N8N webhooks. Should include http/https."
  type        = string
  default     = "http://localhost:5678/"
}

variable "api_service_discovery_name" {
  description = "The service discovery name for the API service (e.g., used for internal communication from N8N)."
  type        = string
  default     = "api"
}

# RabbitMQ Configuration
variable "rabbitmq_user" {
  description = "Username for RabbitMQ. Should be overridden in environment-specific tfvars or use AWS Secrets Manager."
  type        = string
  default     = "guest"
  sensitive   = true
}

variable "rabbitmq_password" {
  description = "Password for RabbitMQ. Should be stored in AWS Secrets Manager for production environments."
  type        = string
  default     = ""
  sensitive   = true
}

# Security Configuration
variable "enable_secrets_manager" {
  description = "Whether to use AWS Secrets Manager for sensitive values like database passwords and API keys"
  type        = bool
  default     = true
}

variable "force_destroy_s3_buckets" {
  description = "Allow destruction of S3 buckets that contain objects. Use with caution in production."
  type        = bool
  default     = false
}

# SMTP Configuration Variables (Production-ready)
variable "smtp_host" {
  description = "SMTP server hostname for application email notifications"
  type        = string
  default     = "smtp.gmail.com"
}

variable "smtp_port" {
  description = "SMTP server port"
  type        = string
  default     = "587"
}

variable "smtp_user" {
  description = "SMTP username for authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "smtp_password" {
  description = "SMTP password for authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "smtp_sender_email" {
  description = "Email address used as sender for application emails"
  type        = string
  default     = "noreply@invoiceb2b.com"
}

# Monitoring Configuration Variables
variable "slack_webhook_url" {
  description = "Slack webhook URL for monitoring alerts"
  type        = string
  sensitive   = true
  default     = ""
}

variable "monitoring_smtp_host" {
  description = "SMTP server for monitoring alerts"
  type        = string
  default     = "smtp.gmail.com"
}

variable "monitoring_smtp_port" {
  description = "SMTP port for monitoring alerts"
  type        = string
  default     = "465"
}

variable "monitoring_smtp_user" {
  description = "SMTP username for monitoring alerts"
  type        = string
  sensitive   = true
  default     = "andrew.sackey@syentia.io"
}

variable "monitoring_smtp_password" {
  description = "SMTP password for monitoring alerts"
  type        = string
  sensitive   = true
  default     = "xyspnvdkrwabrnmb"
}

variable "monitoring_from_email" {
  description = "From email address for monitoring alerts"
  type        = string
  default     = "SmeLoan Financing <no-reply@profundr.io>"
}

variable "monitoring_admin_emails" {
  description = "List of admin email addresses for critical alerts"
  type        = list(string)
  default     = ["andrew.sackey@syentia.io", "admin@invoiceb2b.com"]
}

variable "monitoring_ops_emails" {
  description = "List of operations team email addresses"
  type        = list(string)
  default     = ["andrew.sackey@syentia.io"]
}

variable "monitoring_dba_emails" {
  description = "List of DBA email addresses for database alerts"
  type        = list(string)
  default     = ["andrew.sackey@syentia.io"]
}

# Production Deployment Variables
variable "enable_multi_az" {
  description = "Enable multi-AZ deployment for high availability"
  type        = bool
  default     = true
}

variable "enable_auto_scaling" {
  description = "Enable auto-scaling for ECS services"
  type        = bool
  default     = true
}

variable "enable_backup" {
  description = "Enable automated backups for databases and EFS"
  type        = bool
  default     = true
}

variable "backup_retention_days" {
  description = "Number of days to retain backups"
  type        = number
  default     = 30
}

variable "log_retention_days" {
  description = "CloudWatch log retention period in days"
  type        = number
  default     = 30
}

variable "enable_encryption" {
  description = "Enable encryption at rest for all storage services"
  type        = bool
  default     = true
}

variable "dr_region" {
  description = "Disaster recovery region for cross-region backups"
  type        = string
  default     = "us-west-2"
}