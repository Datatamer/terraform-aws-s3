# Tamr AWS S3 IAM Policy Module
This terraform module creates read-only and/or read-write IAM policies that cover a specified S3 bucket.

# Examples
## Basic
```
module "s3-iam-policy-eg" {
  source = "git::https://github.com/Datatamer/terraform-aws-s3.git//modules/bucket-iam-policy?ref=0.1.0"
  bucket_name = "mybucket"
  read_only_paths   = ["mybucket/path/to/ro-folder"]
  read_write_paths  = ["mybucket/path/to/rw-folder", "mybucket/path/to/another-rw-folder"]
}
```
## Minimal
Smallest complete fully working example. This example might require extra resources to run the example.
- [Minimal](https://github.com/Datatamer/terraform-template-repo/tree/master/examples/minimal)

# Resources Created
This modules creates:
* a null resource

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12 |

## Providers

| Name | Version |
|------|---------|
| aws | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| bucket\_name | Name of S3 bucket resource that IAM policies will be created for. | `string` | `n/a` | yes |
| read\_only\_paths | List of bucket paths that should be attached to a read-only policy. | `list(string)` | `[]` | no |
| read\_write\_paths | List of bucket paths that should be attached to a read-write policy. | `list(string)` | `[]` | no |
| read\_only\_actions | List of actions that should be permitted by a read-only policy. | `list(string)` | <pre>[<br>  "s3:Get*",<br>  "s3:List*",<br>  "s3:PutObject",<br>  "s3:PutObjectTagging"<br>]</pre> | no |
| read\_write\_actions | List of actions that should be permitted by a read-write policy. | `list(string)` | <pre>[<br>  "s3:GetBucketLocation",<br>  "s3:GetBucketCORS",<br>  "s3:GetObjectVersionForReplication",<br>  "s3:GetObject",<br>  "s3:GetBucketTagging",<br>  "s3:GetObjectVersion",<br>  "s3:GetObjectTagging",<br>  "s3:ListMultipartUploadParts",<br>  "s3:ListBucketByTags",<br>  "s3:ListBucket",<br>  "s3:ListObjects",<br>  "s3:ListObjectsV2",<br>  "s3:ListBucketMultipartUploads",<br>  "s3:PutObject",<br>  "s3:PutObjectTagging",<br>  "s3:HeadBucket",<br>  "s3:DeleteObject"<br>]</pre> | no |

## Outputs

| Name | Description |
|------|-------------|
| ro_policy_arn | ARN assigned to read-only IAM policy. |
| rw_policy_arn | ARN assigned to read-write IAM policy. |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

# References
This repo is based on:
* [terraform standard module structure](https://www.terraform.io/docs/modules/index.html#standard-module-structure)
* [templated terraform module](https://github.com/tmknom/template-terraform-module)

# Development
## Generating Docs
Run `make terraform/docs` to generate the section of docs around terraform inputs, outputs and requirements.

## Checkstyles
Run `make lint`, this will run terraform fmt, in addition to a few other checks to detect whitespace issues.
NOTE: this requires having docker working on the machine running the test

## Releasing new versions
* Update version contained in `VERSION`
* Document changes in `CHANGELOG.md`
* Create a tag in github for the commit associated with the version

# License
Apache 2 Licensed. See LICENSE for full details.
