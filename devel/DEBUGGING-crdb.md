## Enable connection, SQL statement execution logs

logs can be found in the container under `cockroach-data/logs/`

```
SHOW CLUSTER SETTING server.auth_log.sql_connections.enabled;
SET CLUSTER SETTING sql.trace.log_statement_execute=true;
```

## Backup database

back will be placed under `cockroach-data/external`
```
backup database defaultdb into 'nodelocal://self/tmp/';
```

