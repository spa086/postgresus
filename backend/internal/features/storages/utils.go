package storages

import (
	"fmt"
	"os"
	"postgresus-backend/internal/config"
)

// EnsureSystemDirectories creates required system directories for all storage types
// This function should be called before any storage operations
func EnsureSystemDirectories() error {
	// Standard permissions for directories: owner can read/write/execute, others can read/execute
	const directoryPermissions = 0755

	dirs := []string{
		config.GetEnv().DataFolder, // /postgresus-data/backups
		config.GetEnv().TempFolder, // /postgresus-data/temp
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, directoryPermissions); err != nil {
			return fmt.Errorf("failed to create system directory %s: %w", dir, err)
		}
	}

	return nil
}
