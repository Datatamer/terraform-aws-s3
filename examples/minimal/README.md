<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

No requirements.

## Providers

No provider.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| test\_bucket\_name | Name of test S3 bucket name. | `string` | n/a | yes |
| read\_only\_paths | n/a | `list(string)` | <pre>[<br>  "path/to/ro-folder"<br>]</pre> | no |
| read\_write\_paths | n/a | `list(string)` | <pre>[<br>  "path/to/rw-folder",<br>  "path/to/another-rw-folder"<br>]</pre> | no |
| s3\_bucket\_logging | The name of S3 bucket where to store server access logs. | `string` | `""` | no |
| tags | A map of tags to add to all resources created by this example. | `map(string)` | <pre>{<br>  "Author": "Tamr",<br>  "Environment": "Example"<br>}</pre> | no |

## Outputs

| Name | Description |
|------|-------------|
| test-bucket | n/a |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
