# Configure Confluence to use an external database

Use this document to configure Confluence to use an external PostgreSQL database, rather than the package's provided Postgres Operator.

## Instructions

1. In `kustomizations/softwarefactoryaddons/confluence/common-values.yaml`, update the values in the `confluence.additionalEnvironmentVariables` block:
   1. `ATL_DB_TYPE` is the database type. This should stay as it is. `postgresql` is the value that tells Confluence to use a Postgres database.
   1. `ATL_JDBC_URL` is the connection string to your external database
   1. `ATL_JDBC_USER` is the database username. You can use your own secret if you want, or you can change the `valueFrom:` to just `value:` and hard code the value in the yaml file. Make sure you encrypt the file with SOPS if you are hard coding secrets in it.
   1. `ATL_JDBC_PASSWORD` is the database password. You can use your own secret if you want, or you can change the `valueFrom:` to just `value:` and hard code the value in the yaml file. Make sure you encrypt the file with SOPS if you are hard coding secrets in it.
