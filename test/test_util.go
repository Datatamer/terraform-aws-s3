package test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// define struct for each object Test Case
type ObjectTestCase struct {
	key             string
	encryption      string
	expectPassRead  bool
	expectPassWrite bool
}

// define struct for each bucket Test Case
type BucketTestCase struct {
	testName     string
	vars         map[string]interface{}
	objTestCases []ObjectTestCase
}

// MaybeCreateObject conditionally creates an S3 Object inside `bucket` with `body` content
// It is called when we want to test ReadOnly path
// Creates object with default AWS session that presumably has PutObject permission
func MaybeCreateObject(t *testing.T, awsRegion string, bucket string, body string, obj ObjectTestCase) {
	// When encryption isn't set, trying to upload with default AWS session would fail as well
	// hence the check here
	if obj.expectPassWrite || obj.encryption == "" {
		return
	}

	sess, err := aws.NewAuthenticatedSession(awsRegion)
	if err != nil {
		assert.FailNow(t, "Failed in creating session")
	}
	uploader := s3manager.NewUploader(sess)

	logger.Log(t, fmt.Sprintf("Default Credential: uploading object %s...", obj.key))

	_, err = UploadObjectWithUploaderE(
		awsRegion,
		bucket,
		obj.key,
		obj.encryption,
		body,
		uploader,
	)

	require.NoError(t, err, "Error raised when trying to upload objects that would be tested in ReadOnly paths")

	logger.Log(t, fmt.Sprintf("Uploaded object %s with default credential", obj.key))

}

func UploadObjectWithUploaderE(awsRegion string, bucketName string, key string, encryption string, body string, uploader *s3manager.Uploader) (*s3manager.UploadOutput, error) {
	s3Input := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   strings.NewReader(body),
	}

	if encryption != "" {
		s3Input.ServerSideEncryption = &encryption
	}

	return uploader.Upload(s3Input)

}

func GetS3ObjectWithSessionE(t *testing.T, bucket string, key string, sess *session.Session) (string, error) {
	s3Client := s3.New(sess)

	res, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}

	contents := buf.String()
	logger.Log(t, fmt.Sprintf("Read contents from s3://%s/%s", bucket, key))

	return contents, nil
}

func validatePutObject(t *testing.T, awsRegion string, bucket string, obj ObjectTestCase, body string, sess *session.Session) {
	_, err := retry.DoWithRetryInterfaceE(t,
		"Trying to upload S3 Object..",
		4, 3*time.Second,
		func() (interface{}, error) {
			return UploadObjectWithUploaderE(awsRegion, bucket, obj.key, obj.encryption, body, s3manager.NewUploader(sess))
		})

	if obj.expectPassWrite {
		logger.Log(t, fmt.Sprintf("Try to upload using module permissions: %s", obj.key))
		require.NoError(t, err)
		logger.Log(t, fmt.Sprintf("%s uploaded using module permissions", obj.key))
	} else {
		require.Error(t, err)
		logger.Log(t, fmt.Sprintf("%s could not be uploaded using module permissions", obj.key))
	}
}

func validateGetObject(t *testing.T, bucket string, obj ObjectTestCase, expectedBody string, sess *session.Session) {
	body, err := GetS3ObjectWithSessionE(t, bucket, obj.key, sess)

	if obj.expectPassRead {
		require.NoError(t, err)
		assert.Equal(t, expectedBody, body)
		logger.Log(t, fmt.Sprintf("%s was read using role", obj.key))
	} else {
		require.Error(t, err)
		logger.Log(t, fmt.Sprintf("%s could not be read using role", obj.key))
	}
}

func assumeRoleWithRetry(t *testing.T, awsRegion string, roleARN string) *session.Session {
	assumedRoleSession, err := retry.DoWithRetryInterfaceE(t,
		"Trying to assume role...",
		3, 5*time.Second,
		func() (interface{}, error) {
			return aws.NewAuthenticatedSessionFromRole(awsRegion, roleARN)
		})

	// Fail test if sts:AssumeRole fails
	require.NoError(t, err)
	logger.Log(t, fmt.Sprintf("Assumed role %s successfully", roleARN))
	return assumedRoleSession.(*session.Session)
}

func getPoliciesArnFromOutput(t *testing.T, testFolder string, outputName string) (rwPolicyARN string, roPolicyARN string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, testFolder)
	bucketMap := terraform.OutputMapOfObjects(t, terraformOptions, outputName)
	return bucketMap["rw_policy_arn"].(string), bucketMap["ro_policy_arn"].(string)
}

func getBucketNameFromOutput(t *testing.T, testFolder string, outputName string) string {
	terraformOptions := test_structure.LoadTerraformOptions(t, testFolder)
	bucketMap := terraform.OutputMapOfObjects(t, terraformOptions, outputName)
	return bucketMap["bucket_name"].(string)
}

// func getBucketName(t *testing.T, testFolder string) string {
// 	uniqueID := test_structure.LoadString(t, testFolder, "unique_id")
// 	return fmt.Sprintf("terratest-s3-%s", uniqueID)
// }
