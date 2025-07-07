package backups

import (
	"errors"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"
	"postgresus-backend/internal/features/users"
	"postgresus-backend/internal/util/logger"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_BackupExecuted_NotificationSent(t *testing.T) {
	user := users.GetTestUser()
	storage := storages.CreateTestStorage(user.UserID)
	notifier := notifiers.CreateTestNotifier(user.UserID)
	database := databases.CreateTestDatabase(user.UserID, storage, notifier)

	defer storages.RemoveTestStorage(storage.ID)
	defer notifiers.RemoveTestNotifier(notifier)
	defer databases.RemoveTestDatabase(database)

	t.Run("BackupFailed_FailNotificationSent", func(t *testing.T) {
		mockNotificationSender := &MockNotificationSender{}
		backupService := &BackupService{
			databases.GetDatabaseService(),
			storages.GetStorageService(),
			backupRepository,
			notifiers.GetNotifierService(),
			mockNotificationSender,
			&CreateFailedBackupUsecase{},
			logger.GetLogger(),
		}

		// Set up expectations
		mockNotificationSender.On("SendNotification",
			mock.Anything,
			mock.MatchedBy(func(title string) bool {
				return strings.Contains(title, "❌ Backup failed")
			}),
			mock.MatchedBy(func(message string) bool {
				return strings.Contains(message, "backup failed")
			}),
		).Once()

		backupService.MakeBackup(database.ID)

		// Verify all expectations were met
		mockNotificationSender.AssertExpectations(t)
	})

	t.Run("BackupSuccess_SuccessNotificationSent", func(t *testing.T) {
		mockNotificationSender := &MockNotificationSender{}

		// Set up expectations
		mockNotificationSender.On("SendNotification",
			mock.Anything,
			mock.MatchedBy(func(title string) bool {
				return strings.Contains(title, "✅ Backup completed")
			}),
			mock.MatchedBy(func(message string) bool {
				return strings.Contains(message, "Backup completed successfully")
			}),
		).Once()

		backupService := &BackupService{
			databases.GetDatabaseService(),
			storages.GetStorageService(),
			backupRepository,
			notifiers.GetNotifierService(),
			mockNotificationSender,
			&CreateSuccessBackupUsecase{},
			logger.GetLogger(),
		}

		backupService.MakeBackup(database.ID)

		// Verify all expectations were met
		mockNotificationSender.AssertExpectations(t)
	})

	t.Run("BackupSuccess_VerifyNotificationContent", func(t *testing.T) {
		mockNotificationSender := &MockNotificationSender{}
		backupService := &BackupService{
			databases.GetDatabaseService(),
			storages.GetStorageService(),
			backupRepository,
			notifiers.GetNotifierService(),
			mockNotificationSender,
			&CreateSuccessBackupUsecase{},
			logger.GetLogger(),
		}

		// capture arguments
		var capturedNotifier *notifiers.Notifier
		var capturedTitle string
		var capturedMessage string

		mockNotificationSender.On("SendNotification",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Run(func(args mock.Arguments) {
			capturedNotifier = args.Get(0).(*notifiers.Notifier)
			capturedTitle = args.Get(1).(string)
			capturedMessage = args.Get(2).(string)
		}).Once()

		backupService.MakeBackup(database.ID)

		// Verify expectations were met
		mockNotificationSender.AssertExpectations(t)

		// Additional detailed assertions
		assert.Contains(t, capturedTitle, "✅ Backup completed")
		assert.Contains(t, capturedTitle, database.Name)
		assert.Contains(t, capturedMessage, "Backup completed successfully")
		assert.Contains(t, capturedMessage, "10.00 MB")
		assert.Equal(t, notifier.ID, capturedNotifier.ID)
	})
}

type CreateFailedBackupUsecase struct {
}

func (uc *CreateFailedBackupUsecase) Execute(
	backupID uuid.UUID,
	database *databases.Database,
	storage *storages.Storage,
	backupProgressListener func(
		completedMBs float64,
	),
) error {
	backupProgressListener(10) // Assume we completed 10MB
	return errors.New("backup failed")
}

type CreateSuccessBackupUsecase struct {
}

func (uc *CreateSuccessBackupUsecase) Execute(
	backupID uuid.UUID,
	database *databases.Database,
	storage *storages.Storage,
	backupProgressListener func(
		completedMBs float64,
	),
) error {
	backupProgressListener(10) // Assume we completed 10MB
	return nil
}
