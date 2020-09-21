# Tamr AWS S3 Module
This terraform module creates a server-side encrypted S3 bucket with a bucket policy enforcing encryption policies.

# Examples
## Basic
```
module "encrypted-s3-eg" {
  source        = "git::https://github.com/Datatamer/terraform-aws-s3.git//modules/encrypted-bucket?ref=0.1.0"
  bucket_name   = "mybucket"
}
```
## Minimal
Smallest complete fully working example. This example might require extra resources to run the example.
- [Minimal](https://github.com/Datatamer/terraform-aws-s3/tree/master/examples/minimal)

# Resources Created
This modules creates:
* a server-side encrypted S3 bucket
* an attached bucket policy enforcing AES256 encryption

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
| additional\_tags | Additional tags to be attached to the S3 bucket. | `map(string)` | `{}` | no |

**TODO**: Add option for KMS S3 encryption.

## Outputs

| Name | Description |
|------|-------------|
| bucket\_name | Name of S3 bucket that was created. |

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
