# 3. Every Application That Needs Buckets Gets Its Own MinIO

Date: 2022-10-21

## Status

Accepted

## Context

Using a single MinIO instance for all applications that need buckets creates a difficult situation for backup and restore by creating a blast radius that overlaps with other applications. Since Velero is the chosen backup and restore solution, and since it does those backups at a namespace level, if I am trying to restore a single application I have to restore the entire MinIO instance. That means I have to restore all the buckets for all the applications that use MinIO. If another application is also using MinIO, it will have its buckets overwritten by the restore.

## Decision

Each application that needs buckets will get its own MinIO instance. This will allow us to backup and restore each application individually without affecting other applications by creating a 1:1 relationship between namespaces in the cluster and applications to be backed up. For example, to restore GitLab, all I will have to do is restore the GitLab namespace from a backup.

## Consequences

* It will be much easier to backup and restore individual applications without affecting other applications in the cluster.
* Having 1 MinIO per application will be more resource intensive and create a higher maintenance burden.
* Complexity of the system will be reduced overall, since while we do have more pods running in the cluster, they will be the same pods with the same (or very similar) configuration.
