-- +goose Up
-- create "sync_checkpoints" table
CREATE TABLE "public"."sync_checkpoints" (
  "id" uuid NOT NULL,
  "block" bigint NOT NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "sync_checkpoints" table
DROP TABLE "public"."sync_checkpoints";
