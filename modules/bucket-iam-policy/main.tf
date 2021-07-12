locals {
  # If neither `read_only_paths` nor `read_write_paths` are provided, default to
  # read-only access to the entire bucket
  ro_paths     = length(var.read_only_paths) + length(var.read_write_paths) == 0 ? [""] : var.read_only_paths
  ro_paths_map = { for idx, val in local.ro_paths : idx => val }
  rw_paths_map = { for idx, val in var.read_write_paths : idx => val }
}

# Policy document for read-only access to entire bucket (bucket, bucket/*)
data "aws_iam_policy_document" "ro_source_policy_doc" {
  count = length(local.ro_paths) == 0 ? 0 : 1

  version = "2012-10-17"

  statement {
    sid     = "ReadOnlyPolicy0"
    effect  = "Allow"
    actions = var.read_only_actions
    resources = [
      "arn:${var.arn_partition}:s3:::${var.bucket_name}",
      "arn:${var.arn_partition}:s3:::${var.bucket_name}/*"
    ]
  }
}

# If any read_only_paths are specified, the read-only source policy doc will be
# overwritten by a scoped down bucket resource (bucket/some/path,
# bucket/some/path/*)
data "aws_iam_policy_document" "path_specific_ro_doc" {
  count = length(local.ro_paths) == 0 ? 0 : 1

  version     = "2012-10-17"
  source_json = data.aws_iam_policy_document.ro_source_policy_doc[0].json

  dynamic "statement" {
    for_each = local.ro_paths_map

    content {
      sid     = "ReadOnlyPolicy${statement.key}"
      effect  = "Allow"
      actions = var.read_only_actions
      resources = [
        "arn:${var.arn_partition}:s3:::${var.bucket_name}/${statement.value}",
        "arn:${var.arn_partition}:s3:::${var.bucket_name}/${statement.value}/*"
      ]
    }
  }
}

# Appended to policy name to allow creation of multiple policies on the same bucket.
resource "random_string" "rand" {
  length  = 6
  special = false
}

# Read-only IAM policy
resource "aws_iam_policy" "ro_policy" {
  count = length(local.ro_paths) == 0 ? 0 : 1

  name = format("%s-read-only-%s", var.bucket_name, random_string.rand.result)
  # If you want read-only access to the entire bucket, path_specific_ro_doc should not overwrite ReadOnlyPolicy0 in ro_source_policy_doc
  policy = local.ro_paths[0] == "" ? data.aws_iam_policy_document.ro_source_policy_doc[0].json : data.aws_iam_policy_document.path_specific_ro_doc[0].json
  tags = var.tags
}

# Policy document for read-write access to entire bucket (bucket, bucket/*)
data "aws_iam_policy_document" "rw_source_policy_doc" {
  count = length(var.read_write_paths) == 0 ? 0 : 1

  version = "2012-10-17"

  statement {
    sid     = "ReadWritePolicy0"
    effect  = "Allow"
    actions = var.read_write_actions
    resources = [
      "arn:${var.arn_partition}:s3:::${var.bucket_name}",
      "arn:${var.arn_partition}:s3:::${var.bucket_name}/*"
    ]
  }
}

# If any read_write_paths are specified, the read-write source policy doc will be
# overwritten by a scoped down bucket resource (bucket/some/path,
# bucket/some/path/*)
data "aws_iam_policy_document" "path_specific_rw_doc" {
  count = length(var.read_write_paths) == 0 ? 0 : 1

  version     = "2012-10-17"
  source_json = data.aws_iam_policy_document.rw_source_policy_doc[0].json

  dynamic "statement" {
    for_each = local.rw_paths_map

    content {
      sid     = "ReadWritePolicy${statement.key}"
      effect  = "Allow"
      actions = var.read_write_actions
      resources = [
        "arn:${var.arn_partition}:s3:::${var.bucket_name}/${statement.value}",
        "arn:${var.arn_partition}:s3:::${var.bucket_name}/${statement.value}/*"
      ]
    }
  }
}

# Read-write IAM policy
resource "aws_iam_policy" "rw_policy" {
  count = length(var.read_write_paths) == 0 ? 0 : 1

  name = format("%s-read-write-%s", var.bucket_name, random_string.rand.result)
  # If you want read-write access to the entire bucket, path_specific_rw_doc should not overwrite ReadWritePolicy0 in rw_source_policy_doc
  policy = var.read_write_paths[0] == "" ? data.aws_iam_policy_document.rw_source_policy_doc[0].json : data.aws_iam_policy_document.path_specific_rw_doc[0].json
  tags = var.tags
}
