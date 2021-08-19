package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformS3Bucket(t *testing.T) {
	t.Parallel()

	// For convenience - uncomment these as well as the "os" import
	// when doing local testing if you need to skip any sections.
	//
	// A common usage for this would be skipping teardown first (keeping infrastructure)
	// and then in your next run skip the setup* and create* steps. This way you can keep testing
	// your Go test against your infrastructure quicker. Be mindful of random-ids, as they would be updated
	// on each run, which would make some assertions fail.

	// os.Setenv("SKIP_", "true")
	// os.Setenv("TERRATEST_REGION", "us-east-1")

	// os.Setenv("SKIP_setup_options", "true")
	// os.Setenv("SKIP_create_bucket", "true")
	// os.Setenv("SKIP_validate_bucket", "true")
	// os.Setenv("SKIP_create_ro_objects", "true")
	// os.Setenv("SKIP_setup_role_options", "true")
	// os.Setenv("SKIP_create_role", "true")
	// os.Setenv("SKIP_validate_bucket_and_policies", "true")

	// os.Setenv("SKIP_teardown", "true")
	// os.Setenv("SKIP_teardown_role", "true")

	// string to be used as body for files to be created
	testBodyString := "test"

	// define struct for each object Test Case
	type objectTestCase struct {
		key             string
		encryption      string
		expectPassRead  bool
		expectPassWrite bool
	}

	// list of different buckets that will be created to be tested
	var bucketTestCases = []struct {
		testName      string
		pathRO        []string
		pathRW        []string
		objTestCases  []objectTestCase
		sleepDuration int
	}{
		{
			"TestBucket1",
			[]string{"path/to/ro-folder"},
			[]string{"path/to/rw-folder"},
			[]objectTestCase{
				{
					key:             "path/to/ro-folder/obj1",
					encryption:      "AES256",
					expectPassRead:  true,
					expectPassWrite: false,
				},
				{
					key:             "path/to/rw-folder/obj2",
					encryption:      "AES256",
					expectPassRead:  true,
					expectPassWrite: true,
				},
				{
					key:             "other/folder/obj3",
					encryption:      "AES256",
					expectPassRead:  false,
					expectPassWrite: false,
				},
			},
			0,
		},
		{
			"TestBucket2",
			[]string{"other/path/to/ro-folder"},
			[]string{"other/path/to/rw-folder"},
			[]objectTestCase{
				{
					key: "path/to/ro-folder/obj1",
					// encryption:      "",
					expectPassRead:  false,
					expectPassWrite: false,
				},
				{
					key: "path/to/rw-folder/obj2",
					// encryption:      "",
					expectPassRead:  false,
					expectPassWrite: false,
				},
				{
					key: "other/rw-folder/obj3",
					// encryption:      "",
					expectPassRead:  false,
					expectPassWrite: false,
				},
			},
			0,
		},
	}

	awsRegion := aws.GetRandomStableRegion(t, []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}, nil)

	defaultSession, err := aws.NewAuthenticatedSession(awsRegion)
	if err != nil {
		assert.FailNow(t, "Failed in creating session")
	}
	defaultUploader := s3manager.NewUploader(defaultSession)

	for _, testCase := range bucketTestCases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase

		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			// This is ugly - but attempt to stagger the test cases to
			// avoid a concurrency issue
			// time.Sleep(time.Duration(testCase.sleepDuration) * time.Second)

			// this creates a tempTestFolder for each bucketTestCase
			tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", "test_examples/minimal")
			roleTempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", "test/helpers/iam_lpp")

			defer test_structure.RunTestStage(t, "teardown", func() {
				teraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)
				terraform.Destroy(t, teraformOptions)
			})

			expectedBucketName := fmt.Sprintf("terratest-aws-s3-example-%s", strings.ToLower(random.UniqueId()))

			test_structure.RunTestStage(t, "setup_options", func() {
				terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
					TerraformDir: tempTestFolder,
					Vars: map[string]interface{}{
						"test_bucket_name": expectedBucketName,
						"read_only_paths":  testCase.pathRO,
						"read_write_paths": testCase.pathRW,
					},
					EnvVars: map[string]string{
						"AWS_REGION": awsRegion,
					},
				})

				test_structure.SaveTerraformOptions(t, tempTestFolder, terraformOptions)
			})

			test_structure.RunTestStage(t, "create_bucket", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)
				terraform.InitAndApply(t, terraformOptions)
			})

			test_structure.RunTestStage(t, "validate_bucket", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)
				bucketMap := terraform.OutputMapOfObjects(t, terraformOptions, "test-bucket")
				bucketName := bucketMap["bucket_name"].(string)
				assert.Equal(t, expectedBucketName, bucketName)
				aws.AssertS3BucketExists(t, awsRegion, bucketName)
				aws.AssertS3BucketPolicyExists(t, awsRegion, bucketName)
			})

			// in here we use Terratest (default AWS env or TERRATEST_IAM_ROLE envvar) user to create objects
			// that should be tested in ReadOnly paths
			// (in other words, we expect PutObject to fail for this object)
			test_structure.RunTestStage(t, "create_ro_objects", func() {
				for _, obj := range testCase.objTestCases {
					obj := obj
					// Creates object with default AWS session that presumably has PutObject permission
					// If encryption isn't set, trying to upload with default AWS session would fail as well
					if !obj.expectPassWrite && obj.encryption != "" {
						logger.Log(t, fmt.Sprintf("Default Credential: uploading object %s...", obj.key))
						err := uploadObjectWithUploaderE(awsRegion, expectedBucketName, obj.key, obj.encryption, testBodyString, defaultUploader)
						require.NoError(t, err, "Error raised when trying to upload objects that would be tested in ReadOnly paths")
						logger.Log(t, fmt.Sprintf("Uploaded object %s with default credential", obj.key))
					}
				}
			})

			defer test_structure.RunTestStage(t, "teardown_role", func() {
				teraformOptions := test_structure.LoadTerraformOptions(t, roleTempTestFolder)
				terraform.Destroy(t, teraformOptions)
			})

			test_structure.RunTestStage(t, "setup_role_options", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)
				bucketMap := terraform.OutputMapOfObjects(t, terraformOptions, "test-bucket")
				rwPolicyARN := bucketMap["rw_policy_arn"].(string)
				roPolicyARN := bucketMap["ro_policy_arn"].(string)
				roleTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
					TerraformDir: roleTempTestFolder,
					Vars: map[string]interface{}{
						"name_prefix":  fmt.Sprintf("a%s", strings.ToLower(random.UniqueId())),
						"policies_arn": []string{rwPolicyARN, roPolicyARN},
					},
					EnvVars: map[string]string{
						"AWS_REGION": awsRegion,
					},
				})

				test_structure.SaveTerraformOptions(t, roleTempTestFolder, roleTerraformOptions)
			})

			test_structure.RunTestStage(t, "create_role", func() {
				roleTerraformOptions := test_structure.LoadTerraformOptions(t, roleTempTestFolder)
				terraform.InitAndApply(t, roleTerraformOptions)
			})

			test_structure.RunTestStage(t, "validate_bucket_and_policies", func() {
				// load terraform environments
				roleTerraformOptions := test_structure.LoadTerraformOptions(t, roleTempTestFolder)
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)

				// grab outputs
				bucketMap := terraform.OutputMapOfObjects(t, terraformOptions, "test-bucket")
				bucketName := bucketMap["bucket_name"].(string)
				assumeRoleARN := terraform.Output(t, roleTerraformOptions, "role_arn")

				// try to sts:AssumeRole the Role we created earlier
				assumedRoleSession, err := retry.DoWithRetryInterfaceE(t, "Trying to assume role...", 3, 5*time.Second, func() (interface{}, error) {
					return aws.NewAuthenticatedSessionFromRole(awsRegion, assumeRoleARN)
				})

				// Fail test if sts:AssumeRole fails
				require.NoError(t, err)
				logger.Log(t, fmt.Sprintf("Assumed role %s successfully", assumeRoleARN))

				// using the same Uploader object as it is safe to use them concurrently if needed.
				uploaderRole := s3manager.NewUploader(assumedRoleSession.(*session.Session))

				for _, obj := range testCase.objTestCases {
					obj := obj
					if obj.expectPassWrite {
						// upload and require no error
						logger.Log(t, fmt.Sprintf("Try to upload using role: %s", obj.key))
						err := uploadObjectWithUploaderE(awsRegion, bucketName, obj.key, obj.encryption, testBodyString, uploaderRole)
						require.NoError(t, err)
						logger.Log(t, fmt.Sprintf("%s uploaded using role", obj.key))
					} else {
						// try to upload, require error
						err := uploadObjectWithUploaderE(awsRegion, bucketName, obj.key, obj.encryption, testBodyString, uploaderRole)
						require.Error(t, err)
						logger.Log(t, fmt.Sprintf("%s could not be uploaded using role", obj.key))
					}

					if obj.expectPassRead {
						// read and require no error
						objContent, err := retry.DoWithRetryE(t, "Trying to read S3 Object. We may not have permission or just be waiting it to be uploaded and available",
							3, 5*time.Second,
							func() (string, error) {
								return GetS3ObjectContentsWithSessionE(t, awsRegion, bucketName, obj.key, assumedRoleSession.(*session.Session))
							})
						require.NoError(t, err)
						assert.Equal(t, testBodyString, objContent)
						logger.Log(t, fmt.Sprintf("%s was read using role", obj.key))
					} else {
						// try to read, require error
						_, err := retry.DoWithRetryE(t, "Trying to read S3 Object. We may not have permission or just be waiting it to be uploaded and available",
							3, 5*time.Second,
							func() (string, error) {
								return GetS3ObjectContentsWithSessionE(t, awsRegion, bucketName, obj.key, assumedRoleSession.(*session.Session))
							})
						require.Error(t, err)
						logger.Log(t, fmt.Sprintf("%s could not be read using role", obj.key))
					}
				}
			})
		})
	}
}
