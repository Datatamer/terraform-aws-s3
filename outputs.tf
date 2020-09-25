output "bucket_name" {
  value       = module.encrypted-bucket.bucket_name
  description = "Name of S3 bucket created by encrypted-bucket module."
}

output "ro_policy_arn" {
  value       = module.bucket-iam-policy.ro_policy_arn
  description = "ARN assigned to read-only IAM policy created by iam-policy module."
}

output "rw_policy_arn" {
  value       = module.bucket-iam-policy.rw_policy_arn
  description = "ARN assigned to read-write IAM policy created by iam-policy module."
}
