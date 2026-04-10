-- +goose Up
-- create "llm_connections" table
CREATE TABLE "public"."llm_connections" (
  "id" uuid NOT NULL,
  "host" text NULL,
  "port" text NULL,
  "model" text NULL,
  "agent_id" uuid NOT NULL,
  "platform" text NOT NULL,
  "setup_status" text NOT NULL,
  "updated_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  "deleted_at" timestamp NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "llm_connections" table
DROP TABLE "public"."llm_connections";
