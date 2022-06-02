output "test-bucket" {
  value = module.minimal.test-bucket
}

output "role_arn" {
  value = aws_iam_role.this.arn
}
