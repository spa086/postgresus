package backups

import (
	"errors"
	"fmt"
	"log/slog"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/notifiers"
	"postgresus-backend/internal/features/storages"
	users_models "postgresus-backend/internal/features/users/models"
	"slices"
	"time"

	"github.com/google/uuid"
)

type BackupService struct {
	databaseService    *databases.DatabaseService
	storageService     *storages.StorageService
	backupRepository   *BackupRepository
	notifierService    *notifiers.NotifierService
	notificationSender NotificationSender

	createBackupUseCase CreateBackupUsecase

	logger *slog.Logger
}

func (s *BackupService) OnBeforeDbStorageChange(
	databaseID uuid.UUID,
	storageID uuid.UUID,
) error {
	// validate no backups in progress
	backups, err := s.backupRepository.FindByStorageIdAndStatus(
		storageID,
		BackupStatusInProgress,
	)
	if err != nil {
		return err
	}

	if len(backups) > 0 {
		return errors.New("backup is in progress, storage cannot")
	}

	backupsWithStorage, err := s.backupRepository.FindByStorageIdAndStatus(
		storageID,
		BackupStatusCompleted,
	)
	if err != nil {
		return err
	}

	if len(backupsWithStorage) > 0 {
		for _, backup := range backupsWithStorage {
			if err := backup.Storage.DeleteFile(backup.ID); err != nil {
				// most likely we cannot do nothing with this,
				// so we just remove the backup model
				s.logger.Error("Failed to delete backup file", "error", err)
			}

			if err := s.backupRepository.DeleteByID(backup.ID); err != nil {
				return err
			}
		}

		// we repeat remove for the case if backup
		// started until we removed all previous backups
		return s.OnBeforeDbStorageChange(databaseID, storageID)
	}

	return nil
}

func (s *BackupService) MakeBackupWithAuth(
	user *users_models.User,
	databaseID uuid.UUID,
) error {
	database, err := s.databaseService.GetDatabaseByID(databaseID)
	if err != nil {
		return err
	}

	if database.UserID != user.ID {
		return errors.New("user does not have access to this database")
	}

	go s.MakeBackup(databaseID)

	return nil
}

func (s *BackupService) GetBackups(
	user *users_models.User,
	databaseID uuid.UUID,
) ([]*Backup, error) {
	database, err := s.databaseService.GetDatabaseByID(databaseID)
	if err != nil {
		return nil, err
	}

	if database.UserID != user.ID {
		return nil, errors.New("user does not have access to this database")
	}

	backups, err := s.backupRepository.FindByDatabaseID(databaseID)
	if err != nil {
		return nil, err
	}

	return backups, nil
}

func (s *BackupService) DeleteBackup(
	user *users_models.User,
	backupID uuid.UUID,
) error {
	backup, err := s.backupRepository.FindByID(backupID)
	if err != nil {
		return err
	}

	if backup.Database.UserID != user.ID {
		return errors.New("user does not have access to this backup")
	}

	if backup.Status == BackupStatusInProgress {
		return errors.New("backup is in progress")
	}

	backup.DeleteBackupFromStorage(s.logger)

	backup.Status = BackupStatusDeleted
	return s.backupRepository.Save(backup)
}

