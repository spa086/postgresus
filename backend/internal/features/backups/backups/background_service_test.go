package backups

import (
	backups_config "postgresus-backend/internal/features/backups/config"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/intervals"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/period"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_MakeBackupForDbHavingBackupDayAgo_BackupCreated(t *testing.T) {
	// setup data
	user := users.GetTestUser()
	storage := storages.CreateTestStorage(user.UserID)
	notifier := notifiers.CreateTestNotifier(user.UserID)
	database := databases.CreateTestDatabase(user.UserID, storage, notifier)

	// Enable backups for the database
	backupConfig, err := backups_config.GetBackupConfigService().GetBackupConfigByDbId(database.ID)
	assert.NoError(t, err)

	timeOfDay := "04:00"
	backupConfig.BackupInterval = &intervals.Interval{
		Interval:  intervals.IntervalDaily,
		TimeOfDay: &timeOfDay,
	}
	backupConfig.IsBackupsEnabled = true
	backupConfig.StorePeriod = period.PeriodWeek
	backupConfig.Storage = storage
	backupConfig.StorageID = &storage.ID

	_, err = backups_config.GetBackupConfigService().SaveBackupConfig(backupConfig)
	assert.NoError(t, err)

	// add old backup
	backupRepository.Save(&Backup{
		Database:   database,
		DatabaseID: database.ID,

		Storage:   storage,
		StorageID: storage.ID,

		Status: BackupStatusCompleted,

		CreatedAt: time.Now().UTC().Add(-24 * time.Hour),
	})

	GetBackupBackgroundService().runPendingBackups()

	time.Sleep(100 * time.Millisecond)

	// assertions
	backups, err := backupRepository.FindByDatabaseID(database.ID)
	assert.NoError(t, err)
	assert.Len(t, backups, 2)

	// cleanup
	for _, backup := range backups {
		err := backupRepository.DeleteByID(backup.ID)
		assert.NoError(t, err)
	}

	databases.RemoveTestDatabase(database)
	storages.RemoveTestStorage(storage.ID)
	notifiers.RemoveTestNotifier(notifier)
}

func Test_MakeBackupForDbHavingHourAgoBackup_BackupSkipped(t *testing.T) {
	// setup data
	user := users.GetTestUser()
	storage := storages.CreateTestStorage(user.UserID)
	notifier := notifiers.CreateTestNotifier(user.UserID)
	database := databases.CreateTestDatabase(user.UserID, storage, notifier)

	// Enable backups for the database
	backupConfig, err := backups_config.GetBackupConfigService().GetBackupConfigByDbId(database.ID)
	assert.NoError(t, err)

	timeOfDay := "04:00"
	backupConfig.BackupInterval = &intervals.Interval{
		Interval:  intervals.IntervalDaily,
		TimeOfDay: &timeOfDay,
	}
	backupConfig.IsBackupsEnabled = true
	backupConfig.StorePeriod = period.PeriodWeek
	backupConfig.Storage = storage
	backupConfig.StorageID = &storage.ID

	_, err = backups_config.GetBackupConfigService().SaveBackupConfig(backupConfig)
	assert.NoError(t, err)

	// add recent backup (1 hour ago)
	backupRepository.Save(&Backup{
		Database:   database,
		DatabaseID: database.ID,

		Storage:   storage,
		StorageID: storage.ID,

		Status: BackupStatusCompleted,

		CreatedAt: time.Now().UTC().Add(-1 * time.Hour),
	})

	GetBackupBackgroundService().runPendingBackups()

	time.Sleep(100 * time.Millisecond)

	// assertions
	backups, err := backupRepository.FindByDatabaseID(database.ID)
	assert.NoError(t, err)
	assert.Len(t, backups, 1) // Should still be 1 backup, no new backup created

	// cleanup
	for _, backup := range backups {
		err := backupRepository.DeleteByID(backup.ID)
		assert.NoError(t, err)
	}

	databases.RemoveTestDatabase(database)
	storages.RemoveTestStorage(storage.ID)
	notifiers.RemoveTestNotifier(notifier)
}

func Test_MakeBackupHavingFailedBackupWithoutRetries_BackupSkipped(t *testing.T) {
	// setup data
	user := users.GetTestUser()
	storage := storages.CreateTestStorage(user.UserID)
	notifier := notifiers.CreateTestNotifier(user.UserID)
	database := databases.CreateTestDatabase(user.UserID, storage, notifier)

	// Enable backups for the database with retries disabled
	backupConfig, err := backups_config.GetBackupConfigService().GetBackupConfigByDbId(database.ID)
	assert.NoError(t, err)

	timeOfDay := "04:00"
	backupConfig.BackupInterval = &intervals.Interval{
		Interval:  intervals.IntervalDaily,
		TimeOfDay: &timeOfDay,
	}
	backupConfig.IsBackupsEnabled = true
	backupConfig.StorePeriod = period.PeriodWeek
	backupConfig.Storage = storage
	backupConfig.StorageID = &storage.ID
	backupConfig.IsRetryIfFailed = false
	backupConfig.MaxFailedTriesCount = 0

	_, err = backups_config.GetBackupConfigService().SaveBackupConfig(backupConfig)
	assert.NoError(t, err)

	// add failed backup
	failMessage := "backup failed"
	backupRepository.Save(&Backup{
		Database:   database,
		DatabaseID: database.ID,

		Storage:   storage,
		StorageID: storage.ID,

		Status:      BackupStatusFailed,
		FailMessage: &failMessage,

		CreatedAt: time.Now().UTC().Add(-1 * time.Hour),
	})

	GetBackupBackgroundService().runPendingBackups()

	time.Sleep(100 * time.Millisecond)

	// assertions
	backups, err := backupRepository.FindByDatabaseID(database.ID)
	assert.NoError(t, err)
	assert.Len(t, backups, 1) // Should still be 1 backup, no retry attempted

	// cleanup
	for _, backup := range backups {
		err := backupRepository.DeleteByID(backup.ID)
		assert.NoError(t, err)
	}

	databases.RemoveTestDatabase(database)
	storages.RemoveTestStorage(storage.ID)
	notifiers.RemoveTestNotifier(notifier)
}

func Test_MakeBackupHavingFailedBackupWithRetries_BackupCreated(t *testing.T) {
	// setup data
	user := users.GetTestUser()
	storage := storages.CreateTestStorage(user.UserID)
	notifier := notifiers.CreateTestNotifier(user.UserID)
	database := databases.CreateTestDatabase(user.UserID, storage, notifier)

	// Enable backups for the database with retries enabled
	backupConfig, err := backups_config.GetBackupConfigService().GetBackupConfigByDbId(database.ID)
	assert.NoError(t, err)

	timeOfDay := "04:00"
	backupConfig.BackupInterval = &intervals.Interval{
		Interval:  intervals.IntervalDaily,
		TimeOfDay: &timeOfDay,
	}
	backupConfig.IsBackupsEnabled = true
	backupConfig.StorePeriod = period.PeriodWeek
	backupConfig.Storage = storage
	backupConfig.StorageID = &storage.ID
	backupConfig.IsRetryIfFailed = true
	backupConfig.MaxFailedTriesCount = 3

	_, err = backups_config.GetBackupConfigService().SaveBackupConfig(backupConfig)
	assert.NoError(t, err)

	// add failed backup
	failMessage := "backup failed"
	backupRepository.Save(&Backup{
		Database:   database,
		DatabaseID: database.ID,

		Storage:   storage,
		StorageID: storage.ID,

		Status:      BackupStatusFailed,
		FailMessage: &failMessage,

		CreatedAt: time.Now().UTC().Add(-1 * time.Hour),
	})

	GetBackupBackgroundService().runPendingBackups()

	time.Sleep(100 * time.Millisecond)

	// assertions
	backups, err := backupRepository.FindByDatabaseID(database.ID)
	assert.NoError(t, err)
	assert.Len(t, backups, 2) // Should have 2 backups, retry was attempted

	// cleanup
	for _, backup := range backups {
		err := backupRepository.DeleteByID(backup.ID)
		assert.NoError(t, err)
	}

	databases.RemoveTestDatabase(database)
	storages.RemoveTestStorage(storage.ID)
	notifiers.RemoveTestNotifier(notifier)
}

func Test_MakeBackupHavingFailedBackupWithRetries_RetriesCountNotExceeded(t *testing.T) {
	// setup data
	user := users.GetTestUser()
	storage := storages.CreateTestStorage(user.UserID)
	notifier := notifiers.CreateTestNotifier(user.UserID)
	database := databases.CreateTestDatabase(user.UserID, storage, notifier)

	// Enable backups for the database with retries enabled
	backupConfig, err := backups_config.GetBackupConfigService().GetBackupConfigByDbId(database.ID)
	assert.NoError(t, err)

	timeOfDay := "04:00"
	backupConfig.BackupInterval = &intervals.Interval{
		Interval:  intervals.IntervalDaily,
		TimeOfDay: &timeOfDay,
	}
	backupConfig.IsBackupsEnabled = true
	backupConfig.StorePeriod = period.PeriodWeek
	backupConfig.Storage = storage
	backupConfig.StorageID = &storage.ID
	backupConfig.IsRetryIfFailed = true
	backupConfig.MaxFailedTriesCount = 3

	_, err = backups_config.GetBackupConfigService().SaveBackupConfig(backupConfig)
	assert.NoError(t, err)

	failMessage := "backup failed"

	for i := 0; i < 3; i++ {
		backupRepository.Save(&Backup{
			Database:   database,
			DatabaseID: database.ID,

			Storage:   storage,
			StorageID: storage.ID,

			Status:      BackupStatusFailed,
			FailMessage: &failMessage,

			CreatedAt: time.Now().UTC().Add(-1 * time.Hour),
		})
	}

	GetBackupBackgroundService().runPendingBackups()

	time.Sleep(100 * time.Millisecond)

	// assertions
	backups, err := backupRepository.FindByDatabaseID(database.ID)
	assert.NoError(t, err)
	assert.Len(t, backups, 3) // Should have 3 backups, not more than max

	// cleanup
	for _, backup := range backups {
		err := backupRepository.DeleteByID(backup.ID)
		assert.NoError(t, err)
	}

	databases.RemoveTestDatabase(database)
	storages.RemoveTestStorage(storage.ID)
	notifiers.RemoveTestNotifier(notifier)
}
