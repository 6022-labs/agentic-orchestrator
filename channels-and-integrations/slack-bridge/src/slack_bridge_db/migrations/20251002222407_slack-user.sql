-- +goose Up
-- modify "slack_agents" table
ALTER TABLE "public"."slack_agents" ADD COLUMN "owner_wallet" text; -- NOT NULL;
UPDATE "public"."slack_agents" SET "owner_wallet" = '' WHERE "owner_wallet" IS NULL;
ALTER TABLE "public"."slack_agents" ALTER COLUMN "owner_wallet" SET NOT NULL;

-- create "slack_users" table
CREATE TABLE "public"."slack_users" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "real_name" text NOT NULL,
  "wallet_address" text NULL,
  "paired_human_agent_collection_address" text NULL,
  "paired_human_agent_collection_token_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_slack_users_wallet_address" UNIQUE ("wallet_address"),
  CONSTRAINT "fk_slack_users_paired_human_agent" FOREIGN KEY ("paired_human_agent_collection_address", "paired_human_agent_collection_token_id") REFERENCES "public"."slack_agents" ("collection_address", "collection_token_id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- create index "idx_paired_slack_user_to_human_agent_collection" to table: "slack_users"
CREATE UNIQUE INDEX "idx_paired_slack_user_to_human_agent_collection" ON "public"."slack_users" ("paired_human_agent_collection_address", "paired_human_agent_collection_token_id");
-- modify "slack_messages" table
ALTER TABLE "public"."slack_messages" DROP COLUMN "author", ADD COLUMN "is_from_user" boolean NOT NULL, ADD COLUMN "is_from_agent" boolean NOT NULL, ADD COLUMN "author_user_id" text NULL, ADD COLUMN "author_collection_address" text NULL, ADD COLUMN "author_collection_token_id" bigint NULL, ADD CONSTRAINT "fk_slack_messages_author_agent" FOREIGN KEY ("author_collection_address", "author_collection_token_id") REFERENCES "public"."slack_agents" ("collection_address", "collection_token_id") ON UPDATE CASCADE ON DELETE SET NULL, ADD CONSTRAINT "fk_slack_messages_author_user" FOREIGN KEY ("author_user_id") REFERENCES "public"."slack_users" ("id") ON UPDATE CASCADE ON DELETE SET NULL;
-- create index "idx_slack_messages_author_user_id" to table: "slack_messages"
CREATE INDEX "idx_slack_messages_author_user_id" ON "public"."slack_messages" ("author_user_id");

-- +goose Down
-- reverse: create index "idx_slack_messages_author_user_id" to table: "slack_messages"
DROP INDEX "public"."idx_slack_messages_author_user_id";
-- reverse: modify "slack_messages" table
ALTER TABLE "public"."slack_messages" DROP CONSTRAINT "fk_slack_messages_author_user", DROP CONSTRAINT "fk_slack_messages_author_agent", DROP COLUMN "author_collection_token_id", DROP COLUMN "author_collection_address", DROP COLUMN "author_user_id", DROP COLUMN "is_from_agent", DROP COLUMN "is_from_user", ADD COLUMN "author" text NOT NULL;
-- reverse: create index "idx_paired_slack_user_to_human_agent_collection" to table: "slack_users"
DROP INDEX "public"."idx_paired_slack_user_to_human_agent_collection";
-- reverse: create "slack_users" table
DROP TABLE "public"."slack_users";
-- reverse: modify "slack_agents" table
ALTER TABLE "public"."slack_agents" DROP COLUMN "owner_wallet";
