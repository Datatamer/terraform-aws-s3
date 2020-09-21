variable "bucket_name" {
  description = "Name of S3 bucket to create."
  type        = string
}

variable "additional_tags" {
  type        = map(string)
  description = "Additional tags to be attached to the S3 bucket."
  default     = {}
}
