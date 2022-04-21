# Tamr S3 Module Repo

## v1.2.1 - April 21st 2022
* Replaces deprecated S3 arguments by resource blocks.
* Replaces deprecated IAM policy document argument names.

## v1.2.0 - April 18th 2022
* Resolves S3 bucket public access block tfsec vulnerability.

## v1.1.1 - February 12th 2022
* Updates version file to prevent the major upgrade to the AWS provider version 4.0.

## v1.1.0 - July 12nd 2021
* Adds tags for IAM policies
* Adds new variable `tags` to set tags for all resources
* Deprecates `additional_tags` in favor of `tags`

## v1.0.0 - April 12th 2021
* Updates minimum Terraform version to 13
* Updates minimum AWS provider version to 3.36.0

## v0.2.0 - April 7th 2021
*  Adds new variable `arn_partition` to set the partition of any ARNs referenced in this module

## v0.1.3 - Nov 17th 2020
* Fixes for issues with the outputs.tf when a resource does not exist

## v0.1.2 - Nov 9th 2020
* Adds force_destroy variable

## v0.1.1 - Nov 2nd 2020
* Allows creation of multiple IAM policies on the same bucket
* Adds outputs to examples

## v0.1.0 - Sep 25th 2020
* Initing project
