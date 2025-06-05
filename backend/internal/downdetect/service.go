package downdetect

import (
	"postgresus-backend/internal/storage"
)

type DowndetectService struct {
}

func (s *DowndetectService) IsDbAvailable() error {
	err := storage.GetDb().Exec("SELECT 1").Error
	if err != nil {
		return err
	}

	return nil
}
