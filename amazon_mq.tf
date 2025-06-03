# Generate a random username for RabbitMQ
resource "random_string" "rabbitmq_username" {
  length  = 16
  special = false
  numeric = false
  upper   = false
}

# Generate a random password for RabbitMQ
resource "random_password" "rabbitmq_password" {
  length           = 24
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_mq_broker" "main" {
  broker_name         = "${var.project_name}-rabbitmq"
  engine_type         = "RabbitMQ"
  engine_version      = "3.13"                  # Check AWS console for latest supported versions
  host_instance_type  = "mq.m5.large"
  deployment_mode      = "CLUSTER_MULTI_AZ"
  publicly_accessible = false                     # Keep it private
  subnet_ids          = aws_subnet.private[*].id  # Using all private subnets for multi-AZ
  security_groups     = [aws_security_group.rabbitmq.id]
  auto_minor_version_upgrade = true


  user {
    username = random_string.rabbitmq_username.result
    password = random_password.rabbitmq_password.result
  }

  # maintenance_window_start_time { # Optional
  #   day_of_week = "SUNDAY"
  #   time_of_day = "03:00"
  #   time_zone   = "UTC"
  # }

  logs { # Optional: enable CloudWatch logs
    general = true
    audit   = false
  }

  tags = {
    Name        = "${var.project_name}-rabbitmq"
    Project     = var.project_name
    Environment = "production"
  }
}