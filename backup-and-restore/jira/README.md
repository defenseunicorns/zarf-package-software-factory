# Backup and Restore Process for Jira

Jira is configured to automatically take backups, but since DI2-ME is designed to be deployable to air gaps, the backup artifacts still reside inside the cluster. They need to be extracted and kept somewhere safe. Please read through this guide in it's entirety before attempting any operations.

## Backup Procedure

1. Get a terminal session on a Linux host that has direct `kubectl` access to the cluster.
1. Download the `zarf.yaml` file and the `files` folder that is in this directory. Put it in a new empty directory on your host, and `cd` to that directory.
1. Ensure you have the `zarf` CLI installed. Use the same version that is listed at the top of the Makefile in the root of this repository.
1. Create the backup package by running:

    ```shell
    zarf package create --set BACKUP_TIMESTAMP=$(date --iso-8601=seconds) --confirm
    ```

This will create a file with a similar name to `zarf-package-di2me-jira-restorable-backup-amd64-1970-01-01T00:00:00+00:00.tar.zst` that contains all necessary items to perform a full restore of Jira. Save it to somewhere safe.

## Restore Procedure

1. Get a terminal session on a Linux host that has direct `kubectl` access to the cluster.
1. Copy the zarf package you wish to restore from to an empty directory on the host. The name of the package file will be similar to `zarf-package-di2me-jira-restorable-backup-amd64-1970-01-01T00:00:00+00:00.tar.zst` except with the timestamp being when the backup package was created.
1. Ensure you have the `zarf` CLI installed. Use the same version that is listed at the top of the Makefile in the root of this repository.
1. **Warning!** The next step will cause downtime until this guide is finished for restoring Jira!
1. Begin the restore operation using zarf:

    ```shell
    zarf package deploy <ThePackageFilename> --components=warning-downtime-begin-restore --confirm
    ```

    This has temporarily stopped Jira, removed the jira-database pods, and uploaded the backup data into the internal minio bucket for them.

1. You will now need to make modifications to the di2me repo and push those changes. 

    **Attention!** The following steps, until noted otherwise, require an environment with access to the repository that you host di2me out of, this may or may not be different to the environment that has access to your cluster depending on how you've decided to deploy di2me.

    1. In `kustomizations/softwarefactoryaddons/base/databases/jira.yaml` There is the following commented section:

        ```yaml
                # clone:
                #   cluster: "acid-jira"
                #   timestamp: ""
                #   uid: ""
        ```

        Notice the indentation, it is crucial that it is correct when deploying.

    1. Inside the quotes after `timestamp` paste in the output of running `date --iso-8601=seconds` in the terminal.
    1. Inside the quotes after `uid` paste in the contents of the `postgres-cluster-uid` file.
    1. On each line with a `#` delete the `#` and the following space, an example of the section being valid:

        ```yaml
                clone:
                  cluster: "acid-jira"
                  timestamp: "2023-03-20T15:57:27-05:00"
                  uid: "ef59344b-cacd-4bd7-9b50-2da144586638"
        ```

    1. Again double check that the indentation is correct, 8 spaces before `clone` and 10 for the following lines.
    1. Next you may need to change the ref used by di2me to deploy depending on how you deploy it.
    1. in `manifests/setup.yaml` at the bottom there is a section called `ref:`, under that is either a branch or a tag. Update the value of the branch or the tag to wherever you will be pushing these changes to.
1. Push and optionally tag these changes.
1. Change the branch/tag of your local copy of the repo to what you pushed/tagged to in the previous step (it may already be correct in the case of a branch)
1. Run the following command and note it's output. You will be watching for this value in a later step.

    ```shell
    echo "$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)/$(git rev-parse HEAD)"
    ```

1. In the day2 folder of the repository there is a zarf.yaml, the package this creates is what will be used to push the changes that have been made into the cluster.
1. While your terminal is in the day2 folder run the following command, replacing `<repo>` with the url to your di2me repository (e.g. `https://github.com/defenseunicorns/zarf-package-software-factory.git`)

    ```shell
    zarf package create --confirm --set DI2ME_REPO="<repo>"
    ```

    This will have created a file called `zarf-package-day-two-update-amd64.tar.zst`

1. If pushing changes to your di2me repository required a different environment, transfer the package that the previous command created to a machine with access to your cluster. 

    **Attention!** All commands for the rest of this guide must be run from a machine with access to the di2me cluster.

1. Place the day2 package into an empty folder and deploy it using the following command.

    ```shell
    zarf package deploy zarf-package-day-two-update-amd64.tar.zst --confirm
    ```

1. Once it has deployed run the following command until the output matches the string you noted from step 9. This should take less than a minute.

    ```shell
    kubectl get gitrepo zarf-package-software-factory -n flux-system -o=jsonpath='{.status.artifact.revision}{"\n"}'
    ```

    Alternatively if you are able to easily copy and paste the value from step 9 you can run the following command, replacing `<revision>` with the value from step 9. This will let you know when the value has changed instead of checking manually.

    ```shell
    kubectl wait --for=jsonpath='{.status.artifact.revision}'="<revision>" -n flux-system gitrepo/zarf-package-software-factory --timeout=300s
    ```

    If this step has taken more than 2 minutes go back to step 6 and try again.

1. Now we will unsuspend the flux kustomization to allow the new changes to propogate through the cluster, run the following command.

    ```shell
    flux resume kustomization softwarefactoryaddons
    ```

1. Now we will wait for the new database to be ready. Run the following command

    ```shell
    kubectl wait --for=jsonpath='{.status.PostgresClusterStatus}'='Running' -n jira postgresql/acid-jira --timeout=300s
    ```

    If you get an error like `Error from server (NotFound): postgresqls.acid.zalan.do "acid-jira" not found` just run the command again.

1. Now that the databases are up we will reenable jira. Run the following commands.

    ```shell
    kubectl scale --replicas=1 -n jira statefulset/jira
    kubectl wait --for=jsonpath='{.status.availableReplicas}'=1 -n jira statefulset/jira --timeout=300s
    ```

1. As long as Jira comes up healthy you have now restored it! You may now delete any files you do not wish to keep from the restore operations.
