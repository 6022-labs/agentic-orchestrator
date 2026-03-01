package common_tests

import (
	"os"
	"path/filepath"
)

func FindMigrationFolderPath() (*string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		gitPath := filepath.Join(currentDir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil && info.IsDir() {
			migrationFolderPath := filepath.Join(currentDir, "agent-reasoning", "src", "llm_connector_db", "migrations")
			return &migrationFolderPath, nil
		} else if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached filesystem root
			return nil, os.ErrNotExist
		}
		currentDir = parentDir
	}
}
