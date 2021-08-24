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

resource "aws_iam_role_policy_attachment" "this" {
  for_each = toset(var.policies_arn)

  role       = aws_iam_role.this.name
  policy_arn = each.value
}
