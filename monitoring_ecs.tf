# monitoring_ecs.tf - Deploys SonarQube, Prometheus, Grafana, and Alertmanager to ECS

# EFS File Systems for Prometheus, Grafana, and Alertmanager
# Note: SonarQube EFS is already defined in efs.tf

resource "aws_efs_file_system" "prometheus_data" {
  creation_token = "${var.project_name}-prometheus-data-efs"
  encrypted      = true
  tags = {
    Name    = "${var.project_name}-prometheus-data"
    Project = var.project_name
  }
}

resource "aws_efs_mount_target" "prometheus_data" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.prometheus_data.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_access_point" "prometheus_data" {
  file_system_id = aws_efs_file_system.prometheus_data.id
  posix_user {
    gid = 65534 # nobody
    uid = 65534 # nobody
  }
  root_directory {
    path = "/prometheus-data"
    creation_info {
      owner_gid   = 65534
      owner_uid   = 65534
      permissions = "755"
    }
  }
  tags = {
    Name    = "${var.project_name}-prometheus-ap"
    Project = var.project_name
  }
}

resource "aws_efs_file_system" "grafana_data" {
  creation_token = "${var.project_name}-grafana-data-efs"
  encrypted      = true
  tags = {
    Name    = "${var.project_name}-grafana-data"
    Project = var.project_name
  }
}

resource "aws_efs_mount_target" "grafana_data" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.grafana_data.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_access_point" "grafana_data" {
  file_system_id = aws_efs_file_system.grafana_data.id
  posix_user {
    gid = 472 # grafana
    uid = 472 # grafana
  }
  root_directory {
    path = "/grafana-data"
    creation_info {
      owner_gid   = 472
      owner_uid   = 472
      permissions = "755"
    }
  }
  tags = {
    Name    = "${var.project_name}-grafana-ap"
    Project = var.project_name
  }
}

resource "aws_efs_file_system" "alertmanager_data" {
  creation_token = "${var.project_name}-alertmanager-data-efs"
  encrypted      = true
  tags = {
    Name    = "${var.project_name}-alertmanager-data"
    Project = var.project_name
  }
}

resource "aws_efs_mount_target" "alertmanager_data" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.alertmanager_data.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_access_point" "alertmanager_data" {
  file_system_id = aws_efs_file_system.alertmanager_data.id
  posix_user {
    gid = 65534 # nobody
    uid = 65534 # nobody
  }
  root_directory {
    path = "/alertmanager-data"
    creation_info {
      owner_gid   = 65534
      owner_uid   = 65534
      permissions = "755"
    }
  }
  tags = {
    Name    = "${var.project_name}-alertmanager-ap"
    Project = var.project_name
  }
}

# Access point for SonarQube (since we need to use the existing EFS)
resource "aws_efs_access_point" "sonarqube_data" {
  file_system_id = aws_efs_file_system.sonarqube_data.id
  posix_user {
    gid = 1000 # sonarqube
    uid = 1000 # sonarqube
  }
  root_directory {
    path = "/sonarqube-data"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "755"
    }
  }
  tags = {
    Name    = "${var.project_name}-sonarqube-data-ap"
    Project = var.project_name
  }
}

resource "aws_efs_access_point" "sonarqube_logs" {
  file_system_id = aws_efs_file_system.sonarqube_logs.id
  posix_user {
    gid = 1000 # sonarqube
    uid = 1000 # sonarqube
  }
  root_directory {
    path = "/sonarqube-logs"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "755"
    }
  }
  tags = {
    Name    = "${var.project_name}-sonarqube-logs-ap"
    Project = var.project_name
  }
}

resource "aws_efs_access_point" "sonarqube_extensions" {
  file_system_id = aws_efs_file_system.sonarqube_extensions.id
  posix_user {
    gid = 1000 # sonarqube
    uid = 1000 # sonarqube
  }
  root_directory {
    path = "/sonarqube-extensions"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "755"
    }
  }
  tags = {
    Name    = "${var.project_name}-sonarqube-extensions-ap"
    Project = var.project_name
  }
}

