package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
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
	testBody := "test"

	// list of different buckets that will be created to be tested
	var bucketTestCases = []BucketTestCase{
		{
			"TestBucket1",
			[]string{"path/to/ro-folder"},
			[]string{"path/to/rw-folder"},
			[]ObjectTestCase{
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
			[]ObjectTestCase{
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

			// in here we use Terratest user (default AWS env or TERRATEST_IAM_ROLE env var)
			// to create objects that should be tested in ReadOnly paths of the policies
			test_structure.RunTestStage(t, "create_ro_objects", func() {
				for _, obj := range testCase.objTestCases {
					obj := obj
					CondCreateObject(CondCreateObjectInput{
						t,
						awsRegion,
						expectedBucketName,
						testBody,
						obj,
						defaultUploader,
					})
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

				assumedRoleSession := assumeRoleWithRetry(assumeRoleWithRetryInput{
					t,
					awsRegion,
					assumeRoleARN,
				})

				for _, obj := range testCase.objTestCases {
					obj := obj

					validatePutObject(validatePutObjectInput{
						t,
						awsRegion,
						bucketName,
						obj,
						testBody,
						assumedRoleSession,
					})

					validateGetObject(validateGetObjectInput{
						t,
						bucketName,
						obj,
						testBody,
						assumedRoleSession,
					})
				}
			})
		})
	}
}
