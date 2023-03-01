# Notes

mc cp local/postgres-backups/spilo/jira-acid-jira/ ./backup --recursive

envdir "/run/etc/wal-e.d/env" /scripts/postgres_backup.sh "/home/postgres/pgdata/pgroot/data"

envdir "/run/etc/wal-e.d/env" wal-g backup-list

kubectl run mc-tool -n jira --image=registry1.dso.mil/ironbank/big-bang/base:2.0.0 --command -- sleep infinity

kubectl cp $(which mc) jira/mc-tool:/home/base/mc

kubectl exec -i -n jira mc-tool -c mc-tool -- /home/base/mc alias list

kubectl exec -i -n jira mc-tool -c mc-tool -- /home/base/mc alias set postgres http://minio.postgres-minio.svc.cluster.local:80 $(kubectl get secret minio-user-creds -n postgres-minio -o jsonpath='{.data.CONSOLE_ACCESS_KEY}' | base64 -d) $(kubectl get secret minio-user-creds -n postgres-minio -o jsonpath='{.data.CONSOLE_SECRET_KEY}' | base64 -d)

kubectl exec -i -n jira mc-tool -c mc-tool -- /home/base/mc cp postgres/postgres-backups/spilo/jira-acid-jira/ /home/base/backup --recursive

kubectl exec -i -n jira mc-tool -c mc-tool -- bash -c 'cd /home/base/backup; tar -cvf ../backup.tar wal'

kubectl cp jira/mc-tool:/home/base/backup.tar ./backup.tar

## backup create

1. force run backup (wal-g) using kubectl exec
1. get backup info from cluster, minio, and wal-g (timestamp, uid, name(?))
1. download s3 bucket
1. zip up downloaded bucket and add to package

## backup restore

1. check with user if they would like to use in-cluster backup or backup from package
1. double check with user since data destruction is possible
1. if from package overwrite files in bucket with files from package
1. double check with user since following automated actions will result in possibly significant downtime
1. suspend swf kustomization
1. suspend jira hr
1. suspend database hr
1. scale jira to 0 pods
1. delete jira postgresql resource
1. ask user if they'd like to restore to specific time or to the latest backed up data
1. prompt user for timestamp if wanted
1. print required modification to kustomizations/softwarefactoryaddons/base/helmrelease.yaml with very specific instructions (line number?)
    - maybe just give them modified file to copy to destination?
1. show user what ref they're using to deploy from (branch, tag, etc.)
1. if tag explain they will need to cut a new tag
1. wait for user to confirm that changes are pushed to source repo
1. create and deploy zarf package to update in-cluster zarf gitea repo
1. if ref has changed patch it
1. resume swf kustomization
1. force reconcile swf kustomization
1. resume database hr
1. force reconcile database hr
1. wait for database pods to be up
1. resume jira hr
1. force reconcile jira hr
1. wait for jira pods to be up
1. restore successful?
