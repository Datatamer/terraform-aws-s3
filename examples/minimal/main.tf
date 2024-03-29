module "minimal" {
  # source           = "git::https://github.com/Datatamer/terraform-aws-s3?ref=1.3.2"
  source            = "../../"
  bucket_name       = var.test_bucket_name
  read_only_paths   = var.read_only_paths  # ["path/to/ro-folder"]
  read_write_paths  = var.read_write_paths # ["path/to/rw-folder", "path/to/another-rw-folder"]
  tags              = var.tags
  s3_bucket_logging = var.s3_bucket_logging
}
