data "aws_caller_identity" "current" {}

resource "aws_iam_oidc_provider" "github" {
  url = "https://${var.github_oidc_provider_url}" # e.g., "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com"
  ]

  thumbprint_list = ["1b511abead59c6ce207077c0bf0e0043b1382612"] # Replace with current GitHub OIDC thumbprint if it changes
}

# IAM Role for GitHub Actions (ECR Push, general AWS access)
resource "aws_iam_role" "github_actions_role" {
  name = var.github_actions_role_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = aws_iam_oidc_provider.github.arn
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringLike = {
            "${var.github_oidc_provider_url}:sub" : "repo:${var.github_repository}:*"
            # Add branch condition if needed: "${var.github_oidc_provider_url}:ref": "refs/heads/main"
          }
        }
      }
    ]
  })

  tags = {
    Name    = var.github_actions_role_name
    Project = var.project_name
  }
}

resource "aws_iam_policy" "github_actions_ecr_policy" {
  name        = "${var.project_name}-GitHubActionsECRPolicy"
  description = "Policy for GitHub Actions to push to ECR and list repositories"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow",
        Action   = [
          "ecr:GetAuthorizationToken"
        ],
        Resource = "*"
      },
      {
        Effect   = "Allow",
        Action   = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetRepositoryPolicy",
          "ecr:DescribeRepositories",
          "ecr:ListImages",
          "ecr:DescribeImages",
          "ecr:BatchGetImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "ecr:PutImage"
        ],
        Resource = "arn:aws:ecr:${var.aws_region}:${data.aws_caller_identity.current.account_id}:repository/${var.ecr_repository_api_name}"
        # Add other ECR repositories if needed
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "github_actions_ecr_attach" {
  role       = aws_iam_role.github_actions_role.name
  policy_arn = aws_iam_policy.github_actions_ecr_policy.arn
}

# IAM Role for GitHub Actions (ECS Deployments)
resource "aws_iam_role" "github_actions_ecs_deploy_role" {
  name = var.github_actions_ecs_deploy_role_name
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = aws_iam_oidc_provider.github.arn
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringLike = {
            "${var.github_oidc_provider_url}:sub" : "repo:${var.github_repository}:ref:refs/heads/main" # Only main branch
          }
        }
      }
    ]
  })
  tags = {
    Name    = var.github_actions_ecs_deploy_role_name
    Project = var.project_name
  }
}

resource "aws_iam_policy" "github_actions_ecs_deploy_policy" {
  name        = "${var.project_name}-GitHubActionsECSDeployPolicy"
  description = "Policy for GitHub Actions to deploy to ECS and related resources"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      { # Permissions for docker ecs compose up (CloudFormation)
        Effect = "Allow",
        Action = [
          "cloudformation:DescribeStacks",
          "cloudformation:DescribeStackEvents",
          "cloudformation:DescribeStackResources",
          "cloudformation:CreateStack",
          "cloudformation:UpdateStack",
          "cloudformation:DeleteStack",
          "cloudformation:ValidateTemplate"
        ],
        Resource = "arn:aws:cloudformation:${var.aws_region}:${data.aws_caller_identity.current.account_id}:stack/${var.project_name}-*"
      },
      { # General ECS permissions
        Effect = "Allow",
        Action = [
          "ecs:CreateCluster", # If docker ecs compose needs to create it (though we create it via TF)
          "ecs:DescribeClusters",
          "ecs:CreateService",
          "ecs:UpdateService",
          "ecs:DeleteService",
          "ecs:DescribeServices",
          "ecs:RegisterTaskDefinition",
          "ecs:DescribeTaskDefinition",
          "ecs:DeregisterTaskDefinition",
          "ecs:ListTasks",
          "ecs:DescribeTasks",
          "ecs:RunTask",
          "ecs:StopTask",
          "ecs:TagResource"
        ],
        Resource = "*" # Scope down if possible
      },
      { # IAM PassRole for ECS tasks
        Effect   = "Allow",
        Action   = "iam:PassRole",
        Resource = [
          aws_iam_role.ecs_task_execution_role.arn,
          # Add task roles for specific services if they need more permissions
        ]
      },
      { # ECR read access for ECS
        Effect = "Allow",
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ],
        Resource = "*"
      },
      { # ELB permissions for service creation/update
        Effect = "Allow",
        Action = [
          "elasticloadbalancing:DescribeLoadBalancers",
          "elasticloadbalancing:DescribeTargetGroups",
          "elasticloadbalancing:DescribeListeners",
          "elasticloadbalancing:CreateRule",
          "elasticloadbalancing:DeleteRule",
          "elasticloadbalancing:ModifyRule",
          "elasticloadbalancing:DescribeRules",
          "elasticloadbalancing:SetRulePriorities",
          "elasticloadbalancing:RegisterTargets",
          "elasticloadbalancing:DeregisterTargets",
          "elasticloadbalancing:CreateTargetGroup",
          "elasticloadbalancing:DeleteTargetGroup",
          "elasticloadbalancing:ModifyTargetGroupAttributes"
        ],
        Resource = "*" # Scope down if possible
      },
      { # Secrets Manager read access
        Effect   = "Allow",
        Action   = "secretsmanager:GetSecretValue",
        Resource = "arn:aws:secretsmanager:${var.aws_region}:${data.aws_caller_identity.current.account_id}:secret:*"
        # Scope down to specific secrets
      },
      { # EFS permissions if docker ecs compose manages EFS volumes (though we create EFS via TF)
        Effect = "Allow",
        Action = [
          "elasticfilesystem:DescribeFileSystems",
          "elasticfilesystem:DescribeMountTargets"
        ],
        Resource = "*"
      }
      # Add other permissions as needed by docker ecs compose up (e.g., CloudWatch Logs, IAM for creating service-linked roles)
    ]
  })
}

resource "aws_iam_role_policy_attachment" "github_actions_ecs_deploy_attach" {
  role       = aws_iam_role.github_actions_ecs_deploy_role.name
  policy_arn = aws_iam_policy.github_actions_ecs_deploy_policy.arn
}


# ECS Task Execution Role (for ECS agent to pull images, publish logs)
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.project_name}-ecs-task-execution-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action    = "sts:AssumeRole",
      Effect    = "Allow",
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
  tags = {
    Name    = "${var.project_name}-ecs-task-execution-role"
    Project = var.project_name
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# (Optional) ECS Task Role for API service (if it needs to interact with other AWS services)
# resource "aws_iam_role" "api_task_role" {
#   name = "${var.project_name}-api-task-role"
#   assume_role_policy = jsonencode({
#     Version = "2012-10-17",
#     Statement = [{
#       Action    = "sts:AssumeRole",
#       Effect    = "Allow",
#       Principal = { Service = "ecs-tasks.amazonaws.com" }
#     }]
#   })
# }
# Add policies to api_task_role for S3, etc. if needed
