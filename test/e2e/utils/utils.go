package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	terratesting "github.com/gruntwork-io/terratest/modules/testing"
)

func InitTestPlatform(t *testing.T) *types.TestPlatform {
	tempFolder := teststructure.CopyTerraformFolderToTemp(t, "..", "tf/public-ec2-instance")
	platform := types.NewTestPlatform(t, tempFolder)
	return platform
}

func SetupTestPlatform(t *testing.T, platform *types.TestPlatform) {
	repoUrl, err := getRepoUrl()
	require.NoError(t, err)
	gitBranch, err := getGitBranch()
	require.NoError(t, err)
	awsRegion, err := getAwsRegion()
	require.NoError(t, err)
	namespace := "di2me"
	stage := "terratest"
	name := fmt.Sprintf("e2e-%s", random.UniqueId())
	instanceType := "m6i.8xlarge"

	// Since Terraform is going to be run with that temp folder as the CWD, we also need our .tool-versions file to be
	// in that directory so that the right version of Terraform is being run there. I can neither confirm nor deny that
	// this took me 2 days to figure out...
	// Since we can't be sure what the working directory is, we are going to walk up one directory at a time until we
	// find a .tool-versions file and then copy it into the temp folder
	found := false
	filePath := ".tool-versions"
	for !found {
		//nolint:gocritic
		if _, err := os.Stat(filePath); err == nil {
			// The file exists
			found = true
		} else if errors.Is(err, os.ErrNotExist) {
			// The file does *not* exist. Add a "../" and try again
			filePath = fmt.Sprintf("../%v", filePath)
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
			require.NoError(t, err)
		}
	}
	tempFolder := platform.TestFolder
	err = copyFile(filePath, fmt.Sprintf("%v/.tool-versions", tempFolder))
	require.NoError(t, err)

	keyPairName := fmt.Sprintf("%s-%s-%s", namespace, stage, name)
	keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, keyPairName)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: tempFolder,
		Vars: map[string]interface{}{
			"aws_region":    awsRegion,
			"namespace":     namespace,
			"stage":         stage,
			"name":          name,
			"key_pair_name": keyPairName,
			"instance_type": instanceType,
		},
	})

	teststructure.RunTestStage(t, "SETUP", func() {
		teststructure.SaveTerraformOptions(t, tempFolder, terraformOptions)
		SaveEc2KeyPair(t, tempFolder, keyPair)
		terraform.InitAndApply(t, terraformOptions)
	})

	// It can take a minute or so for the instance to boot up, so retry a few times
	maxRetries := 15
	timeBetweenRetries, err := time.ParseDuration("5s")
	require.NoError(t, err)
	_, err = retry.DoWithRetryE(t, "Wait for the instance to be ready", maxRetries, timeBetweenRetries, func() (string, error) {
		_, err := platform.RunSSHCommand("whoami")
		if err != nil {
			return "", err
		}
		return "", nil
	})
	require.NoError(t, err)

	// Clone the repo
	output, err := platform.RunSSHCommand(fmt.Sprintf("git clone --depth 1 %v --branch %v --single-branch ~/app", repoUrl, gitBranch))
	require.NoError(t, err, output)
}

func getRepoUrl() (string, error) {
	val, present := os.LookupEnv("REPO_URL")
	if !present {
		return "", fmt.Errorf("expected env var REPO_URL not set")
	} else {
		return val, nil
	}
}

func getGitBranch() (string, error) {
	val, present := os.LookupEnv("GIT_BRANCH")
	if !present {
		return "", fmt.Errorf("expected env var GIT_BRANCH not set")
	} else {
		return val, nil
	}
}

