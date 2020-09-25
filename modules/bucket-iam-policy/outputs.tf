output "ro_policy_arn" {
  value       = length(local.ro_paths) > 0 ? aws_iam_policy.ro_policy[0].arn : ""
  description = "ARN assigned to read-only IAM policy."
}

output "rw_policy_arn" {
  value       = length(var.read_write_paths) > 0 ? aws_iam_policy.rw_policy[0].arn : ""
  description = "ARN assigned to read-write IAM policy."
}
