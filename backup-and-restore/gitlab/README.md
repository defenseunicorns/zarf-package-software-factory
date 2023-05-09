# Backup and Restore Process for GitLab

GitLab is configured to automatically take backups, but since DI2-ME is designed to be deployable to air gaps, the backup artifacts still reside inside the cluster. They need to be extracted and kept somewhere safe. Please read through this guide in it's entirety before attempting any operations.

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
1. Perform the restore and extract files:

    ```shell
    zarf package deploy <ThePackageFilename> --components=warning-downtime-begin-restore --confirm
    ```

1. Clean up files if not needed:

    ```shell
    rm -rf ./1675119631_2023_01_30_15.7.0-ee_gitlab_backup.tar
    rm -rf ./gitlab-gitlab-initial-root-password.yaml
    rm -rf ./gitlab-rails-secret.yaml
    ```
