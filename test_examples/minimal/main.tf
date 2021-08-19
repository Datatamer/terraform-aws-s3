module "minimal" {
  # source           = "git::https://github.com/Datatamer/terraform-aws-s3?ref=1.0.0"
  source           = "../../"
  bucket_name      = var.test_bucket_name
  read_only_paths  = var.read_only_paths
  read_write_paths = var.read_write_paths
}
