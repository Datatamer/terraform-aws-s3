module "minimal" {
  # source           = "git::https://github.com/Datatamer/terraform-aws-s3?ref=1.2.1"
  source           = "../../examples/minimal"
  test_bucket_name = var.test_bucket_name
  read_only_paths  = var.read_only_paths
  read_write_paths = var.read_write_paths
}

data "aws_caller_identity" "current" {}

# Assume Role Policy for the Role
data "aws_iam_policy_document" "account-assume-role-policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "AWS"
      identifiers = [format("arn:aws:iam::%s:root", data.aws_caller_identity.current.account_id)]
    }
  }
}

# Creates role and attaches the policies that came from variable list.
resource "aws_iam_role" "this" {
  name               = format("%s_%s", var.name_prefix, "TerratestRole")
  assume_role_policy = data.aws_iam_policy_document.account-assume-role-policy.json
}

resource "aws_iam_role_policy_attachment" "ro" {
  role       = aws_iam_role.this.name
  policy_arn = module.minimal.test-bucket.ro_policy_arn
}

resource "aws_iam_role_policy_attachment" "rw" {
  role       = aws_iam_role.this.name
  policy_arn = module.minimal.test-bucket.rw_policy_arn
}
