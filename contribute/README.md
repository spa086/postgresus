### Prerequisites

1. Read docs in /docs folder, README.md in /backend and /frontend folders
2. Run both backend and frontend following the instructions in their respective README.md files (for development)
3. Read this file till the end

### How to create a pull request?

We use gitflow approach.

1. Create a new branch from main
2. Make changes
3. Create a pull request to main
4. Wait for review
5. Merge pull request

Commits should be named in the following format depending on the type of change:

- `FEATURE (area): What was done`
- `FIX (area): What was fixed`
- `REFACTOR (area): What was refactored`

To see examples, look at commit history in main branch.

Branches should be named in the following format:

- `feature/what_was_done`
- `fix/what_was_fixed`
- `refactor/what_was_refactored`

Example:

- `feature/add_support_of_kubernetes_helm`
- `fix/make_healthcheck_optional`
- `refactor/refactor_navbar`

Before any commit, make sure:

1. You created critical tests for your changes
2. `golangci-lint fmt` and `golangci-lint run` are passing
3. All tests are passing
4. Project is building successfully

If you need to add some explanation, do it in appropriate place in the code. Or in the /docs folder if it is something general. For charts, use Mermaid.

### Priorities

Before taking anything more than a couple of lines of code, please write Rostislav via Telegram (@rostislav_dugin) and confirm priority. It is possible that we already have something in the works, it is not needed or it's not project priority.

Deploy flow:

- add support of Kubernetes Helm (in progress by Rostislav Dugin)
- add devcontainers for backend and frontend

Backups flow:

- add FTP
- add Dropbox
- add OneDrive
- add NAS
- add Yandex Drive
- think about pg_dumpall / pg_basebackup / WAL backup / incremental backups
- add encryption for backups
- add support of PgBouncer

Notifications flow:

- add Mattermost
- add MS Teams

Extra:

- add prettier labels to GitHub README
- allow to download backup file (via streaming)
- add linters and formatters on each PR
- add versioning instead of :latest
- create pretty website like rybbit.io with demo
- add HTTPS for Postgresus
- add simple SQL queries via UI
- add brute force protection on auth (via local RPS limiter)

Monitoring flow:

- add system metrics (CPU, RAM, disk, IO)
- add queries stats (slowest, most frequent, etc. via pg_stat_statements)
- add alerting for slow queries (listen for slow query and if they reach >100ms - send message)
- add alerting for high resource usage (listen for high resource usage and if they reach >90% - send message)
- add DB size distribution chart (tables, indexes, etc.)
- add performance test for DB (to compare DBs on different clouds and VPS)
- add DB metrics (pg_stat_activity, pg_locks, pg_stat_database)
- add chart of connections (from IPs, apps names, etc.)
- add chart of transactions (TPS)
- deadlocks chart
- chart of connection attempts (to see crash loops)
- add chart of IDLE transactions VS executing transactions
- show queries that take the most IO time (suboptimal indexes)
- show chart by top IO / CPU queries usage (see page 90 of the PostgreSQL monitoring book)

```
exec_time | IO   | CPU | query
105 hrs   | 73%  | 27% | SELECT * FROM users;
```

- chart of read / update / delete / insert queries
- chart with deadlocks, conflicts, rollbacks (see page 115 of the PostgreSQL monitoring book)
- stats of buffer usage
- status of IO (DB, indexes, sequences)
- % of cache hit
- replication stats
