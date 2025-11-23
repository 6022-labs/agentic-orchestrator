-- +goose Up
-- create "database_connection_migrations" table
CREATE TABLE "public"."database_connection_migrations" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "updated_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
-- create "database_connection_migration_connections" table
CREATE TABLE "public"."database_connection_migration_connections" (
  "database_connection_migration_id" uuid NOT NULL,
  "database_connection_id" uuid NOT NULL,
  PRIMARY KEY ("database_connection_migration_id", "database_connection_id"),
  CONSTRAINT "fk_database_connection_migration_connections_database_co31c61dc" FOREIGN KEY ("database_connection_migration_id") REFERENCES "public"."database_connection_migrations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_database_connection_migration_connections_database_connectio" FOREIGN KEY ("database_connection_id") REFERENCES "public"."database_connections" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

INSERT INTO "public"."database_connection_migrations" ("id", "name", "updated_at", "created_at") VALUES
  ('00000000-0000-0000-0000-000000000001', '20250729144720_initial-migration.sql', NOW(), NOW());

INSERT INTO "public"."database_connection_migration_connections" (
  "database_connection_migration_id",
  "database_connection_id"
)
SELECT
  '00000000-0000-0000-0000-000000000001' AS database_connection_migration_id,
  id AS database_connection_id
FROM
  "public"."database_connections";

-- +goose Down
-- reverse: create "database_connection_migration_connections" table
DROP TABLE "public"."database_connection_migration_connections";
-- reverse: create "database_connection_migrations" table
DROP TABLE "public"."database_connection_migrations";
