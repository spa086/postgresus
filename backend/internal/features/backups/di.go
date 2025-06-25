package backups

import (
	"postgresus-backend/internal/features/backups/usecases"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/logger"
	"time"
)

var backupRepository = &BackupRepository{}
var backupService = &BackupService{
	databases.GetDatabaseService(),
	storages.GetStorageService(),
	backupRepository,
	notifiers.GetNotifierService(),
	notifiers.GetNotifierService(),
	usecases.GetCreateBackupUsecase(),
	logger.GetLogger(),
}

var backupBackgroundService = &BackupBackgroundService{
	backupService,
	backupRepository,
	databases.GetDatabaseService(),
	storages.GetStorageService(),
	time.Now().UTC(),
	logger.GetLogger(),
}

var backupController = &BackupController{
	backupService,
	users.GetUserService(),
}

func SetupDependencies() {
	databases.
		GetDatabaseService().
		SetDatabaseStorageChangeListener(backupService)
}

func GetBackupService() *BackupService {
	return backupService
}

func GetBackupController() *BackupController {
	return backupController
}

func GetBackupBackgroundService() *BackupBackgroundService {
	return backupBackgroundService
}
