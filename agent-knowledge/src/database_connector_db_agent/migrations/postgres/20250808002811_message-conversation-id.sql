-- +goose Up
-- modify "agent_context_messages" table
ALTER TABLE "public"."agent_context_messages" ADD COLUMN "conversation_id" character varying(128) NULL;

-- +goose Down
-- reverse: modify "agent_context_messages" table
ALTER TABLE "public"."agent_context_messages" DROP COLUMN "conversation_id";
