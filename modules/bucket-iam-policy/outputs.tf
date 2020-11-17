output "ro_policy_arn" {
  value       = length(aws_iam_policy.ro_policy) > 0 ? aws_iam_policy.ro_policy[0].arn : ""
  description = "ARN assigned to read-only IAM policy."
}

output "rw_policy_arn" {
  value       = length(aws_iam_policy.rw_policy) > 0 ? aws_iam_policy.rw_policy[0].arn : ""
  description = "ARN assigned to read-write IAM policy."
}
