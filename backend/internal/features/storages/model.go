package storages

import (
	"errors"
	"io"
	local_storage "postgresus-backend/internal/features/storages/storages/local"
	s3_storage "postgresus-backend/internal/features/storages/storages/s3"

	"github.com/google/uuid"
)

type Storage struct {
	ID            uuid.UUID   `json:"id"            gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        uuid.UUID   `json:"userId"        gorm:"column:user_id;not null;type:uuid;index"`
	Type          StorageType `json:"type"          gorm:"column:type;not null;type:text"`
	Name          string      `json:"name"          gorm:"column:name;not null;type:text"`
	LastSaveError *string     `json:"lastSaveError" gorm:"column:last_save_error;type:text"`

	// specific storage
	LocalStorage *local_storage.LocalStorage `json:"localStorage" gorm:"foreignKey:StorageID"`
	S3Storage    *s3_storage.S3Storage       `json:"s3Storage"    gorm:"foreignKey:StorageID"`
}

func (s *Storage) SaveFile(fileID uuid.UUID, file io.Reader) error {
	err := s.getSpecificStorage().SaveFile(fileID, file)
	if err != nil {
		lastSaveError := err.Error()
		s.LastSaveError = &lastSaveError
		return err
	}

	s.LastSaveError = nil

	return nil
}

func (s *Storage) GetFile(fileID uuid.UUID) (io.ReadCloser, error) {
	return s.getSpecificStorage().GetFile(fileID)
}

func (s *Storage) DeleteFile(fileID uuid.UUID) error {
	return s.getSpecificStorage().DeleteFile(fileID)
}

func (s *Storage) Validate() error {
	if s.Type == "" {
		return errors.New("storage type is required")
	}

	if s.Name == "" {
		return errors.New("storage name is required")
	}

	return s.getSpecificStorage().Validate()
}

func (s *Storage) TestConnection() error {
	return s.getSpecificStorage().TestConnection()
}

func (s *Storage) getSpecificStorage() StorageFileSaver {
	switch s.Type {
	case StorageTypeLocal:
		return s.LocalStorage
	case StorageTypeS3:
		return s.S3Storage
	default:
		panic("invalid storage type: " + string(s.Type))
	}
}
