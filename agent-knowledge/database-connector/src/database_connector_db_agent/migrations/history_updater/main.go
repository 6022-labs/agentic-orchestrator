package main

import (
	"database/sql"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const CONFIGURATION_FILE = "config.yml"

type config struct {
	Connections []struct {
		Driver           enums.DatabaseConnectionDriver `json:"driver" yaml:"driver"`
		ConnectionString string                         `json:"connection_string" yaml:"connection_string"`
	} `json:"connections" yaml:"connections"`
	Until         int64  `json:"until" yaml:"until"`
	MigrationsDir string `json:"migrations_dir" yaml:"migrations_dir"`
}

func main() {
	cfg := loadConfig(CONFIGURATION_FILE)
	dir := lo.Ternary(cfg.MigrationsDir != "", cfg.MigrationsDir, "..")

	if len(cfg.Connections) == 0 {
		log.Fatal("no connections provided in config")
	}
	if cfg.Until <= 0 {
		log.Fatal("invalid 'until' version in config, must be positive")
	}

	for connIndex, conn := range cfg.Connections {
		connStr := strings.TrimSpace(conn.ConnectionString)
		if len(strings.TrimSpace(connStr)) == 0 {
			log.Printf("skipping empty connection at index %d", connIndex)
			continue
		}

		if !enums.GetAllDatabaseConnectionDrivers().Contains(conn.Driver) {
			log.Printf("skipping connection at index %d with unsupported driver %s", connIndex, conn.Driver)
			continue
		}

		availableVersions := readAvailableVersions(dir, conn.Driver)
		versionsToMark := lo.Filter(availableVersions, func(v int64, _ int) bool {
			return v <= cfg.Until
		})
		if len(versionsToMark) == 0 {
			log.Fatal("no versions to mark")
		}

		sort.Slice(versionsToMark, func(i, j int) bool {
			return versionsToMark[i] < versionsToMark[j]
		})

		if err := mark(connStr, conn.Driver, versionsToMark); err != nil {
			log.Fatalf("error updating %s: %v", conn, err)
		}

		log.Printf("updated %s\n", conn)
	}
}

func loadConfig(path string) config {
	cfgBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed reading config %s: %v", path, err)
	}

	var cfg config
	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		log.Fatalf("failed parsing yaml config: %v", err)
	}

	return cfg
}

func readAvailableVersions(dir string, driver enums.DatabaseConnectionDriver) []int64 {
	entries, err := os.ReadDir(dir + "/" + driver.String())
	if err != nil {
		log.Fatalf("failed reading migrations: %v", err)
	}

	files := lo.Filter(entries, func(e os.DirEntry, _ int) bool {
		return strings.HasSuffix(e.Name(), ".sql")
	})

	versions := lo.FilterMap(files, func(e os.DirEntry, _ int) (int64, bool) {
		parts := strings.SplitN(e.Name(), "_", 2)
		if len(parts) == 0 {
			return 0, false
		}
		v, err := strconv.ParseInt(parts[0], 10, 64)
		return v, err == nil
	})

	return lo.Uniq(versions)
}

func mark(
	connStr string,
	driver enums.DatabaseConnectionDriver,
	versions []int64,
) error {
	if len(versions) == 0 {
		return nil
	}

	db, err := sql.Open(driver.String(), connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Ensure goose uses the correct dialect
	if err := goose.SetDialect(driver.ToGooseDialect()); err != nil {
		return err
	}

	currentVersion, err := goose.EnsureDBVersion(db)
	if err != nil {
		return err
	}

	tableName := goose.TableName()

	for _, v := range versions {
		if v <= currentVersion {
			log.Printf("skipping version %d, already applied in %s", v, connStr)
			continue
		}

		log.Printf("marking version %d in %s", v, connStr)
		updateQuery := "UPDATE " + tableName + " SET is_applied = $1 WHERE version_id = $2"
		res, err := db.Exec(updateQuery, true, v)
		if err != nil {
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			insertQuery := "INSERT INTO " + tableName + " (version_id, is_applied) VALUES ($1, $2)"
			_, err := db.Exec(insertQuery, v, true)
			if err != nil {
				return err
			}
		}
		log.Printf("marked version %d in %s", v, connStr)
	}

	return nil
}
