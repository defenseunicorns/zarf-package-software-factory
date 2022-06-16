package test_test

import (
	"context"
	"testing"
	"time"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/utils"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// TestAllServicesRunning waits until all services report that they are ready.
func TestAllServicesRunning(t *testing.T) {
	// BOILERPLATE, EXPECTED TO BE PRESENT AT THE BEGINNING OF EVERY TEST FUNCTION

	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	go utils.HoldYourDamnHorses(ctx, t, 10*time.Second)
	defer cancel()
	platform := types.NewTestPlatform(t)
	defer platform.Teardown()
	utils.SetupTestPlatform(t, platform)
	// The repo has now been downloaded to /root/app and the software factory package deployment has been initiated.
	teststructure.RunTestStage(platform.T, "TEST", func() {
		// END BOILERPLATE

		// TEST CODE STARTS HERE.

		// Just make sure we can hit the cluster
		output, err := platform.RunSSHCommandAsSudo(`kubectl get nodes`)
		require.NoError(t, err, output)
		// Wait up to 18 minutes for the "bigbang" kustomization to report "Ready==True". Our testing shows if everything goes right this should take 11-13 minutes.
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/bigbang -n flux-system --for=condition=Ready --timeout=1080s`)
		require.NoError(t, err, output)
		// Wait up to 4 minutes for the "softwarefactoryaddons-deps" kustomization to report "Ready==True".
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/softwarefactoryaddons-deps -n flux-system --for=condition=Ready --timeout=240s`)
		require.NoError(t, err, output)
		// Wait up to 4 minutes for the "softwarefactoryaddons" kustomization to report "Ready==True".
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/softwarefactoryaddons -n flux-system --for=condition=Ready --timeout=240s`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the GitLab Webservice Deployment to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 600 bash -c \"while ! kubectl get deployment gitlab-webservice-default -n gitlab; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the GitLab Webservice Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status deployment/gitlab-webservice-default -n gitlab --watch --timeout=600s`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the Jenkins StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 600 bash -c \"while ! kubectl get statefulset jenkins -n jenkins; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the Jenkins StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/jenkins -n jenkins --watch --timeout=600s`)
		require.NoError(t, err, output)
		// Ensure that Jenkins is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 600 bash -c \"while ! kubectl exec statefulset/jenkins -n jenkins -c jenkins -- curl -L -s --fail --show-error https://gitlab.bigbang.dev > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the Jira StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 600 bash -c \"while ! kubectl get statefulset jira -n jira; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the Jira StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/jira -n jira --watch --timeout=600s`)
		require.NoError(t, err, output)
		// Ensure that Jira is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 600 bash -c \"while ! kubectl exec statefulset/jira -n jira -c jira -- curl -L -s --fail --show-error https://gitlab.bigbang.dev > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the Confluence StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 600 bash -c \"while ! kubectl get statefulset confluence -n confluence; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait up to 10 minutes for the Confluence StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/confluence -n confluence --watch --timeout=600s`)
		require.NoError(t, err, output)
		// Ensure that Confluence is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 600 bash -c \"while ! kubectl exec statefulset/confluence -n confluence -c confluence -- curl -L -s --fail --show-error https://gitlab.bigbang.dev > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Make sure flux is present.
		output, err = platform.RunSSHCommandAsSudo("flux --help")
		require.NoError(t, err, output)
	})
}
