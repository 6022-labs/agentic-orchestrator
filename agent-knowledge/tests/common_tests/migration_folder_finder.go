package common_tests

import (
	"os"
	"path/filepath"
)

func findGitRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		gitPath := filepath.Join(currentDir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil && info.IsDir() {
			return currentDir, nil
		} else if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached filesystem root
			return "", os.ErrNotExist
		}
		currentDir = parentDir
	}
}

func findMigrationFolderPath(subPath string) (*string, error) {
	gitRoot, err := findGitRoot()
	if err != nil {
		return nil, err
	}
	migrationFolderPath := filepath.Join(gitRoot, subPath)
	return &migrationFolderPath, nil
}

func FindAgentMigrationFolderPath() (*string, error) {
	return findMigrationFolderPath("agent-knowledge/src/database_connector_db_agent/migrations")

}

func FindOrchestratorMigrationFolderPath() (*string, error) {
	return findMigrationFolderPath("agent-knowledge/src/database_connector_db_orchestrator/migrations")
}
