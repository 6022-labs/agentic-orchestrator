package common_tests

import (
	"os"
	"path/filepath"
)

func FindMigrationFolderPath() (*string, error) {
	return findMigrationFolderFor("agent_compute_db")
}

func FindConversationBrokerMigrationFolderPath() (*string, error) {
	return findMigrationFolderFor("conversation_broker_db")
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

func findMigrationFolderFor(module string) (*string, error) {
	root, err := findModuleRoot()
	if err != nil {
		return nil, err
	}
	migrationFolderPath := filepath.Join(root, "src", module, "migrations")
	return &migrationFolderPath, nil
}
