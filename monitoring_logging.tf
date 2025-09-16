# Production-Ready Logging and Audit Trail Configuration
# monitoring_logging.tf

# Enhanced CloudWatch Log Groups with structured logging
resource "aws_cloudwatch_log_group" "monitoring_audit_logs" {
  name              = "/aws/monitoring/${var.project_name}/audit"
  retention_in_days = var.log_retention_days * 2 # Keep audit logs longer
  kms_key_id        = aws_kms_key.backup_key.arn

  tags = {
    Name        = "${var.project_name}-monitoring-audit-logs"
    Project     = var.project_name
    Environment = var.environment
    Type        = "AuditLog"
    Compliance  = "Required"
  }
}

resource "aws_cloudwatch_log_group" "security_logs" {
  name              = "/aws/monitoring/${var.project_name}/security"
  retention_in_days = 90 # Security logs kept for 90 days minimum
  kms_key_id        = aws_kms_key.backup_key.arn

  tags = {
    Name        = "${var.project_name}-security-logs"
    Project     = var.project_name
    Environment = var.environment
    Type        = "SecurityLog"
    Compliance  = "Required"
  }
}

resource "aws_cloudwatch_log_group" "performance_logs" {
  name              = "/aws/monitoring/${var.project_name}/performance"
  retention_in_days = var.log_retention_days
  kms_key_id        = aws_kms_key.backup_key.arn

  tags = {
    Name        = "${var.project_name}-performance-logs"
    Project     = var.project_name
    Environment = var.environment
    Type        = "PerformanceLog"
  }
}

# CloudWatch Log Streams for different monitoring components
resource "aws_cloudwatch_log_stream" "prometheus_metrics" {
  name           = "prometheus-metrics-stream"
  log_group_name = aws_cloudwatch_log_group.performance_logs.name
}

resource "aws_cloudwatch_log_stream" "grafana_access" {
  name           = "grafana-access-stream"
  log_group_name = aws_cloudwatch_log_group.monitoring_audit_logs.name
}

resource "aws_cloudwatch_log_stream" "alertmanager_notifications" {
  name           = "alertmanager-notifications-stream"
  log_group_name = aws_cloudwatch_log_group.monitoring_audit_logs.name
}

resource "aws_cloudwatch_log_stream" "sonarqube_security" {
  name           = "sonarqube-security-stream"
  log_group_name = aws_cloudwatch_log_group.security_logs.name
}

# CloudWatch Insights Queries for monitoring
resource "aws_cloudwatch_query_definition" "monitoring_errors" {
  name = "${var.project_name}-monitoring-errors"

  log_group_names = [
    aws_cloudwatch_log_group.monitoring_audit_logs.name,
    aws_cloudwatch_log_group.security_logs.name,
    aws_cloudwatch_log_group.performance_logs.name
  ]

  query_string = <<EOF
fields @timestamp, @message, level, service
| filter level = "ERROR"
| sort @timestamp desc
| limit 100
EOF
}

resource "aws_cloudwatch_query_definition" "security_events" {
  name = "${var.project_name}-security-events"

  log_group_names = [
    aws_cloudwatch_log_group.security_logs.name
  ]

  query_string = <<EOF
fields @timestamp, @message, source_ip, user_id, action
| filter action like /login|logout|failed_auth|privilege_escalation/
| sort @timestamp desc
| limit 50
EOF
}

resource "aws_cloudwatch_query_definition" "performance_metrics" {
  name = "${var.project_name}-performance-metrics"

  log_group_names = [
    aws_cloudwatch_log_group.performance_logs.name
  ]

  query_string = <<EOF
fields @timestamp, @message, response_time, cpu_usage, memory_usage
| filter response_time > 1000 or cpu_usage > 80 or memory_usage > 85
| sort @timestamp desc
| limit 100
EOF
}

# CloudTrail for API-level audit logging
resource "aws_cloudtrail" "monitoring_api_trail" {
  name                          = "${var.project_name}-monitoring-api-trail"
  s3_bucket_name               = aws_s3_bucket.cloudtrail_logs.bucket
  include_global_service_events = true
  is_multi_region_trail        = true
  enable_logging               = true

  event_selector {
    read_write_type                 = "All"
    include_management_events       = true
    exclude_management_event_sources = ["kms.amazonaws.com", "rdsdata.amazonaws.com"]

    data_resource {
      type   = "AWS::S3::Object"
      values = ["${aws_s3_bucket.cloudtrail_logs.arn}/*"]
    }

    data_resource {
      type   = "AWS::ECS::Service"
      values = ["*"]
    }

    data_resource {
      type   = "AWS::EFS::FileSystem"
      values = ["*"]
    }
  }

  tags = {
    Name        = "${var.project_name}-monitoring-api-trail"
    Project     = var.project_name
    Environment = var.environment
    Purpose     = "MonitoringAudit"
  }
}

