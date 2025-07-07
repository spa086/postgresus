package databases

import (
	"postgresus-backend/internal/features/databases/databases/postgresql"
	"postgresus-backend/internal/features/intervals"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/util/tools"

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

		Postgresql: &postgresql.PostgresqlDatabase{
			Version:  tools.PostgresqlVersion16,
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
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

func RemoveTestDatabase(database *Database) {
	err := databaseRepository.Delete(database.ID)
	if err != nil {
		panic(err)
	}
}
