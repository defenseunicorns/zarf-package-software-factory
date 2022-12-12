package test_test

import (
	"testing"
	"time"

	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/types"
	"github.com/defenseunicorns/zarf-package-software-factory/test/e2e/utils"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

// TestAllServicesRunning waits until all services report that they are ready.
func TestAllServicesRunning(t *testing.T) { //nolint:funlen
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

		// Wait for the databases to report "PostgresClusterStatus==Running", then create a timestamp for 15 minutes and
		// 10 seconds in the future which is when it should be checked again to make sure that nothing is blocking the
		// operator from talking to the databases.
		// Wait for the postgresql object "acid-artifactory" to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-artifactory -n artifactory; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-artifactory" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=$(kubectl get postgresql acid-artifactory -n artifactory -o jsonpath="{.status.PostgresClusterStatus}"); while [ $DB_STATUS != "Running" ]; do sleep 5; done"`)
		require.NoError(t, err, output)
		timestampArtifactoryDb := time.Now().Add(time.Minute * 15).Add(time.Second * 10)
		// Wait for the "acid-confluence" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-confluence -n confluence; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-confluence" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=$(kubectl get postgresql acid-confluence -n confluence -o jsonpath=:"{.status.PostgresClusterStatus}"); while [ $DB_STATUS != "Running" ]; do sleep 5; done"`)
		require.NoError(t, err, output)
		timestampConfluenceDb := time.Now().Add(time.Minute * 15).Add(time.Second * 10)
		// Wait for the "acid-gitlab" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-gitlab -n gitlab; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-gitlab" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=$(kubectl get postgresql acid-gitlab -n gitlab -o jsonpath="{.status.PostgresClusterStatus}"); while [ $DB_STATUS != "Running" ]; do sleep 5; done"`)
		require.NoError(t, err, output)
		timestampGitlabDb := time.Now().Add(time.Minute * 15).Add(time.Second * 10)
		// Wait for the "acid-jira" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-jira -n jira; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-jira" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=$(kubectl get postgresql acid-jira -n jira -o jsonpath="{.status.PostgresClusterStatus}"); while [ $DB_STATUS != "Running" ]; do sleep 5; done"`)
		require.NoError(t, err, output)
		timestampJiraDb := time.Now().Add(time.Minute * 15).Add(time.Second * 10)

		// Wait for the "bigbang" kustomization to report "Ready==True". Our testing shows if everything goes right this should take 11-13 minutes.
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/bigbang -n flux-system --for=condition=Ready --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the "softwarefactoryaddons" kustomization to report "Ready==True".
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/softwarefactoryaddons -n flux-system --for=condition=Ready --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the bbcore-minio Statefulset "bbcore-minio-minio-instance-ss-0" to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset bbcore-minio-minio-instance-ss-0 -n bbcore-minio; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the bbcore-minio Statefulset "bbcore-minio-minio-instance-ss-0" to report that it is ready.
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/bbcore-minio-minio-instance-ss-0 -n bbcore-minio --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the gitlab-minio Statefulset "gitlab-minio-minio-instance-ss-0" to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset gitlab-minio-minio-instance-ss-0 -n gitlab-minio; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the gitlab-minio Statefulset "gitlab-minio-minio-instance-ss-0" to report that it is ready.
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/gitlab-minio-minio-instance-ss-0 -n gitlab-minio --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the velero-minio Statefulset "velero-minio-minio-instance-ss-0" to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset velero-minio-minio-instance-ss-0 -n velero-minio; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the velero-minio Statefulset "velero-minio-minio-instance-ss-0" to report that it is ready.
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/velero-minio-minio-instance-ss-0 -n velero-minio --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the GitLab Webservice Deployment to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get deployment gitlab-webservice-default -n gitlab; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the GitLab Webservice Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status deployment/gitlab-webservice-default -n gitlab --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the Velero Deployment to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get deployment velero-velero -n velero; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Velero Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status deployment/velero-velero -n velero --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the Restic Daemonset to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get daemonset restic -n velero; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Restic Daemonset to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status daemonset/restic -n velero --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Wait for the Jenkins StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset jenkins -n jenkins; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Jenkins StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/jenkins -n jenkins --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Jenkins is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl exec statefulset/jenkins -n jenkins -c jenkins -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Jira StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset jira -n jira; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Jira StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/jira -n jira --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Jira is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl exec statefulset/jira -n jira -c jira -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Confluence StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset confluence -n confluence; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Confluence StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/confluence -n confluence --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Confluence is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl exec statefulset/confluence -n confluence -c confluence -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Make sure flux is present.
		output, err = platform.RunSSHCommandAsSudo("flux --help")
		require.NoError(t, err, output)
		// Ensure that Jenkins is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://jenkins.bigbang.dev/login > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Confluence is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://confluence.bigbang.dev/status > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Jira is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://jira.bigbang.dev/status > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that GitLab is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Artifactory StatefulSet to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset artifactory -n artifactory; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the Artifactory StatefulSet to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/artifactory -n artifactory --watch --timeout=1200s`)
		require.NoError(t, err, output)
		// Ensure that Artifactory is able to talk to GitLab internally
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl exec statefulset/artifactory -n artifactory -c artifactory -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Artifactory is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://artifactory.bigbang.dev/artifactory/api/system/ping > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)

		// Wait for the Loki write Statefulset to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset logging-loki-write -n logging; do sleep 5; done"`)
		require.NoError(t, err, output)

		// Wait for the Loki wrie Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/logging-loki-write -n logging --watch --timeout=1200s`)
		require.NoError(t, err, output)

		// Wait for the Loki read Statefulset to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset logging-loki-read -n logging; do sleep 5; done"`)
		require.NoError(t, err, output)

		// Wait for the Loki read Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/logging-loki-read -n logging --watch --timeout=1200s`)
		require.NoError(t, err, output)

		// Wait for the Promtail Daemonset to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get daemonset logging-promtail -n logging; do sleep 5; done"`)
		require.NoError(t, err, output)

		// Wait for the Promtail Statefulset to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status daemonset/logging-promtail -n logging --watch --timeout=1200s`)
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

		// Ensure that the databases are still reporting "PostgresClusterStatus==Running"
		// Wait until timestampArtifactoryDb, then verify "PostgresClusterStatus==Running"
		for time.Now().Before(timestampArtifactoryDb) {
			time.Sleep(1 * time.Second)
		}
		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-artifactory -n artifactory -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-artifactory expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)
		// Wait until timestampConfluenceDb, then verify "PostgresClusterStatus==Running"
		for time.Now().Before(timestampConfluenceDb) {
			time.Sleep(1 * time.Second)
		}
		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-confluence -n confluence -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-confluence expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)
		// Wait until timestampGitlabDb, then verify "PostgresClusterStatus==Running"
		for time.Now().Before(timestampGitlabDb) {
			time.Sleep(1 * time.Second)
		}
		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-gitlab -n gitlab -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-gitlab expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)
		// Wait until timestampJiraDb, then verify "PostgresClusterStatus==Running"
		for time.Now().Before(timestampJiraDb) {
			time.Sleep(1 * time.Second)
		}
		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-jira -n jira -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-jira expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)
	})
}
