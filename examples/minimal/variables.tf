variable "test_bucket_name" {
  description = "Name of test S3 bucket name."
  type        = string
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources created by this example."
  default = {
    Author      = "Tamr"
    Environment = "Example"
  }
}

variable "read_only_paths" {
  type    = list(string)
  default = ["path/to/ro-folder"]
}

variable "read_write_paths" {
  type    = list(string)
  default = ["path/to/rw-folder", "path/to/another-rw-folder"]
}