# S3 bucket for CloudTrail logs
resource "aws_s3_bucket" "cloudtrail_logs" {
  bucket = "${var.project_name}-monitoring-cloudtrail-${random_string.bucket_suffix.result}"

  tags = {
    Name        = "${var.project_name}-monitoring-cloudtrail"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "random_string" "bucket_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "aws_s3_bucket_versioning" "cloudtrail_logs_versioning" {
  bucket = aws_s3_bucket.cloudtrail_logs.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_encryption" "cloudtrail_logs_encryption" {
  bucket = aws_s3_bucket.cloudtrail_logs.id

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.backup_key.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }
}

resource "aws_s3_bucket_policy" "cloudtrail_logs_policy" {
  bucket = aws_s3_bucket.cloudtrail_logs.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AWSCloudTrailAclCheck"
        Effect = "Allow"
        Principal = {
          Service = "cloudtrail.amazonaws.com"
        }
        Action   = "s3:GetBucketAcl"
        Resource = aws_s3_bucket.cloudtrail_logs.arn
      },
      {
        Sid    = "AWSCloudTrailWrite"
        Effect = "Allow"
        Principal = {
          Service = "cloudtrail.amazonaws.com"
        }
        Action   = "s3:PutObject"
        Resource = "${aws_s3_bucket.cloudtrail_logs.arn}/*"
        Condition = {
          StringEquals = {
            "s3:x-amz-acl" = "bucket-owner-full-control"
          }
        }
      }
    ]
  })
}

# VPC Flow Logs for network monitoring
resource "aws_flow_log" "monitoring_vpc_flow_log" {
  iam_role_arn    = aws_iam_role.flow_log_role.arn
  log_destination = aws_cloudwatch_log_group.vpc_flow_logs.arn
  traffic_type    = "ALL"
  vpc_id          = aws_vpc.main.id

  tags = {
    Name        = "${var.project_name}-monitoring-vpc-flow-log"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_cloudwatch_log_group" "vpc_flow_logs" {
  name              = "/aws/vpc/${var.project_name}/flowlogs"
  retention_in_days = var.log_retention_days
  kms_key_id        = aws_kms_key.backup_key.arn

  tags = {
    Name        = "${var.project_name}-vpc-flow-logs"
    Project     = var.project_name
    Environment = var.environment
  }
}

# IAM role for VPC Flow Logs
resource "aws_iam_role" "flow_log_role" {
  name = "${var.project_name}-flow-log-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "vpc-flow-logs.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "${var.project_name}-flow-log-role"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_iam_role_policy" "flow_log_policy" {
  name = "${var.project_name}-flow-log-policy"
  role = aws_iam_role.flow_log_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}

# CloudWatch Alarms for log monitoring
resource "aws_cloudwatch_metric_alarm" "high_error_rate_logs" {
  alarm_name          = "${var.project_name}-high-error-rate-logs"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "ErrorCount"
  namespace           = "AWS/Logs"
  period              = "300"
  statistic           = "Sum"
  threshold           = "10"
  alarm_description   = "High error rate detected in monitoring logs"
  alarm_actions       = [aws_sns_topic.backup_notifications.arn]

  dimensions = {
    LogGroupName = aws_cloudwatch_log_group.monitoring_audit_logs.name
  }

  tags = {
    Name        = "${var.project_name}-high-error-rate-logs-alarm"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_cloudwatch_metric_alarm" "security_events_alarm" {
  alarm_name          = "${var.project_name}-security-events"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "SecurityEventCount"
  namespace           = "AWS/Logs"
  period              = "300"
  statistic           = "Sum"
  threshold           = "5"
  alarm_description   = "Suspicious security events detected"
  alarm_actions       = [aws_sns_topic.backup_notifications.arn]

  dimensions = {
    LogGroupName = aws_cloudwatch_log_group.security_logs.name
  }

  tags = {
    Name        = "${var.project_name}-security-events-alarm"
    Project     = var.project_name
    Environment = var.environment
  }
}