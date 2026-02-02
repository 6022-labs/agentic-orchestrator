package repositories

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_db"
)

type SlackThreadRepository struct {
	AbstractRepository[models.SlackThread]
}

func NewSlackThreadRepository(
	client *slack_bridge_db.DatabaseClient,
) *SlackThreadRepository {
	return &SlackThreadRepository{
		AbstractRepository: *NewAbstractRepository[models.SlackThread](client),
	}
}

func (r *SlackThreadRepository) GetByIdAndChannelId(id string, channelId string) (*models.SlackThread, error) {
	var thread models.SlackThread
	err := r.client.GetDB().
		First(&thread, "id = ? AND channel_id = ?", id, channelId).
		Error
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (r *SlackThreadRepository) GetProcessableThreadsForAgent(agent *models.SlackAgent, facilitatorTickInterval uint64) ([]*models.SlackThread, error) {
	query := r.client.GetDB().Model(&models.SlackThread{}).
		Joins("JOIN slack_channels ON slack_threads.channel_id = slack_channels.id").
		Joins("JOIN slack_channel_agents ON slack_channels.id = slack_channel_agents.slack_channel_id").
		Joins("JOIN slack_agents ON slack_channel_agents.slack_agent_collection_address = slack_agents.collection_address AND slack_channel_agents.slack_agent_collection_token_id = slack_agents.collection_token_id").
		Where("slack_agents.collection_address = ? AND slack_agents.collection_token_id = ? AND slack_threads.is_concluded = ?",
			agent.CollectionAddress,
			agent.CollectionTokenId,
			false,
		)

	if agent.Role == enums.AgentRoleFacilitator {
		// Check if last message from facilitator is older than threshold or facilitator has not sent any message yet
		facilitatorMessagesSubquery := r.client.GetDB().Model(&models.SlackMessage{}).
			Select("MAX(CAST(slack_messages.id AS double precision))").
			Where(`
				slack_messages.thread_id = slack_threads.id AND
				slack_messages.channel_id = slack_threads.channel_id AND
				slack_messages.author_collection_address = ? AND
				slack_messages.author_collection_token_id = ?`,
				agent.CollectionAddress,
				agent.CollectionTokenId,
			)

		query = query.Where(
			"COALESCE((?), 0) < (EXTRACT(EPOCH FROM NOW()) - ?)",
			facilitatorMessagesSubquery,
			facilitatorTickInterval,
		)
	} else if agent.Role == enums.AgentRoleHuman {
		// Check if human agent has been mentioned since his latest message or if physical human posted a message
		mentionSubquery := r.client.GetDB().Model(&models.SlackMessage{}).
			Select("1").
			Where("slack_messages.thread_id = slack_threads.id AND slack_messages.channel_id = slack_threads.channel_id").
			Where("slack_messages.content ILIKE ?", "%"+agent.EnsDomain+"%")

		humanMessagesSubquery := r.client.GetDB().Model(&models.SlackMessage{}).
			Select("1").
			Where("slack_messages.thread_id = slack_threads.id AND slack_messages.channel_id = slack_threads.channel_id").
			// Just need to check that the message is attributed to the human agent
			// because nobody apart from the human can post messages with this collection address + token id
			Where(
				"slack_messages.author_collection_address = ? AND slack_messages.author_collection_token_id = ?",
				agent.CollectionAddress,
				agent.CollectionTokenId,
			)

		query = query.Where("EXISTS (?) OR EXISTS (?)", mentionSubquery, humanMessagesSubquery)
	} else {
		// Check if agent has been mentioned since his latest message
		mentionSubquery := r.client.GetDB().Model(&models.SlackMessage{}).
			Select("1").
			Where("slack_messages.thread_id = slack_threads.id AND slack_messages.channel_id = slack_threads.channel_id").
			Where("slack_messages.content ILIKE ?", "%"+agent.EnsDomain+"%").
			Where(`
				CAST(slack_messages.id AS double precision) > COALESCE(
				(
					SELECT MAX(CAST(sm2.id AS double precision))
					FROM slack_messages sm2
					WHERE sm2.thread_id = slack_threads.id
						AND sm2.channel_id = slack_threads.channel_id
						AND sm2.author_collection_address = ?
						AND sm2.author_collection_token_id = ?
				), 0)`,
				agent.CollectionAddress,
				agent.CollectionTokenId,
			)

		query = query.Where("EXISTS (?)", mentionSubquery)
	}

	var ids []*string
	err := query.Distinct("slack_threads.id").Find(&ids).Error
	if err != nil {
		return nil, err
	}

	var threads []*models.SlackThread
	query = r.client.GetDB().Model(&models.SlackThread{}).Where("id IN ?", ids)
	err = query.Find(&threads).Error
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (r *SlackThreadRepository) ConcludeThreadsByChannelId(channelId string) error {
	return r.client.GetDB().
		Model(&models.SlackThread{}).
		Where("channel_id = ?", channelId).
		Update("is_concluded", true).
		Error
}
