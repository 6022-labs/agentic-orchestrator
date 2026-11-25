-- +goose Up
-- drop index "idx_participants_processing" from table: "participants"
DROP INDEX "public"."idx_participants_processing";
-- modify "participants" table
ALTER TABLE "public"."participants" DROP COLUMN "is_processing", ADD COLUMN "lease_id" uuid NULL, ADD COLUMN "lease_until" timestamptz NULL;
-- create index "idx_participants_processing" to table: "participants"
CREATE INDEX "idx_participants_processing" ON "public"."participants" ("lease_until", "last_processed_at");

-- +goose Down
-- reverse: create index "idx_participants_processing" to table: "participants"
DROP INDEX "public"."idx_participants_processing";
-- reverse: modify "participants" table
ALTER TABLE "public"."participants" DROP COLUMN "lease_until", DROP COLUMN "lease_id", ADD COLUMN "is_processing" boolean NOT NULL DEFAULT false;
-- reverse: drop index "idx_participants_processing" from table: "participants"
CREATE INDEX "idx_participants_processing" ON "public"."participants" ("is_processing", "last_processed_at");
