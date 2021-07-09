# Example of creating IAM policy for pre-existing S3 bucket
data "aws_s3_bucket" "existing-bucket" {
  bucket = var.existing_bucket_name
}

module "existing-bucket-iam-0" {
  # source = "git::https://github.com/Datatamer/terraform-aws-s3.git//modules/bucket-iam-policy?ref=1.0.0"
  source           = "../../modules/bucket-iam-policy"
  bucket_name      = data.aws_s3_bucket.existing-bucket.id
  read_write_paths = ["some/read-write-folder"]
  additional_tags  = var.additional_tags
}

module "existing-bucket-iam-1" {
  # source = "git::https://github.com/Datatamer/terraform-aws-s3.git//modules/bucket-iam-policy?ref=1.0.0"
  source           = "../../modules/bucket-iam-policy"
  bucket_name      = data.aws_s3_bucket.existing-bucket.id
  read_write_paths = ["another/read-write-folder"]
  additional_tags  = var.additional_tags
}
