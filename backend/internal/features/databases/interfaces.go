package databases

import "github.com/google/uuid"

type DatabaseValidator interface {
	Validate() error
}

type DatabaseConnector interface {
	TestConnection() error
}

type DatabaseStorageChangeListener interface {
	OnBeforeDbStorageChange(dbID uuid.UUID, storageID uuid.UUID) error
}
