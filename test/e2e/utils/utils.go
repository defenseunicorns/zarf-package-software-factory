package utils

import (
	"fmt"
	"os"
	"testing"
	"time"

	customteststructure "github.com/defenseunicorns/zarf-package-software-factory/test/e2e/terratest/teststructure"
	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// SetupTestPlatform uses Terratest to create an EC2 instance. It then (on the new instance) downloads
// the repo specified by env var REPO_URL at the ref specified by env var GIT_BRANCH, installs Zarf,
// logs into registry1.dso.mil using env vars REGISTRY1_USERNAME and REGISTRY1_PASSWORD, builds all
// the packages, and deploys the init package, the flux package, and the software factory package.
// It is finished when the zarf command returns from deploying the software factory package. It is
// the responsibility of the test being run to do the appropriate waiting for services to come up.
func SetupTestPlatform(t *testing.T, platform *types.TestPlatform) {
	t.Helper()
	repoURL, err := getEnvVar("REPO_URL")
	require.NoError(t, err)
	gitBranch, err := getEnvVar("GIT_BRANCH")
	require.NoError(t, err)
	awsRegion, err := getAwsRegion()
	require.NoError(t, err)
	registry1Username, err := getEnvVar("REGISTRY1_USERNAME")
	require.NoError(t, err)
	registry1Password, err := getEnvVar("REGISTRY1_PASSWORD")
	require.NoError(t, err)
	namespace := "di2me"
	stage := "terratest"
	name := fmt.Sprintf("e2e-%s", random.UniqueId())
	instanceType := "m6i.8xlarge"
	teststructure.RunTestStage(t, "SETUP", func() {
		keyPairName := fmt.Sprintf("%s-%s-%s", namespace, stage, name)
		keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, keyPairName)
		terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
			TerraformDir: platform.TestFolder,
			Vars: map[string]interface{}{
				"aws_region":    awsRegion,
				"namespace":     namespace,
				"stage":         stage,
				"name":          name,
				"key_pair_name": keyPairName,
				"instance_type": instanceType,
			},
		})
		teststructure.SaveTerraformOptions(t, platform.TestFolder, terraformOptions)
		// Use a custom version of this function because the upstream version leaks the private SSH key in the pipeline logs
		customteststructure.SaveEc2KeyPair(t, platform.TestFolder, keyPair)
		terraform.InitAndApply(t, terraformOptions)

		// It can take a minute or so for the instance to boot up, so retry a few times
		err = waitForInstanceReady(t, platform, 5*time.Second, 15) //nolint:gomnd
		require.NoError(t, err)

		// Install dependencies. Doing it here since the instance user-data is being flaky, still saying things like make are not installed
		output, err := platform.RunSSHCommandAsSudo(`apt update && apt install -y jq git make wget sslscan && sysctl -w vm.max_map_count=262144`)
		require.NoError(t, err, output)

		// Clone the repo idempotently
		output, err = platform.RunSSHCommandAsSudo(fmt.Sprintf(`rm -rf ~/app && git clone --depth 1 %v --branch %v --single-branch ~/app`, repoURL, gitBranch))
		require.NoError(t, err, output)

		// Install Zarf
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app && make build/zarf`)
		require.NoError(t, err, output)
		// Log into registry1.dso.mil
		output, err = platform.RunSSHCommandAsSudo(fmt.Sprintf(`~/app/build/zarf tools registry login registry1.dso.mil -u %v -p %v`, registry1Username, registry1Password))
		require.NoError(t, err, output)
		// Build init package
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app && make build/zarf-init-amd64.tar.zst`)
		require.NoError(t, err, output)
		// Build flux package
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app && make build/zarf-package-flux-amd64.tar.zst`)
		require.NoError(t, err, output)
		// Build software factory package
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app && make build/zarf-package-software-factory-amd64.tar.zst`)
		require.NoError(t, err, output)
		// Try to be idempotent
		_, _ = platform.RunSSHCommandAsSudo(`cd ~/app/build && ./zarf destroy --confirm`)
		// Deploy init package
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/build && ./zarf package deploy zarf-init-amd64.tar.zst --components k3s,git-server --confirm`)
		require.NoError(t, err, output)
		// Deploy Flux
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/build && ./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm`)
		require.NoError(t, err, output)
		// Generate a bogus gpg key so it can be applied to flux since flux complains if one isn't present, even if one isn't needed. Only do it if it doesn't already exist.
		output, err = platform.RunSSHCommandAsSudo(`gpg --list-secret-keys user@example.com || gpg --batch --passphrase "" --quick-gen-key user@example.com default default`)
		require.NoError(t, err, output)
		// Apply the bogus gpg key so Flux won't complain
		output, err = platform.RunSSHCommandAsSudo(`gpg --export-secret-keys --armor user@example.com | kubectl create secret generic sops-gpg -n flux-system --from-file=sops.asc=/dev/stdin`)
		require.NoError(t, err, output)
		// Deploy software factory
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/build && ./zarf package deploy zarf-package-software-factory-amd64.tar.zst --components flux-cli --confirm`)
		require.NoError(t, err, output)
		// We have to patch the zarf-package-software-factory GitRepo to point at the right branch
		output, err = platform.RunSSHCommandAsSudo(fmt.Sprintf(`kubectl patch gitrepositories.source.toolkit.fluxcd.io -n flux-system zarf-package-software-factory --type=json -p '"'"'[{"op": "replace", "path": "/spec/ref/branch", "value": "%v"}]'"'"'`, gitBranch))
		require.NoError(t, err, output)
	})
}

// getAwsRegion returns the desired AWS region to use by first checking the env var AWS_REGION, then checking
// AWS_DEFAULT_REGION if AWS_REGION isn't set. If neither is set it returns an error.
func getAwsRegion() (string, error) {
	val, present := os.LookupEnv("AWS_REGION")
	if !present {
		val, present = os.LookupEnv("AWS_DEFAULT_REGION")
	}
	if !present {
		return "", fmt.Errorf("expected either AWS_REGION or AWS_DEFAULT_REGION env var to be set, but they were not")
	}

	return val, nil
}

// getEnvVar gets an environment variable, returning an error if it isn't found.
func getEnvVar(varName string) (string, error) {
	val, present := os.LookupEnv(varName)
	if !present {
		return "", fmt.Errorf("expected env var %v not set", varName)
	}

	return val, nil
}

// waitForInstanceReady tries/retries a simple SSH command until it works successfully, meaning the server is ready to accept connections.
func waitForInstanceReady(t *testing.T, platform *types.TestPlatform, timeBetweenRetries time.Duration, maxRetries int) error {
	t.Helper()
	_, err := retry.DoWithRetryE(t, "Wait for the instance to be ready", maxRetries, timeBetweenRetries, func() (string, error) {
		_, err := platform.RunSSHCommandAsSudo("whoami")
		if err != nil {
			return "", fmt.Errorf("unknown error: %w", err)
		}

		return "", nil
	})
	if err != nil {
		return fmt.Errorf("error while waiting for instance to be ready: %w", err)
	}

	// Wait another 5 seconds because race conditions suck
	time.Sleep(5 * time.Second) //nolint:gomnd

	return nil
}
