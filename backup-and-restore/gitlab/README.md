# Backup and Restore Process for GitLab

GitLab is configured to automatically take backups, but since DI2-ME is designed to be deployable to air gaps, the backup artifacts still reside inside the cluster. They need to be extracted and kept somewhere safe.

## Backup Procedure

1. Get a terminal session on a Linux host that has direct `kubectl` access to the cluster.
1. Get the list of backups that exist in your cluster by running:

    ```shell
    kubectl exec -i -n gitlab -c toolbox $(kubectl get pod -n gitlab -l app=toolbox -o jsonpath='{.items[0].metadata.name}') -- s3cmd ls s3://gitlab-backups
    ```

    Expected output:

    ```shell
    ...
    2023-01-30 23:42       409600  s3://gitlab-backups/1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar
    ...
    ```

1. Take note of the backup filename that you want to extract. In the above example it is `1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar`.
1. Download the `zarf.yaml` file that is in this directory. Put it in a new empty directory on your host, and `cd` to that directory.
1. Ensure you have the `zarf` CLI installed. Use the same version that is listed at the top of the Makefile in the root of this repository.
1. Create the backup package by running:

    > NOTE: Use the filename that you noted above, not this exact command!

    ```shell
    zarf package create --set BACKUP_FILENAME=1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar --set DELETE_REMOTE_BACKUP_FILE=no --confirm
    ```

    > NOTE: If you'd like the in-cluster backup file to be deleted after the package is created, you can set `DELETE_REMOTE_BACKUP_FILE` to `yes`.

This will create a file called `zarf-package-di2me-gitlab-restorable-backup-amd64-1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar.tar.zst` that contains all necessary items to perform a full restore of GitLab. Save it to somewhere safe.

## Restore Procedure

1. Get a terminal session on a Linux host that has direct `kubectl` access to the cluster.
1. Copy the `zarf-package-di2me-gitlab-restorable-backup-amd64-1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar.tar.zst` file to an empty directory on the host.
1. Ensure you have the `zarf` CLI installed. Use the same version that is listed at the top of the Makefile in the root of this repository.
1. Extract the backup files by running:

    ```shell
    zarf package deploy <ThePackageFilename> --confirm
    ```

    A backup tarball and 2 kubernetes yaml files will be extracted to the current directory.

1. Prepare the database for the restore by running:

    ```shell
    kubectl exec -it -n gitlab acid-gitlab-0 -- psql -c "DROP EXTENSION IF EXISTS pg_auth_mon; DROP EXTENSION IF EXISTS pg_stat_kcache; DROP EXTENSION IF EXISTS pg_stat_statements CASCADE;"
    ```

1. Copy the database file to the toolbox pod

    ```shell
    kubectl cp -c toolbox ./1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar gitlab/$(kubectl get pod -n gitlab -l app=toolbox -o jsonpath='{.items[0].metadata.name}'):home/git/1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar
    ```

1. Perform the restore

    ```shell
    # First log into the toolbox pod
    kubectl exec -it -n gitlab -c toolbox $(kubectl get pod -n gitlab -l app=toolbox -o jsonpath='{.items[0].metadata.name}') -- bash
    # Now that you are inside the toolbox pod, run:
    backup-utility --restore -f file:///home/git/1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar
    rm /home/git/1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar
    # Now exit out of the pod
    exit
    ```

1. Apply the 2 yaml files

    ```shell
    kubectl delete -f gitlab-gitlab-initial-root-password.yaml
    kubectl create -f gitlab-gitlab-initial-root-password.yaml
    kubectl delete -f gitlab-rails-secret.yaml
    kubectl create -f gitlab-rails-secret.yaml
    ```

1. Restart the gitlab pods that are affected by the new data

    ```shell
    kubectl delete pods -n gitlab -lapp=sidekiq,release=gitlab
    kubectl delete pods -n gitlab -lapp=webservice,release=gitlab
    ```

1. Re-enable database extensions that were previously disabled (they stop the restore from working if they are enabled)

    ```shell
    kubectl exec -it -n gitlab acid-gitlab-0 -- psql -c "CREATE EXTENSION IF NOT EXISTS pg_auth_mon; CREATE EXTENSION IF NOT EXISTS pg_stat_kcache; CREATE EXTENSION IF NOT EXISTS pg_stat_statements;"
    ```

    > NOTE: The restore process may have already re-enabled these extensions. If it says they already exist, you can ignore the errors.

1. Clean up

    ```shell
    rm -rf ./1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar
    rm -rf ./gitlab-gitlab-initial-root-password.yaml
    rm -rf ./gitlab-rails-secret.yaml
    ```
