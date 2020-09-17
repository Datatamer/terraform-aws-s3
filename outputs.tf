output "bucket_name" {
  value       = aws_s3_bucket.new_bucket.id
  description = "Name of S3 bucket that was created."
}

output "policy_arns" {
  value       = var.read_write_paths == [] ? list(aws_iam_policy.ro_policy.arn) : list(aws_iam_policy.ro_policy.arn, aws_iam_policy.rw_policy[0].arn)
  description = "List of ARNs assigned to bucket IAM policies."
}
