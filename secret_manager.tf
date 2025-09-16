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
    host         = var.smtp_host
    port         = var.smtp_port
    user         = var.smtp_user
    password     = var.smtp_password
    sender_email = var.smtp_sender_email
  })
}

# Monitoring secrets for production-ready alerting
resource "aws_secretsmanager_secret" "monitoring_slack_webhook" {
  name        = "${var.project_name}/monitoring/slack_webhook"
  description = "Slack webhook URL for monitoring alerts"
  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
    Service     = "Monitoring"
  }
}

resource "aws_secretsmanager_secret_version" "monitoring_slack_webhook_version" {
  secret_id = aws_secretsmanager_secret.monitoring_slack_webhook.id
  secret_string = jsonencode({
    webhook_url = var.slack_webhook_url
  })
}

resource "aws_secretsmanager_secret" "monitoring_email_config" {
  name        = "${var.project_name}/monitoring/email_config"
  description = "Email configuration for monitoring alerts"
  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
    Service     = "Monitoring"
  }
}

resource "aws_secretsmanager_secret_version" "monitoring_email_config_version" {
  secret_id = aws_secretsmanager_secret.monitoring_email_config.id
  secret_string = jsonencode({
    smtp_host     = var.monitoring_smtp_host
    smtp_port     = var.monitoring_smtp_port
    smtp_user     = var.monitoring_smtp_user
    smtp_password = var.monitoring_smtp_password
    from_email    = var.monitoring_from_email
    admin_emails  = var.monitoring_admin_emails
    ops_emails    = var.monitoring_ops_emails
    dba_emails    = var.monitoring_dba_emails
  })
}

# SonarQube database credentials
resource "aws_secretsmanager_secret" "sonarqube_db_credentials" {
  name        = "${var.project_name}/sonarqube/db_credentials"
  description = "SonarQube database credentials"
  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
    Service     = "SonarQube"
  }
}

resource "random_password" "sonarqube_db_password" {
  length  = 32
  special = true
}

resource "aws_secretsmanager_secret_version" "sonarqube_db_credentials_version" {
  secret_id = aws_secretsmanager_secret.sonarqube_db_credentials.id
  secret_string = jsonencode({
    username = "sonarqube"
    password = random_password.sonarqube_db_password.result
    host     = aws_db_instance.main.address
    port     = aws_db_instance.main.port
    dbname   = var.sonarqube_db_name
    jdbc_url = "jdbc:postgresql://${aws_db_instance.main.address}:${aws_db_instance.main.port}/${var.sonarqube_db_name}"
  })
}

# Grafana admin credentials
resource "aws_secretsmanager_secret" "grafana_admin_credentials" {
  name        = "${var.project_name}/grafana/admin_credentials"
  description = "Grafana admin user credentials"
  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
    Service     = "Grafana"
  }
}

resource "random_password" "grafana_admin_password" {
  length  = 16
  special = true
}

resource "aws_secretsmanager_secret_version" "grafana_admin_credentials_version" {
  secret_id = aws_secretsmanager_secret.grafana_admin_credentials.id
  secret_string = jsonencode({
    admin_user     = "admin"
    admin_password = random_password.grafana_admin_password.result
  })
}