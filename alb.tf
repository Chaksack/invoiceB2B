resource "aws_lb" "main" {
  name               = "${var.project_name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = aws_subnet.public[*].id # ALB in public subnets

  enable_deletion_protection = false # Enabled for production

  tags = {
    Name        = "${var.project_name}-alb"
    Project     = var.project_name
    Environment = "production"
  }
}

# Target Group for API Service
resource "aws_lb_target_group" "api" {
  name        = "${var.project_name}-api-tg"
  port        = var.app_port
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip" # For Fargate

  health_check {
    enabled             = true
    path                = "/api/" # Your API health check endpoint
    protocol            = "HTTP"
    matcher             = "200" # Expect HTTP 200 for healthy
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  tags = {
    Name        = "${var.project_name}-api-tg"
    Project     = var.project_name
    Environment = "production"
  }
}

# Target Group for N8N Service
resource "aws_lb_target_group" "n8n" {
  name        = "${var.project_name}-n8n-tg"
  port        = 5678 # N8N default port
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    enabled  = true
    path     = "/" # N8N health check endpoint
    protocol = "HTTP"
    matcher  = "200-399" # N8N might redirect, so wider range
  }
  tags = { Name = "${var.project_name}-n8n-tg" }
}

# Target Group for SonarQube Service
resource "aws_lb_target_group" "sonarqube" {
  name        = "${var.project_name}-sq-tg"
  port        = 9000 # SonarQube default port
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    enabled  = true
    path     = "/api/system/health" # SonarQube health check
    protocol = "HTTP"
    matcher  = "200"
  }
  tags = { Name = "${var.project_name}-sq-tg" }
}


resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.main.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  # Using a self-signed certificate for testing
  # In production, you should use a proper certificate from ACM
  certificate_arn = aws_acm_certificate.self_signed.arn

  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Path not found"
      status_code  = "404"
    }
  }
}

# Create a self-signed certificate for testing
resource "aws_acm_certificate" "self_signed" {
  private_key      = tls_private_key.self_signed.private_key_pem
  certificate_body = tls_self_signed_cert.self_signed.cert_pem

  lifecycle {
    create_before_destroy = true
  }
}

resource "tls_private_key" "self_signed" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "self_signed" {
  private_key_pem = tls_private_key.self_signed.private_key_pem

  subject {
    common_name  = "example.com"
    organization = "Example Organization"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

# Listener Rule for API Service (e.g., path-based or host-based)
resource "aws_lb_listener_rule" "api_rule" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }

  condition {
    path_pattern {
      values = ["/api/*"] # Route /api/* to the API service
    }
  }
  # Or use host_header condition if you have different domains
  # condition {
  #   host_header {
  #     values = ["api.yourdomain.com"]
  #   }
  # }
}

# Listener Rule for N8N Service
resource "aws_lb_listener_rule" "n8n_rule" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 110
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.n8n.arn
  }
  condition {
    path_pattern { values = ["/n8n/*"] } # Example path
  }
  # Or host_header: condition { host_header { values = ["n8n.yourdomain.com"] } }
}

# Listener Rule for SonarQube Service
resource "aws_lb_listener_rule" "sonarqube_rule" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 120
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.sonarqube.arn
  }
  condition {
    path_pattern { values = ["/sonarqube/*"] } # Example path
  }
  # Or host_header: condition { host_header { values = ["sonarqube.yourdomain.com"] } }
}
