package usecases

import (
	usecases_postgresql "postgresus-backend/internal/features/backups/backups/usecases/postgresql"
)

var createBackupUsecase = &CreateBackupUsecase{
	usecases_postgresql.GetCreatePostgresqlBackupUsecase(),
}

func GetCreateBackupUsecase() *CreateBackupUsecase {
	return createBackupUsecase
}
