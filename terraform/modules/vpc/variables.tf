variable "name" {
  description = "Name prefix for all resources."
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block of the VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones. The public subnet goes in the first; a private subnet is created in each."
  type        = list(string)

  validation {
    condition     = length(var.availability_zones) == length(var.private_subnet_cidrs)
    error_message = "availability_zones and private_subnet_cidrs must have the same length."
  }
}

variable "public_subnet_cidr" {
  description = "CIDR block of the public subnet (first availability zone)."
  type        = string
  default     = "10.0.1.0/24"
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks of the private subnets, one per availability zone."
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24"]
}

variable "tags" {
  description = "Tags applied to all resources."
  type        = map(string)
  default     = {}
}
