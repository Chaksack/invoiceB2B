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

# ECS Task Execution Role
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.project_name}-ecs-task-execution-role"

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
    Name        = "${var.project_name}-ecs-task-execution-role"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Attach the AmazonECSTaskExecutionRolePolicy to the task execution role
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

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
          aws_secretsmanager_secret.n8n_encryption_key.arn
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
      image     = "${aws_ecr_repository.api.repository_url}:latest"
      essential = true

      portMappings = [
        {
          containerPort = var.app_port
          hostPort      = var.app_port
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "APP_ENV"
          value = var.environment
        }
      ]

      secrets = [
        {
          name      = "DB_CONNECTION_STRING"
          valueFrom = "${aws_secretsmanager_secret.db_credentials.arn}:connectionString::"
        },
        {
          name      = "REDIS_URL"
          valueFrom = "${aws_secretsmanager_secret.redis_config.arn}:url::"
        },
        {
          name      = "RABBITMQ_URL"
          valueFrom = "${aws_secretsmanager_secret.rabbitmq_config.arn}:url::"
        },
        {
          name      = "INTERNAL_API_KEY"
          valueFrom = "${aws_secretsmanager_secret.internal_api_key.arn}:key::"
        }
      ]

      mountPoints = [
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
          "awslogs-stream-prefix" = "api"
        }
      }
    }
  ])

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
        {
          name  = "N8N_HOST"
          value = "localhost"
        },
        {
          name  = "N8N_PORT"
          value = "5678"
        },
        {
          name  = "GENERIC_TIMEZONE"
          value = "UTC"
        },
        {
          name  = "WEBHOOK_URL"
          value = "http://localhost:5678/"
        },
        {
          name  = "GO_API_BASE_URL"
          value = "http://api:3000/api/v1"
        }
      ]

      secrets = [
        {
          name      = "GO_API_INTERNAL_KEY"
          valueFrom = "${aws_secretsmanager_secret.internal_api_key.arn}:key::"
        },
        {
          name      = "N8N_ENCRYPTION_KEY"
          valueFrom = "${aws_secretsmanager_secret.n8n_encryption_key.arn}:key::"
        }
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

# Target Group for API
resource "aws_lb_target_group" "api" {
  name        = "${var.project_name}-api-tg"
  port        = var.app_port
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    enabled             = true
    interval            = 30
    path                = "/health"
    port                = "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    matcher             = "200"
  }

  tags = {
    Name        = "${var.project_name}-api-tg"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Target Group for N8N
resource "aws_lb_target_group" "n8n" {
  name        = "${var.project_name}-n8n-tg"
  port        = 5678
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    enabled             = true
    interval            = 30
    path                = "/"
    port                = "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    matcher             = "200"
  }

  tags = {
    Name        = "${var.project_name}-n8n-tg"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Security Group for ECS Tasks
resource "aws_security_group" "ecs_tasks" {
  name        = "${var.project_name}-ecs-tasks-sg"
  description = "Allow inbound traffic to ECS tasks"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Allow traffic from ALB"
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.project_name}-ecs-tasks-sg"
    Project     = var.project_name
    Environment = var.environment
  }
}