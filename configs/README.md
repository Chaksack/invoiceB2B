# Prometheus and Alertmanager Configuration for InvoiceB2B

## Overview

This directory contains the configuration files for Prometheus and Alertmanager to enable comprehensive monitoring and alerting for the InvoiceB2B application with Slack and email notifications.

## Files Structure

```
configs/
├── prometheus/
│   ├── prometheus.yml      # Main Prometheus configuration
│   └── alert_rules.yml     # Alert rules definitions
├── alertmanager/
│   └── alertmanager.yml    # Alertmanager routing and notification configuration
└── README.md              # This file
```

## Configuration Details

### Prometheus Configuration (`prometheus/prometheus.yml`)

- **Scrape Targets**: Configured to monitor the InvoiceB2B application, system metrics, and AWS ECS tasks
- **Alert Rules**: References `alert_rules.yml` for alert definitions
- **Alertmanager Integration**: Configured to send alerts to Alertmanager on port 9093

### Alert Rules (`prometheus/alert_rules.yml`)

Includes comprehensive monitoring rules for:
- **Application Availability**: Detects when the InvoiceB2B app is down
- **System Resources**: CPU usage, memory usage, disk space
- **Database Issues**: Connection errors and performance problems
- **HTTP Errors**: High error rates and response times
- **ECS Services**: Service health and task counts

### Alertmanager Configuration (`alertmanager/alertmanager.yml`)

Configured with multiple notification channels:

#### Slack Integration
- **Channels**:
  - `#alerts` - Default alerts
  - `#critical-alerts` - Critical severity alerts
  - `#monitoring` - Warning alerts
  - `#database-alerts` - Database-specific alerts
  - `#app-alerts` - Application-specific alerts

#### Email Integration
- **Critical Alerts**: Sent to admin and ops team
- **Database Alerts**: Sent to DBA and admin
- **SMTP Configuration**: Uses Gmail SMTP (configurable)

## Required Environment Variables

To enable Slack and email notifications, you need to set the following environment variables:

### Slack Configuration
```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
```

### Email Configuration
```bash
export SMTP_PASSWORD="your-email-app-password"
```

## Setup Instructions

### 1. Slack Webhook Setup

1. Go to your Slack workspace settings
2. Navigate to "Apps" → "Incoming Webhooks"
3. Create a new webhook for each channel you want to use:
   - `#alerts`
   - `#critical-alerts`
   - `#monitoring`
   - `#database-alerts`
   - `#app-alerts`
4. Copy the webhook URL and set it as `SLACK_WEBHOOK_URL`

### 2. Email Configuration

1. Configure your email provider (example uses Gmail):
   - Enable 2-factor authentication
   - Generate an app-specific password
   - Set the password as `SMTP_PASSWORD`

2. Update email addresses in `alertmanager.yml`:
   - Replace `admin@invoiceb2b.com` with your admin email
   - Replace `ops-team@invoiceb2b.com` with your ops team email
   - Replace `dba@invoiceb2b.com` with your database admin email

### 3. AWS Secrets Manager Integration (Recommended)

Instead of environment variables, store sensitive information in AWS Secrets Manager:

```bash
# Store Slack webhook URL
aws secretsmanager create-secret \
    --name "invoiceb2b/slack-webhook" \
    --secret-string "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"

# Store SMTP password
aws secretsmanager create-secret \
    --name "invoiceb2b/smtp-password" \
    --secret-string "your-email-app-password"
```

### 4. Update ECS Task Definitions

Update the Terraform configuration to mount these config files and set environment variables:

```hcl
# Add to prometheus task definition
environment = [
  { name = "PROMETHEUS_CONFIG_PATH", value = "/etc/prometheus/prometheus.yml" }
]

# Add to alertmanager task definition
environment = [
  { name = "ALERTMANAGER_CONFIG_PATH", value = "/etc/alertmanager/alertmanager.yml" },
  { name = "SLACK_WEBHOOK_URL", value = data.aws_secretsmanager_secret_version.slack_webhook.secret_string },
  { name = "SMTP_PASSWORD", value = data.aws_secretsmanager_secret_version.smtp_password.secret_string }
]
```

## Testing the Configuration

### 1. Validate Configuration Files

```bash
# Test Prometheus configuration
promtool check config configs/prometheus/prometheus.yml

# Test alert rules
promtool check rules configs/prometheus/alert_rules.yml

# Test Alertmanager configuration
amtool check-config configs/alertmanager/alertmanager.yml
```

### 2. Test Alert Routing

```bash
# Send a test alert to Alertmanager
amtool alert add --alertmanager.url=http://your-alertmanager:9093 \
  alertname="TestAlert" \
  severity="warning" \
  service="test" \
  summary="Test alert for configuration verification"
```

### 3. Monitor Logs

Check ECS CloudWatch logs for Prometheus and Alertmanager services to ensure:
- Configuration files are loaded correctly
- No parsing errors
- Successful connections to Slack and email services

## Alert Severity Levels

| Severity | Description | Notification Channels | Repeat Interval |
|----------|-------------|----------------------|----------------|
| **Critical** | Service down, system failures | Slack + Email | 1 hour |
| **Warning** | Performance issues, resource usage | Slack only | 4 hours |
| **Database** | Database-specific issues | Slack + Email (DBA) | 2 hours |
| **Application** | Application-specific alerts | Slack only | 2 hours |

## Customization

### Adding New Alert Rules

1. Edit `configs/prometheus/alert_rules.yml`
2. Add new rules following the existing pattern
3. Redeploy the Prometheus service

### Adding New Notification Channels

1. Edit `configs/alertmanager/alertmanager.yml`
2. Add new receivers and routing rules
3. Update environment variables for new integrations
4. Redeploy the Alertmanager service

### Modifying Alert Thresholds

Alert thresholds can be adjusted in `alert_rules.yml`:
- CPU usage: Currently set to 80%
- Memory usage: Currently set to 85%
- Disk space: Currently set to 85%
- HTTP error rate: Currently set to 5%
- Response time: Currently set to 2 seconds

## Security Considerations

1. **Secrets Management**: Use AWS Secrets Manager instead of environment variables in production
2. **Access Control**: Restrict access to configuration files containing sensitive information
3. **Network Security**: Ensure Alertmanager can only access required external services
4. **Webhook Validation**: Consider implementing webhook signature validation for Slack integration

## Troubleshooting

### Common Issues

1. **Alerts not firing**: Check Prometheus targets and rule evaluation
2. **Slack notifications not working**: Verify webhook URL and channel permissions
3. **Email notifications failing**: Check SMTP credentials and server connectivity
4. **Configuration errors**: Use validation tools mentioned in testing section

### Logs to Check

- Prometheus: `/var/log/prometheus/prometheus.log`
- Alertmanager: `/var/log/alertmanager/alertmanager.log`
- ECS CloudWatch logs for both services

## Maintenance

### Regular Tasks

1. **Review Alert Rules**: Monthly review of alert effectiveness and false positives
2. **Update Thresholds**: Adjust based on application performance patterns
3. **Test Notifications**: Quarterly testing of all notification channels
4. **Security Audit**: Regular review of secrets and access permissions

This configuration provides a robust monitoring and alerting solution with multiple notification channels to ensure critical issues are promptly addressed.