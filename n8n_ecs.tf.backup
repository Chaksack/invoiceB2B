resource "aws_ecs_task_definition" "n8n" {
  family                   = "${var.project_name}-n8n"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024" # 1 vCPU
  memory                   = "2048" # 2 GB
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "n8n"
      image     = "${aws_ecr_repository.n8n.repository_url}:latest" # Assuming ECR repo exists
      essential = true

      portMappings = [
        {
          containerPort = 5678
          hostPort      = 5678
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "N8N_PORT", value = "5678" },
        { name = "N8N_PROTOCOL", value = "https" },
        { name = "N8N_HOST", value = "${aws_lb.main.dns_name}" }, # Using ALB DNS name
        { name = "N8N_PATH", value = "/n8n/" }, # Match the ALB path pattern
        { name = "N8N_SECURE_COOKIE", value = "false" }, # Disable secure cookies
        { name = "DB_TYPE", value = "postgresdb" },
        { name = "DB_POSTGRESDB_DATABASE", value = var.n8n_db_name }, # Database for n8n
        { name = "DB_POSTGRESDB_PORT", value = "5432" },
        { name = "NODE_ENV", value = "production" },
        { name = "N8N_METRICS", value = "true" },
        { name = "N8N_DIAGNOSTICS_ENABLED", value = "false" },
        { name = "N8N_USER_MANAGEMENT_DISABLED", value = "false" },
        { name = "N8N_HIRING_BANNER_ENABLED", value = "false" },
        { name = "WEBHOOK_URL", value = "https://${aws_lb.main.dns_name}/n8n/" },
        { name = "EXECUTIONS_PROCESS", value = "main" },
        { name = "GENERIC_TIMEZONE", value = "UTC" }
      ]

      secrets = [
        {
          name      = "N8N_ENCRYPTION_KEY"
          valueFrom = aws_secretsmanager_secret.n8n_encryption_key.arn
        },
        {
          name      = "DB_POSTGRESDB_HOST"
          valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:host::"
        },
        {
          name      = "DB_POSTGRESDB_USER"
          valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:username::"
        },
        {
          name      = "DB_POSTGRESDB_PASSWORD"
          valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:password::"
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

      mountPoints = [
        {
          sourceVolume  = "n8n-data"
          containerPath = "/home/node/.n8n"
          readOnly      = false
        }
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "curl -f http://localhost:5678/healthz || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
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

  tags = {
    Name        = "${var.project_name}-n8n-task"
    Project     = var.project_name
    Environment = "production"
  }
}

# Create an EFS access point for n8n data
resource "aws_efs_access_point" "n8n_data" {
  file_system_id = aws_efs_file_system.n8n_data.id

  posix_user {
    gid = 1000 # node user in the n8n container
    uid = 1000
  }

  root_directory {
    path = "/n8n-data"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "755"
    }
  }

  tags = {
    Name    = "${var.project_name}-n8n-ap"
    Project = var.project_name
  }
}

# Create an ECR repository for the n8n custom image
resource "aws_ecr_repository" "n8n" {
  name                 = "${var.project_name}-n8n"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name    = "${var.project_name}-n8n-ecr"
    Project = var.project_name
  }
}

# Create the ECS service for n8n
resource "aws_ecs_service" "n8n" {
  name            = "${var.project_name}-n8n"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.n8n.arn
  desired_count   = 1
  launch_type     = "FARGATE"

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

  # Ensure the service depends on the ALB listener rule
  depends_on = [aws_lb_listener_rule.n8n_rule]

  tags = {
    Name        = "${var.project_name}-n8n-service"
    Project     = var.project_name
    Environment = "production"
  }
}

# IAM roles for ECS (if not already defined elsewhere)
resource "aws_iam_role" "ecs_execution_role" {
  name = "${var.project_name}-ecs-execution-role"

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
    Name    = "${var.project_name}-ecs-execution-role"
    Project = var.project_name
  }
}

resource "aws_iam_role_policy_attachment" "ecs_execution_role_policy" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Add permissions for Secrets Manager
resource "aws_iam_policy" "secrets_manager_access" {
  name        = "${var.project_name}-secrets-manager-access"
  description = "Allow ECS tasks to access Secrets Manager"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "secretsmanager:GetSecretValue",
        ]
        Effect   = "Allow"
        Resource = [
          aws_secretsmanager_secret.n8n_encryption_key.arn,
          aws_secretsmanager_secret.db_credentials.arn
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "secrets_manager_access" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = aws_iam_policy.secrets_manager_access.arn
}

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
    Name    = "${var.project_name}-ecs-task-role"
    Project = var.project_name
  }
}

# Add permissions for EFS
resource "aws_iam_policy" "efs_access" {
  name        = "${var.project_name}-efs-access"
  description = "Allow ECS tasks to access EFS"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "elasticfilesystem:ClientMount",
          "elasticfilesystem:ClientWrite"
        ]
        Effect   = "Allow"
        Resource = aws_efs_file_system.n8n_data.arn
        Condition = {
          StringEquals = {
            "elasticfilesystem:AccessPointArn" = aws_efs_access_point.n8n_data.arn
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "efs_access" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.efs_access.arn
}
