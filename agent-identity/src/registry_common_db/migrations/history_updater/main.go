package main

import (
	"database/sql"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

const CONFIGURATION_FILE = "config.yml"

type config struct {
	Connections   []string `json:"connections"    yaml:"connections"`
	Until         int64    `json:"until"          yaml:"until"`
	MigrationsDir string   `json:"migrations_dir" yaml:"migrations_dir"`
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

	availableVersions := readAvailableVersions(dir)
	versionsToMark := lo.Filter(availableVersions, func(v int64, _ int) bool {
		return v <= cfg.Until
	})
	if len(versionsToMark) == 0 {
		log.Fatal("no versions to mark")
	}

	sort.Slice(versionsToMark, func(i, j int) bool {
		return versionsToMark[i] < versionsToMark[j]
	})

	for connIndex, conn := range cfg.Connections {
		conn = strings.TrimSpace(conn)
		if len(strings.TrimSpace(conn)) == 0 {
			log.Printf("skipping empty connection at index %d", connIndex)
			continue
		}

		if err := mark(conn, dir, versionsToMark); err != nil {
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

func readAvailableVersions(dir string) []int64 {
	entries, err := os.ReadDir(dir)
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

func mark(conn, migrationsDir string, versions []int64) error {
	if len(versions) == 0 {
		return nil
	}

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Ensure goose uses the correct dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	currentVersion, err := goose.EnsureDBVersion(db)
	if err != nil {
		return err
	}

	tableName := goose.TableName()

	for _, v := range versions {
		if v <= currentVersion {
			log.Printf("skipping version %d, already applied in %s", v, conn)
			continue
		}

		log.Printf("marking version %d in %s", v, conn)
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
		log.Printf("marked version %d in %s", v, conn)
	}

	return nil
}
