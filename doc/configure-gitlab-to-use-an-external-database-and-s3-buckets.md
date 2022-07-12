# Configure GitLab to use an external database and S3 buckets

Use this document to configure GitLab to use an external PostgreSQL database and external S3 buckets, rather than the package's provided Postgres Operator and MinIO

## External Database

1. In `kustomizations/bigbang/common/values.yaml`, update the values in the `addons.gitlab.values.global.psql` block:
   1. `host:` is the URL or IP address of the database
   1. `port:` is the port to use, usually `5432`
   1. `database:` is the name of the database to use. The recommended database name is `gitlabhq_production`, but a different name may be specified if desired.
   1. `username:` is the username to use. The recommended username is `gitlab`, but a different username may be specified if desired
   1. `password:` has subkeys `secret:` and `key:`. `secret:` is the name of the kubernetes secret where the user's password is kept, and `key:` is the name of the key in the secret where the password value is. This secret is not created by the package and must be added by the user.

## External S3 buckets

TODO: Write this section
