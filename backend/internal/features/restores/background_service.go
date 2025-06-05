package restores

import (
	"postgresus-backend/internal/features/restores/enums"
	"postgresus-backend/internal/util/logger"
)

var log = logger.GetLogger()

type RestoreBackgroundService struct {
	restoreRepository *RestoreRepository
}

func (s *RestoreBackgroundService) Run() {
	if err := s.failRestoresInProgress(); err != nil {
		log.Error("Failed to fail restores in progress", "error", err)
		panic(err)
	}
}

func (s *RestoreBackgroundService) failRestoresInProgress() error {
	restoresInProgress, err := s.restoreRepository.FindByStatus(enums.RestoreStatusInProgress)
	if err != nil {
		return err
	}

	for _, restore := range restoresInProgress {
		failMessage := "Restore failed due to application restart"
		restore.Status = enums.RestoreStatusFailed
		restore.FailMessage = &failMessage

		if err := s.restoreRepository.Save(restore); err != nil {
			return err
		}
	}

	return nil
}
