package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	terratestutils "github.com/Datatamer/go-terratest-functions/pkg/terratest_utils"
	"github.com/Datatamer/go-terratest-functions/pkg/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// string to be used as key of the bucket name output.
const output_bucket = "test-bucket"

// string to be used as body for files to be created.
const test_body = "test"

// initTestCases returns a list of BucketTestCase to be used for tests.
func initTestCases() []BucketTestCase {
	return []BucketTestCase{
		{
			region:           "",
			tfDir:            "test_examples/minimal",
			bucket_name:      "",
			testName:         "TestBucketSinglePath",
			expectApplyError: false,
			vars: map[string]interface{}{
				"test_bucket_name": "",
				"read_only_paths":  []string{"path/to/ro-folder"},
				"read_write_paths": []string{"path/to/rw-folder"},
			},
			objTestCases: []ObjectTestCase{
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
		},
		{
			region:           "",
			tfDir:            "test_examples/minimal",
			bucket_name:      "",
			testName:         "TestBucketMultiplePaths",
			expectApplyError: false,
			vars: map[string]interface{}{
				"test_bucket_name": "",
				"read_only_paths":  []string{"path1/to/ro-folder", "path2/to/ro-folder"},
				"read_write_paths": []string{"path1/to/rw-folder", "path2/to/rw-folder"},
			},
			objTestCases: []ObjectTestCase{
				{
					key:             "random/path/obj1",
					encryption:      "AES256",
					expectPassRead:  false,
					expectPassWrite: false,
				},
				{
					key:             "path1/to/ro-folder/obj1",
					encryption:      "AES256",
					expectPassRead:  true,
					expectPassWrite: false,
				},
				{
					key:             "path2/to/ro-folder/obj1",
					encryption:      "AES256",
					expectPassRead:  true,
					expectPassWrite: false,
				},
				{
					key:             "path1/to/rw-folder/obj1",
					encryption:      "AES256",
					expectPassRead:  true,
					expectPassWrite: true,
				},
				{
					key:             "path2/to/rw-folder/obj1",
					encryption:      "AES256",
					expectPassRead:  true,
					expectPassWrite: true,
				},
				{
					// not setting encryption here to make sure we cannot upload unencrypted objects
					encryption:      "",
					key:             "path1/to/rw-folder/obj2",
					expectPassRead:  false,
					expectPassWrite: false,
				},
			},
		},
	}
}

// TestTerraformS3Module is the main function that will initialize TestCases and run all tests
func TestTerraformS3Module(t *testing.T) {
	const MODULE_NAME = "terraform-aws-s3"
	// Override random region if needed.
	// os.Setenv("TERRATEST_REGION", "us-east-1")

	// list of different buckets that will be created to be tested
	bucketTestCases := initTestCases()
	// Generate file containing GCS URL to be used on Jenkins.
	// TERRATEST_BACKEND_BUCKET_NAME and TERRATEST_URL_FILE_NAME are both set on Jenkins declaration.
	gcsTestExamplesURL := terratestutils.GenerateUrlFile(t, MODULE_NAME, os.Getenv("TERRATEST_BACKEND_BUCKET_NAME"), os.Getenv("TERRATEST_URL_FILE_NAME"))
	for _, testCase := range bucketTestCases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase

		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			// These will create a tempTestFolder for each bucketTestCase.
			tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", testCase.tfDir)

			// this stage will generate a random `awsRegion` and a `uniqueId` to be used in tests.
			test_structure.RunTestStage(t, "pick_new_randoms", func() {
				usRegions := []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}
				// This function will first check for the Env Var TERRATEST_REGION and return its value if != ""
				awsRegion := aws.GetRandomStableRegion(t, usRegions, nil)

				test_structure.SaveString(t, tempTestFolder, "region", awsRegion)
				test_structure.SaveString(t, tempTestFolder, "unique_id", strings.ToLower(random.UniqueId()))
			})

			defer test_structure.RunTestStage(t, "teardown", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)
				terraformOptions.MaxRetries = 5

				_, err := terraform.DestroyE(t, terraformOptions)
				if err != nil {
					// If there is an error on destroy, it will be logged.
					logger.Log(t, err)
				}
			})

			test_structure.RunTestStage(t, "setup_options", func() {
				awsRegion := test_structure.LoadString(t, tempTestFolder, "region")
				uniqueID := test_structure.LoadString(t, tempTestFolder, "unique_id")
				backendConfig := terratestutils.ParseBackendConfig(t, gcsTestExamplesURL, testCase.testName, testCase.tfDir)

				testCase.vars["test_bucket_name"] = fmt.Sprintf("terratest-s3-%s", uniqueID)
				testCase.vars["name_prefix"] = uniqueID

				terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
					TerraformDir: tempTestFolder,
					Vars:         testCase.vars,
					EnvVars: map[string]string{
						"AWS_REGION": awsRegion,
					},
					BackendConfig: backendConfig,
					MaxRetries:    2,
				})
				test_structure.SaveTerraformOptions(t, tempTestFolder, terraformOptions)
			})

			test_structure.RunTestStage(t, "create_bucket", func() {
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)
				terraformConfig := &types.TerraformData{
					TerraformBackendConfig: terraformOptions.BackendConfig,
					TerraformVars:          terraformOptions.Vars,
					TerraformEnvVars:       terraformOptions.EnvVars,
				}
				if _, err := terratestutils.UploadFilesE(t, terraformConfig); err != nil {
					logger.Log(t, err)
				}
				_, err := terraform.InitAndApplyE(t, terraformOptions)

				if testCase.expectApplyError {
					require.Error(t, err)
					// If it failed as expected, we should skip the rest (validate function).
					t.SkipNow()
				}
			})

			test_structure.RunTestStage(t, "validate_bucket", func() {
				testCase.region = test_structure.LoadString(t, tempTestFolder, "region")
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)

				validateBucket(t, terraformOptions, testCase)
			})

			// in here we use Terratest user (default AWS env or TERRATEST_IAM_ROLE env var)
			// to create objects that should be tested in ReadOnly paths of the policies
			test_structure.RunTestStage(t, "create_ro_objects", func() {

				testCase.region = test_structure.LoadString(t, tempTestFolder, "region")
				validateCreateObjects(t, testCase)
			})

			test_structure.RunTestStage(t, "validate_bucket_and_policies", func() {
				testCase.region = test_structure.LoadString(t, tempTestFolder, "region")
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)

				validateBucketAndPolicies(t, terraformOptions, testCase)
			})
		})
	}
}
