-- +goose Up
-- modify "agent_context_messages" table
ALTER TABLE "public"."agent_context_messages" ALTER COLUMN "author" TYPE text;

-- +goose Down
-- reverse: modify "agent_context_messages" table
ALTER TABLE "public"."agent_context_messages" ALTER COLUMN "author" TYPE character varying(64);
