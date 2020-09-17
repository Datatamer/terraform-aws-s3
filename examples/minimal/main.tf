module "example" {
  source           = "git::https://github.com/Datatamer/terraform-aws-s3?ref=0.1.0"
  bucket_name      = "test-bucket"
  read_only_paths  = ["test-bucket/path/to/ro-folder"]
  read_write_paths = ["test-bucket/path/to/rw-folder", "test-bucket/path/to/another-rw-folder"]
}
