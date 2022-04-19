#tfsec:ignore:aws-s3-enable-bucket-logging tfsec:ignore:aws-s3-enable-versioning
resource "aws_s3_bucket" "new_bucket" {
  bucket = var.bucket_name
  acl    = "private"

  force_destroy = var.force_destroy
  tags          = var.tags
}

# Bucket policy to enforce AES256 server-side-encryption
resource "aws_s3_bucket_policy" "sse_bucket_policy" {
  bucket = aws_s3_bucket.new_bucket.id
  policy = templatefile(
    "${path.module}/bucket-policy.json",
    {
      bucket_name   = aws_s3_bucket.new_bucket.id,
      arn_partition = var.arn_partition
    }
  )
}

# Sets S3 bucket default bucket encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "encryption_for_new_bucket" {
  bucket = aws_s3_bucket.new_bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "AES256"
    }
}
}
