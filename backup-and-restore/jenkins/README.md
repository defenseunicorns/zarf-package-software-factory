# Backup and Restore Process for Jenkins

The backup process for Jenkins is completely managed through this zarf package. Backups are not taken unless this zarf package is created. The package itself is the backup artifact, so after creation it should be stored somewhere safe. Please read through this guide in it's entirety before attempting any operations.

## Backup Procedure

1. Get a terminal session on a Linux host that has direct `kubectl` access to the cluster.
1. Download the `zarf.yaml` file and the `files` folder that is in this directory. Put it in a new empty directory on your host, and `cd` to that directory.
1. Ensure you have the `zarf` CLI installed. Use the same version that is listed at the top of the Makefile in the root of this repository.
1. Create the backup package by running:

    ```shell
    zarf package create --set BACKUP_TIMESTAMP="$(date --iso-8601=seconds)" --confirm
    ```

This will create a file with a similar name to `zarf-package-di2me-jenkins-restorable-backup-amd64-1970-01-01T00:00:00+00:00.tar.zst` that contains all necessary items to perform a full restore of Jenkins. Save it to somewhere safe.

## Restore Procedure

1. Get a terminal session on a Linux host that has direct `kubectl` access to the cluster.
1. Copy the zarf package you wish to restore from to an empty directory on the host. The name of the package file will be similar to `zarf-package-di2me-jenkins-restorable-backup-amd64-1970-01-01T00:00:00+00:00.tar.zst` except with the timestamp being when the backup package was created.
1. Ensure you have the `zarf` CLI installed. Use the same version that is listed at the top of the Makefile in the root of this repository.
1. **Warning!** The next step will cause downtime until Jenkins is done being restored!
1. Begin the restore operation using zarf, replacing `<ThePackageFilename>` with the filename of the package you want to restore from:

    ```shell
    zarf package deploy <ThePackageFilename> --components=warning-downtime-begin-restore --confirm
    ```

1. As long as Jenkins comes up healthy you have now restored it! You may now delete any files you do not wish to keep from the restore operations.
