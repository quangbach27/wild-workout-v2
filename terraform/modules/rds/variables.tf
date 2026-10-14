variable "name" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  description = "Private subnets (two AZs, AWS requirement) for the DB subnet group."
  type        = list(string)
}

variable "allowed_security_group_ids" {
  description = "Security groups allowed to connect on the Postgres port."
  type        = list(string)
}

variable "db_name" {
  type    = string
  default = "wildworkout"
}

variable "username" {
  type    = string
  default = "wildworkout"
}

variable "engine_version" {
  type    = string
  default = "17"
}

# Defaults stay within the AWS RDS free tier: single-AZ db.t4g.micro, 20 GB, 1-day backups.
variable "instance_class" {
  type    = string
  default = "db.t4g.micro"
}

variable "allocated_storage" {
  type    = number
  default = 20
}

variable "multi_az" {
  type    = bool
  default = false
}

variable "backup_retention_period" {
  description = "Days to keep automated backups. Free-tier accounts reject values above 1."
  type        = number
  default     = 1
}

variable "deletion_protection" {
  type    = bool
  default = true
}

variable "tags" {
  type    = map(string)
  default = {}
}
