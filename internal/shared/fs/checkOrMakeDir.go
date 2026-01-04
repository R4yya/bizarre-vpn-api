package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckOrMakeDir(filePath string) error {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error when create dir: %w", err)
	}

	return nil
}
