package tools

type PostgresqlVersion string

const (
	PostgresqlVersion13 PostgresqlVersion = "13"
	PostgresqlVersion14 PostgresqlVersion = "14"
	PostgresqlVersion15 PostgresqlVersion = "15"
	PostgresqlVersion16 PostgresqlVersion = "16"
	PostgresqlVersion17 PostgresqlVersion = "17"
)

type PostgresqlExecutable string

const (
	PostgresqlExecutablePgDump PostgresqlExecutable = "pg_dump"
	PostgresqlExecutablePsql   PostgresqlExecutable = "psql"
)
