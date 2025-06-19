package backups

import (
	"postgresus-backend/internal/features/backups/usecases"
	usecases_postgresql "postgresus-backend/internal/features/backups/usecases/postgresql"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/features/users"
	"time"
)

var createPostgresqlBackupUsecase = &usecases_postgresql.CreatePostgresqlBackupUsecase{}
var createBackupUseCase = &usecases.CreateBackupUsecase{
	CreatePostgresqlBackupUsecase: createPostgresqlBackupUsecase,
}
var backupRepository = &BackupRepository{}
var backupService = &BackupService{
	databases.GetDatabaseService(),
	storages.GetStorageService(),
	backupRepository,
	notifiers.GetNotifierService(),
	createBackupUseCase,
}

var backupBackgroundService = &BackupBackgroundService{
	backupService,
	backupRepository,
	databases.GetDatabaseService(),
	storages.GetStorageService(),
	time.Now().UTC(),
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
