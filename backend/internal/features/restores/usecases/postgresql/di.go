package usecases_postgresql

import (
	"postgresus-backend/internal/util/logger"
)

var restorePostgresqlBackupUsecase = &RestorePostgresqlBackupUsecase{
	logger.GetLogger(),
}

func GetRestorePostgresqlBackupUsecase() *RestorePostgresqlBackupUsecase {
	return restorePostgresqlBackupUsecase
}
