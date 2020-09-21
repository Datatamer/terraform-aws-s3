module "encrypted-bucket" {
  source          = "./modules/encrypted-bucket"
  bucket_name     = var.bucket_name
  additional_tags = var.additional_tags
}

module "bucket-iam-policy" {
  source             = "./modules/bucket-iam-policy"
  bucket_name        = var.bucket_name
  read_only_paths    = var.read_only_paths
  read_write_paths   = var.read_write_paths
  read_only_actions  = var.read_only_actions
  read_write_actions = var.read_write_actions
}
