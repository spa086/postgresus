package restores

import (
	"postgresus-backend/internal/features/backups/backups"
	backups_config "postgresus-backend/internal/features/backups/config"
	"postgresus-backend/internal/features/restores/usecases"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/logger"
)

var restoreRepository = &RestoreRepository{}
var restoreService = &RestoreService{
	backups.GetBackupService(),
	restoreRepository,
	storages.GetStorageService(),
	backups_config.GetBackupConfigService(),
	usecases.GetRestoreBackupUsecase(),
	logger.GetLogger(),
}
var restoreController = &RestoreController{
	restoreService,
	users.GetUserService(),
}

var restoreBackgroundService = &RestoreBackgroundService{
	restoreRepository,
	logger.GetLogger(),
}

func GetRestoreController() *RestoreController {
	return restoreController
}

func GetRestoreBackgroundService() *RestoreBackgroundService {
	return restoreBackgroundService
}

func SetupDependencies() {
	backups.GetBackupService().AddBackupRemoveListener(restoreService)
}
