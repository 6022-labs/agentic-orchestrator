-- +goose Up
-- modify "agent_context_messages" table
ALTER TABLE `agent_context_messages` MODIFY COLUMN `author` text NOT NULL;

-- +goose Down
-- reverse: modify "agent_context_messages" table
ALTER TABLE `agent_context_messages` MODIFY COLUMN `author` varchar(64) NOT NULL;
