# Containers only: no aws_secretsmanager_secret_version, so values never enter Terraform state.
# Set a value out of band, e.g.
#   aws secretsmanager put-secret-value --secret-id <name>/<secret> --secret-string '...'
resource "aws_secretsmanager_secret" "this" {
  for_each = var.secret_names

  name                    = "${var.name}/${each.key}"
  recovery_window_in_days = var.recovery_window_in_days

  tags = var.tags
}
