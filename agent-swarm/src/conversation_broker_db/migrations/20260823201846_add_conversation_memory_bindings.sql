-- +goose Up
-- create "conversation_memory_bindings" table
CREATE TABLE "public"."conversation_memory_bindings" (
  "conversation_id" uuid NOT NULL,
  "participant_id" uuid NOT NULL,
  "memory_session_id" uuid NOT NULL,
  "bound_at" timestamptz NOT NULL,
  PRIMARY KEY ("conversation_id", "participant_id"),
  CONSTRAINT "fk_conversation_memory_bindings_conversation" FOREIGN KEY ("conversation_id") REFERENCES "public"."conversations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_conversation_memory_bindings_participant" FOREIGN KEY ("participant_id") REFERENCES "public"."participants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- +goose Down
-- reverse: create "conversation_memory_bindings" table
DROP TABLE "public"."conversation_memory_bindings";
