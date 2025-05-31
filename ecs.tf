resource "aws_ecs_cluster" "main" {
  name = var.ecs_cluster_name

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name        = var.ecs_cluster_name
    Project     = var.project_name
    Environment = "production"
  }
}

# CloudWatch Log Group for ECS tasks
resource "aws_cloudwatch_log_group" "ecs_logs" {
  name              = "/ecs/${var.project_name}" # Matches LOG_GROUP_NAME_ENV in GitHub Actions
  retention_in_days = 30                          # Adjust as needed

  tags = {
    Name        = "${var.project_name}-ecs-logs"
    Project     = var.project_name
    Environment = "production"
  }
}
