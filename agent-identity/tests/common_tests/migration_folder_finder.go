package common_tests

import (
	"os"
	"path/filepath"
)

func FindMigrationFolderPath() (*string, error) {
	root, err := findModuleRoot()
	if err != nil {
		return nil, err
	}
	migrationFolderPath := filepath.Join(root, "src", "registry_common_db", "migrations")
	return &migrationFolderPath, nil
}

func findModuleRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", os.ErrNotExist
		}
		currentDir = parentDir
	}
}
