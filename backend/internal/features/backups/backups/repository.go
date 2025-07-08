package backups

import (
	"postgresus-backend/internal/storage"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackupRepository struct{}

func (r *BackupRepository) Save(backup *Backup) error {
	db := storage.GetDb()

	isNew := backup.ID == uuid.Nil
	if isNew {
		backup.ID = uuid.New()
		return db.Create(backup).
			Omit("Database", "Storage").
			Error
	}

	return db.Save(backup).
		Omit("Database", "Storage").
		Error
}

func (r *BackupRepository) FindByDatabaseID(databaseID uuid.UUID) ([]*Backup, error) {
	var backups []*Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("database_id = ?", databaseID).
		Order("created_at DESC").
		Find(&backups).Error; err != nil {
		return nil, err
	}

	return backups, nil
}

func (r *BackupRepository) FindByStorageID(storageID uuid.UUID) ([]*Backup, error) {
	var backups []*Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("storage_id = ?", storageID).
		Order("created_at DESC").
		Find(&backups).Error; err != nil {
		return nil, err
	}

	return backups, nil
}

func (r *BackupRepository) FindLastByDatabaseID(databaseID uuid.UUID) (*Backup, error) {
	var backup Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("database_id = ?", databaseID).
		Order("created_at DESC").
		First(&backup).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &backup, nil
}

func (r *BackupRepository) FindByID(id uuid.UUID) (*Backup, error) {
	var backup Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("id = ?", id).
		First(&backup).Error; err != nil {
		return nil, err
	}

	return &backup, nil
}

func (r *BackupRepository) FindByStatus(status BackupStatus) ([]*Backup, error) {
	var backups []*Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&backups).Error; err != nil {
		return nil, err
	}

	return backups, nil
}

func (r *BackupRepository) FindByStorageIdAndStatus(
	storageID uuid.UUID,
	status BackupStatus,
) ([]*Backup, error) {
	var backups []*Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("storage_id = ? AND status = ?", storageID, status).
		Order("created_at DESC").
		Find(&backups).Error; err != nil {
		return nil, err
	}

	return backups, nil
}

func (r *BackupRepository) FindByDatabaseIdAndStatus(
	databaseID uuid.UUID,
	status BackupStatus,
) ([]*Backup, error) {
	var backups []*Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("database_id = ? AND status = ?", databaseID, status).
		Order("created_at DESC").
		Find(&backups).Error; err != nil {
		return nil, err
	}

	return backups, nil
}

func (r *BackupRepository) DeleteByID(id uuid.UUID) error {
	return storage.GetDb().Delete(&Backup{}, "id = ?", id).Error
}

func (r *BackupRepository) FindBackupsBeforeDate(
	databaseID uuid.UUID,
	date time.Time,
) ([]*Backup, error) {
	var backups []*Backup

	if err := storage.
		GetDb().
		Preload("Database").
		Preload("Storage").
		Where("database_id = ? AND created_at < ?", databaseID, date).
		Order("created_at DESC").
		Find(&backups).Error; err != nil {
		return nil, err
	}

	return backups, nil
}
