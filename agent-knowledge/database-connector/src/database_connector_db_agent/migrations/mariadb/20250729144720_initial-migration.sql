-- +goose Up
-- create "agent_contexts" table
CREATE TABLE `agent_contexts` (
  `id` char(36) NOT NULL,
  `key` varchar(256) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_agent_contexts_key` (`key`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- create "agent_prompt_templates" table
CREATE TABLE `agent_prompt_templates` (
  `id` char(36) NOT NULL,
  `key` varchar(64) NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uni_agent_prompt_templates_key` (`key`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- create "agent_context_conversation_mappings" table
CREATE TABLE `agent_context_conversation_mappings` (
  `id` char(36) NOT NULL,
  `conversation_id` varchar(1024) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `agent_context_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_agent_context_conversation_mappings_agent_context_id` (`agent_context_id`),
  CONSTRAINT `fk_agent_contexts_conversation_mappings` FOREIGN KEY (`agent_context_id`) REFERENCES `agent_contexts` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- create "agent_context_messages" table
CREATE TABLE `agent_context_messages` (
  `id` varchar(36) NOT NULL,
  `content` text NOT NULL,
  `author` varchar(64) NOT NULL,
  `platform` varchar(64) NOT NULL,
  `role` varchar(64) NOT NULL,
  `hash` varchar(128) NOT NULL,
  `created_at` datetime(3) NOT NULL,
  `agent_context_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_agent_context_messages_agent_context_id` (`agent_context_id`),
  CONSTRAINT `fk_agent_contexts_messages` FOREIGN KEY (`agent_context_id`) REFERENCES `agent_contexts` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;

-- +goose Down
-- reverse: create "agent_context_messages" table
DROP TABLE `agent_context_messages`;
-- reverse: create "agent_context_conversation_mappings" table
DROP TABLE `agent_context_conversation_mappings`;
-- reverse: create "agent_prompt_templates" table
DROP TABLE `agent_prompt_templates`;
-- reverse: create "agent_contexts" table
DROP TABLE `agent_contexts`;
