# Database migrations

This directory contains the SQL migrations used by the compute database (`compute_db` module). Migrations are managed using [Atlas](https://github.com/ariga/atlas) together with the [Goose](https://github.com/pressly/goose) formatter. The `atlas.hcl` file describes how the schema is loaded and where migrations are written.

## Creating a new migration

1. Navigate to this folder:

   ```bash
   cd agent-reasoning/src/compute_db/migrations
   ```

2. Generate the migration with Atlas using the `gorm` environment defined in `atlas.hcl`:

   ```bash
   atlas migrate diff <name> --env gorm
   ```

   The command runs the `loader` program to load the current GORM models and writes a new file named `TIMESTAMP_<name>.sql` together with an updated `atlas.sum` file.

## Customizing the migration

### ⚠️ Edit only migrations that have not been applied to any environment!

After editing the generated SQL file you must recompute the checksums so that Atlas can verify the migration history:

1. Navigate to this folder:

   ```bash
   cd agent-reasoning/src/compute_db/migrations
   ```

2. Regenerate hashes of existing migrations:

   ```bash
   atlas migrate hash --dir file://.
   ```

This command updates `atlas.sum` with the new file hashes.

## Tools

### loader

Located in `loader/`, this small program loads the GORM schema and prints it in a format that Atlas understands. It is executed automatically by Atlas when running the `diff` command.

### history_updater

Located in `history_updater/`, this utility marks migrations as applied in existing databases. Provide a `config.yml` (see `config.yml.example`) containing database connection strings, the highest migration version to mark and optionally the path to the migrations directory. Run it with:

```bash
go run ./history_updater
```

It iterates over the listed connections and marks all migrations up to the specified version as applied.
