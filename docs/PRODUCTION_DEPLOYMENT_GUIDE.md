# Production Deployment Guide - InvoiceB2B Monitoring Infrastructure

## Overview

This guide provides comprehensive instructions for deploying the production-ready monitoring infrastructure for the InvoiceB2B application. The monitoring stack includes Prometheus, Grafana, Alertmanager, and SonarQube with enterprise-grade security, scalability, and reliability features.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Pre-deployment Setup](#pre-deployment-setup)
3. [Security Configuration](#security-configuration)
4. [Deployment Process](#deployment-process)
5. [Post-deployment Verification](#post-deployment-verification)
6. [Operational Procedures](#operational-procedures)
7. [Troubleshooting](#troubleshooting)
8. [Maintenance](#maintenance)

## Prerequisites

### Infrastructure Requirements
- AWS Account with appropriate permissions
- Terraform >= 1.0
- AWS CLI configured with appropriate credentials
- Docker for building custom images
- Valid SSL certificate for HTTPS endpoints

### Access Requirements
- Admin access to AWS Console
- Slack workspace with webhook integration capability
- SMTP server access for email notifications
- GitHub/GitLab repository access for CI/CD

### Network Requirements
- VPC with public and private subnets
- NAT Gateway for outbound internet access from private subnets
- Application Load Balancer with SSL certificate
- Security groups configured for monitoring services

## Pre-deployment Setup

### 1. Environment Variables Configuration

Create a `.env` file with the following production variables:

```bash
# Core Configuration
export PROJECT_NAME="invoiceb2b"
export ENVIRONMENT="prod"
export AWS_REGION="us-east-1"
export DR_REGION="us-west-2"

# Monitoring Configuration
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
export MONITORING_SMTP_HOST="smtp.gmail.com"
export MONITORING_SMTP_PORT="587"
export MONITORING_SMTP_USER="alerts@yourdomain.com"
export MONITORING_SMTP_PASSWORD="your-app-password"
export MONITORING_FROM_EMAIL="alerts@invoiceb2b.com"
export MONITORING_ADMIN_EMAILS='["admin@invoiceb2b.com","ops@invoiceb2b.com"]'
export MONITORING_OPS_EMAILS='["ops@invoiceb2b.com"]'
export MONITORING_DBA_EMAILS='["dba@invoiceb2b.com"]'

# Production Flags
export ENABLE_MULTI_AZ="true"
export ENABLE_AUTO_SCALING="true"
export ENABLE_BACKUP="true"
export ENABLE_ENCRYPTION="true"
export BACKUP_RETENTION_DAYS="30"
export LOG_RETENTION_DAYS="30"
```

### 2. Terraform Variables Configuration

Update the production terraform.tfvars:

```hcl
# environments/prod/terraform.tfvars
project_name = "invoiceb2b"
environment = "prod"
aws_region = "us-east-1"

# Production settings
enable_multi_az = true
enable_auto_scaling = true
enable_backup = true
enable_encryption = true
backup_retention_days = 30
log_retention_days = 30

# Monitoring configuration
monitoring_admin_emails = ["admin@invoiceb2b.com", "ops@invoiceb2b.com"]
monitoring_ops_emails = ["ops@invoiceb2b.com"]
monitoring_dba_emails = ["dba@invoiceb2b.com"]

# Security
enable_secrets_manager = true
```

### 3. Secrets Management Setup

Store sensitive values in AWS Secrets Manager:

```bash
# Store Slack webhook URL
aws secretsmanager create-secret \
    --name "invoiceb2b/monitoring/slack_webhook" \
    --secret-string '{"webhook_url":"https://hooks.slack.com/services/YOUR/WEBHOOK/URL"}' \
    --region us-east-1

# Store SMTP configuration
aws secretsmanager create-secret \
    --name "invoiceb2b/monitoring/email_config" \
    --secret-string '{
        "smtp_host":"smtp.gmail.com",
        "smtp_port":"587",
        "smtp_user":"alerts@yourdomain.com",
        "smtp_password":"your-app-password",
        "from_email":"alerts@invoiceb2b.com",
        "admin_emails":["admin@invoiceb2b.com"],
        "ops_emails":["ops@invoiceb2b.com"],
        "dba_emails":["dba@invoiceb2b.com"]
    }' \
    --region us-east-1
```

## Security Configuration

### 1. IAM Permissions

Ensure the deployment user has the following managed policies:
- `PowerUserAccess` (or more restrictive custom policy)
- `IAMFullAccess` (for service roles creation)

### 2. VPC Security

Configure security groups with minimal required access:
- ALB: Port 443 (HTTPS) from 0.0.0.0/0
- ECS Tasks: Ports 3000, 9000, 9090, 9093 from ALB only
- RDS: Port 5432 from ECS tasks only
- ElastiCache: Port 6379 from ECS tasks only

### 3. Encryption

All resources use encryption at rest:
- EFS volumes encrypted with KMS
- RDS encrypted with KMS
- S3 buckets encrypted with KMS
- CloudWatch logs encrypted with KMS

## Deployment Process

### 1. Infrastructure Deployment

```bash
# Initialize Terraform
terraform init -backend-config="environments/prod/backend.tfvars"

# Plan deployment
terraform plan -var-file="environments/prod/terraform.tfvars" -out=prod.plan

# Apply infrastructure
terraform apply prod.plan
```

### 2. Container Images Build and Push

```bash
# Build and push Prometheus image
docker build -t ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-prometheus:latest \
    -f monitoring/prometheus/Dockerfile .
docker push ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-prometheus:latest

# Build and push Grafana image
docker build -t ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-grafana:latest \
    -f monitoring/grafana/Dockerfile .
docker push ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-grafana:latest

# Build and push Alertmanager image
docker build -t ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-alertmanager:latest \
    -f monitoring/alertmanager/Dockerfile .
docker push ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-alertmanager:latest

# Build and push SonarQube image
docker build -t ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-sonarqube:latest \
    -f monitoring/sonarqube/Dockerfile .
docker push ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/invoiceb2b-sonarqube:latest
```

### 3. Service Deployment Verification

```bash
# Check ECS services status
aws ecs describe-services \
    --cluster invoiceb2b-cluster \
    --services invoiceb2b-prometheus invoiceb2b-grafana invoiceb2b-alertmanager invoiceb2b-sonarqube

# Check ALB target health
aws elbv2 describe-target-health \
    --target-group-arn $(terraform output prometheus_target_group_arn)
```

## Post-deployment Verification

### 1. Service Health Checks

```bash
# Test Prometheus
curl -k https://your-alb-dns/prometheus/-/healthy

# Test Grafana
curl -k https://your-alb-dns/grafana/api/health

# Test Alertmanager
curl -k https://your-alb-dns/alertmanager/-/healthy

# Test SonarQube
curl -k https://your-alb-dns/sonarqube/api/system/health
```

### 2. Monitoring Verification

1. **Prometheus Targets**: Access Prometheus UI and verify all targets are up
2. **Grafana Dashboards**: Login to Grafana and verify dashboards are loading
3. **Alert Rules**: Check Prometheus rules are loaded and evaluating
4. **Alertmanager**: Verify alertmanager is receiving alerts from Prometheus

### 3. Integration Testing

```bash
# Send test alert to verify Slack integration
curl -X POST https://your-alb-dns/alertmanager/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '[{
    "labels": {
      "alertname": "TestAlert",
      "severity": "warning",
      "service": "test"
    },
    "annotations": {
      "summary": "Test alert for deployment verification",
      "description": "This is a test alert to verify monitoring deployment"
    }
  }]'
```

## Operational Procedures

### 1. Daily Operations

#### Morning Health Check
```bash
#!/bin/bash
# daily-health-check.sh

echo "=== Daily Monitoring Health Check ==="
echo "Date: $(date)"

# Check ECS services
echo "Checking ECS services..."
aws ecs describe-services --cluster invoiceb2b-cluster \
    --services invoiceb2b-prometheus invoiceb2b-grafana invoiceb2b-alertmanager invoiceb2b-sonarqube \
    --query 'services[*].[serviceName,runningCount,desiredCount]' \
    --output table

# Check ALB target health
echo "Checking ALB target health..."
for tg in $(aws elbv2 describe-target-groups --names invoiceb2b-prometheus-tg invoiceb2b-grafana-tg invoiceb2b-alertmanager-tg --query 'TargetGroups[*].TargetGroupArn' --output text); do
    aws elbv2 describe-target-health --target-group-arn $tg --query 'TargetHealthDescriptions[*].[Target.Id,TargetHealth.State]' --output table
done

# Check recent backup status
echo "Checking backup status..."
aws backup list-backup-jobs --by-backup-vault-name invoiceb2b-monitoring-backup-vault \
    --max-results 5 \
    --query 'BackupJobs[*].[BackupJobId,State,CreationDate]' \
    --output table
```

#### Weekly Maintenance
```bash
#!/bin/bash
# weekly-maintenance.sh

echo "=== Weekly Monitoring Maintenance ==="

# Update ECS services to latest task definition
echo "Updating ECS services..."
for service in prometheus grafana alertmanager sonarqube; do
    aws ecs update-service \
        --cluster invoiceb2b-cluster \
        --service invoiceb2b-${service} \
        --force-new-deployment
done

# Clean up old logs
echo "Cleaning up old CloudWatch logs..."
aws logs describe-log-groups --log-group-name-prefix "/aws/ecs/invoiceb2b" \
    --query 'logGroups[*].logGroupName' --output text | \
    xargs -I {} aws logs put-retention-policy --log-group-name {} --retention-in-days 30
```

### 2. Scaling Operations

#### Manual Scaling
```bash
# Scale up during high load
aws ecs update-service \
    --cluster invoiceb2b-cluster \
    --service invoiceb2b-prometheus \
    --desired-count 2

# Scale down during low load
aws ecs update-service \
    --cluster invoiceb2b-cluster \
    --service invoiceb2b-prometheus \
    --desired-count 1
```

#### Auto-scaling Configuration
Auto-scaling is enabled by default based on CPU utilization (70% threshold).

### 3. Backup and Recovery

#### Manual Backup
```bash
# Trigger manual backup
aws backup start-backup-job \
    --backup-vault-name invoiceb2b-monitoring-backup-vault \
    --resource-arn $(aws efs describe-file-systems --query 'FileSystems[?Name==`invoiceb2b-prometheus-data`].FileSystemArn' --output text) \
    --iam-role-arn $(terraform output backup_role_arn) \
    --backup-options WindowsVSS=disabled
```

#### Disaster Recovery
```bash
# Restore from backup (example)
aws backup start-restore-job \
    --recovery-point-arn "arn:aws:backup:region:account:recovery-point:backup-vault-name/recovery-point-id" \
    --iam-role-arn $(terraform output backup_role_arn) \
    --resource-type EFS
```

## Troubleshooting

### Common Issues

#### 1. Service Won't Start
```bash
# Check ECS service events
aws ecs describe-services --cluster invoiceb2b-cluster --services invoiceb2b-prometheus \
    --query 'services[0].events[0:5].[createdAt,message]' --output table

# Check task logs
aws logs tail /aws/ecs/invoiceb2b-prometheus --follow
```

#### 2. High Memory Usage
```bash
# Check task resource utilization
aws ecs describe-tasks \
    --cluster invoiceb2b-cluster \
    --tasks $(aws ecs list-tasks --cluster invoiceb2b-cluster --service-name invoiceb2b-prometheus --query 'taskArns[0]' --output text)
```

#### 3. Alert Not Firing
1. Check Prometheus rules: `https://your-alb-dns/prometheus/rules`
2. Check Alertmanager status: `https://your-alb-dns/alertmanager/#/status`
3. Verify metric collection: Query metrics in Prometheus UI

#### 4. Backup Failures
```bash
# Check backup job status
aws backup describe-backup-job --backup-job-id YOUR_JOB_ID

# Check IAM permissions
aws iam simulate-principal-policy \
    --policy-source-arn $(terraform output backup_role_arn) \
    --action-names backup:StartBackupJob \
    --resource-arns "arn:aws:efs:*:*:file-system/*"
```

## Maintenance

### Monthly Tasks
1. Review and update alert thresholds based on application performance
2. Check and rotate secrets in AWS Secrets Manager
3. Review CloudWatch costs and optimize log retention
4. Update container images to latest security patches
5. Review and update backup retention policies

### Quarterly Tasks
1. Disaster recovery testing
2. Security audit of IAM roles and policies
3. Performance review and capacity planning
4. Update and test runbooks
5. Review and update monitoring dashboards

### Annual Tasks
1. Infrastructure cost optimization review
2. Security penetration testing
3. Business continuity plan update
4. Monitoring strategy review and planning
5. Technology stack upgrade planning

## Emergency Procedures

### 1. Complete Monitoring Stack Failure
```bash
# Emergency recovery procedure
terraform apply -target=aws_ecs_service.prometheus -auto-approve
terraform apply -target=aws_ecs_service.grafana -auto-approve
terraform apply -target=aws_ecs_service.alertmanager -auto-approve
terraform apply -target=aws_ecs_service.sonarqube -auto-approve
```

### 2. Data Loss Recovery
1. Identify affected EFS volumes
2. Restore from AWS Backup using recovery point
3. Restart affected services
4. Verify data integrity

### 3. Security Incident Response
1. Isolate affected resources
2. Review CloudTrail logs for suspicious activity
3. Rotate all secrets and credentials
4. Update security groups to restrict access
5. Document incident and lessons learned

## Support Contacts

- **Infrastructure Team**: ops@invoiceb2b.com
- **Security Team**: security@invoiceb2b.com
- **On-call Engineer**: +1-XXX-XXX-XXXX

## Additional Resources

- [AWS ECS Documentation](https://docs.aws.amazon.com/ecs/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Alertmanager Documentation](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [SonarQube Documentation](https://docs.sonarqube.org/)