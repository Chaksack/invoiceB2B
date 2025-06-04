resource "aws_efs_file_system" "uploads" {
  creation_token   = "${var.project_name}-uploads-efs"
  performance_mode = "generalPurpose"
  throughput_mode  = "bursting"
  encrypted        = true

  tags = {
    Name        = "${var.project_name}-uploads"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_efs_mount_target" "uploads" {
  count           = length(aws_subnet.private) # Mount target in each private subnet AZ
  file_system_id  = aws_efs_file_system.uploads.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_access_point" "uploads" {
  file_system_id = aws_efs_file_system.uploads.id

  posix_user {
    gid = 1000
    uid = 1000
  }

  root_directory {
    path = "/uploads"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "755"
    }
  }

  tags = {
    Name        = "${var.project_name}-uploads-ap"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_efs_file_system" "n8n_data" {
  creation_token   = "${var.project_name}-n8n-data-efs"
  performance_mode = "generalPurpose"
  throughput_mode  = "bursting"
  encrypted        = true

  tags = {
    Name        = "${var.project_name}-n8n-data"
    Project     = var.project_name
    Environment = var.environment
  }
}

resource "aws_efs_mount_target" "n8n_data" {
  count           = length(aws_subnet.private)
  file_system_id  = aws_efs_file_system.n8n_data.id
  subnet_id       = aws_subnet.private[count.index].id
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_access_point" "n8n_data" {
  file_system_id = aws_efs_file_system.n8n_data.id

  posix_user {
    gid = 1000
    uid = 1000
  }

  root_directory {
    path = "/n8n-data"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "755"
    }
  }

  tags = {
    Name        = "${var.project_name}-n8n-data-ap"
    Project     = var.project_name
    Environment = var.environment
  }
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