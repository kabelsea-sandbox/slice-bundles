# clickhouse bundle

Bundle provide configured clickhouse native client.

## Environment

- `CLICKHOUSE_ADDRS`                    (required, default=clickhouse:9000, list)
- `CLICKHOUSE_DEBUG`                    (default=false)
- `CLICKHOUSE_DIAL_TIMEOUT`             (default=5s)
- `CLICKHOUSE_MAX_OPEN_CONNS`           (default=2)
- `CLICKHOUSE_MAX_IDLE_CONNS`           (default=2)
- `CLICKHOUSE_CONN_MAX_LIFETIME`        (default=10m)
- `CLICKHOUSE_BLOCK_BUFFER_SIZE`        (default=10)

- `CLICKHOUSE_AUTH_DATABASE`            (required)
- `CLICKHOUSE_AUTH_USERNAME`            (required, default=default)
- `CLICKHOUSE_AUTH_PASSWORD`            ()

- `CLICKHOUSE_SETTINGS_MAX_EXECUTION_TIME`  (default=15)
