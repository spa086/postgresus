package backups

import (
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"

	"github.com/google/uuid"
)

type NotificationSender interface {
	SendNotification(
		notifier *notifiers.Notifier,
		title string,
		message string,
	)
}

type CreateBackupUsecase interface {
	Execute(
		backupID uuid.UUID,
		database *databases.Database,
		storage *storages.Storage,
		backupProgressListener func(
			completedMBs float64,
		),
	) error
}
