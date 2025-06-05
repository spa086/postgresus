package usecases

import (
	"errors"
	usecases_postgresql "postgresus-backend/internal/features/backups/usecases/postgresql"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/storages"

	"github.com/google/uuid"
)

type CreateBackupUsecase struct {
	CreatePostgresqlBackupUsecase *usecases_postgresql.CreatePostgresqlBackupUsecase
}

// Execute creates a backup of the database and returns the backup size in MB
func (uc *CreateBackupUsecase) Execute(
	backupID uuid.UUID,
	database *databases.Database,
	storage *storages.Storage,
	backupProgressListener func(
		completedMBs float64,
	),
) error {
	if database.Type == databases.DatabaseTypePostgres {
		return uc.CreatePostgresqlBackupUsecase.Execute(
			backupID,
			database,
			storage,
			backupProgressListener,
		)
	}

	return errors.New("database type not supported")
}
