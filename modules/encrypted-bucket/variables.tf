variable "bucket_name" {
  description = "Name of S3 bucket to create."
  type        = string
}

variable "additional_tags" {
  type        = map(string)
  description = "Additional tags to be attached to the S3 bucket."
  default     = {}
}

variable "force_destroy" {
  type        = bool
  description = <<EOF
  A boolean that indicates all objects (including any locked objects) should be deleted from the
  bucket so that the bucket can be destroyed without error. These objects are not recoverable.
  EOF
  default     = true
}
