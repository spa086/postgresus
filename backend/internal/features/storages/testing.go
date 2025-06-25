package storages

import (
	local_storage "postgresus-backend/internal/features/storages/models/local"

	"github.com/google/uuid"
)

func CreateTestStorage(userID uuid.UUID) *Storage {
	storage := &Storage{
		UserID:       userID,
		Type:         StorageTypeLocal,
		Name:         "Test Storage " + uuid.New().String(),
		LocalStorage: &local_storage.LocalStorage{},
	}

	storage, err := storageRepository.Save(storage)
	if err != nil {
		panic(err)
	}

	return storage
}
