resource "aws_s3_bucket" "new_bucket" { #tfsec:ignore:AWS002, AWS077 and AWS098
  bucket = var.bucket_name
  acl    = "private"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }

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
