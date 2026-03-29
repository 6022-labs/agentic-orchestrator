-- +goose Up
-- create "sync_checkpoints" table
CREATE TABLE "public"."sync_checkpoints" (
  "block" bigserial NOT NULL,
  PRIMARY KEY ("block")
);

-- +goose Down
-- reverse: create "sync_checkpoints" table
DROP TABLE "public"."sync_checkpoints";
