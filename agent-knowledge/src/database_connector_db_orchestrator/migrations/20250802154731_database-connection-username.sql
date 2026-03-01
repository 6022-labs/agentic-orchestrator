-- +goose Up
-- modify "database_connections" table
-- Step 1: Add column as nullable
ALTER TABLE "public"."database_connections" ADD COLUMN "username" text;

-- Step 2: Backfill username for existing rows
UPDATE "public"."database_connections" SET "username" = agent_id WHERE "username" IS NULL;

-- Step 3: Set NOT NULL constraint
ALTER TABLE "public"."database_connections" ALTER COLUMN "username" SET NOT NULL;

-- +goose Down
-- reverse: modify "database_connections" table
ALTER TABLE "public"."database_connections" DROP COLUMN "username";
