# Database migrations

This directory contains the SQL migrations used by the orchestrator service. Migrations are managed using [Atlas](https://github.com/ariga/atlas) together with the [Goose](https://github.com/pressly/goose) formatter. Each database system has its own folder:

- `mariadb`
- `mysql`
- `postgres`
- `sql-server`

Migration files are generated separately for each database dialect.

## Generating migrations

Use the helper script named `generate_migrations.sh` in the migration folder. Provide a descriptive migration name:

```bash
./scripts/generate_agent_migrations.sh add-new-table
```

The script runs `atlas migrate diff` for every supported database and places the
resulting `TIMESTAMP_<name>.sql` files in the corresponding directories above.

## Customizing the migration

### ⚠️ Edit only migrations that is not already played on any server !

After editing the generated SQL file you must recompute the checksums so that Atlas can verify the migration history:

1. Navigate to this folder:

   ```bash
   cd src/orchestrator_agent_database_connector_db_agent/migrations
   ```

2. Regenerate hashes of existing migrations where `atlas.hcl` are stored :

   ```bash
   atlas migrate hash --dir file://<DRIVER>
   ```

This command updates `atlas.sum` with the new file hashes.

## Tools

### loader

Located in `loader/`, this small program loads the GORM schema and prints it in a format that Atlas understands. It is executed automatically by Atlas when running the `diff` command.

### history_updater

Located in `history_updater/`, this utility marks migrations as applied in existing databases. Provide a `config.yml` (see `config.yml.example`) containing database connection strings with the corresponding driver, the highest migration version to mark and optionally the path to the migrations directory. Run it with:

```bash
go run ./history_updater
```

It iterates over the listed connections and marks all migrations up to the specified version as applied.
