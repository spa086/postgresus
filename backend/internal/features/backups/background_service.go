package backups

import (
	"log/slog"
	"postgresus-backend/internal/config"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/storages"
	"time"
)

type BackupBackgroundService struct {
	backupService    *BackupService
	backupRepository *BackupRepository
	databaseService  *databases.DatabaseService
	storageService   *storages.StorageService

	lastBackupTime time.Time
	logger         *slog.Logger
}

func (s *BackupBackgroundService) Run() {
	s.lastBackupTime = time.Now().UTC()

	if err := s.failBackupsInProgress(); err != nil {
		s.logger.Error("Failed to fail backups in progress", "error", err)
		panic(err)
	}

	if config.IsShouldShutdown() {
		return
	}

	for {
		if config.IsShouldShutdown() {
			return
		}

		if err := s.cleanOldBackups(); err != nil {
			s.logger.Error("Failed to clean old backups", "error", err)
		}

		if err := s.runPendingBackups(); err != nil {
			s.logger.Error("Failed to run pending backups", "error", err)
		}

		s.lastBackupTime = time.Now().UTC()
		time.Sleep(1 * time.Minute)
	}
}

func (s *BackupBackgroundService) IsBackupsRunning() bool {
	// if last backup time is more than 5 minutes ago, return false
	return s.lastBackupTime.After(time.Now().UTC().Add(-5 * time.Minute))
}

func (s *BackupBackgroundService) failBackupsInProgress() error {
	backupsInProgress, err := s.backupRepository.FindByStatus(BackupStatusInProgress)
	if err != nil {
		return err
	}

	for _, backup := range backupsInProgress {
		failMessage := "Backup failed due to application restart"
		backup.FailMessage = &failMessage
		backup.Status = BackupStatusFailed
		backup.BackupSizeMb = 0

		s.backupService.SendBackupNotification(
			backup.Database,
			backup,
			databases.NotificationBackupFailed,
			&failMessage,
		)

		if err := s.backupRepository.Save(backup); err != nil {
			return err
		}
	}

	return nil
}

func (s *BackupBackgroundService) cleanOldBackups() error {
	allDatabases, err := s.databaseService.GetAllDatabases()
	if err != nil {
		return err
	}

	for _, database := range allDatabases {
		backupStorePeriod := database.StorePeriod

		if backupStorePeriod == databases.PeriodForever {
			continue
		}

		storeDuration := backupStorePeriod.ToDuration()
		dateBeforeBackupsShouldBeDeleted := time.Now().UTC().Add(-storeDuration)

		oldBackups, err := s.backupRepository.FindBackupsBeforeDate(
			database.ID,
			dateBeforeBackupsShouldBeDeleted,
		)
		if err != nil {
			s.logger.Error(
				"Failed to find old backups for database",
				"databaseId",
				database.ID,
				"error",
				err,
			)
			continue
		}

		for _, backup := range oldBackups {
			storage, err := s.storageService.GetStorageByID(backup.StorageID)
			if err != nil {
				s.logger.Error(
					"Failed to get storage by ID",
					"storageId",
					backup.StorageID,
					"error",
					err,
				)
				continue
			}

			err = storage.DeleteFile(backup.ID)
			if err != nil {
				s.logger.Error("Failed to delete backup file", "backupId", backup.ID, "error", err)
			}

			if err := s.backupRepository.DeleteByID(backup.ID); err != nil {
				s.logger.Error("Failed to delete old backup", "backupId", backup.ID, "error", err)
				continue
			}

			s.logger.Info("Deleted old backup", "backupId", backup.ID, "databaseId", database.ID)
		}
	}

	return nil
}

func (s *BackupBackgroundService) runPendingBackups() error {
	allDatabases, err := s.databaseService.GetAllDatabases()
	if err != nil {
		return err
	}

	for _, database := range allDatabases {
		if database.BackupInterval == nil {
			continue
		}

		lastBackup, err := s.backupRepository.FindLastByDatabaseID(database.ID)
		if err != nil {
			s.logger.Error(
				"Failed to get last backup for database",
				"databaseId",
				database.ID,
				"error",
				err,
			)
			continue
		}

		var lastBackupTime *time.Time
		if lastBackup != nil {
			lastBackupTime = &lastBackup.CreatedAt
		}

		if database.BackupInterval.ShouldTriggerBackup(time.Now().UTC(), lastBackupTime) {
			s.logger.Info(
				"Triggering scheduled backup",
				"databaseId",
				database.ID,
				"intervalType",
				database.BackupInterval.Interval,
			)

			go s.backupService.MakeBackup(database.ID)
			s.logger.Info("Successfully triggered scheduled backup", "databaseId", database.ID)
		}
	}

	return nil
}
