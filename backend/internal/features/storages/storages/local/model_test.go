package local_storage

import (
	"bytes"
	"io"
	"os"
	"postgresus-backend/internal/config"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LocalStorage_BasicOperations(t *testing.T) {
	storage := &LocalStorage{
		StorageID: uuid.New(),
	}

	t.Run("Test_TestConnection_ConnectionSucceeds", func(t *testing.T) {
		err := storage.TestConnection()
		assert.NoError(t, err, "TestConnection should succeed")
	})

	t.Run("Test_TestValidation_ValidationSucceeds", func(t *testing.T) {
		err := storage.Validate()
		assert.NoError(t, err, "Validate should succeed")
	})

	t.Run("Test_TestSaveAndGetFile_ReturnsCorrectContent", func(t *testing.T) {
		testData, err := os.ReadFile(config.GetEnv().TestFilePath)
		require.NoError(t, err, "Should be able to read test file")

		fileID := uuid.New()

		err = storage.SaveFile(fileID, bytes.NewReader(testData))
		require.NoError(t, err, "SaveFile should succeed")

		file, err := storage.GetFile(fileID)
		assert.NoError(t, err, "GetFile should succeed")
		defer file.Close()

		content, err := io.ReadAll(file)
		assert.NoError(t, err, "Should be able to read file")
		assert.Equal(t, testData, content, "File content should match the original")
	})

	t.Run("Test_TestDeleteFile_RemovesFileFromDisk", func(t *testing.T) {
		// Read test file
		testData, err := os.ReadFile(config.GetEnv().TestFilePath)
		require.NoError(t, err, "Should be able to read test file")

		// Generate a unique file ID
		fileID := uuid.New()

		// Save file first
		err = storage.SaveFile(fileID, bytes.NewReader(testData))
		require.NoError(t, err, "SaveFile should succeed")

		// Delete file
		err = storage.DeleteFile(fileID)
		assert.NoError(t, err, "DeleteFile should succeed")

		file, err := storage.GetFile(fileID)
		assert.Error(t, err, "GetFile should fail for non-existent file")
		assert.Nil(t, file, "File should be nil")
	})

	t.Run("Test_TestDeleteNonExistentFile_DoesNotError", func(t *testing.T) {
		// Try to delete a non-existent file
		nonExistentID := uuid.New()
		err := storage.DeleteFile(nonExistentID)
		assert.NoError(t, err, "DeleteFile should not error for non-existent file")
	})
}
