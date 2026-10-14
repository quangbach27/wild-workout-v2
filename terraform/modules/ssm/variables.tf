variable "name" {
  description = "Name prefix; secrets are named <name>/<secret_name>."
  type        = string
}

variable "secret_names" {
  description = "Secrets to create, e.g. [\"clerk-secret-key\"]."
  type        = set(string)
}

variable "recovery_window_in_days" {
  type    = number
  default = 7
}

variable "tags" {
  type    = map(string)
  default = {}
}
