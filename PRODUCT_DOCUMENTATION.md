# Product Documentation: InvoiceB2B Platform

## 1. Overview

InvoiceB2B is a comprehensive financial platform designed for Business-to-Business (B2B) transactions. The system is built as a containerized Go application with a backend infrastructure defined using Terraform, indicating a robust, cloud-native approach.

Its core functionalities revolve around invoicing, loan application processing, and payment management, supported by a full suite of services for user management, administration, security, and notifications.

## 2. Core Features

Based on the service and data transfer object (DTO) definitions, the platform includes the following key features:

- **Invoice Management**: Create, view, and manage invoices.
- **Loan Application System**: A complete workflow for submitting, processing, and managing loan applications. This appears to be partially automated via an `n8n` workflow.
- **Payment Processing**: Services for handling payments related to invoices or loans.
- **User & Admin Management**: Separate handlers and services for regular users and platform administrators, including role-based access control (RBAC) via policies.
- **Financial Institution Management**: Functionality for managing financial institutions, likely as part of the loan or payment ecosystem.
- **Security & Middleware**:
    - **Authentication**: JWT-based authentication (`jwt_service.go`).
    - **Authorization**: Policy-based authorization (`authorization_service.go`, `policy_repository.go`).
    - **CSRF Protection**: Middleware to prevent Cross-Site Request Forgery (`csrf_middleware.go`).
    - **Rate Limiting**: Middleware to protect against brute-force and denial-of-service attacks.
    - **Secure Secret Management**: Integration with AWS Secrets Manager, including encryption and key rotation (`aws_secrets_manager_service.go`, `secureconfig/`).
- **Reporting & Auditing**: Services for generating reports and logging user/system activities (`reporting_service.go`, `activity_log_service.go`).
- **Notifications**: Email notification system for communicating with users (`email_service.go`).

## 3. Architecture

The project is organized into two main architectural domains: the Go application and the cloud infrastructure.

### 3.1. Application Architecture (Go Backend)

The Go application follows a clean, layered architecture, promoting separation of concerns and maintainability. The logic flows from routes to the database through distinct layers located in the `internal/` directory.

1.  **Routes (`internal/routes`)**: Defines the API endpoints (e.g., `/invoices`, `/users`) and maps them to the appropriate handlers.
2.  **Middleware (`internal/middleware`)**: Intercepts incoming HTTP requests to perform cross-cutting tasks like authentication, authorization, logging, and security checks before they reach the handlers.
3.  **Handlers (`internal/handlers`)**: The controller layer. It parses requests, validates input using DTOs, calls the relevant business logic in the services, and formats the API response.
4.  **Services (`internal/services`)**: Contains the core business logic of the application. For example, `invoice_service.go` orchestrates all actions related to invoices.
5.  **Repositories (`internal/repositories`)**: The data access layer (DAL). It abstracts the database interactions, providing a clean API for services to query and manipulate data without knowing the underlying database implementation.
6.  **Models (`internal/models`)**: Defines the core data structures that represent database entities like `User`, `Invoice`, and `LoanApplication`.
7.  **DTOs (`internal/dtos`)**: Data Transfer Objects define the structure of data for API requests and responses, ensuring a clear contract between the client and the server.

### 3.2. Infrastructure Architecture (Terraform)

The entire cloud infrastructure is managed as code using **Terraform**. The `.tf` files at the root of the project define all the necessary cloud components, including:

-   **Networking**: A custom VPC (`vpc.tf`) with security groups (`security_groups.tf`).
-   **Compute**:
    -   **ECS (Elastic Container Service)** (`ecs.tf`): The Go application is deployed as a Docker container managed by ECS.
    -   **ECR (Elastic Container Registry)** (`ecr.tf`): Stores the Docker images for the application.
-   **Databases & Storage**:
    -   **RDS (Relational Database Service)** (`rds.tf`): A managed SQL database for the application's primary data.
    -   **ElastiCache** (`elasticache.tf`): An in-memory cache to improve performance.
    -   **EFS (Elastic File System)** (`efs.tf`): A shared file system, possibly for storing uploads or shared resources.
-   **Security & Configuration**:
    -   **IAM (Identity and Access Management)** (`iam.tf`): Defines roles and permissions for cloud resources.
    -   **AWS Secrets Manager** (`secret_manager.tf`): For storing sensitive data like database credentials and API keys.
-   **Workflow & Messaging**:
    -   **Amazon MQ** (`amazon_mq.tf`): A managed message broker, likely used for asynchronous communication between services.
    -   **n8n Workflow Automation**: A separate ECS service (`n8n_ecs.tf`) runs n8n, an open-source workflow automation tool. The `n8n_custom/` directory contains custom workflow definitions (`loan_application_workflow.json`).

## 4. Key Directory Breakdown

-   `./`: The root contains the Go application entrypoint (`main.go`), Docker configuration (`Dockerfile`, `docker-compose.yaml`), and Terraform infrastructure definitions (`*.tf`).
-   `internal/`: The heart of the Go application code, organized by architectural layer.
-   `config/`: Global application configuration loading.
-   `docs/`: OpenAPI/Swagger API documentation artifacts.
-   `documentation/`: Manual, high-level documentation on various aspects of the project like API flows and security implementations.
-   `.github/workflows/`: CI/CD pipelines defined for GitHub Actions, used for automated testing, building, and deployment.
-   `environments/`: Contains environment-specific configurations (e.g., for `dev` and `prod`).
-   `n8n_custom/`: Custom workflows for the n8n automation engine.
-   `terraform_test/`: Contains tests for the Terraform infrastructure code, ensuring reliability.

## 5. How to Run

### Local Development

The presence of `docker-compose.yaml` suggests the project can be run locally using Docker Compose. This is the recommended approach as it will likely spin up the Go application and its dependencies (like the RDS database).

```bash
# Start all services in the background
docker-compose up -d
```

Alternatively, the Go application can be run directly, though this may require manual configuration of environment variables for database connections and other services.

```bash
# Run the Go application
go run main.go
```

### Deployment

Deployments are handled automatically via the GitHub Actions workflows defined in `.github/workflows/`. These pipelines likely build the Docker image, push it to ECR, and then apply the Terraform configurations to update the ECS services.
