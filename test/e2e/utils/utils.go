package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func SetupTestPlatform(t *testing.T) *types.TestPlatform {
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

	tempFolder := teststructure.CopyTerraformFolderToTemp(t, "..", "tf/public-ec2-instance")
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
	platform := types.NewTestPlatform(t, tempFolder)

	teststructure.RunTestStage(t, "SETUP", func() {
		teststructure.SaveTerraformOptions(t, tempFolder, terraformOptions)
		teststructure.SaveEc2KeyPair(t, tempFolder, keyPair)
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

	return platform
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
