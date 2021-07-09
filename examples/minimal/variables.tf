variable "test_bucket_name" {
  description = "Name of test S3 bucket name."
  type        = string
}

variable "additional_tags" {
  type        = map(string)
  description = "Additional tags for resources created by this example"
  default     = { 
    Author = "Tamr"
    Environment = "Example"
  }
}