# ECR Repositories for monitoring services
resource "aws_ecr_repository" "sonarqube" {
  name                 = "${var.project_name}-sonarqube"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = {
    Name    = "${var.project_name}-sonarqube-ecr"
    Project = var.project_name
  }
}

resource "aws_ecr_repository" "prometheus" {
  name                 = "${var.project_name}-prometheus"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = {
    Name    = "${var.project_name}-prometheus-ecr"
    Project = var.project_name
  }
}

resource "aws_ecr_repository" "grafana" {
  name                 = "${var.project_name}-grafana"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = {
    Name    = "${var.project_name}-grafana-ecr"
    Project = var.project_name
  }
}

resource "aws_ecr_repository" "alertmanager" {
  name                 = "${var.project_name}-alertmanager"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
  tags = {
    Name    = "${var.project_name}-alertmanager-ecr"
    Project = var.project_name
  }
}

# Update ECS tasks security group to allow traffic to monitoring services
resource "aws_security_group_rule" "prometheus_ingress" {
  security_group_id        = aws_security_group.ecs_tasks.id
  type                     = "ingress"
  from_port                = 9090
  to_port                  = 9090
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.alb.id
  description              = "Allow traffic from ALB to Prometheus"
}

resource "aws_security_group_rule" "grafana_ingress" {
  security_group_id        = aws_security_group.ecs_tasks.id
  type                     = "ingress"
  from_port                = 3000
  to_port                  = 3000
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.alb.id
  description              = "Allow traffic from ALB to Grafana"
}

resource "aws_security_group_rule" "alertmanager_ingress" {
  security_group_id        = aws_security_group.ecs_tasks.id
  type                     = "ingress"
  from_port                = 9093
  to_port                  = 9093
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.alb.id
  description              = "Allow traffic from ALB to Alertmanager"
}

# ALB Target Groups for monitoring services (SonarQube already defined in alb.tf)
resource "aws_lb_target_group" "prometheus" {
  name        = "${var.project_name}-prometheus-tg"
  port        = 9090
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"
  health_check {
    enabled             = true
    path                = "/-/healthy"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
  tags = {
    Name    = "${var.project_name}-prometheus-tg"
    Project = var.project_name
  }
}

resource "aws_lb_target_group" "grafana" {
  name        = "${var.project_name}-grafana-tg"
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"
  health_check {
    enabled             = true
    path                = "/api/health"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
  tags = {
    Name    = "${var.project_name}-grafana-tg"
    Project = var.project_name
  }
}

resource "aws_lb_target_group" "alertmanager" {
  name        = "${var.project_name}-alertmanager-tg"
  port        = 9093
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"
  health_check {
    enabled             = true
    path                = "/-/healthy"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
  tags = {
    Name    = "${var.project_name}-alertmanager-tg"
    Project = var.project_name
  }
}

# ALB Listener Rules for monitoring services (SonarQube already defined in alb.tf)
resource "aws_lb_listener_rule" "prometheus_rule" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 130
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.prometheus.arn
  }
  condition {
    path_pattern {
      values = ["/prometheus/*"]
    }
  }
}

resource "aws_lb_listener_rule" "grafana_rule" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 140
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.grafana.arn
  }
  condition {
    path_pattern {
      values = ["/grafana/*"]
    }
  }
}

resource "aws_lb_listener_rule" "alertmanager_rule" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 150
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.alertmanager.arn
  }
  condition {
    path_pattern {
      values = ["/alertmanager/*"]
    }
  }
}

