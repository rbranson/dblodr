# dblodr

# Install

```bash
$ go install github.com/rbranson/dblodr@latest
```

# Usage

```
dblodr is a simple SQL database load generator

Usage:
  dblodr [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize the target database for a load profile
  run         Execute load against the target database
```

## Example Run

```
$ dblodr run simpleinsert1 \
  --driver mysql \
  --dsn "myuser:mypassword@tcp(mysql.local)/mydatabase" \
  --table-name data \
  --columns host,data \
  --concurrency 50 \
  --iterations 1000
2025-01-30T12:49:17.817-0800	INFO	Starting Run	{"workload": "simpleinsert1"}
2025-01-30T12:49:18.817-0800	INFO	Run Stats	{"avg_runtime": "412.026ms", "errors": "0", "iterations": "102"}
2025-01-30T12:49:19.818-0800	INFO	Run Stats	{"avg_runtime": "236.344ms", "errors": "0", "iterations": "319"}
2025-01-30T12:49:20.818-0800	INFO	Run Stats	{"avg_runtime": "236.199ms", "errors": "0", "iterations": "536"}
2025-01-30T12:49:21.818-0800	INFO	Run Stats	{"avg_runtime": "236.726ms", "errors": "0", "iterations": "750"}
2025-01-30T12:49:22.818-0800	INFO	Run Stats	{"avg_runtime": "237.267ms", "errors": "0", "iterations": "956"}
2025-01-30T12:49:23.038-0800	INFO	Run Stats	{"avg_runtime": "235.599ms", "errors": "0", "iterations": "1000"}
```