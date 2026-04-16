-- +goose Up
-- create "slack_conversation_links" table
CREATE TABLE "public"."slack_conversation_links" (
  "slack_channel_id" text NOT NULL,
  "slack_thread_ts" text NOT NULL,
  "broker_conversation_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("slack_channel_id", "slack_thread_ts")
);
-- create index "idx_slack_conversation_link_broker" to table: "slack_conversation_links"
CREATE UNIQUE INDEX "idx_slack_conversation_link_broker" ON "public"."slack_conversation_links" ("broker_conversation_id");
-- create "slack_message_links" table
CREATE TABLE "public"."slack_message_links" (
  "slack_channel_id" text NOT NULL,
  "slack_message_ts" text NOT NULL,
  "broker_message_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("slack_channel_id", "slack_message_ts")
);

-- +goose Down
-- reverse: create "slack_message_links" table
DROP TABLE "public"."slack_message_links";
-- reverse: create index "idx_slack_conversation_link_broker" to table: "slack_conversation_links"
DROP INDEX "public"."idx_slack_conversation_link_broker";
-- reverse: create "slack_conversation_links" table
DROP TABLE "public"."slack_conversation_links";
