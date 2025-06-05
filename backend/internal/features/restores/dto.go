package restores

import (
	"postgresus-backend/internal/features/databases/databases/postgresql"
)

type RestoreBackupRequest struct {
	PostgresqlDatabase *postgresql.PostgresqlDatabase `json:"postgresqlDatabase"`
}
