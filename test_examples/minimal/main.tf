module "minimal" {
  # source           = "git::https://github.com/Datatamer/terraform-aws-s3?ref=1.0.0"
  source           = "../../examples/minimal"
  test_bucket_name = var.test_bucket_name
  read_only_paths  = var.read_only_paths
  read_write_paths = var.read_write_paths
  # read_write_actions = var.read_write_actions
  # read_only_actions  = var.read_only_actions
}
