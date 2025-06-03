resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-db-subnet-group"
  subnet_ids = aws_subnet.private[*].id # RDS in private subnets

  tags = {
    Name    = "${var.project_name}-db-subnet-group"
    Project = var.project_name
  }
}

# Generate a random username for RDS
resource "random_string" "db_username" {
  length  = 16
  special = false
  numeric = false
  upper   = false
}

# Generate a random password for RDS
resource "random_password" "db_password" {
  length           = 24
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_db_instance" "main" {
  identifier_prefix      = "${var.project_name}-db-" # Prefix for the instance identifier
  allocated_storage      = 20
  storage_type           = "gp2"
  engine                 = "postgres"
  engine_version         = "15"          # Match your local version or desired version
  instance_class         = "db.t3.micro" # Choose appropriate instance size
  db_name                = var.db_name   # Main application DB
  username               = random_string.db_username.result
  password               = random_password.db_password.result
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = "default.postgres15"
  skip_final_snapshot    = false # Enabled for production
  publicly_accessible    = false
  multi_az               = true # Enabled for production HA

  tags = {
    Name        = "${var.project_name}-main-db"
    Project     = var.project_name
    Environment = "production"
  }
}
