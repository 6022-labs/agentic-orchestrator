-- +goose Up
-- create "database_connections" table
CREATE TABLE "public"."database_connections" (
  "id" uuid NOT NULL,
  "host" text NOT NULL,
  "port" integer NOT NULL,
  "driver" text NOT NULL,
  "agent_id" uuid NOT NULL,
  "setup_status" text NOT NULL,
  "updated_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  "deleted_at" timestamp NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "database_connections" table
DROP TABLE "public"."database_connections";