// getAwsRegion returns the desired AWS region to use by first checking the env var AWS_REGION, then checking
// AWS_DEFAULT_REGION if AWS_REGION isn't set. If neither is set it returns an error
func getAwsRegion() (string, error) {
	val, present := os.LookupEnv("AWS_REGION")
	if !present {
		val, present = os.LookupEnv("AWS_DEFAULT_REGION")
	}
	if !present {
		return "", fmt.Errorf("expected either AWS_REGION or AWS_DEFAULT_REGION env var to be set, but they were not")
	} else {
		return val, nil
	}
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherwise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func copyFile(src string, dst string) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return nil
		}
	}
	if err = os.Link(src, dst); err == nil {
		return err
	}
	err = copyFileContents(src, dst)
	return nil
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(in)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return nil
}

// SaveEc2KeyPair serializes and saves an Ec2KeyPair into the given folder. This allows you to create an Ec2KeyPair during setup
// and to reuse that Ec2KeyPair later during validation and teardown.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func SaveEc2KeyPair(t terratesting.TestingT, testFolder string, keyPair *aws.Ec2Keypair) {
	saveTestData(t, formatEc2KeyPairPath(testFolder), keyPair)
}

// SaveTestData serializes and saves a value used at test time to the given path. This allows you to create some sort of test data
// (e.g., TerraformOptions) during setup and to reuse this data later during validation and teardown.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func saveTestData(t terratesting.TestingT, path string, value interface{}) {
	logger.Logf(t, "Storing test data in %s so it can be reused later", path)

	if IsTestDataPresent(t, path) {
		logger.Logf(t, "[WARNING] The named test data at path %s is non-empty. Save operation will overwrite existing value with \"%v\".\n.", path, value)
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to convert value %s to JSON: %v", path, err)
	}

	// Don't log this data, it exposes the EC2 Key Pair's private key in the logs, which are public on GitHub Actions
	// logger.Logf(t, "Marshalled JSON: %s", string(bytes))

	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0777); err != nil {
		t.Fatalf("Failed to create folder %s: %v", parentDir, err)
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		t.Fatalf("Failed to save value %s: %v", path, err)
	}
}

// formatEc2KeyPairPath formats a path to save an Ec2KeyPair in the given folder.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func formatEc2KeyPairPath(testFolder string) string {
	return formatTestDataPath(testFolder, "Ec2KeyPair.json")
}

// FormatTestDataPath formats a path to save test data.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func formatTestDataPath(testFolder string, filename string) string {
	return filepath.Join(testFolder, ".test-data", filename)
}

// IsTestDataPresent returns true if a file exists at $path and the test data there is non-empty.
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func IsTestDataPresent(t terratesting.TestingT, path string) bool {
	exists, err := files.FileExistsE(path)
	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}
	if !exists {
		return false
	}

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}

	if isEmptyJSON(t, bytes) {
		return false
	}

	return true
}

// isEmptyJSON returns true if the given bytes are empty, or in a valid JSON format that can reasonably be considered empty.
// The types used are based on the type possibilities listed at https://golang.org/src/encoding/json/decode.go?s=4062:4110#L51
// This function is directly copied from https://github.com/gruntwork-io/terratest/tree/5913a2925623d3998841cb25de7b26731af9ab13
// due to this issue: https://github.com/gruntwork-io/terratest/issues/1135
func isEmptyJSON(t terratesting.TestingT, bytes []byte) bool {
	var value interface{}

	if len(bytes) == 0 {
		return true
	}

	if err := json.Unmarshal(bytes, &value); err != nil {
		t.Fatalf("Failed to parse JSON while testing whether it is empty: %v", err)
	}

	if value == nil {
		return true
	}

	valueBool, ok := value.(bool)
	if ok && !valueBool {
		return true
	}

	valueFloat64, ok := value.(float64)
	if ok && valueFloat64 == 0 {
		return true
	}

	valueString, ok := value.(string)
	if ok && valueString == "" {
		return true
	}

	valueSlice, ok := value.([]interface{})
	if ok && len(valueSlice) == 0 {
		return true
	}

	valueMap, ok := value.(map[string]interface{})
	if ok && len(valueMap) == 0 {
		return true
	}

	return false
}