# SonarQube ECS Task Definition
resource "aws_ecs_task_definition" "sonarqube" {
  family                   = "${var.project_name}-sonarqube"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "2048" # 2 vCPU
  memory                   = "4096" # 4 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "sonarqube"
      image     = "${aws_ecr_repository.sonarqube.repository_url}:latest"
      essential = true

      portMappings = [
        {
          containerPort = 9000
          hostPort      = 9000
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "SONAR_JDBC_URL", value = "jdbc:postgresql://${aws_db_instance.main.endpoint}/${var.sonarqube_db_name}" },
        { name = "SONAR_JDBC_USERNAME", value = "sonarqube" },
        { name = "SONAR_JDBC_PASSWORD", value = "sonarqube_password" }, # Should use secrets manager in production
        { name = "SONAR_WEB_CONTEXT", value = "/sonarqube" },
        { name = "SONAR_WEB_HOST", value = "0.0.0.0" }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_logs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "sonarqube"
        }
      }

      mountPoints = [
        {
          sourceVolume  = "sonarqube-data"
          containerPath = "/opt/sonarqube/data"
          readOnly      = false
        },
        {
          sourceVolume  = "sonarqube-logs"
          containerPath = "/opt/sonarqube/logs"
          readOnly      = false
        },
        {
          sourceVolume  = "sonarqube-extensions"
          containerPath = "/opt/sonarqube/extensions"
          readOnly      = false
        }
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "curl -f http://localhost:9000/api/system/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 120 # SonarQube can take a while to start
      }
    }
  ])

  volume {
    name = "sonarqube-data"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.sonarqube_data.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.sonarqube_data.id
        iam             = "ENABLED"
      }
    }
  }

  volume {
    name = "sonarqube-logs"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.sonarqube_logs.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.sonarqube_logs.id
        iam             = "ENABLED"
      }
    }
  }

  volume {
    name = "sonarqube-extensions"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.sonarqube_extensions.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.sonarqube_extensions.id
        iam             = "ENABLED"
      }
    }
  }

  tags = {
    Name        = "${var.project_name}-sonarqube-task"
    Project     = var.project_name
    Environment = "production"
  }
}

# Prometheus ECS Task Definition
resource "aws_ecs_task_definition" "prometheus" {
  family                   = "${var.project_name}-prometheus"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024" # 1 vCPU
  memory                   = "2048" # 2 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "prometheus"
      image     = "${aws_ecr_repository.prometheus.repository_url}:latest"
      essential = true

      portMappings = [
        {
          containerPort = 9090
          hostPort      = 9090
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "PROMETHEUS_CONFIG_PATH", value = "/etc/prometheus/prometheus.yml" }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_logs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "prometheus"
        }
      }

      mountPoints = [
        {
          sourceVolume  = "prometheus-data"
          containerPath = "/prometheus"
          readOnly      = false
        }
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:9090/-/healthy || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 30
      }
    }
  ])

  volume {
    name = "prometheus-data"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.prometheus_data.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.prometheus_data.id
        iam             = "ENABLED"
      }
    }
  }

  tags = {
    Name        = "${var.project_name}-prometheus-task"
    Project     = var.project_name
    Environment = "production"
  }
}

# Grafana ECS Task Definition
resource "aws_ecs_task_definition" "grafana" {
  family                   = "${var.project_name}-grafana"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024" # 1 vCPU
  memory                   = "2048" # 2 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "grafana"
      image     = "${aws_ecr_repository.grafana.repository_url}:latest"
      essential = true

      portMappings = [
        {
          containerPort = 3000
          hostPort      = 3000
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "GF_SERVER_ROOT_URL", value = "https://${aws_lb.main.dns_name}/grafana" },
        { name = "GF_SERVER_SERVE_FROM_SUB_PATH", value = "true" },
        { name = "GF_SECURITY_ADMIN_USER", value = "admin" },
        { name = "GF_SECURITY_ADMIN_PASSWORD", value = "admin" } # Should use secrets manager in production
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_logs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "grafana"
        }
      }

      mountPoints = [
        {
          sourceVolume  = "grafana-data"
          containerPath = "/var/lib/grafana"
          readOnly      = false
        }
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 30
      }
    }
  ])

  volume {
    name = "grafana-data"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.grafana_data.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.grafana_data.id
        iam             = "ENABLED"
      }
    }
  }

  tags = {
    Name        = "${var.project_name}-grafana-task"
    Project     = var.project_name
    Environment = "production"
  }
}

