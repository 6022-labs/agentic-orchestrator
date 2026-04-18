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

func findMigrationFolderFor(module string) (*string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for {
		gitPath := filepath.Join(currentDir, ".git")
		_, err := os.Stat(gitPath)
		if err == nil {
			migrationFolderPath := filepath.Join(currentDir, "agent-reasoning", "src", module, "migrations")
			return &migrationFolderPath, nil
		} else if !os.IsNotExist(err) {
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