func (s *BackupService) MakeBackup(databaseID uuid.UUID) {
	database, err := s.databaseService.GetDatabaseByID(databaseID)
	if err != nil {
		s.logger.Error("Failed to get database by ID", "error", err)
		return
	}

	lastBackup, err := s.backupRepository.FindLastByDatabaseID(databaseID)
	if err != nil {
		s.logger.Error("Failed to find last backup by database ID", "error", err)
		return
	}

	if lastBackup != nil && lastBackup.Status == BackupStatusInProgress {
		s.logger.Error("Backup is in progress")
		return
	}

	storage, err := s.storageService.GetStorageByID(database.StorageID)
	if err != nil {
		s.logger.Error("Failed to get storage by ID", "error", err)
		return
	}

	backup := &Backup{
		DatabaseID: databaseID,
		Database:   database,

		StorageID: storage.ID,
		Storage:   storage,

		Status: BackupStatusInProgress,

		BackupSizeMb: 0,

		CreatedAt: time.Now().UTC(),
	}

	if err := s.backupRepository.Save(backup); err != nil {
		s.logger.Error("Failed to save backup", "error", err)
		return
	}

	start := time.Now().UTC()

	backupProgressListener := func(
		completedMBs float64,
	) {
		backup.BackupSizeMb = completedMBs
		backup.BackupDurationMs = time.Since(start).Milliseconds()

		if err := s.backupRepository.Save(backup); err != nil {
			s.logger.Error("Failed to update backup progress", "error", err)
		}
	}

	err = s.createBackupUseCase.Execute(
		backup.ID,
		database,
		storage,
		backupProgressListener,
	)
	if err != nil {
		errMsg := err.Error()
		backup.FailMessage = &errMsg
		backup.Status = BackupStatusFailed
		backup.BackupDurationMs = time.Since(start).Milliseconds()
		backup.BackupSizeMb = 0

		if updateErr := s.databaseService.SetBackupError(databaseID, errMsg); updateErr != nil {
			s.logger.Error(
				"Failed to update database last backup time",
				"databaseId",
				databaseID,
				"error",
				updateErr,
			)
		}

		if err := s.backupRepository.Save(backup); err != nil {
			s.logger.Error("Failed to save backup", "error", err)
		}

		s.SendBackupNotification(
			database,
			backup,
			databases.NotificationBackupFailed,
			&errMsg,
		)

		return
	}

	backup.Status = BackupStatusCompleted
	backup.BackupDurationMs = time.Since(start).Milliseconds()

	if err := s.backupRepository.Save(backup); err != nil {
		s.logger.Error("Failed to save backup", "error", err)
		return
	}

	// Update database last backup time
	now := time.Now().UTC()
	if updateErr := s.databaseService.SetLastBackupTime(databaseID, now); updateErr != nil {
		s.logger.Error(
			"Failed to update database last backup time",
			"databaseId",
			databaseID,
			"error",
			updateErr,
		)
	}

	s.SendBackupNotification(
		database,
		backup,
		databases.NotificationBackupSuccess,
		nil,
	)
}

func (s *BackupService) SendBackupNotification(
	db *databases.Database,
	backup *Backup,
	notificationType databases.BackupNotificationType,
	errorMessage *string,
) {
	database, err := s.databaseService.GetDatabaseByID(db.ID)
	if err != nil {
		return
	}

	for _, notifier := range database.Notifiers {
		if !slices.Contains(
			database.SendNotificationsOn,
			notificationType,
		) {
			continue
		}

		title := ""
		switch notificationType {
		case databases.NotificationBackupFailed:
			title = fmt.Sprintf("❌ Backup failed for database \"%s\"", database.Name)
		case databases.NotificationBackupSuccess:
			title = fmt.Sprintf("✅ Backup completed for database \"%s\"", database.Name)
		}

		message := ""
		if errorMessage != nil {
			message = *errorMessage
		} else {
			// Format size conditionally
			var sizeStr string
			if backup.BackupSizeMb < 1024 {
				sizeStr = fmt.Sprintf("%.2f MB", backup.BackupSizeMb)
			} else {
				sizeGB := backup.BackupSizeMb / 1024
				sizeStr = fmt.Sprintf("%.2f GB", sizeGB)
			}

			// Format duration as "0m 0s 0ms"
			totalMs := backup.BackupDurationMs
			minutes := totalMs / (1000 * 60)
			seconds := (totalMs % (1000 * 60)) / 1000
			durationStr := fmt.Sprintf("%dm %ds", minutes, seconds)

			message = fmt.Sprintf(
				"Backup completed successfully in %s.\nCompressed backup size: %s",
				durationStr,
				sizeStr,
			)
		}

		s.notificationSender.SendNotification(
			&notifier,
			title,
			message,
		)
	}
}

func (s *BackupService) GetBackup(backupID uuid.UUID) (*Backup, error) {
	return s.backupRepository.FindByID(backupID)
}
