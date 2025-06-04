# Bastion Host for RabbitMQ Access
# This creates an EC2 instance in the public subnet that can be used to access RabbitMQ in the private subnet

# Security group for the bastion host
resource "aws_security_group" "bastion" {
  name        = "${var.project_name}-bastion-sg"
  description = "Security group for bastion host"
  vpc_id      = aws_vpc.main.id

  # Allow SSH access from anywhere (you may want to restrict this to your IP)
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Consider restricting this to your IP address
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.project_name}-bastion-sg"
    Project     = var.project_name
    Environment = "production"
  }
}

# Update RabbitMQ security group to allow access from bastion
resource "aws_security_group_rule" "rabbitmq_from_bastion_amqp" {
  type                     = "ingress"
  from_port                = 5672
  to_port                  = 5672
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rabbitmq.id
  source_security_group_id = aws_security_group.bastion.id
  description              = "Allow AMQP access from bastion host"
}

resource "aws_security_group_rule" "rabbitmq_from_bastion_amqps" {
  type                     = "ingress"
  from_port                = 5671
  to_port                  = 5671
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rabbitmq.id
  source_security_group_id = aws_security_group.bastion.id
  description              = "Allow AMQPS access from bastion host"
}

resource "aws_security_group_rule" "rabbitmq_from_bastion_mgmt" {
  type                     = "ingress"
  from_port                = 15672
  to_port                  = 15672
  protocol                 = "tcp"
  security_group_id        = aws_security_group.rabbitmq.id
  source_security_group_id = aws_security_group.bastion.id
  description              = "Allow management UI access from bastion host"
}

# Generate a new SSH key pair if no public key is provided
resource "tls_private_key" "bastion" {
  count     = var.bastion_public_key == "" ? 1 : 0
  algorithm = "RSA"
  rsa_bits  = 4096
}

# Create a key pair for SSH access
resource "aws_key_pair" "bastion" {
  key_name   = "${var.project_name}-bastion-key"
  public_key = var.bastion_public_key != "" ? var.bastion_public_key : tls_private_key.bastion[0].public_key_openssh

  tags = {
    Name        = "${var.project_name}-bastion-key"
    Project     = var.project_name
    Environment = "production"
  }
}

# Create the bastion host EC2 instance
resource "aws_instance" "bastion" {
  ami                         = var.bastion_ami # Amazon Linux 2 AMI, should be provided as a variable
  instance_type               = "t3.micro"      # Small instance type for cost efficiency
  key_name                    = aws_key_pair.bastion.key_name
  vpc_security_group_ids      = [aws_security_group.bastion.id]
  subnet_id                   = aws_subnet.public[0].id # Place in the first public subnet
  associate_public_ip_address = true

  root_block_device {
    volume_size = 8 # GB
    volume_type = "gp3"
  }

  tags = {
    Name        = "${var.project_name}-bastion"
    Project     = var.project_name
    Environment = "production"
  }
}

# Output the public IP of the bastion host
output "bastion_public_ip" {
  value       = aws_instance.bastion.public_ip
  description = "Public IP address of the bastion host"
}

# Output the private key if it was generated
output "bastion_private_key" {
  value       = var.bastion_public_key == "" ? tls_private_key.bastion[0].private_key_pem : "Using provided public key"
  description = "Private key for SSH access to the bastion host (only if generated)"
  sensitive   = true
}