variable "name" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_id" {
  description = "Public subnet for the instance."
  type        = string
}

variable "app_host" {
  description = "Hostname of the React app, proxied to app_port."
  type        = string
}

variable "api_host" {
  description = "Hostname of the backend API, proxied to api_port."
  type        = string
}

variable "certbot_email" {
  description = "Email for Let's Encrypt registration and expiry notices."
  type        = string
}

variable "app_port" {
  description = "Localhost port the web container listens on."
  type        = number
  default     = 3000
}

variable "api_port" {
  description = "Localhost port the backend container listens on."
  type        = number
  default     = 4000
}

variable "grafana_secret_id" {
  description = "Secrets Manager secret (name or ARN) holding LOKI_URL, LOKI_USER and LOKI_TOKEN as env lines; read by Alloy on every start."
  type        = string
}

variable "instance_type" {
  description = "Must be arm64 (Graviton), e.g. t4g.small."
  type        = string
  default     = "t4g.small"
}

variable "ssh_cidr_blocks" {
  description = "IPv4 CIDRs allowed to reach SSH (port 22). Restrict to your own IP, e.g. [\"203.0.113.7/32\"]; empty closes the port."
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "root_volume_size" {
  type    = number
  default = 20
}

variable "secret_arns" {
  description = "Secrets Manager secrets the instance may read."
  type        = list(string)
  default     = []
}

variable "tags" {
  type    = map(string)
  default = {}
}
