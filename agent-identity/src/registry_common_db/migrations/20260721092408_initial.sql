-- +goose Up
-- create "agent_addresses" table
CREATE TABLE "public"."agent_addresses" (
  "agent_id" uuid NOT NULL,
  "type" text NOT NULL,
  "address" text NOT NULL,
  "updated_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("agent_id", "type")
);
-- create "agent_collections" table
CREATE TABLE "public"."agent_collections" (
  "address" text NOT NULL,
  "name" text NOT NULL,
  "symbol" text NOT NULL,
  "updated_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  "deleted_at" timestamp NULL,
  PRIMARY KEY ("address")
);
-- create index "idx_agent_collections_address" to table: "agent_collections"
CREATE UNIQUE INDEX "idx_agent_collections_address" ON "public"."agent_collections" ("address");
-- create "agents" table
CREATE TABLE "public"."agents" (
  "id" uuid NOT NULL,
  "collection_token_id" numeric NULL,
  "collection_address" text NULL,
  "proposal_id" numeric NULL,
  "name" text NOT NULL,
  "clone_of" numeric NULL,
  "ens_domain" text NOT NULL,
  "creator" text NOT NULL,
  "owner" text NULL,
  "images" jsonb NULL,
  "attributes" jsonb NULL,
  "updated_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL,
  "deleted_at" timestamp NULL,
  PRIMARY KEY ("id")
);
-- create "sync_checkpoints" table
CREATE TABLE "public"."sync_checkpoints" (
  "id" uuid NOT NULL,
  "block" bigint NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "sync_checkpoints" table
DROP TABLE "public"."sync_checkpoints";
-- reverse: create "agents" table
DROP TABLE "public"."agents";
-- reverse: create index "idx_agent_collections_address" to table: "agent_collections"
DROP INDEX "public"."idx_agent_collections_address";
-- reverse: create "agent_collections" table
DROP TABLE "public"."agent_collections";
-- reverse: create "agent_addresses" table
DROP TABLE "public"."agent_addresses";
