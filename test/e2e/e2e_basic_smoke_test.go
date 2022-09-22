package test_test

import (
	"testing"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/utils"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// TestAllServicesRunning waits until all services report that they are ready.
func TestAllServicesRunning(t *testing.T) {
	// BOILERPLATE, EXPECTED TO BE PRESENT AT THE BEGINNING OF EVERY TEST FUNCTION

	t.Parallel()
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
		// Wait for the "postgres-operator" kustomization to report "Ready==True".
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/postgres-operator -n flux-system --for=condition=Ready --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the "bigbang" kustomization to report "Ready==True". Our testing shows if everything goes right this should take 11-13 minutes.
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/bigbang -n flux-system --for=condition=Ready --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the "softwarefactoryaddons" kustomization to report "Ready==True".
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/softwarefactoryaddons -n flux-system --for=condition=Ready --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the GitLab Webservice Deployment to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl get deployment gitlab-webservice-default -n gitlab; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the GitLab Webservice Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status deployment/gitlab-webservice-default -n gitlab --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the Jenkins StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl get statefulset jenkins -n jenkins; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the Jenkins StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/jenkins -n jenkins --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Jenkins is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl exec statefulset/jenkins -n jenkins -c jenkins -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the Jira StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl get statefulset jira -n jira; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the Jira StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/jira -n jira --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Jira is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl exec statefulset/jira -n jira -c jira -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the Confluence StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl get statefulset confluence -n confluence; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the Confluence StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/confluence -n confluence --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Confluence is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl exec statefulset/confluence -n confluence -c confluence -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Make sure flux is present.
		output, err = platform.RunSSHCommandAsSudo("flux --help")
		require.NoError(t, err, output)
		// Ensure that Jenkins is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! curl -L -s --fail --show-error https://jenkins.bigbang.dev/login > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Ensure that Confluence is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! curl -L -s --fail --show-error https://confluence.bigbang.dev/status > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Ensure that Jira is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! curl -L -s --fail --show-error https://jira.bigbang.dev/status > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Ensure that GitLab is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the Artifactory StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl get statefulset artifactory -n artifactory; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Wait for the Artifactory StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/artifactory -n artifactory --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Artifactory is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! kubectl exec statefulset/artifactory -n artifactory -c artifactory -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)
		// Ensure that Artifactory is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c \"while ! curl -L -s --fail --show-error https://artifactory.bigbang.dev/artifactory/api/system/ping > /dev/null; do sleep 5; done\"`)
		require.NoError(t, err, output)

		// Ensure that the services do not accept discontinued TLS versions. If they reject TLSv1.1 it is assumed that they also reject anything below TLSv1.1.
		// Ensure that GitLab does not accept TLSv1.1
		output, err = platform.RunSSHCommand(`sslscan gitlab.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Jenkins does not accept TLSv1.1
		output, err = platform.RunSSHCommand(`sslscan jenkins.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Jira does not accept TLSv1.1
		output, err = platform.RunSSHCommand(`sslscan jira.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Confluence does not accept TLSv1.1
		output, err = platform.RunSSHCommand(`sslscan confluence.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Artifactory does not accept TLSv1.1
		output, err = platform.RunSSHCommand(`sslscan artifactory.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
	})
}
