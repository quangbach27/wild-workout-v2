output "instance_id" {
  value = aws_instance.this.id
}

output "public_ip" {
  description = "Elastic IP; point the app and api DNS A records here."
  value       = aws_eip.this.public_ip
}

output "security_group_id" {
  value = aws_security_group.this.id
}
