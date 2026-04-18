-- +goose Up
-- create "slack_agents" table
CREATE TABLE "public"."slack_agents" (
  "collection_address" text NOT NULL,
  "collection_token_id" bigint NOT NULL,
  "ens_domain" text NOT NULL,
  "role" text NOT NULL,
  "icon_image_cid" text NULL,
  "owner_wallet" text NOT NULL,
  "broker_participant_id" uuid NULL,
  PRIMARY KEY ("collection_address", "collection_token_id")
);
-- create "slack_conversations" table
CREATE TABLE "public"."slack_conversations" (
  "slack_channel_id" text NOT NULL,
  "slack_thread_ts" text NOT NULL,
  "broker_conversation_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("slack_channel_id", "slack_thread_ts")
);
-- create index "idx_slack_conversation_broker" to table: "slack_conversations"
CREATE UNIQUE INDEX "idx_slack_conversation_broker" ON "public"."slack_conversations" ("broker_conversation_id");
-- create "slack_messages" table
CREATE TABLE "public"."slack_messages" (
  "slack_channel_id" text NOT NULL,
  "slack_message_ts" text NOT NULL,
  "broker_message_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("slack_channel_id", "slack_message_ts")
);
-- create "slack_channels" table
CREATE TABLE "public"."slack_channels" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- create "slack_channel_agents" table
CREATE TABLE "public"."slack_channel_agents" (
  "slack_channel_id" text NOT NULL,
  "slack_agent_collection_address" text NOT NULL,
  "slack_agent_collection_token_id" bigint NOT NULL,
  PRIMARY KEY ("slack_channel_id", "slack_agent_collection_address", "slack_agent_collection_token_id"),
  CONSTRAINT "fk_slack_channel_agents_slack_agent" FOREIGN KEY ("slack_agent_collection_address", "slack_agent_collection_token_id") REFERENCES "public"."slack_agents" ("collection_address", "collection_token_id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_slack_channel_agents_slack_channel" FOREIGN KEY ("slack_channel_id") REFERENCES "public"."slack_channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- create "slack_users" table
CREATE TABLE "public"."slack_users" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "real_name" text NOT NULL,
  "wallet_address" text NULL,
  "paired_human_agent_collection_address" text NULL,
  "paired_human_agent_collection_token_id" bigint NULL,
  "broker_participant_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_slack_users_wallet_address" UNIQUE ("wallet_address"),
  CONSTRAINT "fk_slack_users_paired_human_agent" FOREIGN KEY ("paired_human_agent_collection_address", "paired_human_agent_collection_token_id") REFERENCES "public"."slack_agents" ("collection_address", "collection_token_id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- create index "idx_paired_slack_user_to_human_agent_collection" to table: "slack_users"
CREATE UNIQUE INDEX "idx_paired_slack_user_to_human_agent_collection" ON "public"."slack_users" ("paired_human_agent_collection_address", "paired_human_agent_collection_token_id");
-- create "slack_channel_users" table
CREATE TABLE "public"."slack_channel_users" (
  "slack_channel_id" text NOT NULL,
  "slack_user_id" text NOT NULL,
  PRIMARY KEY ("slack_channel_id", "slack_user_id"),
  CONSTRAINT "fk_slack_channel_users_slack_channel" FOREIGN KEY ("slack_channel_id") REFERENCES "public"."slack_channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_slack_channel_users_slack_user" FOREIGN KEY ("slack_user_id") REFERENCES "public"."slack_users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

-- +goose Down
-- reverse: create "slack_channel_users" table
DROP TABLE "public"."slack_channel_users";
-- reverse: create index "idx_paired_slack_user_to_human_agent_collection" to table: "slack_users"
DROP INDEX "public"."idx_paired_slack_user_to_human_agent_collection";
-- reverse: create "slack_users" table
DROP TABLE "public"."slack_users";
-- reverse: create "slack_channel_agents" table
DROP TABLE "public"."slack_channel_agents";
-- reverse: create "slack_channels" table
DROP TABLE "public"."slack_channels";
-- reverse: create "slack_messages" table
DROP TABLE "public"."slack_messages";
-- reverse: create index "idx_slack_conversation_broker" to table: "slack_conversations"
DROP INDEX "public"."idx_slack_conversation_broker";
-- reverse: create "slack_conversations" table
DROP TABLE "public"."slack_conversations";
-- reverse: create "slack_agents" table
DROP TABLE "public"."slack_agents";
