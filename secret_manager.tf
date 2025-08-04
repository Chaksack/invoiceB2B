resource "aws_secretsmanager_secret" "db_credentials" {
  name        = "${var.project_name}/db_credentials"
  description = "Database credentials for RDS"
  tags        = { Project = var.project_name }
}

resource "aws_secretsmanager_secret_version" "db_credentials_version" {
  secret_id = aws_secretsmanager_secret.db_credentials.id
  secret_string = jsonencode({
    username = random_string.db_username.result
    password = random_password.db_password.result
    host     = aws_db_instance.main.address # From rds.tf
    port     = aws_db_instance.main.port    # From rds.tf
    dbname   = var.db_name
    sslmode  = "disable"
    # For SonarQube, it might use the same user or a different one.
    # If different, create another secret or add to this JSON.
    # For simplicity, assuming SonarQube uses the main db user for now.
    sonardbname = var.sonarqube_db_name
    # Add connection string for direct use in application
    connectionString = "postgres://${random_string.db_username.result}:${random_password.db_password.result}@${aws_db_instance.main.address}:${aws_db_instance.main.port}/${var.db_name}"
  })
}

resource "aws_secretsmanager_secret" "redis_config" {
  name        = "${var.project_name}/redis_config"
  description = "Redis connection details"
  tags        = { Project = var.project_name }
}

resource "aws_secretsmanager_secret_version" "redis_config_version" {
  secret_id = aws_secretsmanager_secret.redis_config.id
  secret_string = jsonencode({
    host = aws_elasticache_replication_group.main.primary_endpoint_address # From elasticache.tf
    port = aws_elasticache_replication_group.main.port                     # From elasticache.tf
    # password = var.redis_password # If you set one
    # Add URL for direct use in application
    url = "redis://${aws_elasticache_replication_group.main.primary_endpoint_address}:${aws_elasticache_replication_group.main.port}"
  })
}

resource "aws_secretsmanager_secret" "rabbitmq_config" {
  name        = "${var.project_name}/rabbitmq_config"
  description = "RabbitMQ connection details"
  tags        = { Project = var.project_name }
}
resource "aws_secretsmanager_secret_version" "rabbitmq_config_version" {
  secret_id = aws_secretsmanager_secret.rabbitmq_config.id
  secret_string = jsonencode({
    host     = aws_mq_broker.main.instances[0].endpoints[0] # Primary endpoint from amazon_mq.tf (adjust index/protocol if needed)
    port     = 5671                                         # AMQPS default, or 5672 for AMQP (adjust based on broker config)
    username = random_string.rabbitmq_username.result
    password = random_password.rabbitmq_password.result
    # Add URL for direct use in application
    url = "amqps://${random_string.rabbitmq_username.result}:${random_password.rabbitmq_password.result}@${aws_mq_broker.main.instances[0].endpoints[0]}:5671"
  })
}

resource "aws_secretsmanager_secret" "internal_api_key" {
  name        = "${var.project_name}/internal_api_key"
  description = "Internal API Key"
  tags        = { Project = var.project_name }
}

resource "random_password" "internal_api_key" {
  length           = 32
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_secretsmanager_secret_version" "internal_api_key_version" {
  secret_id = aws_secretsmanager_secret.internal_api_key.id
  secret_string = jsonencode({
    key = random_password.internal_api_key.result
  })
}

resource "random_password" "n8n_encryption_key" {
  length           = 32
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_secretsmanager_secret" "n8n_encryption_key" {
  name        = "${var.project_name}/n8n_encryption_key"
  description = "N8N Encryption Key"
  tags        = { Project = var.project_name }
}

resource "aws_secretsmanager_secret_version" "n8n_encryption_key_version" {
  secret_id = aws_secretsmanager_secret.n8n_encryption_key.id
  secret_string = jsonencode({
    key = random_password.n8n_encryption_key.result
  })
}

# Secret for VPC ID to be used by docker-compose.ecs.yml x-aws-vpc extension
resource "aws_secretsmanager_secret" "vpc_id" {
  name        = "${var.project_name}/vpc_id"
  description = "VPC ID for ECS deployment"
  tags        = { Project = var.project_name }
}

resource "aws_secretsmanager_secret_version" "vpc_id_version" {
  secret_id     = aws_secretsmanager_secret.vpc_id.id
  secret_string = aws_vpc.main.id
}

resource "aws_secretsmanager_secret" "smtp_config" {
  name        = "${var.project_name}/smtp_config"
  description = "SMTP configuration details for sending emails"
  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

resource "aws_secretsmanager_secret_version" "smtp_config_version" {
  secret_id = aws_secretsmanager_secret.smtp_config.id
  secret_string = jsonencode({
    host         = "smtp.gmail.com"
    port         = "465"
    user         = "andrew.sackey@syentia.io"
    password     = "xyspnvdkrwabrnmb"
    sender_email = "Profundr Invoice Financing <no-reply@profundr.io>"
  })
}