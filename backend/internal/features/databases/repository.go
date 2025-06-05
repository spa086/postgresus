package databases

import (
	"postgresus-backend/internal/features/databases/databases/postgresql"
	"postgresus-backend/internal/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DatabaseRepository struct{}

func (r *DatabaseRepository) Save(database *Database) error {
	db := storage.GetDb()

	isNew := database.ID == uuid.Nil
	if isNew {
		database.ID = uuid.New()
	}

	database.StorageID = database.Storage.ID

	return db.Transaction(func(tx *gorm.DB) error {
		if database.BackupInterval != nil {
			if database.BackupInterval.ID == uuid.Nil {
				if err := tx.Create(database.BackupInterval).Error; err != nil {
					return err
				}

				database.BackupIntervalID = database.BackupInterval.ID
			} else {
				if err := tx.Save(database.BackupInterval).Error; err != nil {
					return err
				}

				database.BackupIntervalID = database.BackupInterval.ID
			}
		}

		switch database.Type {
		case DatabaseTypePostgres:
			if database.Postgresql != nil {
				database.Postgresql.DatabaseID = &database.ID
			}
		}

		if isNew {
			if err := tx.Create(database).
				Omit("Postgresql", "Storage", "Notifiers", "BackupInterval").
				Error; err != nil {
				return err
			}
		} else {
			if err := tx.Save(database).
				Omit("Postgresql", "Storage", "Notifiers", "BackupInterval").
				Error; err != nil {
				return err
			}
		}

		// Save the specific database type
		switch database.Type {
		case DatabaseTypePostgres:
			if database.Postgresql != nil {
				database.Postgresql.DatabaseID = &database.ID
				if database.Postgresql.ID == uuid.Nil {
					database.Postgresql.ID = uuid.New()
					if err := tx.Create(database.Postgresql).Error; err != nil {
						return err
					}
				} else {
					if err := tx.Save(database.Postgresql).Error; err != nil {
						return err
					}
				}
			}
		}

		if err := tx.Model(database).Association("Notifiers").Replace(database.Notifiers); err != nil {
			return err
		}

		return nil
	})
}

func (r *DatabaseRepository) FindByID(id uuid.UUID) (*Database, error) {
	var database Database

	if err := storage.
		GetDb().
		Preload("BackupInterval").
		Preload("Postgresql").
		Preload("Storage").
		Preload("Notifiers").
		Where("id = ?", id).
		First(&database).Error; err != nil {
		return nil, err
	}

	return &database, nil
}

func (r *DatabaseRepository) FindByUserID(userID uuid.UUID) ([]*Database, error) {
	var databases []*Database

	if err := storage.
		GetDb().
		Preload("BackupInterval").
		Preload("Postgresql").
		Preload("Storage").
		Preload("Notifiers").
		Where("user_id = ?", userID).
		Find(&databases).Error; err != nil {
		return nil, err
	}

	return databases, nil
}

func (r *DatabaseRepository) Delete(id uuid.UUID) error {
	db := storage.GetDb()

	return db.Transaction(func(tx *gorm.DB) error {
		var database Database
		if err := tx.Where("id = ?", id).First(&database).Error; err != nil {
			return err
		}

		if err := tx.Model(&database).Association("Notifiers").Clear(); err != nil {
			return err
		}

		switch database.Type {
		case DatabaseTypePostgres:
			if err := tx.Where("database_id = ?", id).Delete(&postgresql.PostgresqlDatabase{}).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(&Database{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *DatabaseRepository) IsNotifierUsing(notifierID uuid.UUID) (bool, error) {
	var count int64

	if err := storage.
		GetDb().
		Table("database_notifiers").
		Where("notifier_id = ?", notifierID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *DatabaseRepository) IsStorageUsing(storageID uuid.UUID) (bool, error) {
	var count int64

	if err := storage.
		GetDb().
		Table("databases").
		Where("storage_id = ?", storageID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *DatabaseRepository) GetAllDatabases() ([]*Database, error) {
	var databases []*Database

	if err := storage.
		GetDb().
		Preload("BackupInterval").
		Preload("Postgresql").
		Preload("Storage").
		Preload("Notifiers").
		Find(&databases).Error; err != nil {
		return nil, err
	}

	return databases, nil
}
