# Notes

mc cp local/postgres-backups/spilo/jira-acid-jira/ ./backup --recursive

envdir "/run/etc/wal-e.d/env" /scripts/postgres_backup.sh "/home/postgres/pgdata/pgroot/data"

envdir "/run/etc/wal-e.d/env" wal-g backup-list