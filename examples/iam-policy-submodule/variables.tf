variable "existing_bucket_name" {
  description = "Name of existing S3 bucket to create policies for."
  type        = string
}

variable "tags" {
  type        = map(string)
  description = "Additional tags for resources created by this example"
  default     = { 
    Author = "Tamr"
    Environment = "Example"
  }
}
