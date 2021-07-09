<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| aws | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| existing\_bucket\_name | Name of existing S3 bucket to create policies for. | `string` | n/a | yes |
| additional\_tags | Additional tags for resources created by this example | `map(string)` | <pre>{<br>  "Author": "Tamr",<br>  "Environment": "Example"<br>}</pre> | no |

## Outputs

| Name | Description |
|------|-------------|
| iam-policy-0 | n/a |
| iam-policy-1 | n/a |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
