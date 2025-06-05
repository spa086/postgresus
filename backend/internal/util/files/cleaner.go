package files_utils

import "os"

func CleanFolder(folder string) error {
	return os.RemoveAll(folder)
}
