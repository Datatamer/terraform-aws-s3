package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformS3CreateUnencryptedBucket(t *testing.T) {
	t.Parallel()

	expectedName := fmt.Sprintf("terratest-aws-s3-example-%s", strings.ToLower(random.UniqueId()))

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test_minimal",

		Vars: map[string]interface{}{
			"test_bucket_name": expectedName,
			"aws_region":       awsRegion,
			// "tag_bucket_environment": expectedEnvironment,
			// "with_policy": "true",
		},

		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	b, err := terraform.OutputMapOfObjectsE(t, terraformOptions, "test-bucket")

	if err != nil {
		logger.Log(t, err)
	}

	assert.Equal(t, expectedName, b["bucket_name"].(string))
	// aws.AssertS3BucketPolicyExists(t, awsRegion, b["bucket_name"].(string))

}

/*
map[bucket_name:terratest-aws-s3-example-a5qoov ro_policy_arn:arn:aws:iam::131578276461:policy/terratest-aws-s3-example-a5qoov-read-only-J4XMR5 rw_policy_arn:arn:aws:iam::131578276461:policy/terratest-aws-s3-example-a5qoov-read-write-J4XMR5]
*/
