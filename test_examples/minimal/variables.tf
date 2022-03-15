variable "test_bucket_name" {
  description = "Name of test S3 bucket name."
  type        = string
}

variable "read_only_paths" {
  description = "List of paths/prefixes that should be attached to a read-only policy. Listed path(s) should omit the head bucket."
  type        = list(string)
  default     = []
}

variable "read_write_paths" {
  description = "List of paths/prefixes that should be attached to a read-write` policy. Listed path(s) should omit the head bucket."
  type        = list(string)
  default     = []
}
