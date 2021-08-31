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
	testName         string
	vars             map[string]interface{}
	expectApplyError bool
	objTestCases     []ObjectTestCase
	sleepDuration    int
}

// struct used as input for function CondCreateObject
type CondCreateObjectInput struct {
	t         *testing.T
	awsRegion string
	bucket    string
	body      string
	obj       ObjectTestCase
	uploader  *s3manager.Uploader
}

// CondCreateObject conditionally creates an S3 Object inside `bucket` with `body` content
// It is called when we want to test ReadOnly path
// Creates object with default AWS session that presumably has PutObject permission
func CondCreateObject(input CondCreateObjectInput) {
	// When encryption isn't set, trying to upload with default AWS session would fail as well
	// hence the check here
	if input.obj.expectPassWrite || input.obj.encryption == "" {
		return
	}

	logger.Log(input.t, fmt.Sprintf("Default Credential: uploading object %s...", input.obj.key))

	err := UploadObjectWithUploaderE(UploadObjectWithUploaderInput{
		input.awsRegion,
		input.bucket,
		input.obj.key,
		input.obj.encryption,
		input.body,
		input.uploader,
	})

	require.NoError(input.t, err, "Error raised when trying to upload objects that would be tested in ReadOnly paths")

	logger.Log(input.t, fmt.Sprintf("Uploaded object %s with default credential", input.obj.key))

}

// struct used as input for function UploadObjectWithUploader
type UploadObjectWithUploaderInput struct {
	awsRegion  string
	bucketName string
	key        string
	encryption string
	body       string
	uploader   *s3manager.Uploader
}

func UploadObjectWithUploaderE(input UploadObjectWithUploaderInput) error {
	s3Input := &s3manager.UploadInput{
		Bucket: &input.bucketName,
		Key:    &input.key,
		Body:   strings.NewReader(input.body),
	}

	if input.encryption != "" {
		s3Input.ServerSideEncryption = &input.encryption
	}

	_, err := input.uploader.Upload(s3Input)

	return err
}

// struct used as input for function GetS3ObjectWithSession
type GetS3ObjectWithSessionInput struct {
	t *testing.T
	// awsRegion string
	bucket string
	key    string
	sess   *session.Session
}

func GetS3ObjectWithSessionE(input GetS3ObjectWithSessionInput) (string, error) {
	s3Client := s3.New(input.sess)

	res, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: &input.bucket,
		Key:    &input.key,
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
	logger.Log(input.t, fmt.Sprintf("Read contents from s3://%s/%s", input.bucket, input.key))

	return contents, nil
}

// struct used as input for function validatePutObject
type validatePutObjectInput struct {
	t         *testing.T
	awsRegion string
	bucket    string
	obj       ObjectTestCase
	body      string
	sess      *session.Session
}

func validatePutObject(input validatePutObjectInput) {
	up := s3manager.NewUploader(input.sess)

	upObjInput := UploadObjectWithUploaderInput{
		input.awsRegion,
		input.bucket,
		input.obj.key,
		input.obj.encryption,
		input.body,
		up,
	}

	if input.obj.expectPassWrite {
		// upload and require no error
		logger.Log(input.t, fmt.Sprintf("Try to upload using role: %s", input.obj.key))
		err := UploadObjectWithUploaderE(upObjInput)
		require.NoError(input.t, err)
		logger.Log(input.t, fmt.Sprintf("%s uploaded using role", input.obj.key))
	} else {
		// try to upload, require error
		err := UploadObjectWithUploaderE(upObjInput)
		require.Error(input.t, err)
		logger.Log(input.t, fmt.Sprintf("%s could not be uploaded using role", input.obj.key))
	}
}

// struct used as input for function validateGetObject
type validateGetObjectInput struct {
	t            *testing.T
	bucket       string
	obj          ObjectTestCase
	expectedBody string
	sess         *session.Session
}

func validateGetObject(input validateGetObjectInput) {
	getObjInput := GetS3ObjectWithSessionInput{
		input.t,
		input.bucket,
		input.obj.key,
		input.sess,
	}

	if input.obj.expectPassRead {
		// read and require no error
		objContent, err := retry.DoWithRetryE(input.t,
			"Trying to read S3 Object. We may not have permission or just be waiting it to be uploaded and available",
			3, 5*time.Second,
			func() (string, error) {
				return GetS3ObjectWithSessionE(getObjInput)
			})
		require.NoError(input.t, err)
		assert.Equal(input.t, input.expectedBody, objContent)
		logger.Log(input.t, fmt.Sprintf("%s was read using role", input.obj.key))
	} else {
		// try to read, require error
		_, err := retry.DoWithRetryE(input.t,
			"Trying to read S3 Object. We may not have permission or just be waiting it to be uploaded and available",
			3, 5*time.Second,
			func() (string, error) {
				return GetS3ObjectWithSessionE(getObjInput)
			})
		require.Error(input.t, err)
		logger.Log(input.t, fmt.Sprintf("%s could not be read using role", input.obj.key))
	}
}

// struct used as input for function assumeRoleWithRetry
type assumeRoleWithRetryInput struct {
	t         *testing.T
	awsRegion string
	roleARN   string
}

func assumeRoleWithRetry(input assumeRoleWithRetryInput) *session.Session {
	// try to sts:AssumeRole the Role we created earlier
	assumedRoleSession, err := retry.DoWithRetryInterfaceE(input.t,
		"Trying to assume role...",
		3, 5*time.Second,
		func() (interface{}, error) {
			return aws.NewAuthenticatedSessionFromRole(input.awsRegion, input.roleARN)
		})

	// Fail test if sts:AssumeRole fails
	require.NoError(input.t, err)
	logger.Log(input.t, fmt.Sprintf("Assumed role %s successfully", input.roleARN))
	return assumedRoleSession.(*session.Session)
}