# Alertmanager ECS Task Definition
resource "aws_ecs_task_definition" "alertmanager" {
  family                   = "${var.project_name}-alertmanager"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512" # 0.5 vCPU
  memory                   = "1024" # 1 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "alertmanager"
      image     = "${aws_ecr_repository.alertmanager.repository_url}:latest"
      essential = true

      portMappings = [
        {
          containerPort = 9093
          hostPort      = 9093
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "ALERTMANAGER_CONFIG_PATH", value = "/etc/alertmanager/alertmanager.yml" }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs_logs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "alertmanager"
        }
      }

      mountPoints = [
        {
          sourceVolume  = "alertmanager-data"
          containerPath = "/alertmanager"
          readOnly      = false
        }
      ]

      healthCheck = {
        command     = ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:9093/-/healthy || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 30
      }
    }
  ])

  volume {
    name = "alertmanager-data"
    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.alertmanager_data.id
      transit_encryption = "ENABLED"
      authorization_config {
        access_point_id = aws_efs_access_point.alertmanager_data.id
        iam             = "ENABLED"
      }
    }
  }

  tags = {
    Name        = "${var.project_name}-alertmanager-task"
    Project     = var.project_name
    Environment = "production"
  }
}

# ECS Services
resource "aws_ecs_service" "sonarqube" {
  name            = "${var.project_name}-sonarqube"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.sonarqube.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.sonarqube.arn
    container_name   = "sonarqube"
    container_port   = 9000
  }

  depends_on = [aws_lb_listener_rule.sonarqube_rule]

  tags = {
    Name        = "${var.project_name}-sonarqube-service"
    Project     = var.project_name
    Environment = "production"
  }
}

resource "aws_ecs_service" "prometheus" {
  name            = "${var.project_name}-prometheus"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.prometheus.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.prometheus.arn
    container_name   = "prometheus"
    container_port   = 9090
  }

  depends_on = [aws_lb_listener_rule.prometheus_rule]

  tags = {
    Name        = "${var.project_name}-prometheus-service"
    Project     = var.project_name
    Environment = "production"
  }
}

resource "aws_ecs_service" "grafana" {
  name            = "${var.project_name}-grafana"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.grafana.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.grafana.arn
    container_name   = "grafana"
    container_port   = 3000
  }

  depends_on = [aws_lb_listener_rule.grafana_rule]

  tags = {
    Name        = "${var.project_name}-grafana-service"
    Project     = var.project_name
    Environment = "production"
  }
}

resource "aws_ecs_service" "alertmanager" {
  name            = "${var.project_name}-alertmanager"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.alertmanager.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.alertmanager.arn
    container_name   = "alertmanager"
    container_port   = 9093
  }

  depends_on = [aws_lb_listener_rule.alertmanager_rule]

  tags = {
    Name        = "${var.project_name}-alertmanager-service"
    Project     = var.project_name
    Environment = "production"
  }
}

# Add IAM policy for EFS access for monitoring services
resource "aws_iam_policy" "monitoring_efs_access" {
  name        = "${var.project_name}-monitoring-efs-access"
  description = "Allow ECS tasks to access EFS for monitoring services"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "elasticfilesystem:ClientMount",
          "elasticfilesystem:ClientWrite"
        ]
        Effect = "Allow"
        Resource = [
          aws_efs_file_system.sonarqube_data.arn,
          aws_efs_file_system.sonarqube_logs.arn,
          aws_efs_file_system.sonarqube_extensions.arn,
          aws_efs_file_system.prometheus_data.arn,
          aws_efs_file_system.grafana_data.arn,
          aws_efs_file_system.alertmanager_data.arn
        ]
        Condition = {
          StringEquals = {
            "elasticfilesystem:AccessPointArn" = [
              aws_efs_access_point.sonarqube_data.arn,
              aws_efs_access_point.sonarqube_logs.arn,
              aws_efs_access_point.sonarqube_extensions.arn,
              aws_efs_access_point.prometheus_data.arn,
              aws_efs_access_point.grafana_data.arn,
              aws_efs_access_point.alertmanager_data.arn
            ]
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "monitoring_efs_access" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.monitoring_efs_access.arn
}