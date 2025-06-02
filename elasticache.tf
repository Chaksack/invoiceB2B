resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.project_name}-cache-subnet-group"
  subnet_ids = aws_subnet.private[*].id # ElastiCache in private subnets
}

resource "aws_elasticache_replication_group" "main" {
  replication_group_id = "${var.project_name}-redis-repl-group"
  description          = "ElastiCache Redis cluster for ${var.project_name}"
  node_type            = "cache.t3.micro" # Choose appropriate node type
  num_cache_clusters   = 2                # Minimum of 2 nodes required when automatic_failover_enabled is true
  # For multi-AZ replication, set num_node_groups > 0 and automatic_failover_enabled = true
  # num_node_groups             = 1 # Number of shards
  # replicas_per_node_group     = 1 # Number of read replicas per shard
  engine                     = "redis"
  engine_version             = "7.0" # Match your local version or desired version
  port                       = 6379
  parameter_group_name       = "default.redis7"
  subnet_group_name          = aws_elasticache_subnet_group.main.name
  security_group_ids         = [aws_security_group.elasticache.id]
  automatic_failover_enabled = true # Enabled for production HA with replicas
  at_rest_encryption_enabled = true # Enabled for production
  transit_encryption_enabled = true # Enabled for production
  # auth_token                  = var.redis_password # If you need password authentication

  tags = {
    Name        = "${var.project_name}-redis"
    Project     = var.project_name
    Environment = "production"
  }
}
