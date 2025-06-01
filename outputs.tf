output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs of the private subnets"
  value       = aws_subnet.private[*].id
}

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = aws_lb.main.dns_name
}

output "ecr_api_repository_url" {
  description = "URL of the API ECR repository"
  value       = aws_ecr_repository.api.repository_url
}

output "ecs_cluster_name_output" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}


output "efs_uploads_id" {
  description = "ID of the EFS for uploads"
  value       = aws_efs_file_system.uploads.id
}
output "efs_n8n_data_id" {
  description = "ID of the EFS for N8N data"
  value       = aws_efs_file_system.n8n_data.id
}
output "efs_sonarqube_data_id" {
  description = "ID of the EFS for SonarQube data"
  value       = aws_efs_file_system.sonarqube_data.id
}
output "efs_sonarqube_logs_id" {
  description = "ID of the EFS for SonarQube logs"
  value       = aws_efs_file_system.sonarqube_logs.id
}
output "efs_sonarqube_extensions_id" {
  description = "ID of the EFS for SonarQube extensions"
  value       = aws_efs_file_system.sonarqube_extensions.id
}

output "db_credentials_secret_arn" {
  description = "ARN of the Secrets Manager secret for DB credentials"
  value       = aws_secretsmanager_secret.db_credentials.arn
}
output "redis_config_secret_arn" {
  description = "ARN of the Secrets Manager secret for Redis config"
  value       = aws_secretsmanager_secret.redis_config.arn
}
output "rabbitmq_config_secret_arn" {
  description = "ARN of the Secrets Manager secret for RabbitMQ config"
  value       = aws_secretsmanager_secret.rabbitmq_config.arn
}
output "internal_api_key_secret_arn" {
  description = "ARN of the Secrets Manager secret for Internal API Key"
  value       = aws_secretsmanager_secret.internal_api_key.arn
}
output "n8n_encryption_key_secret_arn" {
  description = "ARN of the Secrets Manager secret for N8N Encryption Key"
  value       = aws_secretsmanager_secret.n8n_encryption_key.arn
}
output "vpc_id_secret_arn" {
  description = "ARN of the Secrets Manager secret for VPC ID"
  value       = aws_secretsmanager_secret.vpc_id.arn
}