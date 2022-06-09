variable "bucket_name" {
  description = "Name of S3 bucket to create."
  type        = string
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources."
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

variable "arn_partition" {
  type        = string
  description = <<EOF
  The partition in which the resource is located. A partition is a group of AWS Regions.
  Each AWS account is scoped to one partition.
  The following are the supported partitions:
    aws -AWS Regions
    aws-cn - China Regions
    aws-us-gov - AWS GovCloud (US) Regions
  EOF
  default     = "aws"
}

variable "s3_bucket_logging" {
  description = "The name of S3 bucket where to store server access logs."
  type        = string
  default = ""
}
