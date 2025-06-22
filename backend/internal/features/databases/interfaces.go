package databases

import (
	"log/slog"

	"github.com/google/uuid"
)

type DatabaseValidator interface {
	Validate() error
}

type DatabaseConnector interface {
	TestConnection(logger *slog.Logger) error
}

type DatabaseStorageChangeListener interface {
	OnBeforeDbStorageChange(dbID uuid.UUID, storageID uuid.UUID) error
}
