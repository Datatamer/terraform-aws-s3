# Tamr AWS S3 IAM Policy Module
This terraform module creates read-only and/or read-write IAM policies that cover a specified S3 bucket.

# Examples
## Basic
```
module "s3-iam-policy-eg" {
  source = "git::https://github.com/Datatamer/terraform-aws-s3.git//modules/bucket-iam-policy?ref=1.0.0"
  bucket_name = "mybucket"
  read_only_paths   = ["path/to/ro-folder"]
  read_write_paths  = ["path/to/rw-folder", "path/to/another-rw-folder"]
}
```

## Minimal
Smallest complete fully working example. This example might require extra resources to run the example.
- [Minimal](https://github.com/Datatamer/terraform-aws-s3/tree/master/examples/iam-policy-submodule)

# Resources Created
This modules creates:
* read-only and/or read-write IAM policies
  * IAM policies created by this module are intended to be attached to _service roles_ downstream. S3-related permissions intended for an instance profile should be configured entirely downstream.
  * If neither `read_only_paths` nor `read_write_paths` are provided, the module will default to creating a read-only IAM policy on the entire bucket
  * If you set `read_write_paths` to `[""]`, the module will permit `read_write_actions` on the whole bucket

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12 |
| aws | >= 2.45.0 |

## Providers

| Name | Version |
|------|---------|
| aws | >= 2.45.0 |
| random | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| bucket\_name | Name of S3 bucket resource that IAM policies will be created for. | `string` | n/a | yes |
| arn\_partition | The partition in which the resource is located. A partition is a group of AWS Regions.<br>  Each AWS account is scoped to one partition.<br>  The following are the supported partitions:<br>    aws -AWS Regions<br>    aws-cn - China Regions<br>    aws-us-gov - AWS GovCloud (US) Regions | `string` | `"aws"` | no |
| read\_only\_actions | List of actions that should be permitted by a read-only policy. | `list(string)` | <pre>[<br>  "s3:Get*",<br>  "s3:List*"<br>]</pre> | no |
| read\_only\_paths | List of paths/prefixes that should be attached to a read-only policy. Listed path(s) should omit the head bucket. | `list(string)` | `[]` | no |
| read\_write\_actions | List of actions that should be permitted by a read-write policy. | `list(string)` | <pre>[<br>  "s3:GetBucketLocation",<br>  "s3:GetBucketCORS",<br>  "s3:GetObjectVersionForReplication",<br>  "s3:GetObject",<br>  "s3:GetBucketTagging",<br>  "s3:GetObjectVersion",<br>  "s3:GetObjectTagging",<br>  "s3:ListMultipartUploadParts",<br>  "s3:ListBucketByTags",<br>  "s3:ListBucket",<br>  "s3:ListObjects",<br>  "s3:ListObjectsV2",<br>  "s3:ListBucketMultipartUploads",<br>  "s3:PutObject",<br>  "s3:PutObjectTagging",<br>  "s3:HeadBucket",<br>  "s3:DeleteObject"<br>]</pre> | no |
| read\_write\_paths | List of paths/prefixes that should be attached to a read-write` policy. Listed path(s) should omit the head bucket.` | `list(string)` | `[]` | no |
| tags | A map of tags to add to all resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| ro\_policy\_arn | ARN assigned to read-only IAM policy. |
| rw\_policy\_arn | ARN assigned to read-write IAM policy. |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

# References
This repo is based on:
* [terraform standard module structure](https://www.terraform.io/docs/modules/index.html#standard-module-structure)
* [templated terraform module](https://github.com/tmknom/template-terraform-module)

# License
Apache 2 Licensed. See LICENSE for full details.
