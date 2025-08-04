/*
 * ECS Infrastructure Management:
 *
 * This Terraform configuration creates the ECS cluster, task definitions, services,
 * and CloudWatch log group. Previously, these resources were managed by the GitHub
 * Actions workflow using Docker Compose CLI for ECS, but now they are managed directly
 * by Terraform for better control and visibility.
 */

resource "aws_ecs_cluster" "main" {
  name = var.ecs_cluster_name

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name        = var.ecs_cluster_name
    Project     = var.project_name
    Environment = var.environment
  }
}

# CloudWatch Log Group for ECS tasks
resource "aws_cloudwatch_log_group" "ecs_logs" {
  name              = "/ecs/${var.project_name}" # Matches LOG_GROUP_NAME_ENV in GitHub Actions
  retention_in_days = 30                         # Adjust as needed

  tags = {
    Name        = "${var.project_name}-ecs-logs"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Using the ECS Task Execution Role defined in iam.tf

# ECS Task Role (for task-specific permissions)
resource "aws_iam_role" "ecs_task_role" {
  name = "${var.project_name}-ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-ecs-task-role"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Policy for accessing EFS
resource "aws_iam_policy" "ecs_efs_access_policy" {
  name        = "${var.project_name}-ecs-efs-access-policy"
  description = "Policy for ECS tasks to access EFS"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "elasticfilesystem:ClientMount",
          "elasticfilesystem:ClientWrite"
        ]
        Resource = [
          aws_efs_file_system.uploads.arn,
          aws_efs_file_system.n8n_data.arn
        ]
      }
    ]
  })
}

# Attach EFS access policy to the task role
resource "aws_iam_role_policy_attachment" "ecs_task_role_efs_policy" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.ecs_efs_access_policy.arn
}

# Policy for accessing Secrets Manager
resource "aws_iam_policy" "ecs_secrets_access_policy" {
  name        = "${var.project_name}-ecs-secrets-access-policy"
  description = "Policy for ECS tasks to access Secrets Manager"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = [
          aws_secretsmanager_secret.db_credentials.arn,
          aws_secretsmanager_secret.redis_config.arn,
          aws_secretsmanager_secret.rabbitmq_config.arn,
          aws_secretsmanager_secret.internal_api_key.arn,
          "${aws_secretsmanager_secret.internal_api_key.arn}-*", # Allow access to specific versions if needed
          aws_secretsmanager_secret.n8n_encryption_key.arn,
          "${aws_secretsmanager_secret.n8n_encryption_key.arn}-*" # Allow access to specific versions if needed
        ]
      }
    ]
  })
}

# Attach Secrets Manager access policy to the task role
resource "aws_iam_role_policy_attachment" "ecs_task_role_secrets_policy" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.ecs_secrets_access_policy.arn
}

