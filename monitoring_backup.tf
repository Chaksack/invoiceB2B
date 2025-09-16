# Production-Ready Backup and Disaster Recovery for Monitoring Services
# monitoring_backup.tf

# AWS Backup Vault for monitoring services
resource "aws_backup_vault" "monitoring_backup_vault" {
  name        = "${var.project_name}-monitoring-backup-vault"
  kms_key_arn = aws_kms_key.backup_key.arn

  tags = {
    Name        = "${var.project_name}-monitoring-backup-vault"
    Project     = var.project_name
    Environment = var.environment
    Service     = "Monitoring"
  }
}

# KMS Key for backup encryption
resource "aws_kms_key" "backup_key" {
  description             = "KMS key for monitoring services backup encryption"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  tags = {
    Name        = "${var.project_name}-backup-key"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_kms_alias" "backup_key_alias" {
  name          = "alias/${var.project_name}-monitoring-backup"
  target_key_id = aws_kms_key.backup_key.key_id
}

# IAM Role for AWS Backup
resource "aws_iam_role" "backup_role" {
  name = "${var.project_name}-backup-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "backup.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-backup-role"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Attach AWS managed policy for EFS backups
resource "aws_iam_role_policy_attachment" "backup_policy_efs" {
  role       = aws_iam_role.backup_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSBackupServiceRolePolicyForBackup"
}

resource "aws_iam_role_policy_attachment" "backup_policy_restore" {
  role       = aws_iam_role.backup_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSBackupServiceRolePolicyForRestores"
}

# Backup Plan for Monitoring Services
resource "aws_backup_plan" "monitoring_backup_plan" {
  name = "${var.project_name}-monitoring-backup-plan"

  rule {
    rule_name         = "daily_backup_rule"
    target_vault_name = aws_backup_vault.monitoring_backup_vault.name
    schedule          = "cron(0 5 ? * * *)" # Daily at 5 AM UTC

    start_window = 60  # 1 hour
    completion_window = 300 # 5 hours

    recovery_point_tags = {
      BackupType  = "Daily"
      Project     = var.project_name
      Environment = var.environment
    }

    lifecycle {
      cold_storage_after = 30
      delete_after       = var.backup_retention_days
    }
  }

  rule {
    rule_name         = "weekly_backup_rule"
    target_vault_name = aws_backup_vault.monitoring_backup_vault.name
    schedule          = "cron(0 3 ? * SUN *)" # Weekly on Sunday at 3 AM UTC

    start_window = 60
    completion_window = 300

    recovery_point_tags = {
      BackupType  = "Weekly"
      Project     = var.project_name
      Environment = var.environment
    }

    lifecycle {
      cold_storage_after = 7
      delete_after       = 90 # Keep weekly backups for 3 months
    }
  }

  tags = {
    Name        = "${var.project_name}-monitoring-backup-plan"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Backup Selection for EFS File Systems
resource "aws_backup_selection" "monitoring_efs_backup" {
  count        = var.enable_backup ? 1 : 0
  iam_role_arn = aws_iam_role.backup_role.arn
  name         = "${var.project_name}-monitoring-efs-backup"
  plan_id      = aws_backup_plan.monitoring_backup_plan.id

  resources = [
    aws_efs_file_system.prometheus_data.arn,
    aws_efs_file_system.grafana_data.arn,
    aws_efs_file_system.alertmanager_data.arn,
    aws_efs_file_system.sonarqube_data.arn,
    aws_efs_file_system.sonarqube_logs.arn,
    aws_efs_file_system.sonarqube_extensions.arn
  ]

  selection_tag {
    type  = "STRINGEQUALS"
    key   = "Project"
    value = var.project_name
  }
}

# CloudWatch Log Groups retention configuration
resource "aws_cloudwatch_log_group" "monitoring_logs" {
  for_each          = toset(["prometheus", "grafana", "alertmanager", "sonarqube"])
  name              = "/aws/ecs/${var.project_name}-${each.key}"
  retention_in_days = var.log_retention_days
  kms_key_id        = aws_kms_key.backup_key.arn

  tags = {
    Name        = "${var.project_name}-${each.key}-logs"
    Project     = var.project_name
    Environment = var.environment
    Service     = title(each.key)
  }
}

# SNS Topic for backup notifications
resource "aws_sns_topic" "backup_notifications" {
  name              = "${var.project_name}-backup-notifications"
  kms_master_key_id = aws_kms_key.backup_key.id

  tags = {
    Name        = "${var.project_name}-backup-notifications"
    Project     = var.project_name
    Environment = var.environment
  }
}

# CloudWatch Alarms for backup failures
resource "aws_cloudwatch_metric_alarm" "backup_failure" {
  alarm_name          = "${var.project_name}-backup-failure"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "NumberOfBackupJobsFailed"
  namespace           = "AWS/Backup"
  period              = "300"
  statistic           = "Sum"
  threshold           = "0"
  alarm_description   = "This metric monitors backup job failures"
  alarm_actions       = [aws_sns_topic.backup_notifications.arn]

  dimensions = {
    BackupVaultName = aws_backup_vault.monitoring_backup_vault.name
  }

  tags = {
    Name        = "${var.project_name}-backup-failure-alarm"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Disaster Recovery: Cross-region EFS replication (optional)
resource "aws_efs_replication_configuration" "prometheus_dr" {
  count                       = var.enable_multi_az && var.environment == "prod" ? 1 : 0
  source_file_system_id       = aws_efs_file_system.prometheus_data.id
  destination {
    region = var.dr_region != "" ? var.dr_region : "us-west-2"
    availability_zone_name = null
  }

  tags = {
    Name        = "${var.project_name}-prometheus-dr-replication"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_efs_replication_configuration" "grafana_dr" {
  count                       = var.enable_multi_az && var.environment == "prod" ? 1 : 0
  source_file_system_id       = aws_efs_file_system.grafana_data.id
  destination {
    region = var.dr_region != "" ? var.dr_region : "us-west-2"
    availability_zone_name = null
  }

  tags = {
    Name        = "${var.project_name}-grafana-dr-replication"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Database backup policy (RDS automated backups are handled in rds.tf)
# Here we create additional manual snapshot for disaster recovery
resource "aws_db_snapshot" "monitoring_db_snapshot" {
  count                          = var.enable_backup && var.environment == "prod" ? 1 : 0
  db_instance_identifier         = aws_db_instance.main.id
  db_snapshot_identifier         = "${var.project_name}-monitoring-manual-snapshot-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"

  tags = {
    Name        = "${var.project_name}-monitoring-db-snapshot"
    Project     = var.project_name
    Environment = var.environment
    Type        = "Manual"
  }

  lifecycle {
    ignore_changes = [db_snapshot_identifier]
  }
}