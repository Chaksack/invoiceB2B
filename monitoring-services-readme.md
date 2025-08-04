# Monitoring Services Deployment Guide

This document provides information about the monitoring services deployed to ECS using Terraform.

## Overview

The following monitoring services have been deployed to ECS:

1. **SonarQube** - Code quality and security analysis
2. **Prometheus** - Metrics collection and storage
3. **Grafana** - Metrics visualization and dashboards
4. **Alertmanager** - Alert handling and notifications

## Architecture

- All services are deployed as ECS Fargate tasks in the existing ECS cluster
- Each service has its own EFS volume for persistent storage
- Services are exposed through the Application Load Balancer (ALB) with path-based routing:
  - SonarQube: `/sonarqube/*`
  - Prometheus: `/prometheus/*`
  - Grafana: `/grafana/*`
  - Alertmanager: `/alertmanager/*`
- Services run in private subnets with no direct internet access
- Security groups control traffic flow between services

## Resource Allocation

| Service      | CPU (vCPU) | Memory (GB) | Port |
|--------------|------------|-------------|------|
| SonarQube    | 2          | 4           | 9000 |
| Prometheus   | 1          | 2           | 9090 |
| Grafana      | 1          | 2           | 3000 |
| Alertmanager | 0.5        | 1           | 9093 |

## Deployment Instructions

### Prerequisites

1. AWS CLI configured with appropriate credentials
2. Terraform installed (version 1.7.0 or later)
3. Docker installed (for building and pushing images)

### Step 1: Build and Push Docker Images

Before applying the Terraform configuration, you need to build and push Docker images to the ECR repositories:

```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com

# SonarQube
docker pull sonarqube:9.9-community
docker tag sonarqube:9.9-community <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-sonarqube:latest
docker push <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-sonarqube:latest

# Prometheus
docker pull prom/prometheus:v2.45.0
docker tag prom/prometheus:v2.45.0 <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-prometheus:latest
docker push <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-prometheus:latest

# Grafana
docker pull grafana/grafana:10.0.3
docker tag grafana/grafana:10.0.3 <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-grafana:latest
docker push <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-grafana:latest

# Alertmanager
docker pull prom/alertmanager:v0.25.0
docker tag prom/alertmanager:v0.25.0 <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-alertmanager:latest
docker push <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/invoice-alertmanager:latest
```

### Step 2: Apply Terraform Configuration

```bash
# Initialize Terraform (if not already done)
terraform init

# Plan the deployment
terraform plan -out=tfplan

# Apply the configuration
terraform apply tfplan
```

### Step 3: Configure Services

#### Prometheus Configuration

You'll need to create a Prometheus configuration file and deploy it to the container. Here's a basic example:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'ecs-tasks'
    ec2_sd_configs:
      - region: us-east-1
        port: 9090
    relabel_configs:
      - source_labels: [__meta_ec2_tag_Name]
        regex: .*
        action: keep
```

#### Alertmanager Configuration

Create an Alertmanager configuration file:

```yaml
# alertmanager.yml
global:
  resolve_timeout: 5m
  smtp_from: '${SMTP_FROM}'
  smtp_smarthost: '${SMTP_SMARTHOST}'
  smtp_auth_username: '${SMTP_AUTH_USERNAME}'
  smtp_auth_password: '${SMTP_AUTH_PASSWORD}'
  slack_api_url: '${SLACK_WEBHOOK_URL}'

route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'slack-email-notifications'
  routes:
    - match:
        severity: critical
      receiver: 'slack-email-notifications'
      continue: true
    - match:
        severity: warning
      receiver: 'slack-email-notifications'

receivers:
  - name: 'slack-email-notifications'
    slack_configs:
      - channel: '#alerts'
        send_resolved: true
        title: '[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] Monitoring Alert'
        text: >-
          {{ range .Alerts }}
            *Alert:* {{ .Labels.alertname }}{{ if .Labels.severity }} - {{ .Labels.severity }}{{ end }}
            *Description:* {{ .Annotations.description }}
            *Details:*
            {{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}`
            {{ end }}
          {{ end }}
    email_configs:
      - to: 'alerts@example.com'
        send_resolved: true
        html: |
          {{ range .Alerts }}
            <h3>{{ .Labels.alertname }}{{ if .Labels.severity }} - {{ .Labels.severity }}{{ end }}</h3>
            <p><strong>Description:</strong> {{ .Annotations.description }}</p>
            <p><strong>Details:</strong></p>
            <ul>
              {{ range .Labels.SortedPairs }}
                <li><strong>{{ .Name }}:</strong> {{ .Value }}</li>
              {{ end }}
            </ul>
          {{ end }}
```

#### Grafana Configuration

Grafana is configured with environment variables in the ECS task definition. The default admin username and password are both set to "admin". You should change this after the first login.

To add Prometheus as a data source in Grafana:
1. Log in to Grafana at https://your-alb-dns/grafana/
2. Go to Configuration > Data Sources
3. Add a new data source (Prometheus)
4. Set the URL to https://your-alb-dns/prometheus/

## Access URLs

After deployment, the services will be available at:

- SonarQube: https://your-alb-dns/sonarqube/
- Prometheus: https://your-alb-dns/prometheus/
- Grafana: https://your-alb-dns/grafana/
- Alertmanager: https://your-alb-dns/alertmanager/

## Security Considerations

- Default credentials are used in this deployment for simplicity. In a production environment, you should:
  - Use AWS Secrets Manager for all credentials
  - Implement proper authentication for all services
  - Consider using AWS Cognito or another identity provider for authentication
- The ALB is publicly accessible. Consider implementing IP restrictions or VPN access for production environments
- Implement proper backup strategies for the EFS volumes

### GitHub Secrets for Alertmanager

The Alertmanager configuration uses environment variables that should be set using GitHub secrets. Add the following secrets to your GitHub repository:

- `SMTP_FROM`: The email address to send alerts from (e.g., 'alertmanager@yourdomain.com')
- `SMTP_SMARTHOST`: The SMTP server address and port (e.g., 'smtp.gmail.com:587')
- `SMTP_AUTH_USERNAME`: The SMTP username for authentication
- `SMTP_AUTH_PASSWORD`: The SMTP password for authentication
- `SLACK_WEBHOOK_URL`: The Slack Incoming Webhook URL for sending alerts to Slack

These secrets will be used in the GitHub workflow and passed to the Alertmanager container as environment variables.

## Maintenance

### Updating Services

To update a service:

1. Build and push a new Docker image to ECR
2. Update the ECS task definition to use the new image
3. Deploy the updated task definition

### Scaling

To scale a service, update the `desired_count` parameter in the ECS service definition.

### Monitoring the Monitoring Services

It's recommended to set up CloudWatch alarms to monitor the health of these monitoring services themselves.