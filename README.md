# Tamr AWS S3 Terraform Module
This terraform module creates an encrypted S3 bucket and any associated IAM policies.

# Examples
## Basic
```
module "tamr-s3-eg" {
  source            = "git::https://github.com/Datatamer/terraform-aws-s3?ref=0.1.0"
  bucket_name       = "mybucket"
  read_only_paths   = ["mybucket/path/to/ro-folder"]
  read_write_paths  = ["mybucket/path/to/rw-folder", "mybucket/path/to/another-rw-folder"]
}
```
**Note about`read_only_paths` and `read_write_paths`:**
* Providing a path to a folder like in the example, `"mybucket/path/to/folder"` permits the actions specified in `read_only_actions`/`read_write_actions` on **`["mybucket/pack/to/folder", "mybucket/path/to/folder/*"]`**.

## Minimal
Smallest complete fully working example. This example might require extra resources to run the example.
- [Minimal](https://github.com/Datatamer/terraform-aws-s3/tree/master/examples/minimal)

# Resources Created
This modules creates:
* a s3 bucket
* a s3 bucket policy to enforce AES256 server-side-encryption
* read-only and/or read-write IAM policies
  * IAM policies created by this module are intended to be attached to _service roles_ downstream. S3-related permissions intended for an instance profile should be configured entirely downstream.
  * **NOTE**: If neither `read_only_paths` nor `read_write_paths` are provided, the module will default to creating a read-only IAM policy on the entire bucket

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
| bucket\_name | Name of S3 bucket to create. | `string` | `n/a` | yes |
| read\_only\_paths | List of bucket paths that should be attached to a read-only policy. | `list(string)` | `[]` | no |
| read\_write\_paths | List of bucket paths that should be attached to a read-write policy. | `list(string)` | `[]` | no |
| read\_only\_actions | List of actions that should be permitted by a read-only policy. | `list(string)` | <pre>[<br>  "s3:Get*",<br>  "s3:List*",<br>]</pre> | no |
| read\_write\_actions | List of actions that should be permitted by a read-write policy. | `list(string)` | <pre>[<br>  "s3:GetBucketLocation",<br>  "s3:GetBucketCORS",<br>  "s3:GetObjectVersionForReplication",<br>  "s3:GetObject",<br>  "s3:GetBucketTagging",<br>  "s3:GetObjectVersion",<br>  "s3:GetObjectTagging",<br>  "s3:ListMultipartUploadParts",<br>  "s3:ListBucketByTags",<br>  "s3:ListBucket",<br>  "s3:ListObjects",<br>  "s3:ListObjectsV2",<br>  "s3:ListBucketMultipartUploads",<br>  "s3:PutObject",<br>  "s3:PutObjectTagging",<br>  "s3:HeadBucket",<br>  "s3:DeleteObject"<br>]</pre> | no |
| additional\_tags | Additional tags to be attached to the S3 bucket. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| bucket\_name | Name of S3 bucket created by encrypted-bucket module. |
| ro\_policy\_arn | ARN assigned to read-only IAM policy created by iam-policy module. |
| rw\_policy\_arn | ARN assigned to read-write IAM policy created by iam-policy module. |
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
