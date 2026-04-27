package common_tests

import (
	"os"
	"path/filepath"
)

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

func findMigrationFolderPath(subPath string) (*string, error) {
	root, err := findModuleRoot()
	if err != nil {
		return nil, err
	}
	migrationFolderPath := filepath.Join(root, subPath)
	return &migrationFolderPath, nil
}

func FindAgentMigrationFolderPath() (*string, error) {
	return findMigrationFolderPath("src/database_connector_db_agent/migrations")
}

func FindOrchestratorMigrationFolderPath() (*string, error) {
	return findMigrationFolderPath("src/database_connector_db_orchestrator/migrations")
}
