####################################################
# Example of using complete terraform-aws-s3 module
####################################################
module "minimal" {
  # source           = "git::https://github.com/Datatamer/terraform-aws-s3?ref=0.1.0"
  source           = "../../"
  bucket_name      = "test-bucket"
  read_only_paths  = ["test-bucket/path/to/ro-folder"]
  read_write_paths = ["test-bucket/path/to/rw-folder", "test-bucket/path/to/another-rw-folder"]
}

####################################################
# Example of creating IAM policy for pre-existing S3 bucket
####################################################
data "aws_s3_bucket" "test-bucket" {
  bucket = var.test_bucket_name
}

module "existing-bucket-iam" {
  # source = "git::https://github.com/Datatamer/terraform-aws-s3.git//modules/bucket-iam-policy?ref=0.1.0"
  source      = "../../modules/bucket-iam-policy"
  bucket_name = var.test_bucket_name
  read_only_paths  = ["${var.test_bucket_name}/read-only-folder"]
  read_write_paths = ["${var.test_bucket_name}/read-write-folder"]
}
