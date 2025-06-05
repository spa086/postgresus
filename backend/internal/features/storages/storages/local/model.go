package local_storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"postgresus-backend/internal/config"
	"postgresus-backend/internal/util/logger"

	"github.com/google/uuid"
)

var log = logger.GetLogger()

// LocalStorage uses ./postgresus_local_backups folder as a
// directory for backups and ./postgresus_local_temp folder as a
// directory for temp files
type LocalStorage struct {
	StorageID uuid.UUID `json:"storageId" gorm:"primaryKey;type:uuid;column:storage_id"`
}

func (l *LocalStorage) TableName() string {
	return "local_storages"
}

func (l *LocalStorage) SaveFile(fileID uuid.UUID, file io.Reader) error {
	log.Info("Starting to save file to local storage", "fileId", fileID.String())

	if err := l.ensureDirectories(); err != nil {
		log.Error("Failed to ensure directories", "fileId", fileID.String(), "error", err)
		return err
	}

	tempFilePath := filepath.Join(config.GetEnv().TempFolder, fileID.String())
	log.Debug("Creating temp file", "fileId", fileID.String(), "tempPath", tempFilePath)

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Error(
			"Failed to create temp file",
			"fileId",
			fileID.String(),
			"tempPath",
			tempFilePath,
			"error",
			err,
		)
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = tempFile.Close()
	}()

	log.Debug("Copying file data to temp file", "fileId", fileID.String())
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Error("Failed to write to temp file", "fileId", fileID.String(), "error", err)
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err = tempFile.Sync(); err != nil {
		log.Error("Failed to sync temp file", "fileId", fileID.String(), "error", err)
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	err = tempFile.Close()
	if err != nil {
		log.Error("Failed to close temp file", "fileId", fileID.String(), "error", err)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	finalPath := filepath.Join(config.GetEnv().DataFolder, fileID.String())
	log.Debug(
		"Moving file from temp to final location",
		"fileId",
		fileID.String(),
		"finalPath",
		finalPath,
	)

	// Move the file from temp to backups directory
	if err = os.Rename(tempFilePath, finalPath); err != nil {
		log.Error(
			"Failed to move file from temp to backups",
			"fileId",
			fileID.String(),
			"tempPath",
			tempFilePath,
			"finalPath",
			finalPath,
			"error",
			err,
		)
		return fmt.Errorf("failed to move file from temp to backups: %w", err)
	}

	log.Info(
		"Successfully saved file to local storage",
		"fileId",
		fileID.String(),
		"finalPath",
		finalPath,
	)
	return nil
}

func (l *LocalStorage) GetFile(fileID uuid.UUID) (io.ReadCloser, error) {
	filePath := filepath.Join(config.GetEnv().DataFolder, fileID.String())

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", fileID.String())
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

func (l *LocalStorage) DeleteFile(fileID uuid.UUID) error {
	filePath := filepath.Join(config.GetEnv().DataFolder, fileID.String())

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (l *LocalStorage) Validate() error {
	return l.ensureDirectories()
}

func (l *LocalStorage) TestConnection() error {
	if err := l.ensureDirectories(); err != nil {
		return err
	}

	testFile := filepath.Join(config.GetEnv().TempFolder, "test_connection")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("failed to close test file: %w", err)
	}

	if err = os.Remove(testFile); err != nil {
		return fmt.Errorf("failed to remove test file: %w", err)
	}

	return nil
}

func (l *LocalStorage) ensureDirectories() error {
	// Standard permissions for directories: owner
	// can read/write/execute, others can read/execute
	const directoryPermissions = 0755

	if err := os.MkdirAll(config.GetEnv().DataFolder, directoryPermissions); err != nil {
		return fmt.Errorf("failed to create backups directory: %w", err)
	}

	if err := os.MkdirAll(config.GetEnv().TempFolder, directoryPermissions); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	return nil
}
