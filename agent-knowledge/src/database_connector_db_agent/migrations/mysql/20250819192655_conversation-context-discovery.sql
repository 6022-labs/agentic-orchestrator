-- +goose Up
-- modify "agent_prompt_templates" table
UPDATE `agent_prompt_templates`
SET
    `key` = 'expert-context-found',
    `updated_at` = CURRENT_TIMESTAMP
WHERE
    `key` = 'expert-thread-found';

UPDATE `agent_prompt_templates`
SET
    `key` = 'expert-context-not-found',
    `updated_at` = CURRENT_TIMESTAMP
WHERE
    `key` = 'expert-thread-not-found';

UPDATE `agent_prompt_templates`
SET
    `content` = REPLACE (
        REPLACE (
            REPLACE (
                REPLACE (`content`, '{{.Keys}}', '{{.ContextKeys}}'),
                '{{ .Keys }}',
                '{{.ContextKeys}}'
            ),
            '{{ .Keys}}',
            '{{.ContextKeys}}'
        ),
        '{{.Keys }}',
        '{{.ContextKeys}}'
    ),
    `updated_at` = CURRENT_TIMESTAMP
WHERE
    `content` LIKE '%{{.Keys}}%'
    OR `content` LIKE '%{{ .Keys }}%'
    OR `content` LIKE '%{{ .Keys}}%'
    OR `content` LIKE '%{{.Keys }}%';

-- +goose Down
-- reverse: modify "agent_prompt_templates" table
UPDATE `agent_prompt_templates`
SET
    `key` = 'expert-thread-found',
    `updated_at` = CURRENT_TIMESTAMP
WHERE
    `key` = 'expert-context-found';

UPDATE `agent_prompt_templates`
SET
    `key` = 'expert-thread-not-found',
    `updated_at` = CURRENT_TIMESTAMP
WHERE
    `key` = 'expert-context-not-found';

UPDATE `agent_prompt_templates`
SET
    `content` = REPLACE (
        REPLACE (
            REPLACE (
                REPLACE (`content`, '{{.ContextKeys}}', '{{.Keys}}'),
                '{{ .ContextKeys }}',
                '{{.Keys}}'
            ),
            '{{ .ContextKeys}}',
            '{{.Keys}}'
        ),
        '{{.ContextKeys }}',
        '{{.Keys}}'
    ),
    `updated_at` = CURRENT_TIMESTAMP
WHERE
    `content` LIKE '%{{.ContextKeys}}%'
    OR `content` LIKE '%{{ .ContextKeys }}%'
    OR `content` LIKE '%{{ .ContextKeys}}%'
    OR `content` LIKE '%{{.ContextKeys }}%';