-- +goose Up
-- create "mcp_connections" table
CREATE TABLE "public"."mcp_connections" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "url" text NOT NULL,
  "agent_id" uuid NOT NULL,
  "header_auth_keys" jsonb NULL,
  "updated_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "mcp_connections" table
DROP TABLE "public"."mcp_connections";
