package backups

import (
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/storages"
	"time"

	"github.com/google/uuid"
)

type Backup struct {
	ID uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey"`

	Database   *databases.Database `json:"database"   gorm:"foreignKey:DatabaseID"`
	DatabaseID uuid.UUID           `json:"databaseId" gorm:"column:database_id;type:uuid;not null"`

	Storage   *storages.Storage `json:"storage"   gorm:"foreignKey:StorageID"`
	StorageID uuid.UUID         `json:"storageId" gorm:"column:storage_id;type:uuid;not null"`

	Status      BackupStatus `json:"status"      gorm:"column:status;not null"`
	FailMessage *string      `json:"failMessage" gorm:"column:fail_message"`

	BackupSizeMb float64 `json:"backupSizeMb" gorm:"column:backup_size_mb;default:0"`

	BackupDurationMs int64 `json:"backupDurationMs" gorm:"column:backup_duration_ms;default:0"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (b *Backup) DeleteBackupFromStorage() {
	if b.Status != BackupStatusCompleted {
		return
	}

	err := b.Storage.DeleteFile(b.ID)
	if err != nil {
		log.Error("Failed to delete backup from storage", "error", err)
		// we ignore the error, because the access to the storage
		// may be lost, file already deleted, etc.
	}
}
