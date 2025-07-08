package backups_config

import (
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/features/users"
)

var backupConfigRepository = &BackupConfigRepository{}
var backupConfigService = &BackupConfigService{
	backupConfigRepository,
	databases.GetDatabaseService(),
	storages.GetStorageService(),
	nil,
}
var backupConfigController = &BackupConfigController{
	backupConfigService,
	users.GetUserService(),
}

func GetBackupConfigController() *BackupConfigController {
	return backupConfigController
}

func GetBackupConfigService() *BackupConfigService {
	return backupConfigService
}
