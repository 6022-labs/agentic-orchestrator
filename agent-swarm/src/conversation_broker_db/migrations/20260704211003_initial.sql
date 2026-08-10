-- +goose Up
-- create "agent_prompt_templates" table
CREATE TABLE "public"."agent_prompt_templates" (
  "id" uuid NOT NULL,
  "collection_address" text NOT NULL,
  "collection_token_id" numeric NOT NULL,
  "key" text NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_agent_prompt_templates_agent_key" to table: "agent_prompt_templates"
CREATE UNIQUE INDEX "idx_agent_prompt_templates_agent_key" ON "public"."agent_prompt_templates" ("collection_address", "collection_token_id", "key");
-- create "bridges" table
CREATE TABLE "public"."bridges" (
  "id" uuid NOT NULL,
  "kind" text NOT NULL,
  "base_url" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_bridges_kind" to table: "bridges"
CREATE INDEX "idx_bridges_kind" ON "public"."bridges" ("kind");
-- create "conversations" table
CREATE TABLE "public"."conversations" (
  "id" uuid NOT NULL,
  "bridge_id" uuid NOT NULL,
  "description" text NULL,
  "status" text NOT NULL DEFAULT 'active',
  "concluded_at" timestamptz NULL,
  "unconcluded_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_conversations_bridge" FOREIGN KEY ("bridge_id") REFERENCES "public"."bridges" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_conversations_bridge_id" to table: "conversations"
CREATE INDEX "idx_conversations_bridge_id" ON "public"."conversations" ("bridge_id");
-- create index "idx_conversations_status_updated" to table: "conversations"
CREATE INDEX "idx_conversations_status_updated" ON "public"."conversations" ("status", "updated_at");
-- create "participants" table
CREATE TABLE "public"."participants" (
  "id" uuid NOT NULL,
  "collection_address" text NULL,
  "collection_token_id" numeric NULL,
  "ens_domain" text NULL,
  "display_name" text NULL,
  "is_processing" boolean NOT NULL DEFAULT false,
  "last_processed_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_participants_collection" to table: "participants"
CREATE UNIQUE INDEX "idx_participants_collection" ON "public"."participants" ("collection_address", "collection_token_id");
-- create index "idx_participants_processing" to table: "participants"
CREATE INDEX "idx_participants_processing" ON "public"."participants" ("is_processing", "last_processed_at");
-- create "conversation_messages" table
CREATE TABLE "public"."conversation_messages" (
  "id" uuid NOT NULL,
  "conversation_id" uuid NOT NULL,
  "participant_id" uuid NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_conversation_messages_conversation" FOREIGN KEY ("conversation_id") REFERENCES "public"."conversations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_conversation_messages_participant" FOREIGN KEY ("participant_id") REFERENCES "public"."participants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_messages_conv_time" to table: "conversation_messages"
CREATE INDEX "idx_messages_conv_time" ON "public"."conversation_messages" ("conversation_id", "created_at");
-- create "conversation_participants" table
CREATE TABLE "public"."conversation_participants" (
  "conversation_id" uuid NOT NULL,
  "participant_id" uuid NOT NULL,
  "joined_at" timestamptz NOT NULL,
  PRIMARY KEY ("conversation_id", "participant_id"),
  CONSTRAINT "fk_conversation_participants_conversation" FOREIGN KEY ("conversation_id") REFERENCES "public"."conversations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_conversation_participants_participant" FOREIGN KEY ("participant_id") REFERENCES "public"."participants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- +goose Down
-- reverse: create "conversation_participants" table
DROP TABLE "public"."conversation_participants";
-- reverse: create index "idx_messages_conv_time" to table: "conversation_messages"
DROP INDEX "public"."idx_messages_conv_time";
-- reverse: create "conversation_messages" table
DROP TABLE "public"."conversation_messages";
-- reverse: create index "idx_participants_processing" to table: "participants"
DROP INDEX "public"."idx_participants_processing";
-- reverse: create index "idx_participants_collection" to table: "participants"
DROP INDEX "public"."idx_participants_collection";
-- reverse: create "participants" table
DROP TABLE "public"."participants";
-- reverse: create index "idx_conversations_status_updated" to table: "conversations"
DROP INDEX "public"."idx_conversations_status_updated";
-- reverse: create index "idx_conversations_bridge_id" to table: "conversations"
DROP INDEX "public"."idx_conversations_bridge_id";
-- reverse: create "conversations" table
DROP TABLE "public"."conversations";
-- reverse: create index "idx_bridges_kind" to table: "bridges"
DROP INDEX "public"."idx_bridges_kind";
-- reverse: create "bridges" table
DROP TABLE "public"."bridges";
-- reverse: create index "idx_agent_prompt_templates_agent_key" to table: "agent_prompt_templates"
DROP INDEX "public"."idx_agent_prompt_templates_agent_key";
-- reverse: create "agent_prompt_templates" table
DROP TABLE "public"."agent_prompt_templates";
