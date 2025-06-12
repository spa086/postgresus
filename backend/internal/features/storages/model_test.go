package storages

import (
	"bytes"
	"io"
	"os"
	"postgresus-backend/internal/config"
	local_storage "postgresus-backend/internal/features/storages/storages/local"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Storage_BasicOperations(t *testing.T) {
	testCases := []struct {
		name    string
		storage StorageFileSaver
	}{
		{
			name:    "LocalStorage",
			storage: &local_storage.LocalStorage{StorageID: uuid.New()},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("Test_TestConnection_ConnectionSucceeds", func(t *testing.T) {
				err := tc.storage.TestConnection()
				assert.NoError(t, err, "TestConnection should succeed")
			})

			t.Run("Test_TestValidation_ValidationSucceeds", func(t *testing.T) {
				err := tc.storage.Validate()
				assert.NoError(t, err, "Validate should succeed")
			})

			t.Run("Test_TestSaveAndGetFile_ReturnsCorrectContent", func(t *testing.T) {
				testData, err := os.ReadFile(config.GetEnv().TestFilePath)
				require.NoError(t, err, "Should be able to read test file")

				fileID := uuid.New()

				err = tc.storage.SaveFile(fileID, bytes.NewReader(testData))
				require.NoError(t, err, "SaveFile should succeed")

				file, err := tc.storage.GetFile(fileID)
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
				err = tc.storage.SaveFile(fileID, bytes.NewReader(testData))
				require.NoError(t, err, "SaveFile should succeed")

				// Delete file
				err = tc.storage.DeleteFile(fileID)
				assert.NoError(t, err, "DeleteFile should succeed")

				// Try to get the deleted file
				file, err := tc.storage.GetFile(fileID)
				assert.Error(t, err, "GetFile should fail for non-existent file")
				if file != nil {
					file.Close()
				}
			})

			t.Run("Test_TestDeleteNonExistentFile_DoesNotError", func(t *testing.T) {
				// Try to delete a non-existent file
				nonExistentID := uuid.New()
				err := tc.storage.DeleteFile(nonExistentID)
				assert.NoError(t, err, "DeleteFile should not error for non-existent file")
			})
		})
	}
}
