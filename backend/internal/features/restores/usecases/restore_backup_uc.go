package usecases

import (
	"errors"
	"postgresus-backend/internal/features/backups"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/restores/models"
	usecases_postgresql "postgresus-backend/internal/features/restores/usecases/postgresql"
)

type RestoreBackupUsecase struct {
	RestorePostgresqlBackupUsecase *usecases_postgresql.RestorePostgresqlBackupUsecase
}

func (uc *RestoreBackupUsecase) Execute(
	restore models.Restore,
	backup *backups.Backup,
) error {
	if restore.Backup.Database.Type == databases.DatabaseTypePostgres {
		return uc.RestorePostgresqlBackupUsecase.Execute(restore, backup)
	}

	return errors.New("database type not supported")
}
