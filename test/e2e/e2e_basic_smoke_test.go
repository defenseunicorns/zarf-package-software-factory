package test_test

import (
	"testing"

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

		// DISABLE-ARTIFACTORY
		// // Wait for the postgresql object "acid-artifactory" to exist.
		// output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-artifactory -n artifactory; do sleep 5; done"`)
		// require.NoError(t, err, output)
		// // Wait for the "acid-artifactory" database to report "PostgresClusterStatus==Running", then set the timestamp
		// output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=\$(kubectl get postgresql acid-artifactory -n artifactory -o jsonpath=\"{.status.PostgresClusterStatus}\"); while [ \"\$DB_STATUS\" != \"Running\" ]; do sleep 5; DB_STATUS=\$(kubectl get postgresql acid-artifactory -n artifactory -o jsonpath=\"{.status.PostgresClusterStatus}\"); done"`)
		// require.NoError(t, err, output)

		// Wait for the "acid-confluence" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-confluence -n confluence; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-confluence" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=\$(kubectl get postgresql acid-confluence -n confluence -o jsonpath=\"{.status.PostgresClusterStatus}\"); while [ \"\$DB_STATUS\" != \"Running\" ]; do sleep 5; DB_STATUS=\$(kubectl get postgresql acid-confluence -n confluence -o jsonpath=\"{.status.PostgresClusterStatus}\"); done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-gitlab" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-gitlab -n gitlab; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-gitlab" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=\$(kubectl get postgresql acid-gitlab -n gitlab -o jsonpath=\"{.status.PostgresClusterStatus}\"); while [ \"\$DB_STATUS\" != \"Running\" ]; do sleep 5; DB_STATUS=\$(kubectl get postgresql acid-gitlab -n gitlab -o jsonpath=\"{.status.PostgresClusterStatus}\"); done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-jira" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-jira -n jira; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-jira" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=\$(kubectl get postgresql acid-jira -n jira -o jsonpath=\"{.status.PostgresClusterStatus}\"); while [ \"\$DB_STATUS\" != \"Running\" ]; do sleep 5; DB_STATUS=\$(kubectl get postgresql acid-jira -n jira -o jsonpath=\"{.status.PostgresClusterStatus}\"); done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-sonarqube" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-sonarqube -n sonarqube; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-sonarqube" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=\$(kubectl get postgresql acid-sonarqube -n sonarqube -o jsonpath=\"{.status.PostgresClusterStatus}\"); while [ \"\$DB_STATUS\" != \"Running\" ]; do sleep 5; DB_STATUS=\$(kubectl get postgresql acid-sonarqube -n sonarqube -o jsonpath=\"{.status.PostgresClusterStatus}\"); done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-keycloak" database to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get postgresql acid-keycloak -n keycloak; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the "acid-keycloak" database to report "PostgresClusterStatus==Running", then set the timestamp
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "DB_STATUS=\$(kubectl get postgresql acid-keycloak -n keycloak -o jsonpath=\"{.status.PostgresClusterStatus}\"); while [ \"\$DB_STATUS\" != \"Running\" ]; do sleep 5; DB_STATUS=\$(kubectl get postgresql acid-keycloak -n keycloak -o jsonpath=\"{.status.PostgresClusterStatus}\"); done"`)
		require.NoError(t, err, output)

		// In order for the GitLab HelmRelease to fully reconcile, we need to manually trigger a backup, so that 2 remaining PVCs will bind.
		// Wait for the "gitlab-toolbox-backup" CronJob to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 2400 bash -c "while ! kubectl get cronjob gitlab-toolbox-backup -n gitlab; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Trigger the backup
		output, err = platform.RunSSHCommandAsSudo(`kubectl create job -n gitlab --from=cronjob/gitlab-toolbox-backup gitlab-toolbox-backup-manual`)
		require.NoError(t, err, output)

		// Wait for the "bigbang" kustomization to report "Ready==True". Our testing shows if everything goes right this should take 15-25 minutes.
		output, err = platform.RunSSHCommandAsSudo(`kubectl wait kustomization/bigbang -n flux-system --for=condition=Ready --timeout=2400s`)
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
		// Wait for the GitLab Webservice Deployment to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get deployment gitlab-webservice-default -n gitlab; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Wait for the GitLab Webservice Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status deployment/gitlab-webservice-default -n gitlab --watch --timeout=1200s`)
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
		// Setup DNS records for cluster services
		output, err = platform.RunSSHCommandAsSudo("cd ~/app && test/metallb/dns.sh")
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
		// Ensure that SonarQube is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://sonarqube.bigbang.dev/login > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Neuvector is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://neuvector.bigbang.dev/#/login > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Jaeger is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://tracing.bigbang.dev/search > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that AlertManager is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://alertmanager.bigbang.dev/#/alerts > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Grafana is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://grafana.bigbang.dev/login > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Kiali is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://kiali.bigbang.dev/kiali/ > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Prometheus is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://prometheus.bigbang.dev/graph > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Mattermost is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://chat.bigbang.dev/login > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// Ensure that Nexus is available outside of the cluster.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://nexus.bigbang.dev/ > /dev/null; do sleep 5; done"`)
		require.NoError(t, err, output)
		// // Ensure that keycloak is available outside of the cluster.
		// output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://keycloak.bigbang.dev/login > /dev/null; do sleep 5; done"`)
		// require.NoError(t, err, output)

		// DISABLE-ARTIFACTORY
		// // Wait for the Artifactory StatefulSet to exist.
		// output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset artifactory -n artifactory; do sleep 5; done"`)
		// require.NoError(t, err, output)
		// // Wait for the Artifactory StatefulSet to report that it is ready
		// output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/artifactory -n artifactory --watch --timeout=1200s`)
		// require.NoError(t, err, output)
		// // Ensure that Artifactory is able to talk to GitLab internally
		// output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl exec statefulset/artifactory -n artifactory -c artifactory -- curl -L -s --fail --show-error https://gitlab.bigbang.dev/-/health > /dev/null; do sleep 5; done"`)
		// require.NoError(t, err, output)
		// // Ensure that Artifactory is available outside of the cluster.
		// output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! curl -L -s --fail --show-error https://artifactory.bigbang.dev/artifactory/api/system/ping > /dev/null; do sleep 5; done"`)
		// require.NoError(t, err, output)

		// Wait for the Loki write Statefulset to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get statefulset logging-loki-write -n logging; do sleep 5; done"`)
		require.NoError(t, err, output)

		// Wait for the Loki write Statefulset to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status statefulset/logging-loki-write -n logging --watch --timeout=1200s`)
		require.NoError(t, err, output)

		// Wait for the Loki read Deployment to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get deployment logging-loki-read -n logging; do sleep 5; done"`)
		require.NoError(t, err, output)

		// Wait for the Loki read Deployment to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status deployment/logging-loki-read -n logging --watch --timeout=1200s`)
		require.NoError(t, err, output)

		// Wait for the Promtail Daemonset to exist.
		output, err = platform.RunSSHCommandAsSudo(`timeout 1200 bash -c "while ! kubectl get daemonset promtail-promtail -n promtail; do sleep 5; done"`)
		require.NoError(t, err, output)

		// Wait for the Promtail Daemonset to report that it is ready
		output, err = platform.RunSSHCommandAsSudo(`kubectl rollout status daemonset/promtail-promtail -n promtail --watch --timeout=1200s`)
		require.NoError(t, err, output)

		// Ensure that the services do not accept discontinued TLS versions. If they reject TLSv1.1 it is assumed that they also reject anything below TLSv1.1.
		// Ensure that GitLab does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan gitlab.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Jenkins does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan jenkins.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Jira does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan jira.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Confluence does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan confluence.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)

		// DISABLE-ARTIFACTORY
		// // Ensure that Artifactory does not accept TLSv1.1
		// output, err = platform.RunSSHCommandAsSudo(`sslscan artifactory.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		// require.NoError(t, err, output)

		// Ensure that Sonarqube does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan sonarqube.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Neuvector does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan neuvector.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Jaeger does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan tracing.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Kiali does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan kiali.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that AlertManager does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan alertmanager.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Grafana does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan grafana.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Prometheus does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan prometheus.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Mattermost does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan chat.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// Ensure that Nexus does not accept TLSv1.1
		output, err = platform.RunSSHCommandAsSudo(`sslscan nexus.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		require.NoError(t, err, output)
		// // Ensure that Keycloak does not accept TLSv1.1
		// output, err = platform.RunSSHCommandAsSudo(`sslscan keycloak.bigbang.dev | grep "TLSv1.1" | grep "disabled"`)
		// require.NoError(t, err, output)

		// Ensure that the databases are still reporting "PostgresClusterStatus==Running"

		// DISABLE-ARTIFACTORY
		// output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-artifactory -n artifactory -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-artifactory expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		// require.NoError(t, err, output)

		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-confluence -n confluence -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-confluence expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)

		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-gitlab -n gitlab -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-gitlab expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)

		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-jira -n jira -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-jira expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)

		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-sonarqube -n sonarqube -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-sonarqube expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)

		output, err = platform.RunSSHCommandAsSudo(`DB_STATUS=$(kubectl get postgresql acid-keycloak -n keycloak -o jsonpath="{.status.PostgresClusterStatus}"); if [ "$DB_STATUS" != "Running" ]; then echo "Status of database acid-keycloak expected to be Running, but got $DB_STATUS"; exit 1; fi`)
		require.NoError(t, err, output)

		// DISABLE-ARTIFACTORY
		// // Create backup for Artifactory
		// output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/artifactory; ~/app/build/zarf p c --confirm --set BACKUP_TIMESTAMP=""`)
		// require.NoError(t, err, output)
		// // Start restore process for Artifactory
		// output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/artifactory; mkdir test; mv zarf-package* test; cd test;  ~/app/build/zarf p d zarf-package* --components warning-downtime-begin-restore --confirm`)
		// require.NoError(t, err, output)

		// Create backup for Confluence
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/confluence; ~/app/build/zarf p c --confirm --set BACKUP_TIMESTAMP=""`)
		require.NoError(t, err, output)
		// Start restore process for Confluence
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/confluence; mkdir test; mv zarf-package* test; cd test;  ~/app/build/zarf p d zarf-package* --components warning-downtime-begin-restore --confirm`)
		require.NoError(t, err, output)

		// Create backup for GitLab
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/gitlab; ~/app/build/zarf p c --confirm --set BACKUP_FILENAME="$(kubectl exec -i -n gitlab -c toolbox $(kubectl get pod -n gitlab -l app=toolbox -o jsonpath="{.items[0].metadata.name}") -- s3cmd ls s3://gitlab-backups | awk "{split(\$NF,a,\"/\"); print a[length(a)]; exit}")" --set DELETE_REMOTE_BACKUP_FILE="no"`)
		require.NoError(t, err, output)
		// Start restore process for Gitlab
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/gitlab; mkdir test; mv zarf-package* test; cd test;  ~/app/build/zarf p d zarf-package* --components warning-downtime-begin-restore --confirm`)
		require.NoError(t, err, output)

		// Create backup for Jenkins
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/jenkins; ~/app/build/zarf p c --confirm --set BACKUP_TIMESTAMP=""`)
		require.NoError(t, err, output)
		// Start restore process for Jenkins
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/jenkins; mkdir test; mv zarf-package* test; cd test;  ~/app/build/zarf p d zarf-package* --components warning-downtime-begin-restore --confirm`)
		require.NoError(t, err, output)

		// Create backup for Jira
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/jira; ~/app/build/zarf p c --confirm --set BACKUP_TIMESTAMP=""`)
		require.NoError(t, err, output)
		// Start restore process for Jira
		output, err = platform.RunSSHCommandAsSudo(`cd ~/app/backup-and-restore/jira; mkdir test; mv zarf-package* test; cd test;  ~/app/build/zarf p d zarf-package* --components warning-downtime-begin-restore --confirm`)
		require.NoError(t, err, output)
	})
}
