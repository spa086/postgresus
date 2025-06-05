package restores

import (
	"errors"
	"postgresus-backend/internal/features/backups"
	"postgresus-backend/internal/features/databases"
	"postgresus-backend/internal/features/restores/enums"
	"postgresus-backend/internal/features/restores/models"
	"postgresus-backend/internal/features/restores/usecases"
	users_models "postgresus-backend/internal/features/users/models"
	"time"

	"github.com/google/uuid"
)

type RestoreService struct {
	backupService     *backups.BackupService
	restoreRepository *RestoreRepository

	restoreBackupUsecase *usecases.RestoreBackupUsecase
}

func (s *RestoreService) GetRestores(
	user *users_models.User,
	backupID uuid.UUID,
) ([]*models.Restore, error) {
	backup, err := s.backupService.GetBackup(backupID)
	if err != nil {
		return nil, err
	}

	if backup.Database.UserID != user.ID {
		return nil, errors.New("user does not have access to this backup")
	}

	return s.restoreRepository.FindByBackupID(backupID)
}

func (s *RestoreService) RestoreBackupWithAuth(
	user *users_models.User,
	backupID uuid.UUID,
	requestDTO RestoreBackupRequest,
) error {
	backup, err := s.backupService.GetBackup(backupID)
	if err != nil {
		return err
	}

	if backup.Database.UserID != user.ID {
		return errors.New("user does not have access to this backup")
	}

	go func() {
		if err := s.RestoreBackup(backup, requestDTO); err != nil {
			log.Error("Failed to restore backup", "error", err)
		}
	}()

	return nil
}

func (s *RestoreService) RestoreBackup(
	backup *backups.Backup,
	requestDTO RestoreBackupRequest,
) error {
	if backup.Status != backups.BackupStatusCompleted {
		return errors.New("backup is not completed")
	}

	if backup.Database.Type == databases.DatabaseTypePostgres {
		if requestDTO.PostgresqlDatabase == nil {
			return errors.New("postgresql database is required")
		}
	}

	restore := models.Restore{
		ID:     uuid.New(),
		Status: enums.RestoreStatusInProgress,

		BackupID: backup.ID,
		Backup:   backup,

		CreatedAt:         time.Now().UTC(),
		RestoreDurationMs: 0,

		FailMessage: nil,
	}

	// Save the restore first
	if err := s.restoreRepository.Save(&restore); err != nil {
		return err
	}

	// Set the RestoreID on the PostgreSQL database and save it
	if requestDTO.PostgresqlDatabase != nil {
		requestDTO.PostgresqlDatabase.RestoreID = &restore.ID
		restore.Postgresql = requestDTO.PostgresqlDatabase

		// Save the restore again to include the postgresql database
		if err := s.restoreRepository.Save(&restore); err != nil {
			return err
		}
	}

	start := time.Now().UTC()

	err := s.restoreBackupUsecase.Execute(
		restore,
		backup,
	)
	if err != nil {
		errMsg := err.Error()
		restore.FailMessage = &errMsg
		restore.Status = enums.RestoreStatusFailed
		restore.RestoreDurationMs = time.Since(start).Milliseconds()

		if err := s.restoreRepository.Save(&restore); err != nil {
			return err
		}

		return err
	}

	restore.Status = enums.RestoreStatusCompleted
	restore.RestoreDurationMs = time.Since(start).Milliseconds()

	if err := s.restoreRepository.Save(&restore); err != nil {
		return err
	}

	return nil
}
