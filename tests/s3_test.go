package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformS3Bucket(t *testing.T) {
	t.Parallel()

	expectedName := fmt.Sprintf("terratest-aws-s3-example-%s", strings.ToLower(random.UniqueId()))

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/test_minimal",

		Vars: map[string]interface{}{
			"test_bucket_name": expectedName,
			"aws_region":       awsRegion,
		},

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

	bucket_name := b["bucket_name"].(string)

	assert.Equal(t, expectedName, bucket_name)
	aws.AssertS3BucketPolicyExists(t, awsRegion, bucket_name)

	// Test Unencrypted File
	bodyString := "test-body"
	k := fmt.Sprintf("example-file-%s", strings.ToLower(random.UniqueId()))
	upParams := &s3manager.UploadInput{
		Bucket: &bucket_name,
		Key:    &k,
		Body:   strings.NewReader(bodyString),
	}

	up := aws.NewS3Uploader(t, awsRegion)

	_, err = up.Upload(upParams)
	assert.Error(t, err)

	// Test AES256 Encrypted File
	k = fmt.Sprintf("example-enc-file-%s", strings.ToLower(random.UniqueId()))
	e := "AES256"
	upParams = &s3manager.UploadInput{
		Bucket:               &bucket_name,
		Key:                  &k,
		Body:                 strings.NewReader(bodyString),
		ServerSideEncryption: &e,
	}
	up = aws.NewS3Uploader(t, awsRegion)

	_, err = up.Upload(upParams)

	objContent := retry.DoWithRetry(t, "desc", 5, 3*time.Second, func() (string, error) {
		return aws.GetS3ObjectContentsE(t, awsRegion, bucket_name, k)
	})

	assert.Equal(t, bodyString, objContent)
}
