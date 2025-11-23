-- +goose Up
-- create "slack_agents" table
CREATE TABLE "public"."slack_agents" (
  "collection_address" text NOT NULL,
  "collection_token_id" bigint NOT NULL,
  "name" text NOT NULL,
  "ens_domain" text NOT NULL,
  "role" text NOT NULL,
  "icon_image_cid" text NULL,
  "is_processing" boolean NOT NULL,
  "last_processed_at" timestamp NULL,
  PRIMARY KEY ("collection_address", "collection_token_id")
);
-- create "slack_messages" table
CREATE TABLE "public"."slack_messages" (
  "id" text NOT NULL,
  "author" text NOT NULL,
  "content" text NOT NULL,
  "thread_id" text NOT NULL,
  "channel_id" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create "slack_threads" table
CREATE TABLE "public"."slack_threads" (
  "id" text NOT NULL,
  "channel_id" text NOT NULL,
  "is_concluded" boolean NOT NULL,
  PRIMARY KEY ("id", "channel_id")
);
-- create "slack_channels" table
CREATE TABLE "public"."slack_channels" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "description" text NULL,
  "members" text NULL,
  PRIMARY KEY ("id")
);
-- create "slack_channel_agents" table
CREATE TABLE "public"."slack_channel_agents" (
  "slack_channel_id" text NOT NULL,
  "slack_agent_collection_address" text NOT NULL,
  "slack_agent_collection_token_id" bigint NOT NULL,
  PRIMARY KEY ("slack_channel_id", "slack_agent_collection_address", "slack_agent_collection_token_id"),
  CONSTRAINT "fk_slack_channel_agents_slack_agent" FOREIGN KEY ("slack_agent_collection_address", "slack_agent_collection_token_id") REFERENCES "public"."slack_agents" ("collection_address", "collection_token_id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_slack_channel_agents_slack_channel" FOREIGN KEY ("slack_channel_id") REFERENCES "public"."slack_channels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- +goose Down
-- reverse: create "slack_channel_agents" table
DROP TABLE "public"."slack_channel_agents";
-- reverse: create "slack_channels" table
DROP TABLE "public"."slack_channels";
-- reverse: create "slack_threads" table
DROP TABLE "public"."slack_threads";
-- reverse: create "slack_messages" table
DROP TABLE "public"."slack_messages";
-- reverse: create "slack_agents" table
DROP TABLE "public"."slack_agents";
