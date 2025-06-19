package restores

import (
	"postgresus-backend/internal/features/backups"
	"postgresus-backend/internal/features/restores/usecases"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/features/users"
)

var restoreBackupUsecase = &usecases.RestoreBackupUsecase{}
var restoreRepository = &RestoreRepository{}
var restoreService = &RestoreService{
	backups.GetBackupService(),
	restoreRepository,
	storages.GetStorageService(),
	restoreBackupUsecase,
}
var restoreController = &RestoreController{
	restoreService,
	users.GetUserService(),
}

var restoreBackgroundService = &RestoreBackgroundService{
	restoreRepository,
}

func GetRestoreController() *RestoreController {
	return restoreController
}

func GetRestoreBackgroundService() *RestoreBackgroundService {
	return restoreBackgroundService
}
