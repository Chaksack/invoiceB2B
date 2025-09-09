# Monitoring Services CI/CD Pipeline

This document provides information about the CI/CD pipeline for building and deploying the monitoring services (SonarQube, Prometheus, Grafana, and Alertmanager) to AWS ECR and ECS.

## Overview

The monitoring services CI/CD pipeline is implemented as a GitHub Actions workflow that:

1. Builds Docker images for each monitoring service
2. Pushes these images to their respective ECR repositories
3. Updates the ECS services to use the new images

The workflow is defined in `.github/workflows/monitoring-services.yaml`.

## Workflow Triggers

The workflow is triggered by:

- **Push to main or staging branches** (only when monitoring_ecs.tf or the workflow file itself changes)
- **Pull requests to main or staging branches** (only when monitoring_ecs.tf or the workflow file itself changes)
- **Manual trigger** via GitHub Actions UI with environment selection

## Service Images

The workflow builds and pushes the following images:

| Service | Base Image | ECR Repository |
|---------|------------|----------------|
| SonarQube | sonarqube:9.9-community | ${PROJECT_NAME}-sonarqube |
| Prometheus | prom/prometheus:v2.45.0 | ${PROJECT_NAME}-prometheus |
| Grafana | grafana/grafana:10.0.3 | ${PROJECT_NAME}-grafana |
| Alertmanager | prom/alertmanager:v0.25.0 | ${PROJECT_NAME}-alertmanager |

Where `${PROJECT_NAME}` is defined in the GitHub repository secrets or defaults to 'invoice'.

## Image Tagging

Each image is tagged with:

- `latest` - Always points to the most recent build
- `<version>` - The specific version of the base image (e.g., "9.9-community" for SonarQube)
- `<git-sha>` - The Git commit SHA that triggered the workflow

## ECS Deployment

After building and pushing the images, the workflow automatically updates the ECS services to use the new images by forcing a new deployment of each service:

- ${PROJECT_NAME}-sonarqube
- ${PROJECT_NAME}-prometheus
- ${PROJECT_NAME}-grafana
- ${PROJECT_NAME}-alertmanager

## Running the Workflow Manually

To run the workflow manually:

1. Go to the GitHub repository
2. Click on the "Actions" tab
3. Select "Monitoring Services CI/CD Pipeline" from the list of workflows
4. Click "Run workflow"
5. Select the branch and environment (dev, staging, or prod)
6. Click "Run workflow"

## Customizing the Workflow

### Changing Service Versions

To change the version of a service, update the corresponding environment variable in the workflow file:

```yaml
env:
  SONARQUBE_VERSION: "9.9-community"
  PROMETHEUS_VERSION: "v2.45.0"
  GRAFANA_VERSION: "10.0.3"
  ALERTMANAGER_VERSION: "v0.25.0"
```

### Adding Custom Configuration

To add custom configuration to a service, modify the corresponding Dockerfile section in the workflow. For example, to add a custom Prometheus configuration:

1. Create a `prometheus.yml` file in your repository
2. Uncomment and modify the COPY line in the Prometheus Dockerfile section:

```yaml
dockerfile: |
  FROM prom/prometheus:${{ env.PROMETHEUS_VERSION }}
  COPY prometheus.yml /etc/prometheus/prometheus.yml
  USER nobody
```

## Troubleshooting

If the workflow fails, check the following:

1. **AWS Credentials**: Ensure that the AWS credentials (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) are correctly set in the repository secrets.
2. **ECR Repositories**: Verify that the ECR repositories exist and are accessible with the provided credentials.
3. **ECS Services**: Confirm that the ECS services are correctly named and running.
4. **Logs**: Review the workflow logs for specific error messages.

## Related Resources

- [Terraform Configuration for Monitoring Services](../monitoring_ecs.tf)
- [Monitoring Services Deployment Guide](monitoring-services-readme.md)
- [Monitoring Services Implementation Summary](monitoring-implementation-summary.md)