# API Task Definition
resource "aws_ecs_task_definition" "api" {
  family                   = "${var.project_name}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"  # 0.5 vCPU
  memory                   = "1024" # 1 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "api"
      image     = "${aws_ecr_repository.api.repository_url}:latest" # Ensure this is dynamically updated or use a specific tag/digest
      essential = true

      portMappings = [
        {
          containerPort = var.app_port # Should match the port your Go app listens on
          hostPort      = var.app_port # Not strictly needed for Fargate with awsvpc, but doesn't hurt
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "APP_ENV"
          value = "production"
        },
        {
          name  = "APP_PORT" # Pass the app port as an environment variable
          value = tostring(var.app_port)
        },
        {
          name  = "UPLOADS_DIR"
          value = "/mnt/invoice_uploads" # Align with the EFS mount point's containerPath
        }
        # Add any other non-sensitive environment variables here
      ]

      secrets = [
        # Database Credentials (assuming aws_secretsmanager_secret.db_credentials stores a JSON with keys: host, port, username, password, dbname, sslmode)
        { name = "DB_HOST",     valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:host::" },
        { name = "DB_PORT",     valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:port::" },
        { name = "DB_USER",     valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:username::" },
        { name = "DB_PASSWORD", valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:password::" },
        { name = "DB_NAME",     valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:dbname::" },
        { name = "DB_SSLMODE",  valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:sslmode::" },

        # Redis Configuration (assuming aws_secretsmanager_secret.redis_config stores a JSON with keys: host, port, password, db)
        { name = "REDIS_HOST",     valueFrom = "${aws_secretsmanager_secret.redis_config.arn}:host::" },
        { name = "REDIS_PORT",     valueFrom = "${aws_secretsmanager_secret.redis_config.arn}:port::" },
        { name = "REDIS_PASSWORD", valueFrom = "${aws_secretsmanager_secret.redis_config.arn}:password::" }, # Ensure 'password' key exists in secret if used
        { name = "REDIS_DB",       valueFrom = "${aws_secretsmanager_secret.redis_config.arn}:db::" },       # Go app expects int, will be string here, parsed by Go app

        # RabbitMQ URL (assuming aws_secretsmanager_secret.rabbitmq_config stores a JSON with key: url)
        { name = "RABBITMQ_URL", valueFrom = "${aws_secretsmanager_secret.rabbitmq_config.arn}:url::" },

        # Internal API Key (assuming aws_secretsmanager_secret.internal_api_key stores a JSON with key: key)
        { name = "INTERNAL_API_KEY", valueFrom = "${aws_secretsmanager_secret.internal_api_key.arn}:key::" },

        # SMTP Configuration (assuming you create a secret for these, e.g., aws_secretsmanager_secret.smtp_config)
        { name = "SMTP_HOST", valueFrom = "${aws_secretsmanager_secret.smtp_config.arn}:host::" },
        { name = "SMTP_PORT", valueFrom = "${aws_secretsmanager_secret.smtp_config.arn}:port::" },
        { name = "SMTP_USER", valueFrom = "${aws_secretsmanager_secret.smtp_config.arn}:user::" },
        { name = "SMTP_PASSWORD", valueFrom = "${aws_secretsmanager_secret.smtp_config.arn}:password::" },
        { name = "SMTPSenderEmail", valueFrom = "${aws_secretsmanager_secret.smtp_config.arn}:sender_email::" },

        # JWT Secret (assuming you create a secret for this, e.g., aws_secretsmanager_secret.jwt_secret_config)
        # Example:
        # { name = "JWT_SECRET", valueFrom = "${aws_secretsmanager_secret.jwt_secret_config.arn}:secret::" }
      ]

      mountPoints = [
        {
          sourceVolume  = "uploads"
          containerPath = "/mnt/invoice_uploads" # This is where EFS 'uploads' volume will be mounted inside the container
          readOnly      = false
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_logs.name
          "awslogs-region"        = var.aws_region # Make sure var.aws_region is correctly defined
          "awslogs-stream-prefix" = "api"
        }
      }
    }
  ])

  volume {
    name = "uploads" # This name must match sourceVolume in mountPoints

    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.uploads.id
      transit_encryption = "ENABLED" # Recommended for security

      authorization_config { # Required if EFS is configured with IAM authorization for mount targets
        access_point_id = aws_efs_access_point.uploads.id # Optional, but recommended for fine-grained access control
        iam             = "ENABLED"                       # Set to "ENABLED" if using IAM authorization for the access point
      }
    }
  }

  tags = {
    Name        = "${var.project_name}-api-task"
    Project     = var.project_name
    Environment = var.environment
  }
}

# N8N Task Definition
resource "aws_ecs_task_definition" "n8n" {
  family                   = "${var.project_name}-n8n"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"  # 0.5 vCPU
  memory                   = "1024" # 1 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "n8n"
      image     = "n8nio/n8n:latest"
      essential = true

      portMappings = [
        {
          containerPort = 5678
          hostPort      = 5678
          protocol      = "tcp"
        }
      ]

      environment = [
        { name  = "N8N_HOST", value = "localhost" }, # n8n listens on localhost within its container
        { name  = "N8N_PORT", value = "5678" },
        { name  = "GENERIC_TIMEZONE", value = var.n8n_generic_timezone != "" ? var.n8n_generic_timezone : "UTC" }, # Example: make timezone configurable
        { name  = "WEBHOOK_URL", value = var.n8n_webhook_url }, # This should be the public URL of n8n
        { name  = "GO_API_BASE_URL", value = "http://${var.api_service_discovery_name}:${var.app_port}/api/v1" }, # Example using service discovery
        { name  = "DB_TYPE", value = "postgresdb"},
        { name  = "DB_POSTGRESDB_DATABASE",  value = var.db_name },
        { name  = "DB_POSTGRESDB_SSL_REJECT_UNAUTHORIZED", value = "false" }, # Allow self-signed certificates
        { name  = "DB_POSTGRESDB_SSL", value = "false" }, # Disable SSL for PostgreSQL connections
        { name  = "NODE_ENV", value = "production" } # Good practice for Node.js apps
      ]

      secrets = [
        { name = "GO_API_INTERNAL_KEY", valueFrom = "${aws_secretsmanager_secret.internal_api_key.arn}:key::" },
        { name = "N8N_ENCRYPTION_KEY", valueFrom = "${aws_secretsmanager_secret.n8n_encryption_key.arn}:key::" },
        { name = "DB_POSTGRESDB_HOST",     valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:host::" },
        { name = "DB_POSTGRESDB_PORT",     valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:port::" },
        { name = "DB_POSTGRESDB_USER",     valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:username::" },
        { name = "DB_POSTGRESDB_PASSWORD", valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:password::" },
      ]

      mountPoints = [
        {
          sourceVolume  = "n8n-data"
          containerPath = "/home/node/.n8n"
          readOnly      = false
        },
        {
          sourceVolume  = "uploads"
          containerPath = "/mnt/invoice_uploads"
          readOnly      = false
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_logs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "n8n"
        }
      }
    }
  ])

  volume {
    name = "n8n-data"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.n8n_data.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.n8n_data.id
        iam             = "ENABLED"
      }
    }
  }

  volume {
    name = "uploads"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.uploads.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.uploads.id
        iam             = "ENABLED"
      }
    }
  }

  tags = {
    Name        = "${var.project_name}-n8n-task"
    Project     = var.project_name
    Environment = var.environment
  }
}

# API Service
resource "aws_ecs_service" "api" {
  name            = "${var.project_name}-api-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  # Enable ECS managed tags for better cost tracking and resource identification
  enable_ecs_managed_tags = true
  propagate_tags          = "TASK_DEFINITION" # Propagates tags from task definition to tasks

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = var.app_port
  }

  depends_on = [
    aws_lb_listener.http,
    aws_iam_role_policy_attachment.ecs_task_execution_role_policy
  ]

  tags = {
    Name        = "${var.project_name}-api-service"
    Project     = var.project_name
    Environment = var.environment
  }
}

# N8N Service
resource "aws_ecs_service" "n8n" {
  name            = "${var.project_name}-n8n-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.n8n.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  enable_ecs_managed_tags = true
  propagate_tags          = "TASK_DEFINITION"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.n8n.arn
    container_name   = "n8n"
    container_port   = 5678
  }


  depends_on = [
    aws_lb_listener.http,
    aws_iam_role_policy_attachment.ecs_task_execution_role_policy
  ]

  tags = {
    Name        = "${var.project_name}-n8n-service"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Using the Target Groups defined in alb.tf

# Using the Security Group for ECS Tasks defined in security_groups.tf