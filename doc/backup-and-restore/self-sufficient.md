# Backup and Restore of a "Self Sufficient" Cluster

:construction: **This document is under construction**. Please check back later.

Use this document to perform backup and restore of a cluster that is self sufficient, meaning that it uses MinIO for S3-compatible object storage and Postgres Operator for databases. This is the default functionality.

TODO: Write this doc. Plan right now is to use Velero to backup the whole cluster, with a note that the user will be responsible for shipping the Velero backups out of Minio to somewhere external to the cluster.
