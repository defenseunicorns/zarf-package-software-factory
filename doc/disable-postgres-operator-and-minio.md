# Disable Postgres Operator and Minio

This document will explain how to disable Postgres Operator and Minio, which should be done if you have chosen to use external database and object storage services such as AWS RDS and S3.

## Disable Postgres Operator

1. Comment out or delete the component named `postgres-operator` in `zarf.yaml`

1. Delete `kustomizations/softwarefactoryaddons/base/postgresql.yaml`

## Disable Minio

1. In `zarf.yaml` in the component named `big-bang`:
   1. In the `repos:` section, comment out or delete the git repos `minio-operator.git` and `minio.git`
   1. In the `images:` section, delete any docker images that are only used by MinIO
1. In `kustomizations/bigbang/common/values.yaml`:
   1. Change `addons.minioOperator.enabled` to `false`
   1. Change `addons.minio.enabled` to `false`
