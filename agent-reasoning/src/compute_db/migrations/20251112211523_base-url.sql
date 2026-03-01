-- +goose Up
-- modify "llm_connections" table
ALTER TABLE "public"."llm_connections"
ADD COLUMN "base_url" text NULL;

UPDATE "public"."llm_connections"
SET
    "base_url" = 'https://' || CONCAT_WS (':', "host", NULLIF(port, ''))
WHERE
    "host" IS NOT NULL
    AND "host" <> '';

ALTER TABLE "public"."llm_connections"
DROP COLUMN "host",
DROP COLUMN "port";

-- +goose Down
-- reverse: modify "llm_connections" table
ALTER TABLE "public"."llm_connections"
ADD COLUMN "port" text NULL,
ADD COLUMN "host" text NULL;

UPDATE "public"."llm_connections"
SET
    "host" = NULLIF(
        SPLIT_PART (
            SPLIT_PART (
                regexp_replace (COALESCE("base_url", ''), '^https?://', ''),
                '/',
                1
            ),
            ':',
            1
        ),
        ''
    ),
    "port" = NULLIF(
        SPLIT_PART (
            SPLIT_PART (
                regexp_replace (COALESCE("base_url", ''), '^https?://', ''),
                '/',
                1
            ),
            ':',
            2
        ),
        ''
    )
WHERE
    "base_url" IS NOT NULL
    AND "base_url" <> '';

ALTER TABLE "public"."llm_connections"
DROP COLUMN "base_url";