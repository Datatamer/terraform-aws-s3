output "bucket_name" {
  value       = aws_s3_bucket.new_bucket.id
  description = "Name of S3 bucket that was created."
}
