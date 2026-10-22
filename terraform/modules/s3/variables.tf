variable "bucket_name" {
  description = "Globally unique S3 bucket name"
  type        = string
}

variable "cors_allowed_origins" {
  description = "Allowed origins for browser requests"
  type        = list(string)
  default     = ["*"]
}

variable "force_destroy" {
  description = "Delete bucket even when it is not empty"
  type        = bool
  default     = false
}

variable "enable_versioning" {
  description = "Keep previous versions of uploaded files"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags applied to resources"
  type        = map(string)
  default     = {}
}
