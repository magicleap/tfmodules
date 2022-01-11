variable "message" {
  description = "will simply be displayed as an output"
  type        = string
  default     = "hello world"
}

output "result" {
  value = var.message
}
