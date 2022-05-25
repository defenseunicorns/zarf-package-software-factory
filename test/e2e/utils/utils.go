package utils

import (
	"fmt"
	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

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

//# grab the repo, install everything, and init zarf
//git clone --depth 1 ${var.repo_url} --branch ${var.git_branch} --single-branch /app
//(cd /app && make build/zarf)
///app/build/zarf tools registry login registry1.dso.mil -u ${var.registry1_username} -p ${var.registry1_password}
//(cd /app && make build/zarf-init-amd64.tar.zst build/zarf-package-flux-amd64.tar.zst build/zarf-package-software-factory-amd64.tar.zst)
///app/build/zarf init --

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
