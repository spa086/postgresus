package storages

import (
	"errors"
	users_models "postgresus-backend/internal/features/users/models"

	"github.com/google/uuid"
)

type StorageService struct {
	storageRepository *StorageRepository
}

func (s *StorageService) SaveStorage(
	user *users_models.User,
	storage *Storage,
) error {
	if storage.ID != uuid.Nil {
		existingStorage, err := s.storageRepository.FindByID(storage.ID)
		if err != nil {
			return err
		}

		if existingStorage.UserID != user.ID {
			return errors.New("you have not access to this storage")
		}

		storage.UserID = existingStorage.UserID
	} else {
		storage.UserID = user.ID
	}

	_, err := s.storageRepository.Save(storage)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) DeleteStorage(
	user *users_models.User,
	storageID uuid.UUID,
) error {
	storage, err := s.storageRepository.FindByID(storageID)
	if err != nil {
		return err
	}

	if storage.UserID != user.ID {
		return errors.New("you have not access to this storage")
	}

	return s.storageRepository.Delete(storage)
}

func (s *StorageService) GetStorage(
	user *users_models.User,
	id uuid.UUID,
) (*Storage, error) {
	storage, err := s.storageRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if storage.UserID != user.ID {
		return nil, errors.New("you have not access to this storage")
	}

	return storage, nil
}

func (s *StorageService) GetStorages(
	user *users_models.User,
) ([]*Storage, error) {
	return s.storageRepository.FindByUserID(user.ID)
}

func (s *StorageService) TestStorageConnection(
	user *users_models.User,
	storageID uuid.UUID,
) error {
	storage, err := s.storageRepository.FindByID(storageID)
	if err != nil {
		return err
	}

	if storage.UserID != user.ID {
		return errors.New("you have not access to this storage")
	}

	err = storage.TestConnection()
	if err != nil {
		lastSaveError := err.Error()
		storage.LastSaveError = &lastSaveError
		return err
	}

	storage.LastSaveError = nil
	_, err = s.storageRepository.Save(storage)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) TestStorageConnectionDirect(
	storage *Storage,
) error {
	return storage.TestConnection()
}

func (s *StorageService) GetStorageByID(
	id uuid.UUID,
) (*Storage, error) {
	return s.storageRepository.FindByID(id)
}
