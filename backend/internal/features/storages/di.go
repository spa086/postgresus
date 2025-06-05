package storages

import (
	"postgresus-backend/internal/features/users"
)

var storageRepository = &StorageRepository{}
var storageService = &StorageService{
	storageRepository,
}
var storageController = &StorageController{
	storageService,
	users.GetUserService(),
}

func GetStorageService() *StorageService {
	return storageService
}

func GetStorageController() *StorageController {
	return storageController
}
