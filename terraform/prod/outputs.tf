output "instance_id" {
  description = "EC2 instance; open a shell with `aws ssm start-session --target <id>`."
  value       = module.ec2.instance_id
}

output "public_ip" {
  description = "Elastic IP; point the A records of app_host and api_host here."
  value       = module.ec2.public_ip
}

output "rds_address" {
  value = module.rds.address
}

output "name_prefix" {
  value = local.name
}

output "s3_base_url" {
  value = module.s3.base_url
}
