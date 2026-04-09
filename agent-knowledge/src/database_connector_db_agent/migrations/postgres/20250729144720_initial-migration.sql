-- +goose Up
-- create "agent_contexts" table
CREATE TABLE "public"."agent_contexts" (
  "id" character(36) NOT NULL,
  "key" character varying(256) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_agent_contexts_key" to table: "agent_contexts"
CREATE UNIQUE INDEX "idx_agent_contexts_key" ON "public"."agent_contexts" ("key");
-- create "agent_prompt_templates" table
CREATE TABLE "public"."agent_prompt_templates" (
  "id" character(36) NOT NULL,
  "key" character varying(64) NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_agent_prompt_templates_key" UNIQUE ("key")
);
-- create "agent_context_conversation_mappings" table
CREATE TABLE "public"."agent_context_conversation_mappings" (
  "id" character(36) NOT NULL,
  "conversation_id" character varying(1024) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "agent_context_id" character(36) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_agent_contexts_conversation_mappings" FOREIGN KEY ("agent_context_id") REFERENCES "public"."agent_contexts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_agent_context_conversation_mappings_agent_context_id" to table: "agent_context_conversation_mappings"
CREATE INDEX "idx_agent_context_conversation_mappings_agent_context_id" ON "public"."agent_context_conversation_mappings" ("agent_context_id");
-- create "agent_context_messages" table
CREATE TABLE "public"."agent_context_messages" (
  "id" character varying(36) NOT NULL,
  "content" text NOT NULL,
  "author" character varying(64) NOT NULL,
  "platform" character varying(64) NOT NULL,
  "role" character varying(64) NOT NULL,
  "hash" character varying(128) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "agent_context_id" character(36) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_agent_contexts_messages" FOREIGN KEY ("agent_context_id") REFERENCES "public"."agent_contexts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_agent_context_messages_agent_context_id" to table: "agent_context_messages"
CREATE INDEX "idx_agent_context_messages_agent_context_id" ON "public"."agent_context_messages" ("agent_context_id");

-- +goose Down
-- reverse: create index "idx_agent_context_messages_agent_context_id" to table: "agent_context_messages"
DROP INDEX "public"."idx_agent_context_messages_agent_context_id";
-- reverse: create "agent_context_messages" table
DROP TABLE "public"."agent_context_messages";
-- reverse: create index "idx_agent_context_conversation_mappings_agent_context_id" to table: "agent_context_conversation_mappings"
DROP INDEX "public"."idx_agent_context_conversation_mappings_agent_context_id";
-- reverse: create "agent_context_conversation_mappings" table
DROP TABLE "public"."agent_context_conversation_mappings";
-- reverse: create "agent_prompt_templates" table
DROP TABLE "public"."agent_prompt_templates";
-- reverse: create index "idx_agent_contexts_key" to table: "agent_contexts"
DROP INDEX "public"."idx_agent_contexts_key";
-- reverse: create "agent_contexts" table
DROP TABLE "public"."agent_contexts";
