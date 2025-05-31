resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-db-subnet-group"
  subnet_ids = aws_subnet.private[*].id # RDS in private subnets

  tags = {
    Name    = "${var.project_name}-db-subnet-group"
    Project = var.project_name
  }
}

resource "aws_db_instance" "main" {
  identifier_prefix      = "${var.project_name}-db-" # Prefix for the instance identifier
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "postgres"
  engine_version       = "15" # Match your local version or desired version
  instance_class       = "db.t3.micro" # Choose appropriate instance size
  db_name              = var.db_name # Main application DB
  username             = var.db_username
  password             = var.db_password
  db_subnet_group_name = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name = "default.postgres15"
  skip_final_snapshot  = true # Set to false for production
  publicly_accessible  = false
  multi_az             = false # Set to true for production HA

  tags = {
    Name        = "${var.project_name}-main-db"
    Project     = var.project_name
    Environment = "production"
  }
}