variable "region" {
  type    = string
  default = "ap-southeast-2"
}

variable "app_host" {
  description = "Hostname of the React app, e.g. app.example.com."
  type        = string
  default     = "app.deadun.site"
}

variable "api_host" {
  description = "Hostname of the backend API, e.g. api.example.com."
  type        = string
  default     = "api.deadun.site"
}

variable "certbot_email" {
  description = "Email for Let's Encrypt registration and expiry notices."
  type        = string
  default     = "buiquangbach27@gmail.com"
}

variable "bucket_name" {
  type    = string
  default = "bach-van-wedding"
}
