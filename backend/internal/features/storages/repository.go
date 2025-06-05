package storages

import (
	db "postgresus-backend/internal/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorageRepository struct{}

func (r *StorageRepository) Save(s *Storage) error {
	database := db.GetDb()

	return database.Transaction(func(tx *gorm.DB) error {
		switch s.Type {
		case StorageTypeLocal:
			if s.LocalStorage != nil {
				s.LocalStorage.StorageID = s.ID
			}
		case StorageTypeS3:
			if s.S3Storage != nil {
				s.S3Storage.StorageID = s.ID
			}
		}

		if s.ID == uuid.Nil {
			if err := tx.Create(s).
				Omit("LocalStorage", "S3Storage").
				Error; err != nil {
				return err
			}
		} else {
			if err := tx.Save(s).
				Omit("LocalStorage", "S3Storage").
				Error; err != nil {
				return err
			}
		}

		switch s.Type {
		case StorageTypeLocal:
			if s.LocalStorage != nil {
				s.LocalStorage.StorageID = s.ID // Ensure ID is set
				if err := tx.Save(s.LocalStorage).Error; err != nil {
					return err
				}
			}
		case StorageTypeS3:
			if s.S3Storage != nil {
				s.S3Storage.StorageID = s.ID // Ensure ID is set
				if err := tx.Save(s.S3Storage).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *StorageRepository) FindByID(id uuid.UUID) (*Storage, error) {
	var s Storage

	if err := db.
		GetDb().
		Preload("LocalStorage").
		Preload("S3Storage").
		Where("id = ?", id).
		First(&s).Error; err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StorageRepository) FindByUserID(userID uuid.UUID) ([]*Storage, error) {
	var storages []*Storage

	if err := db.
		GetDb().
		Preload("LocalStorage").
		Preload("S3Storage").
		Where("user_id = ?", userID).
		Find(&storages).Error; err != nil {
		return nil, err
	}

	return storages, nil
}

func (r *StorageRepository) Delete(s *Storage) error {
	return db.GetDb().Transaction(func(tx *gorm.DB) error {
		// Delete specific storage based on type
		switch s.Type {
		case StorageTypeLocal:
			if s.LocalStorage != nil {
				if err := tx.Delete(s.LocalStorage).Error; err != nil {
					return err
				}
			}
		case StorageTypeS3:
			if s.S3Storage != nil {
				if err := tx.Delete(s.S3Storage).Error; err != nil {
					return err
				}
			}
		}

		// Delete the main storage
		return tx.Delete(s).Error
	})
}
