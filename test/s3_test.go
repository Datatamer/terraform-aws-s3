package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

const output_bucket = "test-bucket"

// initTestCases returns a list of BucketTestCase to be used for tests.
func initTestCases() []BucketTestCase {
	return []BucketTestCase{
		{
			testName: "TestBucket1",
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
			testName: "TestBucket2",
			vars: map[string]interface{}{
				"test_bucket_name": "",
				"read_only_paths":  []string{"other/path/to/ro-folder"},
				"read_write_paths": []string{"other/path/to/rw-folder"},
			},
			objTestCases: []ObjectTestCase{
				{
					key:             "other/path/to/ro-folder/obj1",
					encryption:      "AES256",
					expectPassRead:  true,
					expectPassWrite: false,
				},
				{
					key: "path/to/rw-folder/obj2",
					// not setting encryption here to make sure we cannot upload unencrypted objects
					encryption:      "",
					expectPassRead:  false,
					expectPassWrite: false,
				},
			},
		},
	}
}

func TestTerraformS3Module(t *testing.T) {

	// os.Setenv("SKIP_", "true")
	// os.Setenv("TERRATEST_REGION", "us-east-1")

	// For convenience - uncomment these as well as the "os" import
	// when doing local testing if you need to skip any sections.
	//
	// The good usage for skipping stages is local testing, when you want to make sure your test code is fine,
	// and not waste time building new infrastructure every time.
	// A nice approach for doing that would be:
	// 1- Make sure there are no local files in the dirs (`rm -rf .test-dir terraform.tfstate*`)
	// 2- Skip teardown* and run tests. (create and keep infrastructure, tfstate will be into local folders)
	// 3- In your next run skip `pick*` steps (which makes sure random values in Terratest state won't be updated)
	//    obs: if you forget and run tests once, you can make use of `terraform.tfstate.backup` file to restore it
	// 4- Run local Tests; Do your thing;
	// 5- When you are done and want to destroy the infrastructure -> comment back teardown* steps
	// 6- Comment back pick* steps.

	// os.Setenv("SKIP_pick_new_randoms", "true")
	// os.Setenv("SKIP_setup_options", "true")
	// os.Setenv("SKIP_create_bucket", "true")
	// os.Setenv("SKIP_validate_bucket", "true")
	// os.Setenv("SKIP_create_ro_objects", "true")
	// os.Setenv("SKIP_setup_role_options", "true")
	// os.Setenv("SKIP_create_role", "true")
	// os.Setenv("SKIP_validate_bucket_and_policies", "true")

	// Warning: if you skip these steps, Terratest state will be stored in local folder under .test-data inside each TF dir you skip.
	// These are purposedly stored in the known TF dir (not tmp folder) to make sure you can re-run tests on them at any time.
	// Remember removing those folders after you finish your tests so that it won't affect the next time you run local tests.
	// os.Setenv("SKIP_teardown", "true")
	// os.Setenv("SKIP_teardown_role", "true")
	//
	// string to be used as body for files to be created
	const test_body = "test"

	// list of different buckets that will be created to be tested
	bucketTestCases := initTestCases()

	// When using test_structure functions + parallel tests with random IDs we will run into consistency problems.
	// This is an easy/lazy way to deal with it.
	if test_structure.SkipStageEnvVarSet() && len(bucketTestCases) > 1 {
		logger.Log(t,
			"Won't run tests using SKIP_* vars having multiple cases. Temporary folders are disabled when using SKIP_* (local testing). Not having different folders for each testCase will generate conflicts with state files.")
		t.FailNow()
		// Another solution would be to simply truncate the list of cases instead of failing.
		// bucketTestCases = bucketTestCases[:1]
	}

	for _, testCase := range bucketTestCases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase

		t.Run(testCase.testName, func(t *testing.T) {
			t.Parallel()

			// These will create a tempTestFolder for each bucketTestCase.
			// Also, if any of SKIP_* env vars are set, it won't create temp folders in order to store configuration that must be retained
			// through different local tests when skipping stages.
			// (e.g. `awsRegion` and `uniqueId`.
			tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", "test_examples/minimal")
			roleTempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", "test_examples/helpers/iam_lpp")

			// this stage will generate a random `awsRegion` and a `uniqueId` to be used in tests.
			test_structure.RunTestStage(t, "pick_new_randoms", func() {
				// Pick a random AWS region to test in. This helps ensure your code works in all regions.
				usRegions := []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2"}
				// This function will first check for the Env Var TERRATEST_REGION and return its value if != ""
				awsRegion := aws.GetRandomStableRegion(t, usRegions, nil)

				test_structure.SaveString(t, tempTestFolder, "region", awsRegion)
				test_structure.SaveString(t, tempTestFolder, "unique_id", strings.ToLower(random.UniqueId()))
			})

			defer test_structure.RunTestStage(t, "teardown", func() {
				teraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)
				terraform.Destroy(t, teraformOptions)
			})

			test_structure.RunTestStage(t, "setup_options", func() {
				awsRegion := test_structure.LoadString(t, tempTestFolder, "region")
				uniqueID := test_structure.LoadString(t, tempTestFolder, "unique_id")

				testCase.vars["test_bucket_name"] = fmt.Sprintf("terratest-s3-%s", uniqueID)

				terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
					TerraformDir: tempTestFolder,
					Vars:         testCase.vars,
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
				awsRegion := test_structure.LoadString(t, tempTestFolder, "region")
				bucketName := getBucketNameFromOutput(t, tempTestFolder, output_bucket)
				require.Equal(t, testCase.vars["test_bucket_name"], bucketName)
				aws.AssertS3BucketExists(t, awsRegion, bucketName)
				aws.AssertS3BucketPolicyExists(t, awsRegion, bucketName)
			})

			// in here we use Terratest user (default AWS env or TERRATEST_IAM_ROLE env var)
			// to create objects that should be tested in ReadOnly paths of the policies
			test_structure.RunTestStage(t, "create_ro_objects", func() {
				awsRegion := test_structure.LoadString(t, tempTestFolder, "region")

				for _, obj := range testCase.objTestCases {
					obj := obj
					MaybeCreateObject(
						t,
						awsRegion,
						testCase.vars["test_bucket_name"].(string),
						test_body,
						obj,
					)
				}
			})

			defer test_structure.RunTestStage(t, "teardown_role", func() {
				teraformOptions := test_structure.LoadTerraformOptions(t, roleTempTestFolder)
				terraform.Destroy(t, teraformOptions)
			})

			test_structure.RunTestStage(t, "setup_role_options", func() {
				rwPolicyARN, roPolicyARN := getPoliciesArnFromOutput(t, tempTestFolder, output_bucket)
				uniqueID := test_structure.LoadString(t, tempTestFolder, "unique_id")

				roleTerraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
					TerraformDir: roleTempTestFolder,
					Vars: map[string]interface{}{
						"name_prefix":  uniqueID,
						"policies_arn": []string{rwPolicyARN, roPolicyARN},
					},
				})
				test_structure.SaveTerraformOptions(t, roleTempTestFolder, roleTerraformOptions)
			})

			test_structure.RunTestStage(t, "create_role", func() {
				roleTerraformOptions := test_structure.LoadTerraformOptions(t, roleTempTestFolder)
				terraform.InitAndApply(t, roleTerraformOptions)
			})

			test_structure.RunTestStage(t, "validate_bucket_and_policies", func() {
				awsRegion := test_structure.LoadString(t, tempTestFolder, "region")
				// load terraform environments
				roleTerraformOptions := test_structure.LoadTerraformOptions(t, roleTempTestFolder)
				terraformOptions := test_structure.LoadTerraformOptions(t, tempTestFolder)

				// grab outputs
				bucketMap := terraform.OutputMapOfObjects(t, terraformOptions, output_bucket)
				bucketName := bucketMap["bucket_name"].(string)
				assumeRoleARN := terraform.Output(t, roleTerraformOptions, "role_arn")

				assumedRoleSession := assumeRoleWithRetry(t, awsRegion, assumeRoleARN)

				for _, obj := range testCase.objTestCases {
					obj := obj

					t.Run("put_object", func(t *testing.T) {
						// we don't run this in parallel just to win some time for s3 GetObject API to be available right after.
						validatePutObject(t, awsRegion, bucketName, obj, test_body, assumedRoleSession)
					})

					t.Run("get_object", func(t *testing.T) {
						validateGetObject(t, bucketName, obj, test_body, assumedRoleSession)
					})
				}
			})
		})
	}
}
