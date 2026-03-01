package utils

import "fmt"

func GetMigrationFolderPath(migrationsFolderPath string, driver string) string {
	return fmt.Sprintf("%s/%s", migrationsFolderPath, driver)
}
