package restores

import (
	"postgresus-backend/internal/features/backups"
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
