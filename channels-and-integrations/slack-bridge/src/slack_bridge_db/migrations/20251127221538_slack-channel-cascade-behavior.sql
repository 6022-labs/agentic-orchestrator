-- +goose Up
-- modify "slack_channel_agents" table
ALTER TABLE "public"."slack_channel_agents"
DROP CONSTRAINT "fk_slack_channel_agents_slack_agent",
DROP CONSTRAINT "fk_slack_channel_agents_slack_channel",
ADD CONSTRAINT "fk_slack_channel_agents_slack_agent" FOREIGN KEY (
    "slack_agent_collection_address",
    "slack_agent_collection_token_id"
) REFERENCES "public"."slack_agents" ("collection_address", "collection_token_id") ON UPDATE CASCADE ON DELETE CASCADE,
ADD CONSTRAINT "fk_slack_channel_agents_slack_channel" FOREIGN KEY ("slack_channel_id") REFERENCES "public"."slack_channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE;

-- modify "slack_threads" table
ALTER TABLE "public"."slack_threads" ADD CONSTRAINT "fk_slack_threads_channel" FOREIGN KEY ("channel_id") REFERENCES "public"."slack_channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE;

-- modify "slack_messages" table
ALTER TABLE "public"."slack_messages" ADD CONSTRAINT "fk_slack_messages_channel" FOREIGN KEY ("channel_id") REFERENCES "public"."slack_channels" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
ADD CONSTRAINT "fk_slack_messages_thread" FOREIGN KEY ("thread_id", "channel_id") REFERENCES "public"."slack_threads" ("id", "channel_id") ON UPDATE CASCADE ON DELETE CASCADE;

-- +goose Down
-- reverse: modify "slack_messages" table
ALTER TABLE "public"."slack_messages"
DROP CONSTRAINT "fk_slack_messages_thread",
DROP CONSTRAINT "fk_slack_messages_channel";

-- reverse: modify "slack_threads" table
ALTER TABLE "public"."slack_threads"
DROP CONSTRAINT "fk_slack_threads_channel";

-- reverse: modify "slack_channel_agents" table
ALTER TABLE "public"."slack_channel_agents"
DROP CONSTRAINT "fk_slack_channel_agents_slack_channel",
DROP CONSTRAINT "fk_slack_channel_agents_slack_agent",
ADD CONSTRAINT "fk_slack_channel_agents_slack_channel" FOREIGN KEY ("slack_channel_id") REFERENCES "public"."slack_channels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
ADD CONSTRAINT "fk_slack_channel_agents_slack_agent" FOREIGN KEY (
    "slack_agent_collection_address",
    "slack_agent_collection_token_id"
) REFERENCES "public"."slack_agents" ("collection_address", "collection_token_id") ON UPDATE NO ACTION ON DELETE NO ACTION;