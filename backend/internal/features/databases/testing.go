package databases

import (
	"postgresus-backend/internal/features/intervals"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"

	"github.com/google/uuid"
)

func CreateTestDatabase(
	userID uuid.UUID,
	storage *storages.Storage,
	notifier *notifiers.Notifier,
) *Database {
	timeOfDay := "16:00"

	database := &Database{
		UserID:      userID,
		Name:        "test " + uuid.New().String(),
		Type:        DatabaseTypePostgres,
		StorePeriod: PeriodDay,

		BackupInterval: &intervals.Interval{
			Interval:  intervals.IntervalDaily,
			TimeOfDay: &timeOfDay,
		},

		StorageID: storage.ID,
		Storage:   *storage,

		Notifiers: []notifiers.Notifier{
			*notifier,
		},
		SendNotificationsOn: []BackupNotificationType{
			NotificationBackupFailed,
			NotificationBackupSuccess,
		},
	}

	database, err := databaseRepository.Save(database)
	if err != nil {
		panic(err)
	}

	return database
}
