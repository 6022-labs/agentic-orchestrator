-- +goose Up
-- modify "slack_threads" table
ALTER TABLE "public"."slack_threads" ADD COLUMN "unconcluded_at" timestamp NULL;

-- +goose Down
-- reverse: modify "slack_threads" table
ALTER TABLE "public"."slack_threads" DROP COLUMN "unconcluded_at";
