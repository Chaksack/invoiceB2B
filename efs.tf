resource "aws_efs_file_system" "uploads" {
  creation_token   = "${var.project_name}-uploads-efs"
  performance_mode = "generalPurpose"
  throughput_mode  = "bursting"
  encrypted        = true

  tags = {
    Name    = "${var.project_name}-uploads"
    Project = var.project_name
  }
}

resource "aws_efs_mount_target" "uploads" {
  count           = length(aws_subnet.private) # Mount target in each private subnet AZ
  file_system_id  = aws_efs_file_system.uploads.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_file_system" "n8n_data" {
  creation_token = "${var.project_name}-n8n-data-efs"
  tags           = { Name = "${var.project_name}-n8n-data", Project = var.project_name }
}
resource "aws_efs_mount_target" "n8n_data" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.n8n_data.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_file_system" "sonarqube_data" {
  creation_token = "${var.project_name}-sq-data-efs"
  tags           = { Name = "${var.project_name}-sq-data", Project = var.project_name }
}
resource "aws_efs_mount_target" "sonarqube_data" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.sonarqube_data.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_file_system" "sonarqube_logs" {
  creation_token = "${var.project_name}-sq-logs-efs"
  tags           = { Name = "${var.project_name}-sq-logs", Project = var.project_name }
}
resource "aws_efs_mount_target" "sonarqube_logs" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.sonarqube_logs.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_file_system" "sonarqube_extensions" {
  creation_token = "${var.project_name}-sq-ext-efs"
  tags           = { Name = "${var.project_name}-sq-ext", Project = var.project_name }
}
resource "aws_efs_mount_target" "sonarqube_extensions" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.sonarqube_extensions.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}
